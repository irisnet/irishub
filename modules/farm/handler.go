package farm

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"mods.irisnet.org/modules/farm/keeper"
	"mods.irisnet.org/modules/farm/types"
)

// NewProposalHandler returns a handler function for processing farm proposals.
//
// It takes in a context and a content interface and returns an error.
func NewProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.CommunityPoolCreateFarmProposal:
			return handleCreateFarmProposal(ctx, k, c)
		default:
			return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized farm proposal content type: %T", c)
		}
	}
}


func handleCreateFarmProposal(ctx sdk.Context, k keeper.Keeper, p *types.CommunityPoolCreateFarmProposal) error {
	return k.HandleCreateFarmProposal(ctx, p)
}