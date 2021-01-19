package migrate

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmjson "github.com/tendermint/tendermint/libs/json"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	coinswaptypes "github.com/irisnet/irismod/modules/coinswap/types"
	htlctypes "github.com/irisnet/irismod/modules/htlc/types"
	randomtypes "github.com/irisnet/irismod/modules/random/types"
	servicetypes "github.com/irisnet/irismod/modules/service/types"
	"github.com/irisnet/irismod/modules/token"
	tokentypes "github.com/irisnet/irismod/modules/token/types"

	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/migrate/v0_16"
	"github.com/irisnet/irishub/migrate/v0_16/auth"
	"github.com/irisnet/irishub/migrate/v0_16/types"
	guardiantypes "github.com/irisnet/irishub/modules/guardian/types"
	minttypes "github.com/irisnet/irishub/modules/mint/types"
)

const (
	flagGenesisTime = "genesis-time"
	flagChainID     = "chain-id"
	IRISATTO        = "iris-atto"
	UIRIS           = "uiris"
)

var Precision = sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(12), nil))

// MigrateGenesisCmd returns a command to execute genesis state migration.
func MigrateGenesisCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate [genesis-file]",
		Short: "Migrate genesis to a specified target version",
		Long: fmt.Sprintf(`Migrate the source genesis into the target version and print to STDOUT.

Example:
$ %s migrate /path/to/genesis.json --chain-id=cosmoshub-3 --genesis-time=2019-04-22T17:00:00Z
`, version.AppName),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			cdc := clientCtx.JSONMarshaler

			var err error

			importGenesis := args[0]

			genDoc, err := v0_16.GenesisDocFromFile(importGenesis)
			if err != nil {
				return errors.Wrapf(err, "failed to read genesis document from file %s", importGenesis)
			}

			var initialState v0_16.GenesisFileState

			if err := types.CodeC.UnmarshalJSON(genDoc.AppState, &initialState); err != nil {
				return errors.Wrap(err, "failed to JSON unmarshal initial genesis state")
			}

			newGenState := Migrate(cdc, initialState)

			// Set DefaultGenesis if this module does not exist in v0.16
			for _, b := range app.ModuleBasics {
				if _, ok := newGenState[b.Name()]; !ok {
					newGenState[b.Name()] = b.DefaultGenesis(cdc)
				}
			}

			genDoc.AppState, err = json.MarshalIndent(newGenState, "", "  ")
			if err != nil {
				return errors.Wrap(err, "failed to JSON marshal migrated genesis state")
			}

			genesisTime := cmd.Flag(flagGenesisTime).Value.String()
			if genesisTime != "" {
				var t time.Time

				err := t.UnmarshalText([]byte(genesisTime))
				if err != nil {
					return errors.Wrap(err, "failed to unmarshal genesis time")
				}

				genDoc.GenesisTime = t
			}

			chainID := cmd.Flag(flagChainID).Value.String()
			if chainID != "" {
				genDoc.ChainID = chainID
			}

			consensusParams := tmtypes.DefaultConsensusParams()
			consensusParams.Block.MaxBytes = genDoc.ConsensusParams.BlockSize.MaxBytes
			consensusParams.Block.MaxGas = genDoc.ConsensusParams.BlockSize.MaxGas

			consensusParams.Validator.PubKeyTypes = genDoc.ConsensusParams.Validator.PubKeyTypes

			newGenDoc := tmtypes.GenesisDoc{
				GenesisTime:     genDoc.GenesisTime,
				ChainID:         genDoc.ChainID,
				ConsensusParams: consensusParams,
				AppHash:         genDoc.AppHash,
				AppState:        genDoc.AppState,
			}

			bz, err := tmjson.MarshalIndent(newGenDoc, "", "  ")
			if err != nil {
				return errors.Wrap(err, "failed to marshal genesis doc")
			}

			fmt.Println(string(bz))
			return nil
		},
	}

	cmd.Flags().String(flagGenesisTime, "", "override genesis_time with this flag")
	cmd.Flags().String(flagChainID, "", "override chain_id with this flag")

	return cmd
}

