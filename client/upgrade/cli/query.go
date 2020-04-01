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

			resCurrentVersion, err := cliCtx.QueryStore(sdk.CurrentVersionKey, sdk.MainStore)
			if err != nil {
				return err
			}
			var currentVersion uint64
			cdc.MustUnmarshalBinaryLengthPrefixed(resCurrentVersion, &currentVersion)

			resProposalID, err := cliCtx.QueryStore(upgrade.GetSuccessVersionKey(currentVersion), storeName)
			if err != nil {
				return err
			}
			var proposalID uint64
			if len(resProposalID) > 0 {
				cdc.MustUnmarshalBinaryLengthPrefixed(resProposalID, &proposalID)
			}

			resCurrentVersionInfo, err := cliCtx.QueryStore(upgrade.GetProposalIDKey(proposalID), storeName)
			if err != nil {
				return err
			}
			var currentVersionInfo upgrade.VersionInfo
			if len(resCurrentVersionInfo) > 0 {
				cdc.MustUnmarshalBinaryLengthPrefixed(resCurrentVersionInfo, &currentVersionInfo)
			}

			resUpgradeInProgress, err := cliCtx.QueryStore(sdk.UpgradeConfigKey, sdk.MainStore)
			if err != nil {
				return err
			}
			var upgradeInProgress sdk.UpgradeConfig
			if len(resUpgradeInProgress) > 0 {
				cdc.MustUnmarshalBinaryLengthPrefixed(resUpgradeInProgress, &upgradeInProgress)
			}

			resLastFailedVersion, err := cliCtx.QueryStore(sdk.LastFailedVersionKey, sdk.MainStore)
			var lastFailedVersion uint64
			if err == nil && len(resLastFailedVersion) > 0 {
				cdc.MustUnmarshalBinaryLengthPrefixed(resLastFailedVersion, &lastFailedVersion)
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

			resUpgradeConfig, err := cliCtx.QueryStore(sdk.UpgradeConfigKey, sdk.MainStore)
			if err != nil {
				return err
			}
			if len(resUpgradeConfig) == 0 {
				fmt.Println("No Software Upgrade Switch Period is in process.")
				return err
			}

			var upgradeConfig sdk.UpgradeConfig
			if err = cdc.UnmarshalBinaryLengthPrefixed(resUpgradeConfig, &upgradeConfig); err != nil {
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
