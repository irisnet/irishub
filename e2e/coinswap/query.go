package coinswap

import (
	"context"
	"fmt"
	"time"

	"github.com/cosmos/gogoproto/proto"

	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"mods.irisnet.org/e2e"
	coinswaptypes "mods.irisnet.org/modules/coinswap/types"
	tokentypes "mods.irisnet.org/modules/token/types/v1"
	"mods.irisnet.org/simapp"
)

// QueryTestSuite is a suite of end-to-end tests for the nft module
type QueryTestSuite struct {
	e2e.TestSuite
}

// SetupSuite creates a new network for integration tests
func (s *QueryTestSuite) SetupSuite() {
	depInjectOptions := simapp.DepinjectOptions{
		Config:    e2e.AppConfig,
		Providers: []interface{}{
			e2e.ProvideEVMKeeper(),
			e2e.ProvideICS20Keeper(),
		},
	}
		sdk.SetCoinDenomRegex(func() string {
		return `[a-zA-Z][a-zA-Z0-9/\-]{2,127}`
	})

	s.T().Log("setting up integration test suite")
	s.Network = simapp.SetupNetwork(s.T(),depInjectOptions)
}

// TestCoinswap tests all query command in the nft module
func (s *QueryTestSuite) TestCoinswap() {
	val := s.Validators[0]
	clientCtx := val.ClientCtx
	// ---------------------------------------------------------------------------

	from := val.Address
	symbol := "kitty"
	name := "Kitty Token"
	minUnit := "kitty"
	scale := uint32(0)
	initialSupply := uint64(100000000)
	maxSupply := uint64(200000000)
	mintable := true
	baseURL := val.APIAddress
	lptDenom := "lpt-1"

	// issue token
	msgIssueToken := &tokentypes.MsgIssueToken{
		Symbol:        symbol,
		Name:          name,
		Scale:         scale,
		MinUnit:       minUnit,
		InitialSupply: initialSupply,
		MaxSupply:     maxSupply,
		Mintable:      mintable,
		Owner:         from.String(),
	}
	txResult := s.BlockSendMsgs(s.T(), msgIssueToken)
	s.Require().Equal(0, txResult.Code, "send issue token msg failed")

	//_ = tokentestutil.IssueTokenExec(s.T(), s.Network, clientCtx, from.String(), args...)

	balances := simapp.QueryBalancesExec(s.T(), s.Network, clientCtx, from.String())
	s.Require().Equal("100000000", balances.AmountOf(symbol).String())
	s.Require().Equal("399986975", balances.AmountOf(sdk.DefaultBondDenom).String())

	// test add liquidity (poor not exist)
	status, err := clientCtx.Client.Status(context.Background())
	s.Require().NoError(err)
	deadline := status.SyncInfo.LatestBlockTime.Add(time.Minute)

	msgAddLiquidity := &coinswaptypes.MsgAddLiquidity{
		MaxToken:         sdk.NewCoin(symbol, sdk.NewInt(1000)),
		ExactStandardAmt: sdk.NewInt(1000),
		MinLiquidity:     sdk.NewInt(1000),
		Deadline:         deadline.Unix(),
		Sender:           from.String(),
	}
	s.SendMsgs(s.T(), msgAddLiquidity)

	balances = simapp.QueryBalancesExec(s.T(), s.Network, clientCtx, from.String())
	s.Require().Equal("99999000", balances.AmountOf(symbol).String())
	s.Require().Equal("399980965", balances.AmountOf(sdk.DefaultBondDenom).String())
	s.Require().Equal("1000", balances.AmountOf(lptDenom).String())

	queryPoolResponse := proto.Message(&coinswaptypes.QueryLiquidityPoolResponse{})
	url := fmt.Sprintf("%s/irismod/coinswap/pools/%s", baseURL, lptDenom)
	resp, err := testutil.GetRequest(url)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp, queryPoolResponse))

	queryPool := queryPoolResponse.(*coinswaptypes.QueryLiquidityPoolResponse)
	s.Require().Equal("1000", queryPool.Pool.Standard.Amount.String())
	s.Require().Equal("1000", queryPool.Pool.Token.Amount.String())
	s.Require().Equal("1000", queryPool.Pool.Lpt.Amount.String())

	// test add liquidity (poor exist)
	status, err = clientCtx.Client.Status(context.Background())
	s.Require().NoError(err)
	deadline = status.SyncInfo.LatestBlockTime.Add(time.Minute)

	msgAddLiquidity = &coinswaptypes.MsgAddLiquidity{
		MaxToken:         sdk.NewCoin(symbol, sdk.NewInt(2001)),
		ExactStandardAmt: sdk.NewInt(2000),
		MinLiquidity:     sdk.NewInt(2000),
		Deadline:         deadline.Unix(),
		Sender:           from.String(),
	}
	s.SendMsgs(s.T(), msgAddLiquidity)

	balances = simapp.QueryBalancesExec(s.T(), s.Network, clientCtx, from.String())
	s.Require().Equal("99996999", balances.AmountOf(symbol).String())
	s.Require().Equal("399978955", balances.AmountOf(sdk.DefaultBondDenom).String())
	s.Require().Equal("3000", balances.AmountOf(lptDenom).String())

	url = fmt.Sprintf("%s/irismod/coinswap/pools/%s", baseURL, lptDenom)
	resp, err = testutil.GetRequest(url)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp, queryPoolResponse))

	s.Require().Equal("3000", queryPool.Pool.Standard.Amount.String())
	s.Require().Equal("3001", queryPool.Pool.Token.Amount.String())
	s.Require().Equal("3000", queryPool.Pool.Lpt.Amount.String())

	// test sell order
	msgSellOrder := &coinswaptypes.MsgSwapOrder{
		Input: coinswaptypes.Input{
			Address: from.String(),
			Coin:    sdk.NewCoin(symbol, sdk.NewInt(1000)),
		},
		Output: coinswaptypes.Output{
			Address: from.String(),
			Coin:    sdk.NewInt64Coin(s.BondDenom, 748),
		},
		Deadline:   deadline.Unix(),
		IsBuyOrder: false,
	}
	s.SendMsgs(s.T(), msgSellOrder)

	balances = simapp.QueryBalancesExec(s.T(), s.Network, clientCtx, from.String())
	s.Require().Equal("99995999", balances.AmountOf(symbol).String())
	s.Require().Equal("399979693", balances.AmountOf(sdk.DefaultBondDenom).String())
	s.Require().Equal("3000", balances.AmountOf(lptDenom).String())

	url = fmt.Sprintf("%s/irismod/coinswap/pools/%s", baseURL, lptDenom)
	resp, err = testutil.GetRequest(url)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp, queryPoolResponse))

	s.Require().Equal("2252", queryPool.Pool.Standard.Amount.String())
	s.Require().Equal("4001", queryPool.Pool.Token.Amount.String())
	s.Require().Equal("3000", queryPool.Pool.Lpt.Amount.String())

	// test buy order
	msgBuyOrder := &coinswaptypes.MsgSwapOrder{
		Input: coinswaptypes.Input{
			Address: from.String(),
			Coin:    sdk.NewInt64Coin(s.BondDenom, 753),
		},
		Output: coinswaptypes.Output{
			Address: from.String(),
			Coin:    sdk.NewCoin(symbol, sdk.NewInt(1000)),
		},
		Deadline:   deadline.Unix(),
		IsBuyOrder: true,
	}
	s.SendMsgs(s.T(), msgBuyOrder)

	balances = simapp.QueryBalancesExec(s.T(), s.Network, clientCtx, from.String())
	s.Require().Equal("99996999", balances.AmountOf(symbol).String())
	s.Require().Equal("399978930", balances.AmountOf(sdk.DefaultBondDenom).String())
	s.Require().Equal("3000", balances.AmountOf(lptDenom).String())

	url = fmt.Sprintf("%s/irismod/coinswap/pools/%s", baseURL, lptDenom)
	resp, err = testutil.GetRequest(url)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp, queryPoolResponse))

	s.Require().Equal("3005", queryPool.Pool.Standard.Amount.String())
	s.Require().Equal("3001", queryPool.Pool.Token.Amount.String())
	s.Require().Equal("3000", queryPool.Pool.Lpt.Amount.String())

	// Test remove liquidity (remove part)
	msgRemoveLiquidity := &coinswaptypes.MsgRemoveLiquidity{
		WithdrawLiquidity: sdk.NewCoin(lptDenom, sdk.NewInt(2000)),
		MinToken:          sdk.NewInt(2000),
		MinStandardAmt:    sdk.NewInt(2000),
		Deadline:          deadline.Unix(),
		Sender:            from.String(),
	}

	// prepare txBuilder with msg
	s.SendMsgs(s.T(), msgRemoveLiquidity)

	balances = simapp.QueryBalancesExec(s.T(), s.Network, clientCtx, from.String())
	s.Require().Equal("99998999", balances.AmountOf(symbol).String())
	s.Require().Equal("399980923", balances.AmountOf(sdk.DefaultBondDenom).String())
	s.Require().Equal("1000", balances.AmountOf(lptDenom).String())

	url = fmt.Sprintf("%s/irismod/coinswap/pools/%s", baseURL, lptDenom)
	resp, err = testutil.GetRequest(url)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp, queryPoolResponse))

	s.Require().Equal("1002", queryPool.Pool.Standard.Amount.String())
	s.Require().Equal("1001", queryPool.Pool.Token.Amount.String())
	s.Require().Equal("1000", queryPool.Pool.Lpt.Amount.String())

	// Test remove liquidity (remove all)
	msgRemoveLiquidity = &coinswaptypes.MsgRemoveLiquidity{
		WithdrawLiquidity: sdk.NewCoin(lptDenom, sdk.NewInt(1000)),
		MinToken:          sdk.NewInt(1000),
		MinStandardAmt:    sdk.NewInt(1000),
		Deadline:          deadline.Unix(),
		Sender:            from.String(),
	}

	// prepare txBuilder with msg
	s.SendMsgs(s.T(), msgRemoveLiquidity)

	balances = simapp.QueryBalancesExec(s.T(), s.Network, clientCtx, from.String())
	s.Require().Equal("100000000", balances.AmountOf(symbol).String())
	s.Require().Equal("399981915", balances.AmountOf(sdk.DefaultBondDenom).String())
	s.Require().Equal("0", balances.AmountOf(lptDenom).String())

	url = fmt.Sprintf("%s/irismod/coinswap/pools/%s", baseURL, lptDenom)
	resp, err = testutil.GetRequest(url)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp, queryPoolResponse))

	s.Require().Equal("0", queryPool.Pool.Standard.Amount.String())
	s.Require().Equal("0", queryPool.Pool.Token.Amount.String())
	s.Require().Equal("0", queryPool.Pool.Lpt.Amount.String())

	queryPoolsResponse := proto.Message(&coinswaptypes.QueryLiquidityPoolsResponse{})
	url = fmt.Sprintf("%s/irismod/coinswap/pools", baseURL)
	resp, err = testutil.GetRequest(url)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp, queryPoolsResponse))

	queryPools := queryPoolsResponse.(*coinswaptypes.QueryLiquidityPoolsResponse)
	s.Require().Len(queryPools.Pools, 1)
}
