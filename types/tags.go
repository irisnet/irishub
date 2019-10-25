package types

import (
	"bytes"

	cmn "github.com/tendermint/tendermint/libs/common"
)

const (
	separate = "::"
	//source type
	TransferFlow                = "Transfer"
	DelegationFlow              = "Delegation"
	UndelegationFlow            = "Undelegation"
	ValidatorRewardFlow         = "ValidatorReward"
	ValidatorCommissionFlow     = "ValidatorCommission"
	DelegatorRewardFlow         = "DelegatorReward"
	BurnFlow                    = "Burn"
	CommunityTaxCollectFlow     = "CommunityTaxCollect"
	CommunityTaxUseFlow         = "CommunityTaxUse"
	GovDepositFlow              = "GovDeposit"
	GovDepositBurnFlow          = "GovDepositBurn"
	GovDepositRefundFlow        = "GovDepositRefund"
	ServiceDepositFlow          = "ServiceDeposit"
	ServiceDepositRefundFlow    = "ServiceDepositRefund"
	MintTokenFlow               = "MintToken"
	IssueTokenFlow              = "IssueToken"
	CoinSwapInputFlow           = "CoinSwapInput"
	CoinSwapOutputFlow          = "CoinSwapOutput"
	CoinSwapAddLiquidityFlow    = "AddLiquidity"
	CoinSwapRemoveLiquidityFlow = "RemoveLiquidity"
	CoinHTLCCreateFlow          = "CreateHTLC"
	CoinHTLCClaimFlow           = "ClaimHTLC"
	CoinHTLCRefundFlow          = "RefundHTLC"

	//Trigger: transaction hash, module endBlock and beginBlock
	GovEndBlocker            = "govEndBlocker"
	SlashBeginBlocker        = "slashBeginBlocker"
	SlashEndBlocker          = "slashEndBlocker"
	StakeEndBlocker          = "stakeEndBlocker"
	ServiceEndBlocker        = "serviceEndBlocker"
	DistributionBeginBlocker = "distributionBeginBlocker"
)

// ----------------------------------------------------------------------------
// CoinFlowTags
// ----------------------------------------------------------------------------

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

// ----------------------------------------------------------------------------
// Tags
// ----------------------------------------------------------------------------

// Type synonym for convenience
type Tag = cmn.KVPair

// Type synonym for convenience
type Tags cmn.KVPairs

// New empty tags
func EmptyTags() Tags {
	return make(Tags, 0)
}

// Append a single tag
func (t Tags) AppendTag(k string, v []byte) Tags {
	return append(t, MakeTag(k, v))
}

// Append two lists of tags
func (t Tags) AppendTags(tags Tags) Tags {
	return append(t, tags...)
}

// Turn tags into KVPair list
func (t Tags) ToKVPairs() []cmn.KVPair {
	return []cmn.KVPair(t)
}

// New variadic tags, must be k string, v []byte repeating
func NewTags(tags ...interface{}) Tags {
	var ret Tags
	if len(tags)%2 != 0 {
		panic("must specify key-value pairs as varargs")
	}
	i := 0
	for {
		if i == len(tags) {
			break
		}
		ret = append(ret, Tag{Key: []byte(tags[i].(string)), Value: tags[i+1].([]byte)})
		i += 2
	}
	return ret
}

// Make a tag from a key and a value
func MakeTag(k string, v []byte) Tag {
	return Tag{Key: []byte(k), Value: v}
}

//__________________________________________________

// common tags
var (
	TagAction              = "action"
	TagSrcValidator        = "source-validator"
	TagDstValidator        = "destination-validator"
	TagDelegator           = "delegator"
	TagReward              = "withdraw-reward-total"
	TagWithdrawAddr        = "withdraw-address"
	TagRewardFromValidator = "withdraw-reward-from-validator-%s"
	TagRewardCommission    = "withdraw-reward-commission"
)
