package service

import (
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	servicekeeper "github.com/irisnet/irismod/modules/service/keeper"
	servicetypes "github.com/irisnet/irismod/modules/service/types"
)

const (
	// TaxAccName is the root string for the service tax account address
	TaxAccName = "service_tax_account"
)

func Migrate(ctx sdk.Context, k servicekeeper.Keeper, bk bankkeeper.Keeper) error {
	oldAcc := sdk.AccAddress(crypto.AddressHash([]byte(TaxAccName)))
	params := servicetypes.NewParams(
		k.MaxRequestTimeout(ctx),
		k.MinDepositMultiple(ctx),
		k.MinDeposit(ctx),
		k.ServiceFeeTax(ctx),
		k.SlashFraction(ctx),
		k.ComplaintRetrospect(ctx),
		k.ArbitrationTimeLimit(ctx),
		k.TxSizeLimit(ctx),
		k.BaseDenom(ctx),
		false,
	)

	k.SetParams(ctx, params)
	return bk.SendCoinsFromAccountToModule(ctx, oldAcc, servicetypes.FeeCollectorName, bk.GetAllBalances(ctx, oldAcc))
}
