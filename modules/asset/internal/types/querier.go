package types

import (
	"fmt"
	"strings"

	sdk "github.com/irisnet/irishub/types"
)

const (
	QueryToken    = "token"
	QueryTokens   = "tokens"
	QueryGateway  = "gateway"
	QueryGateways = "gateways"
	QueryFees     = "fees"
)

// QueryTokenParams is the query parameters for 'custom/asset/tokens/{id}'
type QueryTokenParams struct {
	TokenId string
}

// QueryTokensParams is the query parameters for 'custom/asset/tokens'
type QueryTokensParams struct {
	Source  string
	Gateway string
	Owner   string
}

// QueryGatewayParams is the query parameters for 'custom/asset/gateway'
type QueryGatewayParams struct {
	Moniker string
}

// QueryGatewaysParams is the query parameters for 'custom/asset/gateways'
type QueryGatewaysParams struct {
	Owner sdk.AccAddress
}

// QueryGatewayFeeParams is the query parameters for 'custom/asset/fees/gateways'
type QueryGatewayFeeParams struct {
	Moniker string
}

// QueryTokenFeesParams is the query parameters for 'custom/asset/fees/tokens'
type QueryTokenFeesParams struct {
	ID string
}

// GatewayFeeOutput is for the gateway fee query output
type GatewayFeeOutput struct {
	Exist bool     `json:"exist"` // indicate if the gateway has existed
	Fee   sdk.Coin `json:"fee"`   // creation fee
}

// String implements stringer
func (gfo GatewayFeeOutput) String() string {
	var out strings.Builder
	if gfo.Exist {
		out.WriteString("The gateway moniker has existed\n")
	}

	out.WriteString(fmt.Sprintf("Fee: %s", gfo.Fee.String()))

	return out.String()
}

// HumanString implements human
func (gfo GatewayFeeOutput) HumanString(converter sdk.CoinsConverter) string {
	var out strings.Builder
	if gfo.Exist {
		out.WriteString("The gateway moniker has existed\n")
	}

	out.WriteString(fmt.Sprintf("Fee: %s", converter.ToMainUnit(sdk.Coins{gfo.Fee})))

	return out.String()
}

// TokenFeesOutput is for the token fees query output
type TokenFeesOutput struct {
	Exist    bool     `json:"exist"`     // indicate if the token has existed
	IssueFee sdk.Coin `json:"issue_fee"` // issue fee
	MintFee  sdk.Coin `json:"mint_fee"`  // mint fee
}

// String implements stringer
func (tfo TokenFeesOutput) String() string {
	var out strings.Builder
	if tfo.Exist {
		out.WriteString("The token id has existed\n")
	}

	out.WriteString(fmt.Sprintf(`Fees:
  IssueFee: %s
  MintFee:  %s`,
		tfo.IssueFee.String(), tfo.MintFee.String()))

	return out.String()
}

// String implements human
func (tfo TokenFeesOutput) HumanString(converter sdk.CoinsConverter) string {
	var out strings.Builder
	if tfo.Exist {
		out.WriteString("The token id has existed\n")
	}

	out.WriteString(fmt.Sprintf(`Fees:
  IssueFee: %s
  MintFee:  %s`,
		converter.ToMainUnit(sdk.Coins{tfo.IssueFee}),
		converter.ToMainUnit(sdk.Coins{tfo.MintFee})))

	return out.String()
}
