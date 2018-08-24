package app

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	txcxt "github.com/cosmos/cosmos-sdk/x/auth/client/context"
	"github.com/irisnet/irishub/types"
	"github.com/pkg/errors"
	client2 "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"io/ioutil"
	"net/http"
	"strings"
)

type Context struct {
	context.CLIContext
	txCtx txcxt.TxContext
	Cdc   *wire.Codec
}

func NewContext() Context {
	return Context{
		CLIContext: context.NewCLIContext(),
		txCtx:      txcxt.NewTxContextFromCLI(),
	}
}
func (c Context) Get() context.CLIContext {
	return c.CLIContext.WithCodec(c.Cdc)
}

func (c Context) GetTxCxt() txcxt.TxContext {
	return c.txCtx.WithCodec(c.Cdc)
}

func (c Context) WithCodeC(cdc *wire.Codec) Context {
	c.Cdc = cdc
	return c
}
func (c Context) WithCLIContext(ctx context.CLIContext) Context {
	c.CLIContext = ctx
	return c
}
func (c Context) WithTxContext(ctx txcxt.TxContext) Context {
	c.txCtx = ctx
	return c
}

func (c Context) NetInfo() (*ctypes.ResultNetInfo, error) {
	client := c.Client.(*client2.HTTP)
	return client.NetInfo()
}

func (c Context) NumUnconfirmedTxs() (*ctypes.ResultUnconfirmedTxs, error) {
	client := &http.Client{}
	url := strings.Replace(c.NodeURI, "tcp", "http", 1)
	reqUri := fmt.Sprintf("%s/%s", url, "num_unconfirmed_txs")

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

func (c Context) GetCoinType(coinName string) (types.CoinType, error) {
	var coinType types.CoinType
	if strings.ToLower(coinName) == denom {
		coinType = IrisCt
	} else {
		key := types.CoinTypeKey(coinName)
		bz, err := c.QueryStore([]byte(key), "iparams")
		if err != nil {
			return coinType, err
		}

		if err = c.Cdc.UnmarshalBinary(bz, &coinType); err != nil {
			return coinType, err
		}
	}

	return coinType, nil
}

func (c Context) ParseCoin(coinStr string) (sdk.Coin, error) {
	mainUnit, err := types.GetCoinName(coinStr)
	coinType, err := c.GetCoinType(mainUnit)
	if err != nil {
		return sdk.Coin{}, err
	}

	coin, err := coinType.ConvertToMinCoin(coinStr)
	if err != nil {
		return sdk.Coin{}, err
	}
	return coin, nil
}

func (c Context) ParseCoins(coinsStr string) (coins sdk.Coins, err error) {
	coinsStr = strings.TrimSpace(coinsStr)
	if len(coinsStr) == 0 {
		return coins, nil
	}

	coinStrs := strings.Split(coinsStr, ",")
	for _, coinStr := range coinStrs {
		coin, err := c.ParseCoin(coinStr)
		if err != nil {
			return coins, err
		}
		coins = append(coins, coin)
	}
	return coins, nil
}

// Build builds a single message to be signed from a TxContext given a set of
// messages. It returns an error if a fee is supplied but cannot be parsed.
func (c Context) Build(msgs []sdk.Msg) (auth.StdSignMsg, error) {
	ctx := c.txCtx
	chainID := ctx.ChainID
	if chainID == "" {
		return auth.StdSignMsg{}, errors.Errorf("chain ID required but not specified")
	}

	fee := sdk.Coin{}
	if ctx.Fee != "" {
		parsedFee, err := c.ParseCoin(ctx.Fee)
		if err != nil {
			return auth.StdSignMsg{}, err
		}

		fee = parsedFee
	}

	return auth.StdSignMsg{
		ChainID:       ctx.ChainID,
		AccountNumber: ctx.AccountNumber,
		Sequence:      ctx.Sequence,
		Memo:          ctx.Memo,
		Msgs:          msgs,
		Fee:           auth.NewStdFee(ctx.Gas, fee),
	}, nil
}
