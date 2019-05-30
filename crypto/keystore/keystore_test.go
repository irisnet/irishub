package keystore

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	ctypes "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func TestKeystore(t *testing.T) {
	secret, _ := hex.DecodeString("55A3160577979EC014A2CE85C430E1FF0FF06EFD230B7CE41AEAE2EF00EDF175")
	var priv [32]byte
	copy(priv[:], secret[0:32])
	keyManager := keyManager{
		privKey: secp256k1.PrivKeySecp256k1(priv),
	}

	pubkey, _ := ctypes.Bech32ifyAccPub(secp256k1.PrivKeySecp256k1(priv).PubKey())
	println(pubkey)

	encryptedKeyJSON, _ := keyManager.ExportAsKeyStore("1234567890")

	jsonStr, _ := json.Marshal(encryptedKeyJSON)
	println(string(jsonStr))

	var prvi32 [32]byte
	prvi32 = keyManager.GetPrivKey().(secp256k1.PrivKeySecp256k1)
	println(hex.EncodeToString(prvi32[:]))
	require.Equal(t, hex.EncodeToString(secret), hex.EncodeToString(prvi32[:]))
}
