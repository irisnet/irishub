package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QueryDefinition = "definition"
	QueryBinding    = "binding"
	QueryBindings   = "bindings"
	QueryRequests   = "requests"
	QueryResponse   = "response"
	QueryFees       = "fees"
)

type QueryDefinitionParams struct {
	DefChainID  string
	ServiceName string
}

type DefinitionOutput struct {
	Definition SvcDef           `json:"definition" yaml:"definition"`
	Methods    []MethodProperty `json:"methods" yaml:"methods"`
}

type QueryBindingParams struct {
	DefChainID  string
	ServiceName string
	BindChainId string
	Provider    sdk.AccAddress
}

type QueryBindingsParams struct {
	DefChainID  string
	ServiceName string
}

type QueryRequestsParams struct {
	DefChainID  string
	ServiceName string
	BindChainId string
	Provider    sdk.AccAddress
}

type QueryResponseParams struct {
	ReqChainId string
	RequestId  string
}

type QueryFeesParams struct {
	Address sdk.AccAddress
}

type FeesOutput struct {
	ReturnedFee sdk.Coins `json:"returned_fee" yaml:"returned_fee"`
	IncomingFee sdk.Coins `json:"incoming_fee" yaml:"incoming_fee"`
}
