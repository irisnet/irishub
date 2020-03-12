package asset

import (
	"github.com/irisnet/irishub/app/v3/asset/internal/keeper"
	"github.com/irisnet/irishub/app/v3/asset/internal/types"
)

type (
	MsgIssueToken         = types.MsgIssueToken
	MsgEditToken          = types.MsgEditToken
	MsgMintToken          = types.MsgMintToken
	MsgTransferTokenOwner = types.MsgTransferTokenOwner
	Tokens                = types.Tokens
	TokenOutput           = types.TokenOutput
	TokensOutput          = types.TokensOutput
	Params                = types.Params
	FungibleToken         = types.FungibleToken
	QueryTokenParams      = types.QueryTokenParams
	QueryTokensParams     = types.QueryTokensParams
	QueryTokenFeesParams  = types.QueryTokenFeesParams
	TokenFeesOutput       = types.TokenFeesOutput
	GenesisState          = types.GenesisState

	Keeper = keeper.Keeper
)

var (
	DefaultCodespace      = types.DefaultCodespace
	DefaultParamSpace     = types.DefaultParamSpace
	MaximumAssetMaxSupply = types.MaximumAssetMaxSupply
	RegisterCodec         = types.RegisterCodec
	CheckSymbol           = types.CheckSymbol
	ParseBool             = types.ParseBool

	NewFungibleToken         = types.NewFungibleToken
	NewMsgEditToken          = types.NewMsgEditToken
	NewMsgMintToken          = types.NewMsgMintToken
	NewMsgTransferTokenOwner = types.NewMsgTransferTokenOwner
	DefaultParams            = types.DefaultParams
	DefaultParamsForTest     = types.DefaultParamsForTest
	ValidateParams           = types.ValidateParams

	QueryToken     = types.QueryToken
	QueryTokens    = types.QueryTokens
	QueryFees      = types.QueryFees
	NewKeeper      = keeper.NewKeeper
	NewQuerier     = keeper.NewQuerier
	NewAnteHandler = keeper.NewAnteHandler
)
