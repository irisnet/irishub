package address

import sdk "github.com/cosmos/cosmos-sdk/types"

const (

	// Bech32ChainPrefix defines the prefix of this chain
	Bech32ChainPrefix = "i"

	// PrefixAcc is the prefix for account
	PrefixAcc = "a"

	// PrefixValidator is the prefix for validator keys
	PrefixValidator = "v"

	// PrefixConsensus is the prefix for consensus keys
	PrefixConsensus = "c"

	// PrefixPublic is the prefix for public
	PrefixPublic = "p"

	// PrefixAddress is the prefix for address
	PrefixAddress = "a"

	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
	Bech32PrefixAccAddr = Bech32ChainPrefix + PrefixAcc + PrefixAddress
	// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key
	Bech32PrefixAccPub = Bech32ChainPrefix + PrefixAcc + PrefixPublic
	// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
	Bech32PrefixValAddr = Bech32ChainPrefix + PrefixValidator + PrefixAddress
	// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
	Bech32PrefixValPub = Bech32ChainPrefix + PrefixValidator + PrefixPublic
	// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
	Bech32PrefixConsAddr = Bech32ChainPrefix + PrefixConsensus + PrefixAddress
	// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
	Bech32PrefixConsPub = Bech32ChainPrefix + PrefixConsensus + PrefixPublic
)

func ConfigureBech32Prefix() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
	config.Seal()
}
