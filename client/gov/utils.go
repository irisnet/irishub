package gov

import (
	"github.com/irisnet/irishub/app/v1/asset"
	"github.com/irisnet/irishub/app/v1/auth"
	distr "github.com/irisnet/irishub/app/v1/distribution"
	"github.com/irisnet/irishub/app/v1/gov"
	"github.com/irisnet/irishub/app/v1/mint"
	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/app/v1/service"
	"github.com/irisnet/irishub/app/v1/slashing"
	"github.com/irisnet/irishub/app/v1/stake"
	"github.com/irisnet/irishub/app/v2/coinswap"
	sdk "github.com/irisnet/irishub/types"
)

var ParamSets = make(map[string]params.ParamSet)

func init() {
	params.RegisterParamSet(ParamSets, &mint.Params{}, &slashing.Params{}, &service.Params{}, &auth.Params{}, &stake.Params{}, &distr.Params{}, &asset.Params{}, &gov.GovParams{}, &coinswap.Params{})
}

// Deposit
type DepositOutput struct {
	Depositor  sdk.AccAddress `json:"depositor"`   //  Address of the depositor
	ProposalID int64          `json:"proposal_id"` //  proposalID of the proposal
	Amount     []string       `json:"amount"`      //  Deposit amount
}

type KvPair struct {
	K string `json:"key"`
	V string `json:"value"`
}

// NormalizeVoteOption - normalize user specified vote option
func NormalizeVoteOption(option string) string {
	switch option {
	case "Yes", "yes":
		return "Yes"
	case "Abstain", "abstain":
		return "Abstain"
	case "No", "no":
		return "No"
	case "NoWithVeto", "no_with_veto":
		return "NoWithVeto"
	}
	return option
}

//NormalizeProposalType - normalize user specified proposal type
func NormalizeProposalType(proposalType string) string {
	switch proposalType {
	case "Parameter", "parameter":
		return "Parameter"
	case "SoftwareUpgrade", "software_upgrade":
		return "SoftwareUpgrade"
	case "SystemHalt", "system_halt":
		return "SystemHalt"
	case "CommunityTaxUsage", "community_tax_usage":
		return "CommunityTaxUsage"
	case "TokenAddition", "token_addition":
		return "TokenAddition"
	}
	return proposalType
}

//NormalizeProposalStatus - normalize user specified proposal status
func NormalizeProposalStatus(status string) string {
	switch status {
	case "DepositPeriod", "deposit_period":
		return "DepositPeriod"
	case "VotingPeriod", "voting_period":
		return "VotingPeriod"
	case "Passed", "passed":
		return "Passed"
	case "Rejected", "rejected":
		return "Rejected"
	}
	return status
}

func ValidateParam(param gov.Param) error {
	if p, ok := ParamSets[param.Subspace]; ok {
		if p.ReadOnly() {
			return gov.ErrInvalidParam(gov.DefaultCodespace, param.Subspace)
		}
		if _, err := p.Validate(param.Key, param.Value); err != nil {
			return err
		}
	} else {
		return gov.ErrInvalidParam(gov.DefaultCodespace, param.Subspace)
	}
	return nil
}
