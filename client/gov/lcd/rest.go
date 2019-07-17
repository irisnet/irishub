package lcd

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/codec"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/gov/proposals", postProposalHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/gov/proposals/{%s}/deposits", RestProposalID), depositHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/gov/proposals/{%s}/votes", RestProposalID), voteHandlerFn(cdc, cliCtx)).Methods("POST")

	r.HandleFunc("/gov/proposals", queryProposalsWithParameterFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/gov/proposals/{%s}", RestProposalID), queryProposalHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/gov/proposals/{%s}/deposits", RestProposalID), queryDepositsHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/gov/proposals/{%s}/deposits/{%s}", RestProposalID, RestDepositor), queryDepositHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/gov/proposals/{%s}/votes", RestProposalID), queryVotesOnProposalHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/gov/proposals/{%s}/votes/{%s}", RestProposalID, RestVoter), queryVoteHandlerFn(cdc, cliCtx)).Methods("GET")

	r.HandleFunc(fmt.Sprintf("/gov/proposals/{%s}/tally_result", RestProposalID), queryTallyOnProposalHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/gov/params/{%s}", Module), queryParamsHandlerFn(cdc, cliCtx)).Methods("GET")
}