func Migrate(cdc codec.JSONMarshaler, initialState v0_16.GenesisFileState) (appState map[string]json.RawMessage) {
	appState = make(map[string]json.RawMessage)

	// ------------------------------------------------------------
	// sdk modules
	// ------------------------------------------------------------
	stakingGenesis, bondedTokens, notBondedTokens := migrateStaking(initialState)
	authGenesisState, bankGenesisState, communityTax, tokenGenesisState := migrateAuth(initialState, bondedTokens, notBondedTokens)
	appState[stakingtypes.ModuleName] = cdc.MustMarshalJSON(stakingGenesis)
	appState[authtypes.ModuleName] = cdc.MustMarshalJSON(&authGenesisState)
	appState[banktypes.ModuleName] = cdc.MustMarshalJSON(&bankGenesisState)
	appState[slashingtypes.ModuleName] = cdc.MustMarshalJSON(migrateSlashing(initialState))
	appState[distributiontypes.ModuleName] = cdc.MustMarshalJSON(migrateDistribution(initialState, communityTax))
	appState[govtypes.ModuleName] = cdc.MustMarshalJSON(migrateGov(initialState))
	appState[crisistypes.ModuleName] = cdc.MustMarshalJSON(genCrisis())
	appState[tokentypes.ModuleName] = cdc.MustMarshalJSON(&tokenGenesisState)

	// ------------------------------------------------------------
	// irishub modules
	// ------------------------------------------------------------
	appState[minttypes.ModuleName] = cdc.MustMarshalJSON(migrateMint(initialState))
	appState[randomtypes.ModuleName] = cdc.MustMarshalJSON(migrateRand(initialState))
	appState[htlctypes.ModuleName] = cdc.MustMarshalJSON(migrateHTLC(initialState))
	appState[coinswaptypes.ModuleName] = cdc.MustMarshalJSON(migrateCoinswap(initialState))
	appState[guardiantypes.ModuleName] = cdc.MustMarshalJSON(migrateGuardian(initialState))
	appState[servicetypes.ModuleName] = cdc.MustMarshalJSON(migrateService(initialState))

	return appState

}

func migrateAuth(initialState v0_16.GenesisFileState, bondedTokens, notBondedTokens sdk.Coins) (authtypes.GenesisState,
	banktypes.GenesisState,
	sdk.Coins,
	tokentypes.GenesisState) {
	params := authtypes.DefaultParams()
	var accounts authtypes.GenesisAccounts
	var balances []banktypes.Balance
	var communityTax sdk.Coins
	var serviceTax sdk.Coins
	var burnedCoins sdk.Coins
	for _, acc := range initialState.Accounts {
		var coins sdk.Coins
		for _, c := range acc.Coins {
			coin, add := convertCoinStr(c)
			if add {
				coins = append(coins, coin)
			}
			coins = sdk.NewCoins(coins...)
		}
		var account authtypes.GenesisAccount
		baseAccount := authtypes.NewBaseAccount(acc.Address, nil, acc.AccountNumber, acc.Sequence)

		switch acc.Address.String() {
		case auth.BurnedCoinsAccAddr.String():
			burnedCoins = coins
			continue
		case auth.GovDepositCoinsAccAddr.String():
			baseAccount.Address = authtypes.NewModuleAddress(govtypes.ModuleName).String()
			account = authtypes.NewModuleAccount(baseAccount, govtypes.ModuleName, authtypes.Burner)
		case auth.ServiceDepositCoinsAccAddr.String():
			baseAccount.Address = authtypes.NewModuleAddress(servicetypes.DepositAccName).String()
			account = authtypes.NewModuleAccount(baseAccount, servicetypes.DepositAccName, authtypes.Burner)
		case auth.ServiceRequestCoinsAccAddr.String():
			baseAccount.Address = authtypes.NewModuleAddress(servicetypes.RequestAccName).String()
			account = authtypes.NewModuleAccount(baseAccount, servicetypes.RequestAccName)
		case auth.ServiceTaxCoinsAccAddr.String():
			serviceTax = coins
			account = baseAccount
			coins = sdk.NewCoins()
		case auth.CommunityTaxCoinsAccAddr.String():
			communityTax = coins
			baseAccount.Address = authtypes.NewModuleAddress(distributiontypes.ModuleName).String()
			account = authtypes.NewModuleAccount(baseAccount, distributiontypes.ModuleName)
		default:
			account = baseAccount
		}
		accounts = append(accounts, account)
		balances = append(balances, banktypes.Balance{Address: account.GetAddress().String(), Coins: coins})
	}

	// add bonded pool account
	bondedPoolAddress := authtypes.NewModuleAddress(stakingtypes.BondedPoolName)
	accounts = append(accounts, authtypes.NewModuleAccount(
		authtypes.NewBaseAccount(bondedPoolAddress, nil, uint64(len(initialState.Accounts)), 0),
		stakingtypes.BondedPoolName, authtypes.Burner, authtypes.Staking),
	)
	balances = append(balances, banktypes.Balance{Address: bondedPoolAddress.String(), Coins: bondedTokens})

	// add not bonded pool account
	notBondedPoolAddress := authtypes.NewModuleAddress(stakingtypes.NotBondedPoolName)
	accounts = append(accounts, authtypes.NewModuleAccount(
		authtypes.NewBaseAccount(notBondedPoolAddress, nil, uint64(len(initialState.Accounts)), 0),
		stakingtypes.NotBondedPoolName, authtypes.Burner, authtypes.Staking),
	)
	balances = append(balances, banktypes.Balance{Address: notBondedPoolAddress.String(), Coins: notBondedTokens})

	feeCollectorAddress := authtypes.NewModuleAddress(authtypes.FeeCollectorName)
	balances = append(balances, banktypes.Balance{Address: feeCollectorAddress.String(), Coins: serviceTax})

	authGenesisState := authtypes.NewGenesisState(
		params, accounts,
	)

	bankGenesisState := banktypes.GenesisState{
		Params:   banktypes.DefaultParams(),
		Balances: balances,
	}

	tokenGenesisState := token.DefaultGenesisState()
	tokenGenesisState.BurnedCoins = burnedCoins

	return *authGenesisState, bankGenesisState, communityTax, *tokenGenesisState
}

