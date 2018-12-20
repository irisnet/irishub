package gov

import (
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/modules/bank"
	"github.com/irisnet/irishub/types/gov/params"
	"github.com/tendermint/tendermint/crypto"
	"time"
	"github.com/irisnet/irishub/modules/params"
	"github.com/irisnet/irishub/modules/distribution"
	"github.com/irisnet/irishub/modules/guardian"
	protocolKeeper "github.com/irisnet/irishub/app/protocol/keeper"
	govtypes "github.com/irisnet/irishub/types/gov"
)

// nolint
var (
	DepositedCoinsAccAddr     = sdk.AccAddress(crypto.AddressHash([]byte("govDepositedCoins")))
	BurnedDepositCoinsAccAddr = sdk.AccAddress(crypto.AddressHash([]byte("govBurnedDepositCoins")))
)

// Governance Keeper
type Keeper struct {
	// The reference to the CoinKeeper to modify balances
	ck bank.Keeper

	dk distribution.Keeper

	gk guardian.Keeper
	// The ValidatorSet to get information about validators
	vs sdk.ValidatorSet

	// The reference to the DelegationSet to get information about delegators
	ds sdk.DelegationSet

	pk protocolKeeper.Keeper
	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey

	// The codec codec for binary encoding/decoding.
	cdc *codec.Codec

	// Reserved codespace
	codespace sdk.CodespaceType
}

// NewKeeper returns a governance keeper. It handles:
// - submitting governance proposals
// - depositing funds into proposals, and activating upon sufficient funds being deposited
// - users voting on proposals, with weight proportional to stake in the system
// - and tallying the result of the vote.
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, dk distribution.Keeper, ck bank.Keeper, gk guardian.Keeper, ds sdk.DelegationSet,pk protocolKeeper.Keeper, codespace sdk.CodespaceType) Keeper {
	return Keeper{
		storeKey:  key,
		ck:        ck,
		dk:        dk,
		gk:        gk,
		ds:        ds,
		vs:        ds.GetValidatorSet(),
		pk:        pk,
		cdc:       cdc,
		codespace: codespace,
	}
}

// =====================================================
// Proposals

////////////////////  iris begin  ///////////////////////////
func (keeper Keeper) NewProposal(ctx sdk.Context, title string, description string, proposalType govtypes.ProposalKind, param govtypes.Param) govtypes.Proposal {
	switch proposalType {
	case govtypes.ProposalTypeText:
		return keeper.NewTextProposal(ctx, title, description, proposalType)
	case govtypes.ProposalTypeParameterChange:
		return keeper.NewParametersProposal(ctx, title, description, proposalType, param)
	case govtypes.ProposalTypeSoftwareHalt:
		return keeper.NewHaltProposal(ctx, title, description, proposalType)
	}
	return nil
}

////////////////////  iris end  /////////////////////////////

// =====================================================
// Proposals

// Creates a NewProposal
func (keeper Keeper) NewTextProposal(ctx sdk.Context, title string, description string, proposalType govtypes.ProposalKind) govtypes.Proposal {
	proposalID, err := keeper.getNewProposalID(ctx)
	if err != nil {
		return nil
	}
	var proposal govtypes.Proposal = &govtypes.TextProposal{
		ProposalID:   proposalID,
		Title:        title,
		Description:  description,
		ProposalType: proposalType,
		Status:       govtypes.StatusDepositPeriod,
		TallyResult:  govtypes.EmptyTallyResult(),
		TotalDeposit: sdk.Coins{},
		SubmitTime:   ctx.BlockHeader().Time,
	}
	depositPeriod := govparams.GetDepositProcedure(ctx).MaxDepositPeriod
	proposal.SetDepositEndTime(proposal.GetSubmitTime().Add(depositPeriod))
	keeper.SetProposal(ctx, proposal)
	keeper.InsertInactiveProposalQueue(ctx, proposal.GetDepositEndTime(), proposalID)
	return proposal
}

