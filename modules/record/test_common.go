package record

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"testing"

	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/modules/stake"
	"github.com/irisnet/irishub/modules/mock"
	"github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
)

func createOverflowedData(rep string, bytesLimit int) string {
	overflowedRecordData := rep

	for binary.Size([]byte(overflowedRecordData)) < bytesLimit {
		overflowedRecordData += overflowedRecordData
	}

	return overflowedRecordData
}

func getDataHash(data string) string {
	sum := sha256.Sum256([]byte(data))
	hash := hex.EncodeToString(sum[:])
	return hash
}

func getRecord(ctx sdk.Context, keeper Keeper, hash string) (error, MsgSubmitRecord) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyRecord(hash))
	msg := MsgSubmitRecord{}
	err := keeper.cdc.UnmarshalBinaryLengthPrefixed(bz, &msg)

	return err, msg
}

// checks if two records are equal
func recordEqual(recordA MsgSubmitRecord, recordB MsgSubmitRecord) bool {
	if recordA.SubmitTime == recordB.SubmitTime &&
		recordA.OwnerAddress.String() == recordB.OwnerAddress.String() &&
		recordA.RecordID == recordB.RecordID &&
		recordA.Description == recordB.Description &&
		recordA.DataHash == recordB.DataHash &&
		recordA.DataSize == recordB.DataSize &&
		recordA.Data == recordB.Data {
		return true
	}
	return false
}

// initialize the mock application for this module
func getMockApp(t *testing.T, numGenAccs int) (*mock.App, Keeper, stake.Keeper, []sdk.AccAddress, []crypto.PubKey, []crypto.PrivKey) {
	mapp := mock.NewApp()

	stake.RegisterCodec(mapp.Cdc)
	RegisterCodec(mapp.Cdc)

	keyGov := sdk.NewKVStoreKey("gov")
	keyRecord := sdk.NewKVStoreKey("record")

	sk := stake.NewKeeper(
		mapp.Cdc,
		mapp.KeyStake, mapp.TkeyStake,
		mapp.BankKeeper, mapp.ParamsKeeper.Subspace(stake.DefaultParamspace),
		mapp.RegisterCodespace(stake.DefaultCodespace))
	rk := NewKeeper(mapp.Cdc, keyRecord, mapp.RegisterCodespace(DefaultCodespace))

	mapp.Router().AddRoute("record", []*sdk.KVStoreKey{keyRecord}, NewHandler(rk))

	require.NoError(t, mapp.CompleteSetup(keyGov, keyRecord))

	coin, _ := types.NewDefaultCoinType("iris").ConvertToMinCoin(fmt.Sprintf("%d%s", 1042, "iris"))
	genAccs, addrs, pubKeys, privKeys := mock.CreateGenAccounts(numGenAccs, sdk.Coins{coin})

	mock.SetGenesis(mapp, genAccs)

	return mapp, rk, sk, addrs, pubKeys, privKeys
}
