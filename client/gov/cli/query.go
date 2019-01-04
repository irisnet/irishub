package cli

import (
	"encoding/json"
	"fmt"

	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/client/context"
	client "github.com/irisnet/irishub/client/gov"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/app/v1/gov"
	"github.com/irisnet/irishub/types/gov/params"
	"github.com/irisnet/irishub/modules/params"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	govtypes "github.com/irisnet/irishub/types/gov"
)

// GetCmdQueryProposal implements the query proposal command.
func GetCmdQueryProposal(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-proposal",
		Short:   "Query details of a single proposal",
		Example: "iriscli gov query-proposal --proposal-id=1",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			proposalID := uint64(viper.GetInt64(flagProposalID))

			params := gov.QueryProposalParams{
				ProposalID: proposalID,
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/proposal", queryRoute), bz)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}

	cmd.Flags().String(flagProposalID, "", "proposalID of proposal being queried")
	cmd.MarkFlagRequired(flagProposalID)
	return cmd
}

// nolint: gocyclo
// GetCmdQueryProposals implements a query proposals command.
func GetCmdQueryProposals(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-proposals",
		Short:   "query proposals with optional filters",
		Example: "iriscli gov query-proposals --status=Passed",
		RunE: func(cmd *cobra.Command, args []string) error {
			bechDepositorAddr := viper.GetString(flagDepositor)
			bechVoterAddr := viper.GetString(flagVoter)
			strProposalStatus := viper.GetString(flagStatus)
			numLimit := uint64(viper.GetInt64(flagNumLimit))

			params := gov.QueryProposalsParams{
				Limit: numLimit,
			}

			if len(bechDepositorAddr) != 0 {
				depositorAddr, err := sdk.AccAddressFromBech32(bechDepositorAddr)
				if err != nil {
					return err
				}
				params.Depositor = depositorAddr
			}

			if len(bechVoterAddr) != 0 {
				voterAddr, err := sdk.AccAddressFromBech32(bechVoterAddr)
				if err != nil {
					return err
				}
				params.Voter = voterAddr
			}

			if len(strProposalStatus) != 0 {
				proposalStatus, err := govtypes.ProposalStatusFromString(client.NormalizeProposalStatus(strProposalStatus))
				if err != nil {
					return err
				}
				params.ProposalStatus = proposalStatus
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/proposals", queryRoute), bz)
			if err != nil {
				return err
			}
			////////////////////  iris begin  ///////////////////////////
			var matchingProposals gov.ProposalOutputs
			err = cdc.UnmarshalJSON(res, &matchingProposals)
			if err != nil {
				return err
			}

			if len(matchingProposals) == 0 {
				fmt.Println("No matching proposals found")
				return nil
			}

			for _, proposal := range matchingProposals {
				fmt.Printf("  %d - %s\n", proposal.ProposalID, proposal.Title)
			}
			////////////////////  iris end  /////////////////////////////
			return nil
		},
	}

	cmd.Flags().String(flagNumLimit, "", "(optional) limit to latest [number] proposals. Defaults to all proposals")
	cmd.Flags().String(flagDepositor, "", "(optional) filter by proposals deposited on by depositor")
	cmd.Flags().String(flagVoter, "", "(optional) filter by proposals voted on by voted")
	cmd.Flags().String(flagStatus, "", "(optional) filter proposals by proposal status")

	return cmd
}

// Command to Get a Proposal Information
// GetCmdQueryVote implements the query proposal vote command.
func GetCmdQueryVote(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-vote",
		Short:   "query vote",
		Example: "iriscli gov query-vote --proposal-id=1 --voter=<voter address>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			proposalID := uint64(viper.GetInt64(flagProposalID))

			voterAddr, err := sdk.AccAddressFromBech32(viper.GetString(flagVoter))
			if err != nil {
				return err
			}

			params := gov.QueryVoteParams{
				Voter:      voterAddr,
				ProposalID: proposalID,
			}
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/vote", queryRoute), bz)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}

	cmd.Flags().String(flagProposalID, "", "proposalID of proposal voting on")
	cmd.Flags().String(flagVoter, "", "bech32 voter address")
	cmd.MarkFlagRequired(flagProposalID)
	cmd.MarkFlagRequired(flagVoter)
	return cmd
}

// GetCmdQueryVotes implements the command to query for proposal votes.
func GetCmdQueryVotes(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-votes",
		Short:   "query votes on a proposal",
		Example: "iriscli gov query-votes --proposal-id=1",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			proposalID := uint64(viper.GetInt64(flagProposalID))

			params := gov.QueryVotesParams{
				ProposalID: proposalID,
			}
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/votes", queryRoute), bz)
			if err != nil {
				return err
			}

			if res == nil {
				fmt.Printf("No one votes for the proposal [%v].\n", proposalID)
				return  nil
			}

			fmt.Println(string(res))
			return nil
		},
	}

	cmd.Flags().String(flagProposalID, "", "proposalID of which proposal's votes are being queried")
	cmd.MarkFlagRequired(flagProposalID)
	return cmd
}

