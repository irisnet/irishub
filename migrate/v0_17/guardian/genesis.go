package guardian

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	
)

type GenesisState struct {
	Profilers []Guardian `json:"profilers"`
	Trustees  []Guardian `json:"trustees"`
}

type Guardian struct {
	Description string         `json:"description"`
	AccountType AccountType    `json:"type"`
	Address     sdk.AccAddress `json:"address"`  // this guardian's address
	AddedBy     sdk.AccAddress `json:"added_by"` // address that initiated the AddGuardian tx
}

type AccountType byte

const (
	Genesis  AccountType = 0x01
	Ordinary AccountType = 0x02
)

// String to AccountType byte, Returns ff if invalid.
func AccountTypeFromString(str string) (AccountType, error) {
	switch str {
	case "Genesis":
		return Genesis, nil
	case "Ordinary":
		return Ordinary, nil
	default:
		return AccountType(0xff), errors.Errorf("'%s' is not a valid account type", str)
	}
}

func (bt AccountType) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(fmt.Sprintf("%s", bt.String())))
	default:
		s.Write([]byte(fmt.Sprintf("%v", byte(bt))))
	}
}

// Turns BindingType byte to String
func (bt AccountType) String() string {
	switch bt {
	case Genesis:
		return "Genesis"
	case Ordinary:
		return "Ordinary"
	default:
		return ""
	}
}

// Marshals to JSON using string
func (bt AccountType) MarshalJSON() ([]byte, error) {
	return json.Marshal(bt.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (bt *AccountType) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return nil
	}

	bz2, err := AccountTypeFromString(s)
	if err != nil {
		return err
	}
	*bt = bz2
	return nil
}
