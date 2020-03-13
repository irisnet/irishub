//nolint:bodyclose
package lcdtest

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/client/flags"
	crkeys "github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/tests"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/irisnet/irishub/modules/coinswap"
	coinswaprest "github.com/irisnet/irishub/modules/coinswap/client/rest"
)

func TestCoinswap(t *testing.T) {
	name := "sender"
	denomStandard := sdk.DefaultBondDenom
	denomBTC := "btc"
	denomETH := "eth"
	uniBTC := "uni:btc"
	uniETH := "uni:eth"
	deadline := "10m0s"
	assets := []string{denomBTC, denomETH}

	kb, err := newKeybase()
	require.NoError(t, err)
	addr, _, err := CreateAddr(name, kb)
	require.NoError(t, err)

	cleanup, _, _, port, err := InitializeLCD(1, []sdk.AccAddress{addr}, true, assets)
	require.NoError(t, err)
	defer cleanup()

	acc := getAccount(t, port, addr)
	require.Equal(t, "100000000", acc.GetCoins().AmountOf(sdk.DefaultBondDenom).String())
	require.Equal(t, "100000000", acc.GetCoins().AmountOf(denomBTC).String())
	require.Equal(t, "100000000", acc.GetCoins().AmountOf(denomETH).String())
	require.Equal(t, "0", acc.GetCoins().AmountOf(uniBTC).String())
	require.Equal(t, "0", acc.GetCoins().AmountOf(uniETH).String())

	// add liquidity
	resultTx := addLiquidity(
		t, port, name, kb, addr, uniBTC,
		"2000", "2000", "2000", deadline,
	)
	tests.WaitForHeight(resultTx.Height+1, port)

	acc = getAccount(t, port, addr)
	require.Equal(t, "99997995", acc.GetCoins().AmountOf(sdk.DefaultBondDenom).String())
	require.Equal(t, "99998000", acc.GetCoins().AmountOf(denomBTC).String())
	require.Equal(t, "2000", acc.GetCoins().AmountOf(uniBTC).String())

	// query liquidity
	liquidityBTC := queryLiquidity(t, port, uniBTC)
	require.Equal(t, "2000stake", liquidityBTC.Standard.String())
	require.Equal(t, "2000btc", liquidityBTC.Token.String())
	require.Equal(t, "2000uni:btc", liquidityBTC.Liquidity.String())
	require.Equal(t, "0.003000000000000000", liquidityBTC.Fee)

	// add liquidity
	resultTx = addLiquidity(
		t, port, name, kb, addr, uniETH,
		"1000", "1000", "1000", deadline,
	)
	tests.WaitForHeight(resultTx.Height+1, port)

	acc = getAccount(t, port, addr)
	require.Equal(t, "99996990", acc.GetCoins().AmountOf(sdk.DefaultBondDenom).String())
	require.Equal(t, "99999000", acc.GetCoins().AmountOf(denomETH).String())
	require.Equal(t, "1000", acc.GetCoins().AmountOf(uniETH).String())

	// query liquidity
	liquidityETH := queryLiquidity(t, port, uniETH)
	require.Equal(t, "1000stake", liquidityETH.Standard.String())
	require.Equal(t, "1000eth", liquidityETH.Token.String())
	require.Equal(t, "1000uni:eth", liquidityETH.Liquidity.String())
	require.Equal(t, "0.003000000000000000", liquidityETH.Fee)

	// remove liquidity
	resultTx = removeLiquidity(
		t, port, name, kb, addr, uniBTC,
		"1000", "1000", "1000", deadline,
	)
	tests.WaitForHeight(resultTx.Height+1, port)

	acc = getAccount(t, port, addr)
	require.Equal(t, "99997985", acc.GetCoins().AmountOf(sdk.DefaultBondDenom).String())
	require.Equal(t, "99999000", acc.GetCoins().AmountOf(denomBTC).String())
	require.Equal(t, "1000", acc.GetCoins().AmountOf(uniBTC).String())

	// buy order
	resultTx = buyOrder(
		t, port, name, kb, addr,
		sdk.NewCoin(denomBTC, sdk.NewInt(200)),
		sdk.NewCoin(denomStandard, sdk.NewInt(100)),
		deadline,
	)
	tests.WaitForHeight(resultTx.Height+1, port)

	acc = getAccount(t, port, addr)
	require.Equal(t, "99998080", acc.GetCoins().AmountOf(sdk.DefaultBondDenom).String())
	require.Equal(t, "99998888", acc.GetCoins().AmountOf(denomBTC).String())

	// sell order
	resultTx = sellOrder(
		t, port, name, kb, addr,
		sdk.NewCoin(denomBTC, sdk.NewInt(200)),
		sdk.NewCoin(denomStandard, sdk.NewInt(100)),
		deadline,
	)
	tests.WaitForHeight(resultTx.Height+1, port)

	acc = getAccount(t, port, addr)
	require.Equal(t, "99998211", acc.GetCoins().AmountOf(sdk.DefaultBondDenom).String())
	require.Equal(t, "99998688", acc.GetCoins().AmountOf(denomBTC).String())

	// double swap order
	resultTx = buyOrder(
		t, port, name, kb, addr,
		sdk.NewCoin(denomBTC, sdk.NewInt(644)),
		sdk.NewCoin(denomETH, sdk.NewInt(200)),
		deadline,
	)
	tests.WaitForHeight(resultTx.Height+1, port)

	acc = getAccount(t, port, addr)
	require.Equal(t, "99998206", acc.GetCoins().AmountOf(sdk.DefaultBondDenom).String())
	require.Equal(t, "99998044", acc.GetCoins().AmountOf(denomBTC).String())
	require.Equal(t, "99999200", acc.GetCoins().AmountOf(denomETH).String())

	// query liquidity
	liquidityBTC = queryLiquidity(t, port, uniBTC)
	require.Equal(t, "513stake", liquidityBTC.Standard.String())
	require.Equal(t, "1956btc", liquidityBTC.Token.String())
	require.Equal(t, "1000uni:btc", liquidityBTC.Liquidity.String())
	liquidityETH = queryLiquidity(t, port, uniETH)
	require.Equal(t, "1251stake", liquidityETH.Standard.String())
	require.Equal(t, "800eth", liquidityETH.Token.String())
	require.Equal(t, "1000uni:eth", liquidityETH.Liquidity.String())
}

