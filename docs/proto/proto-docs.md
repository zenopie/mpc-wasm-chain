<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [mpcchain/tss/v1/types.proto](#mpcchain/tss/v1/types.proto)
    - [DKGRound1Data](#mpcchain.tss.v1.DKGRound1Data)
    - [DKGRound2Data](#mpcchain.tss.v1.DKGRound2Data)
    - [DKGSession](#mpcchain.tss.v1.DKGSession)
    - [KeySet](#mpcchain.tss.v1.KeySet)
    - [KeyShare](#mpcchain.tss.v1.KeyShare)
    - [Params](#mpcchain.tss.v1.Params)
    - [SignatureShare](#mpcchain.tss.v1.SignatureShare)
    - [SigningCommitment](#mpcchain.tss.v1.SigningCommitment)
    - [SigningRequest](#mpcchain.tss.v1.SigningRequest)
    - [SigningSession](#mpcchain.tss.v1.SigningSession)
  
    - [DKGState](#mpcchain.tss.v1.DKGState)
    - [KeySetStatus](#mpcchain.tss.v1.KeySetStatus)
    - [SigningRequestStatus](#mpcchain.tss.v1.SigningRequestStatus)
    - [SigningState](#mpcchain.tss.v1.SigningState)
  
- [mpcchain/tss/v1/genesis.proto](#mpcchain/tss/v1/genesis.proto)
    - [GenesisState](#mpcchain.tss.v1.GenesisState)
  
- [mpcchain/tss/v1/module.proto](#mpcchain/tss/v1/module.proto)
    - [Module](#mpcchain.tss.v1.Module)
  
- [mpcchain/tss/v1/query.proto](#mpcchain/tss/v1/query.proto)
    - [QueryAllDKGSessionsRequest](#mpcchain.tss.v1.QueryAllDKGSessionsRequest)
    - [QueryAllDKGSessionsResponse](#mpcchain.tss.v1.QueryAllDKGSessionsResponse)
    - [QueryAllKeySetsRequest](#mpcchain.tss.v1.QueryAllKeySetsRequest)
    - [QueryAllKeySetsResponse](#mpcchain.tss.v1.QueryAllKeySetsResponse)
    - [QueryAllSigningRequestsRequest](#mpcchain.tss.v1.QueryAllSigningRequestsRequest)
    - [QueryAllSigningRequestsResponse](#mpcchain.tss.v1.QueryAllSigningRequestsResponse)
    - [QueryDKGSessionRequest](#mpcchain.tss.v1.QueryDKGSessionRequest)
    - [QueryDKGSessionResponse](#mpcchain.tss.v1.QueryDKGSessionResponse)
    - [QueryKeySetRequest](#mpcchain.tss.v1.QueryKeySetRequest)
    - [QueryKeySetResponse](#mpcchain.tss.v1.QueryKeySetResponse)
    - [QueryParamsRequest](#mpcchain.tss.v1.QueryParamsRequest)
    - [QueryParamsResponse](#mpcchain.tss.v1.QueryParamsResponse)
    - [QuerySigningRequestRequest](#mpcchain.tss.v1.QuerySigningRequestRequest)
    - [QuerySigningRequestResponse](#mpcchain.tss.v1.QuerySigningRequestResponse)
  
    - [Query](#mpcchain.tss.v1.Query)
  
- [mpcchain/tss/v1/tx.proto](#mpcchain/tss/v1/tx.proto)
    - [MsgCreateKeySet](#mpcchain.tss.v1.MsgCreateKeySet)
    - [MsgCreateKeySetResponse](#mpcchain.tss.v1.MsgCreateKeySetResponse)
    - [MsgInitiateDKG](#mpcchain.tss.v1.MsgInitiateDKG)
    - [MsgInitiateDKGResponse](#mpcchain.tss.v1.MsgInitiateDKGResponse)
    - [MsgRequestSignature](#mpcchain.tss.v1.MsgRequestSignature)
    - [MsgRequestSignatureResponse](#mpcchain.tss.v1.MsgRequestSignatureResponse)
    - [MsgSubmitCommitment](#mpcchain.tss.v1.MsgSubmitCommitment)
    - [MsgSubmitCommitmentResponse](#mpcchain.tss.v1.MsgSubmitCommitmentResponse)
    - [MsgSubmitDKGRound1](#mpcchain.tss.v1.MsgSubmitDKGRound1)
    - [MsgSubmitDKGRound1Response](#mpcchain.tss.v1.MsgSubmitDKGRound1Response)
    - [MsgSubmitDKGRound2](#mpcchain.tss.v1.MsgSubmitDKGRound2)
    - [MsgSubmitDKGRound2Response](#mpcchain.tss.v1.MsgSubmitDKGRound2Response)
    - [MsgSubmitSignatureShare](#mpcchain.tss.v1.MsgSubmitSignatureShare)
    - [MsgSubmitSignatureShareResponse](#mpcchain.tss.v1.MsgSubmitSignatureShareResponse)
    - [MsgUpdateParams](#mpcchain.tss.v1.MsgUpdateParams)
    - [MsgUpdateParamsResponse](#mpcchain.tss.v1.MsgUpdateParamsResponse)
  
    - [Msg](#mpcchain.tss.v1.Msg)
  
- [Scalar Value Types](#scalar-value-types)



<a name="mpcchain/tss/v1/types.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## mpcchain/tss/v1/types.proto



<a name="mpcchain.tss.v1.DKGRound1Data"></a>

### DKGRound1Data



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  |  |
| `commitment` | [bytes](#bytes) |  |  |
| `submitted_height` | [int64](#int64) |  |  |






<a name="mpcchain.tss.v1.DKGRound2Data"></a>

### DKGRound2Data



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  |  |
| `share` | [bytes](#bytes) |  |  |
| `submitted_height` | [int64](#int64) |  |  |






<a name="mpcchain.tss.v1.DKGSession"></a>

### DKGSession
DKG Session and Round data


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `key_set_id` | [string](#string) |  |  |
| `state` | [DKGState](#mpcchain.tss.v1.DKGState) |  |  |
| `threshold` | [uint32](#uint32) |  |  |
| `max_signers` | [uint32](#uint32) |  |  |
| `participants` | [string](#string) | repeated |  |
| `start_height` | [int64](#int64) |  |  |
| `timeout_height` | [int64](#int64) |  |  |






<a name="mpcchain.tss.v1.KeySet"></a>

### KeySet
KeySet represents a threshold signature key set


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `owner` | [string](#string) |  |  |
| `threshold` | [uint32](#uint32) |  |  |
| `max_signers` | [uint32](#uint32) |  |  |
| `participants` | [string](#string) | repeated |  |
| `group_pubkey` | [bytes](#bytes) |  |  |
| `status` | [KeySetStatus](#mpcchain.tss.v1.KeySetStatus) |  |  |
| `description` | [string](#string) |  |  |
| `created_height` | [int64](#int64) |  |  |






<a name="mpcchain.tss.v1.KeyShare"></a>

### KeyShare
KeyShare represents a validator's share of a threshold key


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key_set_id` | [string](#string) |  |  |
| `validator_address` | [string](#string) |  |  |
| `share_data` | [bytes](#bytes) |  |  |
| `group_pubkey` | [bytes](#bytes) |  |  |
| `created_height` | [int64](#int64) |  |  |






<a name="mpcchain.tss.v1.Params"></a>

### Params
Params defines the module parameters






<a name="mpcchain.tss.v1.SignatureShare"></a>

### SignatureShare



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  |  |
| `share` | [bytes](#bytes) |  |  |
| `submitted_height` | [int64](#int64) |  |  |






<a name="mpcchain.tss.v1.SigningCommitment"></a>

### SigningCommitment



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  |  |
| `commitment` | [bytes](#bytes) |  |  |
| `submitted_height` | [int64](#int64) |  |  |






<a name="mpcchain.tss.v1.SigningRequest"></a>

### SigningRequest
Signing Request and Session data


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `key_set_id` | [string](#string) |  |  |
| `requester` | [string](#string) |  |  |
| `message_hash` | [bytes](#bytes) |  |  |
| `callback` | [string](#string) |  |  |
| `status` | [SigningRequestStatus](#mpcchain.tss.v1.SigningRequestStatus) |  |  |
| `signature` | [bytes](#bytes) |  |  |
| `created_height` | [int64](#int64) |  |  |






<a name="mpcchain.tss.v1.SigningSession"></a>

### SigningSession



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `request_id` | [string](#string) |  |  |
| `key_set_id` | [string](#string) |  |  |
| `threshold` | [uint32](#uint32) |  |  |
| `participants` | [string](#string) | repeated |  |
| `state` | [SigningState](#mpcchain.tss.v1.SigningState) |  |  |
| `start_height` | [int64](#int64) |  |  |
| `timeout_height` | [int64](#int64) |  |  |





 <!-- end messages -->


<a name="mpcchain.tss.v1.DKGState"></a>

### DKGState
DKGState defines the state of a DKG session

| Name | Number | Description |
| ---- | ------ | ----------- |
| DKG_STATE_UNSPECIFIED | 0 |  |
| DKG_STATE_ROUND1 | 1 |  |
| DKG_STATE_ROUND2 | 2 |  |
| DKG_STATE_COMPLETE | 3 |  |
| DKG_STATE_FAILED | 4 |  |



<a name="mpcchain.tss.v1.KeySetStatus"></a>

### KeySetStatus
KeySetStatus defines the status of a KeySet

| Name | Number | Description |
| ---- | ------ | ----------- |
| KEY_SET_STATUS_UNSPECIFIED | 0 |  |
| KEY_SET_STATUS_PENDING_DKG | 1 |  |
| KEY_SET_STATUS_ACTIVE | 2 |  |
| KEY_SET_STATUS_FAILED | 3 |  |



<a name="mpcchain.tss.v1.SigningRequestStatus"></a>

### SigningRequestStatus
SigningRequestStatus defines the status of a signing request

| Name | Number | Description |
| ---- | ------ | ----------- |
| SIGNING_REQUEST_STATUS_UNSPECIFIED | 0 |  |
| SIGNING_REQUEST_STATUS_PENDING | 1 |  |
| SIGNING_REQUEST_STATUS_ROUND1 | 2 |  |
| SIGNING_REQUEST_STATUS_ROUND2 | 3 |  |
| SIGNING_REQUEST_STATUS_COMPLETE | 4 |  |
| SIGNING_REQUEST_STATUS_FAILED | 5 |  |



<a name="mpcchain.tss.v1.SigningState"></a>

### SigningState
SigningState defines the state of a signing session

| Name | Number | Description |
| ---- | ------ | ----------- |
| SIGNING_STATE_UNSPECIFIED | 0 |  |
| SIGNING_STATE_ROUND1 | 1 |  |
| SIGNING_STATE_ROUND2 | 2 |  |
| SIGNING_STATE_COMPLETE | 3 |  |
| SIGNING_STATE_FAILED | 4 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="mpcchain/tss/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## mpcchain/tss/v1/genesis.proto



<a name="mpcchain.tss.v1.GenesisState"></a>

### GenesisState
GenesisState defines the TSS module genesis state


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#mpcchain.tss.v1.Params) |  |  |
| `key_sets` | [KeySet](#mpcchain.tss.v1.KeySet) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="mpcchain/tss/v1/module.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## mpcchain/tss/v1/module.proto



<a name="mpcchain.tss.v1.Module"></a>

### Module
Module is the config object of the TSS module


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authority` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="mpcchain/tss/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## mpcchain/tss/v1/query.proto



<a name="mpcchain.tss.v1.QueryAllDKGSessionsRequest"></a>

### QueryAllDKGSessionsRequest
QueryAllDKGSessionsRequest is the request type for the Query/AllDKGSessions RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="mpcchain.tss.v1.QueryAllDKGSessionsResponse"></a>

### QueryAllDKGSessionsResponse
QueryAllDKGSessionsResponse is the response type for the Query/AllDKGSessions RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sessions` | [DKGSession](#mpcchain.tss.v1.DKGSession) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  |  |






<a name="mpcchain.tss.v1.QueryAllKeySetsRequest"></a>

### QueryAllKeySetsRequest
QueryAllKeySetsRequest is the request type for the Query/AllKeySets RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="mpcchain.tss.v1.QueryAllKeySetsResponse"></a>

### QueryAllKeySetsResponse
QueryAllKeySetsResponse is the response type for the Query/AllKeySets RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key_sets` | [KeySet](#mpcchain.tss.v1.KeySet) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  |  |






<a name="mpcchain.tss.v1.QueryAllSigningRequestsRequest"></a>

### QueryAllSigningRequestsRequest
QueryAllSigningRequestsRequest is the request type for the Query/AllSigningRequests RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="mpcchain.tss.v1.QueryAllSigningRequestsResponse"></a>

### QueryAllSigningRequestsResponse
QueryAllSigningRequestsResponse is the response type for the Query/AllSigningRequests RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `requests` | [SigningRequest](#mpcchain.tss.v1.SigningRequest) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  |  |






<a name="mpcchain.tss.v1.QueryDKGSessionRequest"></a>

### QueryDKGSessionRequest
QueryDKGSessionRequest is the request type for the Query/DKGSession RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `session_id` | [string](#string) |  |  |






<a name="mpcchain.tss.v1.QueryDKGSessionResponse"></a>

### QueryDKGSessionResponse
QueryDKGSessionResponse is the response type for the Query/DKGSession RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `session` | [DKGSession](#mpcchain.tss.v1.DKGSession) |  |  |






<a name="mpcchain.tss.v1.QueryKeySetRequest"></a>

### QueryKeySetRequest
QueryKeySetRequest is the request type for the Query/KeySet RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |






<a name="mpcchain.tss.v1.QueryKeySetResponse"></a>

### QueryKeySetResponse
QueryKeySetResponse is the response type for the Query/KeySet RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key_set` | [KeySet](#mpcchain.tss.v1.KeySet) |  |  |






<a name="mpcchain.tss.v1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method






<a name="mpcchain.tss.v1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#mpcchain.tss.v1.Params) |  |  |






<a name="mpcchain.tss.v1.QuerySigningRequestRequest"></a>

### QuerySigningRequestRequest
QuerySigningRequestRequest is the request type for the Query/SigningRequest RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `request_id` | [string](#string) |  |  |






<a name="mpcchain.tss.v1.QuerySigningRequestResponse"></a>

### QuerySigningRequestResponse
QuerySigningRequestResponse is the response type for the Query/SigningRequest RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `request` | [SigningRequest](#mpcchain.tss.v1.SigningRequest) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="mpcchain.tss.v1.Query"></a>

### Query
Query defines the gRPC query service

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#mpcchain.tss.v1.QueryParamsRequest) | [QueryParamsResponse](#mpcchain.tss.v1.QueryParamsResponse) | Params queries module parameters | GET|/mpcchain/tss/v1/params|
| `KeySet` | [QueryKeySetRequest](#mpcchain.tss.v1.QueryKeySetRequest) | [QueryKeySetResponse](#mpcchain.tss.v1.QueryKeySetResponse) | KeySet queries a KeySet by ID | GET|/mpcchain/tss/v1/keyset/{id}|
| `AllKeySets` | [QueryAllKeySetsRequest](#mpcchain.tss.v1.QueryAllKeySetsRequest) | [QueryAllKeySetsResponse](#mpcchain.tss.v1.QueryAllKeySetsResponse) | AllKeySets queries all KeySets | GET|/mpcchain/tss/v1/keysets|
| `DKGSession` | [QueryDKGSessionRequest](#mpcchain.tss.v1.QueryDKGSessionRequest) | [QueryDKGSessionResponse](#mpcchain.tss.v1.QueryDKGSessionResponse) | DKGSession queries a DKG session by ID | GET|/mpcchain/tss/v1/dkg/{session_id}|
| `AllDKGSessions` | [QueryAllDKGSessionsRequest](#mpcchain.tss.v1.QueryAllDKGSessionsRequest) | [QueryAllDKGSessionsResponse](#mpcchain.tss.v1.QueryAllDKGSessionsResponse) | AllDKGSessions queries all DKG sessions | GET|/mpcchain/tss/v1/dkg|
| `SigningRequest` | [QuerySigningRequestRequest](#mpcchain.tss.v1.QuerySigningRequestRequest) | [QuerySigningRequestResponse](#mpcchain.tss.v1.QuerySigningRequestResponse) | SigningRequest queries a signing request by ID | GET|/mpcchain/tss/v1/signing/{request_id}|
| `AllSigningRequests` | [QueryAllSigningRequestsRequest](#mpcchain.tss.v1.QueryAllSigningRequestsRequest) | [QueryAllSigningRequestsResponse](#mpcchain.tss.v1.QueryAllSigningRequestsResponse) | AllSigningRequests queries all signing requests | GET|/mpcchain/tss/v1/signing|

 <!-- end services -->



<a name="mpcchain/tss/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## mpcchain/tss/v1/tx.proto



<a name="mpcchain.tss.v1.MsgCreateKeySet"></a>

### MsgCreateKeySet



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `creator` | [string](#string) |  |  |
| `threshold` | [uint32](#uint32) |  |  |
| `max_signers` | [uint32](#uint32) |  |  |
| `description` | [string](#string) |  |  |
| `timeout_blocks` | [int64](#int64) |  |  |






<a name="mpcchain.tss.v1.MsgCreateKeySetResponse"></a>

### MsgCreateKeySetResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key_set_id` | [string](#string) |  |  |
| `dkg_session_id` | [string](#string) |  |  |






<a name="mpcchain.tss.v1.MsgInitiateDKG"></a>

### MsgInitiateDKG



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authority` | [string](#string) |  |  |
| `key_set_id` | [string](#string) |  |  |






<a name="mpcchain.tss.v1.MsgInitiateDKGResponse"></a>

### MsgInitiateDKGResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `session_id` | [string](#string) |  |  |






<a name="mpcchain.tss.v1.MsgRequestSignature"></a>

### MsgRequestSignature



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `requester` | [string](#string) |  |  |
| `key_set_id` | [string](#string) |  |  |
| `message_hash` | [bytes](#bytes) |  |  |
| `callback` | [string](#string) |  |  |






<a name="mpcchain.tss.v1.MsgRequestSignatureResponse"></a>

### MsgRequestSignatureResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `request_id` | [string](#string) |  |  |






<a name="mpcchain.tss.v1.MsgSubmitCommitment"></a>

### MsgSubmitCommitment



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [string](#string) |  |  |
| `request_id` | [string](#string) |  |  |
| `commitment` | [bytes](#bytes) |  |  |






<a name="mpcchain.tss.v1.MsgSubmitCommitmentResponse"></a>

### MsgSubmitCommitmentResponse







<a name="mpcchain.tss.v1.MsgSubmitDKGRound1"></a>

### MsgSubmitDKGRound1



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [string](#string) |  |  |
| `session_id` | [string](#string) |  |  |
| `commitment` | [bytes](#bytes) |  |  |






<a name="mpcchain.tss.v1.MsgSubmitDKGRound1Response"></a>

### MsgSubmitDKGRound1Response







<a name="mpcchain.tss.v1.MsgSubmitDKGRound2"></a>

### MsgSubmitDKGRound2



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [string](#string) |  |  |
| `session_id` | [string](#string) |  |  |
| `share` | [bytes](#bytes) |  |  |






<a name="mpcchain.tss.v1.MsgSubmitDKGRound2Response"></a>

### MsgSubmitDKGRound2Response







<a name="mpcchain.tss.v1.MsgSubmitSignatureShare"></a>

### MsgSubmitSignatureShare



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [string](#string) |  |  |
| `request_id` | [string](#string) |  |  |
| `share` | [bytes](#bytes) |  |  |






<a name="mpcchain.tss.v1.MsgSubmitSignatureShareResponse"></a>

### MsgSubmitSignatureShareResponse







<a name="mpcchain.tss.v1.MsgUpdateParams"></a>

### MsgUpdateParams
MsgUpdateParams updates module parameters


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authority` | [string](#string) |  |  |
| `params` | [Params](#mpcchain.tss.v1.Params) |  |  |






<a name="mpcchain.tss.v1.MsgUpdateParamsResponse"></a>

### MsgUpdateParamsResponse






 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="mpcchain.tss.v1.Msg"></a>

### Msg
Msg defines the TSS Msg service (merged from MPC + Signing modules)

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `UpdateParams` | [MsgUpdateParams](#mpcchain.tss.v1.MsgUpdateParams) | [MsgUpdateParamsResponse](#mpcchain.tss.v1.MsgUpdateParamsResponse) | UpdateParams updates module parameters | |
| `CreateKeySet` | [MsgCreateKeySet](#mpcchain.tss.v1.MsgCreateKeySet) | [MsgCreateKeySetResponse](#mpcchain.tss.v1.MsgCreateKeySetResponse) | DKG Messages (from x/mpc) | |
| `InitiateDKG` | [MsgInitiateDKG](#mpcchain.tss.v1.MsgInitiateDKG) | [MsgInitiateDKGResponse](#mpcchain.tss.v1.MsgInitiateDKGResponse) |  | |
| `SubmitDKGRound1` | [MsgSubmitDKGRound1](#mpcchain.tss.v1.MsgSubmitDKGRound1) | [MsgSubmitDKGRound1Response](#mpcchain.tss.v1.MsgSubmitDKGRound1Response) |  | |
| `SubmitDKGRound2` | [MsgSubmitDKGRound2](#mpcchain.tss.v1.MsgSubmitDKGRound2) | [MsgSubmitDKGRound2Response](#mpcchain.tss.v1.MsgSubmitDKGRound2Response) |  | |
| `RequestSignature` | [MsgRequestSignature](#mpcchain.tss.v1.MsgRequestSignature) | [MsgRequestSignatureResponse](#mpcchain.tss.v1.MsgRequestSignatureResponse) | Signing Messages (from x/signing) | |
| `SubmitCommitment` | [MsgSubmitCommitment](#mpcchain.tss.v1.MsgSubmitCommitment) | [MsgSubmitCommitmentResponse](#mpcchain.tss.v1.MsgSubmitCommitmentResponse) |  | |
| `SubmitSignatureShare` | [MsgSubmitSignatureShare](#mpcchain.tss.v1.MsgSubmitSignatureShare) | [MsgSubmitSignatureShareResponse](#mpcchain.tss.v1.MsgSubmitSignatureShareResponse) |  | |

 <!-- end services -->



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

