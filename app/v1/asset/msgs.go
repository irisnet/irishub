package asset

import (
	sdk "github.com/irisnet/irishub/types"
	"math"
	"regexp"
)

const (
	// name to identify transaction types
	MsgRoute = "asset"

	MaxInitSupply  = 1e+12
	MaxTotalSupply = math.MaxUint64
	MaxDecimal     = 18
)

var (
	// 00 - fungible; 01 - non-fungible
	MsgIssueFamily = map[string]bool{"00": true, "01": true}
	// Reserved - 00 (native); 01 (external); Gateway IDs
	MsgIssueSource = map[string]bool{"00": true, "01": true}
)

// MsgIssueAsset
type MsgIssueAsset struct {
	Family     string           `json:"family"`
	Name       string           `json:"name"`
	Symbol     string           `json:"symbol"`
	Source     string           `json:"source"`
	InitSupply uint64           `json:"init_supply"`
	MaxSupply  uint64           `json:"max_supply"`
	Decimal    uint8            `json:"decimal"`
	Mintable   bool             `json:"mintable"`
	Owner      sdk.AccAddress   `json:owner`
	Operators  []sdk.AccAddress `json:"operators"`
}

var _ sdk.Msg = MsgIssueAsset{}

// NewMsgIssue - construct asset issue msg.
func NewMsgIssue(family string, name string, symbol string, source string, initSupply uint64, maxSupply uint64, decimal uint8, mintable bool, owner sdk.AccAddress, operators []sdk.AccAddress) MsgIssueAsset {
	return MsgIssueAsset{Family: family, Name: name, Symbol: symbol, Source: source, InitSupply: initSupply, MaxSupply: maxSupply, Decimal: decimal, Mintable: mintable, Owner: owner, Operators: operators}
}

// Implements Msg.
// nolint
func (msg MsgIssueAsset) Route() string { return MsgRoute }
func (msg MsgIssueAsset) Type() string  { return "issue" }

// Implements Msg.
func (msg MsgIssueAsset) ValidateBasic() sdk.Error {

	// only accepts alphanumeric characters, _ and -
	reg := regexp.MustCompile(`[^a-zA-Z0-9_-]`)

	if msg.Owner == nil {
		return ErrNilAssetOwner(DefaultCodespace)
	}

	if _, found := MsgIssueFamily[msg.Family]; len(msg.Family) > 0 && !found {
		return ErrInvalidAssetFamily(DefaultCodespace, msg.Family)

	}

	if len(msg.Name) == 0 || !reg.Match([]byte(msg.Name)) {
		return ErrInvalidAssetName(DefaultCodespace, msg.Name)
	}

	if len(msg.Symbol) == 0 || !reg.Match([]byte(msg.Symbol)) {
		return ErrInvalidAssetSymbol(DefaultCodespace, msg.Symbol)
	}

	if msg.InitSupply == 0 || msg.InitSupply > MaxInitSupply {
		return ErrInvalidAssetInitSupply(DefaultCodespace, msg.InitSupply)
	}

	if msg.MaxSupply > 0 && msg.MaxSupply < msg.InitSupply {
		return ErrInvalidAssetMaxSupply(DefaultCodespace, msg.MaxSupply)
	}

	if msg.Decimal > 18 {
		return ErrInvalidAssetDecimal(DefaultCodespace, msg.Decimal)
	}

	return nil
}

// Implements Msg.
func (msg MsgIssueAsset) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgIssueAsset) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
