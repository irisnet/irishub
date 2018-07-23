package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/cosmos/cosmos-sdk/client/context"
	cmn "github.com/tendermint/tmlibs/common"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"encoding/json"
)

type ExRateResponse struct {
	ExRate string `json:"token_shares_rate"`
}


func RegisterStakeExRate(ctx context.CoreContext, r *mux.Router, cdc *wire.Codec) {
	r.HandleFunc("/stake/validator/{valAddr}/exRate", GetValidatorExRate(ctx, cdc)).Methods("GET")
}

func GetValidatorExRate(ctx context.CoreContext, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		valAddr := vars["valAddr"]
		validatorAddr, err := sdk.GetAccAddressBech32(valAddr)

		// get validator
		validator, err := getValidator(validatorAddr, ctx, cdc)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		if validator.Owner == nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("validator not exist"))
			return
		}

		// get pool
		pool, err := getPool(ctx, cdc)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		// validator exRate
		valExRate := validator.DelegatorShareExRate(pool)

		// pool exRate
		poolExRate := bondedShareExRate(pool)

		exRate := poolExRate.Mul(valExRate)

		res := ExRateResponse{
			ExRate: exRate.String(),
		}

		resRaw, err := json.Marshal(res)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(resRaw)
	}
}

func getPool(ctx context.CoreContext, cdc *wire.Codec) (stake.Pool, error) {
	var (
		res []byte
		pool stake.Pool
	)
	res, err := query(ctx, stake.PoolKey)
	if err != nil {
		return pool, err
	}

	cdc.MustUnmarshalBinary(res, &pool)

	return pool, err
}

func getValidator(address sdk.Address, ctx context.CoreContext, cdc *wire.Codec) (stake.Validator, error)  {
	var (
		res []byte
		validator stake.Validator
	)
	res, err := query(ctx, stake.GetValidatorKey(address))
	if err != nil {
		return validator, err
	}

	cdc.MustUnmarshalBinary(res, &validator)

	return validator, err
}

// get the exchange rate of bonded token per issued share
func bondedShareExRate(p stake.Pool) sdk.Rat {
	if p.BondedShares.IsZero() {
		return sdk.OneRat()
	}
	return sdk.NewRat(p.BondedTokens).Quo(p.BondedShares)
}

func query(ctx context.CoreContext, key cmn.HexBytes) ([]byte, error) {
	res, err := ctx.Query(key, "stake")
	return res, err
}