////////////////////  iris begin  ///////////////////////////
func (keeper Keeper) NewParametersProposal(ctx sdk.Context, title string, description string, proposalType govtypes.ProposalKind, param govtypes.Param) govtypes.Proposal {
	proposalID, err := keeper.getNewProposalID(ctx)
	if err != nil {
		return nil
	}
	var textProposal = govtypes.TextProposal{
		ProposalID:   proposalID,
		Title:        title,
		Description:  description,
		ProposalType: proposalType,
		Status:       govtypes.StatusDepositPeriod,
		TallyResult:  govtypes.EmptyTallyResult(),
		TotalDeposit: sdk.Coins{},
		SubmitTime:   ctx.BlockHeader().Time,
	}

	param.Value = params.ParamMapping[param.Key].ToJson(param.Value)

	var proposal govtypes.Proposal = &govtypes.ParameterProposal{
		textProposal,
		param,
	}

	depositPeriod := govparams.GetDepositProcedure(ctx).MaxDepositPeriod
	proposal.SetDepositEndTime(proposal.GetSubmitTime().Add(depositPeriod))
	keeper.SetProposal(ctx, proposal)
	keeper.InsertInactiveProposalQueue(ctx, proposal.GetDepositEndTime(), proposalID)
	return proposal
}

func (keeper Keeper) NewHaltProposal(ctx sdk.Context, title string, description string, proposalType govtypes.ProposalKind) govtypes.Proposal {
	proposalID, err := keeper.getNewProposalID(ctx)
	if err != nil {
		return nil
	}
	var textProposal = govtypes.TextProposal{
		ProposalID:   proposalID,
		Title:        title,
		Description:  description,
		ProposalType: proposalType,
		Status:       govtypes.StatusDepositPeriod,
		TallyResult:  govtypes.EmptyTallyResult(),
		TotalDeposit: sdk.Coins{},
		SubmitTime:   ctx.BlockHeader().Time,
	}
	var proposal govtypes.Proposal = &govtypes.HaltProposal{
		textProposal,
	}

	depositPeriod := govparams.GetDepositProcedure(ctx).MaxDepositPeriod
	proposal.SetDepositEndTime(proposal.GetSubmitTime().Add(depositPeriod))
	keeper.SetProposal(ctx, proposal)
	keeper.InsertInactiveProposalQueue(ctx, proposal.GetDepositEndTime(), proposalID)
	return proposal
}

func (keeper Keeper) NewUsageProposal(ctx sdk.Context, msg MsgSubmitTxTaxUsageProposal) govtypes.Proposal {
	proposalID, err := keeper.getNewProposalID(ctx)
	if err != nil {
		return nil
	}
	var textProposal = govtypes.TextProposal{
		ProposalID:   proposalID,
		Title:        msg.Title,
		Description:  msg.Description,
		ProposalType: msg.ProposalType,
		Status:       govtypes.StatusDepositPeriod,
		TallyResult:  govtypes.EmptyTallyResult(),
		TotalDeposit: sdk.Coins{},
		SubmitTime:   ctx.BlockHeader().Time,
	}
	var proposal govtypes.Proposal = &govtypes.TaxUsageProposal{
		textProposal,
		msg.Usage,
		msg.DestAddress,
		msg.Percent,
	}
	keeper.saveProposal(ctx, proposal)
	return proposal
}

func (keeper Keeper) NewSoftwareUpgradeProposal(ctx sdk.Context, msg MsgSubmitSoftwareUpgradeProposal) govtypes.Proposal {
	proposalID, err := keeper.getNewProposalID(ctx)
	if err != nil {
		return nil
	}
	var textProposal = govtypes.TextProposal{
		ProposalID:   proposalID,
		Title:        msg.Title,
		Description:  msg.Description,
		ProposalType: msg.ProposalType,
		Status:       govtypes.StatusDepositPeriod,
		TallyResult:  govtypes.EmptyTallyResult(),
		TotalDeposit: sdk.Coins{},
		SubmitTime:   ctx.BlockHeader().Time,
	}
	var proposal govtypes.Proposal = &govtypes.SoftwareUpgradeProposal{
		textProposal,
		msg.Version,
		msg.Software,
		msg.SwitchHeight,
	}
	keeper.saveProposal(ctx, proposal)
	return proposal
}

func (keeper Keeper) saveProposal(ctx sdk.Context, proposal govtypes.Proposal) {
	depositPeriod := govparams.GetDepositProcedure(ctx).MaxDepositPeriod
	proposal.SetDepositEndTime(proposal.GetSubmitTime().Add(depositPeriod))
	keeper.SetProposal(ctx, proposal)
	keeper.InsertInactiveProposalQueue(ctx, proposal.GetDepositEndTime(), proposal.GetProposalID())
}

