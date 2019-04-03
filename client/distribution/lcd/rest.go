package lcd

import (
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/codec"
)

const queryRoute = "distr"

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/distribution/{delegatorAddr}/withdrawAddress", SetWithdrawAddressHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/distribution/{delegatorAddr}/withdrawReward", WithdrawRewardsHandlerFn(cdc, cliCtx)).Methods("POST")

	r.HandleFunc("/distribution/{delegatorAddr}/withdrawAddress",
		QueryWithdrawAddressHandlerFn(queryRoute, cliCtx)).Methods("GET")
	r.HandleFunc("/distribution/{delegatorAddr}/distrInfo/{validatorAddr}",
		QueryDelegationDistInfoHandlerFn(queryRoute, cliCtx)).Methods("GET")
	r.HandleFunc("/distribution/{delegatorAddr}/distrInfos",
		QueryDelegatorDistInfoHandlerFn(queryRoute, cliCtx)).Methods("GET")
	r.HandleFunc("/distribution/{validatorAddr}/valDistrInfo",
		QueryValidatorDistInfoHandlerFn(queryRoute, cliCtx)).Methods("GET")
	r.HandleFunc("/distribution/{address}/rewards",
		QueryRewardsHandlerFn(queryRoute, cliCtx)).Methods("GET")
}
