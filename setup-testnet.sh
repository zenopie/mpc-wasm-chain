#!/bin/bash
# Multi-validator testnet setup for FROST testing
# FROST requires at least 2 participants for threshold signatures

set -e

CHAIN_ID="test-chain"
VAL1_HOME="/Users/zenopie/.wasmd-test"
VAL2_HOME="/Users/zenopie/.wasmd-val2"

echo "=== Setting up 2-validator FROST testnet ==="

# Cleanup
rm -rf $VAL1_HOME $VAL2_HOME

# Initialize both validators
./build/wasmd init validator1 --chain-id $CHAIN_ID --home $VAL1_HOME
./build/wasmd init validator2 --chain-id $CHAIN_ID --home $VAL2_HOME

# Create keys for both validators
./build/wasmd keys add validator1 --keyring-backend test --home $VAL1_HOME
./build/wasmd keys add validator2 --keyring-backend test --home $VAL2_HOME

# Get addresses
VAL1_ADDR=$(./build/wasmd keys show validator1 -a --keyring-backend test --home $VAL1_HOME)
VAL2_ADDR=$(./build/wasmd keys show validator2 -a --keyring-backend test --home $VAL2_HOME)

echo "Validator1 address: $VAL1_ADDR"
echo "Validator2 address: $VAL2_ADDR"

# Add genesis accounts to val1's genesis (we'll copy later)
./build/wasmd genesis add-genesis-account $VAL1_ADDR 100000000000stake --home $VAL1_HOME
./build/wasmd genesis add-genesis-account $VAL2_ADDR 100000000000stake --home $VAL1_HOME

# Create gentx for validator 1
./build/wasmd genesis gentx validator1 50000000000stake \
  --chain-id $CHAIN_ID \
  --keyring-backend test \
  --home $VAL1_HOME

# Copy genesis to val2 and create gentx there
cp $VAL1_HOME/config/genesis.json $VAL2_HOME/config/genesis.json
./build/wasmd genesis gentx validator2 50000000000stake \
  --chain-id $CHAIN_ID \
  --keyring-backend test \
  --home $VAL2_HOME

# Collect gentxs - copy val2's gentx to val1
cp $VAL2_HOME/config/gentx/*.json $VAL1_HOME/config/gentx/
./build/wasmd genesis collect-gentxs --home $VAL1_HOME

# Copy final genesis back to val2
cp $VAL1_HOME/config/genesis.json $VAL2_HOME/config/genesis.json

# Enable vote extensions at height 2
sed -i.bak 's/"vote_extensions_enable_height": "0"/"vote_extensions_enable_height": "2"/' $VAL1_HOME/config/genesis.json
sed -i.bak 's/"vote_extensions_enable_height": "0"/"vote_extensions_enable_height": "2"/' $VAL2_HOME/config/genesis.json

# Configure val1 ports (default)
# Val1: 26656 (p2p), 26657 (rpc), 1317 (api), 9090 (grpc)

# Configure val2 ports (shifted)
# Val2: 26666 (p2p), 26667 (rpc), 1318 (api), 9091 (grpc)
sed -i.bak 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/127.0.0.1:26667"/' $VAL2_HOME/config/config.toml
sed -i.bak 's/laddr = "tcp:\/\/0.0.0.0:26656"/laddr = "tcp:\/\/0.0.0.0:26666"/' $VAL2_HOME/config/config.toml
sed -i.bak 's/pprof_laddr = "localhost:6060"/pprof_laddr = "localhost:6061"/' $VAL2_HOME/config/config.toml
sed -i.bak 's/address = "tcp:\/\/localhost:1317"/address = "tcp:\/\/localhost:1318"/' $VAL2_HOME/config/app.toml
sed -i.bak 's/address = "localhost:9090"/address = "localhost:9091"/' $VAL2_HOME/config/app.toml
sed -i.bak 's/address = "localhost:9091"/address = "localhost:9092"/' $VAL2_HOME/config/app.toml

# Get node IDs
VAL1_NODE_ID=$(./build/wasmd comet show-node-id --home $VAL1_HOME)
VAL2_NODE_ID=$(./build/wasmd comet show-node-id --home $VAL2_HOME)

echo "Val1 Node ID: $VAL1_NODE_ID"
echo "Val2 Node ID: $VAL2_NODE_ID"

# Configure persistent peers
VAL1_PEER="$VAL1_NODE_ID@127.0.0.1:26656"
VAL2_PEER="$VAL2_NODE_ID@127.0.0.1:26666"

sed -i.bak "s/persistent_peers = \"\"/persistent_peers = \"$VAL2_PEER\"/" $VAL1_HOME/config/config.toml
sed -i.bak "s/persistent_peers = \"\"/persistent_peers = \"$VAL1_PEER\"/" $VAL2_HOME/config/config.toml

# Enable API on both
sed -i.bak 's/enable = false/enable = true/' $VAL1_HOME/config/app.toml
sed -i.bak 's/enable = false/enable = true/' $VAL2_HOME/config/app.toml

# Set minimum gas prices
sed -i.bak 's/minimum-gas-prices = ""/minimum-gas-prices = "0stake"/' $VAL1_HOME/config/app.toml
sed -i.bak 's/minimum-gas-prices = ""/minimum-gas-prices = "0stake"/' $VAL2_HOME/config/app.toml

echo ""
echo "=== Setup Complete ==="
echo ""
echo "To start the testnet, run in separate terminals:"
echo ""
echo "Terminal 1 (Validator 1):"
echo "  ./build/wasmd start --home $VAL1_HOME"
echo ""
echo "Terminal 2 (Validator 2):"
echo "  ./build/wasmd start --home $VAL2_HOME"
echo ""
echo "Or run start-validators.sh to start both in background"
