package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/modules/asset/types"
)

// RegisterRoutes registers asset-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	registerQueryRoutes(cliCtx, r, queryRoute)
	registerTxRoutes(cliCtx, r)
}

type issueTokenReq struct {
	BaseTx          rest.BaseReq      `json:"base_tx"`
	Owner           sdk.AccAddress    `json:"owner"` //  Owner of the token
	Family          types.AssetFamily `json:"family"`
	Source          types.AssetSource `json:"source"`
	Symbol          string            `json:"symbol"`
	CanonicalSymbol string            `json:"canonical_symbol"`
	Name            string            `json:"name"`
	Decimal         uint8             `json:"decimal"`
	MinUnitAlias    string            `json:"min_unit_alias"`
	InitialSupply   uint64            `json:"initial_supply"`
	MaxSupply       uint64            `json:"max_supply"`
	Mintable        bool              `json:"mintable"`
}

type editTokenReq struct {
	BaseTx          rest.BaseReq   `json:"base_tx"`
	Owner           sdk.AccAddress `json:"owner"`            //  owner of asset
	CanonicalSymbol string         `json:"canonical_symbol"` //  canonical_symbol of asset
	MinUnitAlias    string         `json:"min_unit_alias"`   //  min_unit_alias of asset
	MaxSupply       uint64         `json:"max_supply"`
	Mintable        string         `json:"mintable"` //  mintable of asset
	Name            string         `json:"name"`
}

type transferTokenOwnerReq struct {
	BaseTx   rest.BaseReq   `json:"base_tx"`
	SrcOwner sdk.AccAddress `json:"src_owner"` // the current owner address of the token
	DstOwner sdk.AccAddress `json:"dst_owner"` // the new owner
}

type mintTokenReq struct {
	BaseTx rest.BaseReq   `json:"base_tx"`
	Owner  sdk.AccAddress `json:"owner"`  // the current owner address of the token
	To     sdk.AccAddress `json:"to"`     // address of mint token to
	Amount uint64         `json:"amount"` // amount of mint token
}
