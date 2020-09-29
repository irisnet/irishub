package migrate

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	tmjson "github.com/tendermint/tendermint/libs/json"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	coinswaptypes "github.com/irisnet/irismod/modules/coinswap/types"
	htlctypes "github.com/irisnet/irismod/modules/htlc/types"
	servicetypes "github.com/irisnet/irismod/modules/service/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/irisnet/irishub/migrate/v0_16"
	"github.com/irisnet/irishub/migrate/v0_16/stake"
	"github.com/irisnet/irishub/migrate/v0_16/types"
	guardiantypes "github.com/irisnet/irishub/modules/guardian/types"
	minttypes "github.com/irisnet/irishub/modules/mint/types"
	randomtypes "github.com/irisnet/irismod/modules/random/types"
)

const (
	flagGenesisTime = "genesis-time"
	flagChainID     = "chain-id"
	IRISATTO        = "iris-atto"
	UIRIS           = "uiris"
	Precision       = 12
)

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

			consensusParams.Evidence.MaxNum = uint32(genDoc.ConsensusParams.Evidence.MaxAge)

			consensusParams.Validator.PubKeyTypes = genDoc.ConsensusParams.Validator.PubKeyTypes

			newGenDoc := tmtypes.GenesisDoc{
				GenesisTime:     genDoc.GenesisTime,
				ChainID:         genDoc.ChainID,
				ConsensusParams: consensusParams,
				Validators:      genDoc.Validators,
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
	authGenesisState, bankGenesisState := migrateAuth(initialState)
	appState[authtypes.ModuleName] = cdc.MustMarshalJSON(&authGenesisState)
	appState[banktypes.ModuleName] = cdc.MustMarshalJSON(&bankGenesisState)
	appState[stakingtypes.ModuleName] = cdc.MustMarshalJSON(migrateStaking(initialState))
	appState[slashingtypes.ModuleName] = cdc.MustMarshalJSON(migrateSlashing(initialState))
	appState[distributiontypes.ModuleName] = cdc.MustMarshalJSON(migrateDistribution(initialState))
	appState[govtypes.ModuleName] = cdc.MustMarshalJSON(migrateGov(initialState))

	// ------------------------------------------------------------
	// irishub modules
	// ------------------------------------------------------------
	appState[minttypes.ModuleName] = cdc.MustMarshalJSON(migrateMint(initialState))
	appState[randomtypes.ModuleName] = cdc.MustMarshalJSON(migrateRand(initialState))
	appState[htlctypes.ModuleName] = cdc.MustMarshalJSON(migrateHTLC(initialState))
	//appState[tokentypes.ModuleName] = cdc.MustMarshalJSON(migrateToken(initialState))
	appState[coinswaptypes.ModuleName] = cdc.MustMarshalJSON(migrateCoinswap(initialState))
	appState[guardiantypes.ModuleName] = cdc.MustMarshalJSON(migrateGuardian(initialState)) // TODO
	appState[servicetypes.ModuleName] = cdc.MustMarshalJSON(migrateService(initialState))   // TODO

	return appState

}

func migrateAuth(initialState v0_16.GenesisFileState) (authtypes.GenesisState, banktypes.GenesisState) {
	params := authtypes.DefaultParams()
	var accounts authtypes.GenesisAccounts
	var balances []banktypes.Balance
	for _, acc := range initialState.Accounts {
		var coins sdk.Coins
		for _, c := range acc.Coins {
			coinStr := strings.ReplaceAll(c, IRISATTO, UIRIS)
			coin, err := sdk.ParseCoin(coinStr)
			if err != nil {
				panic(err)
			}
			coins = append(coins, convertCoin(coin))
		}
		baseAccount := authtypes.NewBaseAccount(acc.Address, nil, acc.AccountNumber, acc.Sequence)
		accounts = append(accounts, baseAccount)
		balances = append(balances, banktypes.Balance{Address: acc.Address, Coins: coins})
	}

	authGenesisState := authtypes.NewGenesisState(
		params, accounts,
	)

	bankGenesisState := banktypes.GenesisState{
		Params:   banktypes.DefaultParams(),
		Balances: balances,
	}

	return *authGenesisState, bankGenesisState
}

