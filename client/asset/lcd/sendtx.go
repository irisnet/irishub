package lcd

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/app/v3/asset"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	// issue a token
	r.HandleFunc(
		"/asset/tokens",
		issueTokenHandlerFn(cdc, cliCtx),
	).Methods("POST")

	// edit a token
	r.HandleFunc(
		"/asset/tokens/{token-id}",
		editTokenHandlerFn(cdc, cliCtx),
	).Methods("PUT")

	// transfer owner
	r.HandleFunc(
		"/asset/tokens/{token-id}/transfer",
		transferOwnerHandlerFn(cdc, cliCtx),
	).Methods("POST")

	// mint token
	r.HandleFunc(
		"/asset/tokens/{token-id}/mint",
		mintTokenHandlerFn(cdc, cliCtx),
	).Methods("POST")
}

type issueTokenReq struct {
	BaseTx        utils.BaseTx   `json:"base_tx"`
	Owner         sdk.AccAddress `json:"owner"` //  Owner of the token
	Symbol        string         `json:"symbol"`
	Name          string         `json:"name"`
	Scale         uint8          `json:"scale"`
	InitialSupply uint64         `json:"initial_supply"`
	MaxSupply     uint64         `json:"max_supply"`
	Mintable      bool           `json:"mintable"`
}

type editTokenReq struct {
	BaseTx    utils.BaseTx   `json:"base_tx"`
	Owner     sdk.AccAddress `json:"owner"` //  owner of asset
	MaxSupply uint64         `json:"max_supply"`
	Mintable  string         `json:"mintable"` //  mintable of asset
	Name      string         `json:"name"`
}

type transferTokenOwnerReq struct {
	BaseTx   utils.BaseTx   `json:"base_tx"`
	SrcOwner sdk.AccAddress `json:"src_owner"` // the current owner address of the token
	DstOwner sdk.AccAddress `json:"dst_owner"` // the new owner
}

type mintTokenReq struct {
	BaseTx    utils.BaseTx   `json:"base_tx"`
	Owner     sdk.AccAddress `json:"owner"`     // the current owner address of the token
	Recipient sdk.AccAddress `json:"recipient"` // address of mint token to
	Amount    uint64         `json:"amount"`    // amount of mint token
}

func issueTokenHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req issueTokenReq
		err := utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the MsgEditGateway message
		msg := asset.MsgIssueToken{
			Symbol:        req.Symbol,
			Name:          req.Name,
			Decimal:       req.Scale,
			InitialSupply: req.InitialSupply,
			MaxSupply:     req.MaxSupply,
			Mintable:      req.Mintable,
			Owner:         req.Owner,
		}
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}

func editTokenHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tokenId := vars["token-id"]

		var req editTokenReq
		err := utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the MsgEditToken message
		mintable, err := asset.ParseBool(req.Mintable)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		msg := asset.NewMsgEditToken(req.Name, tokenId, req.MaxSupply, mintable, req.Owner)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}

func transferOwnerHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tokenId := vars["token-id"]
		var req transferTokenOwnerReq
		err := utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the MsgTransferTokenOwner message
		msg := asset.NewMsgTransferTokenOwner(req.SrcOwner, req.DstOwner, tokenId)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}

func mintTokenHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tokenId := vars["token-id"]
		var req mintTokenReq
		err := utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the MsgMintToken message
		msg := asset.NewMsgMintToken(tokenId, req.Owner, req.Recipient, req.Amount)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}
