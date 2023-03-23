package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	v1 "github.com/irisnet/irismod/modules/token/types/v1"
	"github.com/irisnet/irismod/modules/token/types/v1beta1"
)

type legacyMsgServer struct {
	server v1.MsgServer
	k      Keeper
}

var _ v1beta1.MsgServer = legacyMsgServer{}

// NewLegacyMsgServerImpl returns an implementation of the token MsgServer interface
// for the provided Keeper.
func NewLegacyMsgServerImpl(server v1.MsgServer, k Keeper) v1beta1.MsgServer {
	return &legacyMsgServer{
		server: server,
		k:      k,
	}
}

func (m legacyMsgServer) IssueToken(goCtx context.Context, msg *v1beta1.MsgIssueToken) (*v1beta1.MsgIssueTokenResponse, error) {
	_, err := m.server.IssueToken(goCtx, &v1.MsgIssueToken{
		Symbol:        msg.Symbol,
		Name:          msg.Name,
		Scale:         msg.Scale,
		MinUnit:       msg.MinUnit,
		InitialSupply: msg.InitialSupply,
		MaxSupply:     msg.MaxSupply,
		Mintable:      msg.Mintable,
		Owner:         msg.Owner,
	})
	return &v1beta1.MsgIssueTokenResponse{}, err
}

func (m legacyMsgServer) EditToken(goCtx context.Context, msg *v1beta1.MsgEditToken) (*v1beta1.MsgEditTokenResponse, error) {
	_, err := m.server.EditToken(goCtx, &v1.MsgEditToken{
		Symbol:    msg.Symbol,
		Name:      msg.Name,
		MaxSupply: msg.MaxSupply,
		Mintable:  msg.Mintable,
		Owner:     msg.Owner,
	})
	return &v1beta1.MsgEditTokenResponse{}, err
}

func (m legacyMsgServer) MintToken(goCtx context.Context, msg *v1beta1.MsgMintToken) (*v1beta1.MsgMintTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	token, err := m.k.getTokenBySymbol(ctx, msg.Symbol)
	if err != nil {
		return &v1beta1.MsgMintTokenResponse{}, err
	}

	coin, err := token.ToMinCoin(sdk.NewDecCoin(msg.Symbol, sdkmath.NewIntFromUint64(msg.Amount)))
	if err != nil {
		return &v1beta1.MsgMintTokenResponse{}, err
	}
	_, err = m.server.MintToken(goCtx, &v1.MsgMintToken{
		Coin:  coin,
		To:    msg.To,
		Owner: msg.Owner,
	})
	return &v1beta1.MsgMintTokenResponse{}, err
}

func (m legacyMsgServer) BurnToken(goCtx context.Context, msg *v1beta1.MsgBurnToken) (*v1beta1.MsgBurnTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	token, err := m.k.getTokenBySymbol(ctx, msg.Symbol)
	if err != nil {
		return &v1beta1.MsgBurnTokenResponse{}, err
	}

	coin, err := token.ToMinCoin(sdk.NewDecCoin(msg.Symbol, sdkmath.NewIntFromUint64(msg.Amount)))
	if err != nil {
		return &v1beta1.MsgBurnTokenResponse{}, err
	}
	_, err = m.server.BurnToken(goCtx, &v1.MsgBurnToken{
		Coin:   coin,
		Sender: msg.Sender,
	})
	return &v1beta1.MsgBurnTokenResponse{}, err
}

func (m legacyMsgServer) TransferTokenOwner(goCtx context.Context, msg *v1beta1.MsgTransferTokenOwner) (*v1beta1.MsgTransferTokenOwnerResponse, error) {
	_, err := m.server.TransferTokenOwner(goCtx, &v1.MsgTransferTokenOwner{
		SrcOwner: msg.SrcOwner,
		DstOwner: msg.DstOwner,
		Symbol:   msg.Symbol,
	})
	return &v1beta1.MsgTransferTokenOwnerResponse{}, err
}
