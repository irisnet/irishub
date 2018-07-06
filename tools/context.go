package tools

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"strings"
	"net/http"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

)

type Context struct {
	context.CoreContext
}

func NewContext() Context {
	ctx := context.NewCoreContextFromViper()
	return Context{
		ctx,
	}
}

type JsonRpc interface {
	NetInfo() *ctypes.ResultNetInfo
	NumUnconfirmedTxs() *ctypes.ResultUnconfirmedTxs
}

func (rpc Context) NetInfo() *ctypes.ResultNetInfo{
	client := &http.Client{}

	reqUri := tcpToHttpUrl(rpc.NodeURI) + "/net_info"

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
		JsonRpc string 				 `json:"jsonrpc"`
		Id 		string 				 `json:"id"`
		Result 	ctypes.ResultNetInfo `json:"result"`
	}{}
	if err := json.Unmarshal(body,&res); err != nil {
		fmt.Println(err)
	}

	return &res.Result
}

func (rpc Context) NumUnconfirmedTxs() *ctypes.ResultUnconfirmedTxs{
	client := &http.Client{}
	reqUri := tcpToHttpUrl(rpc.NodeURI) + "/num_unconfirmed_txs"

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
		JsonRpc string 						`json:"jsonrpc"`
		Id 		string 						`json:"id"`
		Result 	ctypes.ResultUnconfirmedTxs `json:"result"`
	}{}
	if err := json.Unmarshal(body,&res); err != nil {
		fmt.Println(err)
	}

	return &res.Result
}

func tcpToHttpUrl(url string) string{
	urls := strings.Replace(url,"tcp","http",1)
	return urls
}