////////////////////  iris end  /////////////////////////////

// Get Proposal from store by ProposalID
func (keeper Keeper) GetProposal(ctx sdk.Context, proposalID uint64) govtypes.Proposal {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyProposal(proposalID))
	if bz == nil {
		return nil
	}

	var proposal govtypes.Proposal
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &proposal)

	return proposal
}

// Implements sdk.AccountKeeper.
func (keeper Keeper) SetProposal(ctx sdk.Context, proposal govtypes.Proposal) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(proposal)
	store.Set(KeyProposal(proposal.GetProposalID()), bz)
}

// Implements sdk.AccountKeeper.
func (keeper Keeper) DeleteProposal(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	proposal := keeper.GetProposal(ctx, proposalID)
	keeper.RemoveFromInactiveProposalQueue(ctx, proposal.GetDepositEndTime(), proposalID)
	keeper.RemoveFromActiveProposalQueue(ctx, proposal.GetVotingEndTime(), proposalID)
	store.Delete(KeyProposal(proposalID))
}

// Get Proposal from store by ProposalID
func (keeper Keeper) GetProposalsFiltered(ctx sdk.Context, voterAddr sdk.AccAddress, depositorAddr sdk.AccAddress, status govtypes.ProposalStatus, numLatest uint64) []govtypes.Proposal {

	maxProposalID, err := keeper.peekCurrentProposalID(ctx)
	if err != nil {
		return nil
	}

	matchingProposals := []govtypes.Proposal{}

	if numLatest == 0 || maxProposalID < numLatest {
		numLatest = maxProposalID
	}

	for proposalID := maxProposalID - numLatest; proposalID < maxProposalID; proposalID++ {
		if voterAddr != nil && len(voterAddr) != 0 {
			_, found := keeper.GetVote(ctx, proposalID, voterAddr)
			if !found {
				continue
			}
		}

		if depositorAddr != nil && len(depositorAddr) != 0 {
			_, found := keeper.GetDeposit(ctx, proposalID, depositorAddr)
			if !found {
				continue
			}
		}

		proposal := keeper.GetProposal(ctx, proposalID)
		if proposal == nil {
			continue
		}

		if govtypes.ValidProposalStatus(status) {
			if proposal.GetStatus() != status {
				continue
			}
		}

		matchingProposals = append(matchingProposals, proposal)
	}
	return matchingProposals
}

func (keeper Keeper) setInitialProposalID(ctx sdk.Context, proposalID uint64) sdk.Error {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyNextProposalID)
	if bz != nil {
		return govtypes.ErrInvalidGenesis(keeper.codespace, "Initial ProposalID already set")
	}
	bz = keeper.cdc.MustMarshalBinaryLengthPrefixed(proposalID)
	store.Set(KeyNextProposalID, bz)
	return nil
}

// Get the last used proposal ID
func (keeper Keeper) GetLastProposalID(ctx sdk.Context) (proposalID uint64) {
	proposalID, err := keeper.peekCurrentProposalID(ctx)
	if err != nil {
		return 0
	}
	proposalID--
	return
}

// Gets the next available ProposalID and increments it
func (keeper Keeper) getNewProposalID(ctx sdk.Context) (proposalID uint64, err sdk.Error) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyNextProposalID)
	if bz == nil {
		return 0, govtypes.ErrInvalidGenesis(keeper.codespace, "InitialProposalID never set")
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &proposalID)
	bz = keeper.cdc.MustMarshalBinaryLengthPrefixed(proposalID + 1)
	store.Set(KeyNextProposalID, bz)
	return proposalID, nil
}

// Peeks the next available ProposalID without incrementing it
func (keeper Keeper) peekCurrentProposalID(ctx sdk.Context) (proposalID uint64, err sdk.Error) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyNextProposalID)
	if bz == nil {
		return 0, govtypes.ErrInvalidGenesis(keeper.codespace, "InitialProposalID never set")
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &proposalID)
	return proposalID, nil
}

