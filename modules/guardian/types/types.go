package types

import (
	"fmt"

	"github.com/pkg/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewSuper constructs a super
func NewSuper(description string, accountType AccountType, address, addedBy sdk.AccAddress) Super {
	return Super{
		Description: description,
		AccountType: accountType,
		Address:     address.String(),
		AddedBy:     addedBy.String(),
	}
}

// Equal returns if the guardian is equal to specified guardian
func (g Super) Equal(super Super) bool {
	return g.Address == super.Address &&
		g.AddedBy == super.AddedBy &&
		g.Description == super.Description &&
		g.AccountType == super.AccountType
}

// AccountTypeFromString converts string to AccountType byte, Returns ff if invalid.
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

// ValidAccountType returns true if the AccountType option is valid and false otherwise.
func ValidAccountType(option AccountType) bool {
	if option == Genesis ||
		option == Ordinary {
		return true
	}
	return false
}

// Marshal needed for protobuf compatibility.
func (at AccountType) Marshal() ([]byte, error) {
	return []byte{byte(at)}, nil
}

// Unmarshal needed for protobuf compatibility.
func (at *AccountType) Unmarshal(data []byte) error {
	*at = AccountType(data[0])
	return nil
}

// Format implements the fmt.Formatter interface.
func (at AccountType) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(at.String()))
	default:
		s.Write([]byte(fmt.Sprintf("%v", byte(at))))
	}
}
