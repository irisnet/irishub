package service

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	servicekeeper "github.com/irisnet/irismod/modules/service/keeper"
	servicetypes "github.com/irisnet/irismod/modules/service/types"
)

const (
	service_tax_account_address = "iaa1t2tk2g9uyp4szna7cd4k0ymvgf8qs4rxcrykgr"
)

func Migrate(ctx sdk.Context, k servicekeeper.Keeper, bk bankkeeper.Keeper) error {
	oldAcc, err := sdk.AccAddressFromBech32(service_tax_account_address)
	if err != nil {
		return err
	}
	return bk.SendCoinsFromAccountToModule(ctx, oldAcc, servicetypes.FeeCollectorName, bk.GetAllBalances(ctx, oldAcc))
}
