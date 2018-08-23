package app

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"io/ioutil"
	"net/http"
	"strings"
	"github.com/irisnet/irishub/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Context struct {
	context.CLIContext
	Cdc *wire.Codec
}

func NewContext() Context {
	return Context{
		CLIContext:context.NewCLIContext(),
	}
}
func (c Context) Get() context.CLIContext {
	return c.CLIContext
}

func (c Context) WithCodeC(cdc *wire.Codec) Context {
	c.Cdc = cdc
	return c
}

func (c Context) BroadcastTxAsync(tx []byte) (*ctypes.ResultBroadcastTx, error) {
	return c.Client.BroadcastTxAsync(tx)
}

func (c Context) BroadcastTxSync(tx []byte) (*ctypes.ResultBroadcastTx, error) {
	return c.Client.BroadcastTxSync(tx)
}

func (c Context) NetInfo() (*ctypes.ResultNetInfo, error) {
	client := &http.Client{}

	reqUri := tcpToHttpUrl(c.NodeURI) + "/net_info"

	resp, err := client.Get(reqUri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res = struct {
		JsonRpc string               `json:"jsonrpc"`
		Id      string               `json:"id"`
		Result  ctypes.ResultNetInfo `json:"result"`
	}{}
	if err := c.Cdc.UnmarshalJSON(body, &res); err != nil {
		return nil, err
	}

	return &res.Result, nil
}

func (c Context) NumUnconfirmedTxs() (*ctypes.ResultUnconfirmedTxs, error){
	client := &http.Client{}
	reqUri := tcpToHttpUrl(c.NodeURI) + "/num_unconfirmed_txs"

	resp, err := client.Get(reqUri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res = struct {
		JsonRpc string                      `json:"jsonrpc"`
		Id      string                      `json:"id"`
		Result  ctypes.ResultUnconfirmedTxs `json:"result"`
	}{}

	if err := c.Cdc.UnmarshalJSON(body, &res); err != nil {
		return nil, err
	}

	return &res.Result, nil
}

func (c Context) GetCoinType(coinName string, cdc *wire.Codec) (types.CoinType, error) {
	var coinType types.CoinType
	if strings.ToLower(coinName) == denom {
		coinType = types.NewDefaultCoinType(denom)
	}else{
		key := types.CoinTypeKey(coinName)
		bz,err := c.QueryStore([]byte(key),"iparams")
		if err != nil {
			return coinType,err
		}

		if err = cdc.UnmarshalBinary(bz,&coinType);err != nil {
			return coinType,err
		}
	}

	return coinType, nil
}

func (c Context) ParseCoin(coinStr string, cdc *wire.Codec) (sdk.Coin, error) {
	mainUnit,err := types.GetCoinName(coinStr)
	coinType,err := c.GetCoinType(mainUnit,cdc)
	if err != nil {
		return sdk.Coin{},err
	}

	coin,err:=coinType.ConvertToMinCoin(coinStr)
	if err != nil {
		return sdk.Coin{},err
	}
	return coin,nil
}

func (c Context) ParseCoins(coinsStr string, cdc *wire.Codec) (coins sdk.Coins, err error) {
	coinsStr = strings.TrimSpace(coinsStr)
	if len(coinsStr) == 0 {
		return coins, nil
	}

	coinStrs := strings.Split(coinsStr, ",")
	for _, coinStr := range coinStrs {
		coin, err := c.ParseCoin(coinStr,cdc)
		if err != nil {
			return coins, err
		}
		coins = append(coins, coin)
	}
	return coins,nil
}

func tcpToHttpUrl(url string) string {
	urls := strings.Replace(url, "tcp", "http", 1)
	return urls
}
