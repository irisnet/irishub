package mocks

import (
	"github.com/ethereum/go-ethereum/common"

	tokentypes "github.com/irisnet/irismod/modules/token/types"
)

// ProvideEVMKeeper returns an instance of tokentypes.EVMKeeper.
//
// No parameters.
// Returns a tokentypes.EVMKeeper.
func ProvideEVMKeeper() tokentypes.EVMKeeper {
	return &evm{
		erc20s: make(map[common.Address]*erc20),
	}
}

// ProvideICS20Keeper returns an instance of tokentypes.ICS20Keeper.
//
// No parameters.
// Returns a tokentypes.ICS20Keeper.
func ProvideICS20Keeper() tokentypes.ICS20Keeper {
	return &transferKeeper{}
}
