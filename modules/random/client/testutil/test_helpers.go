package testutil

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/cometbft/cometbft/libs/cli"

// 	"github.com/cosmos/cosmos-sdk/client"
// 	"github.com/cosmos/cosmos-sdk/client/flags"

// 	"github.com/irisnet/irismod/simapp"
// 	randomcli "github.com/irisnet/irismod/random/client/cli"
// 	randomtypes "github.com/irisnet/irismod/random/types"
// )

// // MsgRedelegateExec creates a redelegate message.
// func RequestRandomExec(t *testing.T,
// 	network simapp.Network,
// 	clientCtx client.Context,
// 	from string,
// 	extraArgs ...string,
// ) *simapp.ResponseTx {
// 	args := []string{
// 		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
// 	}
// 	args = append(args, extraArgs...)

// 	return network.ExecTxCmdWithResult(t, clientCtx, randomcli.GetCmdRequestRandom(), args)
// }

// func QueryRandomExec(t *testing.T,
// 	network simapp.Network,
// 	clientCtx client.Context,
// 	requestID string,
// 	extraArgs ...string) *randomtypes.Random {
// 	args := []string{
// 		requestID,
// 		fmt.Sprintf("--%s=json", cli.OutputFlag),
// 	}
// 	args = append(args, extraArgs...)

// 	response := &randomtypes.Random{}
// 	network.ExecQueryCmd(t, clientCtx, randomcli.GetCmdQueryRandom(), args, response)
// 	return response
// }

// func QueryRandomRequestQueueExec(t *testing.T,
// 	network simapp.Network,
// 	clientCtx client.Context,
// 	genHeight string,
// 	extraArgs ...string) *randomtypes.QueryRandomRequestQueueResponse {
// 	args := []string{
// 		genHeight,
// 		fmt.Sprintf("--%s=json", cli.OutputFlag),
// 	}
// 	args = append(args, extraArgs...)

// 	response := &randomtypes.QueryRandomRequestQueueResponse{}
// 	network.ExecQueryCmd(t, clientCtx, randomcli.GetCmdQueryRandomRequestQueue(), args, response)
// 	return response
// }
