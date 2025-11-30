# TSS Module WASM Integration

The TSS module is now fully integrated with CosmWasm, allowing smart contracts to create threshold signature key sets and request signatures.

## Overview

Smart contracts interact with the TSS module by:
1. **Creating KeySets** - Initiates a DKG ceremony among validators
2. **Requesting Signatures** - Requests a threshold signature for a message

Validators automatically handle the DKG rounds and signing protocol through the EndBlocker. Contracts don't need to manage these low-level operations.

## End-to-End Flow

### Creating a KeySet

**Block N:** Contract initiates DKG
- Contract calls `create_key_set` custom message
- `x/tss` creates a `DKGSession` and `KeySet` with status `PENDING_DKG`
- Event emitted: `KeySetCreated`

**Block N+1:** DKG Round 1
- `x/tss` EndBlocker detects new DKG session
- Validators submit `MsgSubmitDKGRound1` with commitments
- `x/tss` stores commitments in state

**Block N+2:** DKG Round 2
- `x/tss` EndBlocker detects threshold commitments met
- Validators submit `MsgSubmitDKGRound2` with shares
- `x/tss` stores shares in state

**Block N+3:** DKG Complete
- `x/tss` EndBlocker detects threshold shares met
- EndBlocker aggregates shares to compute group public key
- `KeySet` status updated to `ACTIVE`
- Event emitted: `KeySetActive`

### Requesting a Signature

**Block M:** Contract initiates signing
- Contract calls `request_signature` custom message
- `x/tss` creates `SigningRequest` and `SigningSession`
- Event emitted: `SignatureRequested`

**Block M+1:** Signing Round 1 - Commitments
- `x/tss` EndBlocker detects new signing request
- Validators submit `MsgSubmitCommitment` (nonce commitments)
- `x/tss` stores commitments in state

**Block M+2:** Signing Round 2 - Shares
- `x/tss` EndBlocker detects threshold commitments met
- Validators submit `MsgSubmitSignatureShare`
- `x/tss` stores signature shares in state

**Block M+3:** Signature Complete & Callback
- `x/tss` EndBlocker detects threshold shares met
- EndBlocker aggregates shares into final signature
- Signature verified against group public key
- `SigningRequest` updated with final signature
- **Sudo callback executed** to original contract with signature
- Contract processes signature

## Custom Messages

Smart contracts can send these custom messages:

### 1. Create Key Set

Initiates a DKG ceremony to create a new threshold signature key set:

```json
{
  "custom": {
    "create_key_set": {
      "threshold": 2,
      "max_signers": 3,
      "description": "My TSS key set",
      "timeout_blocks": 100
    }
  }
}
```

**Flow:**
- Creates `KeySet` with status `PENDING_DKG`
- Creates `DKGSession`
- Validators automatically participate in DKG rounds
- After completion, `KeySet` becomes `ACTIVE`

### 2. Request Signature

Requests a threshold signature for a message hash:

```json
{
  "custom": {
    "request_signature": {
      "key_set_id": "keyset-123",
      "message_hash": "0x1234567890abcdef...",
      "callback": "optional-callback-data"
    }
  }
}
```

**Flow:**
- Creates `SigningRequest` with status `PENDING`
- Creates `SigningSession`
- Validators automatically participate in signing rounds
- After completion, contract receives sudo callback with signature

## Custom Queries

Smart contracts can query the TSS module state using custom queries:

### 1. Query Key Set

Get details about a specific key set:

```json
{
  "custom": {
    "key_set": {
      "id": "keyset-123"
    }
  }
}
```

Response:
```json
{
  "id": "keyset-123",
  "owner": "wasm1...",
  "threshold": 2,
  "max_signers": 3,
  "participants": ["validator1...", "validator2...", "validator3..."],
  "group_pubkey": "0x...",
  "status": "KEY_SET_STATUS_ACTIVE",
  "description": "My TSS key set",
  "created_height": 12345
}
```

### 2. Query Signing Request

Get details about a signing request:

```json
{
  "custom": {
    "signing_request": {
      "id": "signing-request-123"
    }
  }
}
```

Response:
```json
{
  "id": "signing-request-123",
  "key_set_id": "keyset-123",
  "requester": "wasm1...",
  "message_hash": "0x1234567890abcdef...",
  "callback": "optional-callback-data",
  "status": "SIGNING_REQUEST_STATUS_COMPLETE",
  "signature": "0xabcdef...",
  "created_height": 12350
}
```

### 3. Query DKG Session

Get details about a DKG session:

```json
{
  "custom": {
    "dkg_session": {
      "id": "dkg-session-123"
    }
  }
}
```

Response:
```json
{
  "id": "dkg-session-123",
  "key_set_id": "keyset-123",
  "state": "DKG_STATE_COMPLETE",
  "threshold": 2,
  "max_signers": 3,
  "participants": ["validator1...", "validator2...", "validator3..."],
  "start_height": 12340,
  "timeout_height": 12440
}
```

### 4. Query Signing Session

Get details about a signing session:

```json
{
  "custom": {
    "signing_session": {
      "request_id": "signing-request-123"
    }
  }
}
```

Response:
```json
{
  "request_id": "signing-request-123",
  "key_set_id": "keyset-123",
  "threshold": 2,
  "participants": ["validator1...", "validator2...", "validator3..."],
  "state": "SIGNING_STATE_COMPLETE",
  "start_height": 12350,
  "timeout_height": 12450
}
```

## Example: Bitcoin Signer Contract

A complete working example is available in the [examples/bitcoin-signer](./examples/bitcoin-signer/) directory.

This example demonstrates:
- Creating a threshold signature key set
- Requesting signatures for Bitcoin transaction hashes
- Receiving signature callbacks via sudo
- Querying key set and signing request status

### Quick Start

See the [Bitcoin Signer README](./examples/bitcoin-signer/README.md) for:
- Detailed usage instructions
- Building and testing guide
- Integration patterns

### Key Components

- **[tss.rs](./examples/bitcoin-signer/tss.rs)** - TSS module integration types
- **[contract.rs](./examples/bitcoin-signer/contract.rs)** - Main contract logic
- **[msg.rs](./examples/bitcoin-signer/msg.rs)** - Message definitions
- **[state.rs](./examples/bitcoin-signer/state.rs)** - State management

### Basic Usage

```rust
// Initialize key set (admin only)
ExecuteMsg::InitializeKeySet {
    threshold: 2,
    max_signers: 3,
}

// Request signature (anyone)
ExecuteMsg::SignBitcoinTx {
    tx_hash: binary_hash,
}

// Receive signature (automatic sudo callback)
SudoMsg::SignatureComplete {
    request_id: "...",
    signature: binary_signature,
}
```

## Implementation Details

The WASM integration is implemented in [x/tss/wasm.go](./wasm.go) and includes:

1. **Custom Message Encoder**: Converts JSON messages from WASM contracts into native TSS module messages
2. **Custom Query Plugin**: Handles queries from WASM contracts to the TSS module state

These are automatically registered with the WasmKeeper when the application starts (see [app/app.go](../../app/app.go)).

## Use Cases

Some potential use cases for the TSS module in smart contracts:

1. **Multi-Signature Wallets**: Create threshold signature wallets controlled by smart contract logic
2. **Cross-Chain Bridges**: Generate threshold signatures for cross-chain asset transfers
3. **Decentralized Custody**: Implement decentralized custody solutions with threshold signatures
4. **DAO Treasury Management**: Enable DAOs to manage treasuries with threshold signature security
5. **Oracle Networks**: Sign oracle data with threshold signatures for increased security
