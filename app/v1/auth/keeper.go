package auth

import (
	"fmt"
	"strings"

	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"github.com/tendermint/tendermint/crypto"
)

var (
	// Prefix for account-by-address store
	addressStoreKeyPrefix  = []byte("account:")
	globalAccountNumberKey = []byte("globalAccountNumber")
	TotalLoosenTokenKey    = []byte("totalLoosenToken")
	//BurnedTokenKey = []byte("burnedToken")
	totalSupplyKeyPrefix = []byte("totalSupply:")

	//system default special address
	BurnedCoinsAccAddr       = sdk.AccAddress(crypto.AddressHash([]byte("burnedCoins")))
	GovDepositCoinsAccAddr   = sdk.AccAddress(crypto.AddressHash([]byte("govDepositedCoins")))
	CommunityTaxCoinsAccAddr = sdk.AccAddress(crypto.AddressHash([]byte("communityTaxCoins")))

	ServiceDepositCoinsAccAddr = sdk.AccAddress(crypto.AddressHash([]byte("serviceDepositedCoins")))
	ServiceRequestCoinsAccAddr = sdk.AccAddress(crypto.AddressHash([]byte("serviceRequestCoins")))
	ServiceTaxCoinsAccAddr     = sdk.AccAddress(crypto.AddressHash([]byte("serviceTaxCoins")))

	HTLCLockedCoinsAccAddr = sdk.AccAddress(crypto.AddressHash([]byte("HTLCLockedCoins"))) // HTLCLockedCoinsAccAddr store All HTLC locked coins
)

// This AccountKeeper encodes/decodes accounts using the
// go-amino (binary) encoding/decoding library.
type AccountKeeper struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey

	// The prototypical Account constructor.
	proto func() Account

	// The codec codec for binary encoding/decoding of accounts.
	cdc *codec.Codec
}

// NewAccountKeeper returns a new sdk.AccountKeeper that
// uses go-amino to (binary) encode and decode concrete sdk.Accounts.
// nolint
func NewAccountKeeper(cdc *codec.Codec, key sdk.StoreKey, proto func() Account) AccountKeeper {
	return AccountKeeper{
		key:   key,
		proto: proto,
		cdc:   cdc,
	}
}

// Implaements sdk.AccountKeeper.
func (am AccountKeeper) NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) Account {
	acc := am.proto()
	err := acc.SetAddress(addr)
	if err != nil {
		// Handle w/ #870
		panic(err)
	}
	err = acc.SetAccountNumber(am.GetNextAccountNumber(ctx))
	if err != nil {
		// Handle w/ #870
		panic(err)
	}
	return acc
}

// Turn an address to key used to get it from the account store
func AddressStoreKey(addr sdk.AccAddress) []byte {
	return append(addressStoreKeyPrefix, addr.Bytes()...)
}

// Implements sdk.AccountKeeper.
func (am AccountKeeper) GetAccount(ctx sdk.Context, addr sdk.AccAddress) Account {
	store := ctx.KVStore(am.key)
	bz := store.Get(AddressStoreKey(addr))
	if bz == nil {
		return nil
	}
	acc := am.decodeAccount(bz)
	return acc
}

// Implements sdk.AccountKeeper.
func (am AccountKeeper) SetGenesisAccount(ctx sdk.Context, acc Account) {
	if !acc.GetAddress().Equals(BurnedCoinsAccAddr) {
		am.IncreaseTotalLoosenToken(ctx, acc.GetCoins())
	}
	am.SetAccount(ctx, acc)
}

// Implements sdk.AccountKeeper.
func (am AccountKeeper) SetAccount(ctx sdk.Context, acc Account) {
	addr := acc.GetAddress()
	store := ctx.KVStore(am.key)
	bz := am.encodeAccount(acc)
	store.Set(AddressStoreKey(addr), bz)
}

// RemoveAccount removes an account for the account mapper store.
func (am AccountKeeper) RemoveAccount(ctx sdk.Context, acc Account) {
	addr := acc.GetAddress()
	store := ctx.KVStore(am.key)
	store.Delete(AddressStoreKey(addr))
}

// Implements sdk.AccountKeeper.
func (am AccountKeeper) IterateAccounts(ctx sdk.Context, process func(Account) (stop bool)) {
	store := ctx.KVStore(am.key)
	iter := sdk.KVStorePrefixIterator(store, addressStoreKeyPrefix)
	defer iter.Close()
	for {
		if !iter.Valid() {
			return
		}
		val := iter.Value()
		acc := am.decodeAccount(val)
		if process(acc) {
			return
		}
		iter.Next()
	}
}

