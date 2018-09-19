package cli

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/irisnet/irishub/client/context"
	govClient "github.com/irisnet/irishub/client/gov"
	"github.com/irisnet/irishub/modules/gov"
	"github.com/irisnet/irishub/modules/gov/params"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetCmdQueryProposal implements the query proposal command.
func GetCmdQueryProposal(storeName string, cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-proposal",
		Short: "query proposal details",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			proposalID := viper.GetInt64(flagProposalID)

			res, err := cliCtx.QueryStore(gov.KeyProposal(proposalID), storeName)
			if len(res) == 0 || err != nil {
				return fmt.Errorf("proposalID [%d] is not existed", proposalID)
			}

			var proposal gov.Proposal
			cdc.MustUnmarshalBinary(res, &proposal)

			proposalResponse, err := govClient.ConvertProposalCoins(cliCtx, proposal)
			if err != nil {
				return err
			}
			output, err := wire.MarshalJSONIndent(cdc, proposalResponse)
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
		},
	}

	cmd.Flags().String(flagProposalID, "", "proposalID of proposal being queried")

	return cmd
}

// nolint: gocyclo
// GetCmdQueryProposals implements a query proposals command.
func GetCmdQueryProposals(storeName string, cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-proposals",
		Short: "query proposals with optional filters",
		RunE: func(cmd *cobra.Command, args []string) error {
			bechDepositerAddr := viper.GetString(flagDepositer)
			bechVoterAddr := viper.GetString(flagVoter)
			strProposalStatus := viper.GetString(flagStatus)
			latestProposalsIDs := viper.GetInt64(flagLatestProposalIDs)

			var err error
			var voterAddr sdk.AccAddress
			var depositerAddr sdk.AccAddress
			var proposalStatus gov.ProposalStatus

			if len(bechDepositerAddr) != 0 {
				depositerAddr, err = sdk.AccAddressFromBech32(bechDepositerAddr)
				if err != nil {
					return err
				}
			}

			if len(bechVoterAddr) != 0 {
				voterAddr, err = sdk.AccAddressFromBech32(bechVoterAddr)
				if err != nil {
					return err
				}
			}

			if len(strProposalStatus) != 0 {
				proposalStatus, err = gov.ProposalStatusFromString(strProposalStatus)
				if err != nil {
					return err
				}
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryStore(gov.KeyNextProposalID, storeName)
			if err != nil {
				return err
			}
			var maxProposalID int64
			cdc.MustUnmarshalBinary(res, &maxProposalID)

			matchingProposals := []govClient.TextProposalResponse{}

			if latestProposalsIDs == 0 {
				latestProposalsIDs = maxProposalID
			}

			for proposalID := maxProposalID - latestProposalsIDs; proposalID < maxProposalID; proposalID++ {
				if voterAddr != nil {
					res, err = cliCtx.QueryStore(gov.KeyVote(proposalID, voterAddr), storeName)
					if err != nil || len(res) == 0 {
						continue
					}
				}

				if depositerAddr != nil {
					res, err = cliCtx.QueryStore(gov.KeyDeposit(proposalID, depositerAddr), storeName)
					if err != nil || len(res) == 0 {
						continue
					}
				}

				res, err = cliCtx.QueryStore(gov.KeyProposal(proposalID), storeName)
				if err != nil || len(res) == 0 {
					continue
				}

				var proposal gov.Proposal
				cdc.MustUnmarshalBinary(res, &proposal)

				if len(strProposalStatus) != 0 {
					if proposal.GetStatus() != proposalStatus {
						continue
					}
				}

				proposalResponse, err := govClient.ConvertProposalCoins(cliCtx, proposal)
				if err != nil {
					return err
				}

				matchingProposals = append(matchingProposals, proposalResponse)
			}

			if len(matchingProposals) == 0 {
				fmt.Println("No matching proposals found")
				return nil
			}

			for _, proposal := range matchingProposals {
				fmt.Printf("  %d - %s\n", proposal.ProposalID, proposal.Title)
			}

			return nil
		},
	}

	cmd.Flags().String(flagLatestProposalIDs, "", "(optional) limit to latest [number] proposals. Defaults to all proposals")
	cmd.Flags().String(flagDepositer, "", "(optional) filter by proposals deposited on by depositer")
	cmd.Flags().String(flagVoter, "", "(optional) filter by proposals voted on by voted")
	cmd.Flags().String(flagStatus, "", "(optional) filter proposals by proposal status")

	return cmd
}

