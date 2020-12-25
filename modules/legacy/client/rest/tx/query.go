package tx

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/irisnet/irishub/modules/legacy/types"

	"github.com/cosmos/cosmos-sdk/client"
	clienttx "github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/rest"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	genutilrest "github.com/cosmos/cosmos-sdk/x/genutil/client/rest"
)

// Info is used to prepare info to display
type InfoCoinFlow struct {
	Hash      string            `json:"hash"`
	Height    int64             `json:"height"`
	Tx        sdk.Tx            `json:"tx"`
	Result    ResponseDeliverTx `json:"result"`
	Timestamp string            `json:"timestamp,omitempty"`
	CoinFlow  []string          `json:"coin_flow"`
}

type ResponseDeliverTx struct {
	Code                 uint32
	Data                 string
	Log                  string
	Info                 string
	GasWanted            int64
	GasUsed              int64
	Tags                 []ReadableTag
	Codespace            string
	XXX_NoUnkeyedLiteral struct{}
	XXX_unrecognized     []byte
	XXX_sizecache        int32
}

type ReadableTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// SearchTxsResult defines a structure for querying txs pageable
type SearchTxsResult struct {
	TotalCount uint64         `json:"total_count"` // Count of all txs
	Count      uint64         `json:"count"`       // Count of txs in current page
	PageNumber uint64         `json:"page_number"` // Index of current page, start from 1
	PageTotal  uint64         `json:"page_total"`  // Count of total pages
	Size       uint64         `json:"size"`        // Max count txs per page
	Txs        []InfoCoinFlow `json:"txs"`         // List of txs in current page
}

// QueryTxsRequestHandlerFn implements a REST handler that searches for transactions.
// Genesis transactions are returned if the height parameter is set to zero,
// otherwise the transactions are searched for by events.
func QueryTxsRequestHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			rest.WriteErrorResponse(
				w, http.StatusBadRequest,
				fmt.Sprintf("failed to parse query parameters: %s", err),
			)
			return
		}

		// if the height query param is set to zero, query for genesis transactions
		heightStr := r.FormValue("height")
		if heightStr != "" {
			if height, err := strconv.ParseInt(heightStr, 10, 64); err == nil && height == 0 {
				genutilrest.QueryGenesisTxs(clientCtx, w)
				return
			}
		}

		var (
			events      []string
			txs         []sdk.TxResponse
			page, limit int
		)

		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		if len(r.Form) == 0 {
			rest.PostProcessResponseBare(w, clientCtx, txs)
			return
		}

		events, page, limit, err = rest.ParseHTTPArgs(r)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		searchResult, err := authclient.QueryTxsByEvents(clientCtx, events, page, limit, "")
		if rest.CheckInternalServerError(w, err) {
			return
		}

		txsResult := make([]InfoCoinFlow, len(searchResult.Txs))
		for k, txRes := range searchResult.Txs {
			txResult, err := packStdTxResponse(w, clientCtx, txRes)
			if err != nil {
				return
			}
			txsResult[k] = *txResult
		}

		result := &SearchTxsResult{
			TotalCount: searchResult.PageTotal,
			Count:      searchResult.Count,
			PageNumber: searchResult.PageNumber,
			PageTotal:  searchResult.PageTotal,
			Size:       searchResult.Limit,
			Txs:        txsResult,
		}
		rest.PostProcessResponseBare(w, clientCtx, result)
	}
}

// QueryTxRequestHandlerFn implements a REST handler that queries a transaction
// by hash in a committed block.
func QueryTxRequestHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		hashHexStr := vars["hash"]

		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		output, err := authclient.QueryTx(clientCtx, hashHexStr)
		if err != nil {
			if strings.Contains(err.Error(), hashHexStr) {
				rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
				return
			}
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if output.Empty() {
			rest.WriteErrorResponse(w, http.StatusNotFound, fmt.Sprintf("no transaction found with hash %s", hashHexStr))
		}

		result, err := packStdTxResponse(w, clientCtx, output)
		if err != nil {
			// Error is already returned by packStdTxResponse.
			return
		}

		rest.PostProcessResponseBare(w, clientCtx, result)
	}
}

// packStdTxResponse takes a sdk.TxResponse, converts the Tx into a StdTx, and
// packs the StdTx again into the sdk.TxResponse Any. Amino then takes care of
// seamlessly JSON-outputting the Any.
func packStdTxResponse(w http.ResponseWriter, clientCtx client.Context, txRes *sdk.TxResponse) (*InfoCoinFlow, error) {
	// We just unmarshalled from Tendermint, we take the proto Tx's raw
	// bytes, and convert them into a StdTx to be displayed.
	txBytes := txRes.Tx.Value
	stdTx, err := convertToStdTx(w, clientCtx, txBytes)
	if err != nil {
		return nil, err
	}
	signatures := make([]types.StdSignature, len(stdTx.Signatures))
	for k, v := range stdTx.Signatures {
		signatures[k].Signature = v.Signature
		signatures[k].PubKey = v.PubKey
	}
	result := ResponseDeliverTx{
		Code:      txRes.Code,
		Data:      txRes.Data,
		Log:       txRes.RawLog,
		Info:      txRes.Info,
		GasWanted: txRes.GasWanted,
		GasUsed:   txRes.GasUsed,
		Tags:      ConvertLogsToTags(txRes.Logs),
		Codespace: txRes.Codespace,
	}
	return &InfoCoinFlow{
		Hash:   txRes.TxHash,
		Height: txRes.Height,
		Tx: types.StdTx{
			Msgs:       stdTx.Msgs,
			Fee:        stdTx.Fee,
			Signatures: signatures,
			Memo:       stdTx.Memo,
		},
		Result:    result,
		Timestamp: txRes.Timestamp,
	}, nil
}

// convertToStdTx converts tx proto binary bytes retrieved from Tendermint into
// a StdTx. Returns the StdTx, as well as a flag denoting if the function
// successfully converted or not.
func convertToStdTx(w http.ResponseWriter, clientCtx client.Context, txBytes []byte) (legacytx.StdTx, error) {
	txI, err := clientCtx.TxConfig.TxDecoder()(txBytes)
	if rest.CheckBadRequestError(w, err) {
		return legacytx.StdTx{}, err
	}

	tx, ok := txI.(signing.Tx)
	if !ok {
		rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("%+v is not backwards compatible with %T", tx, legacytx.StdTx{}))
		return legacytx.StdTx{}, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "expected %T, got %T", (signing.Tx)(nil), txI)
	}

	stdTx, err := clienttx.ConvertTxToStdTx(clientCtx.LegacyAmino, tx)
	if rest.CheckBadRequestError(w, err) {
		return legacytx.StdTx{}, err
	}

	return stdTx, nil
}

func ConvertEventToTags(event sdk.StringEvent) []ReadableTag {
	readableTags := make([]ReadableTag, len(event.Attributes))
	for i, kv := range event.Attributes {
		readableTags[i] = ReadableTag{
			Key:   event.Type + "." + string(kv.Key),
			Value: kv.Value,
		}
	}
	return readableTags
}

func ConvertEventsToTags(events sdk.StringEvents) []ReadableTag {
	var readableTags []ReadableTag
	for _, kv := range events {
		readableTags = append(readableTags, ConvertEventToTags(kv)...)
	}
	return readableTags
}

func ConvertLogsToTags(logs sdk.ABCIMessageLogs) []ReadableTag {
	var readableTags []ReadableTag
	for _, kv := range logs {
		readableTags = append(readableTags, ConvertEventsToTags(kv.Events)...)
	}
	return readableTags
}