func migrateStaking(initialState v0_16.GenesisFileState) (*stakingtypes.GenesisState, sdk.Coins, sdk.Coins) {
	bondedTokens := sdk.NewCoins()
	notBondedTokens := sdk.NewCoins()
	params := stakingtypes.Params{
		UnbondingTime:     initialState.StakeData.Params.UnbondingTime,
		MaxValidators:     uint32(initialState.StakeData.Params.MaxValidators),
		MaxEntries:        stakingtypes.DefaultParams().MaxEntries,
		HistoricalEntries: stakingtypes.DefaultParams().HistoricalEntries,
		BondDenom:         UIRIS,
	}

	// lastTotalPower := initialState.StakeData.LastTotalPower
	// var lastValidatorPowers []stakingtypes.LastValidatorPower
	// for _, lvp := range initialState.StakeData.LastValidatorPowers {
	// 	lastValidatorPowers = append(
	// 		lastValidatorPowers,
	// 		stakingtypes.LastValidatorPower{
	// 			Address: lvp.Address.String(),
	// 			Power:   lvp.Power.Quo(Precision).Int64(),
	// 		},
	// 	)
	// }

	validatorTotalShares := make(map[string]sdk.Dec)
	var delegations stakingtypes.Delegations
	for _, b := range initialState.StakeData.Bonds {
		shares := b.Shares.Quo(sdk.NewDecFromInt(Precision))
		delegations = append(
			delegations,
			stakingtypes.Delegation{
				DelegatorAddress: b.DelegatorAddr.String(),
				ValidatorAddress: b.ValidatorAddr.String(),
				Shares:           shares,
			},
		)
		if _, ok := validatorTotalShares[b.ValidatorAddr.String()]; ok {
			validatorTotalShares[b.ValidatorAddr.String()] = validatorTotalShares[b.ValidatorAddr.String()].Add(shares)
		} else {
			validatorTotalShares[b.ValidatorAddr.String()] = shares
		}
	}

	var validators stakingtypes.Validators
	for _, v := range initialState.StakeData.Validators {
		status := stakingtypes.BondStatus(int32(v.Status) + 1)
		bondedToken := v.Tokens.TruncateInt().Quo(Precision)
		consPubKey, err := cryptocodec.FromTmPubKeyInterface(v.ConsPubKey)
		if err != nil {
			panic(err)
		}
		pkAny, err := codectypes.NewAnyWithValue(consPubKey)
		if err != nil {
			panic(err)
		}
		validators = append(
			validators,
			stakingtypes.Validator{
				OperatorAddress: v.OperatorAddr.String(),
				ConsensusPubkey: pkAny,
				Jailed:          v.Jailed,
				Status:          status,
				Tokens:          bondedToken,
				DelegatorShares: validatorTotalShares[v.OperatorAddr.String()],
				Description: stakingtypes.Description{
					Moniker:         v.Description.Moniker,
					Identity:        v.Description.Identity,
					Website:         v.Description.Website,
					SecurityContact: "",
					Details:         v.Description.Details,
				},
				UnbondingHeight: v.UnbondingHeight,
				UnbondingTime:   v.UnbondingMinTime,
				Commission: stakingtypes.Commission{
					CommissionRates: stakingtypes.CommissionRates{
						Rate:          v.Commission.Rate,
						MaxRate:       v.Commission.MaxRate,
						MaxChangeRate: v.Commission.MaxChangeRate,
					},
					UpdateTime: v.Commission.UpdateTime,
				},
				MinSelfDelegation: sdk.ZeroInt(),
			},
		)
		if status == stakingtypes.Bonded {
			bondedTokens = bondedTokens.Add(sdk.NewCoin(params.BondDenom, bondedToken))
		} else {
			notBondedTokens = notBondedTokens.Add(sdk.NewCoin(params.BondDenom, bondedToken))
		}
	}

	var unbondingDelegations []stakingtypes.UnbondingDelegation
	for _, b := range initialState.StakeData.UnbondingDelegations {
		balance := b.Balance.Amount.Quo(Precision)
		unbondingDelegations = append(
			unbondingDelegations,
			stakingtypes.UnbondingDelegation{
				DelegatorAddress: b.DelegatorAddr.String(),
				ValidatorAddress: b.ValidatorAddr.String(),
				Entries: []stakingtypes.UnbondingDelegationEntry{
					{
						CreationHeight: b.CreationHeight,
						CompletionTime: b.MinTime,
						InitialBalance: b.InitialBalance.Amount.Quo(Precision),
						Balance:        balance,
					},
				},
			},
		)
		notBondedTokens = notBondedTokens.Add(sdk.NewCoin(params.BondDenom, balance))
	}
	var redelegations []stakingtypes.Redelegation
	for _, r := range initialState.StakeData.Redelegations {
		redelegations = append(
			redelegations,
			stakingtypes.Redelegation{
				DelegatorAddress:    r.DelegatorAddr.String(),
				ValidatorSrcAddress: r.ValidatorSrcAddr.String(),
				ValidatorDstAddress: r.ValidatorDstAddr.String(),
				Entries: []stakingtypes.RedelegationEntry{
					{
						CreationHeight: r.CreationHeight,
						CompletionTime: r.MinTime,
						InitialBalance: r.InitialBalance.Amount.Quo(Precision),
						SharesDst:      r.SharesDst.Quo(sdk.NewDecFromInt(Precision)),
					},
				},
			},
		)
	}

	return &stakingtypes.GenesisState{
		Params: params,
		// LastTotalPower:       lastTotalPower,
		// LastValidatorPowers:  lastValidatorPowers,
		Validators:           validators,
		Delegations:          delegations,
		UnbondingDelegations: unbondingDelegations,
		Redelegations:        redelegations,
		Exported:             false,
	}, bondedTokens, notBondedTokens
}

