package utils

import (
	"fmt"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/crypto/keys/keyerror"
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/modules/auth"
	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/client/context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const (
	Async                = "async"
	queryArgDryRun       = "simulate"
	queryArgGenerateOnly = "generate-only"
)

//----------------------------------------
// Basic HTTP utilities

// WriteErrorResponse prepares and writes a HTTP error
// given a status code and an error message.
func WriteErrorResponse(w http.ResponseWriter, status int, err string) {
	w.WriteHeader(status)
	w.Write([]byte(err))
}

// WriteSimulationResponse prepares and writes an HTTP
// response for transactions simulations.
type kvPair struct {
	TagKey   string `json:"tag_key"`
	TagValue string `json:"tag_value"`
}
type abciResult struct {
	Code      sdk.CodeType     `json:"code"`
	Data      []byte           `json:"data"`
	Log       string           `json:"log"`
	GasWanted uint64           `json:"gas_wanted"`
	GasUsed   uint64           `json:"gas_used"`
	FeeAmount int64            `json:"fee_amount"`
	FeeDenom  string           `json:"fee_denom"`
	Tags      []kvPair         `json:"tagsy"`
}
type simulateResult struct {
	GasEstimate uint64      `json:"gas_estimate"`
	Result      abciResult `json:"result"`
}

func WriteSimulationResponse(w http.ResponseWriter, cliCtx context.CLIContext, gas uint64, result sdk.Result) {
	w.WriteHeader(http.StatusOK)
	var kvPairs []kvPair
	for _, tag := range result.Tags {
		kvPairs = append(kvPairs, kvPair{
			TagKey:   string(tag.Key),
			TagValue: string(tag.Value),
		})
	}
	abciResult := abciResult{
		Code:      result.Code,
		Data:      result.Data,
		Log:       result.Log,
		GasWanted: result.GasWanted,
		GasUsed:   result.GasUsed,
		FeeAmount: result.FeeAmount,
		FeeDenom:  result.FeeDenom,
		Tags:      kvPairs,
	}

	simulateResult := simulateResult{
		GasEstimate: gas,
		Result:      abciResult,
	}
	var output []byte
	var err error
	if cliCtx.Indent {
		output, err = cliCtx.Codec.MarshalJSONIndent(simulateResult, "", "  ")
	} else {
		output, err = cliCtx.Codec.MarshalJSON(simulateResult)
	}
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Write(output)
}

// HasDryRunArg returns true if the request's URL query contains the dry run
// argument and its value is set to "true".
func HasDryRunArg(r *http.Request) bool {
	return urlQueryHasArg(r.URL, queryArgDryRun)
}

// HasGenerateOnlyArg returns whether a URL's query "generate-only" parameter
// is set to "true".
func HasGenerateOnlyArg(r *http.Request) bool {
	return urlQueryHasArg(r.URL, queryArgGenerateOnly)
}

// AsyncOnlyArg returns whether a URL's query "async" parameter
func AsyncOnlyArg(r *http.Request) bool {
	return urlQueryHasArg(r.URL, Async)
}

// ParseInt64OrReturnBadRequest converts s to a int64 value.
func ParseInt64OrReturnBadRequest(w http.ResponseWriter, s string) (n int64, ok bool) {
	var err error

	n, err = strconv.ParseInt(s, 10, 64)
	if err != nil {
		err := fmt.Errorf("'%s' is not a valid int64", s)
		WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return n, false
	}

	return n, true
}

// ParseUint64OrReturnBadRequest converts s to a uint64 value.
func ParseUint64OrReturnBadRequest(w http.ResponseWriter, s string) (n uint64, ok bool) {
	var err error
	n, err = strconv.ParseUint(s, 10, 64)
	if err != nil {
		err := fmt.Errorf("'%s' is not a valid uint64", s)
		WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return n, false
	}
	return n, true
}

// ParseFloat64OrReturnBadRequest converts s to a float64 value. It returns a
// default value, defaultIfEmpty, if the string is empty.
func ParseFloat64OrReturnBadRequest(w http.ResponseWriter, s string, defaultIfEmpty float64) (n float64, ok bool) {
	if len(s) == 0 {
		return defaultIfEmpty, true
	}

	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return n, false
	}

	return n, true
}

