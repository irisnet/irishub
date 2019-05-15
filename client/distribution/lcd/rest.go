package lcd

import (
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/codec"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/distribution/{delegatorAddr}/withdraw-address", SetWithdrawAddressHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/distribution/{delegatorAddr}/withdraw-reward", WithdrawRewardsHandlerFn(cdc, cliCtx)).Methods("POST")

	r.HandleFunc("/distribution/{delegatorAddr}/withdraw-address",
		QueryWithdrawAddressHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/distribution/{address}/rewards",
		QueryRewardsHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/distribution/community-tax",
		QueryCommunityTaxFn(cliCtx)).Methods("GET")

	//r.HandleFunc("/distribution/{delegatorAddr}/distrInfo/{validatorAddr}",
	//	QueryDelegationDistInfoHandlerFn(cliCtx)).Methods("GET")
	//r.HandleFunc("/distribution/{delegatorAddr}/distrInfo",
	//	QueryDelegatorDistInfoHandlerFn(cliCtx)).Methods("GET")
	//r.HandleFunc("/distribution/{validatorAddr}/valDistrInfo",
	//	QueryValidatorDistInfoHandlerFn(cliCtx)).Methods("GET")
}
