package gov

import (
	"encoding/json"
	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v1/gov/tags"
	"time"

	"github.com/irisnet/irishub/app/v1/bank"
	"github.com/irisnet/irishub/app/v1/distribution"
	stakeTypes "github.com/irisnet/irishub/app/v1/stake/types"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/guardian"
	sdk "github.com/irisnet/irishub/types"

	"github.com/irisnet/irishub/app/v1/params"
	"strconv"
)

// nolint
var (
	BurnRate       = sdk.NewDecWithPrec(2, 1)
	MinDepositRate = sdk.NewDecWithPrec(3, 1)
)

// Governance ProtocolKeeper
type Keeper struct {
	// The (unexposed) keys used to access the stores from the Content.
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

	metrics *Metrics

	ak AssetKeeper
}

// NewProtocolKeeper returns a governance keeper. It handles:
// - submitting governance proposals
// - depositing funds into proposals, and activating upon sufficient funds being deposited
// - users voting on proposals, with weight proportional to stake in the system
// - and tallying the result of the vote.
func NewKeeper(key sdk.StoreKey, cdc *codec.Codec, paramSpace params.Subspace, paramsKeeper params.Keeper, protocolKeeper sdk.ProtocolKeeper, ck bank.Keeper, dk distribution.Keeper, guardianKeeper guardian.Keeper, ds sdk.DelegationSet, codespace sdk.CodespaceType, metrics *Metrics, ak AssetKeeper) Keeper {
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
		metrics,
		ak,
	}
}

