package types

// MUST be loaded before running
import (
	"os"
	"path/filepath"

	"github.com/cometbft/cometbft/crypto"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/evmos/ethermint/ethereum/eip712"
	etherminttypes "github.com/evmos/ethermint/types"

	tokentypes "github.com/irisnet/irismod/modules/token/types"
	tokenv1 "github.com/irisnet/irismod/modules/token/types/v1"
)

const (
	// AppName is the name of the app
	AppName = "IrisApp"
)

var (
	// NativeToken represents the native token
	NativeToken     tokenv1.Token
	// EvmToken represents the EVM token
	EvmToken        tokenv1.Token
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string
)

// Init initializes the legacyAmino and interfaceRegistry.
// Parameters: legacyAmino *codec.LegacyAmino, interfaceRegistry types.InterfaceRegistry.
func Init(legacyAmino *codec.LegacyAmino, interfaceRegistry types.InterfaceRegistry) {
	initApp()
	initEVM(legacyAmino, interfaceRegistry)
}

func initApp() {
	// set bech32 prefix
	ConfigureBech32Prefix()

	// set coin denom regexs
	sdk.SetCoinDenomRegex(func() string {
		return `[a-zA-Z][a-zA-Z0-9/-]{2,127}`
	})

	NativeToken = tokenv1.Token{
		Symbol:        "iris",
		Name:          "Irishub staking token",
		Scale:         6,
		MinUnit:       "uiris",
		InitialSupply: 2000000000,
		MaxSupply:     10000000000,
		Mintable:      true,
		Owner:         sdk.AccAddress(crypto.AddressHash([]byte(tokentypes.ModuleName))).String(),
	}

	EvmToken = tokenv1.Token{
		Symbol:        "eris",
		Name:          "IRISHub EVM Fee Token",
		Scale:         18,
		MinUnit:       "weris",
		InitialSupply: 0,
		MaxSupply:     10000000000,
		Mintable:      true,
		Owner:         sdk.AccAddress(crypto.AddressHash([]byte(tokentypes.ModuleName))).String(),
	}

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, ".iris")
	owner, err := sdk.AccAddressFromBech32(NativeToken.Owner)
	if err != nil {
		panic(err)
	}

	// replace the default token
	tokenv1.SetNativeToken(
		NativeToken.Symbol,
		NativeToken.Name,
		NativeToken.MinUnit,
		NativeToken.Scale,
		NativeToken.InitialSupply,
		NativeToken.MaxSupply,
		NativeToken.Mintable,
		owner,
	)
}

func initEVM(legacyAmino *codec.LegacyAmino, interfaceRegistry types.InterfaceRegistry) {
	etherminttypes.InjectChainIDParser(parseChainID)
	eip712.InjectCodec(eip712.Codec{
		InterfaceRegistry: interfaceRegistry,
		Amino:             legacyAmino,
	})
}
