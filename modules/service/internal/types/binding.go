package types

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SvcBinding
type SvcBinding struct {
	DefName     string         `json:"def_name" yaml:"def_name"`
	DefChainID  string         `json:"def_chain_id" yaml:"def_chain_id"`
	BindChainID string         `json:"bind_chain_id" yaml:"bind_chain_id"`
	Provider    sdk.AccAddress `json:"provider" yaml:"provider"`
	BindingType BindingType    `json:"binding_type" yaml:"binding_type"`
	Deposit     sdk.Coins      `json:"deposit" yaml:"deposit"`
	Prices      []sdk.Coin     `json:"price" yaml:"price"`
	Level       Level          `json:"level" yaml:"level"`
	Available   bool           `json:"available" yaml:"available"`
	DisableTime time.Time      `json:"disable_time" yaml:"disable_time"`
}

// Level
type Level struct {
	AvgRspTime int64 `json:"avg_rsp_time" yaml:"avg_rsp_time"`
	UsableTime int64 `json:"usable_time" yaml:"usable_time"`
}

// NewSvcBinding returns a new SvcBinding with the provided values.
func NewSvcBinding(ctx sdk.Context, defChainID, defName, bindChainID string, provider sdk.AccAddress, bindingType BindingType, deposit sdk.Coins, prices []sdk.Coin, level Level, available bool) SvcBinding {
	return SvcBinding{
		DefChainID:  defChainID,
		DefName:     defName,
		BindChainID: bindChainID,
		Provider:    provider,
		BindingType: bindingType,
		Deposit:     deposit,
		Prices:      prices,
		Level:       level,
		Available:   available,
		DisableTime: ctx.BlockHeader().Time,
	}
}

// SvcBindingEqual
func SvcBindingEqual(bindingA, bindingB SvcBinding) bool {
	if bindingA.DefChainID == bindingB.DefChainID &&
		bindingA.DefName == bindingB.DefName &&
		bindingA.BindChainID == bindingB.BindChainID &&
		bindingA.Provider.String() == bindingB.Provider.String() &&
		bindingA.BindingType == bindingB.BindingType &&
		bindingA.Deposit.IsEqual(bindingB.Deposit) &&
		bindingA.Level.AvgRspTime == bindingB.Level.AvgRspTime &&
		bindingA.Level.UsableTime == bindingB.Level.UsableTime &&
		len(bindingA.Prices) == len(bindingB.Prices) &&
		bindingA.Available == bindingB.Available &&
		bindingA.DisableTime.Equal(bindingB.DisableTime) {
		for j, prices := range bindingA.Prices {
			if !prices.IsEqual(bindingB.Prices[j]) {
				return false
			}
		}
		return true
	}
	return false
}

// is valid level?
func validLevel(lv Level) bool {
	return lv.AvgRspTime > 0 && lv.UsableTime > 0 && lv.UsableTime <= 10000
}

// is valid update level?
func validUpdateLevel(lv Level) bool {
	return lv.AvgRspTime >= 0 && lv.UsableTime >= 0 && lv.UsableTime <= 10000
}

func (svcBind SvcBinding) isValid() bool {
	return svcBind.Available
}

// BindingType
type BindingType byte

const (
	Global BindingType = 0x01 // global type
	Local  BindingType = 0x02 // local type
)

// BindingTypeFromString converts string to BindingType byte, returns ff if invalid.
func BindingTypeFromString(str string) (BindingType, error) {
	switch str {
	case "Local":
		return Local, nil
	case "Global":
		return Global, nil
	default:
		return BindingType(0xff), errors.Errorf("'%s' is not a valid binding type", str)
	}
}

// is defined BindingType?
func validBindingType(bt BindingType) bool {
	return bt == Local || bt == Global
}

// Format for Printf / Sprintf, returns bech32 when using %s
func (bt BindingType) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(fmt.Sprintf("%s", bt.String())))
	default:
		s.Write([]byte(fmt.Sprintf("%v", byte(bt))))
	}
}

// String converts BindingType byte to string
func (bt BindingType) String() string {
	switch bt {
	case Local:
		return "Local"
	case Global:
		return "Global"
	default:
		return ""
	}
}

// MarshalJSON marshals BindingType to JSON using string
func (bt BindingType) MarshalJSON() ([]byte, error) {
	return json.Marshal(bt.String())
}

// UnmarshalJSON unmarshals BindingType from JSON assuming Bech32 encoding
func (bt *BindingType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return nil
	}
	bz2, err := BindingTypeFromString(s)
	if err != nil {
		return err
	}
	*bt = bz2
	return nil
}
