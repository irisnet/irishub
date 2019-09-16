package cli

import (
	"fmt"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v1/gov"
	"github.com/irisnet/irishub/client/context"
	client "github.com/irisnet/irishub/client/gov"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetCmdQueryProposal implements the query proposal command.
func GetCmdQueryProposal(cdc *codec.Codec) *cobra.Command {
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

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/proposal", protocol.GovRoute), bz)
			if err != nil {
				return err
			}

			var proposal gov.Proposal
			err = cdc.UnmarshalJSON(res, &proposal)
			if err != nil {
				return err
			}
			return cliCtx.PrintOutput(proposal)
		},
	}

	cmd.Flags().String(flagProposalID, "", "proposalID of proposal being queried")
	cmd.MarkFlagRequired(flagProposalID)
	return cmd
}

// GetCmdQueryProposals implements a query proposals command.
func GetCmdQueryProposals(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-proposals",
		Short:   "Query proposals with optional filters",
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

			var status = ""
			if len(strProposalStatus) > 0 {
				status = client.NormalizeProposalStatus(strProposalStatus)
				if _, err := gov.ProposalStatusFromString(status); err != nil {
					return err
				}
			}
			params.ProposalStatus = status

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/proposals", protocol.GovRoute), bz)
			if err != nil {
				return err
			}

			var proposals gov.Proposals
			err = cdc.UnmarshalJSON(res, &proposals)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(proposals)
		},
	}

	cmd.Flags().String(flagNumLimit, "", "(optional) limit to latest [number] proposals. Defaults to all proposals")
	cmd.Flags().String(flagDepositor, "", "(optional) filter by proposals deposited on by depositor")
	cmd.Flags().String(flagVoter, "", "(optional) filter by proposals voted on by voted")
	cmd.Flags().String(flagStatus, "", "(optional) filter proposals by proposal status")

	return cmd
}

// GetCmdQueryVote implements the query proposal vote command.
func GetCmdQueryVote(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-vote",
		Short:   "Query vote",
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

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/vote", protocol.GovRoute), bz)
			if err != nil {
				return err
			}

			var vote gov.Vote
			if err := cdc.UnmarshalJSON(res, &vote); err != nil {
				return err
			}

			return cliCtx.PrintOutput(vote)
		},
	}

	cmd.Flags().String(flagProposalID, "", "proposalID of proposal voting on")
	cmd.Flags().String(flagVoter, "", "bech32 voter address")
	cmd.MarkFlagRequired(flagProposalID)
	cmd.MarkFlagRequired(flagVoter)
	return cmd
}

// GetCmdQueryVotes implements the command to query for proposal votes.
func GetCmdQueryVotes(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-votes",
		Short:   "Query votes on a proposal",
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

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/votes", protocol.GovRoute), bz)
			if err != nil {
				return err
			}

			var votes gov.Votes
			if err := cdc.UnmarshalJSON(res, &votes); err != nil {
				return err
			}

			return cliCtx.PrintOutput(votes)
		},
	}

	cmd.Flags().String(flagProposalID, "", "proposalID of which proposal's votes are being queried")
	cmd.MarkFlagRequired(flagProposalID)
	return cmd
}

// GetCmdQueryDeposit implements the query proposal deposit command.
func GetCmdQueryDeposit(cdc *codec.Codec) *cobra.Command {
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

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/deposit", protocol.GovRoute), bz)
			if err != nil {
				return err
			}

			var deposit gov.Deposit
			if err := cdc.UnmarshalJSON(res, &deposit); err != nil {
				return err
			}

			return cliCtx.PrintOutput(deposit)
		},
	}

	cmd.Flags().String(flagProposalID, "", "proposalID of proposal deposited on")
	cmd.Flags().String(flagDepositor, "", "bech32 depositor address")
	cmd.MarkFlagRequired(flagProposalID)
	cmd.MarkFlagRequired(flagDeposit)
	return cmd
}

// GetCmdQueryDeposits implements the command to query for proposal deposits.
func GetCmdQueryDeposits(cdc *codec.Codec) *cobra.Command {
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

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/deposits", protocol.GovRoute), bz)
			if err != nil {
				return err
			}

			var deposits gov.Deposits
			if err := cdc.UnmarshalJSON(res, &deposits); err != nil {
				return err
			}

			return cliCtx.PrintOutput(deposits)
		},
	}

	cmd.Flags().String(flagProposalID, "", "proposalID of which proposal's deposits are being queried")
	cmd.MarkFlagRequired(flagProposalID)
	return cmd
}

// GetCmdQueryTally implements the command to query for proposal deposits.
func GetCmdQueryTally(cdc *codec.Codec) *cobra.Command {
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

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/tally", protocol.GovRoute), bz)
			if err != nil {
				return err
			}

			var tally gov.TallyResult
			if err := cdc.UnmarshalJSON(res, &tally); err != nil {
				return err
			}

			return cliCtx.PrintOutput(tally)
		},
	}

	cmd.Flags().String(flagProposalID, "", "proposalID of which proposal is being tallied")

	return cmd
}
