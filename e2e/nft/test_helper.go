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
	t.Helper()
	args := []string{
		denom,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, nftcli.GetCmdIssueDenom(), args)
}

// BurnNFTExec creates a nft burnt message.
func BurnNFTExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	denomID string,
	tokenID string,
	extraArgs ...string,
) *simapp.ResponseTx {
	t.Helper()
	args := []string{
		denomID,
		tokenID,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, nftcli.GetCmdBurnNFT(), args)
}

// MintNFTExec creates a nft minted message.
func MintNFTExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	denomID string,
	tokenID string,
	extraArgs ...string,
) *simapp.ResponseTx {
	t.Helper()
	args := []string{
		denomID,
		tokenID,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, nftcli.GetCmdMintNFT(), args)
}

// EditNFTExec creates a nft edited message.
func EditNFTExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	denomID string,
	tokenID string,
	extraArgs ...string,
) *simapp.ResponseTx {
	t.Helper()
	args := []string{
		denomID,
		tokenID,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, nftcli.GetCmdEditNFT(), args)
}

// TransferNFTExec creates a nft transferred message.
func TransferNFTExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	recipient string,
	denomID string,
	tokenID string,
	extraArgs ...string,
) *simapp.ResponseTx {
	t.Helper()
	args := []string{
		recipient,
		denomID,
		tokenID,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, nftcli.GetCmdTransferNFT(), args)
}

// TransferDenomExec creates a nft transferred message.
func TransferDenomExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	recipient string,
	denomID string,
	extraArgs ...string,
) *simapp.ResponseTx {
	t.Helper()
	args := []string{
		recipient,
		denomID,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}

	args = append(args, extraArgs...)
	return network.ExecTxCmdWithResult(t, clientCtx, nftcli.GetCmdTransferDenom(), args)
}

// QueryDenomExec query denom.
func QueryDenomExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	denomID string,
	extraArgs ...string,
) *nfttypes.Denom {
	t.Helper()
	args := []string{
		denomID,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	response := &nfttypes.Denom{}
	network.ExecQueryCmd(t, clientCtx, nftcli.GetCmdQueryDenom(), args, response)
	return response
}

// QueryCollectionExec query collection.
func QueryCollectionExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	denomID string,
	extraArgs ...string,
) *nfttypes.QueryCollectionResponse {
	t.Helper()
	args := []string{
		denomID,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	response := &nfttypes.QueryCollectionResponse{}
	network.ExecQueryCmd(t, clientCtx, nftcli.GetCmdQueryCollection(), args, response)
	return response
}

// QueryDenomsExec query denoms.
func QueryDenomsExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	extraArgs ...string,
) *nfttypes.QueryDenomsResponse {
	t.Helper()
	args := []string{
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	response := &nfttypes.QueryDenomsResponse{}
	network.ExecQueryCmd(t, clientCtx, nftcli.GetCmdQueryDenoms(), args, response)
	return response
}

// QuerySupplyExec query supply.
func QuerySupplyExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	denom string,
	extraArgs ...string,
) *nfttypes.QuerySupplyResponse {
	t.Helper()
	args := []string{
		denom,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	response := &nfttypes.QuerySupplyResponse{}
	network.ExecQueryCmd(t, clientCtx, nftcli.GetCmdQuerySupply(), args, response)
	return response
}

// QueryOwnerExec query owner.
func QueryOwnerExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	address string,
	extraArgs ...string,
) *nfttypes.QueryNFTsOfOwnerResponse {
	t.Helper()
	args := []string{
		address,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	response := &nfttypes.QueryNFTsOfOwnerResponse{}
	network.ExecQueryCmd(t, clientCtx, nftcli.GetCmdQueryOwner(), args, response)
	return response
}

// QueryNFTExec query nft.
func QueryNFTExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	denomID string,
	tokenID string,
	extraArgs ...string,
) *nfttypes.BaseNFT {
	t.Helper()
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
