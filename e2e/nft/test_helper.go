package nft

import (
	"fmt"
	"testing"

	"github.com/cometbft/cometbft/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	nftcli "mods.irisnet.org/modules/nft/client/cli"
	nfttypes "mods.irisnet.org/modules/nft/types"
	"mods.irisnet.org/simapp"
)

// IssueDenomExec creates a redelegate message.
func IssueDenomExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	denom string,
	extraArgs ...string,
) *simapp.ResponseTx {
	args := []string{
		denom,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, nftcli.GetCmdIssueDenom(), args)
}

func BurnNFTExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	denomID string,
	tokenID string,
	extraArgs ...string,
) *simapp.ResponseTx {
	args := []string{
		denomID,
		tokenID,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, nftcli.GetCmdBurnNFT(), args)
}

func MintNFTExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	denomID string,
	tokenID string,
	extraArgs ...string,
) *simapp.ResponseTx {
	args := []string{
		denomID,
		tokenID,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, nftcli.GetCmdMintNFT(), args)
}

func EditNFTExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	denomID string,
	tokenID string,
	extraArgs ...string,
) *simapp.ResponseTx {
	args := []string{
		denomID,
		tokenID,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, nftcli.GetCmdEditNFT(), args)
}

func TransferNFTExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	recipient string,
	denomID string,
	tokenID string,
	extraArgs ...string,
) *simapp.ResponseTx {
	args := []string{
		recipient,
		denomID,
		tokenID,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, nftcli.GetCmdTransferNFT(), args)
}

func TransferDenomExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	recipient string,
	denomID string,
	extraArgs ...string,
) *simapp.ResponseTx {
	args := []string{
		recipient,
		denomID,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}

	args = append(args, extraArgs...)
	return network.ExecTxCmdWithResult(t, clientCtx, nftcli.GetCmdTransferDenom(), args)
}

func QueryDenomExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	denomID string,
	extraArgs ...string) *nfttypes.Denom {
	args := []string{
		denomID,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	response := &nfttypes.Denom{}
	network.ExecQueryCmd(t, clientCtx, nftcli.GetCmdQueryDenom(), args, response)
	return response
}

func QueryCollectionExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	denomID string,
	extraArgs ...string) *nfttypes.QueryCollectionResponse {
	args := []string{
		denomID,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	response := &nfttypes.QueryCollectionResponse{}
	network.ExecQueryCmd(t, clientCtx, nftcli.GetCmdQueryCollection(), args, response)
	return response
}

func QueryDenomsExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	extraArgs ...string) *nfttypes.QueryDenomsResponse {
	args := []string{
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	response := &nfttypes.QueryDenomsResponse{}
	network.ExecQueryCmd(t, clientCtx, nftcli.GetCmdQueryDenoms(), args, response)
	return response
}

func QuerySupplyExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	denom string,
	extraArgs ...string) *nfttypes.QuerySupplyResponse {
	args := []string{
		denom,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	response := &nfttypes.QuerySupplyResponse{}
	network.ExecQueryCmd(t, clientCtx, nftcli.GetCmdQuerySupply(), args, response)
	return response
}

func QueryOwnerExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	address string,
	extraArgs ...string) *nfttypes.QueryNFTsOfOwnerResponse {
	args := []string{
		address,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	response := &nfttypes.QueryNFTsOfOwnerResponse{}
	network.ExecQueryCmd(t, clientCtx, nftcli.GetCmdQueryOwner(), args, response)
	return response
}

func QueryNFTExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	denomID string,
	tokenID string,
	extraArgs ...string) *nfttypes.BaseNFT {
	args := []string{
		denomID,
		tokenID,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	response := &nfttypes.BaseNFT{}
	network.ExecQueryCmd(t, clientCtx, nftcli.GetCmdQueryNFT(), args, response)
	return response
}