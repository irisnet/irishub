package rpc

import (
	"fmt"
	"github.com/irisnet/irishub/client/tendermint"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
)

//BlockCommand returns the verified block data for a given heights
func BlockResultCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "block-result [height]",
		Short:   "Get block result at given height",
		Example: "iriscli tendermint block-result",
		Args:    cobra.MaximumNArgs(1),
		RunE:    printBlockResult,
	}
	cmd.Flags().Bool(client.FlagIndentResponse, true, "Add indent to JSON response")
	cmd.Flags().StringP(client.FlagNode, "n", "tcp://localhost:26657", "Node to connect to")
	cmd.Flags().Bool(client.FlagTrustNode, false, "Trust connected full node (don't verify proofs for responses)")
	cmd.Flags().String(client.FlagChainID, "", "Chain ID of Tendermint node")
	return cmd
}



type ResponseDeliverTx struct {
	Code      uint32        `json:"code"`
	Data      []byte        `json:"data"`
	Log       string        `json:"log"`
	Info      string        `json:"info"`
	GasWanted int64         `json:"gas_wanted"`
	GasUsed   int64         `json:"gas_used"`
	Tags      []tendermint.ReadableTag `json:"tags"`
}

type ResponseEndBlock struct {
	ValidatorUpdates      []abci.ValidatorUpdate `json:"validator_updates"`
	ConsensusParamUpdates *abci.ConsensusParams  `json:"consensus_param_updates"`
	Tags                  []tendermint.ReadableTag          `json:"tags"`
}

type ResponseBeginBlock struct {
	Tags []tendermint.ReadableTag `json:"tags"`
}

type ABCIResponses struct {
	DeliverTx  []*ResponseDeliverTx `json:"deliver_tx"`
	EndBlock   *ResponseEndBlock    `json:"end_block"`
	BeginBlock *ResponseBeginBlock  `json:"begin_block"`
}

type ResultBlockResults struct {
	Height  int64         `json:"height"`
	Results ABCIResponses `json:"results"`
}

func getBlockResult(cliCtx context.CLIContext, height *int64) ([]byte, error) {
	// get the node
	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	// header -> BlockchainInfo
	// header, tx -> Block
	// results -> BlockResults
	res, err := node.BlockResults(height)
	if err != nil {
		return nil, err
	}

	var delieverTxResponse []*ResponseDeliverTx
	for _, delieverTx := range res.Results.DeliverTx {
		delieverTxResponse = append(delieverTxResponse, &ResponseDeliverTx{
			Code:      delieverTx.Code,
			Data:      delieverTx.Data,
			Log:       delieverTx.Log,
			Info:      delieverTx.Info,
			GasWanted: delieverTx.GasWanted,
			GasUsed:   delieverTx.GasUsed,
			Tags:      tendermint.MakeTagsHumanReadable(delieverTx.Tags),
		})
	}
	abciResponses := ABCIResponses{
		DeliverTx: delieverTxResponse,
		BeginBlock: &ResponseBeginBlock{
			Tags: tendermint.MakeTagsHumanReadable(res.Results.BeginBlock.Tags),
		},
		EndBlock: &ResponseEndBlock{
			ValidatorUpdates:      res.Results.EndBlock.ValidatorUpdates,
			ConsensusParamUpdates: res.Results.EndBlock.ConsensusParamUpdates,
			Tags:                  tendermint.MakeTagsHumanReadable(res.Results.EndBlock.Tags),
		},
	}
	var response ResultBlockResults
	response.Height = res.Height
	response.Results = abciResponses

	if cliCtx.Indent {
		return cdc.MarshalJSONIndent(response, "", "  ")
	}
	return cdc.MarshalJSON(response)
}

func GetTxCoinFlow(cliCtx context.CLIContext, height *int64, hashStr string) ([]string, error) {
	// get the node
	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	res, err := node.BlockResults(height)
	if err != nil {
		return nil, err
	}

	var coinFlowTags []string
	endBlockTags := tendermint.MakeTagsHumanReadable(res.Results.EndBlock.Tags)
	found := false
	for _,tag := range endBlockTags {
		if tag.Key == hashStr {
			coinFlowTags = append(coinFlowTags, tag.Value)
			found = true
		} else if found {
			//txHash coin flow records are centralized distributed
			break
		}
	}
	return coinFlowTags, nil
}

// CMD

func printBlockResult(cmd *cobra.Command, args []string) error {
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

	output, err := getBlockResult(context.NewCLIContext(), height)
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

// REST handler to get a block
func BlockResultRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
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
		output, err := getBlockResult(cliCtx, &height)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		utils.PostProcessResponse(w, cdc, output, cliCtx.Indent)
	}
}

// REST handler to get the latest block
func LatestBlockResultRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		height, err := GetChainHeight(cliCtx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		output, err := getBlockResult(cliCtx, &height)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		utils.PostProcessResponse(w, cdc, output, cliCtx.Indent)
	}
}
