package gov

import (
	"time"

	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/bank"
	"github.com/irisnet/irishub/modules/distribution"
	"github.com/irisnet/irishub/modules/guardian"
	stakeTypes "github.com/irisnet/irishub/modules/stake/types"
	sdk "github.com/irisnet/irishub/types"

	"github.com/irisnet/irishub/modules/params"
	"github.com/tendermint/tendermint/crypto"
)

// nolint
var (
	DepositedCoinsAccAddr = sdk.AccAddress(crypto.AddressHash([]byte("govDepositedCoins")))
	BurnRate              = sdk.NewDecWithPrec(2, 1)
	MinDepositRate        = sdk.NewDecWithPrec(3,1)
)

// Governance ProtocolKeeper
type Keeper struct {
	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey

	// The codec codec for binary encoding/decoding.
	cdc *codec.Codec

	// The reference to the Param ProtocolKeeper to get and set Global Params
	paramSpace   params.Subspace
	paramsKeeper params.Keeper

	protocolKeeper sdk.ProtocolKeeper

	// The reference to the CoinKeeper to modify balances
	ck bank.Keeper

	dk distribution.Keeper

	guardianKeeper guardian.Keeper

	// The ValidatorSet to get information about validators
	vs sdk.ValidatorSet

	// The reference to the DelegationSet to get information about delegators
	ds sdk.DelegationSet

	// Reserved codespace
	codespace sdk.CodespaceType
}

// NewProtocolKeeper returns a governance keeper. It handles:
// - submitting governance proposals
// - depositing funds into proposals, and activating upon sufficient funds being deposited
// - users voting on proposals, with weight proportional to stake in the system
// - and tallying the result of the vote.
func NewKeeper(key sdk.StoreKey, cdc *codec.Codec, paramSpace params.Subspace, paramsKeeper params.Keeper, protocolKeeper sdk.ProtocolKeeper, ck bank.Keeper, dk distribution.Keeper, guardianKeeper guardian.Keeper, ds sdk.DelegationSet, codespace sdk.CodespaceType) Keeper {
	return Keeper{
		key,
		cdc,
		paramSpace.WithTypeTable(ParamTypeTable()),
		paramsKeeper,
		protocolKeeper,
		ck,
		dk,
		guardianKeeper,
		ds.GetValidatorSet(),
		ds,
		codespace,
	}
}

// =====================================================
// Proposals

////////////////////  iris begin  ///////////////////////////
func (keeper Keeper) NewProposal(ctx sdk.Context, title string, description string, proposalType ProposalKind, param Params) Proposal {
	switch proposalType {
	case ProposalTypeParameterChange:
		return keeper.NewParametersProposal(ctx, title, description, proposalType, param)
	case ProposalTypeSystemHalt:
		return keeper.NewSystemHaltProposal(ctx, title, description, proposalType)
	}
	return nil
}

////////////////////  iris end  /////////////////////////////

// =====================================================
// Proposals

// Creates a NewProposal
////////////////////  iris begin  ///////////////////////////
func (keeper Keeper) NewParametersProposal(ctx sdk.Context, title string, description string, proposalType ProposalKind, params Params) Proposal {
	proposalID, err := keeper.getNewProposalID(ctx)
	if err != nil {
		return nil
	}
	var textProposal = TextProposal{
		ProposalID:   proposalID,
		Title:        title,
		Description:  description,
		ProposalType: proposalType,
		Status:       StatusDepositPeriod,
		TallyResult:  EmptyTallyResult(),
		TotalDeposit: sdk.Coins{},
		SubmitTime:   ctx.BlockHeader().Time,
	}

	var proposal Proposal = &ParameterProposal{
		textProposal,
		params,
	}

	depositPeriod := keeper.GetDepositProcedure(ctx, proposal).MaxDepositPeriod
	proposal.SetDepositEndTime(proposal.GetSubmitTime().Add(depositPeriod))
	keeper.SetProposal(ctx, proposal)
	keeper.InsertInactiveProposalQueue(ctx, proposal.GetDepositEndTime(), proposalID)
	return proposal
}

