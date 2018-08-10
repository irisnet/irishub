package app

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"fmt"
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

func (c Context) NetInfo() *ctypes.ResultNetInfo {
	client := &http.Client{}

	reqUri := tcpToHttpUrl(c.Ctx.NodeURI) + "/net_info"

	resp, err := client.Get(reqUri)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var res = struct {
		JsonRpc string               `json:"jsonrpc"`
		Id      string               `json:"id"`
		Result  ctypes.ResultNetInfo `json:"result"`
	}{}
	if err := c.Cdc.UnmarshalJSON(body,&res); err != nil {
		fmt.Println(err)
	}

	return &res.Result
}

func (c Context) NumUnconfirmedTxs() *ctypes.ResultUnconfirmedTxs {
	client := &http.Client{}
	reqUri := tcpToHttpUrl(c.Ctx.NodeURI) + "/num_unconfirmed_txs"

	resp, err := client.Get(reqUri)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var res = struct {
		JsonRpc string                      `json:"jsonrpc"`
		Id      string                      `json:"id"`
		Result  ctypes.ResultUnconfirmedTxs `json:"result"`
	}{}
	if err := c.Cdc.UnmarshalJSON(body,&res); err != nil {
		fmt.Println(err)
	}

	return &res.Result
}

func tcpToHttpUrl(url string) string {
	urls := strings.Replace(url, "tcp", "http", 1)
	return urls
}
