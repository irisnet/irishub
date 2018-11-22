package service

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type SvcRequest struct {
	DefChainID            string         `json:"def_chain_id"`
	DefName               string         `json:"def_name"`
	BindChainID           string         `json:"bind_chain_id"`
	ReqChainID            string         `json:"req_chain_id"`
	MethodID              int64          `json:"method_id"`
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