func (keeper Keeper) NewSystemHaltProposal(ctx sdk.Context, title string, description string, proposalType ProposalKind) Proposal {
	proposalID, err := keeper.getNewProposalID(ctx)
	if err != nil {
		return nil
	}
	var textProposal = TextProposal{
		ProposalID:   proposalID,
		Title:        title,
		Description:  description,
		ProposalType: proposalType,
		Status:       StatusDepositPeriod,
		TallyResult:  EmptyTallyResult(),
		TotalDeposit: sdk.Coins{},
		SubmitTime:   time.Now(),
	}
	var proposal Proposal = &SystemHaltProposal{
		textProposal,
	}

	depositPeriod := keeper.GetDepositProcedure(ctx, proposal).MaxDepositPeriod
	proposal.SetDepositEndTime(proposal.GetSubmitTime().Add(depositPeriod))
	keeper.SetProposal(ctx, proposal)
	keeper.InsertInactiveProposalQueue(ctx, proposal.GetDepositEndTime(), proposalID)
	return proposal
}

func (keeper Keeper) NewUsageProposal(ctx sdk.Context, msg MsgSubmitTxTaxUsageProposal) Proposal {
	proposalID, err := keeper.getNewProposalID(ctx)
	if err != nil {
		return nil
	}
	var textProposal = TextProposal{
		ProposalID:   proposalID,
		Title:        msg.Title,
		Description:  msg.Description,
		ProposalType: msg.ProposalType,
		Status:       StatusDepositPeriod,
		TallyResult:  EmptyTallyResult(),
		TotalDeposit: sdk.Coins{},
		SubmitTime:   ctx.BlockHeader().Time,
	}
	var proposal Proposal = &TaxUsageProposal{
		textProposal,
		TaxUsage{
			msg.Usage,
			msg.DestAddress,
			msg.Percent},
	}
	keeper.saveProposal(ctx, proposal)
	return proposal
}

func (keeper Keeper) NewSoftwareUpgradeProposal(ctx sdk.Context, msg MsgSubmitSoftwareUpgradeProposal) Proposal {
	proposalID, err := keeper.getNewProposalID(ctx)
	if err != nil {
		return nil
	}
	var textProposal = TextProposal{
		ProposalID:   proposalID,
		Title:        msg.Title,
		Description:  msg.Description,
		ProposalType: msg.ProposalType,
		Status:       StatusDepositPeriod,
		TallyResult:  EmptyTallyResult(),
		TotalDeposit: sdk.Coins{},
		SubmitTime:   ctx.BlockHeader().Time,
	}
	var proposal Proposal = &SoftwareUpgradeProposal{
		textProposal,
		sdk.ProtocolDefinition{
			msg.Version,
			msg.Software,
			msg.SwitchHeight,},
	}
	keeper.saveProposal(ctx, proposal)
	return proposal
}

func (keeper Keeper) saveProposal(ctx sdk.Context, proposal Proposal) {
	depositPeriod := keeper.GetDepositProcedure(ctx, proposal).MaxDepositPeriod
	proposal.SetDepositEndTime(proposal.GetSubmitTime().Add(depositPeriod))
	keeper.SetProposal(ctx, proposal)
	keeper.InsertInactiveProposalQueue(ctx, proposal.GetDepositEndTime(), proposal.GetProposalID())
}

// Get Proposal from store by ProposalID
func (keeper Keeper) GetProposal(ctx sdk.Context, proposalID uint64) Proposal {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyProposal(proposalID))
	if bz == nil {
		return nil
	}

	var proposal Proposal
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &proposal)

	return proposal
}

