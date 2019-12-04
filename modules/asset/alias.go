package asset

import (
	"github.com/irisnet/irishub/modules/asset/keeper"
	"github.com/irisnet/irishub/modules/asset/types"
)

type (
	MsgIssueToken         = types.MsgIssueToken
	MsgEditToken          = types.MsgEditToken
	MsgMintToken          = types.MsgMintToken
	MsgTransferTokenOwner = types.MsgTransferTokenOwner
	Tokens                = types.Tokens
	Params                = types.Params
	FungibleToken         = types.FungibleToken
	AssetFamily           = types.AssetFamily
	AssetSource           = types.AssetSource
	QueryTokenParams      = types.QueryTokenParams
	QueryTokensParams     = types.QueryTokensParams
	QueryGatewayParams    = types.QueryGatewayParams
	QueryGatewaysParams   = types.QueryGatewaysParams
	QueryGatewayFeeParams = types.QueryGatewayFeeParams
	QueryTokenFeesParams  = types.QueryTokenFeesParams
	GatewayFeeOutput      = types.GatewayFeeOutput
	TokenFeesOutput       = types.TokenFeesOutput
	GenesisState          = types.GenesisState

	Keeper = keeper.Keeper
)

var (
	ModuleName             = types.ModuleName
	ModuleCdc              = types.ModuleCdc
	RouterKey              = types.RouterKey
	StoreKey               = types.StoreKey
	QuerierRoute           = types.QuerierRoute
	NATIVE                 = types.NATIVE
	EXTERNAL               = types.EXTERNAL
	FUNGIBLE               = types.FUNGIBLE
	DefaultCodespace       = types.DefaultCodespace
	DefaultParamspace      = types.DefaultParamspace
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

	NewFungibleToken         = types.NewFungibleToken
	NewMsgEditToken          = types.NewMsgEditToken
	NewMsgMintToken          = types.NewMsgMintToken
	NewMsgTransferTokenOwner = types.NewMsgTransferTokenOwner
	NewMsgIssueToken         = types.NewMsgIssueToken
	DefaultParams            = types.DefaultParams

	QueryToken           = types.QueryToken
	QueryTokens          = types.QueryTokens
	QueryGateway         = types.QueryGateway
	QueryGateways        = types.QueryGateways
	QueryFees            = types.QueryFees
	NewKeeper            = keeper.NewKeeper
	TokenIssueFeeHandler = keeper.TokenIssueFeeHandler
	NewQuerier           = keeper.NewQuerier
)
