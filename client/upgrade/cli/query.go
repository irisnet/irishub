package cli

import (
	"fmt"
	"os"

	protocol "github.com/irisnet/irishub/app/protocol/keeper"
	"github.com/irisnet/irishub/client/context"
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

			res_protocolVersion, _ := cliCtx.QueryStore(protocol.CurrentProtocolVersionKey, "protocol")
			var protocolVersion uint64
			cdc.MustUnmarshalBinaryLengthPrefixed(res_protocolVersion, &protocolVersion)

			res_proposalID, _ := cliCtx.QueryStore(upgrade.GetSuccessAppVersionKey(protocolVersion), storeName)
			var proposalID uint64
			cdc.MustUnmarshalBinaryLengthPrefixed(res_proposalID, &proposalID)

			res_appVersion, err := cliCtx.QueryStore(upgrade.GetProposalIDKey(proposalID), storeName)
			var appVersion upgrade.AppVersion
			cdc.MustUnmarshalBinaryLengthPrefixed(res_appVersion, &appVersion)

			res_upgradeConfig, _ := cliCtx.QueryStore(protocol.UpgradeConfigkey, "protocol")
			var upgradeConfig protocol.UpgradeConfig
			if err == nil && len(res_upgradeConfig) != 0 {
				cdc.MustUnmarshalBinaryLengthPrefixed(res_upgradeConfig, &upgradeConfig)
			}

			res_LastFailureVersion, err := cliCtx.QueryStore(protocol.LastFailureVersionKey, "protocol")
			var lastFailureVersion uint64
			if err == nil && len(res_LastFailureVersion) != 0 {
				cdc.MustUnmarshalBinaryLengthPrefixed(res_LastFailureVersion, &lastFailureVersion)
			} else {
				lastFailureVersion = 0
			}

			upgradeInfoOutput := upgcli.ConvertUpgradeInfoToUpgradeOutput(appVersion, upgradeConfig, lastFailureVersion)

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

			res_upgradeConfig, err := cliCtx.QueryStore(protocol.UpgradeConfigkey, "protocol")
			if err != nil {
				return err
			}
			if len(res_upgradeConfig) == 0 {
				fmt.Println("No Software Upgrade Switch Period is in process.")
				return err
			}

			var upgradeConfig protocol.UpgradeConfig
			if err = cdc.UnmarshalBinaryLengthPrefixed(res_upgradeConfig, &upgradeConfig); err != nil {
				return err
			}

			var validatorAddrs []string
			res, err := cliCtx.QuerySubspace(upgrade.GetSignalPrefixKey(upgradeConfig.Definition.Version), storeName)
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
