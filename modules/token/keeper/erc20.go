package keeper

import (
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/irisnet/irismod/contracts"
	"github.com/irisnet/irismod/token/types"
	v1 "github.com/irisnet/irismod/token/types/v1"
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
	scale uint8,
) (common.Address, error) {
	token, err := k.buildERC20Token(ctx, name, symbol, minUnit, uint32(scale))
	if err != nil {
		return common.Address{}, err
	}

	if len(token.Contract) > 0 {
		return common.Address{}, errorsmod.Wrapf(types.ErrERC20AlreadyExists, "token: %s already deployed erc20 contract: %s", token.Symbol, token.Contract)
	}

	params := k.GetParams(ctx)
	if !params.EnableErc20 {
		return common.Address{}, errorsmod.Wrapf(types.ErrERC20Disabled, "erc20 is disabled")
	}

	if len(params.Beacon) == 0 {
		return common.Address{}, errorsmod.Wrapf(types.ErrBeaconNotSet, "beacon not set")
	}

	deployer := k.getModuleEthAddress(ctx)

	initArgs, err := contracts.ERC20TokenContract.ABI.Pack(
		contracts.MethodInitialize,
		name,
		symbol,
		scale,
		deployer,
	)
	if err != nil {
		return common.Address{}, err
	}

	contractArgs, err := contracts.TokenProxyContract.ABI.Pack(
		"",
		common.HexToAddress(params.Beacon),
		initArgs,
	)
	if err != nil {
		return common.Address{}, errorsmod.Wrapf(types.ErrABIPack, "erc20 metadata is invalid %s: %s", name, err.Error())
	}

	data := make([]byte, len(contracts.TokenProxyContract.Bin)+len(contractArgs))
	copy(data[:len(contracts.TokenProxyContract.Bin)], contracts.TokenProxyContract.Bin)
	copy(data[len(contracts.TokenProxyContract.Bin):], contractArgs)

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

	token.Contract = contractAddr.String()
	k.upsertToken(ctx, *token)

	ctx.EventManager().EmitTypedEvent(&v1.EventDeployERC20{
		Symbol:   symbol,
		Name:     name,
		Scale:    uint32(scale),
		MinUnit:  minUnit,
		Contract: contractAddr.String(),
	})
	return contractAddr, nil
}

// SwapFromERC20 executes a swap from an ERC20 token to its native counterpart
//
// Parameters:
//
//	ctx - the context in which the swap is executed
//	sender - the address of the sender
//	receiver - the address of the receiver
//	wantedAmount - the amount of the token to be swapped out
//
// Return type: error
func (k Keeper) SwapFromERC20(
	ctx sdk.Context,
	sender common.Address,
	receiver sdk.AccAddress,
	wantedAmount sdk.Coin,
) error {
	if !k.ERC20Enabled(ctx) {
		return types.ErrERC20Disabled
	}

	token, err := k.getTokenByMinUnit(ctx, wantedAmount.Denom)
	if err != nil {
		return err
	}
	if len(token.Contract) == 0 {
		return errorsmod.Wrapf(types.ErrERC20NotDeployed, "The token %s is not bound to the corresponding erc20 token", wantedAmount.Denom)
	}

	contract := common.HexToAddress(token.Contract)
	if err := k.BurnERC20(ctx, contract, sender, wantedAmount.Amount.BigInt()); err != nil {
		return err
	}

	mintedCoins := sdk.NewCoins(wantedAmount)
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, mintedCoins); err != nil {
		return err
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, mintedCoins); err != nil {
		return err
	}

	ctx.EventManager().EmitTypedEvent(&v1.EventSwapFromERC20{
		FromContract: contract.String(),
		WantedAmount: &wantedAmount,
		Sender:       sender.String(),
		Receiver:     receiver.String(),
	})
	return nil
}

// SwapToERC20 executes a swap from a native token to its ERC20 token counterpart
//
// Parameters:
//   - ctx: the context
//   - sender: the sender of the amount
//   - receiver: the receiver of the erc20 token
//   - amount:  the amount to be swapped
//
// Returns:
//   - error: error if any.
func (k Keeper) SwapToERC20(
	ctx sdk.Context,
	sender sdk.AccAddress,
	receiver common.Address,
	amount sdk.Coin,
) error {
	if !k.ERC20Enabled(ctx) {
		return types.ErrERC20Disabled
	}

	receiverAcc := k.accountKeeper.GetAccount(ctx, sdk.AccAddress(receiver.Bytes()))
	if receiverAcc != nil {
		if !k.evmKeeper.SupportedKey(receiverAcc.GetPubKey()) {
			return errorsmod.Wrapf(types.ErrUnsupportedKey, "key %s", receiverAcc.GetPubKey())
		}
	}

	token, err := k.getTokenByMinUnit(ctx, amount.Denom)
	if err != nil {
		return err
	}
	if len(token.Contract) == 0 {
		return errorsmod.Wrapf(types.ErrERC20NotDeployed, "token: %s is not bound to the corresponding erc20 token", amount.Denom)
	}
	contract := common.HexToAddress(token.Contract)

	amt := sdk.NewCoins(amount)
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, amt); err != nil {
		return err
	}

	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, amt); err != nil {
		return err
	}

	if err := k.MintERC20(ctx, contract, receiver, amount.Amount.BigInt()); err != nil {
		return err
	}

	ctx.EventManager().EmitTypedEvent(&v1.EventSwapToERC20{
		Amount:     amount,
		Sender:     sender.String(),
		Receiver:   receiver.String(),
		ToContract: token.Contract,
	})
	return nil
}

