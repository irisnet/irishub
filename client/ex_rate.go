package client

import (
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/cosmos/cosmos-sdk/x/stake/types"
)

type ExRateResponse struct {
	ExRate float64 `json:"token_shares_rate"`
}

func RegisterStakeExRate(ctx context.CLIContext, r *mux.Router, cdc *wire.Codec) {
	r.HandleFunc("/stake/validator/{valAddr}/exRate", GetValidatorExRate(ctx, cdc)).Methods("GET")
}

func GetValidatorExRate(ctx context.CLIContext, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		valAddr := vars["valAddr"]
		validatorAddr, err := sdk.AccAddressFromBech32(valAddr)

		// get validator
		validator, err := getValidator(validatorAddr, ctx, cdc)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		if validator.Owner == nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("validator not exist"))
			return
		}

		// validator exRate
		valExRate := validator.DelegatorShareExRate()

		floatExRate, _ := valExRate.Float64()
		res := ExRateResponse{
			ExRate: floatExRate,
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

func getValidator(address sdk.AccAddress, ctx context.CLIContext, cdc *wire.Codec) (stake.Validator, error)  {
	var (
		res []byte
		validator stake.Validator
	)
	res, err := query(ctx, stake.GetValidatorKey(address))
	if err != nil {
		return validator, err
	}

	validator = types.MustUnmarshalValidator(cdc, address, res)

	return validator, err
}

func query(ctx context.CLIContext, key cmn.HexBytes) ([]byte, error) {
	res, err := ctx.QueryStore(key, "stake")
	return res, err
}