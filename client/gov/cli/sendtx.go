package cli

import (
	"fmt"
	"os"

	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/codec"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/modules/gov"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	client "github.com/irisnet/irishub/client/gov"
)

// GetCmdSubmitProposal implements submitting a proposal transaction command.
func GetCmdSubmitProposal(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-proposal",
		Short: "Submit a proposal along with an initial deposit",
		RunE: func(cmd *cobra.Command, args []string) error {
			title := viper.GetString(flagTitle)
			description := viper.GetString(flagDescription)
			strProposalType := viper.GetString(flagProposalType)
			initialDeposit := viper.GetString(flagDeposit)
			////////////////////  iris begin  ///////////////////////////
			paramStr := viper.GetString(flagParam)
			////////////////////  iris end  /////////////////////////////

			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			txCtx := context.NewTxContextFromCLI().WithCodec(cdc).
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
			////////////////////  iris begin  ///////////////////////////
			var param gov.Param
			if proposalType == gov.ProposalTypeParameterChange {
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

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagTitle, "", "title of proposal")
	cmd.Flags().String(flagDescription, "", "description of proposal")
	cmd.Flags().String(flagProposalType, "", "proposalType of proposal,eg:Text/ParameterChange/SoftwareUpgrade")
	cmd.Flags().String(flagDeposit, "", "deposit of proposal")
	////////////////////  iris begin  ///////////////////////////
	cmd.Flags().String(flagParam, "", "parameter of proposal,eg. [{key:key,value:value,op:update}]")
	cmd.Flags().String(flagKey, "", "the key of parameter")
	cmd.Flags().String(flagOp, "", "the operation of parameter")
	cmd.Flags().String(flagPath, "", "the path of param.json")
	////////////////////  iris end  /////////////////////////////
	return cmd
}
////////////////////  iris begin  ///////////////////////////
func getParamFromString(paramStr string, pathStr string, keyStr string, opStr string, cdc *codec.Codec) (gov.Param, error) {
	var param gov.Param

	if paramStr != "" {
		err := json.Unmarshal([]byte(paramStr), &param)
		return param, err

	} else if pathStr != "" {
		paramDoc := gov.ParameterConfigFile{}
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
		Use:   "deposit",
		Short: "deposit tokens for activing proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			txCtx := context.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			depositerAddr, err := cliCtx.GetFromAddress()
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

			msg := gov.NewMsgDeposit(depositerAddr, proposalID, amount)

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
		Use:   "vote",
		Short: "vote for an active proposal, options: Yes/No/NoWithVeto/Abstain",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			txCtx := context.NewTxContextFromCLI().WithCodec(cdc).
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
