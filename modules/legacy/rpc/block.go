package rpc

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/modules/legacy/types"
	"github.com/spf13/cobra"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

//BlockCommand returns the verified block data for a given heights
func BlockCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "block [height]",
		Short: "Get verified data for a the block at given height",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

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

			output, err := getBlock(clientCtx, height)
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
		},
	}

	cmd.Flags().StringP(flags.FlagNode, "n", "tcp://localhost:26657", "Node to connect to")

	return cmd
}

func getBlock(clientCtx client.Context, height *int64) ([]byte, error) {
	// get the node
	node, err := clientCtx.GetNode()
	if err != nil {
		return nil, err
	}

	// header -> BlockchainInfo
	// header, tx -> Block
	// results -> BlockResults
	res, err := node.Block(context.Background(), height)
	if err != nil {
		return nil, err
	}

	result := convertResultBlock(res)

	return legacy.Cdc.MarshalJSON(result)
}

func convertResultBlock(tmBlock *ctypes.ResultBlock) *types.ResultBlock{
	precommits := make([]*types.Vote,len(tmBlock.Block.LastCommit.Signatures))
	for k,v := range tmBlock.Block.LastCommit.Signatures{
		precommits[k] = &types.Vote{
			Type:             0x02,
			Height:           tmBlock.Block.LastCommit.Height,
			Round:            tmBlock.Block.LastCommit.Round,
			BlockID:          tmBlock.Block.LastCommit.BlockID,
			Timestamp:        v.Timestamp,
			ValidatorAddress: v.ValidatorAddress,
			ValidatorIndex:   k,
			Signature:        v.Signature,
		}
	}
	header := types.Header{
		Version:            tmBlock.Block.Version,
		ChainID:            tmBlock.Block.ChainID,
		Height:             tmBlock.Block.Height,
		Time:               tmBlock.Block.Time,
		NumTxs:             0,
		TotalTxs:           int64(len(tmBlock.Block.Data.Txs)),
		LastBlockID:        tmBlock.Block.LastBlockID,
		LastCommitHash:     tmBlock.Block.LastCommitHash,
		DataHash:           tmBlock.Block.DataHash,
		ValidatorsHash:     tmBlock.Block.ValidatorsHash,
		NextValidatorsHash: tmBlock.Block.NextValidatorsHash,
		ConsensusHash:      tmBlock.Block.ConsensusHash,
		AppHash:            tmBlock.Block.AppHash,
		LastResultsHash:    tmBlock.Block.LastResultsHash,
		EvidenceHash:       tmBlock.Block.EvidenceHash,
		ProposerAddress:    tmBlock.Block.ProposerAddress,
	}
	return &types.ResultBlock{
		BlockMeta: &types.BlockMeta{
			BlockID: tmBlock.BlockID,
			Header:  header,
		},
		Block:     &types.Block{
			Header:     header,
			Data:       tmBlock.Block.Data,
			Evidence:   types.EvidenceData{
				Evidence: tmBlock.Block.Evidence.Evidence,
			},
			LastCommit: &types.Commit{
				BlockID:    tmBlock.BlockID,
				Precommits: precommits,
			},
		},
	}
}

// get the current blockchain height
func GetChainHeight(clientCtx client.Context) (int64, error) {
	node, err := clientCtx.GetNode()
	if err != nil {
		return -1, err
	}

	status, err := node.Status(context.Background())
	if err != nil {
		return -1, err
	}

	height := status.SyncInfo.LatestBlockHeight
	return height, nil
}

// REST handler to get a block
func BlockRequestHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		height, err := strconv.ParseInt(vars["height"], 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest,
				"couldn't parse block height. Assumed format is '/block/{height}'.")
			return
		}

		chainHeight, err := GetChainHeight(clientCtx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, "failed to parse chain height")
			return
		}

		if height > chainHeight {
			rest.WriteErrorResponse(w, http.StatusNotFound, "requested block height is bigger then the chain length")
			return
		}

		output, err := getBlock(clientCtx, &height)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		rest.PostProcessResponseBare(w, clientCtx, output)
	}
}

// REST handler to get the latest block
func LatestBlockRequestHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		output, err := getBlock(clientCtx, nil)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		rest.PostProcessResponseBare(w, clientCtx, output)
	}
}
