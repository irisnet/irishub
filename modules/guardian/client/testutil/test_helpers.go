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

	guardiancli "github.com/irisnet/irishub/v2/modules/guardian/client/cli"
	"github.com/irisnet/irishub/v2/simapp"
)

// MsgRedelegateExec creates a redelegate message.
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

	return simapp.ExecTxCmdWithResult(t, network, clientCtx, guardiancli.GetCmdCreateSuper(), args)
}

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

	return simapp.ExecTxCmdWithResult(t, network, clientCtx, guardiancli.GetCmdDeleteSuper(), args)
}

func QuerySupersExec(clientCtx client.Context, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, guardiancli.GetCmdQuerySupers(), args)
}
