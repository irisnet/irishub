package types

import (
	"bytes"
	"fmt"
	"strings"

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

	//Trigger: transaction hash, module endBlock
	GovEndBlocker     = "govEndBlocker"
	SlashBeginBlocker = "slashBeginBlocker"
	SlashEndBlocker   = "slashEndBlocker"
	StakeEndBlocker   = "stakeEndBlocker"
	ServiceEndBlocker = "serviceEndBlocker"
)

// ----------------------------------------------------------------------------
// Tags Manager
// ----------------------------------------------------------------------------

// TagsManager implements a simple wrapper around a slice of Tag objects that
// can be added from.
type TagsManager struct {
	tags              Tags
	isCoinFlowEnabled bool
	coinFlowTrigger   string
}

func NewTagsManager(enableCoinFlow bool) *TagsManager {
	return &TagsManager{EmptyTags(), enableCoinFlow, ""}
}

func (tm *TagsManager) Tags() Tags { return tm.tags }

// AddTag stores a single Tag object.
func (tm *TagsManager) AddTag(k string, v []byte) {
	tm.tags = tm.tags.AppendTag(k, v)
}

// AddTags stores a series of Tag objects.
func (tm *TagsManager) AddTags(tags Tags) {
	tm.tags = tm.tags.AppendTags(tags)
}

func (tm *TagsManager) SetCoinFlowTrigger(coinFlowTrigger string) {
	tm.coinFlowTrigger = coinFlowTrigger
}

func (tm *TagsManager) GetCoinFlowTrigger() string {
	return tm.coinFlowTrigger
}

// AddCoinFlow stores a single CoinFlow Tag object.
func (tm *TagsManager) AddCoinFlow(ctx Context, from, to, amount, flowType, desc string) {
	if !tm.isCoinFlowEnabled {
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

	tm.tags = append(tm.tags, MakeTag(tagKeyBuffer.String(), []byte(tagValueBuffer.String())))
}

// ----------------------------------------------------------------------------
// Tag
// ----------------------------------------------------------------------------

// Type synonym for convenience
type Tag = cmn.KVPair

// Make a tag from a key and a value
func MakeTag(k string, v []byte) Tag {
	return Tag{Key: []byte(k), Value: v}
}

// ----------------------------------------------------------------------------
// Tags
// ----------------------------------------------------------------------------

// Type synonym for convenience
type Tags cmn.KVPairs

// New empty tags
func EmptyTags() Tags {
	return make(Tags, 0)
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

func (t Tags) String() string {
	var sb strings.Builder

	for _, e := range t {
		sb.WriteString(fmt.Sprintf("%s=%s, ", e.Key, e.Value))
	}

	return strings.TrimRight(sb.String(), ", ")
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
