package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

// RegisterHandlers register distribution REST routes.
func RegisterHandlers(cliCtx client.Context, r *mux.Router, queryRoute string) {
	registerQueryRoutes(cliCtx, r, queryRoute)
	registerTxRoutes(cliCtx, r, queryRoute)
}

const (
	RestParamDenom   = "denom"
	RestParamTokenID = "id"
	RestParamOwner   = "owner"
)

type issueDenomReq struct {
	BaseReq rest.BaseReq   `json:"base_req"`
	Owner   sdk.AccAddress `json:"owner"`
	ID      string         `json:"id"`
	Name    string         `json:"name"`
	Schema  string         `json:"schema"`
}

type mintNFTReq struct {
	BaseReq   rest.BaseReq   `json:"base_req"`
	Owner     sdk.AccAddress `json:"owner"`
	Recipient sdk.AccAddress `json:"recipient"`
	Denom     string         `json:"denom"`
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	URI       string         `json:"uri"`
	Data      string         `json:"data"`
}

type editNFTReq struct {
	BaseReq rest.BaseReq   `json:"base_req"`
	Owner   sdk.AccAddress `json:"owner"`
	Name    string         `json:"name"`
	URI     string         `json:"uri"`
	Data    string         `json:"data"`
}

type transferNFTReq struct {
	BaseReq   rest.BaseReq   `json:"base_req"`
	Owner     sdk.AccAddress `json:"owner"`
	Recipient string         `json:"recipient"`
	Name      string         `json:"name"`
	URI       string         `json:"uri"`
	Data      string         `json:"data"`
}

type burnNFTReq struct {
	BaseReq rest.BaseReq   `json:"base_req"`
	Owner   sdk.AccAddress `json:"owner"`
}
