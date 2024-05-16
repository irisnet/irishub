package params

// MUST be loaded before running
import (
	"fmt"
	"math/big"
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
	AppName       = "IrisApp"
	EIP155ChainID = "6688"
)

var (
	// BaseToken represents the native token
	BaseToken tokenv1.Token
	// EvmToken represents the EVM token
	EvmToken tokenv1.Token
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string
)

func init() {
	// set bech32 prefix
	ConfigureBech32Prefix()

	// set coin denom regexs
	sdk.SetCoinDenomRegex(func() string {
		return `[a-zA-Z][a-zA-Z0-9/-]{2,127}`
	})

	BaseToken = tokenv1.Token{
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
	sdk.DefaultBondDenom = BaseToken.MinUnit

	if err := sdk.RegisterDenom(BaseToken.Symbol, sdk.OneDec()); err != nil {
		panic(err)
	}

	if err := sdk.RegisterDenom(BaseToken.MinUnit, sdk.NewDecFromIntWithPrec(sdk.OneInt(), int64(BaseToken.Scale))); err != nil {
		panic(err)
	}

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, ".iris")
	owner, err := sdk.AccAddressFromBech32(BaseToken.Owner)
	if err != nil {
		panic(err)
	}

	// replace the default token
	tokenv1.SetNativeToken(
		BaseToken.Symbol,
		BaseToken.Name,
		BaseToken.MinUnit,
		BaseToken.Scale,
		BaseToken.InitialSupply,
		BaseToken.MaxSupply,
		BaseToken.Mintable,
		owner,
	)

	etherminttypes.InjectChainIDParser(parseChainID)
}

// InjectCodec injects an app codec
func InjectCodec(legacyAmino *codec.LegacyAmino, interfaceRegistry types.InterfaceRegistry) {
	eip712.InjectCodec(eip712.Codec{
		InterfaceRegistry: interfaceRegistry,
		Amino:             legacyAmino,
	})
}

func parseChainID(_ string) (*big.Int, error) {
	eip155ChainID, ok := new(big.Int).SetString(EIP155ChainID, 10)
	if !ok {
		return nil, fmt.Errorf("invalid chain-id: %s", EIP155ChainID)
	}
	return eip155ChainID, nil
}
