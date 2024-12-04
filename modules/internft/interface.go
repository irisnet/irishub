package internft

import (
	"context"
	nftkeeper "cosmossdk.io/x/nft/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	nfttypes "mods.irisnet.org/modules/nft/types"
)

type (
	// AccountKeeper defines the contract required for account APIs.
	AccountKeeper interface {
		NewAccountWithAddress(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
		// Set an account in the store.
		SetAccount(context.Context, sdk.AccountI)
		GetModuleAddress(name string) sdk.AccAddress
	}
	// InterNftKeeper defines the ICS721 Keeper
	InterNftKeeper struct {
		nk  nftkeeper.Keeper
		cdc codec.Codec
		ak  AccountKeeper
		cb  nfttypes.ClassBuilder
		tb  nfttypes.TokenBuilder
	}

	InterClass struct {
		ID   string
		URI  string
		Data string
	}

	InterToken struct {
		ClassID string
		ID      string
		URI     string
		Data    string
	}
)

func (d InterClass) GetID() string      { return d.ID }
func (d InterClass) GetURI() string     { return d.URI }
func (d InterClass) GetData() string    { return d.Data }
func (t InterToken) GetClassID() string { return t.ClassID }
func (t InterToken) GetID() string      { return t.ID }
func (t InterToken) GetURI() string     { return t.URI }
func (t InterToken) GetData() string    { return t.Data }
