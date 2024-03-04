package v300

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	// ValidatorBondFactor dictates the cap on the liquid shares
	// for a validator - determined as a multiple to their validator bond
	// (e.g. ValidatorBondShares = 1000, BondFactor = 250 -> LiquidSharesCap: 250,000)
	ValidatorBondFactor = sdk.NewDec(250)
	// ValidatorLiquidStakingCap represents a cap on the portion of stake that
	// comes from liquid staking providers for a specific validator
	ValidatorLiquidStakingCap = sdk.MustNewDecFromStr("0.5") // 50%
	// GlobalLiquidStakingCap represents the percentage cap on
	// the portion of a chain's total stake can be liquid
	GlobalLiquidStakingCap = sdk.MustNewDecFromStr("0.25") // 25%

	allowMessages = []string{
		"/cosmos.authz.v1beta1.MsgExec",
		"/cosmos.authz.v1beta1.MsgGrant",
		"/cosmos.authz.v1beta1.MsgRevoke",
		"/cosmos.bank.v1beta1.MsgSend",
		"/cosmos.bank.v1beta1.MsgMultiSend",
		"/cosmos.distribution.v1beta1.MsgSetWithdrawAddress",
		"/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission",
		"/cosmos.distribution.v1beta1.MsgFundCommunityPool",
		"/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward",
		"/cosmos.feegrant.v1beta1.MsgGrantAllowance",
		"/cosmos.feegrant.v1beta1.MsgRevokeAllowance",
		"/cosmos.gov.v1beta1.MsgVoteWeighted",
		"/cosmos.gov.v1beta1.MsgSubmitProposal",
		"/cosmos.gov.v1beta1.MsgDeposit",
		"/cosmos.gov.v1beta1.MsgVote",
		"/cosmos.gov.v1.MsgVoteWeighted",
		"/cosmos.gov.v1.MsgSubmitProposal",
		"/cosmos.gov.v1.MsgDeposit",
		"/cosmos.gov.v1.MsgVote",
		"/cosmos.staking.v1beta1.MsgEditValidator",
		"/cosmos.staking.v1beta1.MsgDelegate",
		"/cosmos.staking.v1beta1.MsgUndelegate",
		"/cosmos.staking.v1beta1.MsgBeginRedelegate",
		"/cosmos.staking.v1beta1.MsgCreateValidator",
		"/cosmos.vesting.v1beta1.MsgCreateVestingAccount",
		"/ibc.applications.transfer.v1.MsgTransfer",
		"/irismod.nft.MsgIssueDenom",
		"/irismod.nft.MsgTransferDenom",
		"/irismod.nft.MsgMintNFT",
		"/irismod.nft.MsgEditNFT",
		"/irismod.nft.MsgTransferNFT",
		"/irismod.nft.MsgBurnNFT",
		"/irismod.mt.MsgIssueDenom",
		"/irismod.mt.MsgTransferDenom",
		"/irismod.mt.MsgMintMT",
		"/irismod.mt.MsgEditMT",
		"/irismod.mt.MsgTransferMT",
		"/irismod.mt.MsgBurnMT",
	}
)