// Implements sdk.AccountKeeper.
func (keeper Keeper) SetProposal(ctx sdk.Context, proposal Proposal) {

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
func (keeper Keeper) GetProposalsFiltered(ctx sdk.Context, voterAddr sdk.AccAddress, depositorAddr sdk.AccAddress, status ProposalStatus, numLatest uint64) []Proposal {

	maxProposalID, err := keeper.peekCurrentProposalID(ctx)
	if err != nil {
		return nil
	}

	matchingProposals := []Proposal{}

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

		if ValidProposalStatus(status) {
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
		return ErrInvalidGenesis(keeper.codespace, "Initial ProposalID already set")
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
		return 0, ErrInvalidGenesis(keeper.codespace, "InitialProposalID never set")
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
		return 0, ErrInvalidGenesis(keeper.codespace, "InitialProposalID never set")
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &proposalID)
	return proposalID, nil
}

func (keeper Keeper) activateVotingPeriod(ctx sdk.Context, proposal Proposal) {
	proposal.SetVotingStartTime(ctx.BlockHeader().Time)
	votingPeriod := keeper.GetVotingProcedure(ctx, proposal).VotingPeriod
	proposal.SetVotingEndTime(proposal.GetVotingStartTime().Add(votingPeriod))
	proposal.SetStatus(StatusVotingPeriod)
	keeper.SetProposal(ctx, proposal)

	keeper.RemoveFromInactiveProposalQueue(ctx, proposal.GetDepositEndTime(), proposal.GetProposalID())
	keeper.InsertActiveProposalQueue(ctx, proposal.GetVotingEndTime(), proposal.GetProposalID())
	keeper.SetValidatorSet(ctx, proposal.GetProposalID())
}

// =====================================================
// Votes

// Adds a vote on a specific proposal
func (keeper Keeper) AddVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress, option VoteOption) sdk.Error {
	proposal := keeper.GetProposal(ctx, proposalID)
	if proposal == nil {
		return ErrUnknownProposal(keeper.codespace, proposalID)
	}
	if proposal.GetStatus() != StatusVotingPeriod {
		return ErrInactiveProposal(keeper.codespace, proposalID)
	}

	if keeper.vs.Validator(ctx, sdk.ValAddress(voterAddr)) == nil {
		return ErrOnlyValidatorVote(keeper.codespace, voterAddr)
	}

	if _, ok := keeper.GetVote(ctx, proposalID, voterAddr); ok {
		return ErrAlreadyVote(keeper.codespace, voterAddr, proposalID)
	}

	if !ValidVoteOption(option) {
		return ErrInvalidVote(keeper.codespace, option)
	}

	vote := Vote{
		ProposalID: proposalID,
		Voter:      voterAddr,
		Option:     option,
	}
	keeper.setVote(ctx, proposalID, voterAddr, vote)

	return nil
}

// Gets the vote of a specific voter on a specific proposal
func (keeper Keeper) GetVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) (Vote, bool) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyVote(proposalID, voterAddr))
	if bz == nil {
		return Vote{}, false
	}
	var vote Vote
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &vote)
	return vote, true
}

func (keeper Keeper) setVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress, vote Vote) {
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
func (keeper Keeper) GetDeposit(ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress) (Deposit, bool) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyDeposit(proposalID, depositorAddr))
	if bz == nil {
		return Deposit{}, false
	}
	var deposit Deposit
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &deposit)
	return deposit, true
}

func (keeper Keeper) setDeposit(ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress, deposit Deposit) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(deposit)
	store.Set(KeyDeposit(proposalID, depositorAddr), bz)
}

func (keeper Keeper) AddInitialDeposit(ctx sdk.Context, proposal Proposal, depositorAddr sdk.AccAddress, initialDeposit sdk.Coins) (sdk.Error, bool) {

	minDepositInt := sdk.NewDecFromInt(keeper.GetDepositProcedure(ctx, proposal).MinDeposit.AmountOf(stakeTypes.StakeDenom)).Mul(MinDepositRate).RoundInt()
	minInitialDeposit := sdk.Coins{sdk.NewCoin(stakeTypes.StakeDenom,minDepositInt)}
	if !initialDeposit.IsAllGTE(minInitialDeposit) {
		return ErrNotEnoughInitialDeposit(DefaultCodespace,initialDeposit,minInitialDeposit), false
	}

	return keeper.AddDeposit(ctx,proposal.GetProposalID(),depositorAddr, initialDeposit)
}

