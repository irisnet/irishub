package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

// RegisterHandlers registers the NFT REST routes.
func RegisterHandlers(cliCtx client.Context, r *mux.Router, queryRoute string) {
	registerQueryRoutes(cliCtx, r, queryRoute)
	registerTxRoutes(cliCtx, r, queryRoute)
}

const (
	RestParamDenomID = "denom-id"
	RestParamTokenID = "token-id"
	RestParamOwner   = "owner"
)

type issueDenomReq struct {
	BaseReq          rest.BaseReq `json:"base_req"`
	Owner            string       `json:"owner"`
	ID               string       `json:"id"`
	Name             string       `json:"name"`
	Schema           string       `json:"schema"`
	Symbol           string       `json:"symbol"`
	MintRestricted   bool         `json:"mint_restricted"`
	UpdateRestricted bool         `json:"update_restricted"`
	Description      string       `json:"description"`
	Uri              string       `json:"uri"`
	UriHash          string       `json:"uri_hash"`
	Data             string       `json:"data"`
}

type mintNFTReq struct {
	BaseReq   rest.BaseReq `json:"base_req"`
	Owner     string       `json:"owner"`
	Recipient string       `json:"recipient"`
	DenomID   string       `json:"denom_id"`
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	URI       string       `json:"uri"`
	UriHash   string       `json:"uri_hash"`
	Data      string       `json:"data"`
}

type editNFTReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Owner   string       `json:"owner"`
	Name    string       `json:"name"`
	URI     string       `json:"uri"`
	UriHash string       `json:"uri_hash"`
	Data    string       `json:"data"`
}

type transferNFTReq struct {
	BaseReq   rest.BaseReq `json:"base_req"`
	Owner     string       `json:"owner"`
	Recipient string       `json:"recipient"`
	Name      string       `json:"name"`
	URI       string       `json:"uri"`
	UriHash   string       `json:"uri_hash"`
	Data      string       `json:"data"`
}

type burnNFTReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Owner   string       `json:"owner"`
}
