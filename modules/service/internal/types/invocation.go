package types

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type SvcRequest struct {
	DefChainID            string         `json:"def_chain_id" yaml:"def_chain_id"`
	DefName               string         `json:"def_name" yaml:"def_name" `
	BindChainID           string         `json:"bind_chain_id" yaml:"bind_chain_id`
	ReqChainID            string         `json:"req_chain_id" yaml:"req_chain_id"`
	MethodID              int16          `json:"method_id" yaml:"method_id"`
	Provider              sdk.AccAddress `json:"provider" yaml:"provider"`
	Consumer              sdk.AccAddress `json:"consumer" yaml:"consumer"`
	Input                 []byte         `json:"input" yaml:"input"`
	ServiceFee            sdk.Coins      `json:"service_fee" yaml:"service_fee"`
	Profiling             bool           `json:"profiling" yaml:"profiling"`                               // profiling model will be free of service charges
	RequestHeight         int64          `json:"request_height" yaml:"request_height"`                     // block height of service request
	RequestIntraTxCounter int16          `json:"request_intra_tx_counter" yaml:"request_intra_tx_counter"` // block-local tx index of service request
	ExpirationHeight      int64          `json:"expiration_height" yaml:"expiration_height"`               // block height of the service request has expired
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
	ReqChainID            string         `json:"req_chain_id" yaml:"req_chain_id"`
	RequestHeight         int64          `json:"request_height" yaml:"request_height"`
	RequestIntraTxCounter int16          `json:"request_intra_tx_counter" yaml:"request_intra_tx_counter"`
	ExpirationHeight      int64          `json:"expiration_height" yaml:"expiration_height"`
	Provider              sdk.AccAddress `json:"provider" yaml:"provider"`
	Consumer              sdk.AccAddress `json:"consumer" yaml:"consumer"`
	Output                []byte         `json:"output" yaml:"output"`
	ErrorMsg              []byte         `json:"error_msg" yaml:"error_msg"`
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
	Address sdk.AccAddress `json:"address" yaml:"address"`
	Coins   sdk.Coins      `json:"coins" yaml:"coins"`
}

func NewReturnedFee(address sdk.AccAddress, coins sdk.Coins) ReturnedFee {
	return ReturnedFee{
		Address: address,
		Coins:   coins,
	}
}

// incoming fee of a consumer
type IncomingFee struct {
	Address sdk.AccAddress `json:"address" yaml:"address"`
	Coins   sdk.Coins      `json:"coins" yaml:"coins"`
}

func NewIncomingFee(address sdk.AccAddress, coins sdk.Coins) IncomingFee {
	return IncomingFee{
		Address: address,
		Coins:   coins,
	}
}
