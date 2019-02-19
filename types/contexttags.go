package types

import "bytes"

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
	StakeDelegation  = "stakeDelegation"

	//endblock
	GovEndBlocker = "govEndBlocker"
	SlashEndBlocker = "slashEndBlocker"
	StakeEndBlocker = "stakeEndBlocker"
	UpgradeEndBlocker = "upgradeEndBlocker"
	ServiceEndBlocker = "serviceEndBlocker"
)

type CoinFlowRecord interface {
	Append(key, value string)
	GetTags() Tags
	AppendAddCoinTag(trigger, recipient, msgType, amount, timestamp string)
	AppendSubtractCoinTag(trigger, sender, msgType, amount, timestamp string)
	AppendAddCoinSourceTag(trigger, recipient, msgType, sourceType, source, amount, timestamp string)
}

type CoinFlowTags struct {
	tags Tags
}

func NewCoinFlowTags() CoinFlowRecord {
	return &CoinFlowTags{}
}

func (cfTag *CoinFlowTags) Append(key, value string) {
	cfTag.tags = append(cfTag.tags, MakeTag(key, []byte(value)))
}

func (cfTag *CoinFlowTags) GetTags() Tags {
	return cfTag.tags
}

func (cfTag *CoinFlowTags) AppendAddCoinTag(trigger, recipient, msgType, amount, timestamp string) {
	var tagKeyBuffer bytes.Buffer
	tagKeyBuffer.WriteString(trigger)

	var tagValueBuffer bytes.Buffer
	tagValueBuffer.WriteString("add")
	tagValueBuffer.WriteString("&")
	tagValueBuffer.WriteString(recipient)
	tagValueBuffer.WriteString("&")
	tagValueBuffer.WriteString(msgType)
	tagValueBuffer.WriteString("&")
	tagValueBuffer.WriteString(amount)
	tagValueBuffer.WriteString("&")
	tagValueBuffer.WriteString(timestamp)
	cfTag.tags = append(cfTag.tags, MakeTag(tagKeyBuffer.String(), []byte(tagValueBuffer.String())))
}

func (cfTag *CoinFlowTags) AppendSubtractCoinTag(trigger, sender, msgType, amount, timestamp string) {
	var tagKeyBuffer bytes.Buffer
	tagKeyBuffer.WriteString(trigger)

	var tagValueBuffer bytes.Buffer
	tagValueBuffer.WriteString("subtract")
	tagValueBuffer.WriteString("&")
	tagValueBuffer.WriteString(sender)
	tagValueBuffer.WriteString("&")
	tagValueBuffer.WriteString(msgType)
	tagValueBuffer.WriteString("&")
	tagValueBuffer.WriteString(amount)
	tagValueBuffer.WriteString("&")
	tagValueBuffer.WriteString(timestamp)
	cfTag.tags = append(cfTag.tags, MakeTag(tagKeyBuffer.String(), []byte(tagValueBuffer.String())))
}

func (cfTag *CoinFlowTags) AppendAddCoinSourceTag(trigger, recipient, msgType, sourceType, source, amount, timestamp string) {
	var tagKeyBuffer bytes.Buffer
	tagKeyBuffer.WriteString(trigger)

	var tagValueBuffer bytes.Buffer
	tagValueBuffer.WriteString("source")
	tagValueBuffer.WriteString("&")
	tagValueBuffer.WriteString(recipient)
	tagValueBuffer.WriteString("&")
	tagValueBuffer.WriteString(msgType)
	tagValueBuffer.WriteString("&")
	tagValueBuffer.WriteString(sourceType)
	tagValueBuffer.WriteString("&")
	tagValueBuffer.WriteString(source)
	tagValueBuffer.WriteString("&")
	tagValueBuffer.WriteString(amount)
	tagValueBuffer.WriteString("&")
	tagValueBuffer.WriteString(timestamp)
	cfTag.tags = append(cfTag.tags, MakeTag(tagKeyBuffer.String(), []byte(tagValueBuffer.String())))
}
