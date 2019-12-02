package rest

import (
	"github.com/irisnet/irishub/modules/asset/types"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
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

type issueTokenReq struct {
	BaseTx          rest.BaseReq      `json:"base_tx"`
	Owner           sdk.AccAddress    `json:"owner"` //  Owner of the token
	Family          types.AssetFamily `json:"family"`
	Source          types.AssetSource `json:"source"`
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
	BaseTx   rest.BaseReq   `json:"base_tx"`
	Owner    sdk.AccAddress `json:"owner"`    //  Owner of the gateway
	Moniker  string         `json:"moniker"`  //  Name of the gateway
	Identity string         `json:"identity"` //  Identity of the gateway
	Details  string         `json:"details"`  //  Description of the gateway
	Website  string         `json:"website"`  //  Website of the gateway
}

type editGatewayReq struct {
	BaseTx   rest.BaseReq   `json:"base_tx"`
	Owner    sdk.AccAddress `json:"owner"`    //  Owner of the gateway
	Identity string         `json:"identity"` //  Identity of the gateway
	Details  string         `json:"details"`  //  Description of the gateway
	Website  string         `json:"website"`  //  Website of the gateway
}

type transferGatewayOwnerReq struct {
	BaseTx rest.BaseReq   `json:"base_tx"`
	Owner  sdk.AccAddress `json:"owner"` // Current Owner of the gateway
	To     sdk.AccAddress `json:"to"`    // New owner of the gateway
}

type editTokenReq struct {
	BaseTx          rest.BaseReq   `json:"base_tx"`
	Owner           sdk.AccAddress `json:"owner"`            //  owner of asset
	CanonicalSymbol string         `json:"canonical_symbol"` //  canonical_symbol of asset
	MinUnitAlias    string         `json:"min_unit_alias"`   //  min_unit_alias of asset
	MaxSupply       uint64         `json:"max_supply"`
	Mintable        string         `json:"mintable"` //  mintable of asset
	Name            string         `json:"name"`
}

type transferTokenOwnerReq struct {
	BaseTx   rest.BaseReq   `json:"base_tx"`
	SrcOwner sdk.AccAddress `json:"src_owner"` // the current owner address of the token
	DstOwner sdk.AccAddress `json:"dst_owner"` // the new owner
}

type mintTokenReq struct {
	BaseTx rest.BaseReq   `json:"base_tx"`
	Owner  sdk.AccAddress `json:"owner"`  // the current owner address of the token
	To     sdk.AccAddress `json:"to"`     // address of mint token to
	Amount uint64         `json:"amount"` // amount of mint token
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
		msg := types.NewMsgIssueToken(req.Family, req.Source, req.Gateway, req.Symbol, req.CanonicalSymbol, req.Name, req.Decimal, req.MinUnitAlias, req.InitialSupply, req.MaxSupply, req.Mintable, req.Owner)
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
