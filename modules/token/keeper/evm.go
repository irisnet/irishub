package keeper

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	tokentypes "github.com/irisnet/irismod/modules/token/types"
	"github.com/irisnet/irismod/types"
)

// CallEVMWithData executes an Ethereum Virtual Machine (EVM) call with the provided data.
//
// Parameters:
//   - ctx: the context in which the EVM call is executed
//   - from: the address initiating the EVM call
//   - contract: the address of the smart contract
//   - data: the data to be sent with the EVM call
//   - commit: boolean indicating whether the EVM call should be committed
//
// Returns:
//   - *types.Result: the result of the EVM call
//   - error: an error if the EVM call encounters any issues
func (k Keeper) CallEVMWithData(
	ctx sdk.Context,
	from common.Address,
	contract *common.Address,
	data []byte,
	commit bool,
) (*types.Result, error) {
	nonce, err := k.accountKeeper.GetSequence(ctx, from.Bytes())
	if err != nil {
		return nil, err
	}

	gasCap := types.DefaultGasCap
	if commit {
		args, err := json.Marshal(types.TransactionArgs{
			From: &from,
			To:   contract,
			Data: (*hexutil.Bytes)(&data),
		})
		if err != nil {
			return nil, errorsmod.Wrapf(tokentypes.ErrJSONMarshal, "failed to marshal tx args: %s", err.Error())
		}

		gas, err := k.evmKeeper.EstimateGas(sdk.WrapSDKContext(ctx), &types.EthCallRequest{
			Args:    args,
			GasCap:  types.DefaultGasCap,
			ChainID: k.evmKeeper.ChainID().Int64(),
		})
		if err != nil {
			return nil, err
		}
		gasCap = gas
	}

	msg := ethtypes.NewMessage(
		from,
		contract,
		nonce,
		big.NewInt(0), // amount
		gasCap,        // gasLimit
		big.NewInt(0), // gasFeeCap
		big.NewInt(0), // gasTipCap
		big.NewInt(0), // gasPrice
		data,
		ethtypes.AccessList{}, // AccessList
		!commit,               // isFake
	)

	res, err := k.evmKeeper.ApplyMessage(ctx, msg, types.NewNoOpTracer(), commit)
	if err != nil {
		return nil, err
	}

	if res.Failed() {
		return nil, errorsmod.Wrap(tokentypes.ErrVMExecution, res.VMError)
	}

	return res, nil
}