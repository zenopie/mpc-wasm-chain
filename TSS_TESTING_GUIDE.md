# TSS Testing & Debugging Guide

## Overview

Your MPC DKG and threshold signature chain uses **vote extensions** to propagate TSS data between validators. This guide explains how to test and debug the system.

## Architecture Summary

### Components

1. **Vote Extensions** ([x/tss/abci/vote_extensions.go](x/tss/abci/vote_extensions.go))
   - `ExtendVote`: Validators include DKG/signing data in their votes
   - `VerifyVoteExtension`: Validators verify each other's vote extensions

2. **Proposal Handlers** ([x/tss/abci/proposals.go](x/tss/abci/proposals.go))
   - `PrepareProposal`: Block proposer aggregates TSS data from vote extensions
   - `ProcessProposal`: Validators verify the aggregated data

3. **ABCI Hooks** ([x/tss/module/module.go](x/tss/module/module.go))
   - `BeginBlock`: Processes aggregated TSS data and stores it in state
   - `EndBlock`: Advances DKG/signing state machines

4. **TSS Logic** ([x/tss/keeper/](x/tss/keeper/))
   - DKG ceremony management (Round 1 & 2)
   - Threshold signature creation (Round 1 & 2)
   - Uses Binance tss-lib for cryptography

## Quick Start

### 1. Build the Chain

```bash
make build
```

**Note on Password Prompts**: If you see password prompts, use `--keyring-backend test` flag:
```bash
./build/wasmd start --home ./.testnets/node0/wasmd --keyring-backend test
```

### 2. Start the Testnet

```bash
./start-validators.sh
```

This starts 3 validator nodes. Check logs:
```bash
tail -f node0.log
```

### 3. Run Diagnostics

```bash
./diagnose-tss.sh
```

This checks:
- ✅ Chain status
- ✅ Vote extension activity
- ✅ TSS data flow
- ✅ Error conditions

### 4. Test DKG Flow

```bash
./test-tss.sh
```

This will:
1. Create a KeySet
2. Initiate DKG
3. Monitor progress through vote extensions
4. Check final state

## Understanding Vote Extensions

### The Flow

```
Block N:
  1. ExtendVote (each validator)
     - Check for active DKG/signing sessions
     - Generate TSS data (commitments, shares, etc.)
     - Include in vote extension

Block N+1:
  2. PrepareProposal (block proposer)
     - Collect vote extensions from Block N
     - Aggregate TSS data by session/request
     - Include in block proposal

  3. ProcessProposal (all validators)
     - Verify aggregated data matches vote extensions
     - Store in memory for BeginBlock

  4. BeginBlock (all validators)
     - Retrieve aggregated TSS data
     - Store in blockchain state
     - Process DKG/signing submissions

  5. EndBlock (all validators)
     - Check if enough data collected
     - Advance state machines (Round 1 → Round 2 → Complete)
```

### What to Look For in Logs

**Vote Extensions Working:**
```
ExtendVote called height=100
Extended vote with TSS data dkg_r1=3 dkg_r2=0 commitments=0 shares=0
```

**Proposal Aggregation:**
```
Prepared proposal with aggregated TSS data dkg_r1_sessions=1 dkg_r2_sessions=0
```

**BeginBlock Processing:**
```
Processing aggregated TSS data from vote extensions height=101 dkg_r1_sessions=1
```

**EndBlock State Transitions:**
```
DKG session advanced from ROUND1 to ROUND2
DKG completed successfully session=dkg-test-keyset-1-100
```

## Manual Testing

### Create a KeySet

```bash
./build/wasmd tx tss create-keyset \
  "my-keyset" \
  2 \
  3 \
  100 \
  --from node0 \
  --chain-id testing \
  --keyring-backend test \
  --home ./.testnets/node0/wasmd \
  --node tcp://localhost:26657 \
  --yes
```

Parameters:
- `"my-keyset"`: Description
- `2`: Threshold (minimum signatures needed)
- `3`: Max signers
- `100`: Timeout in blocks

### Query KeySet Status

```bash
./build/wasmd query tss keyset <keyset-id> \
  --node tcp://localhost:26657 \
  --output json | jq
```

### Request a Signature

```bash
# First get the keyset ID from the creation response
KEYSET_ID="keyset-1234"
MESSAGE_HASH="0x1234567890abcdef..."

./build/wasmd tx tss request-signature \
  $KEYSET_ID \
  $MESSAGE_HASH \
  "" \
  --from node0 \
  --chain-id testing \
  --keyring-backend test \
  --home ./.testnets/node0/wasmd \
  --node tcp://localhost:26657 \
  --yes
```

### Monitor Logs in Real-Time

