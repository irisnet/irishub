package gov

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/modules/gov"
	"time"
)

// Deposit
type DepositOutput struct {
	Depositer  sdk.AccAddress `json:"depositer"`   //  Address of the depositer
	ProposalID int64          `json:"proposal_id"` //  proposalID of the proposal
	Amount     []string       `json:"amount"`      //  Deposit amount
}

type ProposalOutput struct {
	ProposalID   int64            `json:"proposal_id"`   //  ID of the proposal
	Title        string           `json:"title"`         //  Title of the proposal
	Description  string           `json:"description"`   //  Description of the proposal
	ProposalType gov.ProposalKind `json:"proposal_type"` //  Type of proposal. Initial set {PlainTextProposal, SoftwareUpgradeProposal}

	Status      gov.ProposalStatus `json:"proposal_status"` //  Status of the Proposal {Pending, Active, Passed, Rejected}
	TallyResult gov.TallyResult    `json:"tally_result"`    //  Result of Tallys

	SubmitTime   time.Time `json:"submit_time"`   //  Height of the block where TxGovSubmitProposal was included
	TotalDeposit []string `json:"total_deposit"` //  Current deposit on this proposal. Initial value is set at InitialDeposit

	VotingStartTime time.Time `json:"voting_start_time"` //  Height of the block where MinDeposit was reached. -1 if MinDeposit is not reached

	Param gov.Param `json:"param"`
}

type KvPair struct {
	K string `json:"key"`
	V string `json:"value"`
}

func ConvertProposalToProposalOutput(cliCtx context.CLIContext, proposal gov.Proposal) (ProposalOutput, error) {
	totalDeposit, err := cliCtx.ConvertCoinToMainUnit(proposal.GetTotalDeposit().String())
	if err != nil {
		return ProposalOutput{}, err
	}

	proposalOutput := ProposalOutput{
		ProposalID:   proposal.GetProposalID(),
		Title:        proposal.GetTitle(),
		Description:  proposal.GetDescription(),
		ProposalType: proposal.GetProposalType(),

		Status:      proposal.GetStatus(),
		TallyResult: proposal.GetTallyResult(),

		SubmitTime:  proposal.GetSubmitTime(),
		TotalDeposit: totalDeposit,

		VotingStartTime: proposal.GetVotingStartTime(),
		Param:            gov.Param{},
	}

	if proposal.GetProposalType() == gov.ProposalTypeParameterChange {
		proposalOutput.Param = proposal.(*gov.ParameterProposal).Param
	}

	return proposalOutput, nil
}

func ConvertDepositToDepositOutput(cliCtx context.CLIContext, deposite gov.Deposit) (DepositOutput, error) {
	amount, err := cliCtx.ConvertCoinToMainUnit(deposite.Amount.String())
	if err != nil {
		return DepositOutput{}, err
	}
	return DepositOutput{
		ProposalID: deposite.ProposalID,
		Depositer:  deposite.Depositer,
		Amount:     amount,
	}, nil
}
