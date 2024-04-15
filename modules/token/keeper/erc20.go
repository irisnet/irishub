package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/irisnet/irismod/contracts"
	"github.com/irisnet/irismod/modules/token/types"
	v1 "github.com/irisnet/irismod/modules/token/types/v1"
)

// DeployERC20 deploys an ERC20 token contract.
//
// Parameters:
//   - ctx: the context
//   - name: the name of the token
//   - symbol: the symbol of the token
//   - minUnit: the symbol of the minUnit
//   - scale: the scale of the token
//
// Returns:
//   - Address: the contract address.
//   - error: error if any.
func (k Keeper) DeployERC20(
	ctx sdk.Context,
	name string,
	symbol string,
	minUnit string,
	scale int8,
) (common.Address, error) {
	contractArgs, err := contracts.ERC20TokenContract.ABI.Pack(
		"",
		name,
		symbol,
		scale,
	)
	if err != nil {
		return common.Address{}, errorsmod.Wrapf(types.ErrABIPack, "erc20 metadata is invalid %s: %s", name, err.Error())
	}
	deployer := k.moduleAddress()

	data := make([]byte, len(contracts.ERC20TokenContract.Bin)+len(contractArgs))
	copy(data[:len(contracts.ERC20TokenContract.Bin)], contracts.ERC20TokenContract.Bin)
	copy(data[len(contracts.ERC20TokenContract.Bin):], contractArgs)

	nonce, err := k.accountKeeper.GetSequence(ctx, sdk.AccAddress(deployer.Bytes()))
	if err != nil {
		return common.Address{}, err
	}
	contractAddr := crypto.CreateAddress(deployer, nonce)
	result, err := k.CallEVMWithData(ctx, deployer, nil, data, true)
	if err != nil {
		return common.Address{}, errorsmod.Wrapf(err, "failed to deploy contract for token %s", name)
	}
	if result.Failed() {
		return common.Address{}, errorsmod.Wrapf(types.ErrVMExecution, "failed to deploy contract for %s, reason: %s", name, result.Revert())
	}

	ctx.EventManager().EmitTypedEvent(&v1.EventDeployERC20{
		Symbol:   symbol,
		Name:     name,
		Scale:    uint32(scale),
		MinUnit:  minUnit,
		Contract: contractAddr.String(),
	})
	return contractAddr, nil
}

func (k Keeper) moduleAddress() common.Address {
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	return common.BytesToAddress(moduleAddr.Bytes())
}

func (k Keeper) buildERC20Token(
	ctx sdk.Context,
	name string,
	symbol string,
	minUnit string,
	scale uint32,
) (*v1.Token, error) {
	if !k.HasMinUint(ctx, minUnit) {
		if k.HasSymbol(ctx, symbol) {
			return nil, errorsmod.Wrapf(types.ErrSymbolAlreadyExists, "symbol already exists: %s", symbol)
		}
		if !k.ics20Keeper.HasTrace(ctx, minUnit) {
			return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "token: %s does not exist", minUnit)
		}
		return &v1.Token{
			Symbol:   symbol,
			Name:     name,
			Scale:    scale,
			MinUnit:  minUnit,
			Mintable: true,
			Owner:    k.accountKeeper.GetModuleAddress(types.ModuleName).String(),
		}, nil
	}

	token, err := k.getTokenByMinUnit(ctx, minUnit)
	if err != nil {
		return nil, err
	}
	return &token, nil
}
