package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/irisnet/irishub/modules/asset/internal/types"
)

// Rest variable names
// nolint
const (
	RestTokenID = "token-id"
)

// RegisterRoutes registers asset-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	registerQueryRoutes(cliCtx, r, queryRoute)
	registerTxRoutes(cliCtx, r)
}

type IssueTokenReq struct {
	BaseTx          rest.BaseReq      `json:"base_tx" yaml:"base_tx"`
	Owner           sdk.AccAddress    `json:"owner" yaml:"owner"` //  Owner of the token
	Family          types.AssetFamily `json:"family" yaml:"family"`
	Source          types.AssetSource `json:"source" yaml:"source"`
	Symbol          string            `json:"symbol" yaml:"symbol"`
	CanonicalSymbol string            `json:"canonical_symbol" yaml:"canonical_symbol"`
	Name            string            `json:"name" yaml:"name"`
	Decimal         uint8             `json:"decimal" yaml:"decimal"`
	MinUnitAlias    string            `json:"min_unit_alias" yaml:"min_unit_alias"`
	InitialSupply   uint64            `json:"initial_supply" yaml:"initial_supply"`
	MaxSupply       uint64            `json:"max_supply" yaml:"max_supply"`
	Mintable        bool              `json:"mintable" yaml:"mintable"`
}

type EditTokenReq struct {
	BaseTx          rest.BaseReq   `json:"base_tx" yaml:"base_tx"`
	Owner           sdk.AccAddress `json:"owner" yaml:"owner"`                       //  owner of asset
	CanonicalSymbol string         `json:"canonical_symbol" yaml:"canonical_symbol"` //  canonical_symbol of asset
	MinUnitAlias    string         `json:"min_unit_alias" yaml:"min_unit_alias"`     //  min_unit_alias of asset
	MaxSupply       uint64         `json:"max_supply" yaml:"max_supply"`
	Mintable        string         `json:"mintable" yaml:"mintable"` //  mintable of asset
	Name            string         `json:"name" yaml:"name"`
}

type TransferTokenOwnerReq struct {
	BaseTx   rest.BaseReq   `json:"base_tx" yaml:"base_tx"`
	SrcOwner sdk.AccAddress `json:"src_owner" yaml:"src_owner"` // the current owner address of the token
	DstOwner sdk.AccAddress `json:"dst_owner" yaml:"dst_owner"` // the new owner
}

type MintTokenReq struct {
	BaseTx rest.BaseReq   `json:"base_tx" yaml:"base_tx"`
	Owner  sdk.AccAddress `json:"owner" yaml:"owner"`   // the current owner address of the token
	To     sdk.AccAddress `json:"to" yaml:"to"`         // address of mint token to
	Amount uint64         `json:"amount" yaml:"amount"` // amount of mint token
}
