package lcd

import (
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/codec"
)

const storeName = "distr"
const stakeStoreName = "stake"

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/distribution/{delegatorAddr}/withdrawAddress", SetWithdrawAddressHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/distribution/{delegatorAddr}/withdrawReward", WithdrawRewardsHandlerFn(cdc, cliCtx)).Methods("POST")

	r.HandleFunc("/distribution/{delegatorAddr}/withdrawAddress",
		QueryWithdrawAddressHandlerFn(storeName, cliCtx)).Methods("GET")
	r.HandleFunc("/distribution/{delegatorAddr}/distrInfo/{validatorAddr}",
		QueryDelegationDistInfoHandlerFn(storeName, cliCtx)).Methods("GET")
	r.HandleFunc("/distribution/{delegatorAddr}/distrInfos",
		QueryDelegatorDistInfoHandlerFn(storeName, cliCtx)).Methods("GET")
	r.HandleFunc("/distribution/{validatorAddr}/valDistrInfo",
		QueryValidatorDistInfoHandlerFn(storeName, cliCtx)).Methods("GET")
	r.HandleFunc("/distribution/{address}/rewards",
		QueryRewardsHandlerFn(storeName, stakeStoreName, cliCtx)).Methods("GET")
}
