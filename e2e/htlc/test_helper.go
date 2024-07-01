package htlc

import (
	"fmt"
	"testing"

	"github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	htlccli "mods.irisnet.org/modules/htlc/client/cli"
	htlctypes "mods.irisnet.org/modules/htlc/types"
	"mods.irisnet.org/simapp"
)

// CreateHTLCExec executes the creation of an HTLC with the provided parameters.
//
// Parameters:
// - t: testing.T instance for running test functions
// - network: simapp.Network instance for simulating the network
// - clientCtx: client.Context instance for client context
// - from: string representing the sender of the HTLC
// - extraArgs: variadic string arguments for additional parameters
//
// Returns a simapp.ResponseTx pointer.
func CreateHTLCExec(t *testing.T,
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
	return network.ExecTxCmdWithResult(t, clientCtx, htlccli.GetCmdCreateHTLC(), args)
}

// ClaimHTLCExec executes the claiming of an HTLC with the provided parameters.
//
// Parameters:
// - t: testing.T instance for running test functions
// - network: simapp.Network instance for simulating the network
// - clientCtx: client.Context instance for client context
// - from: string representing the sender of the HTLC
// - id: string representing the ID of the HTLC
// - secret: string representing the secret of the HTLC
// - extraArgs: variadic string arguments for additional parameters
//
// Returns a *simapp.ResponseTx pointer.
func ClaimHTLCExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	id string,
	secret string,
	extraArgs ...string,
) *simapp.ResponseTx {
	t.Helper()
	args := []string{
		id,
		secret,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)
	return network.ExecTxCmdWithResult(t, clientCtx, htlccli.GetCmdClaimHTLC(), args)
}

// QueryHTLCExec executes a query for an HTLC based on the provided ID and additional arguments.
//
// Parameters:
// - t: testing.T instance for running test functions
// - network: simapp.Network instance for simulating the network
// - clientCtx: client.Context instance for client context
// - id: string representing the ID of the HTLC
// - extraArgs: variadic string arguments for additional parameters
//
// Returns an htlctypes.HTLC pointer.
func QueryHTLCExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	id string,
	extraArgs ...string,
) *htlctypes.HTLC {
	t.Helper()
	args := []string{
		id,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)
	response := &htlctypes.HTLC{}
	network.ExecQueryCmd(t, clientCtx, htlccli.GetCmdQueryHTLC(), args, response)
	return response
}
