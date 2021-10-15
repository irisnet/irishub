package v150_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/irisnet/irismod/modules/coinswap/legacy/v150"
	coinswaptypes "github.com/irisnet/irismod/modules/coinswap/types"
	"github.com/irisnet/irismod/simapp"
)

const (
	denomStandard = sdk.DefaultBondDenom
	denomBTC      = "btc"
	denomETH      = "eth"
	denomLptBTC   = "swapbtc"
	denomLptETH   = "swapeth"
)

var (
	addrSender1 = sdk.AccAddress(tmhash.SumTruncated([]byte("addrSender1")))
	addrSender2 = sdk.AccAddress(tmhash.SumTruncated([]byte("addrSender2")))
	poolAddrBTC = v150.GetReservePoolAddr(denomLptBTC)
	poolAddrETH = v150.GetReservePoolAddr(denomLptETH)
)

type (
	verifyFunc = func(ctx sdk.Context, t *testing.T)
)

func TestMigrate(t *testing.T) {
	sdk.SetCoinDenomRegex(func() string {
		return `[a-zA-Z][a-zA-Z0-9/\-]{2,127}`
	})
	app, verify := setupWithGenesisAccounts()
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	err := v150.Migrate(ctx, app.CoinswapKeeper, app.BankKeeper, app.AccountKeeper)
	assert.NoError(t, err)

	//app.BaseApp.Commit()
	verify(ctx, t)
	//perform an Invariants check
	app.CrisisKeeper.AssertInvariants(ctx)
}

func setupWithGenesisAccounts() (*simapp.SimApp, verifyFunc) {
	standardCoin := sdk.NewCoin(denomStandard, sdk.NewIntWithDecimal(1, 18))
	ethCoin := sdk.NewCoin(denomETH, sdk.NewIntWithDecimal(1, 18))
	btcCoin := sdk.NewCoin(denomBTC, sdk.NewIntWithDecimal(1, 18))
	lptBTCCoin := sdk.NewCoin(denomLptBTC, sdk.NewIntWithDecimal(1, 18))
	lptETHCoin := sdk.NewCoin(denomLptETH, sdk.NewIntWithDecimal(1, 18))

	sender1Balances := banktypes.Balance{
		Address: addrSender1.String(),
		Coins: sdk.NewCoins(
			standardCoin,
			lptETHCoin,
		),
	}

	sender2Balances := banktypes.Balance{
		Address: addrSender2.String(),
		Coins: sdk.NewCoins(
			standardCoin,
			lptBTCCoin,
		),
	}

	poolBTCBalances := banktypes.Balance{
		Address: poolAddrBTC.String(),
		Coins: sdk.NewCoins(
			standardCoin,
			btcCoin,
		),
	}

	poolETHBalances := banktypes.Balance{
		Address: poolAddrETH.String(),
		Coins: sdk.NewCoins(
			standardCoin,
			ethCoin,
		),
	}

	senderAcc1 := &authtypes.BaseAccount{
		Address: addrSender1.String(),
	}

	senderAcc2 := &authtypes.BaseAccount{
		Address: addrSender2.String(),
	}

	poolBTCAcc := &authtypes.BaseAccount{
		Address: poolAddrBTC.String(),
	}

	poolETHAcc := &authtypes.BaseAccount{
		Address: poolAddrETH.String(),
	}

	genAccs := []authtypes.GenesisAccount{senderAcc1, senderAcc2, poolBTCAcc, poolETHAcc}
	app := simapp.SetupWithGenesisAccounts(genAccs, sender1Balances, sender2Balances, poolBTCBalances, poolETHBalances)

	verify := func(ctx sdk.Context, t *testing.T) {
		ethPoolId := coinswaptypes.GetPoolId(denomETH)
		ethPool, has := app.CoinswapKeeper.GetPool(ctx, ethPoolId)
		assert.True(t, has)

		btcPoolId := coinswaptypes.GetPoolId(denomBTC)
		btcPool, has := app.CoinswapKeeper.GetPool(ctx, btcPoolId)
		assert.True(t, has)

		// Verify the balance of sender1
		{
			sender1Balances := app.BankKeeper.GetAllBalances(ctx, addrSender1)

			expectsender1Balances := sdk.NewCoins(
				standardCoin,
				sdk.NewCoin(ethPool.LptDenom, lptETHCoin.Amount),
			)
			assert.Equal(t, expectsender1Balances.String(), sender1Balances.String())
		}

		// Verify the balance of sender2
		{
			sender2Balances := app.BankKeeper.GetAllBalances(ctx, addrSender2)

			expectsender2Balances := sdk.NewCoins(
				standardCoin,
				sdk.NewCoin(btcPool.LptDenom, lptBTCCoin.Amount),
			)
			assert.Equal(t, expectsender2Balances.String(), sender2Balances.String())
		}

		// Verify the balance of poolAddrBTC
		{
			srcPoolBTCBalances := app.BankKeeper.GetAllBalances(ctx, poolAddrBTC)
			assert.True(t, srcPoolBTCBalances.IsZero())

			poolBTCAddr, err := sdk.AccAddressFromBech32(btcPool.EscrowAddress)
			assert.NoError(t, err)

			dstPoolBTCBalances := app.BankKeeper.GetAllBalances(ctx, poolBTCAddr)
			assert.Equal(t, poolBTCBalances.Coins.String(), dstPoolBTCBalances.String())
		}

		// Verify the balance of poolAddrETH
		{
			srcPoolETHBalances := app.BankKeeper.GetAllBalances(ctx, poolAddrETH)
			assert.True(t, srcPoolETHBalances.IsZero())

			poolETHAddr, err := sdk.AccAddressFromBech32(ethPool.EscrowAddress)
			assert.NoError(t, err)

			dstPoolETHBalances := app.BankKeeper.GetAllBalances(ctx, poolETHAddr)
			assert.Equal(t, poolETHBalances.Coins.String(), dstPoolETHBalances.String())
		}

	}
	return app, verify
}
