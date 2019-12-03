package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// GenesisState - all guardian state that must be provided at genesis
type GenesisState struct {
	Profilers []Guardian `json:"profilers"`
	Trustees  []Guardian `json:"trustees"`
}

func NewGenesisState(profilers, trustees []Guardian) GenesisState {
	return GenesisState{
		Profilers: profilers,
		Trustees:  trustees,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	guardian := Guardian{Description: "genesis", AccountType: Genesis}
	return NewGenesisState([]Guardian{guardian}, []Guardian{guardian})
}

// get raw genesis raw message for testing
func DefaultGenesisStateForTest() GenesisState {
	return DefaultGenesisState()
}

func ValidateGenesis(data GenesisState) error {
	for _, profiler := range data.Profilers {
		if len(profiler.Address) == 0 {
			return sdk.ErrInvalidAddress(profiler.Address.String())
		}
		if len(profiler.AddedBy) == 0 {
			return sdk.ErrInvalidAddress(profiler.AddedBy.String())
		}
		if !validAccountType(profiler.AccountType) {
			return ErrInvalidOperator(DefaultCodespace, profiler.AddedBy)
		}
	}

	for _, trustee := range data.Trustees {
		if len(trustee.Address) == 0 {
			return sdk.ErrInvalidAddress(trustee.Address.String())
		}
		if len(trustee.AddedBy) == 0 {
			return sdk.ErrInvalidAddress(trustee.AddedBy.String())
		}
		if !validAccountType(trustee.AccountType) {
			return ErrInvalidOperator(DefaultCodespace, trustee.AddedBy)
		}
	}
	return nil
}
