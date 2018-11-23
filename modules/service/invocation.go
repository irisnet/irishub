package service

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"fmt"
	"strings"
	"errors"
	"strconv"
)

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

// RequestID is of format request height-intraTxCounter
func (req SvcRequest) RequestID() string {
	return fmt.Sprintf("%d-%d", req.RequestHeight, req.RequestIntraTxCounter)
}

func TransferRequestID(requestId string) (height int64, counter int16, err error) {
	ss := strings.Split(requestId, "-")
	if len(ss) != 2 {
		return height, counter, errors.New("invalid request id")
	}
	height, err = strconv.ParseInt(ss[0], 10, 64)
	if err != nil {
		return height, counter, err
	}
	counterInt, err := strconv.Atoi(ss[0])
	if err != nil {
		return height, counter, err
	}
	return height, int16(counterInt), err
}

type SvcResponse struct {
	ReqChainID            string         `json:"req_chain_id"`
	RequestHeight         int64          `json:"request_height"`
	RequestIntraTxCounter int16          `json:"request_intra_tx_counter"`
	Provider              sdk.AccAddress `json:"provider"`
	Consumer              sdk.AccAddress `json:"consumer"`
	Output                []byte         `json:"output"`
	ErrorMsg              []byte         `json:"error_msg"`
}

func NewSvcResponse(reqChainID string, height int64, counter int16, provider, consumer sdk.AccAddress, out []byte, errorMsg []byte) SvcResponse {
	return SvcResponse{
		ReqChainID:            reqChainID,
		RequestHeight:         height,
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
