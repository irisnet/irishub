package asset

import (
	"github.com/irisnet/irishub/modules/asset/internal/keeper"
	"github.com/irisnet/irishub/modules/asset/internal/types"
)

type (
	Keeper = keeper.Keeper

	MsgIssueToken         = types.MsgIssueToken
	MsgEditToken          = types.MsgEditToken
	MsgMintToken          = types.MsgMintToken
	MsgTransferTokenOwner = types.MsgTransferTokenOwner
	Tokens                = types.Tokens
	Params                = types.Params
	FungibleToken         = types.FungibleToken
	QueryTokenParams      = types.QueryTokenParams
	QueryTokensParams     = types.QueryTokensParams
	QueryGatewayParams    = types.QueryGatewayParams
	QueryGatewaysParams   = types.QueryGatewaysParams
	QueryGatewayFeeParams = types.QueryGatewayFeeParams
	QueryTokenFeesParams  = types.QueryTokenFeesParams
	GatewayFeeOutput      = types.GatewayFeeOutput
	TokenFeesOutput       = types.TokenFeesOutput
	GenesisState          = types.GenesisState
)

const (
	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
	QuerierRoute      = types.QuerierRoute
	NATIVE            = types.NATIVE
	EXTERNAL          = types.EXTERNAL
	FUNGIBLE          = types.FUNGIBLE
	DefaultCodespace  = types.DefaultCodespace
	DefaultParamspace = types.DefaultParamspace
)

var (
	NewKeeper            = keeper.NewKeeper
	TokenIssueFeeHandler = keeper.TokenIssueFeeHandler
	NewQuerier           = keeper.NewQuerier

	ModuleCdc             = types.ModuleCdc
	RegisterCodec         = types.RegisterCodec
	ErrInvalidAssetFamily = types.ErrInvalidAssetFamily
	NewFungibleToken      = types.NewFungibleToken
	DefaultParams         = types.DefaultParams
	NewParams             = types.NewParams
	NewGenesisState       = types.NewGenesisState
)
