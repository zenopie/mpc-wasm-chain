#!/bin/bash
# Start multi-validator testnet with vote extensions in a single terminal
# All validators run in background, logs are aggregated

set -e

NUM_VALIDATORS=${1:-3}
CHAIN_ID="test-chain"
OUTPUT_DIR="./.testnets"

echo "=== MPC Chain Multi-Validator Testnet ==="
echo "Validators: $NUM_VALIDATORS"
echo ""

# Kill any existing
pkill -f "wasmd start" 2>/dev/null || true
sleep 1

# Clean up
rm -rf $OUTPUT_DIR

# Initialize with vote extensions enabled
echo "Initializing $NUM_VALIDATORS validators..."
./build/wasmd testnet init-files \
  --v $NUM_VALIDATORS \
  --chain-id $CHAIN_ID \
  --output-dir $OUTPUT_DIR \
  --node-daemon-home simd \
  --keyring-backend test \
  --single-host \
  --commit-timeout 5s 2>&1 | grep -v "TSS DEBUG\|TSS ERROR"

# Patch genesis files to enable vote extensions at height 2
echo "Patching genesis for vote extensions..."
for i in $(seq 0 $((NUM_VALIDATORS - 1))); do
    GENESIS="$OUTPUT_DIR/node$i/simd/config/genesis.json"
    jq '.consensus.params.abci.vote_extensions_enable_height = "2"' "$GENESIS" > /tmp/genesis_$i.json
    mv /tmp/genesis_$i.json "$GENESIS"
done

# Fix peer addresses to use localhost instead of external IPs
echo "Fixing peer addresses for localhost..."
for i in $(seq 0 $((NUM_VALIDATORS - 1))); do
    CONFIG="$OUTPUT_DIR/node$i/simd/config/config.toml"
    # Replace any 192.168.x.x addresses with 127.0.0.1
    sed -i '' 's/192\.168\.[0-9]*\.[0-9]*/127.0.0.1/g' "$CONFIG"
done

echo ""
echo "Vote extensions config:"
cat $OUTPUT_DIR/node0/simd/config/genesis.json | jq '.consensus.params.abci'

echo ""
echo "Starting validators..."

# Create a named pipe for aggregated output
FIFO=$(mktemp -u)
mkfifo $FIFO

# Function to prefix output
prefix_output() {
    local prefix=$1
    while IFS= read -r line; do
        echo "[$prefix] $line"
    done
}

# Start validators
PIDS=""
for i in $(seq 0 $((NUM_VALIDATORS - 1))); do
    NODE_DIR="$OUTPUT_DIR/node$i/simd"
    echo "Starting validator $i..."
    ./build/wasmd start --home $NODE_DIR 2>&1 | prefix_output "VAL$i" &
    PIDS="$PIDS $!"
    sleep 2
done

echo ""
echo "=== All validators started ==="
echo "PIDs:$PIDS"
echo ""
echo "RPC: http://localhost:26657"
echo "API: http://localhost:1317"
echo "GRPC: localhost:9090"
echo ""
echo "Press Ctrl+C to stop all validators"
echo ""

# Cleanup function
cleanup() {
    echo ""
    echo "Stopping validators..."
    for pid in $PIDS; do
        kill $pid 2>/dev/null || true
    done
    rm -f $FIFO
    exit 0
}

trap cleanup SIGINT SIGTERM

# Wait for all background processes
wait