```bash
# All TSS activity
tail -f node0.log | grep -E "(TSS|DKG|Signing|vote|proposal)"

# Just vote extensions
tail -f node0.log | grep -E "(ExtendVote|VerifyVoteExtension)"

# Just state transitions
tail -f node0.log | grep -E "(BeginBlock|EndBlock)"
```

## Common Issues & Solutions

### Issue: Vote Extensions Not Being Called

**Symptoms:**
- No "ExtendVote called" in logs
- Vote extension count is 0

**Solutions:**
1. Check that vote extensions are enabled in [app/app.go:988-1002](app/app.go#L988-L1002)
2. Verify CometBFT version supports vote extensions (v0.38+)
3. Check consensus params enable vote extensions

### Issue: Vote Extensions Are Empty

**Symptoms:**
- "ExtendVote called" appears but no TSS data
- `dkg_r1=0 dkg_r2=0` in logs

**This is NORMAL** if:
- No DKG sessions are active
- No signing requests are pending

**To trigger TSS data:**
- Create a KeySet (initiates DKG)
- Request a signature (initiates signing)

### Issue: DKG Not Progressing

**Symptoms:**
- DKG session stuck in ROUND1
- Not enough submissions

**Debug:**
```bash
# Check if validators are generating data
grep "GenerateDKGRound1Data" node*.log

# Check if data is being processed
grep "ProcessDKGRound1" node*.log

# Check submission counts
grep "DKGRound1Count" node*.log
```

**Common causes:**
- Validator addresses mismatch (consensus vs account address)
- Not enough validators participating
- Timeout too short

### Issue: Build Errors

**`RawContractMessage` not found:**

The x/wasm types were accidentally deleted. Don't restore them manually - they're protobuf generated files. Instead:

1. Check git status: `git status x/wasm/types/`
2. The `.pb.go` files should be modified, not deleted
3. If you deleted `.go` files (not `.pb.go`), restore them:
   ```bash
   git checkout HEAD -- x/wasm/types/types.go x/wasm/types/codec.go
   ```

## Testing Checklist

- [ ] Chain starts successfully
- [ ] Vote extensions are called (`ExtendVote` in logs)
- [ ] Vote extensions are verified (`VerifyVoteExtension` in logs)
- [ ] Proposals aggregate TSS data (`PrepareProposal` in logs)
- [ ] BeginBlock processes TSS data
- [ ] Create KeySet transaction succeeds
- [ ] DKG Round 1 data appears in vote extensions
- [ ] DKG transitions to Round 2
- [ ] DKG completes successfully
- [ ] KeySet becomes ACTIVE
- [ ] Request signature transaction succeeds
- [ ] Signing Round 1 commitments appear in vote extensions
- [ ] Signing transitions to Round 2
- [ ] Signature is aggregated successfully

## Key Files Reference

| File | Purpose |
|------|---------|
| [x/tss/abci/vote_extensions.go](x/tss/abci/vote_extensions.go) | Vote extension generation & verification |
| [x/tss/abci/proposals.go](x/tss/abci/proposals.go) | Proposal aggregation logic |
| [x/tss/keeper/block_tss_data.go](x/tss/keeper/block_tss_data.go) | BeginBlock TSS processing |
| [x/tss/keeper/dkg.go](x/tss/keeper/dkg.go) | DKG state machine |
| [x/tss/keeper/signing.go](x/tss/keeper/signing.go) | Signing state machine |
| [x/tss/keeper/frost.go](x/tss/keeper/frost.go) | TSS cryptography (tss-lib) |
| [x/tss/module/module.go](x/tss/module/module.go) | ABCI hooks (BeginBlock/EndBlock) |
| [app/app.go](app/app.go) | Vote extension handler setup |

## Next Steps

1. **Run diagnostics** to check current state: `./diagnose-tss.sh`
2. **Start chain** if not running: `./start-validators.sh`
3. **Run tests** to verify functionality: `./test-tss.sh`
4. **Monitor logs** for any issues: `tail -f node0.log`
5. **Create KeySet** to trigger DKG flow
6. **Request signature** to test signing flow

## Security Notes

⚠️ **Current Implementation:**
- Uses **placeholder cryptography** (not production-ready)
- **Real secret shares** should NEVER be stored on-chain
- Vote extensions are **not validated** cryptographically (TODO in code)

For production:
- Implement real tss-lib integration
- Add cryptographic verification of commitments/shares
- Store key shares in secure off-chain storage (HSM, etc.)
- Add slashing for malicious behavior
- Implement proper validator address mapping

## Need Help?

Check the logs for specific error messages and search for the error in the codebase. Most functions have detailed comments explaining their purpose and limitations.

Common log locations:
- `node0.log`, `node1.log`, `node2.log` - Validator logs
- `chain.log` - Chain startup log
- `testnet.log` - Testnet initialization log