// Command to Get a Proposal Information
// GetCmdQueryVote implements the query proposal vote command.
func GetCmdQueryVote(storeName string, cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-vote",
		Short: "query vote",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			proposalID := viper.GetInt64(flagProposalID)

			voterAddr, err := sdk.AccAddressFromBech32(viper.GetString(flagVoter))
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryStore(gov.KeyVote(proposalID, voterAddr), storeName)
			if len(res) == 0 || err != nil {
				return fmt.Errorf("proposalID [%d] does not exist", proposalID)
			}

			var vote gov.Vote
			cdc.MustUnmarshalBinary(res, &vote)

			output, err := wire.MarshalJSONIndent(cdc, vote)
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
		},
	}

	cmd.Flags().String(flagProposalID, "", "proposalID of proposal voting on")
	cmd.Flags().String(flagVoter, "", "bech32 voter address")

	return cmd
}

// GetCmdQueryVotes implements the command to query for proposal votes.
func GetCmdQueryVotes(storeName string, cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-votes",
		Short: "query votes on a proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			proposalID := viper.GetInt64(flagProposalID)

			res, err := cliCtx.QueryStore(gov.KeyProposal(proposalID), storeName)
			if len(res) == 0 || err != nil {
				return fmt.Errorf("proposalID [%d] does not exist", proposalID)
			}

			var proposal gov.Proposal
			cdc.MustUnmarshalBinary(res, &proposal)

			if proposal.GetStatus() != gov.StatusVotingPeriod {
				fmt.Println("Proposal not in voting period.")
				return nil
			}

			res2, err := cliCtx.QuerySubspace(gov.KeyVotesSubspace(proposalID), storeName)
			if err != nil {
				return err
			}

			var votes []gov.Vote
			for i := 0; i < len(res2); i++ {
				var vote gov.Vote
				cdc.MustUnmarshalBinary(res2[i].Value, &vote)
				votes = append(votes, vote)
			}

			output, err := wire.MarshalJSONIndent(cdc, votes)
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
		},
	}

	cmd.Flags().String(flagProposalID, "", "proposalID of which proposal's votes are being queried")

	return cmd
}

var (
	flagModule = "module"
	flagKey    = "key"
)

type Param struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Op    string `json:"op"`
}

func GetCmdQueryGovConfig(storeName string, cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-params",
		Short: "query parameter proposal's config",
		RunE: func(cmd *cobra.Command, args []string) error {
			moduleStr := viper.GetString(flagModule)
			keyStr := viper.GetString(flagKey)

			ctx := context.NewCLIContext().WithCodec(cdc)

			if moduleStr != "" {
				res, err := ctx.QuerySubspace([]byte("Gov/"+moduleStr), storeName)
				if err == nil {
					var keys []string
					for _, kv := range res {
						keys = append(keys, string(kv.Key))
					}
					output, err := wire.MarshalJSONIndent(cdc, keys)
					//cmn.WriteFile(,output,644)
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
				res, err := ctx.QueryStore([]byte(keyStr), storeName)
				if err == nil {
					switch keyStr {
					case "Gov/gov/depositProcedure":
						var p govparams.DepositProcedure
						cdc.MustUnmarshalBinary(res, &p)
						ToParamStr(cdc, p, keyStr)
					}
				}

			}

			return nil
		},
	}

	cmd.Flags().String(flagModule, "", "the module of parameter ")
	cmd.Flags().String(flagKey, "", "the key of parameter")
	return cmd
}
func ToParamStr(cdc *wire.Codec, p interface{}, keyStr string) {
	var param Param
	param.Key = keyStr
	param.Value = ToJson(cdc, p)
	param.Op = ""
	jsonBytes, _ := cdc.MarshalJSON(param)
	//jsonBytes, _ := wire.MarshalJSONIndent(cdc, param)
	fmt.Println(string(jsonBytes))
}

func ToJson(cdc *wire.Codec, p interface{}) string {
	jsonBytes, _ := cdc.MarshalJSON(p)
	return string(jsonBytes)
}