type SigningInfoSlice []slashingtypes.SigningInfo

func (s SigningInfoSlice) Len() int {
	return len(s)
}
func (s SigningInfoSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SigningInfoSlice) Less(i, j int) bool {
	return s[j].Address < s[i].Address
}

func migrateSlashing(initialState v0_16.GenesisFileState) *slashingtypes.GenesisState {
	params := slashingtypes.Params{
		SignedBlocksWindow:      initialState.SlashingData.Params.SignedBlocksWindow,
		MinSignedPerWindow:      initialState.SlashingData.Params.MinSignedPerWindow,
		DowntimeJailDuration:    initialState.SlashingData.Params.DowntimeJailDuration,
		SlashFractionDoubleSign: initialState.SlashingData.Params.SlashFractionDoubleSign,
		SlashFractionDowntime:   initialState.SlashingData.Params.SlashFractionDowntime,
	}
	var validatorSigningInfos = make(map[string]slashingtypes.ValidatorSigningInfo)
	for ba, vs := range initialState.SlashingData.SigningInfos {
		acc, _ := sdk.ConsAddressFromBech32(ba)
		validatorSigningInfos[ba] = slashingtypes.ValidatorSigningInfo{
			Address:             acc.String(),
			StartHeight:         vs.StartHeight,
			IndexOffset:         vs.IndexOffset,
			JailedUntil:         vs.JailedUntil,
			Tombstoned:          false,
			MissedBlocksCounter: vs.MissedBlocksCounter,
		}
	}

	var signingInfos []slashingtypes.SigningInfo
	for addr, validatorSigningInfo := range validatorSigningInfos {
		signingInfos = append(signingInfos, slashingtypes.SigningInfo{
			Address:              addr,
			ValidatorSigningInfo: validatorSigningInfo,
		})
	}

	sort.Sort(SigningInfoSlice(signingInfos))

	return &slashingtypes.GenesisState{
		Params:       params,
		SigningInfos: signingInfos,
	}
}

