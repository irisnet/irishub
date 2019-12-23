package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/irisnet/irishub/modules/asset/internal/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	// issue a token
	r.HandleFunc("/asset/tokens", issueTokenHandlerFn(cliCtx)).Methods("POST")
	// edit a token
	r.HandleFunc(fmt.Sprintf("/asset/tokens/{%s}", RestTokenID), editTokenHandlerFn(cliCtx)).Methods("PUT")
	// transfer owner
	r.HandleFunc(fmt.Sprintf("/asset/tokens/{%s}/transfer-owner", RestTokenID), transferOwnerHandlerFn(cliCtx)).Methods("POST")
	// mint token
	r.HandleFunc(fmt.Sprintf("/asset/tokens/{%s}/mint", RestTokenID), mintTokenHandlerFn(cliCtx)).Methods("POST")
}

// HTTP request handler to issue token.
func issueTokenHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req IssueTokenReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the NewMsgIssueToken message
		msg := types.NewMsgIssueToken(
			req.Family, req.Source, req.Symbol, req.CanonicalSymbol, req.Name, req.Decimal,
			req.MinUnitAlias, req.InitialSupply, req.MaxSupply, req.Mintable, req.Owner,
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseTx, []sdk.Msg{msg})
	}
}

// HTTP request handler to edit token.
func editTokenHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tokenID := vars[RestTokenID]

		var req EditTokenReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		mintable, err := types.ParseBool(req.Mintable)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the MsgEditToken message
		msg := types.NewMsgEditToken(
			req.Name, req.CanonicalSymbol, req.MinUnitAlias, tokenID, req.MaxSupply, mintable, req.Owner,
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseTx, []sdk.Msg{msg})
	}
}

// HTTP request handler to transfer owner.
func transferOwnerHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tokenID := vars[RestTokenID]
		var req TransferTokenOwnerReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the MsgTransferTokenOwner message
		msg := types.NewMsgTransferTokenOwner(req.SrcOwner, req.DstOwner, tokenID)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseTx, []sdk.Msg{msg})
	}
}

// HTTP request handler to mint token.
func mintTokenHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tokenID := vars[RestTokenID]
		var req MintTokenReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the MsgMintToken message
		msg := types.NewMsgMintToken(tokenID, req.Owner, req.To, req.Amount)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseTx, []sdk.Msg{msg})
	}
}
