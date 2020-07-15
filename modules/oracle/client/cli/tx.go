package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/irisnet/irishub/modules/oracle/types"
)

// GetTxCmd returns the transaction commands for the guardian module.
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "oracle transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		GetCmdCreateFeed(),
		GetCmdStartFeed(),
		GetCmdPauseFeed(),
		GetCmdEditFeed(),
	)
	return txCmd
}

// GetCmdCreateFeed implements defining a feed command
func GetCmdCreateFeed() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: `Create a new feed, the feed will be in "paused" state`,
		Example: fmt.Sprintf(`%s oracle create --chain-id=<chain-id> --from=<key-name> --fee=0.3iris `+
			`--feed-name="test-feed" `+
			`--latest-history=10 `+
			`--service-name="test-service" `+
			`--input=<request-data> `+
			`--providers=<provide1_address>,<provider2_address> `+
			`--service-fee-cap=1iris `+
			`--timeout=2 `+
			`--total=10 `+
			`--threshold=1 `+
			`--aggregate-func="avg" `+
			`--value-json-path="high"`, version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			creator := clientCtx.GetFromAddress()

			var providers []sdk.AccAddress
			for _, addr := range viper.GetStringSlice(FlagProviders) {
				provider, err := sdk.AccAddressFromBech32(addr)
				if err != nil {
					return err
				}
				providers = append(providers, provider)
			}

			serviceFeeCap, err := sdk.ParseCoins(viper.GetString(FlagServiceFeeCap))
			if err != nil {
				return err
			}

			msg := &types.MsgCreateFeed{
				FeedName:          viper.GetString(FlagFeedName),
				AggregateFunc:     viper.GetString(FlagAggregateFunc),
				ValueJsonPath:     viper.GetString(FlagValueJsonPath),
				LatestHistory:     uint64(viper.GetInt64(FlagLatestHistory)),
				Description:       viper.GetString(FlagDescription),
				ServiceName:       viper.GetString(FlagServiceName),
				Providers:         providers,
				Input:             viper.GetString(FlagInput),
				Timeout:           viper.GetInt64(FlagTimeout),
				ServiceFeeCap:     serviceFeeCap,
				RepeatedFrequency: uint64(viper.GetInt64(FlagFrequency)),
				ResponseThreshold: uint32(viper.GetInt32(FlagThreshold)),
				Creator:           creator,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTx(clientCtx, msg)
		},
	}
	cmd.Flags().AddFlagSet(FsCreateFeed)
	_ = cmd.MarkFlagRequired(FlagFeedName)
	_ = cmd.MarkFlagRequired(FlagAggregateFunc)
	_ = cmd.MarkFlagRequired(FlagValueJsonPath)
	_ = cmd.MarkFlagRequired(FlagLatestHistory)
	_ = cmd.MarkFlagRequired(FlagServiceName)
	_ = cmd.MarkFlagRequired(FlagProviders)
	_ = cmd.MarkFlagRequired(FlagServiceFeeCap)
	_ = cmd.MarkFlagRequired(FlagTimeout)
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdStartFeed implements start a feed command
func GetCmdStartFeed() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start [feed-name]",
		Short:   `Start a feed in "paused" state`,
		Args:    cobra.ExactArgs(1),
		Example: fmt.Sprintf(`%s oracle start <feed-name> --chain-id=<chain-id> --from=<key-name> --fee=0.3iris`, version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			creator := clientCtx.GetFromAddress()
			msg := &types.MsgStartFeed{
				FeedName: args[0],
				Creator:  creator,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTx(clientCtx, msg)
		},
	}
	cmd.Flags().AddFlagSet(FsStartFeed)
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdPauseFeed implements pause a running feed command
func GetCmdPauseFeed() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pause [feed-name]",
		Short:   `Pause a feed in "running" state`,
		Args:    cobra.ExactArgs(1),
		Example: fmt.Sprintf(`%s oracle pause <feed-name> --chain-id=<chain-id> --from=<key-name> --fee=0.3iris`, version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			creator := clientCtx.GetFromAddress()
			msg := &types.MsgPauseFeed{
				FeedName: args[0],
				Creator:  creator,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTx(clientCtx, msg)
		},
	}
	cmd.Flags().AddFlagSet(FsStartFeed)
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdEditFeed implements edit a feed command
func GetCmdEditFeed() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit [feed-name]",
		Short: "Modify the feed information and update service invocation parameters by feed creator",
		Args:  cobra.ExactArgs(1),
		Example: fmt.Sprintf(`%s oracle edit <feed-name> --chain-id=<chain-id> --from=<key-name> --fee=0.3iris `+
			`--latest-history=10 `+
			`--providers=<provide1_address>,<provider2_address> `+
			`--service-fee-cap=1iris `+
			`--timeout=2 `+
			`--frequency=10 `+
			`--threshold=5 `+
			`--threshold=1`, version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			creator := clientCtx.GetFromAddress()
			var providers []sdk.AccAddress
			for _, addr := range viper.GetStringSlice(FlagProviders) {
				provider, err := sdk.AccAddressFromBech32(addr)
				if err != nil {
					return err
				}
				providers = append(providers, provider)
			}

			serviceFeeCap, err := sdk.ParseCoins(viper.GetString(FlagServiceFeeCap))
			if err != nil {
				return err
			}

			msg := &types.MsgEditFeed{
				FeedName:          args[0],
				Description:       viper.GetString(FlagDescription),
				LatestHistory:     uint64(viper.GetInt64(FlagLatestHistory)),
				Providers:         providers,
				Timeout:           viper.GetInt64(FlagTimeout),
				ServiceFeeCap:     serviceFeeCap,
				RepeatedFrequency: uint64(viper.GetInt64(FlagFrequency)),
				ResponseThreshold: uint32(viper.GetInt(FlagThreshold)),
				Creator:           creator,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTx(clientCtx, msg)
		},
	}
	cmd.Flags().AddFlagSet(FsEditFeed)
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
