package cli

import (
	"encoding/hex"
	"fmt"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/client/context"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func GetCmdQureyHash(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tx [hash]",
		Short: "Matches this txhash over all committed blocks",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// find the key to look up the account
			hashHexStr := viper.GetString(FlagTxHash)

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			output, err := queryTx(cdc, cliCtx, hashHexStr)
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
	cmd.Flags().String(FlagTxHash, "", "tx hash")

	return cmd
}

func queryTx(cdc *wire.Codec, cliCtx context.CLIContext, hashHexStr string) ([]byte, error) {
	hash, err := hex.DecodeString(hashHexStr)
	if err != nil {
		return nil, err
	}

	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	res, err := node.Tx(hash, !cliCtx.TrustNode)
	if err != nil {
		return nil, err
	}

	if !cliCtx.TrustNode {
		err := ValidateTxResult(cliCtx, res)
		if err != nil {
			return nil, err
		}
	}

	info, err := formatTxResult(cdc, res)
	if err != nil {
		return nil, err
	}

	return cdc.MarshalJSONIndent(info, "", "  ")
}

// ValidateTxResult performs transaction verification
func ValidateTxResult(cliCtx context.CLIContext, res *ctypes.ResultTx) error {
	check, err := cliCtx.Certify(res.Height)
	if err != nil {
		return err
	}

	err = res.Proof.Validate(check.Header.DataHash)
	if err != nil {
		return err
	}
	return nil
}
