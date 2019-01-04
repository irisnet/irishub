package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/client/context"
	client "github.com/irisnet/irishub/client/gov"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/app/v1/gov"
	sdk "github.com/irisnet/irishub/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	govtypes "github.com/irisnet/irishub/types/gov"
)

// GetCmdSubmitProposal implements submitting a proposal transaction command.
func GetCmdSubmitProposal(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "submit-proposal",
		Short:   "Submit a proposal along with an initial deposit",
		Example: "iriscli gov submit-proposal --chain-id=<chain-id> --from=<key name> --fee=0.004iris --type=Text --description=test --title=test-proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			title := viper.GetString(flagTitle)
			description := viper.GetString(flagDescription)
			strProposalType := viper.GetString(flagProposalType)
			initialDeposit := viper.GetString(flagDeposit)
			////////////////////  iris begin  ///////////////////////////
			paramStr := viper.GetString(flagParam)
			////////////////////  iris end  /////////////////////////////

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

			proposalType, err := govtypes.ProposalTypeFromString(strProposalType)
			if err != nil {
				return err
			}
			////////////////////  iris begin  ///////////////////////////
			var param govtypes.Param
			if proposalType == govtypes.ProposalTypeParameterChange {
				pathStr := viper.GetString(flagPath)
				keyStr := viper.GetString(flagKey)
				opStr := viper.GetString(flagOp)
				param, err = getParamFromString(paramStr, pathStr, keyStr, opStr, cdc)
				if err != nil {
					return err
				}
			}
			////////////////////  iris end  /////////////////////////////

			msg := gov.NewMsgSubmitProposal(title, description, proposalType, fromAddr, amount, param)
			if proposalType == govtypes.ProposalTypeTxTaxUsage {
				usageStr := viper.GetString(flagUsage)
				usage, err := govtypes.UsageTypeFromString(usageStr)
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
				taxMsg := gov.NewMsgSubmitTaxUsageProposal(msg, usage, destAddr, percent)
				return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{taxMsg})
			}

			if proposalType == govtypes.ProposalTypeSoftwareUpgrade {
				version := uint64(viper.GetInt64(flagVersion))
				software := viper.GetString(flagSoftware)
				switchHeight := uint64(viper.GetInt64(flagSwitchHeight))
				msg := gov.NewMsgSubmitSoftwareUpgradeProposal(msg, version, software, switchHeight)
				return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
			}
			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagTitle, "", "title of proposal")
	cmd.Flags().String(flagDescription, "", "description of proposal")
	cmd.Flags().String(flagProposalType, "", "proposalType of proposal,eg:Text/ParameterChange/SoftwareUpgrade/SoftwareHalt/TxTaxUsage")
	cmd.Flags().String(flagDeposit, "", "deposit of proposal")
	////////////////////  iris begin  ///////////////////////////
	cmd.Flags().String(flagParam, "", "parameter of proposal,eg. [{key:key,value:value,op:update}]")
	cmd.Flags().String(flagKey, "", "the key of parameter")
	cmd.Flags().String(flagOp, "", "the operation of parameter")
	cmd.Flags().String(flagPath, app.DefaultCLIHome, "the directory of the param.json")
	cmd.Flags().String(flagUsage, "", "the transaction fee tax usage type, valid values can be Burn, Distribute and Grant")
	cmd.Flags().String(flagPercent, "", "percent of transaction fee tax pool to use, integer or decimal >0 and <=1")
	cmd.Flags().String(flagDestAddress, "", "the destination trustee address")

	cmd.Flags().String(flagVersion, "0", "the version of the new protocol")
	cmd.Flags().String(flagSoftware, " ", "the software of the new protocol")
	cmd.Flags().String(flagSwitchHeight, "0", "the switchheight of the new protocol")
	////////////////////  iris end  /////////////////////////////

	cmd.MarkFlagRequired(flagTitle)
	cmd.MarkFlagRequired(flagDescription)
	cmd.MarkFlagRequired(flagProposalType)
	return cmd
}

////////////////////  iris begin  ///////////////////////////
func getParamFromString(paramStr string, pathStr string, keyStr string, opStr string, cdc *codec.Codec) (govtypes.Param, error) {
	var param govtypes.Param

	if paramStr != "" {
		err := json.Unmarshal([]byte(paramStr), &param)
		return param, err

	} else if pathStr != "" {
		paramDoc := govtypes.ParameterConfigFile{}
		err := paramDoc.ReadFile(cdc, pathStr)
		if err != nil {
			return param, err
		}
		param, err := paramDoc.GetParamFromKey(keyStr, opStr)
		return param, err
	} else {

		return param, errors.New("Path and param are both empty")
	}
}

////////////////////  iris end  /////////////////////////////

// GetCmdDeposit implements depositing tokens for an active proposal.
func GetCmdDeposit(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "deposit",
		Short:   "deposit tokens for activing proposal",
		Example: "iriscli gov deposit --chain-id=<chain-id> --from=<key name> --fee=0.004iris --proposal-id=1 --deposit=10iris",
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

			////////////////////  iris begin  ///////////////////////////
			amount, err := cliCtx.ParseCoins(viper.GetString(flagDeposit))
			////////////////////  iris end  /////////////////////////////

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
		Short:   "vote for an active proposal, options: Yes/No/NoWithVeto/Abstain",
		Example: "iriscli gov vote --chain-id=<chain-id> --from=<key name> --fee=0.004iris --proposal-id=1 --option=Yes",
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

			byteVoteOption, err := govtypes.VoteOptionFromString(client.NormalizeVoteOption(option))
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
