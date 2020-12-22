package rpc

import (
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"net/http"
	"strconv"
	"context"

	"github.com/gorilla/mux"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	abci "github.com/tendermint/tendermint/abci/types"
)

type ResponseDeliverTx struct {
	Code      uint32                   `json:"code"`
	Data      []byte                   `json:"data"`
	Log       string                   `json:"log"`
	Info      string                   `json:"info"`
	GasWanted int64                    `json:"gas_wanted"`
	GasUsed   int64                    `json:"gas_used"`
	Tags      []ReadableTag              `json:"tags"`
}

type ReadableTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ResponseEndBlock struct {
	ValidatorUpdates      []abci.ValidatorUpdate   `json:"validator_updates"`
	ConsensusParamUpdates *abci.ConsensusParams    `json:"consensus_param_updates"`
	Tags                  []ReadableTag `json:"tags"`
}

type ResponseBeginBlock struct {
	Tags []ReadableTag `json:"tags"`
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

func ConvertEventToTags(event abci.Event) []ReadableTag {
	readableTags := make([]ReadableTag, len(event.Attributes))
	for i, kv := range event.Attributes {
		readableTags[i] = ReadableTag{
			Key:   event.Type + "." + string(kv.Key),
			Value: string(kv.Value),
		}
	}
	return readableTags
}

func MakeTagsHumanReadable(tags []abci.Event) []ReadableTag {
	var readableTags []ReadableTag
	for _, kv := range tags {
		readableTags = append(readableTags, ConvertEventToTags(kv)...)
	}
	return readableTags
}

func getBlockResult(clientCtx client.Context, height *int64) ([]byte, error) {
	// get the node
	node, err := clientCtx.GetNode()
	if err != nil {
		return nil, err
	}

	// header -> BlockchainInfo
	// header, tx -> Block
	// results -> BlockResults
	res, err := node.BlockResults(context.Background(), height)
	if err != nil {
		return nil, err
	}

	var delieverTxResponse []*ResponseDeliverTx
	for _, delieverTx := range res.TxsResults {
		delieverTxResponse = append(delieverTxResponse, &ResponseDeliverTx{
			Code:      delieverTx.Code,
			Data:      delieverTx.Data,
			Log:       delieverTx.Log,
			Info:      delieverTx.Info,
			GasWanted: delieverTx.GasWanted,
			GasUsed:   delieverTx.GasUsed,
			Tags:      MakeTagsHumanReadable(delieverTx.Events),
		})
	}
	abciResponses := ABCIResponses{
		DeliverTx: delieverTxResponse,
		BeginBlock: &ResponseBeginBlock{
			Tags: MakeTagsHumanReadable(res.BeginBlockEvents),
		},
		EndBlock: &ResponseEndBlock{
			ValidatorUpdates:      res.ValidatorUpdates,
			ConsensusParamUpdates: res.ConsensusParamUpdates,
			Tags:                  MakeTagsHumanReadable(res.EndBlockEvents),
		},
	}
	var response ResultBlockResults
	response.Height = res.Height
	response.Results = abciResponses
	return legacy.Cdc.MarshalJSON(response)
}

// REST handler to get a block
func BlockResultRequestHandlerFn(clientCtx client.Context) http.HandlerFunc {
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
		output, err := getBlockResult(clientCtx, &height)
		if rest.CheckInternalServerError(w, err) {
			return
		}
		rest.PostProcessResponseBare(w, clientCtx, output)
	}
}

// REST handler to get the latest block
func LatestBlockResultRequestHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		height, err := GetChainHeight(clientCtx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, "failed to parse chain height")
			return
		}
		output, err := getBlockResult(clientCtx, &height)
		if rest.CheckInternalServerError(w, err) {
			return
		}
		rest.PostProcessResponseBare(w, clientCtx, output)
	}
}
