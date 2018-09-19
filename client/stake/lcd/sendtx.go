package lcd

import (
	"bytes"
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/cosmos/cosmos-sdk/x/stake/types"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"net/http"
)

type msgDelegationsInput struct {
	ValidatorAddr string   `json:"validator_addr"` // in bech32
	Delegation    sdk.Coin `json:"delegation"`
}
type msgBeginRedelegateInput struct {
	ValidatorSrcAddr string `json:"validator_src_addr"` // in bech32
	ValidatorDstAddr string `json:"validator_dst_addr"` // in bech32
	SharesAmount     string `json:"shares"`
}
type msgCompleteRedelegateInput struct {
	ValidatorSrcAddr string `json:"validator_src_addr"` // in bech32
	ValidatorDstAddr string `json:"validator_dst_addr"` // in bech32
}
type msgBeginUnbondingInput struct {
	ValidatorAddr string `json:"validator_addr"` // in bech32
	SharesAmount  string `json:"shares"`
}
type msgCompleteUnbondingInput struct {
	ValidatorAddr string `json:"validator_addr"` // in bech32
}

// the request body for edit delegations
type EditDelegationsBody struct {
	BaseTx              context.BaseTx               `json:"base_tx"`
	Delegations         []msgDelegationsInput        `json:"delegations"`
	BeginUnbondings     []msgBeginUnbondingInput     `json:"begin_unbondings"`
	CompleteUnbondings  []msgCompleteUnbondingInput  `json:"complete_unbondings"`
	BeginRedelegates    []msgBeginRedelegateInput    `json:"begin_redelegates"`
	CompleteRedelegates []msgCompleteRedelegateInput `json:"complete_redelegates"`
}

// nolint: gocyclo
// TODO: Split this up into several smaller functions, and remove the above nolint
// TODO: use sdk.ValAddress instead of sdk.AccAddress for validators in messages
func delegationsRequestHandlerFn(cdc *wire.Codec, kb keys.Keybase, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		delegatorAddr := vars["delegatorAddr"]

		var m EditDelegationsBody
		err := utils.ReadPostBody(w, r, cdc, &m)
		if err != nil {
			return
		}
		cliCtx = utils.InitRequestClictx(cliCtx, r, m.BaseTx.LocalAccountName, delegatorAddr)
		txCtx, err := context.NewTxContextFromBaseTx(cliCtx, cdc, m.BaseTx)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		delegatorAccAddress, err := sdk.AccAddressFromBech32(delegatorAddr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Couldn't decode delegator. Error: %s", err.Error())))
			return
		}
		if !cliCtx.GenerateOnly {
			fromAddress, err := cliCtx.GetFromAddress()
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			if !bytes.Equal(fromAddress, delegatorAccAddress) {
				utils.WriteErrorResponse(w, http.StatusUnauthorized, "Must use own delegator address")
				return
			}
		}
		// build messages
		messages := make([]sdk.Msg, len(m.Delegations)+
			len(m.BeginRedelegates)+
			len(m.CompleteRedelegates)+
			len(m.BeginUnbondings)+
			len(m.CompleteUnbondings))

		i := 0
		for _, msg := range m.Delegations {

			validatorAddr, err := sdk.AccAddressFromBech32(msg.ValidatorAddr)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Couldn't decode validator. Error: %s", err.Error())))
				return
			}

			delegationToken, err := cliCtx.ParseCoin(msg.Delegation.String())
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			messages[i] = stake.MsgDelegate{
				DelegatorAddr: delegatorAccAddress,
				ValidatorAddr: validatorAddr,
				Delegation:    delegationToken,
			}

			i++
		}

		for _, msg := range m.BeginRedelegates {
			validatorSrcAddr, err := sdk.AccAddressFromBech32(msg.ValidatorSrcAddr)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Couldn't decode validator. Error: %s", err.Error())))
				return
			}
			validatorDstAddr, err := sdk.AccAddressFromBech32(msg.ValidatorDstAddr)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Couldn't decode validator. Error: %s", err.Error())))
				return
			}

			shares, err := sdk.NewRatFromDecimal(msg.SharesAmount, types.MaxBondDenominatorPrecision)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Couldn't decode shares amount. Error: %s", err.Error())))
				return
			}

			messages[i] = stake.MsgBeginRedelegate{
				DelegatorAddr:    delegatorAccAddress,
				ValidatorSrcAddr: validatorSrcAddr,
				ValidatorDstAddr: validatorDstAddr,
				SharesAmount:     shares,
			}

			i++
		}

		for _, msg := range m.CompleteRedelegates {
			validatorSrcAddr, err := sdk.AccAddressFromBech32(msg.ValidatorSrcAddr)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Couldn't decode validator. Error: %s", err.Error())))
				return
			}
			validatorDstAddr, err := sdk.AccAddressFromBech32(msg.ValidatorDstAddr)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Couldn't decode validator. Error: %s", err.Error())))
				return
			}

			messages[i] = stake.MsgCompleteRedelegate{
				DelegatorAddr:    delegatorAccAddress,
				ValidatorSrcAddr: validatorSrcAddr,
				ValidatorDstAddr: validatorDstAddr,
			}

			i++
		}

		for _, msg := range m.BeginUnbondings {
			validatorAddr, err := sdk.AccAddressFromBech32(msg.ValidatorAddr)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Couldn't decode validator. Error: %s", err.Error())))
				return
			}

			shares, err := sdk.NewRatFromDecimal(msg.SharesAmount, types.MaxBondDenominatorPrecision)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Couldn't decode shares amount. Error: %s", err.Error())))
				return
			}

			messages[i] = stake.MsgBeginUnbonding{
				DelegatorAddr: delegatorAccAddress,
				ValidatorAddr: validatorAddr,
				SharesAmount:  shares,
			}

			i++
		}

		for _, msg := range m.CompleteUnbondings {
			validatorAddr, err := sdk.AccAddressFromBech32(msg.ValidatorAddr)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Couldn't decode validator. Error: %s", err.Error())))
				return
			}

			messages[i] = stake.MsgCompleteUnbonding{
				DelegatorAddr: delegatorAccAddress,
				ValidatorAddr: validatorAddr,
			}
			i++
		}

		utils.SendOrReturnUnsignedTx(w, cliCtx, txCtx, m.BaseTx, messages)
	}
}
