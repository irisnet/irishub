package keystore

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func TestKeystore(t *testing.T) {
	defer os.Remove("TestGenerateKeyStoreNoError.json")

	secret, _ := hex.DecodeString("55A3160577979EC014A2CE85C430E1FF0FF06EFD230B7CE41AEAE2EF00EDF175")
	var priv [32]byte
	copy(priv[:], secret[0:32])
	km := NewKeyManager(secp256k1.PrivKeySecp256k1(priv))
	encryPlain1, err := km.GetPrivKey().Sign([]byte("test plain"))
	assert.NoError(t, err)

	encryptedKeyJSON, err := km.ExportAsKeyStore("testpassword")
	assert.NoError(t, err)

	bz, err := json.Marshal(encryptedKeyJSON)
	assert.NoError(t, err)

	err = ioutil.WriteFile("TestGenerateKeyStoreNoError.json", bz, 0660)
	assert.NoError(t, err)

	newkm, err := NewKeyStoreKeyManager("TestGenerateKeyStoreNoError.json", "testpassword")
	assert.NoError(t, err)

	encryPlain2, err := newkm.GetPrivKey().Sign([]byte("test plain"))
	assert.NoError(t, err)
	assert.True(t, bytes.Equal(encryPlain1, encryPlain2))

	assert.Equal(t, km.GetPrivKey().Bytes(), newkm.GetPrivKey().Bytes())
	assert.Equal(t, km.GetPrivKey().PubKey(), newkm.GetPrivKey().PubKey())
}
