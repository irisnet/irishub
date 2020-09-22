package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/oracle/types"
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
		// only check consecutive msgs which are routed to token from the beginning
		if msg.Route() != types.ModuleName {
			break
		}

		switch msg := msg.(type) {
		case *types.MsgCreateFeed:
			if !dtf.ak.Authorized(ctx, msg.Creator) {
				return ctx, sdkerrors.Wrapf(types.ErrUnauthorized, msg.Creator.String())
			}
		}
	}

	// continue
	return next(ctx, tx, simulate)
}
