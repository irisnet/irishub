package rest

import (
	"net/http"
	
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/irisnet/irishub/modules/asset"
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
	BaseReq          utils.BaseReq      `json:"base_req"`
	Owner           sdk.AccAddress    `json:"owner"` //  Owner of the token
	Family          asset.AssetFamily `json:"family"`
	Source          asset.AssetSource `json:"source"`
	Gateway         string            `json:"gateway"`
	Symbol          string            `json:"symbol"`
	CanonicalSymbol string            `json:"canonical_symbol"`
	Name            string            `json:"name"`
	Decimal         uint8             `json:"decimal"`
	MinUnitAlias    string            `json:"min_unit_alias"`
	InitialSupply   uint64            `json:"initial_supply"`
	MaxSupply       uint64            `json:"max_supply"`
	Mintable        bool              `json:"mintable"`
}

type createGatewayReq struct {
	BaseReq   utils.BaseReq   `json:"base_req"`
	Owner    sdk.AccAddress `json:"owner"`    //  Owner of the gateway
	Moniker  string         `json:"moniker"`  //  Name of the gateway
	Identity string         `json:"identity"` //  Identity of the gateway
	Details  string         `json:"details"`  //  Description of the gateway
	Website  string         `json:"website"`  //  Website of the gateway
}

type editGatewayReq struct {
	BaseReq   utils.BaseReq   `json:"base_req"`
	Owner    sdk.AccAddress `json:"owner"`    //  Owner of the gateway
	Identity string         `json:"identity"` //  Identity of the gateway
	Details  string         `json:"details"`  //  Description of the gateway
	Website  string         `json:"website"`  //  Website of the gateway
}

type transferGatewayOwnerReq struct {
	BaseReq utils.BaseReq  `json:"base_req"`
	Owner  sdk.AccAddress `json:"owner"` // Current Owner of the gateway
	To     sdk.AccAddress `json:"to"`    // New owner of the gateway
}

type editTokenReq struct {
	BaseReq          utils.BaseReq   `json:"base_req"`
	Owner           sdk.AccAddress `json:"owner"`            //  owner of asset
	CanonicalSymbol string         `json:"canonical_symbol"` //  canonical_symbol of asset
	MinUnitAlias    string         `json:"min_unit_alias"`   //  min_unit_alias of asset
	MaxSupply       uint64         `json:"max_supply"`
	Mintable        string         `json:"mintable"` //  mintable of asset
	Name            string         `json:"name"`
}

type transferTokenOwnerReq struct {
	BaseReq   utils.BaseReq   `json:"base_req"`
	SrcOwner sdk.AccAddress `json:"src_owner"` // the current owner address of the token
	DstOwner sdk.AccAddress `json:"dst_owner"` // the new owner
}

type mintTokenReq struct {
	BaseReq utils.BaseReq   `json:"base_req"`
	Owner  sdk.AccAddress `json:"owner"`  // the current owner address of the token
	To     sdk.AccAddress `json:"to"`     // address of mint token to
	Amount uint64         `json:"amount"` // amount of mint token
}

func createGatewayHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createGatewayReq
		if !rest.ReadRESTBody(w, r, cliCtx.Codec, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
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

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq,[]sdk.Msg{msg})
	}
}

func editGatewayHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		moniker := vars["moniker"]

		var req editGatewayReq
		if !rest.ReadRESTBody(w, r, cliCtx.Codec, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
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

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func transferGatewayOwnerHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		moniker := vars["moniker"]

		var req transferGatewayOwnerReq
		if !rest.ReadRESTBody(w, r, cliCtx.Codec, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
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

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq,[]sdk.Msg{msg})
	}
}

func issueTokenHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req issueTokenReq
		if !rest.ReadRESTBody(w, r, cliCtx.Codec, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the MsgEditGateway message
		msg := asset.NewMsgIssueToken(req.Family, req.Source, req.Gateway, req.Symbol, req.CanonicalSymbol, req.Name, req.Decimal, req.MinUnitAlias, req.InitialSupply, req.MaxSupply, req.Mintable, req.Owner)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func editTokenHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tokenId := vars["token-id"]

		var req editTokenReq
		if !rest.ReadRESTBody(w, r, cliCtx.Codec, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the MsgEditToken message
		mintable, err := asset.ParseBool(req.Mintable)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		msg := asset.NewMsgEditToken(req.Name, req.CanonicalSymbol, req.MinUnitAlias, tokenId, req.MaxSupply, mintable, req.Owner)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func transferOwnerHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tokenId := vars["token-id"]
		var req transferTokenOwnerReq
		if !rest.ReadRESTBody(w, r, cliCtx.Codec, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
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

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func mintTokenHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tokenId := vars["token-id"]
		
		var req mintTokenReq
		if !rest.ReadRESTBody(w, r, cliCtx.Codec, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
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

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}