// WriteGenerateStdTxResponse writes response for the generate_only mode.
func WriteGenerateStdTxResponse(w http.ResponseWriter, txCtx TxContext, msgs []sdk.Msg) {
	stdMsg, err := txCtx.Build(msgs)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	output, err := txCtx.Codec.MarshalJSON(auth.NewStdTx(stdMsg.Msgs, stdMsg.Fee, nil, stdMsg.Memo))
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Write(output)
	return
}

func urlQueryHasArg(url *url.URL, arg string) bool { return url.Query().Get(arg) == "true" }

// ReadPostBody
func ReadPostBody(w http.ResponseWriter, r *http.Request, cdc *codec.Codec, req interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("invalid post body")
			WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}
	}()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return err
	}

	err = cdc.UnmarshalJSON(body, req)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return err
	}

	return nil
}

// InitReqCliCtx
func InitReqCliCtx(cliCtx context.CLIContext, r *http.Request) context.CLIContext {
	cliCtx.GenerateOnly = HasGenerateOnlyArg(r)
	cliCtx.Async = AsyncOnlyArg(r)
	cliCtx.DryRun = HasDryRunArg(r)
	return cliCtx
}

// SendOrReturnUnsignedTx implements a utility function that facilitates
// sending a series of messages in a signed transaction given a TxBuilder and a
// QueryContext. It ensures that the account exists, has a proper number and
// sequence set. In addition, it builds and signs a transaction with the
// supplied messages. Finally, it broadcasts the signed transaction to a node.
//
// NOTE: Also see SendOrPrintTx.
// NOTE: Also see x/stake/client/rest/tx.go delegationsRequestHandlerFn.
func SendOrReturnUnsignedTx(w http.ResponseWriter, cliCtx context.CLIContext, baseTx BaseTx, msgs []sdk.Msg) {

	simulateGas, gas, err := client.ReadGasFlag(baseTx.Gas)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	adjustment, ok := ParseFloat64OrReturnBadRequest(w, baseTx.GasAdjustment, client.DefaultGasAdjustment)
	if !ok {
		return
	}

	txCtx := TxContext{
		Codec:         cliCtx.Codec,
		Gas:           gas,
		Fee:           baseTx.Fee,
		GasAdjustment: adjustment,
		SimulateGas:   simulateGas,
		ChainID:       baseTx.ChainID,
		AccountNumber: baseTx.AccountNumber,
		Sequence:      baseTx.Sequence,
	}
	txCtx = txCtx.WithCliCtx(cliCtx)

	if cliCtx.GenerateOnly {
		WriteGenerateStdTxResponse(w, txCtx, msgs)
		return
	}

	if cliCtx.DryRun || txCtx.SimulateGas {
		newTxCtx, result, err := EnrichCtxWithGas(txCtx, cliCtx, baseTx.Name, msgs)
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if cliCtx.DryRun {
			WriteSimulationResponse(w, cliCtx, newTxCtx.Gas, result)
			return
		}

		txCtx = newTxCtx
	}

	txBytes, err := txCtx.BuildAndSign(baseTx.Name, baseTx.Password, msgs)
	if keyerror.IsErrKeyNotFound(err) {
		WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	} else if keyerror.IsErrWrongPassword(err) {
		WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	} else if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	res, err := cliCtx.BroadcastTx(txBytes)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	PostProcessResponse(w, cliCtx.Codec, res, cliCtx.Indent)
}

// PostProcessResponse performs post process for rest response
func PostProcessResponse(w http.ResponseWriter, cdc *codec.Codec, response interface{}, indent bool) {
	var output []byte
	switch response.(type) {
	default:
		var err error
		if indent {
			output, err = cdc.MarshalJSONIndent(response, "", "  ")
		} else {
			output, err = cdc.MarshalJSON(response)
		}
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	case []byte:
		output = response.([]byte)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}
