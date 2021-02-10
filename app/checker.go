package app

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// ExportStateToCSV export init state to csv file
func (app *IrisApp) ExportStateToCSV(ctx sdk.Context) {
	var waiter sync.WaitGroup
	waiter.Add(5)

	go app.ExportAccounts(ctx, &waiter)
	go app.ExportDelegations(ctx, &waiter)
	go app.ExportUnbondingDelegations(ctx, &waiter)
	go app.ExportRedelegations(ctx, &waiter)
	go app.ExportValidators(ctx, &waiter)

	waiter.Wait()
	panic("ExportStateToCSV Stop")
}

// ExportAccounts export all the account state to csv file
func (app *IrisApp) ExportAccounts(ctx sdk.Context, waiter *sync.WaitGroup) {
	defer waiter.Done()
	app.Logger().Info("=========ExportAccounts Start=========")
	var data = [][]string{
		{"Address", "Balance", "Type", "Module"},
	}
	for _, acc := range app.accountKeeper.GetAllAccounts(ctx) {
		address := acc.GetAddress()
		balances := app.bankKeeper.GetAllBalances(ctx, address)

		var typ string
		var moduleNm string

		switch acc.(type) {
		case *authtypes.BaseAccount:
			typ = "BaseAccount"
		case *authtypes.ModuleAccount:
			typ = "ModuleAccount"
			moduleNm = acc.(*authtypes.ModuleAccount).Name
		}
		data = append(data, []string{
			address.String(),
			balances.AmountOf("uiris").String(),
			typ,
			moduleNm,
		})

	}

	file := filepath.Join(app.homePath, "account.csv")
	writeCSV(file, data)
	app.Logger().Info("=========ExportAccounts End=========")
}

// ExportDelegations export all the delegation state to csv file
func (app *IrisApp) ExportDelegations(ctx sdk.Context, waiter *sync.WaitGroup) {
	defer waiter.Done()
	app.Logger().Info("=========ExportDelegations Start=========")

	var data = [][]string{
		{"DelegatorAddress", "ValidatorAddress", "Shares"},
	}
	for _, delegation := range app.stakingKeeper.GetAllDelegations(ctx) {
		data = append(data, []string{
			delegation.DelegatorAddress,
			delegation.ValidatorAddress,
			delegation.Shares.String(),
		})
	}

	file := filepath.Join(app.homePath, "delegation.csv")
	writeCSV(file, data)
	app.Logger().Info("=========ExportDelegations End=========")
}

// ExportUnbondingDelegations export all the unbonding delegation state to csv file
func (app *IrisApp) ExportUnbondingDelegations(ctx sdk.Context, waiter *sync.WaitGroup) {
	defer waiter.Done()
	app.Logger().Info("=========ExportUnbondingDelegations start=========")

	var data = [][]string{
		{"DelegatorAddress", "ValidatorAddress", "CreationHeight", "CompletionTime", "InitialBalance", "Balance"},
	}
	app.stakingKeeper.IterateUnbondingDelegations(ctx, func(_ int64, ubd types.UnbondingDelegation) (stop bool) {
		for _, entry := range ubd.Entries {
			data = append(data, []string{
				ubd.DelegatorAddress,
				ubd.ValidatorAddress,
				fmt.Sprintf("%d", entry.CreationHeight),
				entry.CompletionTime.String(),
				entry.InitialBalance.String(),
				entry.Balance.String(),
			})
		}
		return false
	})

	file := filepath.Join(app.homePath, "unbondingDelegation.csv")
	writeCSV(file, data)
	app.Logger().Info("=========ExportUnbondingDelegations End=========")
}

// ExportRedelegations export all the unbonding delegation state to csv file
func (app *IrisApp) ExportRedelegations(ctx sdk.Context, waiter *sync.WaitGroup) {
	defer waiter.Done()
	app.Logger().Info("=========ExportRedelegations start=========")

	var data = [][]string{
		{"DelegatorAddress", "ValidatorSrcAddress", "ValidatorDstAddress", "CreationHeight", "CompletionTime", "InitialBalance", "SharesDst"},
	}
	app.stakingKeeper.IterateRedelegations(ctx, func(_ int64, red types.Redelegation) (stop bool) {
		for _, entry := range red.Entries {
			data = append(data, []string{
				red.DelegatorAddress,
				red.ValidatorSrcAddress,
				red.ValidatorDstAddress,
				fmt.Sprintf("%d", entry.CreationHeight),
				entry.CompletionTime.String(),
				entry.InitialBalance.String(),
				entry.SharesDst.String(),
			})
		}
		return false
	})

	file := filepath.Join(app.homePath, "redelegation.csv")
	writeCSV(file, data)
	app.Logger().Info("=========ExportRedelegations End=========")
}

// ExportValidators export all the validator state to csv file
func (app *IrisApp) ExportValidators(ctx sdk.Context, waiter *sync.WaitGroup) {
	defer waiter.Done()
	app.Logger().Info("=========ExportValidators start=========")

	var data = [][]string{
		{
			"OperatorAddress",
			"ConsensusPubkey",
			"Status",
			"Tokens",
			"DelegatorShares",
			"Moniker",
			"Identity",
			"Details",
			"Website",
			"SecurityContact",
			"UnbondingHeight",
			"Rate",
			"MaxRate",
			"MaxChangeRate",
			"MinSelfDelegation",
		},
	}
	app.stakingKeeper.IterateValidators(ctx, func(_ int64, v types.ValidatorI) (stop bool) {
		validator := v.(stakingtypes.Validator)
		pubkey, err := v.ConsPubKey()
		if err != nil {
			panic(err)
		}
		data = append(data, []string{
			validator.OperatorAddress,
			sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, pubkey),
			validator.Status.String(),
			validator.Tokens.String(),
			validator.DelegatorShares.String(),
			validator.Description.Moniker,
			validator.Description.Identity,
			validator.Description.Details,
			validator.Description.Website,
			validator.Description.SecurityContact,
			fmt.Sprintf("%d", validator.UnbondingHeight),
			validator.Commission.Rate.String(),
			validator.Commission.MaxRate.String(),
			validator.Commission.MaxChangeRate.String(),
			validator.MinSelfDelegation.String(),
		})
		return false
	})

	file := filepath.Join(app.homePath, "validator.csv")
	writeCSV(file, data)
	app.Logger().Info("=========ExportValidators End=========")
}

func writeCSV(fileNm string, data [][]string) {
	file, err := os.OpenFile(fileNm, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("open file is failed, err: ", err)
	}
	defer file.Close()
	// 写入UTF-8 BOM，防止中文乱码
	file.WriteString("\xEF\xBB\xBF")
	w := csv.NewWriter(file)
	w.WriteAll(data)
	w.Flush()
}
