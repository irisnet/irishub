package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/irisnet/irishub/modules/gov"
	"github.com/irisnet/irishub/modules/upgrade"
	"github.com/irisnet/irishub/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"time"
	"github.com/irisnet/irishub/modules/gov/params"
)

var (
	Denom             = "iris"
	feeAmt            = int64(100)
	IrisCt            = types.NewDefaultCoinType(Denom)
	FreeFermionVal, _ = IrisCt.ConvertToMinCoin(fmt.Sprintf("%d%s", feeAmt, Denom))
	FreeFermionAcc, _ = IrisCt.ConvertToMinCoin(fmt.Sprintf("%d%s", int64(150), Denom))
)

const (
	defaultUnbondingTime time.Duration = 60 * 10 * time.Second
	// DefaultKeyPass contains the default key password for genesis transactions
	DefaultKeyPass = "1234567890"
)

// State to Unmarshal
type GenesisState struct {
	Accounts     []GenesisAccount      `json:"accounts"`
	StakeData    stake.GenesisState    `json:"stake"`
	MintData     mint.GenesisState     `json:"mint"`
	DistrData    distr.GenesisState    `json:"distr"`
	GovData      gov.GenesisState      `json:"gov"`
	UpgradeData  upgrade.GenesisState  `json:"upgrade"`
	SlashingData slashing.GenesisState `json:"slashing"`
	GenTxs       []json.RawMessage     `json:"gentxs"`
}

