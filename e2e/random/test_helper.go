package random

import (
	"fmt"
	"testing"

	"github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	randomcli "mods.irisnet.org/modules/random/client/cli"
	randomtypes "mods.irisnet.org/modules/random/types"
	"mods.irisnet.org/simapp"
)

// RequestRandomExec creates a random request execution message.
//
// Parameters:
// - t: The testing context.
// - network: The simulation network.
// - clientCtx: The client context.
// - from: The sender address.
// - extraArgs: Additional arguments.
// Returns a pointer to a simapp.ResponseTx.
func RequestRandomExec(t *testing.T,
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

	return network.ExecTxCmdWithResult(t, clientCtx, randomcli.GetCmdRequestRandom(), args)
}

// QueryRandomExec queries the random number for the given request ID using the provided network, client context, and optional extra arguments.
//
// Parameters:
// - t: The testing context.
// - network: The simulation network.
// - clientCtx: The client context.
// - requestID: The ID of the random request.
// - extraArgs: Optional extra arguments.
//
// Returns:
// - *randomtypes.Random: The response containing the queried random number.
func QueryRandomExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	requestID string,
	extraArgs ...string,
) *randomtypes.Random {
	t.Helper()
	args := []string{
		requestID,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	response := &randomtypes.Random{}
	network.ExecQueryCmd(t, clientCtx, randomcli.GetCmdQueryRandom(), args, response)
	return response
}

// QueryRandomRequestQueueExec queries the random request queue for a given height using the provided network, client context, and optional extra arguments.
//
// Parameters:
// - t: The testing context.
// - network: The simulation network.
// - clientCtx: The client context.
// - genHeight: The height at which to query the random request queue.
// - extraArgs: Optional extra arguments.
//
// Returns:
// - *randomtypes.QueryRandomRequestQueueResponse: The response containing the queried random request queue.
func QueryRandomRequestQueueExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	genHeight string,
	extraArgs ...string,
) *randomtypes.QueryRandomRequestQueueResponse {
	t.Helper()
	args := []string{
		genHeight,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	response := &randomtypes.QueryRandomRequestQueueResponse{}
	network.ExecQueryCmd(t, clientCtx, randomcli.GetCmdQueryRandomRequestQueue(), args, response)
	return response
}
