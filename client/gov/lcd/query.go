package lcd

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/context"
	govClient "github.com/irisnet/irishub/client/gov"
	"github.com/irisnet/irishub/modules/gov"
	"net/http"
	"strconv"
	"github.com/irisnet/irishub/client/utils"
)

func queryProposalHandlerFn(cdc *wire.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]

		if len(strProposalID) == 0 {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "proposalId required but not specified")
			return
		}

		proposalID, err := strconv.ParseInt(strProposalID, 10, 64)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("proposalID [%d] is not positive", proposalID))
			return
		}

		res, err := cliCtx.QueryStore(gov.KeyProposal(proposalID), storeName)
		if err != nil || len(res) == 0 {
			utils.WriteErrorResponse(w, http.StatusNotFound, fmt.Sprintf("proposalID [%d] does not exist", proposalID))
			return
		}

		var proposal gov.Proposal
		cdc.MustUnmarshalBinary(res, &proposal)
		proposalResponse, err := govClient.ConvertProposalCoins(cliCtx, proposal)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		output, err := wire.MarshalJSONIndent(cdc, proposalResponse)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Write(output)
	}
}

func queryDepositHandlerFn(cdc *wire.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]
		bechDepositerAddr := vars[RestDepositer]

		if len(strProposalID) == 0 {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "proposalId required but not specified")
			return
		}

		proposalID, err := strconv.ParseInt(strProposalID, 10, 64)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("proposalID [%d] is not positive", proposalID))
			return
		}

		if len(bechDepositerAddr) == 0 {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "depositer address required but not specified")
			return
		}

		depositerAddr, err := sdk.AccAddressFromBech32(bechDepositerAddr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("'%s' needs to be bech32 encoded", RestDepositer))
			return
		}

		res, err := cliCtx.QueryStore(gov.KeyDeposit(proposalID, depositerAddr), storeName)
		if err != nil || len(res) == 0 {
			res, err := cliCtx.QueryStore(gov.KeyProposal(proposalID), storeName)
			if err != nil || len(res) == 0 {
				utils.WriteErrorResponse(w, http.StatusNotFound, fmt.Sprintf("proposalID [%d] does not exist", proposalID))
				return
			}

			utils.WriteErrorResponse(w, http.StatusNotFound, fmt.Sprintf("depositer [%s] did not deposit on proposalID [%d]",
				bechDepositerAddr, proposalID))
			return
		}

		var deposit gov.Deposit
		cdc.MustUnmarshalBinary(res, &deposit)

		depositeResponse, err := govClient.ConvertDepositeCoins(cliCtx, deposit)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		output, err := wire.MarshalJSONIndent(cdc, depositeResponse)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Write(output)
	}
}

func queryVoteHandlerFn(cdc *wire.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]
		bechVoterAddr := vars[RestVoter]

		if len(strProposalID) == 0 {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "proposalId required but not specified")
			return
		}

		proposalID, err := strconv.ParseInt(strProposalID, 10, 64)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("proposalID [%s] is not positive", proposalID))
			return
		}

		if len(bechVoterAddr) == 0 {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "voter address required but not specified")
			return
		}

		voterAddr, err := sdk.AccAddressFromBech32(bechVoterAddr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("'%s' needs to be bech32 encoded", RestVoter))
			return
		}

		res, err := cliCtx.QueryStore(gov.KeyVote(proposalID, voterAddr), storeName)
		if err != nil || len(res) == 0 {
			res, err := cliCtx.QueryStore(gov.KeyProposal(proposalID), storeName)
			if err != nil || len(res) == 0 {
				utils.WriteErrorResponse(w, http.StatusNotFound, fmt.Sprintf("proposalID [%d] does not exist", proposalID))
				return
			}
			utils.WriteErrorResponse(w, http.StatusNotFound, fmt.Sprintf("voter [%s] did not vote on proposalID [%d]",
				bechVoterAddr, proposalID))
			return
		}

		var vote gov.Vote
		cdc.MustUnmarshalBinary(res, &vote)

		output, err := wire.MarshalJSONIndent(cdc, vote)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Write(output)
	}
}

