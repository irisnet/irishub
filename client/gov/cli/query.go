package cli

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/client/context"
	client "github.com/irisnet/irishub/client/gov"
	"github.com/irisnet/irishub/iparam"
	"github.com/irisnet/irishub/modules/gov"
	"github.com/irisnet/irishub/modules/gov/params"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
				proposalStatus, err := gov.ProposalStatusFromString(client.NormalizeProposalStatus(strProposalStatus))
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

			fmt.Println(string(res))
			return nil
		},
	}

	cmd.Flags().String(flagProposalID, "", "proposalID of which proposal's votes are being queried")

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
				res, err := ctx.QuerySubspace([]byte("Gov/"+moduleStr), storeName)
				if err == nil {

					if len(res) == 0 {
						return sdk.NewError(iparam.DefaultCodespace, iparam.CodeInvalidModule, fmt.Sprintf("The GovParameter of the module %s is not existed", moduleStr))
					}

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
				} else {
					return nil
				}
			}

			if keyStr != "" {
				iparam.RegisterGovParamMapping(&govparams.DepositProcedureParameter,
					&govparams.VotingProcedureParameter,
					&govparams.TallyingProcedureParameter)

				res, err := ctx.QueryStore([]byte(keyStr), storeName)
				if err == nil {
					if p, ok := iparam.ParamMapping[keyStr]; ok {

						if len(res) == 0 {
							return sdk.NewError(iparam.DefaultCodespace, iparam.CodeInvalidKey, fmt.Sprintf(keyStr+" is not existed"))
						}

						p.GetValueFromRawData(cdc, res)
						printParamStr(p, keyStr)
						return nil

					} else {
						return sdk.NewError(iparam.DefaultCodespace, iparam.CodeInvalidKey, fmt.Sprintf(keyStr+" is not found"))
					}
				} else {
					return err
				}

			}
			return sdk.NewError(iparam.DefaultCodespace, iparam.CodeInvalidQueryParams, fmt.Sprintf("--module and --key can't both be empty"))
		},
	}

	cmd.Flags().String(flagModule, "", "module name")
	cmd.Flags().String(flagKey, "", "key name of parameter")
	return cmd
}

func printParamStr(p iparam.GovParameter, keyStr string) {
	var param gov.Param
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
				var pd gov.ParameterConfigFile
				pathStr := viper.GetString(flagPath)
				err := pd.WriteFile(cdc, res, pathStr)
				return err
			} else {
				fmt.Println("No GovParams can be found")
				return err
			}
		},
	}
	cmd.Flags().String(flagPath, app.DefaultNodeHome, "directory of iris home")
	return cmd
}
