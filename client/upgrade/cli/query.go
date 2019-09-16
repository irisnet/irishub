package cli

import (
	"fmt"
	"os"

	"github.com/irisnet/irishub/app/v1/stake"
	"github.com/irisnet/irishub/app/v1/stake/types"
	"github.com/irisnet/irishub/app/v1/upgrade"
	"github.com/irisnet/irishub/client/context"
	upgcli "github.com/irisnet/irishub/client/upgrade"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagDetail = "detail"
)

func GetInfoCmd(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "info",
		Short:   "Query the information of upgrade module",
		Example: "iriscli upgrade info",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))

			res_currentVersion, _ := cliCtx.QueryStore(sdk.CurrentVersionKey, sdk.MainStore)
			var currentVersion uint64
			cdc.MustUnmarshalBinaryLengthPrefixed(res_currentVersion, &currentVersion)

			res_proposalID, _ := cliCtx.QueryStore(upgrade.GetSuccessVersionKey(currentVersion), storeName)
			var proposalID uint64
			cdc.MustUnmarshalBinaryLengthPrefixed(res_proposalID, &proposalID)

			res_currentVersionInfo, err := cliCtx.QueryStore(upgrade.GetProposalIDKey(proposalID), storeName)
			var currentVersionInfo upgrade.VersionInfo
			cdc.MustUnmarshalBinaryLengthPrefixed(res_currentVersionInfo, &currentVersionInfo)

			res_upgradeInProgress, _ := cliCtx.QueryStore(sdk.UpgradeConfigKey, sdk.MainStore)
			var upgradeInProgress sdk.UpgradeConfig
			if err == nil && len(res_upgradeInProgress) != 0 {
				cdc.MustUnmarshalBinaryLengthPrefixed(res_upgradeInProgress, &upgradeInProgress)
			}

			res_LastFailedVersion, err := cliCtx.QueryStore(sdk.LastFailedVersionKey, sdk.MainStore)
			var lastFailedVersion uint64
			if err == nil && len(res_LastFailedVersion) != 0 {
				cdc.MustUnmarshalBinaryLengthPrefixed(res_LastFailedVersion, &lastFailedVersion)
			} else {
				lastFailedVersion = 0
			}

			upgradeInfoOutput := upgcli.NewUpgradeInfoOutput(currentVersionInfo, lastFailedVersion, upgradeInProgress)

			return cliCtx.PrintOutput(upgradeInfoOutput)
		},
	}
	return cmd
}

func GetCmdQuerySignals(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-signals",
		Short:   "Query the information of signals",
		Example: "iriscli upgrade query-signals",
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))

			res_upgradeConfig, err := cliCtx.QueryStore(sdk.UpgradeConfigKey, sdk.MainStore)
			if err != nil {
				return err
			}
			if len(res_upgradeConfig) == 0 {
				fmt.Println("No Software Upgrade Switch Period is in process.")
				return err
			}

			var upgradeConfig sdk.UpgradeConfig
			if err = cdc.UnmarshalBinaryLengthPrefixed(res_upgradeConfig, &upgradeConfig); err != nil {
				return err
			}

			validatorConsAddrs := make(map[string]bool)
			res, err := cliCtx.QuerySubspace(upgrade.GetSignalPrefixKey(upgradeConfig.Protocol.Version), storeName)
			if err != nil {
				return err
			}

			for _, kv := range res {
				validatorConsAddrs[upgrade.GetAddressFromSignalKey(kv.Key)] = true
			}

			if len(validatorConsAddrs) == 0 {
				fmt.Println("No validator has started the new version.")
				return nil
			}

			key := stake.ValidatorsKey
			resKVs, err := cliCtx.QuerySubspace(key, "stake")
			if err != nil {
				return err
			}

			isDetail := viper.GetBool(flagDetail)
			totalVotingPower := sdk.ZeroDec()
			signalsVotingPower := sdk.ZeroDec()

			for _, kv := range resKVs {
				addr := kv.Key[1:]
				validator := types.MustUnmarshalValidator(cdc, addr, kv.Value)
				totalVotingPower = totalVotingPower.Add(validator.GetPower())
				if _, ok := validatorConsAddrs[validator.GetConsAddr().String()]; ok {
					signalsVotingPower = signalsVotingPower.Add(validator.GetPower())
					if isDetail {
						fmt.Println(validator.GetOperator().String(), " ", validator.GetPower().String())
					}
				}
			}
			fmt.Println("signalsVotingPower/totalVotingPower = " + signalsVotingPower.Quo(totalVotingPower).String())
			return nil
		},
	}
	cmd.Flags().Bool(flagDetail, false, "details of siganls")
	return cmd
}