// nolint: gocyclo
// todo: Split this functionality into helper functions to remove the above
func queryVotesOnProposalHandlerFn(cdc *wire.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]

		if len(strProposalID) == 0 {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "proposalId required but not specified")
			return
		}

		proposalID, err := strconv.ParseInt(strProposalID, 10, 64)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("proposalID [%s] is not positive", proposalID))
			return
		}

		res, err := cliCtx.QueryStore(gov.KeyProposal(proposalID), storeName)
		if err != nil || len(res) == 0 {
			utils.WriteErrorResponse(w, http.StatusNotFound, fmt.Sprintf("proposalID [%d] does not exist", proposalID))
			return
		}

		var proposal gov.Proposal
		cdc.MustUnmarshalBinary(res, &proposal)

		if proposal.GetStatus() != gov.StatusVotingPeriod {
			utils.WriteErrorResponse(w, http.StatusNotFound, fmt.Sprintf("proposal is not in Voting Period", proposalID))
			return
		}

		res2, err := cliCtx.QuerySubspace(gov.KeyVotesSubspace(proposalID), storeName)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusNotFound, "ProposalID doesn't exist")
			return
		}

		var votes []gov.Vote

		for i := 0; i < len(res2); i++ {
			var vote gov.Vote
			cdc.MustUnmarshalBinary(res2[i].Value, &vote)
			votes = append(votes, vote)
		}

		output, err := wire.MarshalJSONIndent(cdc, votes)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		w.Write(output)
	}
}

// nolint: gocyclo
// todo: Split this functionality into helper functions to remove the above
func queryProposalsWithParameterFn(cdc *wire.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bechVoterAddr := r.URL.Query().Get(RestVoter)
		bechDepositerAddr := r.URL.Query().Get(RestDepositer)
		strProposalStatus := r.URL.Query().Get(RestProposalStatus)

		var err error
		var voterAddr sdk.AccAddress
		var depositerAddr sdk.AccAddress
		var proposalStatus gov.ProposalStatus

		if len(bechVoterAddr) != 0 {
			voterAddr, err = sdk.AccAddressFromBech32(bechVoterAddr)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("'%s' needs to be bech32 encoded", RestVoter))
				return
			}
		}

		if len(bechDepositerAddr) != 0 {
			depositerAddr, err = sdk.AccAddressFromBech32(bechDepositerAddr)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("'%s' needs to be bech32 encoded", RestDepositer))
				return
			}
		}

		if len(strProposalStatus) != 0 {
			proposalStatus, err = gov.ProposalStatusFromString(strProposalStatus)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("'%s' is not a valid Proposal Status", strProposalStatus))
				return
			}
		}

		res, err := cliCtx.QueryStore(gov.KeyNextProposalID, storeName)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusNotFound, "no proposals exist yet and proposalID has not been set")
			return
		}

		var maxProposalID int64
		cdc.MustUnmarshalBinary(res, &maxProposalID)

		matchingProposals := []govClient.TextProposalResponse{}

		for proposalID := int64(0); proposalID < maxProposalID; proposalID++ {
			if voterAddr != nil {
				res, err = cliCtx.QueryStore(gov.KeyVote(proposalID, voterAddr), storeName)
				if err != nil || len(res) == 0 {
					continue
				}
			}

			if depositerAddr != nil {
				res, err = cliCtx.QueryStore(gov.KeyDeposit(proposalID, depositerAddr), storeName)
				if err != nil || len(res) == 0 {
					continue
				}
			}

			res, err = cliCtx.QueryStore(gov.KeyProposal(proposalID), storeName)
			if err != nil || len(res) == 0 {
				continue
			}

			var proposal gov.Proposal
			cdc.MustUnmarshalBinary(res, &proposal)

			if len(strProposalStatus) != 0 {
				if proposal.GetStatus() != proposalStatus {
					continue
				}
			}
			proposalResponse, err := govClient.ConvertProposalCoins(cliCtx, proposal)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			matchingProposals = append(matchingProposals, proposalResponse)
		}

		output, err := wire.MarshalJSONIndent(cdc, matchingProposals)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Write(output)
	}
}

// nolint: gocyclo
func queryConfigHandlerFn(cdc *wire.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := cliCtx.QuerySubspace([]byte(gov.Prefix), storeName)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		var kvs []govClient.KvPair
		for _, kv := range res {
			var v string
			cdc.UnmarshalBinary(kv.Value, &v)
			kv := govClient.KvPair{
				K: string(kv.Key),
				V: v,
			}
			kvs = append(kvs, kv)
		}
		output, err := wire.MarshalJSONIndent(cdc, kvs)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Write(output)
	}
}