// Command to Get a specific Deposit Information
// GetCmdQueryDeposit implements the query proposal deposit command.
func GetCmdQueryDeposit(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-deposit",
		Short:   "Query details of a deposit",
		Example: "iriscli gov query-deposit --proposal-id=1 --depositor=<depositor address>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			proposalID := uint64(viper.GetInt64(flagProposalID))

			depositorAddr, err := sdk.AccAddressFromBech32(viper.GetString(flagDepositor))
			if err != nil {
				return err
			}

			params := gov.QueryDepositParams{
				Depositor:  depositorAddr,
				ProposalID: proposalID,
			}
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/deposit", queryRoute), bz)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}

	cmd.Flags().String(flagProposalID, "", "proposalID of proposal deposited on")
	cmd.Flags().String(flagDepositor, "", "bech32 depositor address")
	cmd.MarkFlagRequired(flagProposalID)
	cmd.MarkFlagRequired(flagDeposit)
	return cmd
}

// GetCmdQueryDeposits implements the command to query for proposal deposits.
func GetCmdQueryDeposits(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-deposits",
		Short:   "Query deposits on a proposal",
		Example: "iriscli gov query-deposits --proposal-id=4",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			proposalID := uint64(viper.GetInt64(flagProposalID))

			params := gov.QueryDepositsParams{
				ProposalID: proposalID,
			}
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/deposits", queryRoute), bz)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}

	cmd.Flags().String(flagProposalID, "", "proposalID of which proposal's deposits are being queried")
	cmd.MarkFlagRequired(flagProposalID)
	return cmd
}

// GetCmdQueryDeposits implements the command to query for proposal deposits.
func GetCmdQueryTally(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-tally",
		Short:   "Get the tally of a proposal vote",
		Example: "iriscli gov query-tally --proposal-id=4",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			proposalID := uint64(viper.GetInt64(flagProposalID))

			params := gov.QueryTallyParams{
				ProposalID: proposalID,
			}
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/tally", queryRoute), bz)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}

	cmd.Flags().String(flagProposalID, "", "proposalID of which proposal is being tallied")

	return cmd
}

