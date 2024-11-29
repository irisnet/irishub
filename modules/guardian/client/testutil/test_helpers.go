package testutil

import (
	"fmt"
	"testing"

	"github.com/cometbft/cometbft/libs/cli"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/network"

	guardiancli "github.com/irisnet/irishub/v4/modules/guardian/client/cli"
	apptestutil "github.com/irisnet/irishub/v4/testutil"
)

// CreateSuperExec creates a new super
func CreateSuperExec(
	t *testing.T,
	network *network.Network,
	clientCtx client.Context,
	from string,
	extraArgs ...string,
) *coretypes.ResultTx {
	args := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return apptestutil.ExecCommand(t, network, clientCtx, guardiancli.GetCmdCreateSuper(), args)
}

// DeleteSuperExec deletes a super
func DeleteSuperExec(
	t *testing.T,
	network *network.Network,
	clientCtx client.Context,
	from string,
	extraArgs ...string,
) *coretypes.ResultTx {
	args := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return apptestutil.ExecCommand(t, network, clientCtx, guardiancli.GetCmdDeleteSuper(), args)
}

// QuerySupersExec queries supers
func QuerySupersExec(clientCtx client.Context, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, guardiancli.GetCmdQuerySupers(), args)
}
