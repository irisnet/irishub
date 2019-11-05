package config

import (
	"fmt"

	sdk "github.com/irisnet/irishub/types"
)

const (
	defaultMinimumFees = ""
)

// BaseConfig defines the server's basic configuration
type BaseConfig struct {
	// Tx minimum fee
	MinFees string `mapstructure:"minimum_fees"`

	// Enable invariant check, ignore this flag on testnet
	CheckInvariant bool `mapstructure:"check_invariant"`

	// Enable track coin flow
	TrackCoinFlow bool `mapstructure:"track_coin_flow"`
}

// Config defines the server's top level configuration
type Config struct {
	BaseConfig `mapstructure:",squash"`
}

// SetMinimumFee sets the minimum fee.
func (c *Config) SetMinimumFees(fees sdk.Coins) { c.MinFees = fees.String() }

// SetMinimumFee sets the minimum fee.
func (c *Config) MinimumFees() sdk.Coins {
	fees, err := sdk.ParseCoins(c.MinFees)
	if err != nil {
		panic(fmt.Sprintf("invalid minimum fees: %v", err))
	}
	return fees
}

// DefaultConfig returns server's default configuration.
func DefaultConfig() *Config {
	return &Config{BaseConfig{MinFees: defaultMinimumFees, CheckInvariant: false, TrackCoinFlow: false}}
}
