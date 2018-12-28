package v0

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/arbitration"
	"github.com/irisnet/irishub/modules/auth"
	distr "github.com/irisnet/irishub/modules/distribution"
	"github.com/irisnet/irishub/modules/gov"
	"github.com/irisnet/irishub/modules/guardian"
	"github.com/irisnet/irishub/modules/mint"
	"github.com/irisnet/irishub/modules/service"
	"github.com/irisnet/irishub/modules/slashing"
	"github.com/irisnet/irishub/modules/stake"
	stakeTypes "github.com/irisnet/irishub/modules/stake/types"
	"github.com/irisnet/irishub/modules/upgrade"
	"github.com/irisnet/irishub/types"
	sdk "github.com/irisnet/irishub/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

var (
	FeeAmt            = int64(100)
	IrisCt            = types.NewDefaultCoinType(stakeTypes.StakeDenomName)
	FreeFermionVal, _ = IrisCt.ConvertToMinCoin(fmt.Sprintf("%d%s", FeeAmt, stakeTypes.StakeDenomName))
	FreeFermionAcc, _ = IrisCt.ConvertToMinCoin(fmt.Sprintf("%d%s", int64(150), stakeTypes.StakeDenomName))
)

const (
	defaultUnbondingTime time.Duration = 60 * 10 * time.Second
	// DefaultKeyPass contains the default key password for genesis transactions
	DefaultKeyPass = "1234567890"
)

// State to Unmarshal
type GenesisState struct {
	Accounts        []GenesisAccount         `json:"accounts"`
	AuthData        auth.GenesisState        `json:"auth"`
	StakeData       stake.GenesisState       `json:"stake"`
	MintData        mint.GenesisState        `json:"mint"`
	DistrData       distr.GenesisState       `json:"distr"`
	GovData         gov.GenesisState         `json:"gov"`
	UpgradeData     upgrade.GenesisState     `json:"upgrade"`
	SlashingData    slashing.GenesisState    `json:"slashing"`
	ServiceData     service.GenesisState     `json:"service"`
	ArbitrationData arbitration.GenesisState `json:"arbitration"`
	GuardianData    guardian.GenesisState    `json:"guardian"`
	GenTxs          []json.RawMessage        `json:"gentxs"`
}

func NewGenesisState(accounts []GenesisAccount, authData auth.GenesisState, stakeData stake.GenesisState, mintData mint.GenesisState,
	distrData distr.GenesisState, govData gov.GenesisState, upgradeData upgrade.GenesisState, serviceData service.GenesisState,
	arbitrationData arbitration.GenesisState, guardianData guardian.GenesisState, slashingData slashing.GenesisState) GenesisState {

	return GenesisState{
		Accounts:        accounts,
		AuthData:        authData,
		StakeData:       stakeData,
		MintData:        mintData,
		DistrData:       distrData,
		GovData:         govData,
		UpgradeData:     upgradeData,
		ServiceData:     serviceData,
		ArbitrationData: arbitrationData,
		GuardianData:    guardianData,
		SlashingData:    slashingData,
	}
}

// GenesisAccount doesn't need pubkey or sequence
type GenesisAccount struct {
	Address       sdk.AccAddress `json:"address"`
	Coins         sdk.Coins      `json:"coins"`
	Sequence      uint64         `json:"sequence_number"`
	AccountNumber uint64         `json:"account_number"`
}

func NewGenesisAccount(acc *auth.BaseAccount) GenesisAccount {
	return GenesisAccount{
		Address:       acc.Address,
		Coins:         acc.Coins,
		AccountNumber: acc.AccountNumber,
		Sequence:      acc.Sequence,
	}
}

func NewGenesisAccountI(acc auth.Account) GenesisAccount {
	return GenesisAccount{
		Address:       acc.GetAddress(),
		Coins:         acc.GetCoins(),
		AccountNumber: acc.GetAccountNumber(),
		Sequence:      acc.GetSequence(),
	}
}

// convert GenesisAccount to auth.BaseAccount
func (ga *GenesisAccount) ToAccount() (acc *auth.BaseAccount) {
	return &auth.BaseAccount{
		Address:       ga.Address,
		Coins:         ga.Coins.Sort(),
		AccountNumber: ga.AccountNumber,
		Sequence:      ga.Sequence,
	}
}

