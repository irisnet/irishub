package bank

import (
	"fmt"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the auth Querier
const (
	QueryAccount    = "account"
	QueryTokenStats = "tokenStats"
)

// creates a querier for auth REST endpoints
func NewQuerier(keeper Keeper, cdc *codec.Codec) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryAccount:
			return queryAccount(ctx, req, keeper, cdc)
		case QueryTokenStats:
			return queryTokenStats(ctx, req, keeper, cdc)

		default:
			return nil, sdk.ErrUnknownRequest("unknown auth query endpoint")
		}
	}
}

// defines the params for query: "custom/acc/account"
type QueryAccountParams struct {
	Address sdk.AccAddress
}

func NewQueryAccountParams(addr sdk.AccAddress) QueryAccountParams {
	return QueryAccountParams{
		Address: addr,
	}
}

func queryAccount(ctx sdk.Context, req abci.RequestQuery, keeper Keeper, cdc *codec.Codec) ([]byte, sdk.Error) {
	var params QueryAccountParams
	bk := keeper.(BaseKeeper)
	if err := cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	account := bk.am.GetAccount(ctx, params.Address)
	if account == nil {
		return nil, sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", params.Address))
	}

	bz, err := codec.MarshalJSONIndent(cdc, account)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}

	return bz, nil
}

// defines the params for query: "custom/bank/token-stats"
type QueryTokenStatsParams struct {
	TokenId string
}

func queryTokenStats(ctx sdk.Context, req abci.RequestQuery, keeper Keeper, cdc *codec.Codec) ([]byte, sdk.Error) {
	var params QueryTokenStatsParams
	bk := keeper.(BaseKeeper)
	if err := cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	// TODO: query total supply by id

	irisBurnedToken := sdk.Coin{}
	irisBurnedToken.Denom = sdk.NativeTokenMinDenom
	irisBurnedToken.Amount = bk.GetCoins(ctx, BurnedCoinsAccAddr).AmountOf(sdk.NativeTokenMinDenom)
	tokenStats := TokenStats{
		LooseTokens:  bk.GetLoosenCoins(ctx),
		BurnedTokens: sdk.Coins{irisBurnedToken},
	}
	bz, err := codec.MarshalJSONIndent(cdc, tokenStats)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}

	return bz, nil
}

type TokenStats struct {
	LooseTokens  sdk.Coins `json:"loose_tokens"`
	BurnedTokens sdk.Coins `json:"burned_tokens"`
	BondedTokens sdk.Coins `json:"bonded_tokens"`
	TotalSupply  sdk.Coins `json:"total_supply"`
}
