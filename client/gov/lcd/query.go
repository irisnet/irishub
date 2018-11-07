package lcd

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/modules/gov"
	"net/http"
	"github.com/irisnet/irishub/client/utils"
	"github.com/pkg/errors"
	client "github.com/irisnet/irishub/client/gov"
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

		proposalID, ok := utils.ParseInt64OrReturnBadRequest(w, strProposalID)
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

		proposalID, ok := utils.ParseInt64OrReturnBadRequest(w, strProposalID)
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
		bechDepositerAddr := vars[RestDepositer]

		if len(strProposalID) == 0 {
			err := errors.New("proposalId required but not specified")
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		proposalID, ok := utils.ParseInt64OrReturnBadRequest(w, strProposalID)
		if !ok {
			return
		}

		if len(bechDepositerAddr) == 0 {
			err := errors.New("depositer address required but not specified")
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		depositerAddr, err := sdk.AccAddressFromBech32(bechDepositerAddr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := gov.QueryDepositParams{
			ProposalID: proposalID,
			Depositer:  depositerAddr,
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
			res, err := cliCtx.QueryWithData("custom/gov/proposal", cdc.MustMarshalBinary(gov.QueryProposalParams{params.ProposalID}))
			if err != nil || len(res) == 0 {
				err := errors.Errorf("proposalID [%d] does not exist", proposalID)
				utils.WriteErrorResponse(w, http.StatusNotFound, err.Error())
				return
			}
			err = errors.Errorf("depositer [%s] did not deposit on proposalID [%d]", bechDepositerAddr, proposalID)
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

		proposalID, ok := utils.ParseInt64OrReturnBadRequest(w, strProposalID)
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

		proposalID, ok := utils.ParseInt64OrReturnBadRequest(w, strProposalID)
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

		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

// todo: Split this functionality into helper functions to remove the above
func queryProposalsWithParameterFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bechVoterAddr := r.URL.Query().Get(RestVoter)
		bechDepositerAddr := r.URL.Query().Get(RestDepositer)
		strProposalStatus := r.URL.Query().Get(RestProposalStatus)
		strNumLatest := r.URL.Query().Get(RestNumLatest)

		params := gov.QueryProposalsParams{}

		if len(bechVoterAddr) != 0 {
			voterAddr, err := sdk.AccAddressFromBech32(bechVoterAddr)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			params.Voter = voterAddr
		}

		if len(bechDepositerAddr) != 0 {
			depositerAddr, err := sdk.AccAddressFromBech32(bechDepositerAddr)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			params.Depositer = depositerAddr
		}

		if len(strProposalStatus) != 0 {
			proposalStatus, err := gov.ProposalStatusFromString(client.NormalizeProposalStatus(strProposalStatus))
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			params.ProposalStatus = proposalStatus
		}
		if len(strNumLatest) != 0 {
			numLatest, ok := utils.ParseInt64OrReturnBadRequest(w, strNumLatest)
			if !ok {
				return
			}
			params.NumLatestProposals = numLatest
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

		proposalID, ok := utils.ParseInt64OrReturnBadRequest(w, strProposalID)
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
		res, err := cliCtx.QuerySubspace([]byte("Gov/"), "params")
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		var pd gov.ParameterConfigFile
		for _, kv := range res {
			switch string(kv.Key) {
			case "Gov/govDepositProcedure":
				cdc.UnmarshalJSON(kv.Value, &pd.Govparams.DepositProcedure)
			case "Gov/govVotingProcedure":
				cdc.UnmarshalJSON(kv.Value, &pd.Govparams.VotingProcedure)
			case "Gov/govTallyingProcedure":
				cdc.UnmarshalJSON(kv.Value, &pd.Govparams.TallyingProcedure)
			default:
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
		}
		utils.PostProcessResponse(w, cdc, pd, cliCtx.Indent)
	}
}
