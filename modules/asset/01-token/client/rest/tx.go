package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/irisnet/irishub/modules/asset/01-token/internal/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	// issue a token
	r.HandleFunc(
		"/asset/tokens",
		issueTokenHandlerFn(cliCtx),
	).Methods("POST")

	// edit a token
	r.HandleFunc(
		fmt.Sprintf("/asset/tokens/{%s}", RestParamSymbol),
		editTokenHandlerFn(cliCtx),
	).Methods("PUT")

	// transfer owner
	r.HandleFunc(
		fmt.Sprintf("/asset/tokens/{%s}/transfer", RestParamSymbol),
		transferTokenHandlerFn(cliCtx),
	).Methods("POST")

	// mint token
	r.HandleFunc(
		fmt.Sprintf("/asset/tokens/{%s}/mint", RestParamSymbol),
		mintTokenHandlerFn(cliCtx),
	).Methods("POST")

	// burn token
	r.HandleFunc(
		"/asset/tokens/burn",
		burnTokenHandlerFn(cliCtx),
	).Methods("POST")
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
			req.Symbol, req.Name, req.Scale,
			req.MinUnit, req.InitialSupply, req.MaxSupply, req.Mintable, req.Owner,
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
		symbol := vars[RestParamSymbol]

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
			req.Name, symbol, req.MaxSupply, mintable, req.Owner,
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseTx, []sdk.Msg{msg})
	}
}

// HTTP request handler to transfer owner.
func transferTokenHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		symbol := vars[RestParamSymbol]
		var req TransferTokenReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the MsgTransferToken message
		msg := types.NewMsgTransferToken(req.SrcOwner, req.DstOwner, symbol)
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
		symbol := vars[RestParamSymbol]
		var req MintTokenReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the MsgMintToken message
		msg := types.NewMsgMintToken(symbol, req.Owner, req.Recipient, req.Amount)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseTx, []sdk.Msg{msg})
	}
}

func burnTokenHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req BurnTokenReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the MsgBurnToken message
		msg := types.NewMsgBurnToken(req.Sender, req.Amount)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseTx, []sdk.Msg{msg})
	}
}