// POST /coinswap/liquidities/{pool-id}/deposit
func addLiquidity(
	t *testing.T, port string, name string, kb crkeys.Keybase,
	addrSender sdk.AccAddress, poolID string, maxToken string,
	exactStandardAmt string, minLiquidity string, deadline string,
) (txResp sdk.TxResponse) {
	acc := getAccount(t, port, addrSender)
	accnum := acc.GetAccountNumber()
	sequence := acc.GetSequence()
	chainID := viper.GetString(flags.FlagChainID)
	from := acc.GetAddress().String()

	baseReq := rest.NewBaseReq(from, "", chainID, "", "", accnum, sequence, fees, nil, false)
	addLiquidityReq := coinswaprest.AddLiquidityReq{
		BaseTx:           baseReq,
		ID:               poolID,
		MaxToken:         maxToken,
		ExactStandardAmt: exactStandardAmt,
		MinLiquidity:     minLiquidity,
		Deadline:         deadline,
		Sender:           addrSender.String(),
	}

	req, err := cdc.MarshalJSON(addLiquidityReq)
	require.NoError(t, err)

	// generate tx
	res, body := Request(t, port, "POST", fmt.Sprintf("/coinswap/liquidities/%s/deposit", poolID), req)
	require.Equal(t, http.StatusOK, res.StatusCode, body)

	// sign and broadcast
	resp, body := signAndBroadcastGenTx(t, port, name, body, acc, 0, false, kb)
	require.Equal(t, http.StatusOK, resp.StatusCode, body)

	err = cdc.UnmarshalJSON([]byte(body), &txResp)
	require.NoError(t, err)

	return
}

// POST /coinswap/liquidities/{pool-id}/withdraw
func removeLiquidity(
	t *testing.T, port string, name string, kb crkeys.Keybase,
	addrSender sdk.AccAddress, poolID string, minToken string,
	withdrawLiquidity string, minStandardAmt string, deadline string,
) (txResp sdk.TxResponse) {
	acc := getAccount(t, port, addrSender)
	accnum := acc.GetAccountNumber()
	sequence := acc.GetSequence()
	chainID := viper.GetString(flags.FlagChainID)
	from := acc.GetAddress().String()

	baseReq := rest.NewBaseReq(from, "", chainID, "", "", accnum, sequence, fees, nil, false)
	removeLiquidityReq := coinswaprest.RemoveLiquidityReq{
		BaseTx:            baseReq,
		ID:                poolID,
		MinToken:          minToken,
		WithdrawLiquidity: withdrawLiquidity,
		MinStandardAmt:    minStandardAmt,
		Deadline:          deadline,
		Sender:            addrSender.String(),
	}

	req, err := cdc.MarshalJSON(removeLiquidityReq)
	require.NoError(t, err)

	// generate tx
	res, body := Request(t, port, "POST", fmt.Sprintf("/coinswap/liquidities/%s/withdraw", poolID), req)
	require.Equal(t, http.StatusOK, res.StatusCode, body)

	// sign and broadcast
	resp, body := signAndBroadcastGenTx(t, port, name, body, acc, 0, false, kb)
	require.Equal(t, http.StatusOK, resp.StatusCode, body)

	err = cdc.UnmarshalJSON([]byte(body), &txResp)
	require.NoError(t, err)

	return
}