func NewGenesisState(accounts []GenesisAccount, stakeData stake.GenesisState, mintData mint.GenesisState,
	distrData distr.GenesisState, govData gov.GenesisState, upgradeData upgrade.GenesisState, slashingData slashing.GenesisState) GenesisState {

	return GenesisState{
		Accounts:     accounts,
		StakeData:    stakeData,
		MintData:     mintData,
		DistrData:    distrData,
		GovData:      govData,
		UpgradeData:  upgradeData,
		SlashingData: slashingData,
	}
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

// get app init parameters for server init command
func IrisAppInit() server.AppInit {
	return server.AppInit{
		AppGenState: IrisAppGenStateJSON,
	}
}

// Create the core parameters for genesis initialization for iris
// note that the pubkey input is this machines pubkey
func IrisAppGenState(cdc *codec.Codec, appGenTxs []json.RawMessage) (genesisState GenesisState, err error) {
	if len(appGenTxs) == 0 {
		err = errors.New("must provide at least genesis transaction")
		return
	}

	// start with the default staking genesis state
	stakeData := createGenesisState()
	slashingData := slashing.DefaultGenesisState()

	// get genesis flag account information
	genaccs := make([]GenesisAccount, len(appGenTxs))

	for i, appGenTx := range appGenTxs {
		var tx auth.StdTx
		err = cdc.UnmarshalJSON(appGenTx, &tx)
		if err != nil {
			return
		}
		msgs := tx.GetMsgs()
		if len(msgs) != 1 {
			err = errors.New("must provide genesis StdTx with exactly 1 CreateValidator message")
			return
		}
		msg := msgs[0].(stake.MsgCreateValidator)

		// create the genesis account, give'm few iris token and a buncha token with there name
		genaccs[i] = genesisAccountFromMsgCreateValidator(msg, FreeFermionAcc.Amount)
		stakeData.Pool.LooseTokens = stakeData.Pool.LooseTokens.Add(sdk.NewDecFromInt(FreeFermionAcc.Amount)) // increase the supply
	}

	// create the final app state
	genesisState = GenesisState{
		Accounts:     genaccs,
		StakeData:    stakeData,
		MintData:     mint.GenesisState{
			Minter: mint.InitialMinter(),
			Params: mint.Params{
				MintDenom:           "iris",
				InflationRateChange: sdk.NewDecWithPrec(13, 2),
				InflationMax:        sdk.NewDecWithPrec(20, 2),
				InflationMin:        sdk.NewDecWithPrec(7, 2),
				GoalBonded:          sdk.NewDecWithPrec(67, 2),
			},
		},
		DistrData:    distr.DefaultGenesisState(),
		GovData:      gov.GenesisState{
			StartingProposalID: 1,
			DepositProcedure: govparams.DepositProcedure{
				MinDeposit:       sdk.Coins{sdk.NewInt64Coin("iris-atto", 10)},
				MaxDepositPeriod: time.Duration(172800) * time.Second,
			},
			VotingProcedure: govparams.VotingProcedure{
				VotingPeriod: time.Duration(172800) * time.Second,
			},
			TallyingProcedure: govparams.TallyingProcedure{
				Threshold:         sdk.NewDecWithPrec(5, 1),
				Veto:              sdk.NewDecWithPrec(334, 3),
				GovernancePenalty: sdk.NewDecWithPrec(1, 2),
			},
		},
		UpgradeData:  upgrade.DefaultGenesisState(),
		SlashingData: slashingData,
		GenTxs:       appGenTxs,
	}
	return
}

func genesisAccountFromMsgCreateValidator(msg stake.MsgCreateValidator, amount sdk.Int) GenesisAccount {
	accAuth := auth.NewBaseAccountWithAddress(sdk.AccAddress(msg.ValidatorAddr))
	accAuth.Coins = []sdk.Coin{
		{msg.Description.Moniker + "Token", sdk.NewInt(1000)},
		{"iris-atto", amount},
	}
	return NewGenesisAccount(&accAuth)
}

// IrisValidateGenesisState ensures that the genesis state obeys the expected invariants
// TODO: No validators are both bonded and jailed (#2088)
// TODO: Error if there is a duplicate validator (#1708)
// TODO: Ensure all state machine parameters are in genesis (#1704)
func IrisValidateGenesisState(genesisState GenesisState) (err error) {
	err = validateGenesisStateAccounts(genesisState.Accounts)
	if err != nil {
		return
	}
	// skip stakeData validation as genesis is created from txs
	if len(genesisState.GenTxs) > 0 {
		return nil
	}
	return stake.ValidateGenesis(genesisState.StakeData)
}

// Ensures that there are no duplicate accounts in the genesis state,
func validateGenesisStateAccounts(accs []GenesisAccount) (err error) {
	addrMap := make(map[string]bool, len(accs))
	for i := 0; i < len(accs); i++ {
		acc := accs[i]
		strAddr := string(acc.Address)
		if _, ok := addrMap[strAddr]; ok {
			return fmt.Errorf("Duplicate account in genesis state: Address %v", acc.Address)
		}
		addrMap[strAddr] = true
	}
	return
}

// IrisAppGenState but with JSON
func IrisAppGenStateJSON(cdc *codec.Codec, appGenTxs []json.RawMessage) (appState json.RawMessage, err error) {

	// create the final app state
	genesisState, err := IrisAppGenState(cdc, appGenTxs)
	if err != nil {
		return nil, err
	}
	appState, err = codec.MarshalJSONIndent(cdc, genesisState)
	return
}

// CollectStdTxs processes and validates application's genesis StdTxs and returns the list of validators,
// appGenTxs, and persistent peers required to generate genesis.json.
func CollectStdTxs(moniker string, genTxsDir string, cdc *codec.Codec) (
	validators []tmtypes.GenesisValidator, appGenTxs []auth.StdTx, persistentPeers string, err error) {
	var fos []os.FileInfo
	fos, err = ioutil.ReadDir(genTxsDir)
	if err != nil {
		return
	}

	var addresses []string
	for _, fo := range fos {
		filename := filepath.Join(genTxsDir, fo.Name())
		if !fo.IsDir() && (filepath.Ext(filename) != ".json") {
			continue
		}

		// get the genStdTx
		var jsonRawTx []byte
		jsonRawTx, err = ioutil.ReadFile(filename)
		if err != nil {
			return
		}
		var genStdTx auth.StdTx
		err = cdc.UnmarshalJSON(jsonRawTx, &genStdTx)
		if err != nil {
			return
		}
		appGenTxs = append(appGenTxs, genStdTx)

		nodeAddr := genStdTx.GetMemo()
		if len(nodeAddr) == 0 {
			err = fmt.Errorf("couldn't find node's address in %s", fo.Name())
			return
		}

		msgs := genStdTx.GetMsgs()
		if len(msgs) != 1 {
			err = errors.New("each genesis transaction must provide a single genesis message")
			return
		}

		// TODO: this could be decoupled from stake.MsgCreateValidator
		// TODO: and we likely want to do it for real world Gaia
		msg := msgs[0].(stake.MsgCreateValidator)
		validators = append(validators, tmtypes.GenesisValidator{
			PubKey: msg.PubKey,
			Power:  FreeFermionVal.Amount.Int64(),
			Name:   msg.Description.Moniker,
		})

		// exclude itself from persistent peers
		if msg.Description.Moniker != moniker {
			addresses = append(addresses, nodeAddr)
		}
	}

	sort.Strings(addresses)
	persistentPeers = strings.Join(addresses, ",")

	return
}

func NewDefaultGenesisAccount(addr sdk.AccAddress) GenesisAccount {
	accAuth := auth.NewBaseAccountWithAddress(addr)
	accAuth.Coins = []sdk.Coin{
		FreeFermionAcc,
	}
	return NewGenesisAccount(&accAuth)
}

func createGenesisState() stake.GenesisState {
	return stake.GenesisState{
		Pool: stake.Pool{
			LooseTokens:  sdk.ZeroDec(),
			BondedTokens: sdk.ZeroDec(),
		},
		Params: stake.Params{
			UnbondingTime: defaultUnbondingTime,
			MaxValidators: 100,
			BondDenom:     Denom + "-" + types.Atto,
		},
	}
}