func (keeper Keeper) activateVotingPeriod(ctx sdk.Context, proposal govtypes.Proposal) {
	proposal.SetVotingStartTime(ctx.BlockHeader().Time)
	votingPeriod := govparams.GetVotingProcedure(ctx).VotingPeriod
	proposal.SetVotingEndTime(proposal.GetVotingStartTime().Add(votingPeriod))
	proposal.SetStatus(govtypes.StatusVotingPeriod)
	keeper.SetProposal(ctx, proposal)

	keeper.RemoveFromInactiveProposalQueue(ctx, proposal.GetDepositEndTime(), proposal.GetProposalID())
	keeper.InsertActiveProposalQueue(ctx, proposal.GetVotingEndTime(), proposal.GetProposalID())
}

// =====================================================
// Votes

// Adds a vote on a specific proposal
func (keeper Keeper) AddVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress, option govtypes.VoteOption) sdk.Error {
	proposal := keeper.GetProposal(ctx, proposalID)
	if proposal == nil {
		return govtypes.ErrUnknownProposal(keeper.codespace, proposalID)
	}
	if proposal.GetStatus() != govtypes.StatusVotingPeriod {
		return govtypes.ErrInactiveProposal(keeper.codespace, proposalID)
	}

	if !govtypes.ValidVoteOption(option) {
		return govtypes.ErrInvalidVote(keeper.codespace, option)
	}

	vote := govtypes.Vote{
		ProposalID: proposalID,
		Voter:      voterAddr,
		Option:     option,
	}
	keeper.setVote(ctx, proposalID, voterAddr, vote)

	return nil
}

// Gets the vote of a specific voter on a specific proposal
func (keeper Keeper) GetVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) (govtypes.Vote, bool) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyVote(proposalID, voterAddr))
	if bz == nil {
		return govtypes.Vote{}, false
	}
	var vote govtypes.Vote
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &vote)
	return vote, true
}

func (keeper Keeper) setVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress, vote govtypes.Vote) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(vote)
	store.Set(KeyVote(proposalID, voterAddr), bz)
}

// Gets all the votes on a specific proposal
func (keeper Keeper) GetVotes(ctx sdk.Context, proposalID uint64) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return sdk.KVStorePrefixIterator(store, KeyVotesSubspace(proposalID))
}

func (keeper Keeper) deleteVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(KeyVote(proposalID, voterAddr))
}

// =====================================================
// Deposits

// Gets the deposit of a specific depositor on a specific proposal
func (keeper Keeper) GetDeposit(ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress) (govtypes.Deposit, bool) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyDeposit(proposalID, depositorAddr))
	if bz == nil {
		return govtypes.Deposit{}, false
	}
	var deposit govtypes.Deposit
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &deposit)
	return deposit, true
}

func (keeper Keeper) setDeposit(ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress, deposit govtypes.Deposit) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(deposit)
	store.Set(KeyDeposit(proposalID, depositorAddr), bz)
}

// Adds or updates a deposit of a specific depositor on a specific proposal
// Activates voting period when appropriate
func (keeper Keeper) AddDeposit(ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress, depositAmount sdk.Coins) (sdk.Error, bool) {
	// Checks to see if proposal exists
	proposal := keeper.GetProposal(ctx, proposalID)
	if proposal == nil {
		return govtypes.ErrUnknownProposal(keeper.codespace, proposalID), false
	}

	// Check if proposal is still depositable
	if (proposal.GetStatus() != govtypes.StatusDepositPeriod) && (proposal.GetStatus() != govtypes.StatusVotingPeriod) {
		return govtypes.ErrAlreadyFinishedProposal(keeper.codespace, proposalID), false
	}

	// Send coins from depositor's account to DepositedCoinsAccAddr account
	_, err := keeper.ck.SendCoins(ctx, depositorAddr, DepositedCoinsAccAddr, depositAmount)
	if err != nil {
		return err, false
	}

	// Update Proposal
	proposal.SetTotalDeposit(proposal.GetTotalDeposit().Plus(depositAmount))
	keeper.SetProposal(ctx, proposal)

	// Check if deposit tipped proposal into voting period
	// Active voting period if so
	activatedVotingPeriod := false
	if proposal.GetStatus() == govtypes.StatusDepositPeriod && proposal.GetTotalDeposit().IsAllGTE(govparams.GetDepositProcedure(ctx).MinDeposit) {
		keeper.activateVotingPeriod(ctx, proposal)
		activatedVotingPeriod = true
	}

	// Add or update deposit object
	currDeposit, found := keeper.GetDeposit(ctx, proposalID, depositorAddr)
	if !found {
		newDeposit := govtypes.Deposit{depositorAddr, proposalID, depositAmount}
		keeper.setDeposit(ctx, proposalID, depositorAddr, newDeposit)
	} else {
		currDeposit.Amount = currDeposit.Amount.Plus(depositAmount)
		keeper.setDeposit(ctx, proposalID, depositorAddr, currDeposit)
	}

	return nil, activatedVotingPeriod
}

