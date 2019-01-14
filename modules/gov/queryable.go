package gov

import (
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"time"
)

// query endpoints supported by the governance Querier
const (
	QueryProposals = "proposals"
	QueryProposal  = "proposal"
	QueryDeposits  = "deposits"
	QueryDeposit   = "deposit"
	QueryVotes     = "votes"
	QueryVote      = "vote"
	QueryTally     = "tally"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryProposals:
			return queryProposals(ctx, path[1:], req, keeper)
		case QueryProposal:
			return queryProposal(ctx, path[1:], req, keeper)
		case QueryDeposits:
			return queryDeposits(ctx, path[1:], req, keeper)
		case QueryDeposit:
			return queryDeposit(ctx, path[1:], req, keeper)
		case QueryVotes:
			return queryVotes(ctx, path[1:], req, keeper)
		case QueryVote:
			return queryVote(ctx, path[1:], req, keeper)
		case QueryTally:
			return queryTally(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown gov query endpoint")
		}
	}
}

type ProposalOutput struct {
	ProposalID   uint64       `json:"proposal_id"`   //  ID of the proposal
	Title        string       `json:"title"`         //  Title of the proposal
	Description  string       `json:"description"`   //  Description of the proposal
	ProposalType ProposalKind `json:"proposal_type"` //  Type of proposal. Initial set {PlainTextProposal, SoftwareUpgradeProposal}

	Status      ProposalStatus `json:"proposal_status"` //  Status of the Proposal {Pending, Active, Passed, Rejected}
	TallyResult TallyResult    `json:"tally_result"`    //  Result of Tallys

	SubmitTime     time.Time `json:"submit_time"`      //  Time of the block where TxGovSubmitProposal was included
	DepositEndTime time.Time `json:"deposit_end_time"` // Time that the Proposal would expire if deposit amount isn't met
	TotalDeposit   sdk.Coins `json:"total_deposit"`    //  Current deposit on this proposal. Initial value is set at InitialDeposit

	VotingStartTime time.Time              `json:"voting_start_time"` //  Time of the block where MinDeposit was reached. -1 if MinDeposit is not reached
	VotingEndTime   time.Time              `json:"voting_end_time"`   // Time that the VotingPeriod for this proposal will end and votes will be tallied
	Params          Params                 `json:"param"`
	TaxUsage        TaxUsage               `json:"tax_usage"`
	Upgrade         sdk.ProtocolDefinition `json:"protocol_definition"`
}

type ProposalOutputs []ProposalOutput

func ConvertProposalToProposalOutput(proposal Proposal) ProposalOutput {

	proposalOutput := ProposalOutput{
		ProposalID:   proposal.GetProposalID(),
		Title:        proposal.GetTitle(),
		Description:  proposal.GetDescription(),
		ProposalType: proposal.GetProposalType(),

		Status:      proposal.GetStatus(),
		TallyResult: proposal.GetTallyResult(),

		SubmitTime:     proposal.GetSubmitTime(),
		DepositEndTime: proposal.GetDepositEndTime(),
		TotalDeposit:   proposal.GetTotalDeposit(),

		VotingStartTime: proposal.GetVotingStartTime(),
		VotingEndTime:   proposal.GetVotingEndTime(),
		Params:          Params{},
		TaxUsage:        proposal.GetTaxUsage(),
		Upgrade:         proposal.GetProtocolDefinition(),
	}

	if proposal.GetProposalType() == ProposalTypeParameterChange {
		proposalOutput.Params = proposal.(*ParameterProposal).Params
	}
	return proposalOutput
}

func ConvertProposalsToProposalOutputs(proposals []Proposal) ProposalOutputs {

	var proposalOutputs ProposalOutputs
	for _, proposal := range proposals {
		proposalOutputs = append(proposalOutputs, ConvertProposalToProposalOutput(proposal))
	}
	return proposalOutputs
}

// Params for query 'custom/gov/proposal'
type QueryProposalParams struct {
	ProposalID uint64
}

// nolint: unparam
func queryProposal(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	var params QueryProposalParams
	err2 := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err2 != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err2.Error()))
	}

	proposal := keeper.GetProposal(ctx, params.ProposalID)
	if proposal == nil {
		return nil, ErrUnknownProposal(DefaultCodespace, params.ProposalID)
	}

	proposalOutput := ConvertProposalToProposalOutput(proposal)

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, proposalOutput)
	if err2 != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err2.Error()))
	}
	return bz, nil
}

// Params for query 'custom/gov/deposit'
type QueryDepositParams struct {
	ProposalID uint64
	Depositor  sdk.AccAddress
}

