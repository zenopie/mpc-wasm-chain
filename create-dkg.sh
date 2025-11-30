#!/bin/bash

# Script to create a 2/3 DKG (2-of-3 threshold signature scheme)
# Prerequisites: Chain must be running and producing blocks

BINARY="./build/wasmd"
NODE="tcp://localhost:26657"
CHAIN_ID="chain-FR1eOT"  # Update this to match your chain
HOME_DIR="./.testnets/node0/wasmd"
KEYRING="test"

echo "========================================="
echo "Creating 2/3 DKG (Distributed Key Generation)"
echo "========================================="
echo ""

# Check if chain is running
echo "1. Checking if chain is running..."
if ! curl -s $NODE/status > /dev/null 2>&1; then
    echo "   ✗ Chain is NOT running on $NODE"
    echo "   Please start the chain first with: ./start-validators.sh"
    exit 1
fi

HEIGHT=$(curl -s $NODE/status | jq -r '.result.sync_info.latest_block_height')
echo "   ✓ Chain is running at block $HEIGHT"
echo ""

# Step 1: Create a KeySet with threshold 2 and max_signers 3
echo "2. Creating KeySet with threshold=2, max_signers=3..."
echo "   This defines a 2-of-3 threshold signature scheme"
echo ""

# Note: You'll need to have a funded account in the keyring
# The testnet should have created validator accounts automatically

$BINARY tx tss create-key-set \
  --threshold 2 \
  --max-signers 3 \
  --description "2-of-3 DKG for testing" \
  --from validator \
  --chain-id $CHAIN_ID \
  --home $HOME_DIR \
  --keyring-backend $KEYRING \
  --node $NODE \
  --gas auto \
  --gas-adjustment 1.5 \
  --fees 1000stake \
  --yes

if [ $? -eq 0 ]; then
    echo ""
    echo "   ✓ KeySet creation transaction submitted"
    echo ""
    echo "3. Waiting for transaction to be included in a block..."
    sleep 6

    # Query all key sets to get the ID
    echo ""
    echo "4. Querying created KeySets..."
    $BINARY query tss all-key-sets --node $NODE

    echo ""
    echo "========================================="
    echo "Next Steps:"
    echo "========================================="
    echo "1. Get the key_set_id from the output above"
    echo "2. Initiate DKG with:"
    echo "   $BINARY tx tss initiate-dkg \\"
    echo "     --key-set-id <KEY_SET_ID> \\"
    echo "     --from validator \\"
    echo "     --chain-id $CHAIN_ID \\"
    echo "     --node $NODE \\"
    echo "     --yes"
    echo ""
    echo "3. Monitor logs for DKG progress:"
    echo "   tail -f node0.log | grep -E '(DKG|TSS|vote)'"
    echo ""
else
    echo ""
    echo "   ✗ Failed to create KeySet"
    echo "   Check that you have a funded 'validator' account"
    exit 1
fi