func GetCmdQueryGovConfig(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-params",
		Short:   "query parameter proposal's config",
		Example: "iriscli gov query-params --module=<module name> --key=<key name>",
		RunE: func(cmd *cobra.Command, args []string) error {
			moduleStr := viper.GetString(flagModule)
			keyStr := viper.GetString(flagKey)

			ctx := context.NewCLIContext().WithCodec(cdc)

			if moduleStr != "" {
				// There are four possible outputs if the --module parameter is not empty:
				// 1.List of the module;
				// 2.List of keys in the module;
				// 3.Error: GovParameter of the module does not exist;
				// 4.Error: The key in the module does not exist;
				res, err := ctx.QuerySubspace([]byte("Gov/"+moduleStr), storeName)
				if err == nil {
					if len(res) == 0 {
						// Return an error directly if the --module parameter is incorrect.
						return sdk.NewError(params.DefaultCodespace, params.CodeInvalidModule, fmt.Sprintf("The GovParameter of the module %s is not existed", moduleStr))
					}

					if keyStr != "" {
						// There are two possible outputs if the --key parameter is not empty:
						// 1.List of keys in the module;
						// 2.Error: The key in the module does not exist;
						params.RegisterGovParamMapping(&govparams.DepositProcedureParameter,
							&govparams.VotingProcedureParameter,
							&govparams.TallyingProcedureParameter)

						res, err := ctx.QueryStore([]byte(keyStr), storeName)
						return printKeyJsonIfExists(err, keyStr, res, cdc)
					}

					// Print module list
					err := printModuleList(res)
					if err != nil {
						return err
					}
					return nil
				} else {
					// Throw RPC client query exception
					return err
				}
			}

			// Check --key parameter if the --module parameter is empty.
			if keyStr != "" {
				// There are two possible outputs if the --key parameter is not empty:
				// 1.List of keys in the module;
				// 2.Error: The key in the module does not exist;
				params.RegisterGovParamMapping(&govparams.DepositProcedureParameter,
					&govparams.VotingProcedureParameter,
					&govparams.TallyingProcedureParameter)

				res, err := ctx.QueryStore([]byte(keyStr), storeName)
				return printKeyJsonIfExists(err, keyStr, res, cdc)
			}

			// Return error if the --module & --key parameters are all empty.
			return sdk.NewError(params.DefaultCodespace, params.CodeInvalidQueryParams, fmt.Sprintf("--module and --key can't both be empty"))
		},
	}

	cmd.Flags().String(flagModule, "", "module name")
	cmd.Flags().String(flagKey, "", "key name of parameter")
	return cmd
}

func printKeyJsonIfExists(e error, keyStr string, res []byte, cdc *codec.Codec) (err error) {
	if e == nil {
		if p, ok := params.ParamMapping[keyStr]; ok {
			if len(res) == 0 {
				// Return an error directly if the --key parameter is incorrect.
				return sdk.NewError(params.DefaultCodespace, params.CodeInvalidKey, fmt.Sprintf(keyStr+" is not existed"))
			}
			// Print key json in the module
			p.GetValueFromRawData(cdc, res)
			printParamStr(p, keyStr)
			return nil
		} else {
			//
			return sdk.NewError(params.DefaultCodespace, params.CodeInvalidKey, fmt.Sprintf(keyStr+" is not found"))
		}
	} else {
		// Throw RPC client query exception
		return e
	}
}

func printModuleList(res []sdk.KVPair) (err error) {
	var keys []string
	for _, kv := range res {
		keys = append(keys, string(kv.Key))
	}

	output, err := json.MarshalIndent(keys, "", " ")

	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

func printParamStr(p params.GovParameter, keyStr string) {
	var param govtypes.Param
	param.Key = keyStr
	param.Value = p.ToJson("")
	param.Op = ""
	jsonBytes, _ := json.Marshal(param)
	fmt.Println(string(jsonBytes))
}

func GetCmdPullGovConfig(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pull-params",
		Short:   "generate param.json file",
		Example: "iriscli gov pull-params",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx := context.NewCLIContext().WithCodec(cdc)
			res, err := ctx.QuerySubspace([]byte("Gov/"), storeName)
			if err == nil && len(res) != 0 {
				var pd govtypes.ParameterConfigFile
				pathStr := viper.GetString(flagPath)
				err := pd.WriteFile(cdc, res, pathStr)
				return err
			} else {
				fmt.Println("No GovParams can be found")
				return err
			}
		},
	}
	cmd.Flags().String(flagPath, app.DefaultCLIHome, "the directory of the param.json")
	return cmd
}
