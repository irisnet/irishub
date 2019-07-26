package bank

import (
	"fmt"
	"github.com/irisnet/irishub/app/v1/auth"
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
		return nil, sdk.ErrUnknownAddress(params.Address.String())
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

	if err := cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	bk := keeper.(BaseKeeper)
	looseTokens := sdk.Coins{}
	burnedTokens := sdk.Coins{}
	totalSupplies := sdk.Coins{}
	// TODO: bonded tokens for iris

	if params.TokenId == "" { // query all
		looseTokens = bk.GetLoosenCoins(ctx)
		burnedTokens = bk.GetCoins(ctx, auth.BurnedCoinsAccAddr)
		iter := bk.am.GetTotalSupplies(ctx)
		defer iter.Close()
		for ; iter.Valid(); iter.Next() {
			var ts sdk.Coin
			cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &ts)
			totalSupplies = append(totalSupplies, ts)
		}

	} else if params.TokenId == sdk.Iris { // query iris
		looseTokens = bk.GetLoosenCoins(ctx)
		burnedTokens = sdk.Coins{sdk.NewCoin(sdk.IrisAtto, bk.GetCoins(ctx, auth.BurnedCoinsAccAddr).AmountOf(sdk.IrisAtto))}
	} else { // query !iris
		denom, err := sdk.GetCoinMinDenom(params.TokenId)
		if err != nil {
			return nil, sdk.ParseParamsErr(err)
		}
		burnedTokens = sdk.Coins{sdk.NewCoin(denom, bk.GetCoins(ctx, auth.BurnedCoinsAccAddr).AmountOf(denom))}

		ts, found := bk.GetTotalSupply(ctx, denom)
		if !found {
			return nil, sdk.ErrUnknownRequest("unknown token id")
		}

		totalSupplies = append(totalSupplies, ts)
	}

	tokenStats := TokenStats{
		LooseTokens:  looseTokens,
		BurnedTokens: burnedTokens,
		TotalSupply:  totalSupplies,
	}
	bz, err := codec.MarshalJSONIndent(cdc, tokenStats)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}

	return bz, nil
}

type TokenStats struct {
	LooseTokens  sdk.Coins `json:"loose_tokens"`
	BondedTokens sdk.Coins `json:"bonded_tokens"`
	BurnedTokens sdk.Coins `json:"burned_tokens"`
	TotalSupply  sdk.Coins `json:"total_supply"`
}

func (ts TokenStats) String() string {
	return fmt.Sprintf(`TokenStats:
  Loose Tokens:             %s
  Bonded Tokens:            %s
  Burned Tokens:            %s
  Total Supply:             %s`,
		ts.LooseTokens.String(), ts.BondedTokens.String(), ts.BurnedTokens.String(), ts.TotalSupply.String())
}

func (ts TokenStats) HumanString(converter sdk.CoinsConverter) string {
	return fmt.Sprintf(`TokenStats:
  Loose Tokens:             %s
  Bonded Tokens:            %s
  Burned Tokens:            %s
  Total Supply:             %s`,
		converter.ToMainUnit(ts.LooseTokens),
		converter.ToMainUnit(ts.BondedTokens),
		converter.ToMainUnit(ts.BurnedTokens),
		converter.ToMainUnit(ts.TotalSupply))
}
