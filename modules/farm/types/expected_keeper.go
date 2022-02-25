package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// BankKeeper defines the expected bank keeper (noalias)
type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) error
}

type ValidateLPToken func(ctx sdk.Context, lpTokenDenom string) error

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, moduleName string) authtypes.ModuleAccountI
	SetModuleAccount(ctx sdk.Context, macc authtypes.ModuleAccountI)
}

// DistrKeeper defines the expected distribution keeper (noalias)
type DistrKeeper interface {
	GetFeePool(ctx sdk.Context) (feePool distrtypes.FeePool)
	SetFeePool(ctx sdk.Context, feePool distrtypes.FeePool)
}

// GovKeeper defines the expected gov keeper (noalias)
type GovKeeper interface {
	SubmitProposal(ctx sdk.Context, content govtypes.Content) (govtypes.Proposal, error)
	AddDeposit(ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress, depositAmount sdk.Coins) (bool, error)
	GetProposal(ctx sdk.Context, proposalID uint64) (govtypes.Proposal, bool)
}