func migrateDistribution(initialState v0_16.GenesisFileState, communityTax sdk.Coins) *distributiontypes.GenesisState {
	v016params := initialState.DistrData.Params
	params := distributiontypes.Params{
		CommunityTax:        v016params.CommunityTax,
		BaseProposerReward:  v016params.BaseProposerReward,
		BonusProposerReward: v016params.BonusProposerReward,
		WithdrawAddrEnabled: true,
	}
	feePool := distributiontypes.FeePool{CommunityPool: sdk.NewDecCoinsFromCoins(communityTax...)}

	var delegatorWithdrawInfos []distributiontypes.DelegatorWithdrawInfo
	previousProposer := initialState.DistrData.PreviousProposer
	var outstandingRewards []distributiontypes.ValidatorOutstandingRewardsRecord
	var validatorAccumulatedCommissions []distributiontypes.ValidatorAccumulatedCommissionRecord
	var validatorHistoricalRewards []distributiontypes.ValidatorHistoricalRewardsRecord
	var validatorCurrentRewards []distributiontypes.ValidatorCurrentRewardsRecord
	var delegatorStartingInfos []distributiontypes.DelegatorStartingInfoRecord
	var validatorSlashEvents []distributiontypes.ValidatorSlashEventRecord

	return &distributiontypes.GenesisState{
		Params:                          params,
		FeePool:                         feePool,
		DelegatorWithdrawInfos:          delegatorWithdrawInfos,
		PreviousProposer:                previousProposer.String(),
		OutstandingRewards:              outstandingRewards,
		ValidatorAccumulatedCommissions: validatorAccumulatedCommissions,
		ValidatorHistoricalRewards:      validatorHistoricalRewards,
		ValidatorCurrentRewards:         validatorCurrentRewards,
		DelegatorStartingInfos:          delegatorStartingInfos,
		ValidatorSlashEvents:            validatorSlashEvents,
	}
}

func migrateGov(initialState v0_16.GenesisFileState) *govtypes.GenesisState {
	var deposits govtypes.Deposits
	var votes govtypes.Votes
	var proposals govtypes.Proposals
	depositParams := govtypes.DepositParams{
		MinDeposit:       convertCoins(initialState.GovData.Params.NormalMinDeposit),
		MaxDepositPeriod: initialState.GovData.Params.NormalDepositPeriod,
	}
	votingParams := govtypes.VotingParams{
		VotingPeriod: initialState.GovData.Params.NormalVotingPeriod,
	}
	tallyParams := govtypes.TallyParams{
		Quorum:        initialState.GovData.Params.NormalParticipation,
		Threshold:     initialState.GovData.Params.NormalThreshold,
		VetoThreshold: initialState.GovData.Params.NormalVeto,
	}

	return &govtypes.GenesisState{
		StartingProposalId: govtypes.DefaultStartingProposalID,
		Deposits:           deposits,
		Votes:              votes,
		Proposals:          proposals,
		DepositParams:      depositParams,
		VotingParams:       votingParams,
		TallyParams:        tallyParams,
	}
}

func genCrisis() *crisistypes.GenesisState {
	return crisistypes.NewGenesisState(sdk.NewCoin(UIRIS, sdk.NewInt(1000)))
}

func migrateMint(initialState v0_16.GenesisFileState) *minttypes.GenesisState {
	minter := minttypes.Minter{
		LastUpdate:    initialState.MintData.Minter.LastUpdate,
		InflationBase: initialState.MintData.Minter.InflationBase.Quo(Precision),
	}
	params := minttypes.Params{
		Inflation: initialState.MintData.Params.Inflation,
		MintDenom: UIRIS,
	}

	return &minttypes.GenesisState{
		Minter: minter,
		Params: params,
	}
}

