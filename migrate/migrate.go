package migrate

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"

	v0_17 "github.com/irisnet/irishub/migrate/v0_17"
	"github.com/irisnet/irishub/modules/asset"
	token "github.com/irisnet/irishub/modules/asset/01-token"
	"github.com/irisnet/irishub/modules/coinswap"
	"github.com/irisnet/irishub/modules/guardian"
	"github.com/irisnet/irishub/modules/htlc"
	"github.com/irisnet/irishub/modules/mint"
	"github.com/irisnet/irishub/modules/rand"
	"github.com/irisnet/irishub/modules/service"
)

const (
	flagGenesisTime = "genesis-time"
	flagChainID     = "chain-id"
)

// MigrateGenesisCmd returns a command to execute genesis state migration.
// nolint: funlen
func MigrateGenesisCmd(_ *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate [target-version] [genesis-file]",
		Short: "Migrate genesis to a specified target version",
		Long: fmt.Sprintf(`Migrate the source genesis into the target version and print to STDOUT.

Example:
$ %s migrate /path/to/genesis.json --chain-id=cosmoshub-3 --genesis-time=2019-04-22T17:00:00Z
`, version.ServerName),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			importGenesis := args[0]

			genDoc, err := types.GenesisDocFromFile(importGenesis)
			if err != nil {
				return errors.Wrapf(err, "failed to read genesis document from file %s", importGenesis)
			}

			var initialState v0_17.GenesisFileState
			if err := cdc.UnmarshalJSON(genDoc.AppState, &initialState); err != nil {
				return errors.Wrap(err, "failed to JSON unmarshal initial genesis state")
			}

			newGenState := Migrate(cdc, initialState)

			genDoc.AppState, err = cdc.MarshalJSON(newGenState)
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

			bz, err := cdc.MarshalJSONIndent(genDoc, "", "  ")
			if err != nil {
				return errors.Wrap(err, "failed to marshal genesis doc")
			}

			sortedBz, err := sdk.SortJSON(bz)
			if err != nil {
				return errors.Wrap(err, "failed to sort JSON genesis doc")
			}

			fmt.Println(string(sortedBz))
			return nil
		},
	}

	cmd.Flags().String(flagGenesisTime, "", "override genesis_time with this flag")
	cmd.Flags().String(flagChainID, "", "override chain_id with this flag")

	return cmd
}

func Migrate(cdc *codec.Codec, initialState v0_17.GenesisFileState) (appState genutil.AppMap) {
	v016Codec := codec.New()
	codec.RegisterCrypto(v016Codec)

	appState[auth.ModuleName] = cdc.MustMarshalJSON(migrateAuth(initialState))
	appState[distribution.ModuleName] = cdc.MustMarshalJSON(migrateDistribution(initialState))
	// appState[gov.ModuleName] = cdc.MustMarshalJSON(migrateGov(initialState))
	appState[slashing.ModuleName] = cdc.MustMarshalJSON(migrateSlashing(initialState))
	appState[staking.ModuleName] = cdc.MustMarshalJSON(migrateStaking(initialState))

	appState[rand.ModuleName] = cdc.MustMarshalJSON(migrateRand(initialState))
	appState[htlc.ModuleName] = cdc.MustMarshalJSON(migrateHTLC(initialState))
	appState[mint.ModuleName] = cdc.MustMarshalJSON(migrateMint(initialState))
	// appState[guardian.ModuleName] = cdc.MustMarshalJSON(migrateGuardian(initialState))
	// appState[coinswap.ModuleName] = cdc.MustMarshalJSON(migrateCoinswap(initialState))
	// appState[asset.ModuleName] = cdc.MustMarshalJSON(migrateAsset(initialState))
	// appState[service.ModuleName] = cdc.MustMarshalJSON(migrateService(initialState))

	return appState
}

func migrateAuth(initialState v0_17.GenesisFileState) auth.GenesisState {
	authParams := auth.DefaultParams()
	var accounts exported.GenesisAccounts
	for _, acc := range initialState.Accounts {
		var coins sdk.Coins
		for _, c := range acc.Coins {
			coin, err := sdk.ParseCoin(c)
			if err != nil {
				panic(err)
			}
			coins = append(coins, coin)
		}
		baseAccount := auth.NewBaseAccount(acc.Address, coins.Sort(), nil, acc.AccountNumber, acc.Sequence)
		accounts = append(accounts, baseAccount)
	}
	return auth.GenesisState{
		Params:   authParams,
		Accounts: accounts,
	}
}

