package keeper

import (
	"os"
	"testing"

	"github.com/irisnet/irishub/modules/guardian/internal/types"

	"encoding/hex"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

var (
	pks = []crypto.PubKey{
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB50"),
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB51"),
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB52"),
	}
	addrs = []sdk.AccAddress{
		sdk.AccAddress(pks[0].Address()),
		sdk.AccAddress(pks[1].Address()),
		sdk.AccAddress(pks[2].Address()),
	}
)

func newPubKey(pk string) (res crypto.PubKey) {
	pkBytes, err := hex.DecodeString(pk)
	if err != nil {
		panic(err)
	}
	var pkEd ed25519.PubKeyEd25519
	copy(pkEd[:], pkBytes[:])
	return pkEd
}

func createTestCodec() *codec.Codec {
	cdc := codec.New()
	sdk.RegisterCodec(cdc)
	types.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

func createTestInput(t *testing.T) (sdk.Context, Keeper) {
	keyProf := sdk.NewKVStoreKey("guardian")

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyProf, sdk.StoreTypeIAVL, db)

	err := ms.LoadLatestVersion()
	require.Nil(t, err)
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewTMLogger(os.Stdout))
	cdc := createTestCodec()

	keeper := NewKeeper(cdc, keyProf, types.DefaultCodespace)

	return ctx, keeper
}
