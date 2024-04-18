package wrapper

import (
	"context"
	"math/big"
	"strings"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	evmkeeper "github.com/evmos/ethermint/x/evm/keeper"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	ibctransferkeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"

	tokentypes "github.com/irisnet/irismod/modules/token/types"
	irismodtypes "github.com/irisnet/irismod/types"
)

var (
	_ tokentypes.EVMKeeper   = (*evmKeeper)(nil)
	_ tokentypes.ICS20Keeper = (*ics20Keeper)(nil)
)

// NewEVMKeeper wraps the given evmkeeper.Keeper and returns a new evmKeeper.
//
// ek: The evmkeeper.Keeper to be wrapped.
// Returns a pointer to the wrapped evmKeeper.
func NewEVMKeeper(ek *evmkeeper.Keeper) tokentypes.EVMKeeper {
	return &evmKeeper{ek: ek}
}

type evmKeeper struct {
	ek *evmkeeper.Keeper
}

// ApplyMessage implements types.EVMKeeper.
func (e *evmKeeper) ApplyMessage(ctx sdk.Context, msg core.Message, tracer vm.EVMLogger, commit bool) (*irismodtypes.Result, error) {
	res, err := e.ek.ApplyMessage(ctx, msg, tracer, commit)
	if err != nil {
		return nil, err
	}
	return &irismodtypes.Result{
		Hash:    res.Hash,
		Logs:    evmtypes.LogsToEthereum(res.Logs),
		Ret:     res.Ret,
		VMError: res.VmError,
		GasUsed: res.GasUsed,
	}, nil
}

// ChainID implements types.EVMKeeper.
func (e *evmKeeper) ChainID() *big.Int {
	return e.ek.ChainID()
}

// EstimateGas implements types.EVMKeeper.
func (e *evmKeeper) EstimateGas(ctx context.Context, req *irismodtypes.EthCallRequest) (uint64, error) {
	res, err := e.ek.EstimateGas(ctx, &evmtypes.EthCallRequest{
		Args:            req.Args,
		GasCap:          req.GasCap,
		ProposerAddress: req.ProposerAddress,
		ChainId:         req.ChainID,
	})
	if err != nil {
		return 0, err
	}
	return res.Gas, nil
}

// SupportedKey implements types.EVMKeeper.
func (e *evmKeeper) SupportedKey(pubKey cryptotypes.PubKey) bool {
	_, ok := pubKey.(*ethsecp256k1.PubKey)
	return ok
}

// NewICS20Keeper wraps the given ibctransferkeeper.Keeper into an ics20Keeper.
//
// Parameters:
// - ik: the ibctransferkeeper.Keeper to be wrapped.
//
// Return:
// - *ics20Keeper: the wrapped ics20Keeper.
func NewICS20Keeper(ik ibctransferkeeper.Keeper) tokentypes.ICS20Keeper {
	return &ics20Keeper{ik: ik}
}

type ics20Keeper struct {
	ik ibctransferkeeper.Keeper
}

// HasTrace implements types.ICS20Keeper.
func (i *ics20Keeper) HasTrace(ctx sdk.Context, denom string) bool {
	hash, err := ibctransfertypes.ParseHexHash(strings.TrimPrefix(denom, "ibc/"))
	if err != nil {
		return false
	}
	_, has := i.ik.GetDenomTrace(ctx, hash)
	return has
}