func (am AccountKeeper) GetAllAccounts(ctx sdk.Context) []Account {
	accounts := []Account{}
	appendAccount := func(acc Account) (stop bool) {
		accounts = append(accounts, acc)
		return false
	}
	am.IterateAccounts(ctx, appendAccount)
	return accounts
}

// Returns the PubKey of the account at address
func (am AccountKeeper) GetPubKey(ctx sdk.Context, addr sdk.AccAddress) (crypto.PubKey, sdk.Error) {
	acc := am.GetAccount(ctx, addr)
	if acc == nil {
		return nil, sdk.ErrUnknownAddress(addr.String())
	}
	return acc.GetPubKey(), nil
}

// Returns the Sequence of the account at address
func (am AccountKeeper) GetSequence(ctx sdk.Context, addr sdk.AccAddress) (uint64, sdk.Error) {
	acc := am.GetAccount(ctx, addr)
	if acc == nil {
		return 0, sdk.ErrUnknownAddress(addr.String())
	}
	return acc.GetSequence(), nil
}

func (am AccountKeeper) setSequence(ctx sdk.Context, addr sdk.AccAddress, newSequence uint64) sdk.Error {
	acc := am.GetAccount(ctx, addr)
	if acc == nil {
		return sdk.ErrUnknownAddress(addr.String())
	}
	err := acc.SetSequence(newSequence)
	if err != nil {
		// Handle w/ #870
		panic(err)
	}
	am.SetAccount(ctx, acc)
	return nil
}

// Returns and increments the global account number counter
func (am AccountKeeper) GetNextAccountNumber(ctx sdk.Context) uint64 {
	var accNumber uint64
	store := ctx.KVStore(am.key)
	bz := store.Get(globalAccountNumberKey)
	if bz == nil {
		accNumber = 0
	} else {
		am.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &accNumber)
	}

	bz = am.cdc.MustMarshalBinaryLengthPrefixed(accNumber + 1)
	store.Set(globalAccountNumberKey, bz)

	return accNumber
}

func (am AccountKeeper) GetTotalLoosenToken(ctx sdk.Context) sdk.Coins {
	// read from db
	var totalLoosenToken sdk.Coins
	store := ctx.KVStore(am.key)
	bz := store.Get(TotalLoosenTokenKey)
	if bz == nil {
		totalLoosenToken = nil
	} else {
		am.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &totalLoosenToken)
	}
	return totalLoosenToken
}

func (am AccountKeeper) IncreaseTotalLoosenToken(ctx sdk.Context, coins sdk.Coins) {
	if coins == nil || coins.Empty() {
		return
	}

	// loose token only contains iris-atto
	deltaCoin, err := coins.GetCoin(sdk.IrisAtto)
	if err != nil {
		panic(fmt.Sprintf("invalid coins [%s]", coins))
	}

	if !deltaCoin.IsPositive() {
		return
	}

	increaseCoins := sdk.Coins{deltaCoin}

	// read from db
	totalLoosenToken := am.GetTotalLoosenToken(ctx)
	// increase totalLoosenToken
	totalLoosenToken = totalLoosenToken.Add(increaseCoins)
	// write back to db
	bzNew := am.cdc.MustMarshalBinaryLengthPrefixed(totalLoosenToken)
	store := ctx.KVStore(am.key)
	store.Set(TotalLoosenTokenKey, bzNew)

	ctx.Logger().Info("Execute IncreaseTotalLoosenToken Successed",
		"increaseCoins", increaseCoins.String(), "totalLoosenToken", totalLoosenToken.String())
}

func (am AccountKeeper) DecreaseTotalLoosenToken(ctx sdk.Context, coins sdk.Coins) {
	if coins == nil || coins.Empty() {
		return
	}

	// loose token only contains iris-atto
	deltaCoin, err := coins.GetCoin(sdk.IrisAtto)
	if err != nil {
		panic(fmt.Sprintf("invalid coins [%s]", coins))
	}

	if !deltaCoin.IsPositive() {
		return
	}

	decreaseCoins := sdk.Coins{deltaCoin}

	// read from db
	totalLoosenToken := am.GetTotalLoosenToken(ctx)
	// decrease totalLoosenToken
	totalLoosenToken, hasNeg := totalLoosenToken.SafeSub(decreaseCoins)
	if hasNeg {
		panic(fmt.Errorf("total loose token is negative"))
	}
	// write back to db
	bzNew := am.cdc.MustMarshalBinaryLengthPrefixed(totalLoosenToken)
	store := ctx.KVStore(am.key)
	store.Set(TotalLoosenTokenKey, bzNew)

	ctx.Logger().Info("Execute DecreaseTotalLoosenToken Successed",
		"decreaseCoins", decreaseCoins.String(), "totalLoosenToken", totalLoosenToken.String())
}

