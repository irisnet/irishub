### Sleep is needed otherwise the relayer crashes when trying to init
sleep 1s

## Remove exist keys 
$HERMES_BINARY --config ./network/hermes/config.toml keys delete --chain test-1 --all
$HERMES_BINARY --config ./network/hermes/config.toml keys delete --chain test-2 --all

### Mnemonic Files
echo "alley afraid soup fall idea toss can goose become valve initial strong forward bright dish figure check leopard decide warfare hub unusual join cart" > mnemonic_test_1
echo "record gift you once hip style during joke field prize dust unique length more pencil transfer quit train device arrive energy sort steak upset" > mnemonic_test_2

### Restore Keys
$HERMES_BINARY --config ./network/hermes/config.toml keys add --chain test-1 --mnemonic-file mnemonic_test_1
sleep 5s

$HERMES_BINARY --config ./network/hermes/config.toml keys add --chain test-2 --mnemonic-file mnemonic_test_2
sleep 5s

### Remove Mnemonic
rm mnemonic_test_1 mnemonic_test_2