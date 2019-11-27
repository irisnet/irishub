package types

import (
	"fmt"
	"testing"

	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/auth"
	"github.com/irisnet/irishub/modules/bank"
	"github.com/irisnet/irishub/store"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

func setupMultiStore() (sdk.MultiStore, *sdk.KVStoreKey) {
	db := dbm.NewMemDB()
	authKey := sdk.NewKVStoreKey("authkey")
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	return ms, authKey
}

func TestValidatorEqual(t *testing.T) {
	val1 := NewValidator(addr1, pk1, Description{})
	val2 := NewValidator(addr1, pk1, Description{})

	ok := val1.Equal(val2)
	require.True(t, ok)

	val2 = NewValidator(addr2, pk2, Description{})

	ok = val1.Equal(val2)
	require.False(t, ok)
}

func TestUpdateDescription(t *testing.T) {
	d1 := Description{
		Website: "https://validator.cosmos",
		Details: "Test validator",
	}

	d2 := Description{
		Moniker:  DoNotModifyDesc,
		Identity: DoNotModifyDesc,
		Website:  DoNotModifyDesc,
		Details:  DoNotModifyDesc,
	}

	d3 := Description{
		Moniker:  "",
		Identity: "",
		Website:  "",
		Details:  "",
	}

	d, err := d1.UpdateDescription(d2)
	require.Nil(t, err)
	require.Equal(t, d, d1)

	d, err = d1.UpdateDescription(d3)
	require.Nil(t, err)
	require.Equal(t, d, d3)
}

func TestABCIValidatorUpdate(t *testing.T) {
	validator := NewValidator(addr1, pk1, Description{})

	abciVal := validator.ABCIValidatorUpdate()
	require.Equal(t, tmtypes.TM2PB.PubKey(validator.ConsPubKey), abciVal.PubKey)
	require.Equal(t, validator.BondedTokens().RoundInt64(), abciVal.Power)
}

func TestABCIValidatorUpdateZero(t *testing.T) {
	validator := NewValidator(addr1, pk1, Description{})

	abciVal := validator.ABCIValidatorUpdateZero()
	require.Equal(t, tmtypes.TM2PB.PubKey(validator.ConsPubKey), abciVal.PubKey)
	require.Equal(t, int64(0), abciVal.Power)
}

func TestRemoveTokens(t *testing.T) {

	ms, authKey := setupMultiStore()
	cdc := codec.New()
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	accountKeeper := auth.NewAccountKeeper(cdc, authKey, auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper)

	validator := Validator{
		OperatorAddr:    addr1,
		ConsPubKey:      pk1,
		Status:          sdk.Bonded,
		Tokens:          sdk.NewDec(100),
		DelegatorShares: sdk.NewDec(100),
	}

	bondedPool := InitialBondedPool()
	poolA := Pool{
		BondedPool: bondedPool,
		BankKeeper: bankKeeper,
	}
	poolA.BondedPool.BondedTokens = validator.BondedTokens()
	poolA.BankKeeper.IncreaseLoosenToken(ctx, sdk.Coins{sdk.NewCoin(StakeDenom, sdk.NewInt(10))})

	validator, poolA = validator.UpdateStatus(ctx, poolA, sdk.Bonded)
	require.Equal(t, sdk.Bonded, validator.Status)

	// remove tokens and test check everything
	validator, poolA = validator.RemoveTokens(ctx, poolA, sdk.NewDec(10))
	require.Equal(t, int64(90), validator.Tokens.RoundInt64())
	require.Equal(t, int64(90), poolA.BondedPool.BondedTokens.RoundInt64())
	require.Equal(t, int64(20), poolA.BankKeeper.GetLoosenCoins(ctx).AmountOf(StakeDenom).Int64())

	// update validator to unbonded and remove some more tokens
	validator, poolA = validator.UpdateStatus(ctx, poolA, sdk.Unbonded)
	require.Equal(t, sdk.Unbonded, validator.Status)
	require.Equal(t, int64(0), poolA.BondedPool.BondedTokens.RoundInt64())
	require.Equal(t, int64(110), poolA.BankKeeper.GetLoosenCoins(ctx).AmountOf(StakeDenom).Int64())

	validator, poolA = validator.RemoveTokens(ctx, poolA, sdk.NewDec(10))
	require.Equal(t, int64(80), validator.Tokens.RoundInt64())
	require.Equal(t, int64(0), poolA.BondedPool.BondedTokens.RoundInt64())
	require.Equal(t, int64(110), poolA.BankKeeper.GetLoosenCoins(ctx).AmountOf(StakeDenom).Int64())
}

func TestAddTokensValidatorBonded(t *testing.T) {

	ms, authKey := setupMultiStore()
	cdc := codec.New()
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	accountKeeper := auth.NewAccountKeeper(cdc, authKey, auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper)

	bondedPool := InitialBondedPool()
	poolA := Pool{
		BondedPool: bondedPool,
		BankKeeper: bankKeeper,
	}
	poolA.BankKeeper.IncreaseLoosenToken(ctx, sdk.Coins{sdk.NewCoin(StakeDenom, sdk.NewInt(10))})

	validator := NewValidator(addr1, pk1, Description{})
	validator, poolA = validator.UpdateStatus(ctx, poolA, sdk.Bonded)
	validator, poolA, delShares := validator.AddTokensFromDel(ctx, poolA, sdk.NewInt(10))

	require.Equal(t, sdk.OneDec(), validator.DelegatorShareExRate())

	assert.True(sdk.DecEq(t, sdk.NewDec(10), delShares))
	assert.True(sdk.DecEq(t, sdk.NewDec(10), validator.BondedTokens()))
}

func TestAddTokensValidatorUnbonding(t *testing.T) {

	ms, authKey := setupMultiStore()
	cdc := codec.New()
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	accountKeeper := auth.NewAccountKeeper(cdc, authKey, auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper)

	bondedPool := InitialBondedPool()
	poolA := Pool{
		BondedPool: bondedPool,
		BankKeeper: bankKeeper,
	}
	poolA.BankKeeper.IncreaseLoosenToken(ctx, sdk.Coins{sdk.NewCoin(StakeDenom, sdk.NewInt(10))})

	validator := NewValidator(addr1, pk1, Description{})
	validator, poolA = validator.UpdateStatus(ctx, poolA, sdk.Unbonding)
	validator, poolA, delShares := validator.AddTokensFromDel(ctx, poolA, sdk.NewInt(10))

	require.Equal(t, sdk.OneDec(), validator.DelegatorShareExRate())

	assert.True(sdk.DecEq(t, sdk.NewDec(10), delShares))
	assert.Equal(t, sdk.Unbonding, validator.Status)
	assert.True(sdk.DecEq(t, sdk.NewDec(10), validator.Tokens))
}

func TestAddTokensValidatorUnbonded(t *testing.T) {

	ms, authKey := setupMultiStore()
	cdc := codec.New()
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	accountKeeper := auth.NewAccountKeeper(cdc, authKey, auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper)

	bondedPool := InitialBondedPool()
	poolA := Pool{
		BondedPool: bondedPool,
		BankKeeper: bankKeeper,
	}
	poolA.BankKeeper.IncreaseLoosenToken(ctx, sdk.Coins{sdk.NewCoin(StakeDenom, sdk.NewInt(10))})

	validator := NewValidator(addr1, pk1, Description{})
	validator, poolA = validator.UpdateStatus(ctx, poolA, sdk.Unbonded)
	validator, poolA, delShares := validator.AddTokensFromDel(ctx, poolA, sdk.NewInt(10))

	require.Equal(t, sdk.OneDec(), validator.DelegatorShareExRate())

	assert.True(sdk.DecEq(t, sdk.NewDec(10), delShares))
	assert.Equal(t, sdk.Unbonded, validator.Status)
	assert.True(sdk.DecEq(t, sdk.NewDec(10), validator.Tokens))
}

// TODO refactor to make simpler like the AddToken tests above
func TestRemoveDelShares(t *testing.T) {
	ms, authKey := setupMultiStore()
	cdc := codec.New()
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	accountKeeper := auth.NewAccountKeeper(cdc, authKey, auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper)

	valA := Validator{
		OperatorAddr:    addr1,
		ConsPubKey:      pk1,
		Status:          sdk.Bonded,
		Tokens:          sdk.NewDec(100),
		DelegatorShares: sdk.NewDec(100),
	}
	bondedPool := InitialBondedPool()
	poolA := Pool{
		BondedPool: bondedPool,
		BankKeeper: bankKeeper,
	}
	poolA.BondedPool.BondedTokens = valA.BondedTokens()
	poolA.BankKeeper.IncreaseLoosenToken(ctx, sdk.Coins{sdk.NewCoin(StakeDenom, sdk.NewInt(10))})
	require.Equal(t, valA.DelegatorShareExRate(), sdk.OneDec())

	// Remove delegator shares
	valB, poolB, coinsB := valA.RemoveDelShares(ctx, poolA, sdk.NewDec(10))
	assert.Equal(t, int64(10), coinsB.RoundInt64())
	assert.Equal(t, int64(90), valB.DelegatorShares.RoundInt64())
	assert.Equal(t, int64(90), valB.BondedTokens().RoundInt64())
	assert.Equal(t, int64(90), poolB.BondedPool.BondedTokens.RoundInt64())
	assert.Equal(t, int64(20), poolB.BankKeeper.GetLoosenCoins(ctx).AmountOf(StakeDenom).Int64())

	poolA.BondedPool.BondedTokens = poolA.BondedPool.BondedTokens.Sub(sdk.NewDec(10))
	// conservation of tokens
	require.True(sdk.DecEq(t,
		sdk.NewDecFromInt(poolB.BankKeeper.GetLoosenCoins(ctx).AmountOf(StakeDenom)).Add(poolB.BondedPool.BondedTokens),
		sdk.NewDecFromInt(poolA.BankKeeper.GetLoosenCoins(ctx).AmountOf(StakeDenom)).Add(poolA.BondedPool.BondedTokens)))

	// specific case from random tests
	poolTokens := sdk.NewDec(5102)
	delShares := sdk.NewDec(115)
	validator := Validator{
		OperatorAddr:    addr1,
		ConsPubKey:      pk1,
		Status:          sdk.Bonded,
		Tokens:          poolTokens,
		DelegatorShares: delShares,
	}
	poolC := Pool{
		BondedPool: BondedPool{
			BondedTokens: sdk.NewDec(248305),
		},
		BankKeeper: bankKeeper,
	}
	shares := sdk.NewDec(29)
	_, newPool, tokens := validator.RemoveDelShares(ctx, poolC, shares)

	exp, err := sdk.NewDecFromStr("1286")
	require.NoError(t, err)

	require.True(sdk.DecEq(t, exp, tokens))
	poolC.BondedPool.BondedTokens = poolC.BondedPool.BondedTokens.Sub(sdk.NewDec(1286))
	require.True(sdk.DecEq(t,
		sdk.NewDecFromInt(newPool.BankKeeper.GetLoosenCoins(ctx).AmountOf(StakeDenom)).Add(newPool.BondedPool.BondedTokens),
		sdk.NewDecFromInt(poolC.BankKeeper.GetLoosenCoins(ctx).AmountOf(StakeDenom)).Add(poolC.BondedPool.BondedTokens)))
}

func TestUpdateStatus(t *testing.T) {
	ms, authKey := setupMultiStore()
	cdc := codec.New()
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	accountKeeper := auth.NewAccountKeeper(cdc, authKey, auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper)

	bondedPool := InitialBondedPool()
	pool := Pool{
		BondedPool: bondedPool,
		BankKeeper: bankKeeper,
	}
	pool.BankKeeper.IncreaseLoosenToken(ctx, sdk.Coins{sdk.NewCoin(StakeDenom, sdk.NewInt(100))})

	validator := NewValidator(addr1, pk1, Description{})
	validator, pool, _ = validator.AddTokensFromDel(ctx, pool, sdk.NewInt(100))
	require.Equal(t, sdk.Unbonded, validator.Status)
	require.Equal(t, int64(100), validator.Tokens.RoundInt64())
	require.Equal(t, int64(0), pool.BondedPool.BondedTokens.RoundInt64())
	require.Equal(t, int64(100), pool.BankKeeper.GetLoosenCoins(ctx).AmountOf(StakeDenom).Int64())

	validator, pool = validator.UpdateStatus(ctx, pool, sdk.Bonded)
	require.Equal(t, sdk.Bonded, validator.Status)
	require.Equal(t, int64(100), validator.Tokens.RoundInt64())
	require.Equal(t, int64(100), pool.BondedPool.BondedTokens.RoundInt64())
	require.Equal(t, int64(0), pool.BankKeeper.GetLoosenCoins(ctx).AmountOf(StakeDenom).Int64())

	validator, pool = validator.UpdateStatus(ctx, pool, sdk.Unbonding)
	require.Equal(t, sdk.Unbonding, validator.Status)
	require.Equal(t, int64(100), validator.Tokens.RoundInt64())
	require.Equal(t, int64(0), pool.BondedPool.BondedTokens.RoundInt64())
	require.Equal(t, int64(100), pool.BankKeeper.GetLoosenCoins(ctx).AmountOf(StakeDenom).Int64())
}

func TestPossibleOverflow(t *testing.T) {
	ms, authKey := setupMultiStore()
	cdc := codec.New()
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	accountKeeper := auth.NewAccountKeeper(cdc, authKey, auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper)

	poolTokens := sdk.NewDec(2159)
	delShares := sdk.NewDec(391432570689183511).Quo(sdk.NewDec(40113011844664))
	validator := Validator{
		OperatorAddr:    addr1,
		ConsPubKey:      pk1,
		Status:          sdk.Bonded,
		Tokens:          poolTokens,
		DelegatorShares: delShares,
	}
	bondedPool := InitialBondedPool()
	poolA := Pool{
		BondedPool: bondedPool,
		BankKeeper: bankKeeper,
	}
	poolA.BondedPool.BondedTokens = poolTokens
	poolA.BankKeeper.IncreaseLoosenToken(ctx, sdk.Coins{sdk.NewCoin(StakeDenom, poolTokens.TruncateInt())})

	tokens := int64(71)
	msg := fmt.Sprintf("validator %#v", validator)
	newValidator, _, _ := validator.AddTokensFromDel(ctx, poolA, sdk.NewInt(tokens))

	msg = fmt.Sprintf("Added %d tokens to %s", tokens, msg)
	require.False(t, newValidator.DelegatorShareExRate().LT(sdk.ZeroDec()),
		"Applying operation \"%s\" resulted in negative DelegatorShareExRate(): %v",
		msg, newValidator.DelegatorShareExRate())
}

func TestHumanReadableString(t *testing.T) {
	validator := NewValidator(addr1, pk1, Description{})

	// NOTE: Being that the validator's keypair is random, we cannot test the
	// actual contents of the string.
	valStr, err := validator.HumanReadableString()
	require.Nil(t, err)
	require.NotEmpty(t, valStr)
}

func TestValidatorMarshalUnmarshalJSON(t *testing.T) {
	validator := NewValidator(addr1, pk1, Description{})
	js, err := codec.Cdc.MarshalJSON(validator)
	require.NoError(t, err)
	require.NotEmpty(t, js)
	require.Contains(t, string(js), "\"consensus_pubkey\":\"fcp")
	got := &Validator{}
	err = codec.Cdc.UnmarshalJSON(js, got)
	assert.NoError(t, err)
	assert.Equal(t, validator, *got)
}

func TestValidatorSetInitialCommission(t *testing.T) {
	val := NewValidator(addr1, pk1, Description{})
	testCases := []struct {
		validator   Validator
		commission  Commission
		expectedErr bool
	}{
		{val, NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()), false},
		{val, NewCommission(sdk.ZeroDec(), sdk.NewDecWithPrec(-1, 1), sdk.ZeroDec()), true},
		{val, NewCommission(sdk.ZeroDec(), sdk.NewDec(15000000000), sdk.ZeroDec()), true},
		{val, NewCommission(sdk.NewDecWithPrec(-1, 1), sdk.ZeroDec(), sdk.ZeroDec()), true},
		{val, NewCommission(sdk.NewDecWithPrec(2, 1), sdk.NewDecWithPrec(1, 1), sdk.ZeroDec()), true},
		{val, NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.NewDecWithPrec(-1, 1)), true},
		{val, NewCommission(sdk.ZeroDec(), sdk.NewDecWithPrec(1, 1), sdk.NewDecWithPrec(2, 1)), true},
	}

	for i, tc := range testCases {
		val, err := tc.validator.SetInitialCommission(tc.commission)

		if tc.expectedErr {
			require.Error(t, err,
				"expected error for test case #%d with commission: %s", i, tc.commission,
			)
		} else {
			require.NoError(t, err,
				"unexpected error for test case #%d with commission: %s", i, tc.commission,
			)
			require.Equal(t, tc.commission, val.Commission,
				"invalid validator commission for test case #%d with commission: %s", i, tc.commission,
			)
		}
	}
}
