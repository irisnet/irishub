package testutil

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/cometbft/cometbft/libs/cli"

// 	"github.com/cosmos/cosmos-sdk/client"
// 	"github.com/cosmos/cosmos-sdk/client/flags"

// 	"mods.irisnet.org/simapp"
// 	oraclecli "mods.irisnet.org/modules/oracle/client/cli"
// 	oracletypes "mods.irisnet.org/modules/oracle/types"
// )

// // MsgRedelegateExec creates a redelegate message.
// func CreateFeedExec(t *testing.T,
// 	network simapp.Network,
// 	clientCtx client.Context,
// 	from string,
// 	extraArgs ...string) *simapp.ResponseTx {
// 	args := []string{
// 		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
// 	}
// 	args = append(args, extraArgs...)

// 	return network.ExecTxCmdWithResult(t, clientCtx, oraclecli.GetCmdCreateFeed(), args)
// }

// func EditFeedExec(t *testing.T,
// 	network simapp.Network,
// 	clientCtx client.Context,
// 	from string,
// 	feedName string,
// 	extraArgs ...string) *simapp.ResponseTx {
// 	args := []string{
// 		feedName,
// 		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
// 	}
// 	args = append(args, extraArgs...)

// 	return network.ExecTxCmdWithResult(t, clientCtx, oraclecli.GetCmdEditFeed(), args)
// }

// func StartFeedExec(t *testing.T,
// 	network simapp.Network,
// 	clientCtx client.Context,
// 	from string,
// 	feedName string,
// 	extraArgs ...string) *simapp.ResponseTx {
// 	args := []string{
// 		feedName,
// 		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
// 	}
// 	args = append(args, extraArgs...)

// 	return network.ExecTxCmdWithResult(t, clientCtx, oraclecli.GetCmdStartFeed(), args)
// }

// func PauseFeedExec(t *testing.T,
// 	network simapp.Network,
// 	clientCtx client.Context,
// 	from string,
// 	feedName string,
// 	extraArgs ...string) *simapp.ResponseTx {
// 	args := []string{
// 		feedName,
// 		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
// 	}
// 	args = append(args, extraArgs...)

// 	return network.ExecTxCmdWithResult(t, clientCtx, oraclecli.GetCmdPauseFeed(), args)
// }

// func QueryFeedExec(t *testing.T,
// 	network simapp.Network,
// 	clientCtx client.Context,
// 	feedName string,
// 	extraArgs ...string) *oracletypes.FeedContext {
// 	args := []string{
// 		feedName,
// 		fmt.Sprintf("--%s=json", cli.OutputFlag),
// 	}
// 	args = append(args, extraArgs...)

// 	response := &oracletypes.FeedContext{}
// 	network.ExecQueryCmd(t, clientCtx, oraclecli.GetCmdQueryFeed(), args, response)
// 	return response
// }

// func QueryFeedsExec(t *testing.T,
// 	network simapp.Network,
// 	clientCtx client.Context,
// 	extraArgs ...string) *oracletypes.QueryFeedsResponse {
// 	args := []string{
// 		fmt.Sprintf("--%s=json", cli.OutputFlag),
// 	}
// 	args = append(args, extraArgs...)

// 	response := &oracletypes.QueryFeedsResponse{}
// 	network.ExecQueryCmd(t, clientCtx, oraclecli.GetCmdQueryFeeds(), args, response)
// 	return response
// }

// func QueryFeedValueExec(t *testing.T,
// 	network simapp.Network,
// 	clientCtx client.Context,
// 	feedName string,
// 	extraArgs ...string) *oracletypes.QueryFeedValueResponse {
// 	args := []string{
// 		feedName,
// 		fmt.Sprintf("--%s=json", cli.OutputFlag),
// 	}
// 	args = append(args, extraArgs...)

// 	response := &oracletypes.QueryFeedValueResponse{}
// 	network.ExecQueryCmd(t, clientCtx, oraclecli.GetCmdQueryFeedValue(), args, response)
// 	return response
// }
