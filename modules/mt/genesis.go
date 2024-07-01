package mt

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"mods.irisnet.org/modules/mt/keeper"
	"mods.irisnet.org/modules/mt/types"
)

// InitGenesis stores the MT genesis.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	if err := types.ValidateGenesis(data); err != nil {
		panic(err.Error())
	}

	// ---------- init infos ---------- //

	// set denom sequence
	k.SetDenomSequence(ctx, uint64(len(data.Collections)+1))

	var mtSequence uint64 = 1
	for _, c := range data.Collections {
		// store denom
		k.SetDenom(ctx, *c.Denom)

		for _, m := range c.Mts {
			// increase denom supply
			k.IncreaseDenomSupply(ctx, c.Denom.Id)
			// store mt
			k.SetMT(ctx, c.Denom.Id, m)
			mtSequence++
		}
	}

	// set mt sequence
	k.SetMTSequence(ctx, mtSequence)

	// ---------- init balances ---------- //

	for _, o := range data.Owners {
		addr, err := sdk.AccAddressFromBech32(o.Address)
		if err != nil {
			panic(errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", addr))
		}

		for _, d := range o.Denoms {
			for _, b := range d.Balances {
				// increase supply
				if err := k.IncreaseMTSupply(ctx, d.DenomId, b.MtId, b.Amount); err != nil {
					panic(err)
				}
				// add balance to account
				if err := k.AddBalance(ctx, d.DenomId, b.MtId, b.Amount, addr); err != nil {
					panic(err)
				}
			}
		}
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return k.ExportGenesisState(ctx)
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() *types.GenesisState {
	return types.NewGenesisState([]types.Collection{}, []types.Owner{})
}