// POST /coinswap/liquidities/buy
func buyOrder(
	t *testing.T, port string, name string, kb crkeys.Keybase,
	addrSender sdk.AccAddress, inputCoin sdk.Coin, outputCoin sdk.Coin,
	deadline string,
) (txResp sdk.TxResponse) {
	acc := getAccount(t, port, addrSender)
	accnum := acc.GetAccountNumber()
	sequence := acc.GetSequence()
	chainID := viper.GetString(flags.FlagChainID)
	from := acc.GetAddress().String()

	baseReq := rest.NewBaseReq(from, "", chainID, "", "", accnum, sequence, fees, nil, false)
	swapOrderReq := coinswaprest.SwapOrderReq{
		BaseTx:   baseReq,
		Input:    coinswaprest.Input{Address: addrSender.String(), Coin: inputCoin},
		Output:   coinswaprest.Output{Address: addrSender.String(), Coin: outputCoin},
		Deadline: deadline,
	}

	req, err := cdc.MarshalJSON(swapOrderReq)
	require.NoError(t, err)

	// generate tx
	res, body := Request(t, port, "POST", "/coinswap/liquidities/buy", req)
	require.Equal(t, http.StatusOK, res.StatusCode, body)

	// sign and broadcast
	resp, body := signAndBroadcastGenTx(t, port, name, body, acc, 0, false, kb)
	require.Equal(t, http.StatusOK, resp.StatusCode, body)

	err = cdc.UnmarshalJSON([]byte(body), &txResp)
	require.NoError(t, err)

	return
}

// POST /coinswap/liquidities/sell
func sellOrder(
	t *testing.T, port string, name string, kb crkeys.Keybase,
	addrSender sdk.AccAddress, inputCoin sdk.Coin, outputCoin sdk.Coin,
	deadline string,
) (txResp sdk.TxResponse) {
	acc := getAccount(t, port, addrSender)
	accnum := acc.GetAccountNumber()
	sequence := acc.GetSequence()
	chainID := viper.GetString(flags.FlagChainID)
	from := acc.GetAddress().String()

	baseReq := rest.NewBaseReq(from, "", chainID, "", "", accnum, sequence, fees, nil, false)
	swapOrderReq := coinswaprest.SwapOrderReq{
		BaseTx:   baseReq,
		Input:    coinswaprest.Input{Address: addrSender.String(), Coin: inputCoin},
		Output:   coinswaprest.Output{Address: addrSender.String(), Coin: outputCoin},
		Deadline: deadline,
	}

	req, err := cdc.MarshalJSON(swapOrderReq)
	require.NoError(t, err)

	// generate tx
	res, body := Request(t, port, "POST", "/coinswap/liquidities/sell", req)
	require.Equal(t, http.StatusOK, res.StatusCode, body)

	// sign and broadcast
	resp, body := signAndBroadcastGenTx(t, port, name, body, acc, 0, false, kb)
	require.Equal(t, http.StatusOK, resp.StatusCode, body)

	err = cdc.UnmarshalJSON([]byte(body), &txResp)
	require.NoError(t, err)

	return
}

// GET /coinswap/liquidities/{pool-id}
func queryLiquidity(t *testing.T, port string, poolID string) (response coinswap.QueryLiquidityResponse) {
	res, body := Request(t, port, "GET", fmt.Sprintf("/coinswap/liquidities/%s", poolID), nil)
	require.Equal(t, http.StatusOK, res.StatusCode, body)

	var resp rest.ResponseWithHeight
	require.NoError(t, cdc.UnmarshalJSON([]byte(body), &resp))

	err := cdc.UnmarshalJSON(resp.Result, &response)
	require.NoError(t, err)

	return
}
