package keystore

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
)

var (
	errDecrypt = errors.New("could not decrypt key with given passphrase")
)

// PlainKeyJSON define a struct TODO
type PlainKeyJSON struct {
	Address    string `json:"address"`
	PrivateKey string `json:"privatekey"`
	ID         string `json:"id"`
	Version    int    `json:"version"`
}

// EncryptedKeyJSON define a struct TODO
type EncryptedKeyJSON struct {
	Address string     `json:"address"`
	Crypto  CryptoJSON `json:"crypto"`
	ID      string     `json:"id"`
	Version string     `json:"version"`
}

// CryptoJSON define a struct TODO
type CryptoJSON struct {
	Cipher       string                 `json:"cipher"`
	CipherText   string                 `json:"ciphertext"`
	CipherParams cipherparamsJSON       `json:"cipherparams"`
	KDF          string                 `json:"kdf"`
	KDFParams    map[string]interface{} `json:"kdfparams"`
	MAC          string                 `json:"mac"`
}

type cipherparamsJSON struct {
	IV string `json:"iv"`
}

func decryptKey(keyProtected *EncryptedKeyJSON, auth string) ([]byte, error) {
	mac, err := hex.DecodeString(keyProtected.Crypto.MAC)
	if err != nil {
		return nil, err
	}

	iv, err := hex.DecodeString(keyProtected.Crypto.CipherParams.IV)
	if err != nil {
		return nil, err
	}

	cipherText, err := hex.DecodeString(keyProtected.Crypto.CipherText)
	if err != nil {
		return nil, err
	}

	derivedKey, err := getKDFKey(keyProtected.Crypto, auth)
	if err != nil {
		return nil, err
	}

	bufferValue := make([]byte, len(cipherText)+16)
	copy(bufferValue[0:16], derivedKey[16:32])
	copy(bufferValue[16:], cipherText[:])
	calculatedMAC := sha256.Sum256((bufferValue))
	if !bytes.Equal(calculatedMAC[:], mac) {
		return nil, errDecrypt
	}

	plainText, err := aesCTRXOR(derivedKey[:16], cipherText, iv)
	if err != nil {
		return nil, err
	}
	return plainText, err
}

func getKDFKey(cryptoJSON CryptoJSON, auth string) ([]byte, error) {
	authArray := []byte(auth)
	if cryptoJSON.KDFParams["salt"] == nil || cryptoJSON.KDFParams["dklen"] == nil ||
		cryptoJSON.KDFParams["c"] == nil || cryptoJSON.KDFParams["prf"] == nil {
		return nil, errors.New("invalid KDF params, must contains c, dklen, prf and salt")
	}
	salt, err := hex.DecodeString(cryptoJSON.KDFParams["salt"].(string))
	if err != nil {
		return nil, err
	}
	dkLen := ensureInt(cryptoJSON.KDFParams["dklen"])

	if cryptoJSON.KDF == "pbkdf2" {
		c := ensureInt(cryptoJSON.KDFParams["c"])
		prf := cryptoJSON.KDFParams["prf"].(string)
		if prf != "hmac-sha256" {
			return nil, fmt.Errorf("Unsupported PBKDF2 PRF: %s", prf)
		}
		key := pbkdf2.Key(authArray, salt, c, dkLen, sha256.New)
		return key, nil
	}
	return nil, fmt.Errorf("Unsupported KDF: %s", cryptoJSON.KDF)
}

func ensureInt(x interface{}) int {
	res, ok := x.(int)
	if !ok {
		res = int(x.(float64))
	}
	return res
}

func aesCTRXOR(key, inText, iv []byte) ([]byte, error) {
	// AES-128 is selected due to size of encryptKey.
	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(aesBlock, iv)
	outText := make([]byte, len(inText))
	stream.XORKeyStream(outText, inText)
	return outText, err
}
