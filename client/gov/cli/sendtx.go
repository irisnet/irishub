package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/irisnet/irishub/app/v1/gov"
	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/client/context"
	client "github.com/irisnet/irishub/client/gov"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetCmdSubmitProposal implements submitting a proposal transaction command.
func GetCmdSubmitProposal(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "submit-proposal",
		Short:   "Submit a proposal along with an initial deposit",
		Example: "iriscli gov submit-proposal --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --type=Parameter --description=test --title=test-proposal --param='mint/Inflation=0.050'",
		RunE: func(cmd *cobra.Command, args []string) error {
			title := viper.GetString(flagTitle)
			description := viper.GetString(flagDescription)
			strProposalType := viper.GetString(flagProposalType)
			initialDeposit := viper.GetString(flagDeposit)

			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			fromAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			amount, err := cliCtx.ParseCoins(initialDeposit)
			if err != nil {
				return err
			}

			proposalType, err := gov.ProposalTypeFromString(strProposalType)
			if err != nil {
				return err
			}
			var params gov.Params
			if proposalType == gov.ProposalTypeParameter {
				paramStr := viper.GetString(flagParam)
				params, err = getParamFromString(paramStr)
				if err != nil {
					return err
				}
				if len(params) != 1 {
					return errors.New("the length of ParameterProposal's param should be one")
				}
				if err := client.ValidateParam(params[0]); err != nil {
					return err
				}
			}
			msg := gov.NewMsgSubmitProposal(title, description, proposalType, fromAddr, amount, params)
			if proposalType == gov.ProposalTypeCommunityTaxUsage {
				usageStr := viper.GetString(flagUsage)
				usage, err := gov.UsageTypeFromString(usageStr)
				if err != nil {
					return err
				}
				var destAddr sdk.AccAddress
				if usage.String() != "Burn" {
					destAddrStr := viper.GetString(flagDestAddress)
					destAddr, err = sdk.AccAddressFromBech32(destAddrStr)
					if err != nil {
						return err
					}
				}
				percentStr := viper.GetString(flagPercent)
				percent, err := sdk.NewDecFromStr(percentStr)
				if err != nil {
					return err
				}
				taxMsg := gov.NewMsgSubmitCommunityTaxUsageProposal(msg, usage, destAddr, percent)
				return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{taxMsg})
			}

			if proposalType == gov.ProposalTypeSoftwareUpgrade {

				versionInt := viper.GetInt64(flagVersion)
				if versionInt < 0 {
					return errors.Errorf("Version must greater than or equal to zero")
				}

				version := uint64(versionInt)
				software := viper.GetString(flagSoftware)

				switchHeightInt := viper.GetInt64(flagSwitchHeight)
				if switchHeightInt < 0 {
					return errors.Errorf("SwitchHeight must greater than or equal to zero")
				}
				switchHeight := uint64(switchHeightInt)

				thresholdStr := viper.GetString(flagThreshold)
				threshold, err := sdk.NewDecFromStr(thresholdStr)
				if err != nil {
					return err
				}
				msg := gov.NewMsgSubmitSoftwareUpgradeProposal(msg, version, software, switchHeight, threshold)
				return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
			}

			if proposalType == gov.ProposalTypeTokenAddition {
				symbol := viper.GetString(flagTokenSymbol)
				canonicalSymbol := viper.GetString(flagTokenCanonicalSymbol)
				name := viper.GetString(flagTokenName)
				decimal := uint8(viper.GetInt(flagTokenDecimal))
				alias := viper.GetString(flagTokenMinUnitAlias)

				msg := gov.NewMsgSubmitTokenAdditionProposal(msg, symbol, canonicalSymbol, name, alias, decimal)
				return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
			}
			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagTitle, "", "title of proposal")
	cmd.Flags().String(flagDescription, "", "description of proposal")
	cmd.Flags().String(flagProposalType, "", "proposalType of proposal,eg:PlainText/Parameter/SoftwareUpgrade/SystemHalt/CommunityTaxUsage/TokenAddition")
	cmd.Flags().String(flagDeposit, "", "deposit of proposal(at least 30% of MinDeposit)")
	cmd.Flags().String(flagParam, "", "parameter of proposal,eg. key=value")
	cmd.Flags().String(flagUsage, "", "the transaction fee tax usage type, valid values can be Burn, Distribute and Grant")
	cmd.Flags().String(flagPercent, "", "percent of transaction fee tax pool to use, integer or decimal >0 and <=1")
	cmd.Flags().String(flagDestAddress, "", "the destination trustee address")

	cmd.Flags().String(flagVersion, "0", "the version of the new protocol")
	cmd.Flags().String(flagSoftware, " ", "the software of the new protocol")
	cmd.Flags().String(flagSwitchHeight, "0", "the switchheight of the new protocol")
	cmd.Flags().String(flagThreshold, "0.8", "the upgrade signal threshold of the software upgrade")

	//for TokenAdditionProposal
	cmd.Flags().String(flagTokenSymbol, "", "the asset symbol. Once created, it cannot be modified")
	cmd.Flags().String(flagTokenCanonicalSymbol, "", "the source symbol of a external asset")
	cmd.Flags().String(flagTokenName, "", "the asset name")
	cmd.Flags().Uint8(flagTokenDecimal, 0, "the asset decimal. The maximum value is 18")
	cmd.Flags().String(flagTokenMinUnitAlias, "", "the asset symbol minimum alias")

	cmd.MarkFlagRequired(flagTitle)
	cmd.MarkFlagRequired(flagDescription)
	cmd.MarkFlagRequired(flagProposalType)
	return cmd
}

