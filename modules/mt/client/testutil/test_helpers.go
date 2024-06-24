package testutil

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/cometbft/cometbft/libs/cli"
// 	"github.com/cosmos/gogoproto/proto"

// 	"github.com/cosmos/cosmos-sdk/client"
// 	"github.com/cosmos/cosmos-sdk/client/flags"

// 	"mods.irisnet.org/simapp"
// 	mtcli "mods.irisnet.org/modules/mt/client/cli"
// 	mttypes "mods.irisnet.org/modules/mt/types"
// )

// // IssueDenomExec creates a redelegate message.
// func IssueDenomExec(
// 	t *testing.T,
// 	network simapp.Network,
// 	clientCtx client.Context,
// 	from string,
// 	extraArgs ...string,
// ) *simapp.ResponseTx {
// 	args := []string{
// 		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
// 	}
// 	args = append(args, extraArgs...)

// 	return network.ExecTxCmdWithResult(t, clientCtx, mtcli.GetCmdIssueDenom(), args)
// }

// func BurnMTExec(t *testing.T,
// 	network simapp.Network,
// 	clientCtx client.Context,
// 	from string,
// 	denomID string,
// 	mtID string,
// 	amount string,
// 	extraArgs ...string,
// ) *simapp.ResponseTx {
// 	args := []string{
// 		denomID,
// 		mtID,
// 		amount,
// 		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
// 	}
// 	args = append(args, extraArgs...)

// 	return network.ExecTxCmdWithResult(t, clientCtx, mtcli.GetCmdBurnMT(), args)
// }

// func MintMTExec(t *testing.T,
// 	network simapp.Network,
// 	clientCtx client.Context,
// 	from string,
// 	denomID string,
// 	extraArgs ...string,
// ) *simapp.ResponseTx {
// 	args := []string{
// 		denomID,
// 		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
// 	}
// 	args = append(args, extraArgs...)

// 	return network.ExecTxCmdWithResult(t, clientCtx, mtcli.GetCmdMintMT(), args)
// }

// func EditMTExec(t *testing.T,
// 	network simapp.Network,
// 	clientCtx client.Context,
// 	from string,
// 	denomID string,
// 	mtID string,
// 	extraArgs ...string,
// ) *simapp.ResponseTx {
// 	args := []string{
// 		denomID,
// 		mtID,
// 		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
// 	}
// 	args = append(args, extraArgs...)

// 	return network.ExecTxCmdWithResult(t, clientCtx, mtcli.GetCmdEditMT(), args)
// }

// func TransferMTExec(t *testing.T,
// 	network simapp.Network,
// 	clientCtx client.Context,
// 	from string,
// 	recipient string,
// 	denomID string,
// 	mtID string,
// 	amount string,
// 	extraArgs ...string,
// ) *simapp.ResponseTx {
// 	args := []string{
// 		from,
// 		recipient,
// 		denomID,
// 		mtID,
// 		amount,
// 	}
// 	args = append(args, extraArgs...)

// 	return network.ExecTxCmdWithResult(t, clientCtx, mtcli.GetCmdTransferMT(), args)
// }

// func QueryDenomExec(t *testing.T,
// 	network simapp.Network,
// 	clientCtx client.Context,
// 	denomID string,
// 	extraArgs ...string) *mttypes.Denom {
// 	args := []string{
// 		denomID,
// 		fmt.Sprintf("--%s=json", cli.OutputFlag),
// 	}
// 	args = append(args, extraArgs...)

// 	response := &mttypes.Denom{}
// 	network.ExecQueryCmd(t, clientCtx, mtcli.GetCmdQueryDenom(), args, response)
// 	return response
// }

// func QueryDenomsExec(t *testing.T,
// 	network simapp.Network,
// 	clientCtx client.Context,
// 	extraArgs ...string) *mttypes.QueryDenomsResponse {
// 	args := []string{
// 		fmt.Sprintf("--%s=json", cli.OutputFlag),
// 	}
// 	args = append(args, extraArgs...)

// 	response := &mttypes.QueryDenomsResponse{}
// 	network.ExecQueryCmd(t, clientCtx, mtcli.GetCmdQueryDenoms(), args, response)
// 	return response
// }

// func QueryMTsExec(t *testing.T,
// 	network simapp.Network,
// 	clientCtx client.Context,
// 	denomID string,
// 	resp proto.Message,
// 	extraArgs ...string,
// ) {
// 	args := []string{
// 		denomID,
// 		fmt.Sprintf("--%s=json", cli.OutputFlag),
// 	}
// 	args = append(args, extraArgs...)

// 	network.ExecQueryCmd(t, clientCtx, mtcli.GetCmdQueryMTs(), args, resp)
// }

// func QueryMTExec(t *testing.T,
// 	network simapp.Network,
// 	clientCtx client.Context,
// 	denomID string,
// 	mtID string,
// 	extraArgs ...string) *mttypes.MT {
// 	args := []string{
// 		denomID,
// 		mtID,
// 		fmt.Sprintf("--%s=json", cli.OutputFlag),
// 	}
// 	args = append(args, extraArgs...)
// 	response := &mttypes.MT{}
// 	network.ExecQueryCmd(t, clientCtx, mtcli.GetCmdQueryMT(), args, response)
// 	return response
// }

// func QueryBlancesExec(t *testing.T,
// 	network simapp.Network,
// 	clientCtx client.Context,
// 	from string,
// 	denomID string,
// 	extraArgs ...string) *mttypes.QueryBalancesResponse {
// 	args := []string{
// 		from,
// 		denomID,
// 		fmt.Sprintf("--%s=json", cli.OutputFlag),
// 	}
// 	args = append(args, extraArgs...)

// 	response := &mttypes.QueryBalancesResponse{}
// 	network.ExecQueryCmd(t, clientCtx, mtcli.GetCmdQueryBalances(), args, response)
// 	return response
// }

// func TransferDenomExec(t *testing.T,
// 	network simapp.Network,
// 	clientCtx client.Context,
// 	from string,
// 	recipient string,
// 	denomID string,
// 	extraArgs ...string,
// ) *simapp.ResponseTx {
// 	args := []string{
// 		from,
// 		recipient,
// 		denomID,
// 	}

// 	args = append(args, extraArgs...)
// 	return network.ExecTxCmdWithResult(t, clientCtx, mtcli.GetCmdTransferDenom(), args)
// }
