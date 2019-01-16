package cli

import (
	"fmt"
	"os"

	"github.com/irisnet/irishub/client/context"
	sdk "github.com/irisnet/irishub/types"
	upgcli "github.com/irisnet/irishub/client/upgrade"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/upgrade"
	"github.com/spf13/cobra"
)

func GetInfoCmd(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "info",
		Short:   "query the information of upgrade module",
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

			output, err := codec.MarshalJSONIndent(cdc, upgradeInfoOutput)
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
		},
	}
	return cmd
}

func GetCmdQuerySignals(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-signals",
		Short:   "query the information of signals",
		Example: "iriscli upgrade status",
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

			var validatorAddrs []string
			res, err := cliCtx.QuerySubspace(upgrade.GetSignalPrefixKey(upgradeConfig.Protocol.Version), storeName)
			if err != nil {
				return err
			}

			for _, kv := range res {
				validatorAddrs = append(validatorAddrs, upgrade.GetAddressFromSignalKey(kv.Key))
			}

			if len(validatorAddrs) == 0 {
				fmt.Println("No validators have started the new version.")
				return nil
			}

			output, err := codec.MarshalJSONIndent(cdc, validatorAddrs)
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
		},
	}
	return cmd
}
