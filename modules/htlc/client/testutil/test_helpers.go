package testutil

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/cometbft/cometbft/libs/cli"

// 	"github.com/cosmos/cosmos-sdk/client"
// 	"github.com/cosmos/cosmos-sdk/client/flags"

// 	"mods.irisnet.org/simapp"
// 	htlccli "mods.irisnet.org/htlc/client/cli"
// 	htlctypes "mods.irisnet.org/htlc/types"
// )

// // MsgRedelegateExec creates a redelegate message.
// func CreateHTLCExec(t *testing.T,
// 	network simapp.Network,
// 	clientCtx client.Context,
// 	from string,
// 	extraArgs ...string) *simapp.ResponseTx {
// 	args := []string{
// 		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
// 	}
// 	args = append(args, extraArgs...)
// 	return network.ExecTxCmdWithResult(t, clientCtx, htlccli.GetCmdCreateHTLC(), args)
// }

// func ClaimHTLCExec(t *testing.T,
// 	network simapp.Network,
// 	clientCtx client.Context,
// 	from string,
// 	id string,
// 	secret string,
// 	extraArgs ...string) *simapp.ResponseTx {
// 	args := []string{
// 		id,
// 		secret,
// 		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
// 	}
// 	args = append(args, extraArgs...)
// 	return network.ExecTxCmdWithResult(t, clientCtx, htlccli.GetCmdClaimHTLC(), args)
// }

// func QueryHTLCExec(t *testing.T,
// 	network simapp.Network,
// 	clientCtx client.Context,
// 	id string,
// 	extraArgs ...string) *htlctypes.HTLC {
// 	args := []string{
// 		id,
// 		fmt.Sprintf("--%s=json", cli.OutputFlag),
// 	}
// 	args = append(args, extraArgs...)
// 	response := &htlctypes.HTLC{}
// 	network.ExecQueryCmd(t, clientCtx, htlccli.GetCmdQueryHTLC(), args, response)
// 	return response
// }
