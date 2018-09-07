package app

import (
	"encoding/json"
	"errors"

	"github.com/spf13/pflag"
	"github.com/tendermint/tendermint/crypto"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"time"
	"github.com/irisnet/irishub/types"
	"fmt"
)

// State to Unmarshal
type GenesisState struct {
	Accounts  []GenesisAccount   `json:"accounts"`
	StakeData stake.GenesisState `json:"stake"`
}

// GenesisAccount doesn't need pubkey or sequence
type GenesisAccount struct {
	Address sdk.AccAddress `json:"address"`
	Coins   sdk.Coins      `json:"coins"`
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
	denom          = "iris-atto"
	feeAmt   = int64(100)
	IrisCt = types.NewDefaultCoinType(denom)
	freeFermionVal ,_ = IrisCt.ConvertToMinCoin(fmt.Sprintf("%d%s",feeAmt,denom))
)

const defaultUnbondingTime time.Duration = 60 * 10 * time.Second

// get app init parameters for server init command
func IrisAppInit() server.AppInit {
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
type IrisGenTx struct {
	Name    string         `json:"name"`
	Address sdk.AccAddress `json:"address"`
	PubKey  string         `json:"pub_key"`
}

// Generate a gaia genesis transaction with flags
func IrisAppGenTx(cdc *wire.Codec, pk crypto.PubKey, genTxConfig config.GenTx) (
	appGenTx, cliPrint json.RawMessage, validator tmtypes.GenesisValidator, err error) {

	if genTxConfig.Name == "" {
		return nil, nil, tmtypes.GenesisValidator{}, errors.New("Must specify --name (validator moniker)")
	}

	var addr sdk.AccAddress
	var secret string
	addr, secret, err = server.GenerateSaveCoinKey(genTxConfig.CliRoot, genTxConfig.Name, "1234567890", genTxConfig.Overwrite)
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
	appGenTx, _, validator, err = IrisAppGenTxNF(cdc, pk, addr, genTxConfig.Name)
	return
}

// Generate a gaia genesis transaction without flags
func IrisAppGenTxNF(cdc *wire.Codec, pk crypto.PubKey, addr sdk.AccAddress, name string) (
	appGenTx, cliPrint json.RawMessage, validator tmtypes.GenesisValidator, err error) {

	var bz []byte
	gaiaGenTx := IrisGenTx{
		Name:    name,
		Address: addr,
		PubKey:  sdk.MustBech32ifyAccPub(pk),
	}
	bz, err = wire.MarshalJSONIndent(cdc, gaiaGenTx)
	if err != nil {
		return
	}
	appGenTx = json.RawMessage(bz)

	validator = tmtypes.GenesisValidator{
		PubKey: pk,
		Power:  feeAmt,
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

	stakeData := createGenesisState()
	genaccs := make([]GenesisAccount, len(appGenTxs))
	for i, appGenTx := range appGenTxs {

		var genTx IrisGenTx
		err = cdc.UnmarshalJSON(appGenTx, &genTx)
		if err != nil {
			return
		}

		// create the genesis account, give'm few steaks and a buncha token with there name
		accAuth := auth.NewBaseAccountWithAddress(genTx.Address)
		accAuth.Coins = sdk.Coins{
			freeFermionVal,
		}
		acc := NewGenesisAccount(&accAuth)
		genaccs[i] = acc
		stakeData.Pool.LooseTokens = stakeData.Pool.LooseTokens.Add(sdk.NewRatFromInt(freeFermionVal.Amount)) // increase the supply

		// add the validator
		if len(genTx.Name) > 0 {
			desc := stake.NewDescription(genTx.Name, "", "", "")
			validator := stake.NewValidator(genTx.Address,
				sdk.MustGetAccPubKeyBech32(genTx.PubKey), desc)

			stakeData.Pool.LooseTokens = stakeData.Pool.LooseTokens.Add(sdk.NewRatFromInt(freeFermionVal.Amount))

			// add some new shares to the validator
			var issuedDelShares sdk.Rat
			validator, stakeData.Pool, issuedDelShares = validator.AddTokensFromDel(stakeData.Pool, freeFermionVal.Amount)
			//validator.TokenPrecision = stakeData.Params.DenomPrecision
			stakeData.Validators = append(stakeData.Validators, validator)

			// create the self-delegation from the issuedDelShares
			delegation := stake.Delegation{
				DelegatorAddr: validator.Owner,
				ValidatorAddr: validator.Owner,
				Shares:        issuedDelShares,
				Height:        0,
			}

			stakeData.Bonds = append(stakeData.Bonds, delegation)
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

func createGenesisState() stake.GenesisState {
	return stake.GenesisState{
		Pool: stake.Pool{
			LooseTokens:             sdk.ZeroRat(),
			BondedTokens:            sdk.ZeroRat(),
			InflationLastTime:       time.Unix(0, 0),
			Inflation:               sdk.NewRat(7, 100),
			DateLastCommissionReset: 0,
			PrevBondedShares:        sdk.ZeroRat(),
		},
		Params: stake.Params{
			InflationRateChange: sdk.NewRat(13, 100),
			InflationMax:        sdk.NewRat(20, 100),
			InflationMin:        sdk.NewRat(7, 100),
			GoalBonded:          sdk.NewRat(67, 100),
			UnbondingTime:       defaultUnbondingTime,
			MaxValidators:       100,
			BondDenom:           denom,
		},
	}
}