// MintERC20 mints ERC20 tokens to an account.
//
// Parameters:
//   - ctx: the sdk.Context for the function
//   - contract: the address of the contract
//   - to: the address of the receiver
//   - amount: the amount to mint
//
// Returns:
//   - err : error if any
func (k Keeper) MintERC20(
	ctx sdk.Context,
	contract, to common.Address,
	amount *big.Int,
) error {
	balanceBefore, err := k.BalanceOf(ctx, contract, to)
	if err != nil {
		return err
	}

	abi := contracts.ERC20TokenContract.ABI
	res, err := k.CallEVM(ctx, abi, k.getModuleEthAddress(ctx), contract, true, contracts.MethodMint, to, amount)
	if err != nil {
		return err
	}

	if res.Failed() {
		return errorsmod.Wrapf(
			types.ErrVMExecution, "failed to mint contract: %s, reason: %s",
			contract.String(),
			res.Revert(),
		)
	}

	balanceAfter, err := k.BalanceOf(ctx, contract, to)
	if err != nil {
		return err
	}
	expectBalance := big.NewInt(0).Add(balanceBefore, amount)
	if r := expectBalance.Cmp(balanceAfter); r != 0 {
		return errorsmod.Wrapf(
			types.ErrVMExecution, "failed to mint token correctly, expected after-mint amount is incorrect: %s, expected %d, actual %d",
			contract.String(),
			expectBalance.Int64(),
			balanceAfter.Int64(),
		)
	}
	return nil
}

// BurnERC20 burns a specific amount of ERC20 tokens from a given contract and address.
//
// Parameters:
//   - ctx: the context in which the transaction is executed
//   - contract: the contract address of the ERC20 token
//   - from: the address from which the tokens are burned
//   - amount: the amount of tokens to burn
//
// Returns an error.
func (k Keeper) BurnERC20(
	ctx sdk.Context,
	contract, from common.Address,
	amount *big.Int,
) error {
	balanceBefore, err := k.BalanceOf(ctx, contract, from)
	if err != nil {
		return err
	}

	if r := balanceBefore.Cmp(amount); r < 0 {
		return errorsmod.Wrapf(
			sdkerrors.ErrInsufficientFunds,
			"balance: %d, swap: %d",
			balanceBefore.Int64(),
			amount,
		)
	}

	abi := contracts.ERC20TokenContract.ABI
	res, err := k.CallEVM(ctx, abi, k.getModuleEthAddress(ctx), contract, true, contracts.MethodBurn, from, amount)
	if err != nil {
		return err
	}

	if res.Failed() {
		return errorsmod.Wrapf(types.ErrVMExecution, "failed to burn %d", amount)
	}

	balanceAfter, err := k.BalanceOf(ctx, contract, from)
	if err != nil {
		return err
	}
	expectBalance := big.NewInt(0).Sub(balanceBefore, amount)
	if r := expectBalance.Cmp(balanceAfter); r != 0 {
		return errorsmod.Wrapf(
			types.ErrVMExecution, "failed to burn token correctly, expected after-burn amount is incorrect: %s, expected %d, actual %d",
			contract.String(),
			expectBalance.Int64(),
			balanceAfter.Int64(),
		)
	}
	return nil
}

// UpgradeERC20 upgrades the ERC20 contract to a new implementation.
//
// Parameters:
// - ctx: the SDK context.
// - implementation: the address of the new implementation contract.
//
// Returns:
// - error: an error if the upgrade fails.
func (k Keeper) UpgradeERC20(
	ctx sdk.Context,
	implementation common.Address,
) error {
	params := k.GetParams(ctx)
	if !params.EnableErc20 {
		return errorsmod.Wrapf(types.ErrERC20Disabled, "erc20 is disabled")
	}

	if len(params.Beacon) == 0 {
		return errorsmod.Wrapf(types.ErrBeaconNotSet, "beacon not set")
	}

	beacon := common.HexToAddress(params.Beacon)
	abi := contracts.BeaconContract.ABI
	res, err := k.CallEVM(ctx, abi, k.getModuleEthAddress(ctx), beacon, true, contracts.MethodUpgradeTo, implementation)
	if err != nil {
		return err
	}
	if res.Failed() {
		return errorsmod.Wrapf(
			types.ErrVMExecution, "failed to upgrade contract reason: %s",
			res.Revert(),
		)
	}
	return nil
}

// BalanceOf retrieves the balance of a specific account in the contract.
//
// Parameters:
//   - ctx: the sdk.Context for the function
//   - contract: the address of the contract
//   - account: the address of the account to retrieve the balance for
//
// Returns:
//   - *big.Int: the balance of the specified account
func (k Keeper) BalanceOf(
	ctx sdk.Context,
	contract, account common.Address,
) (*big.Int, error) {
	abi := contracts.ERC20TokenContract.ABI
	res, err := k.CallEVM(ctx, abi, k.getModuleEthAddress(ctx), contract, false, contracts.MethodBalanceOf, account)
	if err != nil {
		return nil, err
	}

	unpacked, err := abi.Unpack(contracts.MethodBalanceOf, res.Ret)
	if err != nil || len(unpacked) == 0 {
		return nil, err
	}

	balance, ok := unpacked[0].(*big.Int)
	if !ok {
		return nil, err
	}

	return balance, nil
}

func (k Keeper) getModuleEthAddress(ctx sdk.Context) common.Address {
	moduleAccount := k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
	return common.BytesToAddress(moduleAccount.GetAddress().Bytes())
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
