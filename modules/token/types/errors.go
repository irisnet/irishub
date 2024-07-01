package types

import (
	errorsmod "cosmossdk.io/errors"
)

// token module sentinel errors
var (
	ErrInvalidName          = errorsmod.Register(ModuleName, 2, "invalid token name")
	ErrInvalidMinUnit       = errorsmod.Register(ModuleName, 3, "invalid token min unit")
	ErrInvalidSymbol        = errorsmod.Register(ModuleName, 4, "invalid standard denom")
	ErrInvalidInitSupply    = errorsmod.Register(ModuleName, 5, "invalid token initial supply")
	ErrInvalidMaxSupply     = errorsmod.Register(ModuleName, 6, "invalid token maximum supply")
	ErrInvalidScale         = errorsmod.Register(ModuleName, 7, "invalid token scale")
	ErrSymbolAlreadyExists  = errorsmod.Register(ModuleName, 8, "symbol already exists")
	ErrMinUnitAlreadyExists = errorsmod.Register(ModuleName, 9, "min unit already exists")
	ErrTokenNotExists       = errorsmod.Register(ModuleName, 10, "token does not exist")
	ErrInvalidToAddress     = errorsmod.Register(ModuleName, 11, "the new owner must not be same as the original owner")
	ErrInvalidOwner         = errorsmod.Register(ModuleName, 12, "invalid token owner")
	ErrNotMintable          = errorsmod.Register(ModuleName, 13, "token is not mintable")
	ErrNotFoundTokenAmt     = errorsmod.Register(ModuleName, 14, "burned token amount not found")
	ErrInvalidAmount        = errorsmod.Register(ModuleName, 15, "invalid amount")
	ErrInvalidBaseFee       = errorsmod.Register(ModuleName, 16, "invalid base fee")
	ErrInvalidSwap          = errorsmod.Register(ModuleName, 17, "unregistered swapable fee token")
	ErrInsufficientFee      = errorsmod.Register(ModuleName, 18, "the amount of tokens after swap is less than 1")
	ErrJSONMarshal          = errorsmod.Register(ModuleName, 19, "failed to marshal JSON bytes")
	ErrVMExecution          = errorsmod.Register(ModuleName, 20, "evm transaction execution failed")
	ErrABIPack              = errorsmod.Register(ModuleName, 21, "contract ABI pack failed")
	ErrERC20AlreadyExists   = errorsmod.Register(ModuleName, 22, "erc20 contract already exists")
	ErrERC20NotDeployed     = errorsmod.Register(ModuleName, 23, "erc20 contract not deployed")
	ErrUnsupportedKey       = errorsmod.Register(ModuleName, 24, "evm not supported public key")
	ErrInvalidContract      = errorsmod.Register(ModuleName, 25, "invalid contract")
	ErrERC20Disabled        = errorsmod.Register(ModuleName, 26, "erc20 swap is disabled")
	ErrBeaconNotSet        = errorsmod.Register(ModuleName, 27, "beacon contract not set")
)
