package asset

import (
	"github.com/irisnet/irishub/modules/asset/internal/keeper"
	"github.com/irisnet/irishub/modules/asset/internal/types"
)

// nolint
type (
	MsgIssueToken           = types.MsgIssueToken
	MsgCreateGateway        = types.MsgCreateGateway
	MsgEditGateway          = types.MsgEditGateway
	MsgEditToken            = types.MsgEditToken
	MsgTransferGatewayOwner = types.MsgTransferGatewayOwner
	MsgMintToken            = types.MsgMintToken
	MsgTransferTokenOwner   = types.MsgTransferTokenOwner
	Tokens                  = types.Tokens
	Gateway                 = types.Gateway
	Gateways                = types.Gateways
	Params                  = types.Params
	FungibleToken           = types.FungibleToken
	AssetFamily             = types.AssetFamily
	AssetSource             = types.AssetSource
	QueryTokenParams        = types.QueryTokenParams
	QueryTokensParams       = types.QueryTokensParams
	QueryGatewayParams      = types.QueryGatewayParams
	QueryGatewaysParams     = types.QueryGatewaysParams
	QueryGatewayFeeParams   = types.QueryGatewayFeeParams
	QueryTokenFeesParams    = types.QueryTokenFeesParams
	GatewayFeeOutput        = types.GatewayFeeOutput
	TokenFeesOutput         = types.TokenFeesOutput
	GenesisState            = types.GenesisState

	Keeper = keeper.Keeper
)

// nolint
var (
	ModuleName   = types.ModuleName
	StoreKey     = types.StoreKey
	TStoreKey    = types.TStoreKey
	QuerierRoute = types.QuerierRoute
	RouterKey    = types.RouterKey
	ModuleCdc    = types.ModuleCdc

	NATIVE                 = types.NATIVE
	EXTERNAL               = types.EXTERNAL
	GATEWAY                = types.GATEWAY
	FUNGIBLE               = types.FUNGIBLE
	DefaultCodespace       = types.DefaultCodespace
	DefaultParamSpace      = types.DefaultParamSpace
	DoNotModify            = types.DoNotModify
	CodeInvalidAssetSource = types.CodeInvalidAssetSource
	MaximumAssetMaxSupply  = types.MaximumAssetMaxSupply
	RegisterCodec          = types.RegisterCodec
	ErrInvalidAssetFamily  = types.ErrInvalidAssetFamily
	ErrAssetAlreadyExists  = types.ErrAssetAlreadyExists
	CheckTokenID           = types.CheckTokenID
	ValidateMoniker        = types.ValidateMoniker
	StringToAssetFamilyMap = types.StringToAssetFamilyMap
	StringToAssetSourceMap = types.StringToAssetSourceMap
	GetTokenID             = types.GetTokenID
	ParseBool              = types.ParseBool

	NewFungibleToken           = types.NewFungibleToken
	NewMsgCreateGateway        = types.NewMsgCreateGateway
	NewGateway                 = types.NewGateway
	NewMsgEditGateway          = types.NewMsgEditGateway
	NewMsgEditToken            = types.NewMsgEditToken
	NewMsgTransferGatewayOwner = types.NewMsgTransferGatewayOwner
	NewMsgMintToken            = types.NewMsgMintToken
	NewMsgTransferTokenOwner   = types.NewMsgTransferTokenOwner
	NewMsgIssueToken           = types.NewMsgIssueToken
	DefaultParams              = types.DefaultParams

	QueryToken                  = types.QueryToken
	QueryTokens                 = types.QueryTokens
	QueryGateway                = types.QueryGateway
	QueryGateways               = types.QueryGateways
	QueryFees                   = types.QueryFees
	NewKeeper                   = keeper.NewKeeper
	TokenIssueFeeHandler        = keeper.TokenIssueFeeHandler
	GatewayTokenIssueFeeHandler = keeper.GatewayTokenIssueFeeHandler
	GatewayCreateFeeHandler     = keeper.GatewayCreateFeeHandler
	NewQuerier                  = keeper.NewQuerier
	NewAnteHandler              = keeper.NewAnteHandler
)
