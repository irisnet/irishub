package types

import (
	"bytes"
)

const (
	separate = "::"
	//source type
	TransferFlow             = "Transfer"
	DelegationFlow           = "Delegation"
	UndelegationFlow         = "Undelegation"
	ValidatorRewardFlow      = "ValidatorReward"
	ValidatorCommissionFlow  = "ValidatorCommission"
	DelegatorRewardFlow      = "DelegatorReward"
	BurnFlow                 = "Burn"
	CommunityTaxUseFlow      = "CommunityTaxUse"
	GovDepositFlow           = "GovDeposit"
	GovDepositBurnFlow       = "GovDepositBurn"
	GovDepositRefundFlow     = "GovDepositRefund"
	ServiceDepositFlow       = "ServiceDeposit"
	ServiceDepositRefundFlow = "ServiceDepositRefund"

	//Trigger: transaction hash, module endBlock
	GovEndBlocker     = "govEndBlocker"
	SlashBeginBlocker = "slashBeginBlocker"
	SlashEndBlocker   = "slashEndBlocker"
	StakeEndBlocker   = "stakeEndBlocker"
	ServiceEndBlocker = "serviceEndBlocker"
)

type CoinFlowTags interface {
	GetTags() Tags
	//Append temporary tags to persistent tags
	TagWrite()
	//Clean temporary tags
	TagClean()
	//Add new tag to temporary tags
	AppendCoinFlowTag(ctx Context, from, to, amount, flowType, desc string)
}

type CoinFlowRecord struct {
	tags     Tags
	tempTags Tags
	enable   bool
}

func NewCoinFlowRecord(enable bool) CoinFlowTags {
	return &CoinFlowRecord{
		enable: enable,
	}
}

func (cfRecord *CoinFlowRecord) GetTags() Tags {
	return cfRecord.tags
}

func (cfRecord *CoinFlowRecord) AppendCoinFlowTag(ctx Context, from, to, amount, flowType, desc string) {
	if !cfRecord.enable {
		return
	}
	var tagKeyBuffer bytes.Buffer
	tagKeyBuffer.WriteString(ctx.CoinFlowTrigger())

	var tagValueBuffer bytes.Buffer
	tagValueBuffer.WriteString(from)
	tagValueBuffer.WriteString(separate)
	tagValueBuffer.WriteString(to)
	tagValueBuffer.WriteString(separate)
	tagValueBuffer.WriteString(amount)
	tagValueBuffer.WriteString(separate)
	tagValueBuffer.WriteString(flowType)
	tagValueBuffer.WriteString(separate)
	tagValueBuffer.WriteString(desc)
	tagValueBuffer.WriteString(separate)
	tagValueBuffer.WriteString(ctx.BlockHeader().Time.String())
	cfRecord.tempTags = append(cfRecord.tempTags, MakeTag(tagKeyBuffer.String(), []byte(tagValueBuffer.String())))
}

func (cfRecord *CoinFlowRecord) TagWrite() {
	cfRecord.tags = cfRecord.tags.AppendTags(cfRecord.tempTags)
	cfRecord.tempTags = nil
}

func (cfRecord *CoinFlowRecord) TagClean() {
	cfRecord.tempTags = nil
}
