package types

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"strings"

	sdk "github.com/irisnet/irishub/types"
)

// RequestContext defines a context which holds request-related data
type RequestContext struct {
	ServiceName       string              `json:"service_name"`
	Providers         []sdk.AccAddress    `json:"providers"`
	Consumer          sdk.AccAddress      `json:"consumer"`
	Input             string              `json:"input"`
	ServiceFeeCap     sdk.Coins           `json:"service_fee_cap"`
	Timeout           int64               `json:"timeout"`
	Repeated          bool                `json:"repeated"`
	RepeatedFrequency uint64              `json:"repeated_frequency"`
	RepeatedTotal     int64               `json:"repeated_total"`
	BatchCounter      uint64              `json:"batch_counter"`
	State             RequestContextState `json:"state"`
	ResponseThreshold uint16              `json:"response_threshold"`
	ModuleName        string              `json:"module_name"`
}

// NewRequestContext creates a new RequestContext instance
func NewRequestContext(
	serviceName string,
	providers []sdk.AccAddress,
	consumer sdk.AccAddress,
	input string,
	serviceFeeCap sdk.Coins,
	timeout int64,
	repeated bool,
	repeatedFrequency uint64,
	repeatedTotal int64,
	batchCounter uint64,
	state RequestContextState,
	responseThreshold uint16,
	moduleName string,
) RequestContext {
	return RequestContext{
		ServiceName:       serviceName,
		Providers:         providers,
		Consumer:          consumer,
		Input:             input,
		ServiceFeeCap:     serviceFeeCap,
		Timeout:           timeout,
		Repeated:          repeated,
		RepeatedFrequency: repeatedFrequency,
		RepeatedTotal:     repeatedTotal,
		BatchCounter:      batchCounter,
		State:             state,
		ResponseThreshold: responseThreshold,
		ModuleName:        moduleName,
	}
}

// RequestContextState represents the request context state
type RequestContextState byte

const (
	STARTED RequestContextState = 0x00 // started
	PAUSED  RequestContextState = 0x01 // paused
)

// ResponseCallback defines the response callback interface
type ResponseCallback func(ctx sdk.Context, requestContextID []byte, reponses []string)

// GenerateRequestContextID generates a unique ID for the request context from the specified params
func GenerateRequestContextID(blockHeight int64, intraCounter int16) []byte {
	bz := make([]byte, 10)

	binary.BigEndian.PutUint64(bz, uint64(blockHeight))
	binary.BigEndian.PutUint16(bz[8:], uint16(intraCounter))

	return sdk.SHA256(bz)
}

type SvcRequest struct {
	DefChainID            string         `json:"def_chain_id"`
	DefName               string         `json:"def_name"`
	BindChainID           string         `json:"bind_chain_id"`
	ReqChainID            string         `json:"req_chain_id"`
	MethodID              int16          `json:"method_id"`
	Provider              sdk.AccAddress `json:"provider"`
	Consumer              sdk.AccAddress `json:"consumer"`
	Input                 []byte         `json:"input"`
	ServiceFee            sdk.Coins      `json:"service_fee"`
	Profiling             bool           `json:"profiling"`                // profiling model will be free of service charges
	RequestHeight         int64          `json:"request_height"`           // block height of service request
	RequestIntraTxCounter int16          `json:"request_intra_tx_counter"` // block-local tx index of service request
	ExpirationHeight      int64          `json:"expiration_height"`        // block height of the service request has expired
}

func NewSvcRequest(defChainID, defName, bindChainID, reqChainID string, consumer, provider sdk.AccAddress, methodID int16, input []byte, serviceFee sdk.Coins, profiling bool) SvcRequest {
	return SvcRequest{
		DefChainID:  defChainID,
		DefName:     defName,
		BindChainID: bindChainID,
		ReqChainID:  reqChainID,
		MethodID:    methodID,
		Provider:    provider,
		Consumer:    consumer,
		Input:       input,
		ServiceFee:  serviceFee,
		Profiling:   profiling,
	}
}

// RequestID is of format request expirationHeight-requestHeight-intraTxCounter
func (req SvcRequest) RequestID() string {
	return fmt.Sprintf("%d-%d-%d", req.ExpirationHeight, req.RequestHeight, req.RequestIntraTxCounter)
}

func ConvertRequestID(requestId string) (eHeight int64, rHeight int64, counter int16, err error) {
	ss := strings.Split(requestId, "-")
	if len(ss) != 3 {
		return eHeight, rHeight, counter, errors.New("invalid request id")
	}
	eHeight, err = strconv.ParseInt(ss[0], 10, 64)
	if err != nil {
		return eHeight, rHeight, counter, err
	}
	rHeight, err = strconv.ParseInt(ss[1], 10, 64)
	if err != nil {
		return eHeight, rHeight, counter, err
	}
	counterInt, err := strconv.Atoi(ss[2])
	if err != nil {
		return eHeight, rHeight, counter, err
	}
	return eHeight, rHeight, int16(counterInt), err
}

type SvcResponse struct {
	ReqChainID            string         `json:"req_chain_id"`
	RequestHeight         int64          `json:"request_height"`
	RequestIntraTxCounter int16          `json:"request_intra_tx_counter"`
	ExpirationHeight      int64          `json:"expiration_height"`
	Provider              sdk.AccAddress `json:"provider"`
	Consumer              sdk.AccAddress `json:"consumer"`
	Output                []byte         `json:"output"`
	ErrorMsg              []byte         `json:"error_msg"`
}

func NewSvcResponse(reqChainID string, eheight int64, rheight int64, counter int16, provider, consumer sdk.AccAddress, out []byte, errorMsg []byte) SvcResponse {
	return SvcResponse{
		ReqChainID:            reqChainID,
		ExpirationHeight:      eheight,
		RequestHeight:         rheight,
		RequestIntraTxCounter: counter,
		Provider:              provider,
		Consumer:              consumer,
		Output:                out,
		ErrorMsg:              errorMsg,
	}
}

// return fee of a consumer
type ReturnedFee struct {
	Address sdk.AccAddress `json:"address"`
	Coins   sdk.Coins      `json:"coins"`
}

func NewReturnedFee(address sdk.AccAddress, coins sdk.Coins) ReturnedFee {
	return ReturnedFee{
		Address: address,
		Coins:   coins,
	}
}

// incoming fee of a consumer
type IncomingFee struct {
	Address sdk.AccAddress `json:"address"`
	Coins   sdk.Coins      `json:"coins"`
}

func NewIncomingFee(address sdk.AccAddress, coins sdk.Coins) IncomingFee {
	return IncomingFee{
		Address: address,
		Coins:   coins,
	}
}