// =====================================================
// Proposals
func (keeper Keeper) SubmitProposal(ctx sdk.Context, msg sdk.Msg) (sdk.Tags, sdk.Error) {
	content := msg.(Content)

	// construct a proposal
	proposal, err := content.GetProposalType().NewProposal(content)
	if err != nil {
		return nil, err
	}
	pLevel := proposal.GetProposalLevel()

	// validate MinInitialDeposit
	initialDeposit := content.GetInitialDeposit()
	minInitialDeposit := keeper.getMinInitialDeposit(ctx, pLevel)
	if !initialDeposit.IsAllGTE(minInitialDeposit) {
		return nil, ErrNotEnoughInitialDeposit(DefaultCodespace, initialDeposit, minInitialDeposit)
	}

	// validate proposal
	if err := proposal.Validate(ctx, keeper, true); err != nil {
		return nil, err
	}

	//fill proposal field
	proposalID, err := keeper.getNewProposalID(ctx)
	if err != nil {
		return nil, err
	}
	proposal.SetProposalID(proposalID)
	proposal.SetSubmitTime(ctx.BlockHeader().Time)
	depositPeriod := keeper.GetDepositProcedure(ctx, pLevel).MaxDepositPeriod
	proposal.SetDepositEndTime(proposal.GetSubmitTime().Add(depositPeriod))

	// store proposal
	keeper.SetProposal(ctx, proposal)
	// put proposal to InactiveProposalQueue
	keeper.InsertInactiveProposalQueue(ctx, proposal.GetDepositEndTime(), proposal.GetProposalID())

	//store init deposit
	err, votingStarted := keeper.AddDeposit(ctx, proposal.GetProposalID(), proposal.GetProposer(), initialDeposit)
	if err != nil {
		return nil, err
	}

	//add proposal number
	keeper.AddProposalNum(ctx, pLevel, proposal.GetProposalID())

	//return tag
	proposalIDBytes := []byte(strconv.FormatUint(proposal.GetProposalID(), 10))
	resTags := sdk.NewTags(
		tags.Proposer, []byte(proposal.GetProposer().String()),
		tags.ProposalID, proposalIDBytes,
	)
	if votingStarted {
		resTags = resTags.AppendTag(tags.VotingPeriodStart, proposalIDBytes)
	}

	switch proposal.GetProposalType() {
	case ProposalTypeParameter:
		var paramBytes []byte
		paramBytes, _ = json.Marshal(content.GetParams())
		resTags = resTags.AppendTag(tags.Param, paramBytes)
	case ProposalTypeCommunityTaxUsage:
		msg := msg.(MsgSubmitCommunityTaxUsageProposal)
		resTags = resTags.AppendTag(tags.DestAddress, []byte(msg.DestAddress.String()))
	case ProposalTypeTokenAddition:
		tokenId := proposal.(*TokenAdditionProposal).FToken.GetUniqueID()
		resTags = resTags.AppendTag(tags.TokenId, []byte(tokenId))
	}
	return resTags, nil
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
	votingPeriod := keeper.GetVotingProcedure(ctx, proposal.GetProposalLevel()).VotingPeriod
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

	validator := keeper.vs.Validator(ctx, sdk.ValAddress(voterAddr))
	if validator == nil {
		isDelegator := false
		keeper.ds.IterateDelegations(ctx, voterAddr, func(index int64, delegation sdk.Delegation) (stop bool) {
			isDelegator = true
			return isDelegator
		})
		if !isDelegator {
			return ErrOnlyValidatorOrDelegatorVote(keeper.codespace, voterAddr)
		}
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
	if validator != nil {
		keeper.metrics.AddVote(validator.GetOperator().String(), proposalID, option)
	}
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
	ctx.CoinFlowTags().AppendCoinFlowTag(ctx, depositorAddr.String(), auth.GovDepositCoinsAccAddr.String(), depositAmount.String(), sdk.GovDepositFlow, "")
	_, err := keeper.ck.SendCoins(ctx, depositorAddr, auth.GovDepositCoinsAccAddr, depositAmount)
	if err != nil {
		return err, false
	}

	// Update Proposal
	proposal.SetTotalDeposit(proposal.GetTotalDeposit().Add(depositAmount))
	keeper.SetProposal(ctx, proposal)

	// Check if deposit tipped proposal into voting period
	// Active voting period if so
	activatedVotingPeriod := false
	if proposal.GetTotalDeposit().IsAllGTE(keeper.GetDepositProcedure(ctx, proposal.GetProposalLevel()).MinDeposit) {
		keeper.activateVotingPeriod(ctx, proposal)
		activatedVotingPeriod = true
	}

	// Add or update deposit object
	currDeposit, found := keeper.GetDeposit(ctx, proposalID, depositorAddr)
	if !found {
		newDeposit := Deposit{depositorAddr, proposalID, depositAmount}
		keeper.setDeposit(ctx, proposalID, depositorAddr, newDeposit)
	} else {
		currDeposit.Amount = currDeposit.Amount.Add(depositAmount)
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
	defer depositsIterator.Close()
	depositSum := sdk.Coins{}
	deposits := []*Deposit{}
	for ; depositsIterator.Valid(); depositsIterator.Next() {
		deposit := &Deposit{}
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(depositsIterator.Value(), deposit)
		deposits = append(deposits, deposit)
		depositSum = depositSum.Add(deposit.Amount)
		store.Delete(depositsIterator.Key())
	}

	proposal := keeper.GetProposal(ctx, proposalID)
	BurnAmountDec := sdk.NewDecFromInt(keeper.GetDepositProcedure(ctx, proposal.GetProposalLevel()).MinDeposit.AmountOf(stakeTypes.StakeDenom)).Mul(BurnRate)
	DepositSumInt := depositSum.AmountOf(stakeTypes.StakeDenom)
	rate := BurnAmountDec.Quo(sdk.NewDecFromInt(DepositSumInt))
	RefundSumInt := sdk.NewInt(0)
	for _, deposit := range deposits {
		AmountDec := sdk.NewDecFromInt(deposit.Amount.AmountOf(stakeTypes.StakeDenom))
		RefundAmountInt := AmountDec.Sub(AmountDec.Mul(rate)).RoundInt()
		RefundSumInt = RefundSumInt.Add(RefundAmountInt)
		deposit.Amount = sdk.Coins{sdk.NewCoin(stakeTypes.StakeDenom, RefundAmountInt)}

		ctx.CoinFlowTags().AppendCoinFlowTag(ctx, auth.GovDepositCoinsAccAddr.String(), deposit.Depositor.String(), deposit.Amount.String(), sdk.GovDepositRefundFlow, "")
		_, err := keeper.ck.SendCoins(ctx, auth.GovDepositCoinsAccAddr, deposit.Depositor, deposit.Amount)
		if err != nil {
			panic(err)
		}
	}

	burnCoin := sdk.NewCoin(stakeTypes.StakeDenom, DepositSumInt.Sub(RefundSumInt))
	ctx.CoinFlowTags().AppendCoinFlowTag(ctx, auth.GovDepositCoinsAccAddr.String(), "", burnCoin.String(), sdk.GovDepositBurnFlow, "")
	_, err := keeper.ck.BurnCoins(ctx, auth.GovDepositCoinsAccAddr, sdk.Coins{burnCoin})
	if err != nil {
		panic(err)
	}

}

// Deletes all the deposits on a specific proposal without refunding them
func (keeper Keeper) DeleteDeposits(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	depositsIterator := keeper.GetDeposits(ctx, proposalID)
	defer depositsIterator.Close()

	for ; depositsIterator.Valid(); depositsIterator.Next() {
		deposit := &Deposit{}
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(depositsIterator.Value(), deposit)

		ctx.CoinFlowTags().AppendCoinFlowTag(ctx, auth.GovDepositCoinsAccAddr.String(), "", deposit.Amount.String(), sdk.GovDepositBurnFlow, "")
		_, err := keeper.ck.BurnCoins(ctx, auth.GovDepositCoinsAccAddr, deposit.Amount)
		if err != nil {
			panic(err)
		}

		store.Delete(depositsIterator.Key())
	}
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
	keeper.metrics.SetProposalStatus(proposalID, StatusVotingPeriod)
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
	keeper.metrics.SetProposalStatus(proposalID, StatusDepositPeriod)
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

func (keeper Keeper) GetSystemHaltPeriod(ctx sdk.Context) (SystemHaltPeriod int64) {
	keeper.paramSpace.Get(ctx, KeySystemHaltPeriod, &SystemHaltPeriod)
	return
}

// get inflation params from the global param store
func (keeper Keeper) GetParamSet(ctx sdk.Context) GovParams {
	var params GovParams
	keeper.paramSpace.GetParamSet(ctx, &params)
	return params
}

// set inflation params from the global param store
func (keeper Keeper) SetParamSet(ctx sdk.Context, params GovParams) {
	keeper.paramSpace.SetParamSet(ctx, &params)
}

func (keeper Keeper) getMinInitialDeposit(ctx sdk.Context, proposalLevel ProposalLevel) sdk.Coins {
	minDeposit := keeper.GetDepositProcedure(ctx, proposalLevel).MinDeposit
	minDepositInt := sdk.NewDecFromInt(minDeposit.AmountOf(stakeTypes.StakeDenom)).Mul(MinDepositRate).RoundInt()
	return sdk.Coins{sdk.NewCoin(stakeTypes.StakeDenom, minDepositInt)}
}

// Returns the current Deposit Procedure from the global param store
func (keeper Keeper) GetDepositProcedure(ctx sdk.Context, p ProposalLevel) DepositProcedure {
	params := keeper.GetParamSet(ctx)
	switch p {
	case ProposalLevelCritical:
		return DepositProcedure{
			MinDeposit:       params.CriticalMinDeposit,
			MaxDepositPeriod: params.CriticalDepositPeriod,
		}
	case ProposalLevelImportant:
		return DepositProcedure{
			MinDeposit:       params.ImportantMinDeposit,
			MaxDepositPeriod: params.ImportantDepositPeriod,
		}
	case ProposalLevelNormal:
		return DepositProcedure{
			MinDeposit:       params.NormalMinDeposit,
			MaxDepositPeriod: params.NormalDepositPeriod,
		}
	default:
		panic("There is no level for this proposal which type is " + p.string())
	}
}

// Returns the current Voting Procedure from the global param store
func (keeper Keeper) GetVotingProcedure(ctx sdk.Context, p ProposalLevel) VotingProcedure {
	params := keeper.GetParamSet(ctx)
	switch p {
	case ProposalLevelCritical:
		return VotingProcedure{
			VotingPeriod: params.CriticalVotingPeriod,
		}
	case ProposalLevelImportant:
		return VotingProcedure{
			VotingPeriod: params.ImportantVotingPeriod,
		}
	case ProposalLevelNormal:
		return VotingProcedure{
			VotingPeriod: params.NormalVotingPeriod,
		}
	default:
		panic("There is no level for this proposal which type is " + p.string())
	}
}

func (keeper Keeper) GetMaxNumByProposalLevel(ctx sdk.Context, p ProposalLevel) uint64 {
	params := keeper.GetParamSet(ctx)
	switch p {
	case ProposalLevelCritical:
		return params.CriticalMaxNum

	case ProposalLevelImportant:
		return params.ImportantMaxNum

	case ProposalLevelNormal:
		return params.NormalMaxNum
	default:
		panic("There is no level for this proposal which type is " + p.string())
	}
}

// Returns the current Tallying Procedure from the global param store
func (keeper Keeper) GetTallyingProcedure(ctx sdk.Context, p ProposalLevel) TallyingProcedure {
	params := keeper.GetParamSet(ctx)
	switch p {
	case ProposalLevelCritical:
		return TallyingProcedure{
			Threshold:     params.CriticalThreshold,
			Veto:          params.CriticalVeto,
			Participation: params.CriticalParticipation,
			Penalty:       params.CriticalPenalty,
		}
	case ProposalLevelImportant:
		return TallyingProcedure{
			Threshold:     params.ImportantThreshold,
			Veto:          params.ImportantVeto,
			Participation: params.ImportantParticipation,
			Penalty:       params.ImportantPenalty,
		}
	case ProposalLevelNormal:
		return TallyingProcedure{
			Threshold:     params.NormalThreshold,
			Veto:          params.NormalVeto,
			Participation: params.NormalParticipation,
			Penalty:       params.NormalPenalty,
		}
	default:
		panic("There is no level for this proposal which type is " + p.string())
	}
}

func (keeper Keeper) AddProposalNum(ctx sdk.Context, p ProposalLevel, args ...interface{}) {
	switch p {
	case ProposalLevelCritical:
		proposalID := args[0].(uint64)
		keeper.AddCriticalProposalNum(ctx, proposalID)
	case ProposalLevelImportant:
		keeper.AddImportantProposalNum(ctx)
	case ProposalLevelNormal:
		keeper.AddNormalProposalNum(ctx)
	default:
		panic("There is no level for this proposal which type is " + p.string())
	}
}

func (keeper Keeper) SubProposalNum(ctx sdk.Context, p ProposalLevel) {
	switch p {
	case ProposalLevelCritical:
		keeper.SubCriticalProposalNum(ctx)
	case ProposalLevelImportant:
		keeper.SubImportantProposalNum(ctx)
	case ProposalLevelNormal:
		keeper.SubNormalProposalNum(ctx)
	default:
		panic("There is no level for this proposal which type is " + p.string())
	}
}

func (keeper Keeper) HasReachedTheMaxProposalNum(ctx sdk.Context, p ProposalLevel) (uint64, bool) {
	ctx.Logger().Debug("Proposals Distribution",
		"CriticalProposalNum", keeper.GetCriticalProposalNum(ctx),
		"ImportantProposalNum", keeper.GetImportantProposalNum(ctx),
		"NormalProposalNum", keeper.GetNormalProposalNum(ctx))

	maxNum := keeper.GetMaxNumByProposalLevel(ctx, p)
	switch p {
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

func (keeper Keeper) Init(ctx sdk.Context) {
	if sdk.NetworkType != sdk.Mainnet {
		return
	}
	keeper.SetParamSet(ctx, DefaultParams())
}
