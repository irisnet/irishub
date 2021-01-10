package keystore

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRecoveryAndExportPrivKeyArmor(t *testing.T) {
	keystore := `{"version":"1","id":"65177bc2-8240-4024-8180-dd0b2d888903","address":"faa1ljemm0yznz58qxxs8xyak7fashcfxf5lssn6jm","crypto":{"ciphertext":"793acc81ed7d3f8aead7872f81cc7297e0527ab9ee87a24f8aa7de6a6b4072e9","cipherparams":{"iv":"7ebe22befa6b278f0f348fe9e3f7c524"},"cipher":"aes-128-ctr","kdf":"pbkdf2","kdfparams":{"dklen":32,"salt":"0fa96f07f73d3dfe2bff410b708de347080a326c898e2d5631af4d598e851401","c":262144,"prf":"hmac-sha256"},"mac":"15467c52ade57fd59200544612cccd2310825f8378d3f52228b52d07b56fbdba"}}`
	_, err := RecoveryAndExportPrivKeyArmor([]byte(keystore), "1234567890")
	require.NoError(t, err)
}
