package types

import (
	"fmt"
	"strings"

	sdk "github.com/irisnet/irishub/types"
)

const (
	QueryToken  = "token"
	QueryTokens = "tokens"
	QueryFees   = "fees"
)

// QueryTokenParams is the query parameters for 'custom/asset/tokens/{id}'
type QueryTokenParams struct {
	TokenId string
}

// QueryTokensParams is the query parameters for 'custom/asset/tokens'
type QueryTokensParams struct {
	TokenID string
	Owner   string
}

// QueryTokenFeesParams is the query parameters for 'custom/asset/fees/tokens'
type QueryTokenFeesParams struct {
	Symbol string
}

type TokenOutput struct {
	Id            string         `json:"id"`
	Symbol        string         `json:"symbol"`
	Name          string         `json:"name"`
	Scale         uint8          `json:"scale"`
	MinUnit       string         `json:"min_unit"`
	InitialSupply sdk.Int        `json:"initial_supply"`
	MaxSupply     sdk.Int        `json:"max_supply"`
	Mintable      bool           `json:"mintable"`
	Owner         sdk.AccAddress `json:"owner"`
}

// String implements stringer
func (top TokenOutput) String() string {
	token := top.ToFungibleToken()
	return token.String()
}

func (top TokenOutput) ToFungibleToken() FungibleToken {
	return FungibleToken{BaseToken{
		Id:            top.Id,
		Symbol:        top.Symbol,
		Name:          top.Name,
		Decimal:       top.Scale,
		MinUnitAlias:  top.MinUnit,
		InitialSupply: top.InitialSupply,
		MaxSupply:     top.MaxSupply,
		Mintable:      top.Mintable,
		Owner:         top.Owner,
	}}
}

func NewTokenOutputFrom(token FungibleToken) TokenOutput {
	return TokenOutput{
		Id:            token.Id,
		Symbol:        token.Symbol,
		Name:          token.Name,
		Scale:         token.Decimal,
		MinUnit:       token.MinUnitAlias,
		InitialSupply: token.InitialSupply,
		MaxSupply:     token.MaxSupply,
		Mintable:      token.Mintable,
		Owner:         token.Owner,
	}
}

type TokensOutput []TokenOutput

func (tsop TokensOutput) String() string {
	var tokens Tokens
	for _, t := range tsop {
		tokens = append(tokens, t.ToFungibleToken())
	}
	if len(tokens) == 0 {
		return ""
	}
	return tokens.String()
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
