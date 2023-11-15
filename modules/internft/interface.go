package internft

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	nftkeeper "github.com/cosmos/cosmos-sdk/x/nft/keeper"
	nfttypes "github.com/irisnet/irismod/modules/nft/types"
)

type (
	// AccountKeeper defines the contract required for account APIs.
	AccountKeeper interface {
		NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
		// Set an account in the store.
		SetAccount(sdk.Context, authtypes.AccountI)
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
