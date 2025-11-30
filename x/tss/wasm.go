package tss

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"

	"mpc-wasm-chain/x/tss/keeper"
	"mpc-wasm-chain/x/tss/types"
)

// CustomMessageEncoders returns message encoders for TSS module
func CustomMessageEncoders(k *keeper.Keeper) *wasmkeeper.MessageEncoders {
	return &wasmkeeper.MessageEncoders{
		Custom: CustomEncoder(k),
	}
}

// CustomEncoder creates a custom message encoder for TSS messages
func CustomEncoder(k *keeper.Keeper) wasmkeeper.CustomEncoder {
	return func(sender sdk.AccAddress, msg json.RawMessage) ([]sdk.Msg, error) {
		var tssMsg TSSMsg
		err := json.Unmarshal(msg, &tssMsg)
		if err != nil {
			return nil, errorsmod.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
		}

		// Handle CreateKeySet - initiates DKG ceremony
		if tssMsg.CreateKeySet != nil {
			return []sdk.Msg{&types.MsgCreateKeySet{
				Creator:       sender.String(),
				Threshold:     tssMsg.CreateKeySet.Threshold,
				MaxSigners:    tssMsg.CreateKeySet.MaxSigners,
				Description:   tssMsg.CreateKeySet.Description,
				TimeoutBlocks: tssMsg.CreateKeySet.TimeoutBlocks,
			}}, nil
		}

		// Handle RequestSignature - requests threshold signature
		if tssMsg.RequestSignature != nil {
			return []sdk.Msg{&types.MsgRequestSignature{
				Requester:   sender.String(),
				KeySetId:    tssMsg.RequestSignature.KeySetId,
				MessageHash: tssMsg.RequestSignature.MessageHash,
				Callback:    tssMsg.RequestSignature.Callback,
			}}, nil
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "unknown TSS message variant")
	}
}

// CustomQueryPlugins returns query plugins for TSS module
func CustomQueryPlugins(k *keeper.Keeper) *wasmkeeper.QueryPlugins {
	return &wasmkeeper.QueryPlugins{
		Custom: CustomQuerier(k),
	}
}

// CustomQuerier creates a custom query handler for TSS queries
func CustomQuerier(k *keeper.Keeper) wasmkeeper.CustomQuerier {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		var tssQuery TSSQuery
		err := json.Unmarshal(request, &tssQuery)
		if err != nil {
			return nil, errorsmod.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
		}

		// Handle KeySet query
		if tssQuery.KeySet != nil {
			keySet, err := k.KeySetStore.Get(ctx, tssQuery.KeySet.Id)
			if err != nil {
				return nil, errorsmod.Wrap(sdkerrors.ErrNotFound, "keyset not found")
			}
			return json.Marshal(KeySetResponse{
				Id:            keySet.Id,
				Owner:         keySet.Owner,
				Threshold:     keySet.Threshold,
				MaxSigners:    keySet.MaxSigners,
				Participants:  keySet.Participants,
				GroupPubkey:   keySet.GroupPubkey,
				Status:        keySet.Status.String(),
				Description:   keySet.Description,
				CreatedHeight: keySet.CreatedHeight,
			})
		}

		// Handle SigningRequest query
		if tssQuery.SigningRequest != nil {
			req, err := k.SigningRequestStore.Get(ctx, tssQuery.SigningRequest.Id)
			if err != nil {
				return nil, errorsmod.Wrap(sdkerrors.ErrNotFound, "signing request not found")
			}
			return json.Marshal(SigningRequestResponse{
				Id:            req.Id,
				KeySetId:      req.KeySetId,
				Requester:     req.Requester,
				MessageHash:   req.MessageHash,
				Callback:      req.Callback,
				Status:        req.Status.String(),
				Signature:     req.Signature,
				CreatedHeight: req.CreatedHeight,
			})
		}

		// Handle DKGSession query
		if tssQuery.DKGSession != nil {
			session, err := k.DKGSessionStore.Get(ctx, tssQuery.DKGSession.Id)
			if err != nil {
				return nil, errorsmod.Wrap(sdkerrors.ErrNotFound, "dkg session not found")
			}
			return json.Marshal(DKGSessionResponse{
				Id:            session.Id,
				KeySetId:      session.KeySetId,
				State:         session.State.String(),
				Threshold:     session.Threshold,
				MaxSigners:    session.MaxSigners,
				Participants:  session.Participants,
				StartHeight:   session.StartHeight,
				TimeoutHeight: session.TimeoutHeight,
			})
		}

		// Handle SigningSession query
		if tssQuery.SigningSession != nil {
			session, err := k.SigningSessionStore.Get(ctx, tssQuery.SigningSession.RequestId)
			if err != nil {
				return nil, errorsmod.Wrap(sdkerrors.ErrNotFound, "signing session not found")
			}
			return json.Marshal(SigningSessionResponse{
				RequestId:     session.RequestId,
				KeySetId:      session.KeySetId,
				Threshold:     session.Threshold,
				Participants:  session.Participants,
				State:         session.State.String(),
				StartHeight:   session.StartHeight,
				TimeoutHeight: session.TimeoutHeight,
			})
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "unknown TSS query variant")
	}
}

