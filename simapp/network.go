package simapp

import (
	"context"
	"encoding/hex"
	"strings"
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clienttx "github.com/cosmos/cosmos-sdk/client/tx"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	"github.com/cosmos/gogoproto/proto"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

type Network struct {
	*network.Network
	network.Config
}

type ResponseTx struct {
	abci.ResponseDeliverTx
	Height int64
}

func SetupNetwork(t *testing.T, depInjectOptions DepinjectOptions) Network {
	t.Helper()
	cfg, err := NewConfig(depInjectOptions)
	require.NoError(t, err)

	cfg.NumValidators = 4

	network, err := network.New(t, t.TempDir(), cfg)
	require.NoError(t, err, "SetupNetwork failed")

	n := Network{
		Network: network,
		Config:  cfg,
	}
	require.NoError(t, n.WaitForNBlock(2), "WaitForNBlock failed")
	return n
}

func SetupNetworkWithConfig(t *testing.T, cfg network.Config) Network {
	t.Helper()
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
	t.Helper()
	buf, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, extraArgs)
	require.NoError(t, err, "ExecTestCLICmd failed")

	require.NoError(t, n.WaitForNextBlock(), "WaitForNextBlock failed")

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
	t.Helper()
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
	t.Helper()
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
		require.NoError(t, n.WaitForNextBlock(), "WaitForNextBlock failed")

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

func (n Network) SendMsgs(
	t *testing.T,
	msgs ...sdk.Msg,
) *sdk.TxResponse {
	t.Helper()
	val := n.Validators[0]
	client := val.ClientCtx.WithBroadcastMode(flags.BroadcastSync)

	// prepare txBuilder with msg
	txBuilder := client.TxConfig.NewTxBuilder()
	err := txBuilder.SetMsgs(msgs...)
	require.NoError(t, err, "txBuilder.SetMsgs failed")

	txBuilder.SetFeeAmount(sdk.Coins{sdk.NewInt64Coin(n.BondDenom, 10)})
	txBuilder.SetGasLimit(1000000)
	// setup txFactory
	txFactory := clienttx.Factory{}.
		WithChainID(client.ChainID).
		WithKeybase(client.Keyring).
		WithTxConfig(client.TxConfig).
		WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)

	// Sign Tx.
	err = authclient.SignTx(txFactory, client, val.Moniker, txBuilder, false, true)
	require.NoError(t, err, "authclient.SignTx failed")

	txBytes, err := client.TxConfig.TxEncoder()(txBuilder.GetTx())
	require.NoError(t, err, "TxConfig.TxEncoder failed")
	res, err := client.BroadcastTx(txBytes)
	require.NoError(t, err, "BroadcastTx failed")
	require.Equal(t, uint32(0), res.Code, res.RawLog)
	require.NoError(t, n.WaitForNBlock(2), "WaitForNextBlock failed")
	return res
}

func (n Network) BlockSendMsgs(t *testing.T,
	msgs ...sdk.Msg,
) *ResponseTx {
	t.Helper()
	response := n.SendMsgs(t, msgs...)
	return n.QueryTx(t, n.Validators[0].ClientCtx, response.TxHash)
}
