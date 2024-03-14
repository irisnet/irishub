package keystore

import (
	"encoding/json"
	"fmt"

	tmcrypto "github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/crypto/secp256k1"

	"github.com/cosmos/cosmos-sdk/crypto"
	sdksecp256k1 "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
)

// RecoveryAndExportPrivKeyArmor return the new private key armor from a old keystoreFile
func RecoveryAndExportPrivKeyArmor(keystore []byte, password string) (armor string, err error) {
	priv, err := recoveryFromKeyStore(keystore, password)
	if err != nil {
		return "", err
	}
	return exportPrivKeyArmor(priv, password)
}

func recoveryFromKeyStore(keystore []byte, auth string) (tmcrypto.PrivKey, error) {
	if auth == "" {
		return nil, fmt.Errorf("Password is missing ")
	}

	var encryptedKey EncryptedKeyJSON
	if err := json.Unmarshal(keystore, &encryptedKey); err != nil {
		return nil, err
	}

	keyBytes, err := decryptKey(&encryptedKey, auth)
	if err != nil {
		return nil, err
	}

	if len(keyBytes) != 32 {
		return nil, fmt.Errorf("Len of Keybytes is not equal to 32 ")
	}

	return secp256k1.PrivKey(keyBytes), nil
}

func exportPrivKeyArmor(privKey tmcrypto.PrivKey, password string) (armor string, err error) {
	priv := sdksecp256k1.PrivKey{
		Key: privKey.Bytes(),
	}
	return crypto.EncryptArmorPrivKey(&priv, password, "secp256k1"), nil
}