func migrateRand(initialState v0_16.GenesisFileState) *randomtypes.GenesisState {
	var pendingRandomRequests = make(map[string]randomtypes.Requests)
	for lh, rs := range initialState.RandData.PendingRandRequests {
		var requests []randomtypes.Request
		for _, r := range rs {
			requests = append(
				requests,
				randomtypes.Request{
					Height:   r.Height,
					Consumer: r.Consumer.String(),
					TxHash:   tmbytes.HexBytes(r.TxHash).String(),
				},
			)
		}
		pendingRandomRequests[lh] = randomtypes.Requests{Requests: requests}
	}

	return &randomtypes.GenesisState{
		PendingRandomRequests: pendingRandomRequests,
	}
}

func migrateHTLC(initialState v0_16.GenesisFileState) *htlctypes.GenesisState {
	var pendingHTLCs = make(map[string]htlctypes.HTLC)
	for hk, h := range initialState.HtlcData.PendingHTLCs {
		pendingHTLCs[hk] = htlctypes.NewHTLC(
			h.Sender,
			h.To,
			h.ReceiverOnOtherChain,
			convertCoins(h.Amount),
			h.Secret,
			h.Timestamp,
			h.ExpireHeight,
			htlctypes.HTLCState(h.State),
		)
	}

	return &htlctypes.GenesisState{
		PendingHtlcs: pendingHTLCs,
	}
}

func migrateCoinswap(initialState v0_16.GenesisFileState) *coinswaptypes.GenesisState {
	fee, _ := sdk.NewDecFromStr(initialState.SwapData.Params.Fee.FloatString(sdk.Precision))
	params := coinswaptypes.Params{
		Fee: fee,
	}

	return &coinswaptypes.GenesisState{
		Params:        params,
		StandardDenom: UIRIS,
	}
}

func migrateGuardian(initialState v0_16.GenesisFileState) *guardiantypes.GenesisState {
	var supers []guardiantypes.Super

	for _, profiler := range initialState.GuardianData.Profilers {
		accountType, err := guardiantypes.AccountTypeFromString(profiler.AccountType.String())
		if err != nil {
			panic(err.Error())
		}
		supers = append(supers, guardiantypes.Super{
			Description: profiler.Description,
			AccountType: accountType,
			Address:     profiler.Address.String(),
			AddedBy:     profiler.AddedBy.String(),
		})
	}

	return &guardiantypes.GenesisState{
		Supers: supers,
	}
}

func migrateService(initialState v0_16.GenesisFileState) *servicetypes.GenesisState {
	params := servicetypes.Params{
		MaxRequestTimeout:    initialState.ServiceData.Params.MaxRequestTimeout,
		MinDepositMultiple:   initialState.ServiceData.Params.MinDepositMultiple,
		MinDeposit:           convertCoins(initialState.ServiceData.Params.MinDeposit),
		ServiceFeeTax:        initialState.ServiceData.Params.ServiceFeeTax,
		SlashFraction:        initialState.ServiceData.Params.SlashFraction,
		ComplaintRetrospect:  initialState.ServiceData.Params.ComplaintRetrospect,
		ArbitrationTimeLimit: initialState.ServiceData.Params.ArbitrationTimeLimit,
		TxSizeLimit:          initialState.ServiceData.Params.TxSizeLimit,
		BaseDenom:            UIRIS,
	}

	return &servicetypes.GenesisState{
		Params: params,
	}
}

// ignore token that cannot be converted
func convertCoinStr(coinStr string) (sdk.Coin, bool) {
	c := strings.ReplaceAll(coinStr, IRISATTO, UIRIS)
	coin, err := sdk.ParseCoinNormalized(c)
	if err != nil {
		return sdk.Coin{}, false
	}
	return sdk.Coin{
		Denom:  coin.Denom,
		Amount: coin.Amount.Quo(Precision),
	}, true
}

func convertCoins(coins sdk.Coins) sdk.Coins {
	for i, c := range coins {
		if c.Denom == IRISATTO {
			coins[i].Amount = c.Amount.Quo(Precision)
			coins[i].Denom = UIRIS
		}

	}
	return coins
}