// Message types for WASM contract integration
// Note: Only contract-facing operations are exposed here.
// Validator operations (DKG rounds, signing commitments/shares) are handled
// automatically by the EndBlocker and should not be called by contracts.

type TSSMsg struct {
	CreateKeySet     *CreateKeySetMsg     `json:"create_key_set,omitempty"`
	RequestSignature *RequestSignatureMsg `json:"request_signature,omitempty"`
}

type CreateKeySetMsg struct {
	Threshold     uint32 `json:"threshold"`
	MaxSigners    uint32 `json:"max_signers"`
	Description   string `json:"description"`
	TimeoutBlocks int64  `json:"timeout_blocks,omitempty"`
}

type RequestSignatureMsg struct {
	KeySetId    string `json:"key_set_id"`
	MessageHash []byte `json:"message_hash"`
	Callback    string `json:"callback,omitempty"`
}

// Query types for WASM contract integration

type TSSQuery struct {
	KeySet         *KeySetQuery         `json:"key_set,omitempty"`
	SigningRequest *SigningRequestQuery `json:"signing_request,omitempty"`
	DKGSession     *DKGSessionQuery     `json:"dkg_session,omitempty"`
	SigningSession *SigningSessionQuery `json:"signing_session,omitempty"`
}

type KeySetQuery struct {
	Id string `json:"id"`
}

type SigningRequestQuery struct {
	Id string `json:"id"`
}

type DKGSessionQuery struct {
	Id string `json:"id"`
}

type SigningSessionQuery struct {
	RequestId string `json:"request_id"`
}

// Response types

type KeySetResponse struct {
	Id            string   `json:"id"`
	Owner         string   `json:"owner"`
	Threshold     uint32   `json:"threshold"`
	MaxSigners    uint32   `json:"max_signers"`
	Participants  []string `json:"participants"`
	GroupPubkey   []byte   `json:"group_pubkey"`
	Status        string   `json:"status"`
	Description   string   `json:"description"`
	CreatedHeight int64    `json:"created_height"`
}

type SigningRequestResponse struct {
	Id            string `json:"id"`
	KeySetId      string `json:"key_set_id"`
	Requester     string `json:"requester"`
	MessageHash   []byte `json:"message_hash"`
	Callback      string `json:"callback"`
	Status        string `json:"status"`
	Signature     []byte `json:"signature"`
	CreatedHeight int64  `json:"created_height"`
}

type DKGSessionResponse struct {
	Id            string   `json:"id"`
	KeySetId      string   `json:"key_set_id"`
	State         string   `json:"state"`
	Threshold     uint32   `json:"threshold"`
	MaxSigners    uint32   `json:"max_signers"`
	Participants  []string `json:"participants"`
	StartHeight   int64    `json:"start_height"`
	TimeoutHeight int64    `json:"timeout_height"`
}

type SigningSessionResponse struct {
	RequestId     string   `json:"request_id"`
	KeySetId      string   `json:"key_set_id"`
	Threshold     uint32   `json:"threshold"`
	Participants  []string `json:"participants"`
	State         string   `json:"state"`
	StartHeight   int64    `json:"start_height"`
	TimeoutHeight int64    `json:"timeout_height"`
}