// Turn a token id to key used to get it from the account store
func TotalSupplyStoreKey(denom string) []byte {
	keyId, _ := sdk.ConvertDenomToTokenKeyId(denom)
	return append(totalSupplyKeyPrefix, keyId...)
}

func (am AccountKeeper) IncreaseTotalSupply(ctx sdk.Context, coin sdk.Coin) sdk.Error {
	if !strings.HasSuffix(coin.Denom, sdk.MinDenomSuffix) {
		return sdk.ErrInvalidCoins(fmt.Sprintf("invalid coin [%s]", coin))
	}

	if !coin.IsPositive() {
		return nil
	}

	// read from db
	totalSupply, found := am.GetTotalSupply(ctx, coin.Denom)
	if !found {
		return sdk.ErrInvalidCoins(fmt.Sprintf("unable to get total supply for denom %s", coin.Denom))
	}

	// increase totalSupply
	totalSupply = totalSupply.Add(coin)

	// write back to db
	am.SetTotalSupply(ctx, totalSupply)

	ctx.Logger().Info("Execute IncreaseTotalSupply Succeeded",
		"increaseCoins", coin.String(), "totalSupply", totalSupply.String())

	return nil
}

func (am AccountKeeper) DecreaseTotalSupply(ctx sdk.Context, coin sdk.Coin) sdk.Error {
	if !strings.HasSuffix(coin.Denom, sdk.MinDenomSuffix) {
		return sdk.ErrInvalidCoins(fmt.Sprintf("invalid coin [%s]", coin))
	}

	if !coin.IsPositive() {
		return nil
	}

	// read from db
	totalSupply, found := am.GetTotalSupply(ctx, coin.Denom)
	if !found {
		return sdk.ErrInvalidCoins(fmt.Sprintf("unable to get total supply for denom %s", coin.Denom))
	}

	// decrease totalSupply
	totalSupply = totalSupply.Sub(coin)
	if totalSupply.IsNegative() {
		panic(fmt.Errorf("total supply is negative"))
	}

	// write back to db
	am.SetTotalSupply(ctx, totalSupply)

	ctx.Logger().Info("Execute DecreaseTotalSupply Succeeded",
		"decreaseCoins", coin.String(), "totalSupply", totalSupply.String())

	return nil
}

func (am AccountKeeper) GetTotalSupply(ctx sdk.Context, denom string) (coin sdk.Coin, found bool) {
	store := ctx.KVStore(am.key)
	bz := store.Get(TotalSupplyStoreKey(denom))
	if bz == nil {
		return coin, false
	}

	am.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &coin)
	return coin, true
}

func (am AccountKeeper) GetTotalSupplies(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(am.key)
	return sdk.KVStorePrefixIterator(store, TotalSupplyStoreKey(""))
}

func (am AccountKeeper) SetTotalSupply(ctx sdk.Context, totalSupply sdk.Coin) {
	// write back to db
	bzNew := am.cdc.MustMarshalBinaryLengthPrefixed(totalSupply)
	store := ctx.KVStore(am.key)
	store.Set(TotalSupplyStoreKey(totalSupply.Denom), bzNew)
}

func (am AccountKeeper) InitTotalSupply(ctx sdk.Context) {
	tsMap := make(map[string]sdk.Coin)
	am.IterateAccounts(ctx, func(account Account) (stop bool) {
		for _, coin := range account.GetCoins() {
			if sdk.IrisAtto == coin.Denom || sdk.Iris == coin.Denom || strings.HasPrefix(coin.Denom, sdk.FormatUniABSPrefix) {
				continue
			}
			totalSupply, ok := tsMap[coin.Denom]
			if !ok {
				tsMap[coin.Denom] = coin
			} else {
				tsMap[coin.Denom] = coin.Add(totalSupply)
			}
		}
		return false
	})

	// Defense against empty token's keyId
	var coins sdk.Coins
	for _, coin := range tsMap {
		coins = append(coins, coin)
	}

	coins = coins.Sort()
	for _, coin := range coins {
		am.SetTotalSupply(ctx, coin)
	}
}

//----------------------------------------
// misc.

func (am AccountKeeper) encodeAccount(acc Account) []byte {
	bz := am.cdc.MustMarshalBinaryBare(acc)
	return bz
}

func (am AccountKeeper) decodeAccount(bz []byte) (acc Account) {
	am.cdc.MustUnmarshalBinaryBare(bz, &acc)
	return
}
