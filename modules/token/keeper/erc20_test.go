package keeper_test

import (
	"math/big"

	"github.com/cometbft/cometbft/crypto/tmhash"
	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"

	v1 "irismod.io/token/types/v1"
)

func (suite *KeeperTestSuite) TestDeployERC20() {
	token := v1.NewToken("btc", "Bitcoin Network", "satoshi", 18, 21000000, 21000000, false, owner)

	err := suite.keeper.IssueToken(
		suite.ctx, token.Symbol, token.Name,
		token.MinUnit, token.Scale, token.InitialSupply,
		token.MaxSupply, token.Mintable, token.GetOwner(),
	)
	suite.NoError(err)

	hash, err := suite.keeper.DeployERC20(suite.ctx, token.Name, token.Symbol, token.MinUnit, uint8(token.Scale))
	suite.NoError(err)
	token.Contract = hash.Hex()

	actual, err := suite.keeper.GetToken(suite.ctx, token.Symbol)
	suite.NoError(err)

	suite.EqualValues(&token, actual.(*v1.Token))
}

func (suite *KeeperTestSuite) TestSwapFromERC20() {
	token := v1.NewToken("btc", "Bitcoin Network", "satoshi", 18, 21000000, 21000000, false, owner)

	err := suite.keeper.IssueToken(
		suite.ctx, token.Symbol, token.Name,
		token.MinUnit, token.Scale, token.InitialSupply,
		token.MaxSupply, token.Mintable, token.GetOwner(),
	)
	suite.NoError(err)

	contract, err := suite.keeper.DeployERC20(suite.ctx, token.Name, token.Symbol, token.MinUnit, uint8(token.Scale))
	suite.NoError(err)
	token.Contract = contract.Hex()

	actual, err := suite.keeper.GetToken(suite.ctx, token.Symbol)
	suite.NoError(err)

	suite.EqualValues(&token, actual.(*v1.Token))

	cosmosAddr := sdk.AccAddress(tmhash.SumTruncated([]byte("TestSwapFromERC20")))
	amount := big.NewInt(2e18)
	evmAddr := common.BytesToAddress(cosmosAddr.Bytes())

	suite.Run("mint erc20", func() {
		err = suite.keeper.MintERC20(suite.ctx, contract, evmAddr, amount)
		suite.NoError(err)
	})

	suite.Run("swap from erc20", func() {
		wantedAmount := sdk.NewCoin(token.MinUnit, sdk.NewInt(1e18))

		err = suite.keeper.SwapFromERC20(suite.ctx, evmAddr, cosmosAddr, wantedAmount)
		suite.NoError(err)

		actual := suite.bk.GetBalance(suite.ctx, cosmosAddr, token.MinUnit)
		suite.True(wantedAmount.Equal(actual), "SwapFromERC20 failed: %s != %s", wantedAmount.String(), actual.String())

		balance, err := suite.keeper.BalanceOf(suite.ctx, contract, evmAddr)
		suite.NoError(err)

		expect := big.NewInt(0).Sub(amount, wantedAmount.Amount.BigInt())
		suite.True(expect.Cmp(balance) == 0, "SwapFromERC20 failed balance: %s != %s", expect.String(), balance.String())
	})
}

func (suite *KeeperTestSuite) TestSwapToERC20() {
	token := v1.NewToken("btc", "Bitcoin Network", "satoshi", 18, 21000000, 21000000, false, owner)

	err := suite.keeper.IssueToken(
		suite.ctx, token.Symbol, token.Name,
		token.MinUnit, token.Scale, token.InitialSupply,
		token.MaxSupply, token.Mintable, token.GetOwner(),
	)
	suite.NoError(err)

	contract, err := suite.keeper.DeployERC20(suite.ctx, token.Name, token.Symbol, token.MinUnit, uint8(token.Scale))
	suite.NoError(err)

	sender := token.GetOwner()
	receiver := common.BytesToAddress(sender.Bytes())

	balanceBefore := suite.bk.GetBalance(suite.ctx, sender, token.MinUnit)
	suite.Run("swap to erc20", func() {
		amount := sdk.NewCoin(token.MinUnit, sdk.NewInt(1e18))

		err = suite.keeper.SwapToERC20(suite.ctx, sender, receiver, amount)
		suite.NoError(err)

		balance, err := suite.keeper.BalanceOf(suite.ctx, contract, receiver)
		suite.NoError(err)
		suite.True(amount.Amount.BigInt().Cmp(balance) == 0, "SwapToERC20 failed %s != %s", amount.String(), balance.String())

		actual := suite.bk.GetBalance(suite.ctx, sender, token.MinUnit)
		suite.True(balanceBefore.Sub(amount).IsEqual(actual), "SwapToERC20 failed %s != %s", balanceBefore.String(), actual.String())
	})
}
