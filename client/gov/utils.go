package gov

import (
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/modules/params"
	"github.com/irisnet/irishub/modules/mint"
	"github.com/irisnet/irishub/modules/slashing"
	"github.com/irisnet/irishub/modules/service"
	"github.com/irisnet/irishub/modules/auth"
	"github.com/irisnet/irishub/modules/stake"
	"github.com/irisnet/irishub/app/v1/gov"
	distr "github.com/irisnet/irishub/modules/distribution"
)

var ParamSets = make(map[string]params.ParamSet)

func init() {
	params.RegisterParamSet(ParamSets, &mint.Params{}, &slashing.Params{}, &service.Params{}, &auth.Params{}, &stake.Params{}, &distr.Params{})
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
	return ""
}

//NormalizeProposalType - normalize user specified proposal type
func NormalizeProposalType(proposalType string) string {
	switch proposalType {
	case "ParameterChange", "parameter_change":
		return "ParameterChange"
	case "SoftwareUpgrade", "software_upgrade":
		return "SoftwareUpgrade"
	case "TxTaxUsage", "tx_tax_usage":
		return "TxTaxUsage"
	}
	return ""
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
	return ""
}

func ValidateParam(params gov.Params) error {
	for _, param := range params {
		if p, ok := ParamSets[param.Subspace]; ok {
			if _, err := p.Validate(param.Key, param.Value); err != nil {
				return err
			}
		} else {
			return gov.ErrInvalidParam(gov.DefaultCodespace, param.Subspace)
		}
	}
	return nil
}
