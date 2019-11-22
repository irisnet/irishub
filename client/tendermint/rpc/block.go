package rpc

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/spf13/cobra"
	tmliteProxy "github.com/tendermint/tendermint/lite/proxy"
	"net/http"
	"strconv"
)

//BlockCommand returns the verified block data for a given heights
func BlockCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "block [height]",
		Short:   "Get verified data for the block at given height",
		Example: "iriscli tendermint block",
		Args:    cobra.MaximumNArgs(1),
		RunE:    printBlock,
	}
	cmd.Flags().Bool(client.FlagIndentResponse, true, "Add indent to JSON response")
	cmd.Flags().StringP(client.FlagNode, "n", "tcp://localhost:26657", "Node to connect to")
	cmd.Flags().Bool(client.FlagTrustNode, false, "Trust connected full node (don't verify proofs for responses)")
	cmd.Flags().String(client.FlagChainID, "", "Chain ID of Tendermint node")
	return cmd
}

func getBlock(cliCtx context.CLIContext, height *int64) ([]byte, error) {
	// get the node
	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	// header -> BlockchainInfo
	// header, tx -> Block
	// results -> BlockResults
	res, err := node.Block(height)
	if err != nil || res.Block == nil {
		return nil, fmt.Errorf("block %d not found", *height)
	}

	if !cliCtx.TrustNode {
		check, err := cliCtx.Verify(res.Block.Height)
		if err != nil {
			return nil, err
		}

		err = tmliteProxy.ValidateBlockMeta(res.BlockMeta, check)
		if err != nil {
			return nil, err
		}

		err = tmliteProxy.ValidateBlock(res.Block, check)
		if err != nil {
			return nil, err
		}
	}

	if cliCtx.Indent {
		return cdc.MarshalJSONIndent(res, "", "  ")
	}
	return cdc.MarshalJSON(res)
}

// get the current blockchain height
func GetChainHeight(cliCtx context.CLIContext) (int64, error) {
	node, err := cliCtx.GetNode()
	if err != nil {
		return -1, err
	}
	status, err := node.Status()
	if err != nil {
		return -1, err
	}
	height := status.SyncInfo.LatestBlockHeight
	return height, nil
}

// CMD

func printBlock(cmd *cobra.Command, args []string) error {
	var height *int64
	// optional height
	if len(args) > 0 {
		h, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		if h > 0 {
			tmp := int64(h)
			height = &tmp
		}
	}

	output, err := getBlock(context.NewCLIContext(), height)
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

// REST handler to get a block
func BlockRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		height, err := strconv.ParseInt(vars["height"], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("ERROR: Couldn't parse block height. Assumed format is '/block/{height}'."))
			return
		}
		chainHeight, err := GetChainHeight(cliCtx)
		if height > chainHeight {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("ERROR: Requested block height is bigger then the chain length."))
			return
		}
		output, err := getBlock(cliCtx, &height)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		utils.PostProcessResponse(w, cdc, output, cliCtx.Indent)
	}
}

// REST handler to get the latest block
func LatestBlockRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		height, err := GetChainHeight(cliCtx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		output, err := getBlock(cliCtx, &height)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		utils.PostProcessResponse(w, cdc, output, cliCtx.Indent)
	}
}