// Adds or updates a deposit of a specific depositor on a specific proposal
// Activates voting period when appropriate
func (keeper Keeper) AddDeposit(ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress, depositAmount sdk.Coins) (sdk.Error, bool) {
	// Checks to see if proposal exists
	proposal := keeper.GetProposal(ctx, proposalID)
	if proposal == nil {
		return ErrUnknownProposal(keeper.codespace, proposalID), false
	}

	// Check if proposal is still depositable
	if proposal.GetStatus() != StatusDepositPeriod {
		return ErrNotInDepositPeriod(keeper.codespace, proposalID), false
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
	if proposal.GetStatus() == StatusDepositPeriod && proposal.GetTotalDeposit().IsAllGTE(keeper.GetDepositProcedure(ctx, proposal).MinDeposit) {
		keeper.activateVotingPeriod(ctx, proposal)
		activatedVotingPeriod = true
	}

	// Add or update deposit object
	currDeposit, found := keeper.GetDeposit(ctx, proposalID, depositorAddr)
	if !found {
		newDeposit := Deposit{depositorAddr, proposalID, depositAmount}
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

func (keeper Keeper) RefundDepositsWithoutFee(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	depositsIterator := keeper.GetDeposits(ctx, proposalID)

	for ; depositsIterator.Valid(); depositsIterator.Next() {
		deposit := &Deposit{}
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(depositsIterator.Value(), deposit)

		_, err := keeper.ck.SendCoins(ctx, DepositedCoinsAccAddr, deposit.Depositor, deposit.Amount)
		if err != nil {
			panic("should not happen")
		}

		store.Delete(depositsIterator.Key())
	}

	depositsIterator.Close()
}

// Returns and deletes all the deposits on a specific proposal
func (keeper Keeper) RefundDeposits(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	depositsIterator := keeper.GetDeposits(ctx, proposalID)
	defer depositsIterator.Close()
	depositSum := sdk.Coins{}
	deposits := []*Deposit{}
	for ; depositsIterator.Valid(); depositsIterator.Next() {
		deposit := &Deposit{}
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(depositsIterator.Value(), deposit)
		deposits = append(deposits, deposit)
		depositSum = depositSum.Plus(deposit.Amount)
		store.Delete(depositsIterator.Key())
	}

	proposal := keeper.GetProposal(ctx, proposalID)
	BurnAmountDec := sdk.NewDecFromInt(keeper.GetDepositProcedure(ctx, proposal).MinDeposit.AmountOf(stakeTypes.StakeDenom)).Mul(BurnRate)
	DepositSumInt := depositSum.AmountOf(stakeTypes.StakeDenom)
	rate := BurnAmountDec.Quo(sdk.NewDecFromInt(DepositSumInt))
	RefundSumInt := sdk.NewInt(0)
	for _, deposit := range deposits {
		AmountDec := sdk.NewDecFromInt(deposit.Amount.AmountOf(stakeTypes.StakeDenom))
		RefundAmountInt := AmountDec.Sub(AmountDec.Mul(rate)).RoundInt()
		RefundSumInt = RefundSumInt.Add(RefundAmountInt)
		deposit.Amount = sdk.Coins{sdk.NewCoin(stakeTypes.StakeDenom, RefundAmountInt)}

		_, err := keeper.ck.SendCoins(ctx, DepositedCoinsAccAddr, deposit.Depositor, deposit.Amount)
		if err != nil {
			panic(err)
		}
	}

	_, err := keeper.ck.BurnCoinsFromAddr(ctx, DepositedCoinsAccAddr, sdk.Coins{sdk.NewCoin(stakeTypes.StakeDenom, DepositSumInt.Sub(RefundSumInt))})
	if err != nil {
		panic(err)
	}

}

// Deletes all the deposits on a specific proposal without refunding them
func (keeper Keeper) DeleteDeposits(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	depositsIterator := keeper.GetDeposits(ctx, proposalID)

	for ; depositsIterator.Valid(); depositsIterator.Next() {
		deposit := &Deposit{}
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(depositsIterator.Value(), deposit)

		_, err := keeper.ck.BurnCoinsFromAddr(ctx, DepositedCoinsAccAddr, deposit.Amount)
		if err != nil {
			panic(err)
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

func (keeper Keeper) GetSystemHaltHeight(ctx sdk.Context) int64 {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeySystemHaltHeight)
	if bz == nil {
		return -1
	}
	var height int64
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &height)

	return height
}

func (keeper Keeper) SetSystemHaltHeight(ctx sdk.Context, height int64) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(height)
	store.Set(KeySystemHaltHeight, bz)
}

func (keeper Keeper) GetCriticalProposalID(ctx sdk.Context) (uint64, bool) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyCriticalProposal)
	if bz == nil {
		return 0, false
	}
	var proposalID uint64
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &proposalID)
	return proposalID, true
}

func (keeper Keeper) SetCriticalProposalID(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(proposalID)
	store.Set(KeyCriticalProposal, bz)
}

func (keeper Keeper) GetCriticalProposalNum(ctx sdk.Context) uint64 {
	if _, ok := keeper.GetCriticalProposalID(ctx); ok {
		return 1
	}
	return 0
}

func (keeper Keeper) AddCriticalProposalNum(ctx sdk.Context, proposalID uint64) {
	keeper.SetCriticalProposalID(ctx, proposalID)
}

