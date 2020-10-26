package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/irisnet/irismod/modules/oracle/types"
)

// NewTxCmd returns the transaction commands for the oracle module.
func NewTxCmd() *cobra.Command {
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
		Example: fmt.Sprintf(
			`%s tx oracle create --chain-id=<chain-id> --from=<key-name> --fees=0.3iris `+
				`--feed-name="test-feed" `+
				`--latest-history=10 `+
				`--service-name="test-service" `+
				`--input=<request data or path/to/input.json> `+
				`--providers=<provide1_address>,<provider2_address> `+
				`--service-fee-cap=1iris `+
				`--timeout=2 `+
				`--total=10 `+
				`--threshold=1 `+
				`--aggregate-func="avg" `+
				`--value-json-path="high"`,
			version.AppName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			creator := clientCtx.GetFromAddress()
			rawProviders, err := cmd.Flags().GetStringSlice(FlagProviders)
			if err != nil {
				return err
			}

			var providers []sdk.AccAddress
			for _, addr := range rawProviders {
				provider, err := sdk.AccAddressFromBech32(addr)
				if err != nil {
					return err
				}
				providers = append(providers, provider)
			}

			rawServiceFeeCap, err := cmd.Flags().GetString(FlagServiceFeeCap)
			if err != nil {
				return err
			}
			serviceFeeCap, err := sdk.ParseCoins(rawServiceFeeCap)
			if err != nil {
				return err
			}

			input, err := cmd.Flags().GetString(FlagInput)
			if err != nil {
				return err
			}
			if !json.Valid([]byte(input)) {
				inputContent, err := ioutil.ReadFile(input)
				if err != nil {
					return fmt.Errorf("invalid input data: neither JSON input nor path to .json file were provided")
				}

				if !json.Valid(inputContent) {
					return fmt.Errorf("invalid input data: .json file content is invalid JSON")
				}

				input = string(inputContent)
			}

			buf := bytes.NewBuffer([]byte{})
			if err := json.Compact(buf, []byte(input)); err != nil {
				return fmt.Errorf("failed to compact the input data")
			}

			input = buf.String()

			feedName, err := cmd.Flags().GetString(FlagFeedName)
			if err != nil {
				return err
			}
			aggregateFunc, err := cmd.Flags().GetString(FlagAggregateFunc)
			if err != nil {
				return err
			}
			valueJsonPath, err := cmd.Flags().GetString(FlagValueJsonPath)
			if err != nil {
				return err
			}
			latestHistory, err := cmd.Flags().GetUint64(FlagLatestHistory)
			if err != nil {
				return err
			}
			description, err := cmd.Flags().GetString(FlagDescription)
			if err != nil {
				return err
			}
			serviceName, err := cmd.Flags().GetString(FlagServiceName)
			if err != nil {
				return err
			}
			timeout, err := cmd.Flags().GetInt64(FlagTimeout)
			if err != nil {
				return err
			}
			frequency, err := cmd.Flags().GetUint64(FlagFrequency)
			if err != nil {
				return err
			}
			threshold, err := cmd.Flags().GetUint32(FlagThreshold)
			if err != nil {
				return err
			}
			msg := &types.MsgCreateFeed{
				FeedName:          feedName,
				AggregateFunc:     aggregateFunc,
				ValueJsonPath:     valueJsonPath,
				LatestHistory:     latestHistory,
				Description:       description,
				ServiceName:       serviceName,
				Providers:         providers,
				Input:             input,
				Timeout:           timeout,
				ServiceFeeCap:     serviceFeeCap,
				RepeatedFrequency: frequency,
				ResponseThreshold: threshold,
				Creator:           creator,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
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
		Use:   "start [feed-name]",
		Short: `Start a feed in "paused" state.`,
		Args:  cobra.ExactArgs(1),
		Example: fmt.Sprintf(
			`%s tx oracle start <feed-name> --chain-id=<chain-id> --from=<key-name> --fees=0.3iris`,
			version.AppName,
		),
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

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FsStartFeed)
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdPauseFeed implements pause a running feed command
func GetCmdPauseFeed() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pause [feed-name]",
		Short: `Pause a feed in "running" state`,
		Args:  cobra.ExactArgs(1),
		Example: fmt.Sprintf(
			`%s tx oracle pause <feed-name> --chain-id=<chain-id> --from=<key-name> --fees=0.3iris`,
			version.AppName,
		),
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

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
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
		Short: "Modify the feed information and update service invocation parameters by feed creator.",
		Args:  cobra.ExactArgs(1),
		Example: fmt.Sprintf(
			`%s tx oracle edit <feed-name> --chain-id=<chain-id> --from=<key-name> --fees=0.3iris `+
				`--latest-history=10 `+
				`--providers=<provide1_address>,<provider2_address> `+
				`--service-fee-cap=1iris `+
				`--timeout=2 `+
				`--frequency=10 `+
				`--threshold=5 `+
				`--threshold=1`,
			version.AppName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			creator := clientCtx.GetFromAddress()
			rawProviders, err := cmd.Flags().GetStringSlice(FlagProviders)
			if err != nil {
				return err
			}
			var providers []sdk.AccAddress
			for _, addr := range rawProviders {
				provider, err := sdk.AccAddressFromBech32(addr)
				if err != nil {
					return err
				}
				providers = append(providers, provider)
			}

			rawServiceFeeCap, err := cmd.Flags().GetString(FlagServiceFeeCap)
			if err != nil {
				return err
			}
			serviceFeeCap, err := sdk.ParseCoins(rawServiceFeeCap)
			if err != nil {
				return err
			}

			latestHistory, err := cmd.Flags().GetUint64(FlagLatestHistory)
			if err != nil {
				return err
			}
			description, err := cmd.Flags().GetString(FlagDescription)
			if err != nil {
				return err
			}
			timeout, err := cmd.Flags().GetInt64(FlagTimeout)
			if err != nil {
				return err
			}
			frequency, err := cmd.Flags().GetUint64(FlagFrequency)
			if err != nil {
				return err
			}
			threshold, err := cmd.Flags().GetUint32(FlagThreshold)
			if err != nil {
				return err
			}
			msg := &types.MsgEditFeed{
				FeedName:          args[0],
				Description:       description,
				LatestHistory:     latestHistory,
				Providers:         providers,
				Timeout:           timeout,
				ServiceFeeCap:     serviceFeeCap,
				RepeatedFrequency: frequency,
				ResponseThreshold: threshold,
				Creator:           creator,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FsEditFeed)
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
