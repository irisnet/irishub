package types

import (
	"fmt"
)

const (
	Iris           = "iris"
	IrisAtto       = "iris-atto"
	DefaultKeyPass = "1234567890"
	Testnet        = "testnet"
	Mainnet        = "mainnet"
	InvariantPanic = "panic"
	InvariantError = "error"
)

var (
	InitialIssue     = NewIntWithDecimal(2, 9) // 2 billion
	FreeToken4Val, _ = IrisCoinType.ConvertToMinDenomCoin(fmt.Sprintf("%d%s", int64(100), Iris))
	FreeToken4Acc, _ = IrisCoinType.ConvertToMinDenomCoin(fmt.Sprintf("%d%s", int64(150), Iris))
)

// Can be configured through environment variables
var (
	NetworkType    = Testnet
	InvariantLevel = InvariantPanic
)

var (
	testnetConfig = &Config{
		bech32AddressPrefix: map[string]string{
			"account_addr":   "faa",
			"validator_addr": "fva",
			"consensus_addr": "fca",
			"account_pub":    "fap",
			"validator_pub":  "fvp",
			"consensus_pub":  "fcp",
		},
	}
	mainnetConfig = &Config{
		bech32AddressPrefix: map[string]string{
			"account_addr":   "iaa",
			"validator_addr": "iva",
			"consensus_addr": "ica",
			"account_pub":    "iap",
			"validator_pub":  "ivp",
			"consensus_pub":  "icp",
		},
	}
)

type Config struct {
	bech32AddressPrefix map[string]string
}

// An Invariant is a function which tests a particular invariant.
// If the invariant has been broken, it should return an error
// containing a descriptive message about what happened.
type Invariant func(ctx Context) error

func SetNetworkType(networkType string) {
	NetworkType = networkType
}

// GetConfig returns the config instance for the corresponding network type
func GetConfig() *Config {
	if NetworkType == Mainnet {
		return mainnetConfig
	}
	return testnetConfig
}

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
