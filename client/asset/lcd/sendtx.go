package lcd

import (
	"fmt"
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
		fmt.Sprintf("/asset/tokens/{%s}", RestParamSymbol),
		editTokenHandlerFn(cdc, cliCtx),
	).Methods("PUT")

	// transfer owner
	r.HandleFunc(
		fmt.Sprintf("/asset/tokens/{%s}/transfer", RestParamSymbol),
		transferOwnerHandlerFn(cdc, cliCtx),
	).Methods("POST")

	// mint token
	r.HandleFunc(
		fmt.Sprintf("/asset/tokens/{%s}/mint", RestParamSymbol),
		mintTokenHandlerFn(cdc, cliCtx),
	).Methods("POST")
}

type issueTokenReq struct {
	BaseTx        utils.BaseTx   `json:"base_tx"`
	Owner         sdk.AccAddress `json:"owner"` // owner of the token
	Symbol        string         `json:"symbol"`
	Name          string         `json:"name"`
	Scale         uint8          `json:"scale"`
	MinUnit       string         `json:"min_unit"`
	InitialSupply uint64         `json:"initial_supply"`
	MaxSupply     uint64         `json:"max_supply"`
	Mintable      bool           `json:"mintable"`
}

type editTokenReq struct {
	BaseTx    utils.BaseTx   `json:"base_tx"`
	Owner     sdk.AccAddress `json:"owner"` //  owner of the token
	MaxSupply uint64         `json:"max_supply"`
	Mintable  string         `json:"mintable"` // mintable of the token
	Name      string         `json:"name"`
}

type transferTokenOwnerReq struct {
	BaseTx   utils.BaseTx   `json:"base_tx"`
	SrcOwner sdk.AccAddress `json:"src_owner"` // the current owner address of the token
	DstOwner sdk.AccAddress `json:"dst_owner"` // the new owner
}

type mintTokenReq struct {
	BaseTx utils.BaseTx   `json:"base_tx"`
	Owner  sdk.AccAddress `json:"owner"`  // the current owner address of the token
	To     sdk.AccAddress `json:"to"`     // address of minting token to
	Amount uint64         `json:"amount"` // amount of minting token
}

func issueTokenHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req issueTokenReq
		if err := utils.ReadPostBody(w, r, cdc, &req); err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the MsgIssueToken message
		msg := asset.MsgIssueToken{
			Symbol:        req.Symbol,
			Name:          req.Name,
			Decimal:       req.Scale,
			MinUnitAlias:  req.MinUnit,
			InitialSupply: req.InitialSupply,
			MaxSupply:     req.MaxSupply,
			Mintable:      req.Mintable,
			Owner:         req.Owner,
		}
		if err := msg.ValidateBasic(); err != nil {
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
		symbol := vars[RestParamSymbol]

		var req editTokenReq
		if err := utils.ReadPostBody(w, r, cdc, &req); err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		mintable, err := asset.ParseBool(req.Mintable)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the MsgEditToken message
		msg := asset.NewMsgEditToken(req.Name, symbol, req.MaxSupply, mintable, req.Owner)
		if err := msg.ValidateBasic(); err != nil {
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
		symbol := vars[RestParamSymbol]

		var req transferTokenOwnerReq
		if err := utils.ReadPostBody(w, r, cdc, &req); err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the MsgTransferTokenOwner message
		msg := asset.NewMsgTransferTokenOwner(req.SrcOwner, req.DstOwner, symbol)
		if err := msg.ValidateBasic(); err != nil {
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
		symbol := vars[RestParamSymbol]

		var req mintTokenReq
		if err := utils.ReadPostBody(w, r, cdc, &req); err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the MsgMintToken message
		msg := asset.NewMsgMintToken(symbol, req.Owner, req.To, req.Amount)
		if err := msg.ValidateBasic(); err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}
