package asset

import (
	"github.com/irisnet/irishub/app/v1/asset/internal/keeper"
	"github.com/irisnet/irishub/app/v1/asset/internal/types"
)

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

	Keeper                = keeper.Keeper
	GatewayFeeOutput      = keeper.GatewayFeeOutput
	TokenFeesOutput       = keeper.TokenFeesOutput
	QueryTokenParams      = keeper.QueryTokenParams
	QueryTokensParams     = keeper.QueryTokensParams
	QueryGatewayParams    = keeper.QueryGatewayParams
	QueryGatewaysParams   = keeper.QueryGatewaysParams
	QueryGatewayFeeParams = keeper.QueryGatewayFeeParams
	QueryTokenFeesParams  = keeper.QueryTokenFeesParams
)

var (
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
	DefaultParamsForTest       = types.DefaultParamsForTest
	ValidateParams             = types.ValidateParams

	QueryToken                  = keeper.QueryToken
	QueryTokens                 = keeper.QueryTokens
	QueryGateway                = keeper.QueryGateway
	QueryGateways               = keeper.QueryGateways
	QueryFees                   = keeper.QueryFees
	NewKeeper                   = keeper.NewKeeper
	TokenIssueFeeHandler        = keeper.TokenIssueFeeHandler
	GatewayTokenIssueFeeHandler = keeper.GatewayTokenIssueFeeHandler
	GatewayCreateFeeHandler     = keeper.GatewayCreateFeeHandler
	NewQuerier                  = keeper.NewQuerier
	NewAnteHandler              = keeper.NewAnteHandler
)
