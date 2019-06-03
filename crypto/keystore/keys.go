package keystore

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/irisnet/irishub/crypto/keystore/uuid"
	ctypes "github.com/irisnet/irishub/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"golang.org/x/crypto/pbkdf2"
)

type KeyManager interface {
	ExportAsKeyStore(password string) (*EncryptedKeyJSON, error)
	GetPrivKey() crypto.PrivKey
}

type keyManager struct {
	privKey crypto.PrivKey
	addr    ctypes.AccAddress
}

func (m *keyManager) GetPrivKey() crypto.PrivKey {
	return m.privKey
}

func (m *keyManager) ExportAsKeyStore(password string) (*EncryptedKeyJSON, error) {
	return generateKeyStore(m.GetPrivKey(), password)
}

func NewKeyManager(privKey crypto.PrivKey) KeyManager {
	k := keyManager{}
	k.privKey = privKey
	k.addr = ctypes.AccAddress(privKey.PubKey().Address())
	return &k
}

func NewKeyStoreKeyManager(file string, auth string) (KeyManager, error) {
	k := keyManager{}
	err := k.recoveryFromKeyStore(file, auth)
	return &k, err
}

func (m *keyManager) recoveryFromKeyStore(keystoreFile string, auth string) error {
	if auth == "" {
		return fmt.Errorf("Password is missing ")
	}
	keyJson, err := ioutil.ReadFile(keystoreFile)
	if err != nil {
		return err
	}
	var encryptedKey EncryptedKeyJSON
	err = json.Unmarshal(keyJson, &encryptedKey)
	if err != nil {
		return err
	}
	keyBytes, err := decryptKey(&encryptedKey, auth)
	if err != nil {
		return err
	}
	if len(keyBytes) != 32 {
		return fmt.Errorf("Len of Keybytes is not equal to 32 ")
	}
	var keyBytesArray [32]byte
	copy(keyBytesArray[:], keyBytes[:32])
	privKey := secp256k1.PrivKeySecp256k1(keyBytesArray)
	addr := ctypes.AccAddress(privKey.PubKey().Address())
	m.addr = addr
	m.privKey = privKey
	return nil
}

func generateKeyStore(privateKey crypto.PrivKey, password string) (*EncryptedKeyJSON, error) {
	addr := ctypes.AccAddress(privateKey.PubKey().Address())
	salt, err := GenerateRandomBytes(32)
	if err != nil {
		return nil, err
	}
	iv, err := GenerateRandomBytes(16)
	if err != nil {
		return nil, err
	}
	scryptParamsJSON := make(map[string]interface{}, 4)
	scryptParamsJSON["prf"] = "hmac-sha256"
	scryptParamsJSON["dklen"] = 32
	scryptParamsJSON["salt"] = hex.EncodeToString(salt)
	scryptParamsJSON["c"] = 262144

	cipherParamsJSON := cipherparamsJSON{IV: hex.EncodeToString(iv)}
	derivedKey := pbkdf2.Key([]byte(password), salt, 262144, 32, sha256.New)
	encryptKey := derivedKey[:16]
	secpPrivateKey, ok := privateKey.(secp256k1.PrivKeySecp256k1)
	if !ok {
		return nil, fmt.Errorf(" Only PrivKeySecp256k1 key is supported ")
	}
	cipherText, err := aesCTRXOR(encryptKey, secpPrivateKey[:], iv)
	if err != nil {
		return nil, err
	}

	hasher := sha256.New()
	hasher.Write(derivedKey[16:32])
	hasher.Write(cipherText)
	mac := hasher.Sum(nil)

	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	cryptoStruct := CryptoJSON{
		Cipher:       "aes-256-ctr",
		CipherText:   hex.EncodeToString(cipherText),
		CipherParams: cipherParamsJSON,
		KDF:          "pbkdf2",
		KDFParams:    scryptParamsJSON,
		MAC:          hex.EncodeToString(mac),
	}
	return &EncryptedKeyJSON{
		Address: addr.String(),
		Crypto:  cryptoStruct,
		Id:      id.String(),
		Version: "1",
	}, nil
}
