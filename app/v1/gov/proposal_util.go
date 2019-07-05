package gov

import (
	"github.com/irisnet/irishub/app/v1/asset"
	sdk "github.com/irisnet/irishub/types"
)

type pTypeInfo struct {
	Type           ProposalKind
	Level          ProposalLevel
	createProposal func(content Context) Proposal
}

func createPlainTextInfo() pTypeInfo {
	return pTypeInfo{
		ProposalTypePlainText,
		ProposalLevelNormal,
		func(content Context) Proposal {
			return buildProposal(content, func(p BasicProposal, content Context) Proposal {
				return &PlainTextProposal{
					p,
				}
			})
		},
	}
}
func createParameterChangeInfo() pTypeInfo {
	return pTypeInfo{
		ProposalTypeParameterChange,
		ProposalLevelImportant,
		func(content Context) Proposal {
			return buildProposal(content, func(p BasicProposal, content Context) Proposal {
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
		func(content Context) Proposal {
			return buildProposal(content, func(p BasicProposal, content Context) Proposal {
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
		func(content Context) Proposal {
			return buildProposal(content, func(p BasicProposal, content Context) Proposal {
				return &SystemHaltProposal{
					p,
				}
			})
		},
	}
}

func createTxTaxUsageInfo() pTypeInfo {
	return pTypeInfo{
		ProposalTypeCommunityTaxUsage,
		ProposalLevelImportant,
		func(content Context) Proposal {
			return buildProposal(content, func(p BasicProposal, content Context) Proposal {
				taxMsg := content.(MsgSubmitTxTaxUsageProposal)
				proposal := &TaxUsageProposal{
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

func createAddTokenInfo() pTypeInfo {
	return pTypeInfo{
		ProposalTypeTokenAddition,
		ProposalLevelImportant,
		func(content Context) Proposal {
			return buildProposal(content, func(p BasicProposal, content Context) Proposal {
				addTokenMsg := content.(MsgSubmitAddTokenProposal)
				decimal := int(addTokenMsg.Decimal)
				initialSupply := sdk.NewIntWithDecimal(int64(addTokenMsg.InitialSupply), decimal)
				maxSupply := sdk.NewIntWithDecimal(int64(asset.MaximumAssetMaxSupply), decimal)

				fToken := asset.NewFungibleToken(asset.EXTERNAL, "", addTokenMsg.Symbol, addTokenMsg.Name, addTokenMsg.Decimal, addTokenMsg.SymbolAtSource, addTokenMsg.SymbolMinAlias, initialSupply, maxSupply, false, nil)
				proposal := &AddTokenProposal{
					p,
					fToken,
				}
				return proposal
			})
		},
	}
}

func buildProposal(content Context, callback func(p BasicProposal, content Context) Proposal) Proposal {
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
