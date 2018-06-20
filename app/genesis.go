package app

import (
	"encoding/json"
	"errors"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	crypto "github.com/tendermint/go-crypto"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/stake"
)

// State to Unmarshal
type GenesisState struct {
	Accounts  []GenesisAccount   `json:"accounts"`
	StakeData stake.GenesisState `json:"stake"`
}

// GenesisAccount doesn't need pubkey or sequence
type GenesisAccount struct {
	Address sdk.Address `json:"address"`
	Coins   sdk.Coins   `json:"coins"`
}

func NewGenesisAccount(acc *auth.BaseAccount) GenesisAccount {
	return GenesisAccount{
		Address: acc.Address,
		Coins:   acc.Coins,
	}
}

func NewGenesisAccountI(acc auth.Account) GenesisAccount {
	return GenesisAccount{
		Address: acc.GetAddress(),
		Coins:   acc.GetCoins(),
	}
}

// convert GenesisAccount to auth.BaseAccount
func (ga *GenesisAccount) ToAccount() (acc *auth.BaseAccount) {
	return &auth.BaseAccount{
		Address: ga.Address,
		Coins:   ga.Coins.Sort(),
	}
}

var (
	flagName       = "name"
	flagClientHome = "home-client"
	flagOWK        = "owk"
	denom          = "iris"


	// bonded tokens given to genesis validators/accounts
	freeFermionVal  = int64(100)

	totalTokenAmt 	= int64(200000000)
)

// get app init parameters for server init command
func GaiaAppInit() server.AppInit {
	fsAppGenState := pflag.NewFlagSet("", pflag.ContinueOnError)

	fsAppGenTx := pflag.NewFlagSet("", pflag.ContinueOnError)
	fsAppGenTx.String(flagName, "", "validator moniker, required")
	fsAppGenTx.String(flagClientHome, DefaultCLIHome,
		"home directory for the client, used for key generation")
	fsAppGenTx.Bool(flagOWK, false, "overwrite the accounts created")

	return server.AppInit{
		FlagsAppGenState: fsAppGenState,
		FlagsAppGenTx:    fsAppGenTx,
		AppGenTx:         IrisAppGenTx,
		AppGenState:      IrisAppGenStateJSON,
	}
}

// simple genesis tx
type GaiaGenTx struct {
	Name    string        `json:"name"`
	Address sdk.Address   `json:"address"`
	PubKey  crypto.PubKey `json:"pub_key"`
}

// Generate a gaia genesis transaction with flags
func IrisAppGenTx(cdc *wire.Codec, pk crypto.PubKey) (
	appGenTx, cliPrint json.RawMessage, validator tmtypes.GenesisValidator, err error) {
	clientRoot := viper.GetString(flagClientHome)
	overwrite := viper.GetBool(flagOWK)
	name := viper.GetString(flagName)
	if name == "" {
		return nil, nil, tmtypes.GenesisValidator{}, errors.New("Must specify --name (validator moniker)")
	}

	var addr sdk.Address
	var secret string
	addr, secret, err = server.GenerateSaveCoinKey(clientRoot, name, "1234567890", overwrite)
	if err != nil {
		return
	}
	mm := map[string]string{"secret": secret}
	var bz []byte
	bz, err = cdc.MarshalJSON(mm)
	if err != nil {
		return
	}
	cliPrint = json.RawMessage(bz)
	appGenTx,_,validator,err = IrisAppGenTxNF(cdc, pk, addr, name, overwrite)
	return
}

// Generate a gaia genesis transaction without flags
func IrisAppGenTxNF(cdc *wire.Codec, pk crypto.PubKey, addr sdk.Address, name string, overwrite bool) (
	appGenTx, cliPrint json.RawMessage, validator tmtypes.GenesisValidator, err error) {

	var bz []byte
	gaiaGenTx := GaiaGenTx{
		Name:    name,
		Address: addr,
		PubKey:  pk,
	}
	bz, err = wire.MarshalJSONIndent(cdc, gaiaGenTx)
	if err != nil {
		return
	}
	appGenTx = json.RawMessage(bz)

	validator = tmtypes.GenesisValidator{
		PubKey: pk,
		Power:  freeFermionVal,
	}
	return
}

// Create the core parameters for genesis initialization for gaia
// note that the pubkey input is this machines pubkey
func IrisAppGenState(cdc *wire.Codec, appGenTxs []json.RawMessage) (genesisState GenesisState, err error) {

	if len(appGenTxs) == 0 {
		err = errors.New("must provide at least genesis transaction")
		return
	}

	// start with the default staking genesis state
	//stakeData := stake.DefaultGenesisState()
	stakeData := stake.DefaultGenesisState()

	// get genesis flag account information
	genaccs := make([]GenesisAccount, len(appGenTxs))
	for i, appGenTx := range appGenTxs {

		var genTx GaiaGenTx
		err = cdc.UnmarshalJSON(appGenTx, &genTx)
		if err != nil {
			return
		}

		// create the genesis account, give'm few steaks and a buncha token with there name
		accAuth := auth.NewBaseAccountWithAddress(genTx.Address)
		accAuth.Coins = sdk.Coins{
			{denom, totalTokenAmt},
			{"steak", freeFermionVal},
		}
		acc := NewGenesisAccount(&accAuth)
		genaccs[i] = acc
		stakeData.Pool.LooseUnbondedTokens += freeFermionVal // increase the supply

		// add the validator
		if len(genTx.Name) > 0 {
			desc := stake.NewDescription(genTx.Name, "", "", "")
			validator := stake.NewValidator(genTx.Address, genTx.PubKey, desc)
			validator.PoolShares = stake.NewBondedShares(sdk.NewRat(freeFermionVal))
			stakeData.Validators = append(stakeData.Validators, validator)

			// pool logic
			stakeData.Pool.BondedTokens += freeFermionVal
			stakeData.Pool.BondedShares = sdk.NewRat(stakeData.Pool.BondedTokens)
		}
	}

	// create the final app state
	genesisState = GenesisState{
		Accounts:  genaccs,
		StakeData: stakeData,
	}
	return
}

// IrisAppGenState but with JSON
func IrisAppGenStateJSON(cdc *wire.Codec, appGenTxs []json.RawMessage) (appState json.RawMessage, err error) {

	// create the final app state
	genesisState, err := IrisAppGenState(cdc, appGenTxs)
	if err != nil {
		return nil, err
	}
	appState, err = wire.MarshalJSONIndent(cdc, genesisState)
	return
}

// TODO
func createGenesisState() stake.GenesisState{
	return stake.GenesisState{
		Pool:   stake.Pool{
			LooseUnbondedTokens:     0,
			BondedTokens:            0,
			UnbondingTokens:         0,
			UnbondedTokens:          0,
			BondedShares:            sdk.ZeroRat(),
			UnbondingShares:         sdk.ZeroRat(),
			UnbondedShares:          sdk.ZeroRat(),
			InflationLastTime:       0,
			Inflation:               sdk.NewRat(7, 100),
			DateLastCommissionReset: 0,
			PrevBondedShares:        sdk.ZeroRat(),
		},
		Params: stake.Params{
			InflationRateChange: sdk.NewRat(13, 100),
			InflationMax:        sdk.NewRat(20, 100),
			InflationMin:        sdk.NewRat(7, 100),
			GoalBonded:          sdk.NewRat(67, 100),
			MaxValidators:       100,
			BondDenom:           "steak",
		},
	}
}