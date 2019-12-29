package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

// Rest variable names
// nolint
const (
	RestParamSymbol = "symbol"
	RestParamOwner  = "owner"
)

// RegisterRoutes registers asset-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	registerQueryRoutes(cliCtx, r, queryRoute)
	registerTxRoutes(cliCtx, r)
}

// IssueTokenReq defines the properties of a issue token request's body
type IssueTokenReq struct {
	BaseTx        rest.BaseReq   `json:"base_tx" yaml:"base_tx"`
	Owner         sdk.AccAddress `json:"owner" yaml:"owner"` //  Owner of the token
	Symbol        string         `json:"symbol" yaml:"symbol"`
	Name          string         `json:"name" yaml:"name"`
	Scale         uint8          `json:"scale" yaml:"scale"`
	MinUnit       string         `json:"min_unit" yaml:"min_unit"`
	InitialSupply uint64         `json:"initial_supply" yaml:"initial_supply"`
	MaxSupply     uint64         `json:"max_supply" yaml:"max_supply"`
	Mintable      bool           `json:"mintable" yaml:"mintable"`
}

// EditTokenReq defines the properties of a edit token request's body
type EditTokenReq struct {
	BaseTx    rest.BaseReq   `json:"base_tx" yaml:"base_tx"`
	Owner     sdk.AccAddress `json:"owner" yaml:"owner"` //  owner of asset
	MaxSupply uint64         `json:"max_supply" yaml:"max_supply"`
	Mintable  string         `json:"mintable" yaml:"mintable"` //  mintable of asset
	Name      string         `json:"name" yaml:"name"`
}

// TransferTokenReq defines the properties of a transfer token owner request's body
type TransferTokenReq struct {
	BaseTx   rest.BaseReq   `json:"base_tx" yaml:"base_tx"`
	Symbol   string         `json:"symbol" yaml:"symbol"`
	SrcOwner sdk.AccAddress `json:"src_owner" yaml:"src_owner"` // the current owner address of the token
	DstOwner sdk.AccAddress `json:"dst_owner" yaml:"dst_owner"` // the new owner
}

// MintTokenReq defines the properties of a mint token request's body
type MintTokenReq struct {
	BaseTx    rest.BaseReq   `json:"base_tx" yaml:"base_tx"`
	Owner     sdk.AccAddress `json:"owner" yaml:"owner"`         // the current owner address of the token
	Recipient sdk.AccAddress `json:"recipient" yaml:"recipient"` // address of mint token to
	Amount    uint64         `json:"amount" yaml:"amount"`       // amount of mint token
}

// BurnTokenReq defines the properties of a burn token request's body
type BurnTokenReq struct {
	BaseTx rest.BaseReq   `json:"base_tx"`
	Sender sdk.AccAddress `json:"sender"` // the current sender address of the token
	Amount sdk.Coins      `json:"amount"` // amount of burn token
}
