package lcd

import (
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/modules/gov"
	"net/http"
	"strconv"
)

type postProposalReq struct {
	BaseTx         context.BaseTx   `json:"base_tx"`
	Param          gov.Param        `json:"param"`
	Title          string           `json:"title"`           //  Title of the proposal
	Description    string           `json:"description"`     //  Description of the proposal
	ProposalType   gov.ProposalKind `json:"proposal_type"`   //  Type of proposal. Initial set {PlainTextProposal, SoftwareUpgradeProposal}
	Proposer       string           `json:"proposer"`        //  Address of the proposer
	InitialDeposit string           `json:"initial_deposit"` // Coins to add to the proposal's deposit
}

type depositReq struct {
	BaseTx    context.BaseTx `json:"base_tx"`
	Depositer string         `json:"depositer"` // Address of the depositer
	Amount    string         `json:"amount"`    // Coins to add to the proposal's deposit
}

type voteReq struct {
	BaseTx context.BaseTx `json:"base_tx"`
	Voter  string         `json:"voter"`  //  address of the voter
	Option gov.VoteOption `json:"option"` //  option from OptionSet chosen by the voter
}

func postProposalHandlerFn(cdc *wire.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req postProposalReq
		err := utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}
		cliCtx = utils.InitRequestClictx(cliCtx, r, req.BaseTx.LocalAccountName, req.Proposer)
		txCtx, err := context.NewTxContextFromBaseTx(cliCtx, cdc, req.BaseTx)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		proposer, err := sdk.AccAddressFromBech32(req.Proposer)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		initDepositAmount, err := cliCtx.ParseCoins(req.InitialDeposit)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		// create the message
		msg := gov.NewMsgSubmitProposal(req.Title, req.Description, req.ProposalType, proposer, initDepositAmount, req.Param)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.SendOrReturnUnsignedTx(w, cliCtx, txCtx, req.BaseTx, []sdk.Msg{msg})
	}
}

func depositHandlerFn(cdc *wire.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]

		if len(strProposalID) == 0 {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "proposalId required but not specified")
			return
		}

		proposalID, err := strconv.ParseInt(strProposalID, 10, 64)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req depositReq
		err = utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		cliCtx = utils.InitRequestClictx(cliCtx, r, req.BaseTx.LocalAccountName, req.Depositer)
		txCtx, err := context.NewTxContextFromBaseTx(cliCtx, cdc, req.BaseTx)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		depositer, err := sdk.AccAddressFromBech32(req.Depositer)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		depositAmount, err := cliCtx.ParseCoins(req.Amount)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		// create the message
		msg := gov.NewMsgDeposit(depositer, proposalID, depositAmount)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.SendOrReturnUnsignedTx(w, cliCtx, txCtx, req.BaseTx, []sdk.Msg{msg})
	}
}

func voteHandlerFn(cdc *wire.Codec, cliCtx context.CLIContext) http.HandlerFunc {
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

		var req voteReq
		err = utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		cliCtx = utils.InitRequestClictx(cliCtx, r, req.BaseTx.LocalAccountName, req.Voter)
		txCtx, err := context.NewTxContextFromBaseTx(cliCtx, cdc, req.BaseTx)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		voter, err := sdk.AccAddressFromBech32(req.Voter)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		// create the message
		msg := gov.NewMsgVote(voter, proposalID, req.Option)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.SendOrReturnUnsignedTx(w, cliCtx, txCtx, req.BaseTx, []sdk.Msg{msg})
	}
}