// Gets all the deposits on a specific proposal
func (keeper Keeper) GetDeposits(ctx sdk.Context, proposalID uint64) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return sdk.KVStorePrefixIterator(store, KeyDepositsSubspace(proposalID))
}

// Returns and deletes all the deposits on a specific proposal
func (keeper Keeper) RefundDeposits(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	depositsIterator := keeper.GetDeposits(ctx, proposalID)

	for ; depositsIterator.Valid(); depositsIterator.Next() {
		deposit := &govtypes.Deposit{}
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(depositsIterator.Value(), deposit)

		_, err := keeper.ck.SendCoins(ctx, DepositedCoinsAccAddr, deposit.Depositor, deposit.Amount)
		if err != nil {
			panic("should not happen")
		}

		store.Delete(depositsIterator.Key())
	}

	depositsIterator.Close()
}

// Deletes all the deposits on a specific proposal without refunding them
func (keeper Keeper) DeleteDeposits(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	depositsIterator := keeper.GetDeposits(ctx, proposalID)

	for ; depositsIterator.Valid(); depositsIterator.Next() {
		deposit := &govtypes.Deposit{}
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(depositsIterator.Value(), deposit)

		_, err := keeper.ck.SendCoins(ctx, DepositedCoinsAccAddr, BurnedDepositCoinsAccAddr, deposit.Amount)
		if err != nil {
			panic("should not happen")
		}

		store.Delete(depositsIterator.Key())
	}

	depositsIterator.Close()
}

// =====================================================
// ProposalQueues

// Returns an iterator for all the proposals in the Active Queue that expire by endTime
func (keeper Keeper) ActiveProposalQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return store.Iterator(PrefixActiveProposalQueue, sdk.PrefixEndBytes(PrefixActiveProposalQueueTime(endTime)))
}

// Inserts a ProposalID into the active proposal queue at endTime
func (keeper Keeper) InsertActiveProposalQueue(ctx sdk.Context, endTime time.Time, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(proposalID)
	store.Set(KeyActiveProposalQueueProposal(endTime, proposalID), bz)
}

// removes a proposalID from the Active Proposal Queue
func (keeper Keeper) RemoveFromActiveProposalQueue(ctx sdk.Context, endTime time.Time, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(KeyActiveProposalQueueProposal(endTime, proposalID))
}

// Returns an iterator for all the proposals in the Inactive Queue that expire by endTime
func (keeper Keeper) InactiveProposalQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return store.Iterator(PrefixInactiveProposalQueue, sdk.PrefixEndBytes(PrefixInactiveProposalQueueTime(endTime)))
}

// Inserts a ProposalID into the inactive proposal queue at endTime
func (keeper Keeper) InsertInactiveProposalQueue(ctx sdk.Context, endTime time.Time, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(proposalID)
	store.Set(KeyInactiveProposalQueueProposal(endTime, proposalID), bz)
}

// removes a proposalID from the Inactive Proposal Queue
func (keeper Keeper) RemoveFromInactiveProposalQueue(ctx sdk.Context, endTime time.Time, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(KeyInactiveProposalQueueProposal(endTime, proposalID))
}

func (keeper Keeper) GetTerminatorHeight(ctx sdk.Context) int64 {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyTerminatorHeight)
	if bz == nil {
		return -1
	}
	var height int64
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &height)

	return height
}

func (keeper Keeper) SetTerminatorHeight(ctx sdk.Context, height int64) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(height)
	store.Set(KeyTerminatorHeight, bz)
}

func (keeper Keeper) GetTerminatorPeriod(ctx sdk.Context) int64 {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyTerminatorPeriod)
	if bz == nil {
		return -1
	}
	var height int64
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &height)

	return height
}

func (keeper Keeper) SetTerminatorPeriod(ctx sdk.Context, height int64) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(height)
	store.Set(KeyTerminatorPeriod, bz)
}