func (keeper Keeper) SubCriticalProposalNum(ctx sdk.Context) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(KeyCriticalProposal)
}

func (keeper Keeper) GetImportantProposalNum(ctx sdk.Context) uint64 {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyImportantProposalNum)
	if bz == nil {
		keeper.SetImportantProposalNum(ctx, 0)
		return 0
	}
	var num uint64
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &num)
	return num
}

func (keeper Keeper) SetImportantProposalNum(ctx sdk.Context, num uint64) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(num)
	store.Set(KeyImportantProposalNum, bz)
}

func (keeper Keeper) AddImportantProposalNum(ctx sdk.Context) {
	keeper.SetImportantProposalNum(ctx, keeper.GetImportantProposalNum(ctx)+1)
}

func (keeper Keeper) SubImportantProposalNum(ctx sdk.Context) {
	keeper.SetImportantProposalNum(ctx, keeper.GetImportantProposalNum(ctx)-1)
}

func (keeper Keeper) GetNormalProposalNum(ctx sdk.Context) uint64 {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyNormalProposalNum)
	if bz == nil {
		keeper.SetImportantProposalNum(ctx, 0)
		return 0
	}
	var num uint64
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &num)
	return num
}

func (keeper Keeper) SetNormalProposalNum(ctx sdk.Context, num uint64) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(num)
	store.Set(KeyNormalProposalNum, bz)
}

func (keeper Keeper) AddNormalProposalNum(ctx sdk.Context) {
	keeper.SetNormalProposalNum(ctx, keeper.GetNormalProposalNum(ctx)+1)
}

func (keeper Keeper) SubNormalProposalNum(ctx sdk.Context) {
	keeper.SetNormalProposalNum(ctx, keeper.GetNormalProposalNum(ctx)-1)
}

func (keeper Keeper) HasReachedTheMaxProposalNum(ctx sdk.Context, pl ProposalLevel) (uint64, bool) {
	maxNum := keeper.GetMaxNumByProposalLevel(ctx, pl)
	switch pl {
	case ProposalLevelCritical:
		return keeper.GetCriticalProposalNum(ctx), keeper.GetCriticalProposalNum(ctx) == maxNum
	case ProposalLevelImportant:
		return keeper.GetImportantProposalNum(ctx), keeper.GetImportantProposalNum(ctx) == maxNum
	case ProposalLevelNormal:
		return keeper.GetNormalProposalNum(ctx), keeper.GetNormalProposalNum(ctx) == maxNum
	default:
		panic("There is no level for this proposal")
	}
}

func (keeper Keeper) AddProposalNum(ctx sdk.Context, p Proposal) {
	switch GetProposalLevel(p) {
	case ProposalLevelCritical:
		keeper.AddCriticalProposalNum(ctx, p.GetProposalID())
	case ProposalLevelImportant:
		keeper.AddImportantProposalNum(ctx)
	case ProposalLevelNormal:
		keeper.AddNormalProposalNum(ctx)
	default:
		panic("There is no level for this proposal which type is " + p.GetProposalType().String())
	}
}

func (keeper Keeper) SubProposalNum(ctx sdk.Context, p Proposal) {
	switch GetProposalLevel(p) {
	case ProposalLevelCritical:
		keeper.SubCriticalProposalNum(ctx)
	case ProposalLevelImportant:
		keeper.SubImportantProposalNum(ctx)
	case ProposalLevelNormal:
		keeper.SubNormalProposalNum(ctx)
	default:
		panic("There is no level for this proposal which type is " + p.GetProposalType().String())
	}
}

func (keeper Keeper) SetValidatorSet(ctx sdk.Context, proposalID uint64) {

	valAddrs := []sdk.ValAddress{}
	keeper.vs.IterateBondedValidatorsByPower(ctx, func(index int64, validator sdk.Validator) (stop bool) {
		valAddrs = append(valAddrs, validator.GetOperator())
		return false
	})
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(valAddrs)
	store.Set(KeyValidatorSet(proposalID), bz)
}

func (keeper Keeper) GetValidatorSet(ctx sdk.Context, proposalID uint64) []sdk.ValAddress {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyValidatorSet(proposalID))
	if bz == nil {
		return []sdk.ValAddress{}
	}
	valAddrs := []sdk.ValAddress{}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &valAddrs)
	return valAddrs
}

func (keeper Keeper) DeleteValidatorSet(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(KeyValidatorSet(proposalID))
}
