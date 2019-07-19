package lcd

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/app/v1/asset"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	// create a gateway
	r.HandleFunc(
		"/asset/gateways",
		createGatewayHandlerFn(cdc, cliCtx),
	).Methods("POST")

	// edit a gateway
	r.HandleFunc(
		"/asset/gateways/{moniker}",
		editGatewayHandlerFn(cdc, cliCtx),
	).Methods("PUT")

	// transfer the ownership of a gateway
	r.HandleFunc(
		"/asset/gateways/{moniker}/transfer",
		transferGatewayOwnerHandlerFn(cdc, cliCtx),
	).Methods("POST")

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
		"/asset/tokens/{token-id}/transfer-owner",
		transferOwnerHandlerFn(cdc, cliCtx),
	).Methods("POST")

	// mint token
	r.HandleFunc(
		"/asset/tokens/{token-id}/mint",
		mintTokenHandlerFn(cdc, cliCtx),
	).Methods("POST")
}

type issueTokenReq struct {
	BaseTx         utils.BaseTx      `json:"base_tx"`
	Owner          sdk.AccAddress    `json:"owner"` //  Owner of the token
	Family         asset.AssetFamily `json:"family"`
	Source         asset.AssetSource `json:"source"`
	Gateway        string            `json:"gateway"`
	Symbol         string            `json:"symbol"`
	SymbolAtSource string            `json:"symbol_at_source"`
	Name           string            `json:"name"`
	Decimal        uint8             `json:"decimal"`
	SymbolMinAlias string            `json:"symbol_min_alias"`
	InitialSupply  uint64            `json:"initial_supply"`
	MaxSupply      uint64            `json:"max_supply"`
	Mintable       bool              `json:"mintable"`
}

type createGatewayReq struct {
	BaseTx   utils.BaseTx   `json:"base_tx"`
	Owner    sdk.AccAddress `json:"owner"`    //  Owner of the gateway
	Moniker  string         `json:"moniker"`  //  Name of the gateway
	Identity string         `json:"identity"` //  Identity of the gateway
	Details  string         `json:"details"`  //  Description of the gateway
	Website  string         `json:"website"`  //  Website of the gateway
}

type editGatewayReq struct {
	BaseTx   utils.BaseTx   `json:"base_tx"`
	Owner    sdk.AccAddress `json:"owner"`    //  Owner of the gateway
	Identity string         `json:"identity"` //  Identity of the gateway
	Details  string         `json:"details"`  //  Description of the gateway
	Website  string         `json:"website"`  //  Website of the gateway
}

type transferGatewayOwnerReq struct {
	BaseTx utils.BaseTx   `json:"base_tx"`
	Owner  sdk.AccAddress `json:"owner"` // Current Owner of the gateway
	To     sdk.AccAddress `json:"to"`    // New owner of the gateway
}

type editTokenReq struct {
	BaseTx         utils.BaseTx   `json:"base_tx"`
	Owner          sdk.AccAddress `json:"owner"`            //  owner of asset
	SymbolAtSource string         `json:"symbol_at_source"` //  symbol_at_source of asset
	SymbolMinAlias string         `json:"symbol_min_alias"` //  symbol_min_alias of asset
	MaxSupply      uint64         `json:"max_supply"`
	Mintable       *bool          `json:"mintable"` //  mintable of asset
	Name           string         `json:"name"`
}

type transferTokenOwnerReq struct {
	BaseTx   utils.BaseTx   `json:"base_tx"`
	SrcOwner sdk.AccAddress `json:"src_owner"` // the current owner address of the token
	DstOwner sdk.AccAddress `json:"dst_owner"` // the new owner
}

type mintTokenReq struct {
	BaseTx utils.BaseTx   `json:"base_tx"`
	Owner  sdk.AccAddress `json:"owner"`  // the current owner address of the token
	To     sdk.AccAddress `json:"to"`     // address of mint token to
	Amount uint64         `json:"amount"` // amount of mint token
}

func createGatewayHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createGatewayReq
		err := utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the MsgCreateGateway message
		msg := asset.NewMsgCreateGateway(req.Owner, req.Moniker, req.Identity, req.Details, req.Website)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}

func editGatewayHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		moniker := vars["moniker"]

		var req editGatewayReq
		err := utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the MsgEditGateway message
		msg := asset.NewMsgEditGateway(req.Owner, moniker, req.Identity, req.Details, req.Website)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}

func transferGatewayOwnerHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		moniker := vars["moniker"]

		var req transferGatewayOwnerReq
		err := utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the MsgTransferGatewayOwner message
		msg := asset.NewMsgTransferGatewayOwner(req.Owner, moniker, req.To)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
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
		msg := asset.NewMsgIssueToken(req.Family, req.Source, req.Gateway, req.Symbol, req.SymbolAtSource, req.Name, req.Decimal, req.SymbolMinAlias, req.InitialSupply, req.MaxSupply, req.Mintable, req.Owner)
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
		msg := asset.NewMsgEditToken(req.Name, req.SymbolAtSource, req.SymbolMinAlias, tokenId, req.MaxSupply, req.Mintable, req.Owner)
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
		msg := asset.NewMsgMintToken(tokenId, req.Owner, req.To, req.Amount)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}
