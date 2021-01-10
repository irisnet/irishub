package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

// Rest variable names
// nolint
const (
	RestParamDenom  = "denom"
	RestParamSymbol = "symbol"
	RestParamOwner  = "owner"
)

// RegisterHandlers registers token-related REST handlers to a router
func RegisterHandlers(cliCtx client.Context, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}

type issueTokenReq struct {
	BaseReq       rest.BaseReq `json:"base_req"`
	Owner         string       `json:"owner"` // owner of the token
	Symbol        string       `json:"symbol"`
	Name          string       `json:"name"`
	Scale         uint32       `json:"scale"`
	MinUnit       string       `json:"min_unit"`
	InitialSupply uint64       `json:"initial_supply"`
	MaxSupply     uint64       `json:"max_supply"`
	Mintable      bool         `json:"mintable"`
}

type editTokenReq struct {
	BaseReq   rest.BaseReq `json:"base_req"`
	Owner     string       `json:"owner"` //  owner of the token
	MaxSupply uint64       `json:"max_supply"`
	Mintable  string       `json:"mintable"` // mintable of the token
	Name      string       `json:"name"`
}

type transferTokenOwnerReq struct {
	BaseReq  rest.BaseReq `json:"base_req"`
	SrcOwner string       `json:"src_owner"` // the current owner address of the token
	DstOwner string       `json:"dst_owner"` // the new owner
}

type mintTokenReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Owner   string       `json:"owner"`  // the current owner address of the token
	To      string       `json:"to"`     // address of minting token to
	Amount  uint64       `json:"amount"` // amount of minting token
}

type burnTokenReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Sender  string       `json:"owner"`  // the current owner address of the token
	Amount  uint64       `json:"amount"` // amount of burning token
}
