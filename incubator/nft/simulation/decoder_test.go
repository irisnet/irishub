package simulation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto/ed25519"
	tmkv "github.com/tendermint/tendermint/libs/kv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/modules/incubator/nft/internal/types"
)

var (
	delPk1 = ed25519.GenPrivKey().PubKey()
	addr   = sdk.AccAddress(delPk1.Address())
)

func makeTestCodec() (cdc *codec.Codec) {
	cdc = codec.New()
	sdk.RegisterCodec(cdc)
	types.RegisterCodec(cdc)
	return
}

func TestDecodeStore(t *testing.T) {
	cdc := makeTestCodec()
	nft := types.NewBaseNFT("1", addr, "token URI")
	collection := types.NewCollection("kitties", types.NFTs{&nft})
	idCollection := types.NewIDCollection("kitties", []string{"1", "2", "3"})

	kvPairs := tmkv.Pairs{
		tmkv.Pair{Key: types.GetCollectionKey("kitties"), Value: cdc.MustMarshalBinaryLengthPrefixed(collection)},
		tmkv.Pair{Key: types.GetOwnerKey(addr, "kitties"), Value: cdc.MustMarshalBinaryLengthPrefixed(idCollection)},
		tmkv.Pair{Key: []byte{0x99}, Value: []byte{0x99}},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"collections", fmt.Sprintf("%v\n%v", collection, collection)},
		{"owners", fmt.Sprintf("%v\n%v", idCollection, idCollection)},
		{"other", ""},
	}

	for i, tt := range tests {
		tt, i := tt, i
		t.Run(tt.name, func(t *testing.T) {
			switch i {
			case len(tests) - 1:
				require.Panics(t, func() { DecodeStore(cdc, kvPairs[i], kvPairs[i]) }, tt.name)
			default:
				require.Equal(t, tt.expectedLog, DecodeStore(cdc, kvPairs[i], kvPairs[i]), tt.name)
			}
		})
	}
}
