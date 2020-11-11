package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/irisnet/irismod/modules/record/types"
)

// Rest variable names
// nolint
const (
	RestRecordID = "record-id"
)

// RegisterHandlers defines routes that get registered by the main application
func RegisterHandlers(cliCtx client.Context, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}

type RecordCreateReq struct {
	BaseReq  rest.BaseReq    `json:"base_req" yaml:"base_req"` // base req
	Contents []types.Content `json:"contents" yaml:"contents"`
	Creator  string          `json:"creator" yaml:"creator"`
}