func getParamFromString(paramsStr string) (gov.Params, error) {
	var govParams gov.Params
	str := strings.Split(paramsStr, "=")
	if len(str) != 2 {
		return gov.Params{}, fmt.Errorf("%s is not valid", paramsStr)
	}
	//str = []string{"mint/Inflation","0.0000000000"}
	//params.GetParamSpaceFromKey(str[0]) == "mint"
	//params.GetParamKey(str[0])          == "Inflation"
	govParams = append(govParams,
		gov.Param{Subspace: params.GetParamSpaceFromKey(str[0]),
			Key:   params.GetParamKey(str[0]),
			Value: str[1]})
	return govParams, nil
}

// GetCmdDeposit implements depositing tokens for an active proposal.
func GetCmdDeposit(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "deposit",
		Short:   "Deposit tokens for activing proposal",
		Example: "iriscli gov deposit --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --proposal-id=1 --deposit=10iris",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			depositorAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			proposalID := uint64(viper.GetInt64(flagProposalID))

			amount, err := cliCtx.ParseCoins(viper.GetString(flagDeposit))

			if err != nil {
				return err
			}

			msg := gov.NewMsgDeposit(depositorAddr, proposalID, amount)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagProposalID, "", "proposalID of proposal depositing on")
	cmd.Flags().String(flagDeposit, "", "amount of deposit")
	cmd.MarkFlagRequired(flagProposalID)
	cmd.MarkFlagRequired(flagDeposit)
	return cmd
}

// GetCmdVote implements creating a new vote command.
func GetCmdVote(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "vote",
		Short:   "Vote for an active proposal, options: Yes/No/NoWithVeto/Abstain",
		Example: "iriscli gov vote --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --proposal-id=1 --option=Yes",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			voterAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			proposalID := uint64(viper.GetInt64(flagProposalID))
			option := viper.GetString(flagOption)

			byteVoteOption, err := gov.VoteOptionFromString(client.NormalizeVoteOption(option))
			if err != nil {
				return err
			}

			msg := gov.NewMsgVote(voterAddr, proposalID, byteVoteOption)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			fmt.Printf("Vote[Voter:%s,ProposalID:%d,Option:%s]",
				voterAddr.String(), msg.ProposalID, msg.Option.String(),
			)
			// Build and sign the transaction, then broadcast to a Tendermint
			// node.
			cliCtx.PrintResponse = true

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagProposalID, "", "proposalID of proposal voting on")
	cmd.Flags().String(flagOption, "", "vote option {Yes, No, NoWithVeto, Abstain}")
	cmd.MarkFlagRequired(flagProposalID)
	cmd.MarkFlagRequired(flagOption)
	return cmd
}