func migrateStaking(initialState v0_16.GenesisFileState) *stakingtypes.GenesisState {
	params := stakingtypes.Params{
		UnbondingTime:     initialState.StakeData.Params.UnbondingTime,
		MaxValidators:     uint32(initialState.StakeData.Params.MaxValidators),
		MaxEntries:        stakingtypes.DefaultParams().MaxEntries,
		HistoricalEntries: stakingtypes.DefaultParams().HistoricalEntries,
		BondDenom:         sdk.DefaultBondDenom,
	}
	lastTotalPower := initialState.StakeData.LastTotalPower
	var lastValidatorPowers []stakingtypes.LastValidatorPower
	for _, lvp := range initialState.StakeData.LastValidatorPowers {
		lastValidatorPowers = append(
			lastValidatorPowers,
			stakingtypes.LastValidatorPower{
				Address: lvp.Address,
				Power:   lvp.Power.Quo(sdk.NewInt(sdk.Precision)).Int64(),
			},
		)
	}
	var validators stakingtypes.Validators
	for _, v := range initialState.StakeData.Validators {
		var status sdk.BondStatus
		switch stake.BondStatusToString(v.Status) {
		case sdk.BondStatusBonded:
			status = sdk.Bonded
		case sdk.BondStatusUnbonding:
			status = sdk.Unbonding
		case sdk.BondStatusUnbonded:
			status = sdk.Unbonded
		default:
			panic("unknown bond status")
		}
		validators = append(
			validators,
			stakingtypes.Validator{
				OperatorAddress: v.OperatorAddr,
				ConsensusPubkey: sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, v.ConsPubKey),
				Jailed:          v.Jailed,
				Status:          status,
				Tokens:          sdk.NewIntFromBigInt(v.Tokens.BigInt()),
				DelegatorShares: v.DelegatorShares,
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
				MinSelfDelegation: sdk.ZeroInt(), // TODO
			},
		)
	}
	var delegations stakingtypes.Delegations
	for _, b := range initialState.StakeData.Bonds {
		delegations = append(
			delegations,
			stakingtypes.Delegation{
				DelegatorAddress: b.DelegatorAddr,
				ValidatorAddress: b.ValidatorAddr,
				Shares:           b.Shares,
			},
		)
	}
	var unbondingDelegations []stakingtypes.UnbondingDelegation
	for _, b := range initialState.StakeData.UnbondingDelegations {
		unbondingDelegations = append(
			unbondingDelegations,
			stakingtypes.UnbondingDelegation{
				DelegatorAddress: b.DelegatorAddr,
				ValidatorAddress: b.ValidatorAddr,
				Entries: []stakingtypes.UnbondingDelegationEntry{
					{
						CreationHeight: b.CreationHeight,
						CompletionTime: b.MinTime,
						InitialBalance: b.InitialBalance.Amount,
						Balance:        b.Balance.Amount,
					},
				},
			},
		)
	}
	var redelegations []stakingtypes.Redelegation
	for _, r := range initialState.StakeData.Redelegations {
		redelegations = append(
			redelegations,
			stakingtypes.Redelegation{
				DelegatorAddress:    r.DelegatorAddr,
				ValidatorSrcAddress: r.ValidatorSrcAddr,
				ValidatorDstAddress: r.ValidatorDstAddr,
				Entries: []stakingtypes.RedelegationEntry{
					{
						CreationHeight: r.CreationHeight,
						CompletionTime: r.MinTime,
						InitialBalance: r.InitialBalance.Amount,
						SharesDst:      r.SharesDst,
					},
				},
			},
		)
	}
	exported := initialState.StakeData.Exported

	return &stakingtypes.GenesisState{
		Params:               params,
		LastTotalPower:       lastTotalPower,
		LastValidatorPowers:  lastValidatorPowers,
		Validators:           validators,
		Delegations:          delegations,
		UnbondingDelegations: unbondingDelegations,
		Redelegations:        redelegations,
		Exported:             exported,
	}
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
			Address:             acc,
			StartHeight:         vs.StartHeight,
			IndexOffset:         vs.IndexOffset,
			JailedUntil:         vs.JailedUntil,
			Tombstoned:          false, // TODO
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

	var mMissedBlocks = make(map[string][]slashingtypes.MissedBlock)
	for ba, mbs := range initialState.SlashingData.MissedBlocks {
		var missedBlocks []slashingtypes.MissedBlock
		for _, mb := range mbs {
			missedBlocks = append(
				missedBlocks,
				slashingtypes.MissedBlock{
					Index:  mb.Index,
					Missed: mb.Missed,
				},
			)
		}
		mMissedBlocks[ba] = missedBlocks
	}

	var missedBlocks []slashingtypes.ValidatorMissedBlocks
	for addr, missedBlock := range mMissedBlocks {
		missedBlocks = append(missedBlocks, slashingtypes.ValidatorMissedBlocks{
			Address:      addr,
			MissedBlocks: missedBlock,
		})
	}

	return &slashingtypes.GenesisState{
		Params:       params,
		SigningInfos: signingInfos,
		MissedBlocks: missedBlocks,
	}
}

