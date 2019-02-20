package types

import (
	"bytes"
	"strconv"
)

const (
	//source type
	ValidatorDelegationReward = "validatorDelegationReward"
	ValidatorCommissionReward = "validatorCommissionReward"
	StakeDelegationRefund     = "stakeDelegationRefund"
	TokenTransfer             = "tokenTransfer"
	GovDeposit                = "govDeposit"
	GovDepositRefund          = "govDepositRefund"
	ServiceDeposit            = "serviceDeposit"
	ServiceDepositRefund      = "serviceDepositRefund"
	CommunityTax              = "communityTax"
	//source name
	TokenTransferTx  = "tokenTransferTransaction"
	CommunityTaxPool = "communityTaxPool"

	//trigger: transaction, endBlock
	EndBlockTrigger = "endBlocker"
	//Msg type: transaction msg type, module endBlock and txFee
	GovEndBlocker = "govEndBlocker"
	SlashEndBlocker = "slashEndBlocker"
	StakeEndBlocker = "stakeEndBlocker"
	ServiceEndBlocker = "serviceEndBlocker"
	TxFee = "txFee"
)

type CoinFlowTags interface {
	Append(key, value string)
	GetTags() Tags
	AppendAddCoinTag(ctx Context, recipient, amount string)
	AppendSubtractCoinTag(ctx Context, sender, amount string)
	AppendAddCoinSourceTag(ctx Context, recipient, sourceType, source, amount string)
}

type CoinFlowRecord struct {
	tags   Tags
	enable bool
}

func NewCoinFlowRecord(enable bool) CoinFlowTags {
	return &CoinFlowRecord{
		enable: enable,
	}
}

func (cfRecord *CoinFlowRecord) Append(key, value string) {
	if !cfRecord.enable {
		return
	}
	cfRecord.tags = append(cfRecord.tags, MakeTag(key, []byte(value)))
}

func (cfRecord *CoinFlowRecord) GetTags() Tags {
	return cfRecord.tags
}

func (cfRecord *CoinFlowRecord) AppendAddCoinTag(ctx Context, recipient, amount string) {
	if !cfRecord.enable {
		return
	}
	var tagKeyBuffer bytes.Buffer
	tagKeyBuffer.WriteString(ctx.CoinFlowTrigger())

	var tagValueBuffer bytes.Buffer
	tagValueBuffer.WriteString("add")
	tagValueBuffer.WriteString("&")
	tagValueBuffer.WriteString(recipient)
	tagValueBuffer.WriteString("&")
	tagValueBuffer.WriteString(ctx.CoinFlowMsgType())
	tagValueBuffer.WriteString("&")
	tagValueBuffer.WriteString(amount)
	tagValueBuffer.WriteString("&")
	tagValueBuffer.WriteString(strconv.Itoa(int(ctx.BlockHeight())))
	tagValueBuffer.WriteString("&")
	tagValueBuffer.WriteString(ctx.BlockHeader().Time.String())
	cfRecord.tags = append(cfRecord.tags, MakeTag(tagKeyBuffer.String(), []byte(tagValueBuffer.String())))
}

func (cfRecord *CoinFlowRecord) AppendSubtractCoinTag(ctx Context, sender, amount string) {
	if !cfRecord.enable {
		return
	}
	var tagKeyBuffer bytes.Buffer
	tagKeyBuffer.WriteString(ctx.CoinFlowTrigger())

	var tagValueBuffer bytes.Buffer
	tagValueBuffer.WriteString("subtract")
	tagValueBuffer.WriteString("&")
	tagValueBuffer.WriteString(sender)
	tagValueBuffer.WriteString("&")
	tagValueBuffer.WriteString(ctx.CoinFlowMsgType())
	tagValueBuffer.WriteString("&")
	tagValueBuffer.WriteString(amount)
	tagValueBuffer.WriteString("&")
	tagValueBuffer.WriteString(strconv.Itoa(int(ctx.BlockHeight())))
	tagValueBuffer.WriteString("&")
	tagValueBuffer.WriteString(ctx.BlockHeader().Time.String())
	cfRecord.tags = append(cfRecord.tags, MakeTag(tagKeyBuffer.String(), []byte(tagValueBuffer.String())))
}

func (cfRecord *CoinFlowRecord) AppendAddCoinSourceTag(ctx Context, recipient, sourceType, source, amount string) {
	if !cfRecord.enable {
		return
	}
	var tagKeyBuffer bytes.Buffer
	tagKeyBuffer.WriteString(ctx.CoinFlowTrigger())

	var tagValueBuffer bytes.Buffer
	tagValueBuffer.WriteString("source")
	tagValueBuffer.WriteString("&")
	tagValueBuffer.WriteString(recipient)
	tagValueBuffer.WriteString("&")
	tagValueBuffer.WriteString(ctx.CoinFlowMsgType())
	tagValueBuffer.WriteString("&")
	tagValueBuffer.WriteString(sourceType)
	tagValueBuffer.WriteString("&")
	tagValueBuffer.WriteString(source)
	tagValueBuffer.WriteString("&")
	tagValueBuffer.WriteString(amount)
	tagValueBuffer.WriteString("&")
	tagValueBuffer.WriteString(strconv.Itoa(int(ctx.BlockHeight())))
	tagValueBuffer.WriteString("&")
	tagValueBuffer.WriteString(ctx.BlockHeader().Time.String())
	cfRecord.tags = append(cfRecord.tags, MakeTag(tagKeyBuffer.String(), []byte(tagValueBuffer.String())))
}
