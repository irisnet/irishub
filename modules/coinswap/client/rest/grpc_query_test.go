package rest_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/tidwall/gjson"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	codectype "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	coinswaptypes "github.com/irisnet/irismod/modules/coinswap/types"
	tokencli "github.com/irisnet/irismod/modules/token/client/cli"
	tokentestutil "github.com/irisnet/irismod/modules/token/client/testutil"
	"github.com/irisnet/irismod/simapp"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	cfg := simapp.NewConfig()
	cfg.NumValidators = 1

	s.cfg = cfg
	s.network = network.New(s.T(), cfg)

	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) TestCoinswap() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	// ---------------------------------------------------------------------------

	from := val.Address
	symbol := "kitty"
	name := "Kitty Token"
	minUnit := "kitty"
	scale := 0
	initialSupply := int64(100000000)
	maxSupply := int64(200000000)
	mintable := true
	baseURL := val.APIAddress
	uniKitty := "uni:kitty"

	//------test GetCmdIssueToken()-------------
	args := []string{
		fmt.Sprintf("--%s=%s", tokencli.FlagSymbol, symbol),
		fmt.Sprintf("--%s=%s", tokencli.FlagName, name),
		fmt.Sprintf("--%s=%s", tokencli.FlagMinUnit, minUnit),
		fmt.Sprintf("--%s=%d", tokencli.FlagScale, scale),
		fmt.Sprintf("--%s=%d", tokencli.FlagInitialSupply, initialSupply),
		fmt.Sprintf("--%s=%d", tokencli.FlagMaxSupply, maxSupply),
		fmt.Sprintf("--%s=%t", tokencli.FlagMintable, mintable),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}
	respType := proto.Message(&sdk.TxResponse{})
	expectedCode := uint32(0)
	bz, err := tokentestutil.IssueTokenExec(clientCtx, from.String(), args...)

	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp := respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	respType = proto.Message(&banktypes.QueryAllBalancesResponse{})
	out, err := simapp.QueryBalancesExec(clientCtx, from.String())
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), respType))
	balances := respType.(*banktypes.QueryAllBalancesResponse)
	s.Require().Equal("100000000", balances.Balances[0].Amount.String())
	s.Require().Equal("399986975", balances.Balances[2].Amount.String())

	var account authtypes.AccountI
	respType = proto.Message(&codectype.Any{})
	out, err = simapp.QueryAccountExec(clientCtx, from.String())
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), respType))
	err = clientCtx.InterfaceRegistry.UnpackAny(respType.(*codectype.Any), &account)
	s.Require().NoError(err)

	// test add liquidity (poor not exist)
	status, err := clientCtx.Client.Status(context.Background())
	s.Require().NoError(err)
	deadline := status.SyncInfo.LatestBlockTime.Add(time.Minute)

	txConfig := legacytx.StdTxConfig{Cdc: s.cfg.LegacyAmino}
	msgAddLiquidity := &coinswaptypes.MsgAddLiquidity{
		MaxToken:         sdk.NewCoin(symbol, sdk.NewInt(1000)),
		ExactStandardAmt: sdk.NewInt(1000),
		MinLiquidity:     sdk.NewInt(1000),
		Deadline:         deadline.Unix(),
		Sender:           from.String(),
	}

	// prepare txBuilder with msg
	txBuilder := txConfig.NewTxBuilder()
	feeAmount := sdk.Coins{sdk.NewInt64Coin(s.cfg.BondDenom, 10)}
	err = txBuilder.SetMsgs(msgAddLiquidity)
	s.Require().NoError(err)
	txBuilder.SetFeeAmount(feeAmount)
	txBuilder.SetGasLimit(1000000)

	// setup txFactory
	txFactory := tx.Factory{}.
		WithChainID(val.ClientCtx.ChainID).
		WithKeybase(val.ClientCtx.Keyring).
		WithTxConfig(txConfig).
		WithSignMode(signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON).
		WithSequence(account.GetSequence())

	// sign Tx (offline mode so we can manually set sequence number)
	err = authclient.SignTx(txFactory, val.ClientCtx, val.Moniker, txBuilder, false, true)
	s.Require().NoError(err)

	stdTx := txBuilder.GetTx().(legacytx.StdTx)
	req := authrest.BroadcastReq{
		Tx:   stdTx,
		Mode: "block",
	}
	reqBz, err := val.ClientCtx.LegacyAmino.MarshalJSON(req)
	s.Require().NoError(err)
	_, err = rest.PostRequest(fmt.Sprintf("%s/txs", baseURL), "application/json", reqBz)
	s.Require().NoError(err)

	respType = proto.Message(&banktypes.QueryAllBalancesResponse{})
	out, err = simapp.QueryBalancesExec(clientCtx, from.String())
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), respType))
	balances = respType.(*banktypes.QueryAllBalancesResponse)
	s.Require().Equal("99999000", balances.Balances[0].Amount.String())
	s.Require().Equal("399985965", balances.Balances[2].Amount.String())
	s.Require().Equal("1000", balances.Balances[3].Amount.String())

	url := fmt.Sprintf("%s/coinswap/liquidities/%s", baseURL, uniKitty)
	resp, err := rest.GetRequest(url)
	s.Require().NoError(err)
	s.Require().Equal("1000", gjson.Get(string(resp), "result.standard.amount").String())
	s.Require().Equal("1000", gjson.Get(string(resp), "result.token.amount").String())
	s.Require().Equal("1000", gjson.Get(string(resp), "result.liquidity.amount").String())

	// test add liquidity (poor exist)
	status, err = clientCtx.Client.Status(context.Background())
	s.Require().NoError(err)
	deadline = status.SyncInfo.LatestBlockTime.Add(time.Minute)

	txConfig = legacytx.StdTxConfig{Cdc: s.cfg.LegacyAmino}
	msgAddLiquidity = &coinswaptypes.MsgAddLiquidity{
		MaxToken:         sdk.NewCoin(symbol, sdk.NewInt(2001)),
		ExactStandardAmt: sdk.NewInt(2000),
		MinLiquidity:     sdk.NewInt(2000),
		Deadline:         deadline.Unix(),
		Sender:           from.String(),
	}

	// prepare txBuilder with msg
	txBuilder = txConfig.NewTxBuilder()
	feeAmount = sdk.Coins{sdk.NewInt64Coin(s.cfg.BondDenom, 10)}
	err = txBuilder.SetMsgs(msgAddLiquidity)
	s.Require().NoError(err)
	txBuilder.SetFeeAmount(feeAmount)
	txBuilder.SetGasLimit(1000000)

	// setup txFactory
	txFactory = tx.Factory{}.
		WithChainID(val.ClientCtx.ChainID).
		WithKeybase(val.ClientCtx.Keyring).
		WithTxConfig(txConfig).
		WithSignMode(signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON).
		WithSequence(account.GetSequence() + 1)

	// sign Tx (offline mode so we can manually set sequence number)
	err = authclient.SignTx(txFactory, val.ClientCtx, val.Moniker, txBuilder, false, true)
	s.Require().NoError(err)

	stdTx = txBuilder.GetTx().(legacytx.StdTx)
	req = authrest.BroadcastReq{
		Tx:   stdTx,
		Mode: "block",
	}
	reqBz, err = val.ClientCtx.LegacyAmino.MarshalJSON(req)
	s.Require().NoError(err)
	_, err = rest.PostRequest(fmt.Sprintf("%s/txs", baseURL), "application/json", reqBz)
	s.Require().NoError(err)

	respType = proto.Message(&banktypes.QueryAllBalancesResponse{})
	out, err = simapp.QueryBalancesExec(clientCtx, from.String())
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), respType))
	balances = respType.(*banktypes.QueryAllBalancesResponse)
	s.Require().Equal("99996999", balances.Balances[0].Amount.String())
	s.Require().Equal("399983955", balances.Balances[2].Amount.String())
	s.Require().Equal("3000", balances.Balances[3].Amount.String())

	url = fmt.Sprintf("%s/coinswap/liquidities/%s", baseURL, uniKitty)
	resp, err = rest.GetRequest(url)
	s.Require().NoError(err)
	s.Require().Equal("3000", gjson.Get(string(resp), "result.standard.amount").String())
	s.Require().Equal("3001", gjson.Get(string(resp), "result.token.amount").String())
	s.Require().Equal("3000", gjson.Get(string(resp), "result.liquidity.amount").String())

	// test sell order
	msgSellOrder := &coinswaptypes.MsgSwapOrder{
		Input: coinswaptypes.Input{
			Address: from.String(),
			Coin:    sdk.NewCoin(symbol, sdk.NewInt(1000)),
		},
		Output: coinswaptypes.Output{
			Address: from.String(),
			Coin:    sdk.NewInt64Coin(s.cfg.BondDenom, 748),
		},
		Deadline:   deadline.Unix(),
		IsBuyOrder: false,
	}

	// prepare txBuilder with msg
	txBuilder = txConfig.NewTxBuilder()
	feeAmount = sdk.Coins{sdk.NewInt64Coin(s.cfg.BondDenom, 10)}
	err = txBuilder.SetMsgs(msgSellOrder)
	s.Require().NoError(err)
	txBuilder.SetFeeAmount(feeAmount)
	txBuilder.SetGasLimit(1000000)

	// setup txFactory
	txFactory = tx.Factory{}.
		WithChainID(val.ClientCtx.ChainID).
		WithKeybase(val.ClientCtx.Keyring).
		WithTxConfig(txConfig).
		WithSignMode(signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON).
		WithSequence(account.GetSequence() + 2)

	// sign Tx (offline mode so we can manually set sequence number)
	err = authclient.SignTx(txFactory, val.ClientCtx, val.Moniker, txBuilder, false, true)
	s.Require().NoError(err)

	stdTx = txBuilder.GetTx().(legacytx.StdTx)
	req = authrest.BroadcastReq{
		Tx:   stdTx,
		Mode: "block",
	}
	reqBz, err = val.ClientCtx.LegacyAmino.MarshalJSON(req)
	s.Require().NoError(err)
	_, err = rest.PostRequest(fmt.Sprintf("%s/txs", baseURL), "application/json", reqBz)
	s.Require().NoError(err)

	respType = proto.Message(&banktypes.QueryAllBalancesResponse{})
	out, err = simapp.QueryBalancesExec(clientCtx, from.String())
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), respType))
	balances = respType.(*banktypes.QueryAllBalancesResponse)
	s.Require().Equal("99995999", balances.Balances[0].Amount.String())
	s.Require().Equal("399984693", balances.Balances[2].Amount.String())
	s.Require().Equal("3000", balances.Balances[3].Amount.String())

	url = fmt.Sprintf("%s/coinswap/liquidities/%s", baseURL, uniKitty)
	resp, err = rest.GetRequest(url)
	s.Require().NoError(err)
	s.Require().Equal("2252", gjson.Get(string(resp), "result.standard.amount").String())
	s.Require().Equal("4001", gjson.Get(string(resp), "result.token.amount").String())
	s.Require().Equal("3000", gjson.Get(string(resp), "result.liquidity.amount").String())

	// test buy order
	msgBuyOrder := &coinswaptypes.MsgSwapOrder{
		Input: coinswaptypes.Input{
			Address: from.String(),
			Coin:    sdk.NewInt64Coin(s.cfg.BondDenom, 753),
		},
		Output: coinswaptypes.Output{
			Address: from.String(),
			Coin:    sdk.NewCoin(symbol, sdk.NewInt(1000)),
		},
		Deadline:   deadline.Unix(),
		IsBuyOrder: true,
	}

	// prepare txBuilder with msg
	txBuilder = txConfig.NewTxBuilder()
	feeAmount = sdk.Coins{sdk.NewInt64Coin(s.cfg.BondDenom, 10)}
	err = txBuilder.SetMsgs(msgBuyOrder)
	s.Require().NoError(err)
	txBuilder.SetFeeAmount(feeAmount)
	txBuilder.SetGasLimit(1000000)

	// setup txFactory
	txFactory = tx.Factory{}.
		WithChainID(val.ClientCtx.ChainID).
		WithKeybase(val.ClientCtx.Keyring).
		WithTxConfig(txConfig).
		WithSignMode(signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON).
		WithSequence(account.GetSequence() + 3)

	// sign Tx (offline mode so we can manually set sequence number)
	err = authclient.SignTx(txFactory, val.ClientCtx, val.Moniker, txBuilder, false, true)
	s.Require().NoError(err)

	stdTx = txBuilder.GetTx().(legacytx.StdTx)
	req = authrest.BroadcastReq{
		Tx:   stdTx,
		Mode: "block",
	}
	reqBz, err = val.ClientCtx.LegacyAmino.MarshalJSON(req)
	s.Require().NoError(err)
	_, err = rest.PostRequest(fmt.Sprintf("%s/txs", baseURL), "application/json", reqBz)
	s.Require().NoError(err)

	respType = proto.Message(&banktypes.QueryAllBalancesResponse{})
	out, err = simapp.QueryBalancesExec(clientCtx, from.String())
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), respType))
	balances = respType.(*banktypes.QueryAllBalancesResponse)
	s.Require().Equal("99996999", balances.Balances[0].Amount.String())
	s.Require().Equal("399983930", balances.Balances[2].Amount.String())
	s.Require().Equal("3000", balances.Balances[3].Amount.String())

	url = fmt.Sprintf("%s/coinswap/liquidities/%s", baseURL, uniKitty)
	resp, err = rest.GetRequest(url)
	s.Require().NoError(err)
	s.Require().Equal("3005", gjson.Get(string(resp), "result.standard.amount").String())
	s.Require().Equal("3001", gjson.Get(string(resp), "result.token.amount").String())
	s.Require().Equal("3000", gjson.Get(string(resp), "result.liquidity.amount").String())

	// Test remove liquidity (remove part)
	msgRemoveLiquidity := &coinswaptypes.MsgRemoveLiquidity{
		WithdrawLiquidity: sdk.NewCoin(uniKitty, sdk.NewInt(2000)),
		MinToken:          sdk.NewInt(2000),
		MinStandardAmt:    sdk.NewInt(2000),
		Deadline:          deadline.Unix(),
		Sender:            from.String(),
	}

	// prepare txBuilder with msg
	txBuilder = txConfig.NewTxBuilder()
	feeAmount = sdk.Coins{sdk.NewInt64Coin(s.cfg.BondDenom, 10)}
	err = txBuilder.SetMsgs(msgRemoveLiquidity)
	s.Require().NoError(err)
	txBuilder.SetFeeAmount(feeAmount)
	txBuilder.SetGasLimit(1000000)

	// setup txFactory
	txFactory = tx.Factory{}.
		WithChainID(val.ClientCtx.ChainID).
		WithKeybase(val.ClientCtx.Keyring).
		WithTxConfig(txConfig).
		WithSignMode(signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON).
		WithSequence(account.GetSequence() + 4)

	// sign Tx (offline mode so we can manually set sequence number)
	err = authclient.SignTx(txFactory, val.ClientCtx, val.Moniker, txBuilder, false, true)
	s.Require().NoError(err)

	stdTx = txBuilder.GetTx().(legacytx.StdTx)
	req = authrest.BroadcastReq{
		Tx:   stdTx,
		Mode: "block",
	}
	reqBz, err = val.ClientCtx.LegacyAmino.MarshalJSON(req)
	s.Require().NoError(err)
	_, err = rest.PostRequest(fmt.Sprintf("%s/txs", baseURL), "application/json", reqBz)
	s.Require().NoError(err)

	respType = proto.Message(&banktypes.QueryAllBalancesResponse{})
	out, err = simapp.QueryBalancesExec(clientCtx, from.String())
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), respType))
	balances = respType.(*banktypes.QueryAllBalancesResponse)
	s.Require().Equal("99998999", balances.Balances[0].Amount.String())
	s.Require().Equal("399985923", balances.Balances[2].Amount.String())
	s.Require().Equal("1000", balances.Balances[3].Amount.String())

	url = fmt.Sprintf("%s/coinswap/liquidities/%s", baseURL, uniKitty)
	resp, err = rest.GetRequest(url)
	s.Require().NoError(err)
	s.Require().Equal("1002", gjson.Get(string(resp), "result.standard.amount").String())
	s.Require().Equal("1001", gjson.Get(string(resp), "result.token.amount").String())
	s.Require().Equal("1000", gjson.Get(string(resp), "result.liquidity.amount").String())

	// Test remove liquidity (remove all)
	msgRemoveLiquidity = &coinswaptypes.MsgRemoveLiquidity{
		WithdrawLiquidity: sdk.NewCoin(uniKitty, sdk.NewInt(1000)),
		MinToken:          sdk.NewInt(1000),
		MinStandardAmt:    sdk.NewInt(1000),
		Deadline:          deadline.Unix(),
		Sender:            from.String(),
	}

	// prepare txBuilder with msg
	txBuilder = txConfig.NewTxBuilder()
	feeAmount = sdk.Coins{sdk.NewInt64Coin(s.cfg.BondDenom, 10)}
	err = txBuilder.SetMsgs(msgRemoveLiquidity)
	s.Require().NoError(err)
	txBuilder.SetFeeAmount(feeAmount)
	txBuilder.SetGasLimit(1000000)

	// setup txFactory
	txFactory = tx.Factory{}.
		WithChainID(val.ClientCtx.ChainID).
		WithKeybase(val.ClientCtx.Keyring).
		WithTxConfig(txConfig).
		WithSignMode(signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON).
		WithSequence(account.GetSequence() + 5)

	// sign Tx (offline mode so we can manually set sequence number)
	err = authclient.SignTx(txFactory, val.ClientCtx, val.Moniker, txBuilder, false, true)
	s.Require().NoError(err)

	stdTx = txBuilder.GetTx().(legacytx.StdTx)
	req = authrest.BroadcastReq{
		Tx:   stdTx,
		Mode: "block",
	}
	reqBz, err = val.ClientCtx.LegacyAmino.MarshalJSON(req)
	s.Require().NoError(err)
	_, err = rest.PostRequest(fmt.Sprintf("%s/txs", baseURL), "application/json", reqBz)
	s.Require().NoError(err)

	respType = proto.Message(&banktypes.QueryAllBalancesResponse{})
	out, err = simapp.QueryBalancesExec(clientCtx, from.String())
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), respType))
	balances = respType.(*banktypes.QueryAllBalancesResponse)
	s.Require().Equal("100000000", balances.Balances[0].Amount.String())
	s.Require().Equal("399986915", balances.Balances[2].Amount.String())
	s.Require().Equal("0", balances.Balances[3].Amount.String())

	url = fmt.Sprintf("%s/coinswap/liquidities/%s", baseURL, uniKitty)
	resp, err = rest.GetRequest(url)
	s.Require().NoError(err)
	s.Require().Equal("0", gjson.Get(string(resp), "result.standard.amount").String())
	s.Require().Equal("0", gjson.Get(string(resp), "result.token.amount").String())
	s.Require().Equal("0", gjson.Get(string(resp), "result.liquidity.amount").String())
}
