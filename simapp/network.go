package simapp

import (
	"context"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/client"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
)

type Network struct {
	*network.Network
	network.Config
}

type ResponseTx struct {
	abci.ResponseDeliverTx
	Height int64
}

func SetupNetwork(t *testing.T) Network {
	cfg := NewConfig()
	cfg.NumValidators = 4

	network, err := network.New(t, t.TempDir(), cfg)
	require.NoError(t, err, "SetupNetwork failed")

	n := Network{
		Network: network,
		Config:  cfg,
	}
	n.WaitForNBlock(2)
	return n
}

func SetupNetworkWithConfig(t *testing.T, cfg network.Config) Network {
	network, err := network.New(t, t.TempDir(), cfg)
	require.NoError(t, err, "SetupNetwork failed")

	_, err = network.WaitForHeight(1)
	require.NoError(t, err)
	return Network{
		Network: network,
		Config:  cfg,
	}
}

func (n Network) ExecTxCmdWithResult(t *testing.T,
	clientCtx client.Context,
	cmd *cobra.Command,
	extraArgs []string,
) *ResponseTx {
	buf, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, extraArgs)
	require.NoError(t, err, "ExecTestCLICmd failed")

	n.WaitForNextBlock()

	respType := proto.Message(&sdk.TxResponse{})
	require.NoError(t, clientCtx.Codec.UnmarshalJSON(buf.Bytes(), respType), buf.String())

	txResp := respType.(*sdk.TxResponse)
	require.Equal(t, uint32(0), txResp.Code)
	return n.QueryTx(t, clientCtx, txResp.TxHash)
}

func (n Network) ExecQueryCmd(t *testing.T,
	clientCtx client.Context,
	cmd *cobra.Command,
	extraArgs []string,
	resp proto.Message,
) {
	buf, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, extraArgs)
	require.NoError(t, err, "ExecTestCLICmd failed")
	require.NoError(t, clientCtx.Codec.UnmarshalJSON(buf.Bytes(), resp), buf.String())
}

func (n Network) WaitForNBlock(wait int64) error {
	lastBlock, err := n.LatestHeight()
	if err != nil {
		return err
	}

	_, err = n.WaitForHeight(lastBlock + wait)
	if err != nil {
		return err
	}

	return err
}

func (n Network) QueryTx(t *testing.T,
	clientCtx client.Context,
	txHash string,
) *ResponseTx {
	var (
		result *coretypes.ResultTx
		err    error
		tryCnt = 3
	)

	txHashBz, err := hex.DecodeString(txHash)
	require.NoError(t, err, "hex.DecodeString failed")

reTry:
	result, err = clientCtx.Client.Tx(context.Background(), txHashBz, false)
	if err != nil && strings.Contains(err.Error(), "not found") && tryCnt > 0 {
		n.WaitForNextBlock()
		tryCnt--
		goto reTry
	}

	require.NoError(t, err, "query tx failed")
	return &ResponseTx{result.TxResult, result.Height}
}

func (n Network) GetAttribute(typ, key string, events []abci.Event) string {
	for _, event := range events {
		if event.Type == typ {
			for _, attribute := range event.Attributes {
				if attribute.Key == key {
					return attribute.Value
				}
			}
		}
	}
	return ""
}
