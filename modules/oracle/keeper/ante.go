package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"mods.irisnet.org/oracle/types"
)

type ValidateOracleAuthDecorator struct {
	k  Keeper
	ak types.AuthKeeper
}

func NewValidateOracleAuthDecorator(k Keeper, ak types.AuthKeeper) ValidateOracleAuthDecorator {
	return ValidateOracleAuthDecorator{
		k:  k,
		ak: ak,
	}
}

// AnteHandle returns an AnteHandler that checks if the creator is authorized
func (dtf ValidateOracleAuthDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	for _, msg := range tx.GetMsgs() {
		switch msg := msg.(type) {
		case *types.MsgCreateFeed:
			creator, err := sdk.AccAddressFromBech32(msg.Creator)
			if err != nil {
				return ctx, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid creator")
			}

			if !dtf.ak.Authorized(ctx, creator) {
				return ctx, errorsmod.Wrapf(types.ErrUnauthorized, msg.Creator)
			}
		}
	}

	// continue
	return next(ctx, tx, simulate)
}