// nolint: unparam
func queryDeposit(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	var params QueryDepositParams
	err2 := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err2 != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err2.Error()))
	}

	proposal := keeper.GetProposal(ctx, params.ProposalID)
	if proposal == nil {
		return nil, ErrUnknownProposal(DefaultCodespace, params.ProposalID)
	}

	if proposal.GetStatus() == StatusPassed || proposal.GetStatus() == StatusRejected {
		return nil, ErrCodeDepositDeleted(DefaultCodespace, params.ProposalID)
	}

	deposit, bool := keeper.GetDeposit(ctx, params.ProposalID, params.Depositor)
	if !bool {
		return nil, ErrCodeDepositNotExisted(DefaultCodespace, params.Depositor, params.ProposalID)
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, deposit)
	if err2 != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err2.Error()))
	}
	return bz, nil
}

// Params for query 'custom/gov/vote'
type QueryVoteParams struct {
	ProposalID uint64
	Voter      sdk.AccAddress
}

// nolint: unparam
func queryVote(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	var params QueryVoteParams
	err2 := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err2 != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err2.Error()))
	}

	proposal := keeper.GetProposal(ctx, params.ProposalID)
	if proposal == nil {
		return nil, ErrUnknownProposal(DefaultCodespace, params.ProposalID)
	}

	vote, bool := keeper.GetVote(ctx, params.ProposalID, params.Voter)
	if !bool {
		return nil, ErrCodeVoteNotExisted(DefaultCodespace, params.Voter, params.ProposalID)
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, vote)
	if err2 != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err2.Error()))
	}
	return bz, nil
}

// Params for query 'custom/gov/deposits'
type QueryDepositsParams struct {
	ProposalID uint64
}

// nolint: unparam
func queryDeposits(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	var params QueryDepositsParams
	err2 := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err2 != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err2.Error()))
	}

	proposal := keeper.GetProposal(ctx, params.ProposalID)
	if proposal == nil {
		return nil, ErrUnknownProposal(DefaultCodespace, params.ProposalID)
	}

	if proposal.GetStatus() == StatusPassed || proposal.GetStatus() == StatusRejected {
		return nil, ErrCodeDepositDeleted(DefaultCodespace, params.ProposalID)
	}

	var deposits []Deposit
	depositsIterator := keeper.GetDeposits(ctx, params.ProposalID)
	defer depositsIterator.Close()
	for ; depositsIterator.Valid(); depositsIterator.Next() {
		deposit := Deposit{}
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(depositsIterator.Value(), &deposit)
		deposits = append(deposits, deposit)
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, deposits)
	if err2 != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err2.Error()))
	}
	return bz, nil
}

// Params for query 'custom/gov/votes'
type QueryVotesParams struct {
	ProposalID uint64
}

// nolint: unparam
func queryVotes(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	var params QueryVotesParams
	err2 := keeper.cdc.UnmarshalJSON(req.Data, &params)

	if err2 != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err2.Error()))
	}

	proposal := keeper.GetProposal(ctx, params.ProposalID)
	if proposal == nil {
		return nil, ErrUnknownProposal(DefaultCodespace, params.ProposalID)
	}

	var votes []Vote
	votesIterator := keeper.GetVotes(ctx, params.ProposalID)
	defer votesIterator.Close()
	for ; votesIterator.Valid(); votesIterator.Next() {
		vote := Vote{}
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(votesIterator.Value(), &vote)
		votes = append(votes, vote)
	}

	if len(votes) == 0 {
		return nil, nil
	}
	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, votes)
	if err2 != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err2.Error()))
	}
	return bz, nil
}

// Params for query 'custom/gov/proposals'
type QueryProposalsParams struct {
	Voter          sdk.AccAddress
	Depositor      sdk.AccAddress
	ProposalStatus ProposalStatus
	Limit          uint64
}

// nolint: unparam
func queryProposals(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	var params QueryProposalsParams
	err2 := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err2 != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err2.Error()))
	}

	proposals := keeper.GetProposalsFiltered(ctx, params.Voter, params.Depositor, params.ProposalStatus, params.Limit)

	proposalOutputs := ConvertProposalsToProposalOutputs(proposals)

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, proposalOutputs)
	if err2 != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err2.Error()))
	}
	return bz, nil
}

// Params for query 'custom/gov/tally'
type QueryTallyParams struct {
	ProposalID uint64
}

// nolint: unparam
func queryTally(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	// TODO: Dependant on #1914

	var param QueryTallyParams
	err2 := keeper.cdc.UnmarshalJSON(req.Data, &param)
	if err2 != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err2.Error()))
	}

	proposal := keeper.GetProposal(ctx, param.ProposalID)
	if proposal == nil {
		return nil, ErrUnknownProposal(DefaultCodespace, param.ProposalID)
	}

	var tallyResult TallyResult

	if proposal.GetStatus() == StatusDepositPeriod {
		tallyResult = EmptyTallyResult()
	} else if proposal.GetStatus() == StatusPassed || proposal.GetStatus() == StatusRejected {
		tallyResult = proposal.GetTallyResult()
	} else {
		_, tallyResult, _ = tally(ctx, keeper, proposal)
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, tallyResult)
	if err2 != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err2.Error()))
	}
	return bz, nil
}
