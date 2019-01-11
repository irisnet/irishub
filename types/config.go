package types

import (
	"sync"
	"fmt"
)

const (
	NativeTokenName     = "iris"
	NativeTokenMinDenom = "iris-atto"
	DefaultKeyPass      = "1234567890"
	Testnet             = "testnet"
	Mainnet             = "mainnet"
	InvariantPanic      = "panic"
	InvariantError      = "error"
)

var (
	IRIS = NewDefaultCoinType(NativeTokenName)
	InitialIssue = NewIntWithDecimal(2, 9)	// 2 billion
	FreeToken4Val, _ = IRIS.ConvertToMinCoin(fmt.Sprintf("%d%s", int64(100), NativeTokenName))
	FreeToken4Acc, _ = IRIS.ConvertToMinCoin(fmt.Sprintf("%d%s", int64(150), NativeTokenName))
)

// Can be configured through environment variables
var (
	NetworkType = Testnet
	InvariantLevel = InvariantPanic
	Bech32PrefixAccAddr = "faa"		// account address
	Bech32PrefixAccPub = "fap"		// account public key
	Bech32PrefixValAddr = "fva"		// validator operator address
	Bech32PrefixValPub = "fvp"		// validator operator public key
	Bech32PrefixConsAddr = "fca"	// consensus node address
	Bech32PrefixConsPub = "fcp"		// consensus node public key
)

var (
	// Initializing an instance of Config
	config = &Config{
		sealed: true,
		bech32AddressPrefix: map[string]string{
			"account_addr":   Bech32PrefixAccAddr,
			"validator_addr": Bech32PrefixValAddr,
			"consensus_addr": Bech32PrefixConsAddr,
			"account_pub":    Bech32PrefixAccPub,
			"validator_pub":  Bech32PrefixValPub,
			"consensus_pub":  Bech32PrefixConsPub,
		},
	}
)

type Config struct {
	mtx                 sync.RWMutex
	sealed              bool
	bech32AddressPrefix map[string]string
}

// An Invariant is a function which tests a particular invariant.
// If the invariant has been broken, it should return an error
// containing a descriptive message about what happened.
type Invariant func(ctx Context) error

/*
func InitBech32Prefix() {
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
	config.Seal()
}
*/
// GetConfig returns the config instance for the SDK.
func GetConfig() *Config {
	return config
}
/*
func (config *Config) assertNotSealed() {
	config.mtx.Lock()
	defer config.mtx.Unlock()

	if config.sealed {
		panic("Config is sealed")
	}
}

// SetBech32PrefixForAccount builds the Config with Bech32 addressPrefix and publKeyPrefix for accounts
// and returns the config instance
func (config *Config) SetBech32PrefixForAccount(addressPrefix, pubKeyPrefix string) {
	config.assertNotSealed()
	config.bech32AddressPrefix["account_addr"] = addressPrefix
	config.bech32AddressPrefix["account_pub"] = pubKeyPrefix
}

// SetBech32PrefixForValidator builds the Config with Bech32 addressPrefix and publKeyPrefix for validators
//  and returns the config instance
func (config *Config) SetBech32PrefixForValidator(addressPrefix, pubKeyPrefix string) {
	config.assertNotSealed()
	config.bech32AddressPrefix["validator_addr"] = addressPrefix
	config.bech32AddressPrefix["validator_pub"] = pubKeyPrefix
}

// SetBech32PrefixForConsensusNode builds the Config with Bech32 addressPrefix and publKeyPrefix for consensus nodes
// and returns the config instance
func (config *Config) SetBech32PrefixForConsensusNode(addressPrefix, pubKeyPrefix string) {
	config.assertNotSealed()
	config.bech32AddressPrefix["consensus_addr"] = addressPrefix
	config.bech32AddressPrefix["consensus_pub"] = pubKeyPrefix
}

// Seal seals the config such that the config state could not be modified further
func (config *Config) Seal() *Config {
	config.mtx.Lock()
	defer config.mtx.Unlock()

	config.sealed = true
	return config
}
*/
// GetBech32AccountAddrPrefix returns the Bech32 prefix for account address
func (config *Config) GetBech32AccountAddrPrefix() string {
	return config.bech32AddressPrefix["account_addr"]
}

// GetBech32ValidatorAddrPrefix returns the Bech32 prefix for validator address
func (config *Config) GetBech32ValidatorAddrPrefix() string {
	return config.bech32AddressPrefix["validator_addr"]
}

// GetBech32ConsensusAddrPrefix returns the Bech32 prefix for consensus node address
func (config *Config) GetBech32ConsensusAddrPrefix() string {
	return config.bech32AddressPrefix["consensus_addr"]
}

// GetBech32AccountPubPrefix returns the Bech32 prefix for account public key
func (config *Config) GetBech32AccountPubPrefix() string {
	return config.bech32AddressPrefix["account_pub"]
}

// GetBech32ValidatorPubPrefix returns the Bech32 prefix for validator public key
func (config *Config) GetBech32ValidatorPubPrefix() string {
	return config.bech32AddressPrefix["validator_pub"]
}

// GetBech32ConsensusPubPrefix returns the Bech32 prefix for consensus node public key
func (config *Config) GetBech32ConsensusPubPrefix() string {
	return config.bech32AddressPrefix["consensus_pub"]
}
