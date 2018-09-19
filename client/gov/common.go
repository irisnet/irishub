package gov

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/modules/gov"
)

// Deposit
type DepositResponse struct {
	Depositer  sdk.AccAddress `json:"depositer"`   //  Address of the depositer
	ProposalID int64          `json:"proposal_id"` //  proposalID of the proposal
	Amount     []string       `json:"amount"`      //  Deposit amount
}

type TextProposalResponse struct {
	ProposalID   int64            `json:"proposal_id"`   //  ID of the proposal
	Title        string           `json:"title"`         //  Title of the proposal
	Description  string           `json:"description"`   //  Description of the proposal
	ProposalType gov.ProposalKind `json:"proposal_type"` //  Type of proposal. Initial set {PlainTextProposal, SoftwareUpgradeProposal}

	Status      gov.ProposalStatus `json:"proposal_status"` //  Status of the Proposal {Pending, Active, Passed, Rejected}
	TallyResult gov.TallyResult    `json:"tally_result"`    //  Result of Tallys

	SubmitBlock  int64    `json:"submit_block"`  //  Height of the block where TxGovSubmitProposal was included
	TotalDeposit []string `json:"total_deposit"` //  Current deposit on this proposal. Initial value is set at InitialDeposit

	VotingStartBlock int64 `json:"voting_start_block"` //  Height of the block where MinDeposit was reached. -1 if MinDeposit is not reached
}

type KvPair struct {
	K string `json:"key"`
	V string `json:"value"`
}

func ConvertProposalCoins(cliCtx context.CLIContext, proposal gov.Proposal) (TextProposalResponse, error) {
	totalDeposit, err := cliCtx.ConvertCoinToMainUnit(proposal.GetTotalDeposit().String())
	if err != nil {
		return TextProposalResponse{}, err
	}
	return TextProposalResponse{
		ProposalID:   proposal.GetProposalID(),
		Title:        proposal.GetTitle(),
		Description:  proposal.GetDescription(),
		ProposalType: proposal.GetProposalType(),

		Status:      proposal.GetStatus(),
		TallyResult: proposal.GetTallyResult(),

		SubmitBlock:  proposal.GetSubmitBlock(),
		TotalDeposit: totalDeposit,

		VotingStartBlock: proposal.GetVotingStartBlock(),
	}, nil
}

func ConvertDepositeCoins(cliCtx context.CLIContext, deposite gov.Deposit) (DepositResponse, error) {
	amount, err := cliCtx.ConvertCoinToMainUnit(deposite.Amount.String())
	if err != nil {
		return DepositResponse{}, err
	}
	return DepositResponse{
		ProposalID: deposite.ProposalID,
		Depositer:  deposite.Depositer,
		Amount:     amount,
	}, nil
}
