package types

import (
	"encoding/binary"
	"errors"
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
	MaxServiceFee     sdk.Coins           `json:"max_service_fee"`
	Profiling         bool                `json:"profiling"`
	Timeout           int64               `json:"timeout"`
	Repeated          bool                `json:"repeated"`
	RepeatedFrequency uint64              `json:"repeated_frequency"`
	RepeatedTotal     int64               `json:"repeated_total"`
	BatchCounter      uint64              `json:"batch_counter"`
	State             RequestContextState `json:"state"`
	ResponseThreshold uint16              `json:"response_threshold"`
	ResponseHandler   string              `json:"response_handler"`
}

// NewRequestContext creates a new RequestContext instance
func NewRequestContext(
	serviceName string,
	providers []sdk.AccAddress,
	consumer sdk.AccAddress,
	input string,
	maxServiceFee sdk.Coins,
	profiling bool,
	timeout int64,
	repeated bool,
	repeatedFrequency uint64,
	repeatedTotal int64,
	batchCounter uint64,
	state RequestContextState,
	responseThreshold uint16,
	responseHandler string,
) RequestContext {
	return RequestContext{
		ServiceName:       serviceName,
		Providers:         providers,
		Consumer:          consumer,
		Input:             input,
		MaxServiceFee:     maxServiceFee,
		Profiling:         profiling,
		Timeout:           timeout,
		Repeated:          repeated,
		RepeatedFrequency: repeatedFrequency,
		RepeatedTotal:     repeatedTotal,
		BatchCounter:      batchCounter,
		State:             state,
		ResponseThreshold: responseThreshold,
		ResponseHandler:   responseHandler,
	}
}

// CompactRequest defines a compact request with a request context ID
type CompactRequest struct {
	RequestContextID           []byte
	RequestContextBatchCounter uint64
	Provider                   sdk.AccAddress
	ServiceFee                 sdk.Coins
	RequestHeight              int64
}

// NewCompactRequest creates a new CompactRequest instance
func NewCompactRequest(
	requestContextID []byte,
	batchCounter uint64,
	provider sdk.AccAddress,
	serviceFee sdk.Coins,
	requestHeight int64,
) CompactRequest {
	return CompactRequest{
		RequestContextID:           requestContextID,
		RequestContextBatchCounter: batchCounter,
		Provider:                   provider,
		ServiceFee:                 serviceFee,
		RequestHeight:              requestHeight,
	}
}

// Request defines a request which contains the detailed request data
type Request struct {
	ServiceName                string         `json:"service_name"`
	Provider                   sdk.AccAddress `json:"provider"`
	Consumer                   sdk.AccAddress `json:"consumer"`
	Input                      string         `json:"input"`
	ServiceFee                 sdk.Coins      `json:"service_fee"`
	Profiling                  bool           `json:"profiling"`
	RequestHeight              int64          `json:"request_height"`
	ExpirationHeight           int64          `json:"expiration_height"`
	RequestContextID           []byte         `json:"request_context_id"`
	RequestContextBatchCounter uint64         `json:"request_context_batch_counter"`
}

// NewRequest creates a new Request instance
func NewRequest(
	serviceName string,
	provider,
	consumer sdk.AccAddress,
	input string,
	serviceFee sdk.Coins,
	profiling bool,
	requestHeight int64,
	expirationHeight int64,
	requestContextID []byte,
	batchCounter uint64,
) Request {
	return Request{
		ServiceName:                serviceName,
		Provider:                   provider,
		Consumer:                   consumer,
		Input:                      input,
		ServiceFee:                 serviceFee,
		Profiling:                  profiling,
		RequestHeight:              requestHeight,
		ExpirationHeight:           expirationHeight,
		RequestContextID:           requestContextID,
		RequestContextBatchCounter: batchCounter,
	}
}

// Response defines a response
type Response struct {
	Provider                   sdk.AccAddress `json:"provider"`
	Consumer                   sdk.AccAddress `json:"consumer"`
	Output                     string         `json:"output"`
	Error                      string         `json:"error"`
	RequestContextID           []byte         `json:"request_context_id"`
	RequestContextBatchCounter uint64         `json:"request_context_batch_counter"`
}

// NewResponse creates a new Response instance
func NewResponse(
	provider,
	consumer sdk.AccAddress,
	output,
	err string,
	requestContextID []byte,
	batchCounter uint64,
) Response {
	return Response{
		Provider:                   provider,
		Consumer:                   consumer,
		Output:                     output,
		Error:                      err,
		RequestContextID:           requestContextID,
		RequestContextBatchCounter: batchCounter,
	}
}

// EarnedFees defines a struct for the fees earned by the provider
type EarnedFees struct {
	Address sdk.AccAddress `json:"address"`
	Coins   sdk.Coins      `json:"coins"`
}

// NewEarnedFees creates a new EarnedFees instance
func NewEarnedFees(address sdk.AccAddress, coins sdk.Coins) EarnedFees {
	return EarnedFees{
		Address: address,
		Coins:   coins,
	}
}

// RequestContextState defines the state for the request context
type RequestContextState byte

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

// GenerateRequestContextID generates a unique ID for the request context from the specified params
func GenerateRequestContextID(blockHeight int64, intraCounter int16) []byte {
	bz := make([]byte, 10)

	binary.BigEndian.PutUint64(bz, uint64(blockHeight))
	binary.BigEndian.PutUint16(bz[8:], uint16(intraCounter))

	return sdk.SHA256(bz)
}
