package lcd

import (
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/modules/gov"
	"net/http"
	"strconv"
	"github.com/irisnet/irishub/client/gov/cli"
)

func queryProposalHandlerFn(cdc *wire.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]

		if len(strProposalID) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			err := errors.New("proposalId required but not specified")
			w.Write([]byte(err.Error()))

			return
		}

		proposalID, err := strconv.ParseInt(strProposalID, 10, 64)
		if err != nil {
			err := fmt.Errorf("proposalID [%d] is not positive", proposalID)
			w.Write([]byte(err.Error()))

			return
		}

		res, err := cliCtx.QueryStore(gov.KeyProposal(proposalID), storeName)
		if err != nil || len(res) == 0 {
			err := fmt.Errorf("proposalID [%d] does not exist", proposalID)
			w.Write([]byte(err.Error()))

			return
		}

		var proposal gov.Proposal
		cdc.MustUnmarshalBinary(res, &proposal)

		output, err := wire.MarshalJSONIndent(cdc, proposal)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))

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
			w.WriteHeader(http.StatusBadRequest)
			err := errors.New("proposalId required but not specified")
			w.Write([]byte(err.Error()))

			return
		}

		proposalID, err := strconv.ParseInt(strProposalID, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err := fmt.Errorf("proposalID [%d] is not positive", proposalID)
			w.Write([]byte(err.Error()))

			return
		}

		if len(bechDepositerAddr) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			err := errors.New("depositer address required but not specified")
			w.Write([]byte(err.Error()))

			return
		}

		depositerAddr, err := sdk.AccAddressFromBech32(bechDepositerAddr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err := fmt.Errorf("'%s' needs to be bech32 encoded", RestDepositer)
			w.Write([]byte(err.Error()))

			return
		}

		res, err := cliCtx.QueryStore(gov.KeyDeposit(proposalID, depositerAddr), storeName)
		if err != nil || len(res) == 0 {
			res, err := cliCtx.QueryStore(gov.KeyProposal(proposalID), storeName)
			if err != nil || len(res) == 0 {
				w.WriteHeader(http.StatusNotFound)
				err := fmt.Errorf("proposalID [%d] does not exist", proposalID)
				w.Write([]byte(err.Error()))

				return
			}

			w.WriteHeader(http.StatusNotFound)
			err = fmt.Errorf("depositer [%s] did not deposit on proposalID [%d]", bechDepositerAddr, proposalID)
			w.Write([]byte(err.Error()))

			return
		}

		var deposit gov.Deposit
		cdc.MustUnmarshalBinary(res, &deposit)

		output, err := wire.MarshalJSONIndent(cdc, deposit)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))

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
			w.WriteHeader(http.StatusBadRequest)
			err := errors.New("proposalId required but not specified")
			w.Write([]byte(err.Error()))
			return
		}

		proposalID, err := strconv.ParseInt(strProposalID, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err := fmt.Errorf("proposalID [%s] is not positive", proposalID)
			w.Write([]byte(err.Error()))

			return
		}

		if len(bechVoterAddr) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			err := errors.New("voter address required but not specified")
			w.Write([]byte(err.Error()))
			return
		}

		voterAddr, err := sdk.AccAddressFromBech32(bechVoterAddr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err := fmt.Errorf("'%s' needs to be bech32 encoded", RestVoter)
			w.Write([]byte(err.Error()))

			return
		}

		res, err := cliCtx.QueryStore(gov.KeyVote(proposalID, voterAddr), storeName)
		if err != nil || len(res) == 0 {
			res, err := cliCtx.QueryStore(gov.KeyProposal(proposalID), storeName)
			if err != nil || len(res) == 0 {
				w.WriteHeader(http.StatusNotFound)
				err := fmt.Errorf("proposalID [%d] does not exist", proposalID)
				w.Write([]byte(err.Error()))

				return
			}

			w.WriteHeader(http.StatusNotFound)
			err = fmt.Errorf("voter [%s] did not vote on proposalID [%d]", bechVoterAddr, proposalID)
			w.Write([]byte(err.Error()))

			return
		}

		var vote gov.Vote
		cdc.MustUnmarshalBinary(res, &vote)

		output, err := wire.MarshalJSONIndent(cdc, vote)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))

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
			w.WriteHeader(http.StatusBadRequest)
			err := errors.New("proposalId required but not specified")
			w.Write([]byte(err.Error()))

			return
		}

		proposalID, err := strconv.ParseInt(strProposalID, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err := fmt.Errorf("proposalID [%s] is not positive", proposalID)
			w.Write([]byte(err.Error()))

			return
		}

		res, err := cliCtx.QueryStore(gov.KeyProposal(proposalID), storeName)
		if err != nil || len(res) == 0 {
			err := fmt.Errorf("proposalID [%d] does not exist", proposalID)
			w.Write([]byte(err.Error()))

			return
		}

		var proposal gov.Proposal
		cdc.MustUnmarshalBinary(res, &proposal)

		if proposal.GetStatus() != gov.StatusVotingPeriod {
			err := fmt.Errorf("proposal is not in Voting Period", proposalID)
			w.Write([]byte(err.Error()))

			return
		}

		res2, err := cliCtx.QuerySubspace(gov.KeyVotesSubspace(proposalID), storeName)
		if err != nil {
			err = errors.New("ProposalID doesn't exist")
			w.Write([]byte(err.Error()))
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
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))

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
				w.WriteHeader(http.StatusBadRequest)
				err := fmt.Errorf("'%s' needs to be bech32 encoded", RestVoter)
				w.Write([]byte(err.Error()))
				return
			}
		}

		if len(bechDepositerAddr) != 0 {
			depositerAddr, err = sdk.AccAddressFromBech32(bechDepositerAddr)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				err := fmt.Errorf("'%s' needs to be bech32 encoded", RestDepositer)
				w.Write([]byte(err.Error()))

				return
			}
		}

		if len(strProposalStatus) != 0 {
			proposalStatus, err = gov.ProposalStatusFromString(strProposalStatus)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				err := fmt.Errorf("'%s' is not a valid Proposal Status", strProposalStatus)
				w.Write([]byte(err.Error()))

				return
			}
		}

		res, err := cliCtx.QueryStore(gov.KeyNextProposalID, storeName)
		if err != nil {
			err = errors.New("no proposals exist yet and proposalID has not been set")
			w.Write([]byte(err.Error()))

			return
		}

		var maxProposalID int64
		cdc.MustUnmarshalBinary(res, &maxProposalID)

		matchingProposals := []gov.Proposal{}

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

			matchingProposals = append(matchingProposals, proposal)
		}

		output, err := wire.MarshalJSONIndent(cdc, matchingProposals)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))

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
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))

			return
		}
		var kvs []cli.KvPair
		for _, kv := range res {
			var v string
			cdc.UnmarshalBinary(kv.Value, &v)
			kv := cli.KvPair{
				K: string(kv.Key),
				V: v,
			}
			kvs = append(kvs, kv)
		}
		output, err := wire.MarshalJSONIndent(cdc, kvs)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))

			return
		}

		w.Write(output)
	}
}