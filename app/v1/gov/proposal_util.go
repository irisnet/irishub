package gov

import (
	"github.com/irisnet/irishub/app/v1/asset/exported"
	sdk "github.com/irisnet/irishub/types"
)

type pTypeInfo struct {
	Type           ProposalKind
	Level          ProposalLevel
	createProposal func(content Content) Proposal
}

func createPlainTextInfo() pTypeInfo {
	return pTypeInfo{
		ProposalTypePlainText,
		ProposalLevelNormal,
		func(content Content) Proposal {
			return buildProposal(content, func(p BasicProposal, content Content) Proposal {
				return &PlainTextProposal{
					p,
				}
			})
		},
	}
}
func createParameterInfo() pTypeInfo {
	return pTypeInfo{
		ProposalTypeParameter,
		ProposalLevelImportant,
		func(content Content) Proposal {
			return buildProposal(content, func(p BasicProposal, content Content) Proposal {
				return &ParameterProposal{
					p,
					content.GetParams(),
				}
			})
		},
	}
}
func createSoftwareUpgradeInfo() pTypeInfo {
	return pTypeInfo{
		ProposalTypeSoftwareUpgrade,
		ProposalLevelCritical,
		func(content Content) Proposal {
			return buildProposal(content, func(p BasicProposal, content Content) Proposal {
				upgradeMsg := content.(MsgSubmitSoftwareUpgradeProposal)
				proposal := &SoftwareUpgradeProposal{
					p,
					sdk.ProtocolDefinition{
						Version:   upgradeMsg.Version,
						Software:  upgradeMsg.Software,
						Height:    upgradeMsg.SwitchHeight,
						Threshold: upgradeMsg.Threshold},
				}
				return proposal
			})
		},
	}
}

func createSystemHaltInfo() pTypeInfo {
	return pTypeInfo{
		ProposalTypeSystemHalt,
		ProposalLevelCritical,
		func(content Content) Proposal {
			return buildProposal(content, func(p BasicProposal, content Content) Proposal {
				return &SystemHaltProposal{
					p,
				}
			})
		},
	}
}

func createCommunityTaxUsageInfo() pTypeInfo {
	return pTypeInfo{
		ProposalTypeCommunityTaxUsage,
		ProposalLevelImportant,
		func(content Content) Proposal {
			return buildProposal(content, func(p BasicProposal, content Content) Proposal {
				taxMsg := content.(MsgSubmitCommunityTaxUsageProposal)
				proposal := &CommunityTaxUsageProposal{
					p,
					TaxUsage{
						taxMsg.Usage,
						taxMsg.DestAddress,
						taxMsg.Percent},
				}
				return proposal
			})
		},
	}
}

func createTokenAdditionInfo() pTypeInfo {
	return pTypeInfo{
		ProposalTypeTokenAddition,
		ProposalLevelImportant,
		func(content Content) Proposal {
			return buildProposal(content, func(p BasicProposal, content Content) Proposal {
				addTokenMsg := content.(MsgSubmitTokenAdditionProposal)
				decimal := int(addTokenMsg.Decimal)
				maxSupply := sdk.NewIntWithDecimal(int64(exported.MaximumAssetMaxSupply), decimal)

				fToken := exported.NewFungibleToken(exported.EXTERNAL, "", addTokenMsg.Symbol, addTokenMsg.Name, addTokenMsg.Decimal, addTokenMsg.CanonicalSymbol, addTokenMsg.MinUnitAlias, sdk.ZeroInt(), maxSupply, true, nil)
				proposal := &TokenAdditionProposal{
					p,
					fToken,
				}
				return proposal
			})
		},
	}
}

func buildProposal(content Content, callback func(p BasicProposal, content Content) Proposal) Proposal {
	var p = BasicProposal{
		Title:        content.GetTitle(),
		Description:  content.GetDescription(),
		ProposalType: content.GetProposalType(),
		Status:       StatusDepositPeriod,
		TallyResult:  EmptyTallyResult(),
		TotalDeposit: sdk.Coins{},
		Proposer:     content.GetProposer(),
	}
	return callback(p, content)
}
