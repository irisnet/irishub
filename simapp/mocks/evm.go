package mocks

import (
	"context"
	"math/big"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"

	tokentypes "github.com/irisnet/irismod/modules/token/types"
	"github.com/irisnet/irismod/types"
)

var (
	_ tokentypes.EVMKeeper   = (*evm)(nil)
	_ tokentypes.ICS20Keeper = (*transferKeeper)(nil)
)

type evm struct {
	erc20s map[common.Address]*erc20
}

// ApplyMessage implements types.EVMKeeper.
func (e *evm) ApplyMessage(ctx sdk.Context, msg core.Message, tracer vm.EVMLogger, commit bool) (*types.Result, error) {
	panic("unimplemented")
}

// ChainID implements types.EVMKeeper.
func (e *evm) ChainID() *big.Int {
	return big.NewInt(16688)
}

// EstimateGas implements types.EVMKeeper.
func (e *evm) EstimateGas(ctx context.Context, req *types.EthCallRequest) (uint64, error) {
	return 3000000, nil
}

// FeeDenom implements types.EVMKeeper.
func (e *evm) FeeDenom() string {
	return "eris"
}

// SupportedKey implements types.EVMKeeper.
func (e *evm) SupportedKey(pubKey cryptotypes.PubKey) bool {
	return true
}

type erc20 struct {
	scale        int8
	name, symbol string

	balance map[common.Address]*big.Int
}

type transferKeeper struct{}

// HasTrace implements types.ICS20Keeper.
func (t *transferKeeper) HasTrace(ctx sdk.Context, denom string) bool {
	return true
}
