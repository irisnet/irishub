package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/modules/asset/internal/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	// issue a token
	r.HandleFunc(
		"/asset/tokens",
		issueTokenHandlerFn(cliCtx),
	).Methods("POST")

	// edit a token
	r.HandleFunc(
		"/asset/tokens/{token-id}",
		editTokenHandlerFn(cliCtx),
	).Methods("PUT")

	// transfer owner
	r.HandleFunc(
		"/asset/tokens/{token-id}/transfer-owner",
		transferOwnerHandlerFn(cliCtx),
	).Methods("POST")

	// mint token
	r.HandleFunc(
		"/asset/tokens/{token-id}/mint",
		mintTokenHandlerFn(cliCtx),
	).Methods("POST")
}

func issueTokenHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req issueTokenReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the MsgEditGateway message
		msg := types.NewMsgIssueToken(req.Family, req.Source, req.Symbol, req.CanonicalSymbol, req.Name, req.Decimal,
			req.MinUnitAlias, req.InitialSupply, req.MaxSupply, req.Mintable, req.Owner)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseTx, []sdk.Msg{msg})
	}
}

func editTokenHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tokenId := vars["token-id"]

		var req editTokenReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the MsgEditToken message
		mintable, err := types.ParseBool(req.Mintable)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		msg := types.NewMsgEditToken(req.Name, req.CanonicalSymbol, req.MinUnitAlias, tokenId, req.MaxSupply, mintable, req.Owner)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseTx, []sdk.Msg{msg})
	}
}

func transferOwnerHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tokenId := vars["token-id"]
		var req transferTokenOwnerReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the MsgTransferTokenOwner message
		msg := types.NewMsgTransferTokenOwner(req.SrcOwner, req.DstOwner, tokenId)

		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseTx, []sdk.Msg{msg})
	}
}

func mintTokenHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tokenId := vars["token-id"]
		var req mintTokenReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the MsgMintToken message
		msg := types.NewMsgMintToken(tokenId, req.Owner, req.To, req.Amount)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseTx, []sdk.Msg{msg})
	}
}
