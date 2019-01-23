package lcd

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/context"
	client "github.com/irisnet/irishub/client/gov"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/app/v1/gov"
	sdk "github.com/irisnet/irishub/types"
	"github.com/pkg/errors"
)

func queryProposalHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]

		if len(strProposalID) == 0 {
			err := errors.New("proposalId required but not specified")
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		proposalID, ok := utils.ParseUint64OrReturnBadRequest(w, strProposalID)
		if !ok {
			return
		}

		params := gov.QueryProposalParams{
			ProposalID: proposalID,
		}

		bz, err := cdc.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := cliCtx.QueryWithData("custom/gov/proposal", bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryDepositsHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]

		proposalID, ok := utils.ParseUint64OrReturnBadRequest(w, strProposalID)
		if !ok {
			return
		}

		params := gov.QueryDepositsParams{
			ProposalID: proposalID,
		}

		bz, err := cdc.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := cliCtx.QueryWithData("custom/gov/deposits", bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryDepositHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]
		bechDepositorAddr := vars[RestDepositor]

		if len(strProposalID) == 0 {
			err := errors.New("proposalId required but not specified")
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		proposalID, ok := utils.ParseUint64OrReturnBadRequest(w, strProposalID)
		if !ok {
			return
		}

		if len(bechDepositorAddr) == 0 {
			err := errors.New("depositor address required but not specified")
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		depositorAddr, err := sdk.AccAddressFromBech32(bechDepositorAddr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := gov.QueryDepositParams{
			ProposalID: proposalID,
			Depositor:  depositorAddr,
		}

		bz, err := cdc.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := cliCtx.QueryWithData("custom/gov/deposit", bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var deposit gov.Deposit
		cdc.UnmarshalJSON(res, &deposit)
		if deposit.Empty() {
			res, err := cliCtx.QueryWithData("custom/gov/proposal", cdc.MustMarshalBinaryLengthPrefixed(gov.QueryProposalParams{params.ProposalID}))
			if err != nil || len(res) == 0 {
				err := errors.Errorf("proposalID [%d] does not exist", proposalID)
				utils.WriteErrorResponse(w, http.StatusNotFound, err.Error())
				return
			}
			err = errors.Errorf("depositor [%s] did not deposit on proposalID [%d]", bechDepositorAddr, proposalID)
			utils.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryVoteHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]
		bechVoterAddr := vars[RestVoter]

		if len(strProposalID) == 0 {
			err := errors.New("proposalId required but not specified")
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		proposalID, ok := utils.ParseUint64OrReturnBadRequest(w, strProposalID)
		if !ok {
			return
		}

		if len(bechVoterAddr) == 0 {
			err := errors.New("voter address required but not specified")
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		voterAddr, err := sdk.AccAddressFromBech32(bechVoterAddr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := gov.QueryVoteParams{
			Voter:      voterAddr,
			ProposalID: proposalID,
		}
		bz, err := cdc.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := cliCtx.QueryWithData("custom/gov/vote", bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var vote gov.Vote
		cdc.UnmarshalJSON(res, &vote)
		if vote.Empty() {
			bz, err := cdc.MarshalJSON(gov.QueryProposalParams{params.ProposalID})
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			res, err := cliCtx.QueryWithData("custom/gov/proposal", bz)
			if err != nil || len(res) == 0 {
				err := errors.Errorf("proposalID [%d] does not exist", proposalID)
				utils.WriteErrorResponse(w, http.StatusNotFound, err.Error())
				return
			}
			err = errors.Errorf("voter [%s] did not deposit on proposalID [%d]", bechVoterAddr, proposalID)
			utils.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

// todo: Split this functionality into helper functions to remove the above
func queryVotesOnProposalHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]

		if len(strProposalID) == 0 {
			err := errors.New("proposalId required but not specified")
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		proposalID, ok := utils.ParseUint64OrReturnBadRequest(w, strProposalID)
		if !ok {
			return
		}

		params := gov.QueryVotesParams{
			ProposalID: proposalID,
		}
		bz, err := cdc.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := cliCtx.QueryWithData("custom/gov/votes", bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if res == nil {
			res = []byte(fmt.Sprintf("No one votes for the proposal [%v].\n", proposalID))
		}

		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

// todo: Split this functionality into helper functions to remove the above
func queryProposalsWithParameterFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bechVoterAddr := r.URL.Query().Get(RestVoter)
		bechDepositorAddr := r.URL.Query().Get(RestDepositor)
		strProposalStatus := r.URL.Query().Get(RestProposalStatus)
		strNumLimit := r.URL.Query().Get(RestNumLimit)

		params := gov.QueryProposalsParams{}

		if len(bechVoterAddr) != 0 {
			voterAddr, err := sdk.AccAddressFromBech32(bechVoterAddr)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			params.Voter = voterAddr
		}

		if len(bechDepositorAddr) != 0 {
			depositorAddr, err := sdk.AccAddressFromBech32(bechDepositorAddr)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			params.Depositor = depositorAddr
		}

		if len(strProposalStatus) != 0 {
			proposalStatus, err := gov.ProposalStatusFromString(client.NormalizeProposalStatus(strProposalStatus))
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			params.ProposalStatus = proposalStatus
		}
		if len(strNumLimit) != 0 {
			numLatest, ok := utils.ParseUint64OrReturnBadRequest(w, strNumLimit)
			if !ok {
				return
			}
			params.Limit = numLatest
		}

		bz, err := cdc.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := cliCtx.QueryWithData("custom/gov/proposals", bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

// todo: Split this functionality into helper functions to remove the above
func queryTallyOnProposalHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]

		if len(strProposalID) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			err := errors.New("proposalId required but not specified")
			w.Write([]byte(err.Error()))

			return
		}

		proposalID, ok := utils.ParseUint64OrReturnBadRequest(w, strProposalID)
		if !ok {
			return
		}

		params := gov.QueryTallyParams{
			ProposalID: proposalID,
		}
		bz, err := cdc.MarshalJSON(params)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		res, err := cliCtx.QueryWithData("custom/gov/tally", bz)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

// nolint: gocyclo
func queryParamsHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.PostProcessResponse(w, cdc, "", cliCtx.Indent)
	}
}
