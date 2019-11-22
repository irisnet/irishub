package lcd

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/app/v1/gov"
	"github.com/irisnet/irishub/client/context"
	client "github.com/irisnet/irishub/client/gov"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

type postProposalReq struct {
	BaseTx         utils.BaseTx   `json:"base_tx"`
	Title          string         `json:"title"`           //  Title of the proposal
	Description    string         `json:"description"`     //  Description of the proposal
	ProposalType   string         `json:"proposal_type"`   //  Type of proposal. Initial set {PlainTextProposal, SoftwareUpgradeProposal}
	Proposer       sdk.AccAddress `json:"proposer"`        //  Address of the proposer
	InitialDeposit string         `json:"initial_deposit"` // Coins to add to the proposal's deposit
	Param          gov.Param      `json:"param"`
	CommTax        commTax        `json:"comm_tax"`
	Token          token          `json:"token"`
	Upgrade        upgrade        `json:"upgrade"`
}

type token struct {
	Symbol          string `json:"symbol"`
	CanonicalSymbol string `json:"canonical_symbol"`
	Name            string `json:"name"`
	Decimal         uint8  `json:"decimal"`
	MinUnitAlias    string `json:"min_unit_alias"`
}

type upgrade struct {
	Version      uint64  `json:"version"`
	Software     string  `json:"software"`
	SwitchHeight uint64  `json:"switch_height"`
	Threshold    sdk.Dec `json:"threshold"`
}

type commTax struct {
	Usage       gov.UsageType  `json:"usage"`
	DestAddress sdk.AccAddress `json:"dest_address"`
	Percent     sdk.Dec        `json:"percent"`
}

type depositReq struct {
	BaseTx    utils.BaseTx   `json:"base_tx"`
	Depositor sdk.AccAddress `json:"depositor"` // Address of the depositor
	Amount    string         `json:"amount"`    // Coins to add to the proposal's deposit
}

type voteReq struct {
	BaseTx utils.BaseTx   `json:"base_tx"`
	Voter  sdk.AccAddress `json:"voter"`  //  address of the voter
	Option string         `json:"option"` //  option from OptionSet chosen by the voter
}

func postProposalHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req postProposalReq
		err := utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		proposalType, err := gov.ProposalTypeFromString(client.NormalizeProposalType(req.ProposalType))
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if proposalType == gov.ProposalTypeParameter {
			if err := client.ValidateParam(req.Param); err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		initDepositAmount, err := cliCtx.ParseCoins(req.InitialDeposit)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		// create the message
		msgs := make([]sdk.Msg, 1)
		msg := gov.NewMsgSubmitProposal(req.Title, req.Description, proposalType, req.Proposer, initDepositAmount, gov.Params{req.Param})
		switch msg.ProposalType {
		case gov.ProposalTypeParameter, gov.ProposalTypeSystemHalt, gov.ProposalTypePlainText:
			msgs[0] = msg
			break
		case gov.ProposalTypeSoftwareUpgrade:
			msgs[0] = gov.NewMsgSubmitSoftwareUpgradeProposal(msg, req.Upgrade.Version, req.Upgrade.Software, req.Upgrade.SwitchHeight, req.Upgrade.Threshold)
			break
		case gov.ProposalTypeCommunityTaxUsage:
			msgs[0] = gov.NewMsgSubmitCommunityTaxUsageProposal(msg, req.CommTax.Usage, req.CommTax.DestAddress, req.CommTax.Percent)
			break
		case gov.ProposalTypeTokenAddition:
			msgs[0] = gov.NewMsgSubmitTokenAdditionProposal(msg, req.Token.Symbol, req.Token.CanonicalSymbol, req.Token.Name, req.Token.MinUnitAlias, req.Token.Decimal)
			break
		default:
			utils.WriteErrorResponse(w, http.StatusBadRequest, "not a valid proposal type")
			return
		}

		err = msgs[0].ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.WriteGenerateStdTxResponse(w, txCtx, msgs)
	}
}

func depositHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]

		if len(strProposalID) == 0 {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "proposalId required but not specified")
			return
		}

		proposalID, ok := utils.ParseUint64OrReturnBadRequest(w, strProposalID)
		if !ok {
			return
		}

		var req depositReq
		err := utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		depositAmount, err := cliCtx.ParseCoins(req.Amount)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := gov.NewMsgDeposit(req.Depositor, proposalID, depositAmount)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}

func voteHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
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

		var req voteReq
		err := utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		voteOption, err := gov.VoteOptionFromString(client.NormalizeVoteOption(req.Option))
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := gov.NewMsgVote(req.Voter, proposalID, voteOption)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}
