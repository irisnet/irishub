package types

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/irisnet/irismod/types"
)

// BankKeeper defines the expected bank keeper (noalias)
type BankKeeper interface {
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error

	GetSupply(ctx sdk.Context, denom string) sdk.Coin
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin

	SendCoinsFromModuleToAccount(
		ctx sdk.Context,
		senderModule string,
		recipientAddr sdk.AccAddress,
		amt sdk.Coins,
	) error
	SendCoinsFromAccountToModule(
		ctx sdk.Context,
		senderAddr sdk.AccAddress,
		recipientModule string,
		amt sdk.Coins,
	) error
	SendCoinsFromModuleToModule(
		ctx sdk.Context,
		senderModule, recipientModule string,
		amt sdk.Coins,
	) error

	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins

	SetDenomMetaData(ctx sdk.Context, denomMetaData banktypes.Metadata)
	GetDenomMetaData(ctx sdk.Context, denom string) (banktypes.Metadata, bool)
	GetBlockedAddresses() map[string]bool
}

// AccountKeeper defines the expected account keeper
type AccountKeeper interface {
	GetModuleAddress(moduleName string) sdk.AccAddress
	GetSequence(sdk.Context, sdk.AccAddress) (uint64, error)
	GetAccount(sdk.Context, sdk.AccAddress) authtypes.AccountI
}

// EVMKeeper defines the expected keeper of the evm module
type EVMKeeper interface {
	FeeDenom() string
	ChainID() *big.Int

	SupportedKey(pubKey cryptotypes.PubKey) bool
	EstimateGas(ctx context.Context, req *types.EthCallRequest) (uint64, error)
	ApplyMessage(ctx sdk.Context, msg core.Message, tracer vm.EVMLogger, commit bool) (*types.Result, error)
}

// ICS20Keeper defines the expected keeper of ICS20
type ICS20Keeper interface{
	HasTrace(ctx sdk.Context, denom string) bool
}

