package oracle

import (
	"fmt"
	"testing"

	"github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	oraclecli "mods.irisnet.org/modules/oracle/client/cli"
	oracletypes "mods.irisnet.org/modules/oracle/types"
	"mods.irisnet.org/simapp"
)

// CreateFeedExec creates a feed execution message.
//
// Parameters:
// - t: The testing context.
// - network: The simulation network.
// - clientCtx: The client context.
// - from: The sender address.
// - extraArgs: Additional arguments.
// Returns a response transaction.
func CreateFeedExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	extraArgs ...string,
) *simapp.ResponseTx {
	t.Helper()
	args := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, oraclecli.GetCmdCreateFeed(), args)
}

// EditFeedExec creates a feed execution message.
//
// Parameters:
// - t: The testing context.
// - network: The simulation network.
// - clientCtx: The client context.
// - from: The sender address.
// - feedName: The name of the feed.
// - extraArgs: Additional arguments.
//
// Returns:
// - A pointer to a simapp.ResponseTx.
func EditFeedExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	feedName string,
	extraArgs ...string,
) *simapp.ResponseTx {
	t.Helper()
	args := []string{
		feedName,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, oraclecli.GetCmdEditFeed(), args)
}

// StartFeedExec starts a feed execution message.
//
// Parameters:
// - t: The testing context.
// - network: The simulation network.
// - clientCtx: The client context.
// - from: The sender address.
// - feedName: The name of the feed.
// - extraArgs: Additional arguments.
//
// Returns:
// - A pointer to a simapp.ResponseTx.
func StartFeedExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	feedName string,
	extraArgs ...string,
) *simapp.ResponseTx {
	t.Helper()
	args := []string{
		feedName,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, oraclecli.GetCmdStartFeed(), args)
}

// PauseFeedExec creates a transaction to pause a feed.
//
// t: The testing context.
// network: The simulation network.
// clientCtx: The client context.
// from: The sender address.
// feedName: The name of the feed.
// extraArgs: Additional arguments.
// Returns a pointer to a simapp.ResponseTx.
func PauseFeedExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	feedName string,
	extraArgs ...string,
) *simapp.ResponseTx {
	t.Helper()
	args := []string{
		feedName,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, oraclecli.GetCmdPauseFeed(), args)
}

// QueryFeedExec creates a transaction to query a feed.
//
// Parameters:
// - t: The testing context.
// - network: The simulation network.
// - clientCtx: The client context.
// - feedName: The name of the feed.
// - extraArgs: Additional arguments.
// Returns a pointer to an oracletypes.FeedContext.
func QueryFeedExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	feedName string,
	extraArgs ...string,
) *oracletypes.FeedContext {
	t.Helper()
	args := []string{
		feedName,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	response := &oracletypes.FeedContext{}
	network.ExecQueryCmd(t, clientCtx, oraclecli.GetCmdQueryFeed(), args, response)
	return response
}

// QueryFeedsExec queries the feeds using the provided network, client context, and optional extra arguments.
//
// Parameters:
// - t: The testing context.
// - network: The simulation network.
// - clientCtx: The client context.
// - extraArgs: Optional extra arguments.
//
// Returns:
// - *oracletypes.QueryFeedsResponse: The response containing the queried feeds.
func QueryFeedsExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	extraArgs ...string,
) *oracletypes.QueryFeedsResponse {
	t.Helper()
	args := []string{
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	response := &oracletypes.QueryFeedsResponse{}
	network.ExecQueryCmd(t, clientCtx, oraclecli.GetCmdQueryFeeds(), args, response)
	return response
}

// QueryFeedValueExec creates a transaction to query a feed value.
//
// Parameters:
// - t: The testing context.
// - network: The simulation network.
// - clientCtx: The client context.
// - feedName: The name of the feed.
// - extraArgs: Additional arguments.
// Returns a pointer to an oracletypes.QueryFeedValueResponse.
func QueryFeedValueExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	feedName string,
	extraArgs ...string,
) *oracletypes.QueryFeedValueResponse {
	t.Helper()
	args := []string{
		feedName,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	response := &oracletypes.QueryFeedValueResponse{}
	network.ExecQueryCmd(t, clientCtx, oraclecli.GetCmdQueryFeedValue(), args, response)
	return response
}
