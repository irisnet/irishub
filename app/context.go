package app

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"io/ioutil"
	"strings"
	"net/http"
	"github.com/cosmos/cosmos-sdk/wire"
)

type Context struct {
	Ctx context.CoreContext
	Cdc *wire.Codec
}


func NewContext() Context {
	return Context{
		Ctx:context.NewCoreContextFromViper(),
	}
}

func (c Context) WithCodeC(cdc *wire.Codec)  Context{
	c.Cdc = cdc
	return c
}

func (c Context) BroadcastTxAsync(tx []byte) (*ctypes.ResultBroadcastTx, error) {
	return c.Ctx.Client.BroadcastTxAsync(tx)
}

func (c Context) BroadcastTxSync(tx []byte) (*ctypes.ResultBroadcastTx, error) {
	return c.Ctx.Client.BroadcastTxSync(tx)
}

func (c Context) NetInfo() (*ctypes.ResultNetInfo, error) {
	client := &http.Client{}

	reqUri := tcpToHttpUrl(c.Ctx.NodeURI) + "/net_info"

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
	reqUri := tcpToHttpUrl(c.Ctx.NodeURI) + "/num_unconfirmed_txs"

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

func tcpToHttpUrl(url string) string {
	urls := strings.Replace(url, "tcp", "http", 1)
	return urls
}
