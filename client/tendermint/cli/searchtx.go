package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/client/context"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

const (
	flagTags = "tag"
	flagAny  = "any"
)

// default client command to search through tagged transactions
func SearchTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "txs",
		Short: "Search for all transactions that match the given tags.",
		Long: strings.TrimSpace(`
Search for transactions that match the given tags. By default, transactions must match ALL tags 
passed to the --tags option. To match any transaction, use the --any option.

For example:

$ gaiacli tendermint txs --tag test1,test2

will match any transaction tagged with both test1,test2. To match a transaction tagged with either
test1 or test2, use:

$ gaiacli tendermint txs --tag test1,test2 --any
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			tags := viper.GetStringSlice(flagTags)

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txs, err := searchTxs(cliCtx, cdc, tags)
			if err != nil {
				return err
			}

			output, err := cdc.MarshalJSONIndent(txs,"", "  ")
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
		},
	}

	cmd.Flags().StringP(client.FlagNode, "n", "tcp://localhost:26657", "Node to connect to")
	cmd.Flags().Bool(client.FlagTrustNode, false, "Trust connected full node (don't verify proofs for responses)")
	cmd.Flags().String(client.FlagChainID, "", "Chain ID of Tendermint node")
	cmd.Flags().StringSlice(flagTags, nil, "Comma-separated list of tags that must match")
	cmd.Flags().Bool(flagAny, false, "Return transactions that match ANY tag, rather than ALL")
	return cmd
}

func searchTxs(cliCtx context.CLIContext, cdc *wire.Codec, tags []string) ([]Info, error) {
	if len(tags) == 0 {
		return nil, errors.New("must declare at least one tag to search")
	}

	// XXX: implement ANY
	query := strings.Join(tags, " AND ")

	// get the node
	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	prove := !cliCtx.TrustNode

	// TODO: take these as args
	page := 0
	perPage := 100
	res, err := node.TxSearch(query, prove, page, perPage)
	if err != nil {
		return nil, err
	}

	if prove {
		for _, tx := range res.Txs {
			err := ValidateTxResult(cliCtx, tx)
			if err != nil {
				return nil, err
			}
		}
	}

	info, err := FormatTxResults(cdc, res.Txs)
	if err != nil {
		return nil, err
	}

	return info, nil
}

// parse the indexed txs into an array of Info
func FormatTxResults(cdc *wire.Codec, res []*ctypes.ResultTx) ([]Info, error) {
	var err error
	out := make([]Info, len(res))
	for i := range res {
		out[i], err = formatTxResult(cdc, res[i])
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}