func migrateDistribution(initialState v0_16.GenesisFileState) *distributiontypes.GenesisState {
	v016params := initialState.DistrData.Params
	params := distributiontypes.Params{
		CommunityTax:        v016params.CommunityTax,
		BaseProposerReward:  v016params.BaseProposerReward,
		BonusProposerReward: v016params.BonusProposerReward,
		WithdrawAddrEnabled: true,
	}
	feePool := distributiontypes.FeePool{CommunityPool: initialState.DistrData.FeePool.CommunityPool}

	var delegatorWithdrawInfos []distributiontypes.DelegatorWithdrawInfo

	for _, withdrawInfo := range initialState.DistrData.DelegatorWithdrawInfos {
		delegatorWithdrawInfos = append(delegatorWithdrawInfos, distributiontypes.DelegatorWithdrawInfo{
			DelegatorAddress: withdrawInfo.DelegatorAddr,
			WithdrawAddress:  withdrawInfo.WithdrawAddr,
		})
	}
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
		PreviousProposer:                previousProposer,
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
		MinDeposit:       initialState.GovData.Params.NormalMinDeposit,
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

func migrateMint(initialState v0_16.GenesisFileState) *minttypes.GenesisState {
	minter := minttypes.Minter{
		LastUpdate:    initialState.MintData.Minter.LastUpdate,
		InflationBase: initialState.MintData.Minter.InflationBase,
	}
	params := minttypes.Params{
		Inflation: initialState.MintData.Params.Inflation,
		MintDenom: sdk.DefaultBondDenom,
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
					Consumer: r.Consumer,
					TxHash:   r.TxHash,
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
			h.Amount,
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
		Fee:           fee,
		StandardDenom: sdk.DefaultBondDenom,
	}

	return &coinswaptypes.GenesisState{
		Params: params,
	}
}

func migrateGuardian(initialState v0_16.GenesisFileState) *guardiantypes.GenesisState {
	var profilers guardiantypes.Profilers
	var trustees guardiantypes.Trustees

	for _, profiler := range initialState.GuardianData.Profilers {
		accountType, err := guardiantypes.AccountTypeFromString(profiler.AccountType.String())
		if err != nil {
			panic(err.Error())
		}
		profilers = append(profilers, guardiantypes.Guardian{
			Description: profiler.Description,
			AccountType: accountType,
			Address:     profiler.Address,
			AddedBy:     profiler.AddedBy,
		})
	}

	return &guardiantypes.GenesisState{
		Profilers: profilers,
		Trustees:  trustees,
	}
}

func migrateService(initialState v0_16.GenesisFileState) *servicetypes.GenesisState {
	params := servicetypes.Params{
		MaxRequestTimeout:    initialState.ServiceData.Params.MaxRequestTimeout,
		MinDepositMultiple:   initialState.ServiceData.Params.MinDepositMultiple,
		MinDeposit:           initialState.ServiceData.Params.MinDeposit,
		ServiceFeeTax:        initialState.ServiceData.Params.ServiceFeeTax,
		SlashFraction:        initialState.ServiceData.Params.SlashFraction,
		ComplaintRetrospect:  initialState.ServiceData.Params.ComplaintRetrospect,
		ArbitrationTimeLimit: initialState.ServiceData.Params.ArbitrationTimeLimit,
		TxSizeLimit:          initialState.ServiceData.Params.TxSizeLimit,
		BaseDenom:            servicetypes.DefaultBaseDenom,
	}

	return &servicetypes.GenesisState{
		Params: params,
	}
}

func convertCoin(coin sdk.Coin) sdk.Coin {
	return sdk.Coin{
		Denom:  coin.Denom,
		Amount: coin.Amount.Quo(sdk.NewInt(Precision)),
	}
}
