#!/bin/bash
# Initialize testnet for Docker containers

set -e

NUM_VALIDATORS=2
CHAIN_ID="docker-testnet"
OUTPUT_DIR="./.docker-testnet"

echo "=== Docker Testnet Initialization ==="

# Clean up
rm -rf $OUTPUT_DIR

# Initialize validators
echo "Initializing $NUM_VALIDATORS validators..."
./build/wasmd testnet init-files \
  --v $NUM_VALIDATORS \
  --chain-id $CHAIN_ID \
  --output-dir $OUTPUT_DIR \
  --node-daemon-home "" \
  --keyring-backend test \
  --commit-timeout 2s 2>&1 | grep -v "TSS DEBUG\|TSS ERROR"

# Patch genesis for vote extensions
echo "Patching genesis for vote extensions..."
for i in $(seq 0 $((NUM_VALIDATORS - 1))); do
    GENESIS="$OUTPUT_DIR/node$i/config/genesis.json"
    jq '.consensus.params.abci.vote_extensions_enable_height = "2"' "$GENESIS" > /tmp/genesis_$i.json
    mv /tmp/genesis_$i.json "$GENESIS"
done

# Update peer addresses to use Docker network IPs
echo "Updating peer addresses for Docker network..."

# Get node IDs (filter out TSS debug messages)
NODE0_ID=$(./build/wasmd tendermint show-node-id --home $OUTPUT_DIR/node0 2>/dev/null | grep -v "TSS" | tail -1)
NODE1_ID=$(./build/wasmd tendermint show-node-id --home $OUTPUT_DIR/node1 2>/dev/null | grep -v "TSS" | tail -1)

echo "Node0 ID: $NODE0_ID"
echo "Node1 ID: $NODE1_ID"

# Update config.toml for each node to use Docker IPs (macOS compatible sed)
# Node0 connects to Node1
sed -i '' "s|persistent_peers = .*|persistent_peers = \"${NODE1_ID}@172.20.0.11:26656\"|" $OUTPUT_DIR/node0/config/config.toml
# Node1 connects to Node0
sed -i '' "s|persistent_peers = .*|persistent_peers = \"${NODE0_ID}@172.20.0.10:26656\"|" $OUTPUT_DIR/node1/config/config.toml

# Enable RPC from any host
for i in 0 1; do
    sed -i '' 's|laddr = "tcp://127.0.0.1:26657"|laddr = "tcp://0.0.0.0:26657"|' $OUTPUT_DIR/node$i/config/config.toml
    # Enable API
    sed -i '' 's|enable = false|enable = true|' $OUTPUT_DIR/node$i/config/app.toml
    sed -i '' 's|address = "tcp://localhost:1317"|address = "tcp://0.0.0.0:1317"|' $OUTPUT_DIR/node$i/config/app.toml
done

echo ""
echo "Vote extensions config:"
cat $OUTPUT_DIR/node0/config/genesis.json | jq '.consensus.params.abci'

echo ""
echo "=== Initialization complete ==="
echo ""
echo "To start the testnet:"
echo "  docker-compose up --build"
echo ""
echo "To query TSS state:"
echo "  ./build/wasmd query tss dkg-state --node http://localhost:26657"
