package bank

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/tendermint/tendermint/crypto"

	"github.com/irisnet/irishub/modules/legacy/types"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type BaseAccount struct {
	Address       sdk.AccAddress `json:"address"`
	Coins         sdk.Coins      `json:"coins"`
	PubKey        crypto.PubKey  `json:"public_key"`
	AccountNumber uint64         `json:"account_number"`
	Sequence      uint64         `json:"sequence"`
	MemoRegexp    string         `json:"memo_regexp"`
}

// QueryAccountRequestHandlerFn performs account information query
func QueryAccountRequestHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)
		bech32addr := vars["address"]

		addr, err := sdk.AccAddressFromBech32(bech32addr)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		accGetter := authtypes.AccountRetriever{}

		account, height, err := accGetter.GetAccountWithHeight(clientCtx, addr)
		if err != nil {
			// TODO: Handle more appropriately based on the error type.
			// Ref: https://github.com/cosmos/cosmos-sdk/issues/4923
			if err := accGetter.EnsureExists(clientCtx, addr); err != nil {
				clientCtx = clientCtx.WithHeight(height)
				rest.PostProcessResponse(w, clientCtx, types.BaseAccount{})
				return
			}

			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		params := banktypes.NewQueryAllBalancesRequest(addr, nil)
		queryClient := banktypes.NewQueryClient(clientCtx)
		balances, err := queryClient.AllBalances(context.Background(), params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		legacyAccount := types.BaseAccount{
			Address:       account.GetAddress(),
			Coins:         balances.Balances,
			PubKey:        account.GetPubKey(),
			AccountNumber: account.GetAccountNumber(),
			Sequence:      account.GetSequence(),
		}

		rest.PostProcessResponseBare(w, clientCtx, types.Account(legacyAccount))
	}
}