func migrateDistribution(initialState v0_17.GenesisFileState) distribution.GenesisState {
	var feePool distribution.FeePool
	var communityTax sdk.Dec
	var baseProposerReward sdk.Dec
	var bonusProposerReward sdk.Dec
	var withdrawAddrEnabled bool
	var delegatorWithdrawInfos []distribution.DelegatorWithdrawInfo
	var previousProposer sdk.ConsAddress
	var outstandingRewards []distribution.ValidatorOutstandingRewardsRecord
	var validatorAccumulatedCommissions []distribution.ValidatorAccumulatedCommissionRecord
	var validatorHistoricalRewards []distribution.ValidatorHistoricalRewardsRecord
	var validatorCurrentRewards []distribution.ValidatorCurrentRewardsRecord
	var delegatorStartingInfos []distribution.DelegatorStartingInfoRecord
	var validatorSlashEvents []distribution.ValidatorSlashEventRecord
	// TODO

	return distribution.GenesisState{
		FeePool:                         feePool,
		CommunityTax:                    communityTax,
		BaseProposerReward:              baseProposerReward,
		BonusProposerReward:             bonusProposerReward,
		WithdrawAddrEnabled:             withdrawAddrEnabled,
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

func migrateGov(initialState v0_17.GenesisFileState) gov.GenesisState {
	var startingProposalID uint64
	var deposits gov.Deposits
	var votes gov.Votes
	var proposals gov.Proposals
	var depositParams gov.DepositParams
	var votingParams gov.VotingParams
	var tallyParams gov.TallyParams
	// TODO

	return gov.GenesisState{
		StartingProposalID: startingProposalID,
		Deposits:           deposits,
		Votes:              votes,
		Proposals:          proposals,
		DepositParams:      depositParams,
		VotingParams:       votingParams,
		TallyParams:        tallyParams,
	}
}

func migrateSlashing(initialState v0_17.GenesisFileState) slashing.GenesisState {
	slashingParams := slashing.Params{
		SignedBlocksWindow:      initialState.SlashingData.Params.SignedBlocksWindow,
		MinSignedPerWindow:      initialState.SlashingData.Params.MinSignedPerWindow,
		DowntimeJailDuration:    initialState.SlashingData.Params.DowntimeJailDuration,
		SlashFractionDoubleSign: initialState.SlashingData.Params.SlashFractionDoubleSign,
		SlashFractionDowntime:   initialState.SlashingData.Params.SlashFractionDowntime,
	}
	var signingInfos map[string]slashing.ValidatorSigningInfo
	for ba, vs := range initialState.SlashingData.SigningInfos {
		signingInfos[ba] = slashing.ValidatorSigningInfo{
			// Address: , // TODO
			StartHeight: vs.StartHeight,
			IndexOffset: vs.IndexOffset,
			JailedUntil: vs.JailedUntil,
			// Tombstoned: , // TODO
			MissedBlocksCounter: vs.MissedBlocksCounter,
		}

	}
	var mMissedBlocks map[string][]slashing.MissedBlock
	for ba, mbs := range initialState.SlashingData.MissedBlocks {
		var missedBlocks []slashing.MissedBlock
		for _, mb := range mbs {
			missedBlocks = append(
				missedBlocks,
				slashing.MissedBlock{
					Index:  mb.Index,
					Missed: mb.Missed,
				},
			)
		}
		mMissedBlocks[ba] = missedBlocks
	}

	return slashing.GenesisState{
		Params:       slashingParams,
		SigningInfos: signingInfos,
		MissedBlocks: mMissedBlocks,
	}
}

func migrateStaking(initialState v0_17.GenesisFileState) staking.GenesisState {
	stakingParams := staking.Params{
		UnbondingTime:     initialState.StakeData.Params.UnbondingTime,
		MaxValidators:     initialState.StakeData.Params.MaxValidators,
		MaxEntries:        staking.DefaultParams().MaxEntries,        // TODO
		HistoricalEntries: staking.DefaultParams().HistoricalEntries, // TODO
		BondDenom:         staking.DefaultParams().BondDenom,         // TODO
	}
	lastTotalPower := initialState.StakeData.LastTotalPower
	var lastValidatorPowers []staking.LastValidatorPower
	for _, lvp := range initialState.StakeData.LastValidatorPowers {
		lastValidatorPowers = append(
			lastValidatorPowers,
			staking.LastValidatorPower{
				Address: lvp.Address,
				Power:   lvp.Power.Int64(), // TODO
			},
		)
	}
	var validators staking.Validators
	for _, v := range initialState.StakeData.Validators {
		validators = append(
			validators,
			staking.Validator{
				OperatorAddress: v.OperatorAddr,
				ConsPubKey:      v.ConsPubKey,
				Jailed:          v.Jailed,
				Status:          v.Status,
				Tokens:          v.Tokens.TruncateInt(), // TODO
				DelegatorShares: v.DelegatorShares,
				Description: staking.Description{
					Moniker:         v.Description.Moniker,
					Identity:        v.Description.Identity,
					Website:         v.Description.Website,
					SecurityContact: "", // TODO
					Details:         v.Description.Details,
				},
				UnbondingHeight:         v.UnbondingHeight,
				UnbondingCompletionTime: v.UnbondingMinTime,
				Commission: staking.Commission{
					CommissionRates: staking.CommissionRates{
						Rate:          v.Commission.Rate,
						MaxRate:       v.Commission.MaxRate,
						MaxChangeRate: v.Commission.MaxChangeRate,
					},
					UpdateTime: v.Commission.UpdateTime,
				},
				MinSelfDelegation: sdk.OneInt(), // TODO
			},
		)
	}
	var delegations staking.Delegations
	for _, b := range initialState.StakeData.Bonds {
		delegations = append(
			delegations,
			staking.Delegation{
				DelegatorAddress: b.DelegatorAddr,
				ValidatorAddress: b.ValidatorAddr,
				Shares:           b.Shares,
			},
		)
	}
	var unbondingDelegations []staking.UnbondingDelegation
	for _, b := range initialState.StakeData.UnbondingDelegations {
		unbondingDelegations = append(
			unbondingDelegations,
			staking.UnbondingDelegation{
				DelegatorAddress: b.DelegatorAddr,
				ValidatorAddress: b.ValidatorAddr,
				Entries:          []staking.UnbondingDelegationEntry{}, // TODO
			},
		)
	}
	var redelegations []staking.Redelegation
	for _, r := range initialState.StakeData.Redelegations {
		redelegations = append(
			redelegations,
			staking.Redelegation{
				DelegatorAddress:    r.DelegatorAddr,
				ValidatorSrcAddress: r.ValidatorSrcAddr,
				ValidatorDstAddress: r.ValidatorDstAddr,
				Entries:             []staking.RedelegationEntry{}, // TODO
			},
		)
	}
	exported := initialState.StakeData.Exported

	return staking.GenesisState{
		Params:               stakingParams,
		LastTotalPower:       lastTotalPower,
		LastValidatorPowers:  lastValidatorPowers,
		Validators:           validators,
		Delegations:          delegations,
		UnbondingDelegations: unbondingDelegations,
		Redelegations:        redelegations,
		Exported:             exported,
	}
}

func migrateRand(initialState v0_17.GenesisFileState) rand.GenesisState {
	var pendingRandRequests map[string][]rand.Request
	for lh, rs := range initialState.RandData.PendingRandRequests {
		var requests []rand.Request
		for _, r := range rs {
			requests = append(
				requests,
				rand.Request{
					Height:   r.Height,
					Consumer: r.Consumer,
					TxHash:   r.TxHash,
				},
			)
		}
		pendingRandRequests[lh] = requests
	}

	return rand.GenesisState{
		PendingRandRequests: pendingRandRequests,
	}
}

func migrateHTLC(initialState v0_17.GenesisFileState) htlc.GenesisState {
	var pendingHTLCs map[string]htlc.HTLC
	for hk, h := range initialState.HtlcData.PendingHTLCs {
		pendingHTLCs[hk] = htlc.NewHTLC(
			h.Sender,
			h.To,
			h.ReceiverOnOtherChain,
			h.Amount,
			htlc.HTLCSecret(h.Secret),
			h.Timestamp,
			h.ExpireHeight,
			htlc.HTLCState(h.State),
		)
	}

	return htlc.GenesisState{
		PendingHTLCs: pendingHTLCs,
	}
}

func migrateMint(initialState v0_17.GenesisFileState) mint.GenesisState {
	minter := mint.Minter{
		LastUpdate:    initialState.MintData.Minter.LastUpdate,
		InflationBase: initialState.MintData.Minter.InflationBase,
	}
	mintParams := mint.Params{
		Inflation: initialState.MintData.Params.Inflation,
		MintDenom: mint.DefaultParams().MintDenom, // TODO
	}

	return mint.GenesisState{
		Minter: minter,
		Params: mintParams,
	}
}

func migrateGuardian(initialState v0_17.GenesisFileState) guardian.GenesisState {
	var profilers guardian.Profilers
	var trustees guardian.Trustees
	// TODO

	return guardian.GenesisState{
		Profilers: profilers,
		Trustees:  trustees,
	}
}

func migrateCoinswap(initialState v0_17.GenesisFileState) coinswap.GenesisState {
	var coinswapParams coinswap.Params
	// TODO

	return coinswap.GenesisState{
		Params: coinswapParams,
	}
}

func migrateAsset(initialState v0_17.GenesisFileState) asset.GenesisState {
	var tokenState token.GenesisState
	// TODO

	return asset.GenesisState{
		TokenState: tokenState,
	}
}

func migrateService(initialState v0_17.GenesisFileState) service.GenesisState {
	var serviceParams service.Params
	// TODO

	return service.GenesisState{
		Params: serviceParams,
	}
}