// Create the core parameters for genesis initialization for iris
// note that the pubkey input is this machines pubkey
func IrisAppGenState(cdc *codec.Codec, genDoc tmtypes.GenesisDoc, appGenTxs []json.RawMessage) (
	genesisState GenesisFileState, err error) {
	if err = cdc.UnmarshalJSON(genDoc.AppState, &genesisState); err != nil {
		return genesisState, err
	}

	// if there are no gen txs to be processed, return the default empty state
	if len(appGenTxs) == 0 {
		return genesisState, errors.New("there must be at least one genesis tx")
	}

	stakeData := genesisState.StakeData
	for i, genTx := range appGenTxs {
		var tx auth.StdTx
		if err := cdc.UnmarshalJSON(genTx, &tx); err != nil {
			return genesisState, err
		}
		msgs := tx.GetMsgs()
		if len(msgs) != 1 {
			return genesisState, errors.New(
				"must provide genesis StdTx with exactly 1 CreateValidator message")
		}
		if _, ok := msgs[0].(stake.MsgCreateValidator); !ok {
			return genesisState, fmt.Errorf(
				"Genesis transaction %v does not contain a MsgCreateValidator", i)
		}
	}

	for _, acc := range genesisState.Accounts {
		// create the genesis account, give'm few stake token and a buncha token with there name
		for _, coin := range acc.Coins {
			coinName, err := types.GetCoinName(coin)
			if err != nil {
				panic(fmt.Sprintf("fatal error: failed pick out demon from coin: %s", coin))
			}
			if coinName != stakeTypes.StakeDenomName {
				continue
			}
			stakeToken, err := IrisCt.ConvertToMinCoin(coin)
			if err != nil {
				panic(fmt.Sprintf("fatal error: failed to convert %s to stake token: %s", stakeTypes.StakeDenom, coin))
			}
			stakeData.Pool.LooseTokens = stakeData.Pool.LooseTokens.
				Add(sdk.NewDecFromInt(stakeToken.Amount)) // increase the supply
		}
	}
	genesisState.StakeData = stakeData
	genesisState.GenTxs = appGenTxs
	genesisState.UpgradeData = genesisState.UpgradeData
	return genesisState, nil
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
func IrisAppGenStateJSON(cdc *codec.Codec, genDoc tmtypes.GenesisDoc, appGenTxs []json.RawMessage) (
	appState json.RawMessage, err error) {

	// create the final app state
	genesisState, err := IrisAppGenState(cdc, genDoc, appGenTxs)
	if err != nil {
		return nil, err
	}
	appState, err = codec.MarshalJSONIndent(cdc, genesisState)
	return
}

// CollectStdTxs processes and validates application's genesis StdTxs and returns
// the list of appGenTxs, and persistent peers required to generate genesis.json.
func CollectStdTxs(cdc *codec.Codec, moniker string, genTxsDir string, genDoc tmtypes.GenesisDoc) (
	appGenTxs []auth.StdTx, persistentPeers string, err error) {

	var fos []os.FileInfo
	fos, err = ioutil.ReadDir(genTxsDir)
	if err != nil {
		return appGenTxs, persistentPeers, err
	}

	// prepare a map of all accounts in genesis state to then validate
	// against the validators addresses
	var appFileState GenesisFileState
	if err := cdc.UnmarshalJSON(genDoc.AppState, &appFileState); err != nil {
		return appGenTxs, persistentPeers, err
	}
	appState := convertToGenesisState(appFileState)
	addrMap := make(map[string]GenesisAccount, len(appState.Accounts))
	for i := 0; i < len(appState.Accounts); i++ {
		acc := appState.Accounts[i]
		strAddr := acc.Address.String()
		addrMap[strAddr] = acc
	}

	// addresses and IPs (and port) validator server info
	var addressesIPs []string

	for _, fo := range fos {
		filename := filepath.Join(genTxsDir, fo.Name())
		if !fo.IsDir() && (filepath.Ext(filename) != ".json") {
			continue
		}

		// get the genStdTx
		var jsonRawTx []byte
		if jsonRawTx, err = ioutil.ReadFile(filename); err != nil {
			return appGenTxs, persistentPeers, err
		}
		var genStdTx auth.StdTx
		if err = cdc.UnmarshalJSON(jsonRawTx, &genStdTx); err != nil {
			return appGenTxs, persistentPeers, err
		}
		appGenTxs = append(appGenTxs, genStdTx)

		// the memo flag is used to store
		// the ip and node-id, for example this may be:
		// "528fd3df22b31f4969b05652bfe8f0fe921321d5@192.168.2.37:26656"
		nodeAddrIP := genStdTx.GetMemo()
		if len(nodeAddrIP) == 0 {
			return appGenTxs, persistentPeers, fmt.Errorf(
				"couldn't find node's address and IP in %s", fo.Name())
		}

		// genesis transactions must be single-message
		msgs := genStdTx.GetMsgs()
		if len(msgs) != 1 {

			return appGenTxs, persistentPeers, errors.New(
				"each genesis transaction must provide a single genesis message")
		}

		// validate the validator address and funds against the accounts in the state
		msg := msgs[0].(stake.MsgCreateValidator)
		addr := sdk.AccAddress(msg.ValidatorAddr).String()
		acc, ok := addrMap[addr]
		if !ok {
			return appGenTxs, persistentPeers, fmt.Errorf(
				"account %v in gentx(%s) is not exist in genesis.json accounts: %+v", addr, fo.Name(), addrMap)
		}
		if acc.Coins.AmountOf(msg.Delegation.Denom).LT(msg.Delegation.Amount) {
			err = fmt.Errorf("insufficient fund for the delegation: %v < %v",
				acc.Coins.AmountOf(msg.Delegation.Denom), msg.Delegation.Amount)
		}

		// exclude itself from persistent peers
		if msg.Description.Moniker != moniker {
			addressesIPs = append(addressesIPs, nodeAddrIP)
		}
	}

	sort.Strings(addressesIPs)
	persistentPeers = strings.Join(addressesIPs, ",")

	return appGenTxs, persistentPeers, nil
}

func createStakeGenesisState() stake.GenesisState {
	return stake.GenesisState{
		Pool: stake.Pool{
			LooseTokens:  sdk.ZeroDec(),
			BondedTokens: sdk.ZeroDec(),
		},
		Params: stake.Params{
			UnbondingTime: defaultUnbondingTime,
			MaxValidators: 100,
			BondDenom:     stakeTypes.StakeDenom,
		},
	}
}

func createMintGenesisState() mint.GenesisState {
	return mint.GenesisState{
		Minter: mint.InitialMinter(),
		Params: mint.Params{
			MintDenom:           stakeTypes.StakeDenom,
			InflationRateChange: sdk.NewDecWithPrec(13, 2),
			InflationMax:        sdk.NewDecWithPrec(20, 2),
			InflationMin:        sdk.NewDecWithPrec(7, 2),
			GoalBonded:          sdk.NewDecWithPrec(67, 2),
		},
	}
}

// normalize stake token to mini-unit
func normalizeNativeToken(coins []string) sdk.Coins {
	var accountCoins sdk.Coins
	nativeCoin := sdk.NewInt64Coin(stakeTypes.StakeDenom, 0)
	for _, coin := range coins {
		coinName, err := types.GetCoinName(coin)
		if err != nil {
			panic(fmt.Sprintf("fatal error: failed pick out demon from coin: %s", coin))
		}
		if coinName == stakeTypes.StakeDenomName {
			normalizeNativeToken, err := IrisCt.ConvertToMinCoin(coin)
			if err != nil {
				panic(fmt.Sprintf("fatal error in converting %s to %s", coin, stakeTypes.StakeDenom))
			}
			nativeCoin = nativeCoin.Plus(normalizeNativeToken)
		} else {
			// not native token
			denom, amount, err := types.GetCoin(coin)
			if err != nil {
				panic(fmt.Sprintf("fatal error: genesis file contains invalid coin: %s", coin))
			}

			amt, ok := sdk.NewIntFromString(amount)
			if !ok {
				panic(fmt.Sprintf("non-native coin(%s) amount should be integer ", coin))
			}
			denom = strings.ToLower(denom)
			accountCoins = append(accountCoins, sdk.NewCoin(denom, amt))
		}
	}
	accountCoins = append(accountCoins, nativeCoin)
	if accountCoins.IsZero() {
		panic("invalid genesis file, found account without any token")
	}
	return accountCoins
}

func convertToGenesisState(genesisFileState GenesisFileState) GenesisState {
	var genesisAccounts []GenesisAccount
	for _, gacc := range genesisFileState.Accounts {
		acc := GenesisAccount{
			Address:       gacc.Address,
			Coins:         normalizeNativeToken(gacc.Coins),
			AccountNumber: gacc.AccountNumber,
			Sequence:      gacc.Sequence,
		}
		genesisAccounts = append(genesisAccounts, acc)
	}
	return GenesisState{
		Accounts:        genesisAccounts,
		AuthData:        genesisFileState.AuthData,
		StakeData:       genesisFileState.StakeData,
		MintData:        genesisFileState.MintData,
		DistrData:       genesisFileState.DistrData,
		GovData:         genesisFileState.GovData,
		UpgradeData:     genesisFileState.UpgradeData,
		SlashingData:    genesisFileState.SlashingData,
		ServiceData:     genesisFileState.ServiceData,
		ArbitrationData: genesisFileState.ArbitrationData,
		GuardianData:    genesisFileState.GuardianData,
		GenTxs:          genesisFileState.GenTxs,
	}
}

type GenesisFileState struct {
	Accounts        []GenesisFileAccount     `json:"accounts"`
	AuthData        auth.GenesisState        `json:"auth"`
	StakeData       stake.GenesisState       `json:"stake"`
	MintData        mint.GenesisState        `json:"mint"`
	DistrData       distr.GenesisState       `json:"distr"`
	GovData         gov.GenesisState         `json:"gov"`
	UpgradeData     upgrade.GenesisState     `json:"upgrade"`
	SlashingData    slashing.GenesisState    `json:"slashing"`
	ServiceData     service.GenesisState     `json:"service"`
	GuardianData    guardian.GenesisState    `json:"guardian"`
	ArbitrationData arbitration.GenesisState `json:"arbitration"`
	GenTxs          []json.RawMessage        `json:"gentxs"`
}

type GenesisFileAccount struct {
	Address       sdk.AccAddress `json:"address"`
	Coins         []string       `json:"coins"`
	Sequence      uint64         `json:"sequence_number"`
	AccountNumber uint64         `json:"account_number"`
}

func NewGenesisFileAccount(acc *auth.BaseAccount) GenesisFileAccount {
	var coins []string
	for _, coin := range acc.Coins {
		coins = append(coins, coin.String())
	}
	return GenesisFileAccount{
		Address:       acc.Address,
		Coins:         coins,
		AccountNumber: acc.AccountNumber,
		Sequence:      acc.Sequence,
	}
}

func NewGenesisFileState(accounts []GenesisFileAccount, authData auth.GenesisState, stakeData stake.GenesisState, mintData mint.GenesisState,
	distrData distr.GenesisState, govData gov.GenesisState, upgradeData upgrade.GenesisState, serviceData service.GenesisState,
	arbitrationData arbitration.GenesisState, guardianData guardian.GenesisState, slashingData slashing.GenesisState) GenesisFileState {

	return GenesisFileState{
		Accounts:        accounts,
		AuthData:        authData,
		StakeData:       stakeData,
		MintData:        mintData,
		DistrData:       distrData,
		GovData:         govData,
		UpgradeData:     upgradeData,
		ServiceData:     serviceData,
		ArbitrationData: arbitrationData,
		GuardianData:    guardianData,
		SlashingData:    slashingData,
	}
}

// NewDefaultGenesisState generates the default state for iris.
func NewDefaultGenesisFileState() GenesisFileState {
	return GenesisFileState{
		Accounts:        nil,
		StakeData:       createStakeGenesisState(),
		MintData:        createMintGenesisState(),
		DistrData:       distr.DefaultGenesisState(),
		GovData:         gov.DefaultGenesisState(),
		UpgradeData:     upgrade.DefaultGenesisState(),
		ServiceData:     service.DefaultGenesisState(),
		GuardianData:    guardian.DefaultGenesisState(),
		ArbitrationData: arbitration.DefaultGenesisState(),
		SlashingData:    slashing.DefaultGenesisState(),
		GenTxs:          nil,
	}
}

func NewDefaultGenesisFileAccount(addr sdk.AccAddress) GenesisFileAccount {
	accAuth := auth.NewBaseAccountWithAddress(addr)
	accAuth.Coins = []sdk.Coin{
		FreeFermionAcc,
	}
	return NewGenesisFileAccount(&accAuth)
}
