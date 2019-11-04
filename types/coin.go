package types

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

//-----------------------------------------------------------------------------
// Coin

// Coin hold some amount of one currency.
//
// CONTRACT: A coin will never hold a negative amount of any denomination.
//
// TODO: Make field members private for further safety.
type Coin struct {
	Denom string `json:"denom"`

	// To allow the use of unsigned integers (see: #1273) a larger refactor will
	// need to be made. So we use signed integers for now with safety measures in
	// place preventing negative values being used.
	Amount Int `json:"amount"`
}

// NewCoin returns a new coin with a denomination and amount. It will panic if
// the amount is negative.
func NewCoin(denom string, amount Int) Coin {
	if amount.i == nil {
		amount = ZeroInt()
	}

	if amount.IsNegative() {
		panic("negative coin amount")
	}

	return Coin{
		Denom:  denom,
		Amount: amount,
	}
}

// NewInt64Coin returns a new coin with a denomination and amount. It will panic
// if the amount is negative.
func NewInt64Coin(denom string, amount int64) Coin {
	return NewCoin(denom, NewInt(amount))
}

// String provides a human-readable representation of a coin
func (coin Coin) String() string {
	return fmt.Sprintf("%v%v", coin.Amount, coin.Denom)
}

// IsValid returns true if the coin amount is non-negative
// and the coin is denominated in its minimum unit
func (coin Coin) IsValid() bool {
	if coin.IsNegative() {
		return false
	}
	return IsCoinMinDenomValid(coin.Denom)
}

func (coin Coin) IsValidIrisAtto() bool {
	return coin.Denom == IrisAtto && coin.IsPositive()
}

// IsZero returns if this coin has zero amount
func (coin Coin) IsZero() bool {
	return coin.Amount.i == nil || coin.Amount.IsZero()
}

// IsGTE returns true if they are the same type and the receiver is
// an equal or greater value
func (coin Coin) IsGTE(other Coin) bool {
	return coin.Denom == other.Denom && coin.Amount.GTE(other.Amount)
}

// IsLT returns true if they are the same type and the receiver is
// a smaller value
func (coin Coin) IsLT(other Coin) bool {
	return coin.Denom == other.Denom && coin.Amount.LT(other.Amount)
}

// IsEqual returns true if the two sets of Coins have the same value
func (coin Coin) IsEqual(other Coin) bool {
	return coin.Denom == other.Denom && coin.Amount.Equal(other.Amount)
}

// Adds amounts of two coins with same denom. If the coins differ in denom then
// it panics.
func (coin Coin) Add(coinB Coin) Coin {
	if coin.Denom != coinB.Denom {
		panic(fmt.Sprintf("invalid coin denominations; %s, %s", coin.Denom, coinB.Denom))
	}

	return Coin{coin.Denom, coin.Amount.Add(coinB.Amount)}
}

// Subtracts amounts of two coins with same denom. If the coins differ in denom
// then it panics.
func (coin Coin) Sub(coinB Coin) Coin {
	if coin.Denom != coinB.Denom {
		panic(fmt.Sprintf("invalid coin denominations; %s, %s", coin.Denom, coinB.Denom))
	}

	return NewCoin(coin.Denom, coin.Amount.Sub(coinB.Amount))
}

// IsPositive returns true if coin amount is positive.
//
// TODO: Remove once unsigned integers are used.
func (coin Coin) IsPositive() bool {
	return coin.Amount.i != nil && coin.Amount.IsPositive()
}

// IsNegative returns true if the coin amount is negative and false otherwise.
//
// TODO: Remove once unsigned integers are used.
func (coin Coin) IsNegative() bool {
	return coin.Amount.i != nil && coin.Amount.IsNegative()
}

//-----------------------------------------------------------------------------
// Coins

// Coins is a set of Coin, one per currency
type Coins []Coin

// NewCoins constructs a new coin set.
func NewCoins(coins ...Coin) Coins {
	// remove zeroes
	newCoins := removeZeroCoins(Coins(coins))
	if len(newCoins) == 0 {
		return Coins{}
	}

	newCoins.Sort()

	if !newCoins.IsValid() {
		panic(fmt.Sprintf("invalid coin set: %s", newCoins))
	}

	return newCoins
}

func (coins Coins) String() string {
	if len(coins) == 0 {
		return ""
	}

	out := ""
	for _, coin := range coins {
		out += fmt.Sprintf("%v,", coin.String())
	}
	return out[:len(out)-1]
}

// MainUnitString() returns a string representation of coins,
// with iris-atto coin converted to its corresponding iris denomination
func (coins Coins) MainUnitString() string {
	if len(coins) == 0 {
		return ""
	}
	out := ""
	for _, coin := range coins {
		// only convert iris now
		if coin.Denom == IrisAtto {
			destCoinStr, err := IrisCoinType.Convert(coin.String(), Iris)
			if err == nil {
				out += fmt.Sprintf("%v,", destCoinStr)
				continue
			}
		}
		out += fmt.Sprintf("%v,", coin.String())
	}
	if len(out) > 0 {
		out = out[:len(out)-1]
	}
	return out
}

// IsValid asserts the coins are valid and sorted.
func (coins Coins) IsValid() bool {
	switch len(coins) {
	case 0:
		return true
	case 1:
		return coins[0].IsValid()
	default:
		// check first coin
		if !coins[0].IsValid() {
			return false
		}

		lowDenom := coins[0].Denom
		for _, coin := range coins[1:] {
			if !coin.IsValid() {
				return false
			}
			if coin.Denom <= lowDenom {
				return false
			}

			// we compare each coin against the last denom
			lowDenom = coin.Denom
		}

		return true
	}
}

func (coins Coins) IsValidIrisAtto() bool {
	if coins == nil || len(coins) != 1 {
		return false
	}
	return coins[0].IsValidIrisAtto()
}

// Add adds two sets of coins.
//
// e.g.
// {2A} + {A, 2B} = {3A, 2B}
// {2A} + {0B} = {2A}
//
// NOTE: Add operates under the invariant that coins are sorted by
// denominations.
//
// CONTRACT: Add will never return Coins where one Coin has a negative
// amount. In other words, IsValid will always return true.
func (coins Coins) Add(coinsB Coins) Coins {
	sum, hasNeg := coins.SafeAdd(coinsB)
	if hasNeg {
		panic("negative coin amount")
	}

	return sum
}

// SafeAdd performs the same arithmetic as Add but returns a boolean if any
// negative coin amount was returned.
func (coins Coins) SafeAdd(coinsB Coins) (Coins, bool) {
	sum := coins.safeAdd(coinsB)
	return sum, sum.IsAnyNegative()
}

// safeAdd will perform addition of two coins sets. If both coin sets are
// empty, then an empty set is returned. If only a single set is empty, the
// other set is returned. Otherwise, the coins are compared in order of their
// denomination and addition only occurs when the denominations match, otherwise
// the coin is simply added to the sum assuming it's not zero.
func (coins Coins) safeAdd(coinsB Coins) Coins {
	sum := ([]Coin)(nil)
	indexA, indexB := 0, 0
	lenA, lenB := len(coins), len(coinsB)

	for {
		if indexA == lenA {
			if indexB == lenB {
				// return nil coins if both sets are empty
				return sum
			}

			// return set B (excluding zero coins) if set A is empty
			return append(sum, removeZeroCoins(coinsB[indexB:])...)
		} else if indexB == lenB {
			// return set A (excluding zero coins) if set B is empty
			return append(sum, removeZeroCoins(coins[indexA:])...)
		}

		coinA, coinB := coins[indexA], coinsB[indexB]

		switch strings.Compare(coinA.Denom, coinB.Denom) {
		case -1: // coin A denom < coin B denom
			if !coinA.IsZero() {
				sum = append(sum, coinA)
			}

			indexA++

		case 0: // coin A denom == coin B denom
			res := coinA.Add(coinB)
			if !res.IsZero() {
				sum = append(sum, res)
			}

			indexA++
			indexB++

		case 1: // coin A denom > coin B denom
			if !coinB.IsZero() {
				sum = append(sum, coinB)
			}

			indexB++
		}
	}
}

// Sub subtracts a set of coins from another.
//
// e.g.
// {2A, 3B} - {A} = {A, 3B}
// {2A} - {0B} = {2A}
// {A, B} - {A} = {B}
//
// CONTRACT: Sub will never return Coins where one Coin has a negative
// amount. In other words, IsValid will always return true.
func (coins Coins) Sub(coinsB Coins) Coins {
	diff, hasNeg := coins.SafeSub(coinsB)
	if hasNeg {
		panic("negative coin amount")
	}

	return diff
}

// SafeSub performs the same arithmetic as Sub but returns a boolean if any
// negative coin amount was returned.
func (coins Coins) SafeSub(coinsB Coins) (Coins, bool) {
	diff := coins.safeAdd(coinsB.negative())
	return diff, diff.IsAnyNegative()
}

// IsAllGT returns true if for every denom in coinsB,
// the denom is present at a greater amount in coins.
func (coins Coins) IsAllGT(coinsB Coins) bool {
	if coins.IsZero() {
		return false
	}

	if coinsB.IsZero() {
		return true
	}

	for _, coinB := range coinsB {
		if coins.AmountOf(coinB.Denom).LTE(coinB.Amount) {
			return false
		}
	}

	return true
}

// IsAllGTE returns true if for every denom in coinsB,
// the denom is present at a greater or equal amount in coins.
func (coins Coins) IsAllGTE(coinsB Coins) bool {
	if coinsB.IsZero() {
		return true
	}

	if coins.IsZero() {
		return false
	}

	for _, coinB := range coinsB {
		if coins.AmountOf(coinB.Denom).LT(coinB.Amount) {
			return false
		}
	}

	return true
}

// IsAllLT returns true if for every denom in coins, the denom is present at
// a greater amount in coinsB.
func (coins Coins) IsAllLT(coinsB Coins) bool {
	return coinsB.IsAllGT(coins)
}

// IsAllLTE returns true if for every denom in coins, the denom is present at
// a smaller or equal amount in coinsB.
func (coins Coins) IsAllLTE(coinsB Coins) bool {
	return coinsB.IsAllGTE(coins)
}

// IsAnyGT returns true if there exists at least one denom in coins
// that is present in coinsB with a smaller amount.
//
// e.g.
// {2A, 3B}.IsAnyGT{A} = true
// {2A, 3B}.IsAnyGT{5C} = false
// {}.IsAnyGT{5C} = false
// {2A, 3B}.IsAnyGT{} = false
func (coins Coins) IsAnyGT(coinsB Coins) bool {
	if len(coinsB) == 0 {
		return false
	}

	for _, coin := range coins {
		amt := coinsB.AmountOf(coin.Denom)
		if coin.Amount.GT(amt) && !amt.IsZero() {
			return true
		}
	}

	return false
}

// IsAnyGT returns true if there exists at least one denom in coins
// that is present in coinsB with a smaller or equal amount.
//
func (coins Coins) IsAnyGTE(coinsB Coins) bool {
	if len(coinsB) == 0 {
		return false
	}

	for _, coin := range coins {
		amt := coinsB.AmountOf(coin.Denom)
		if coin.Amount.GTE(amt) && !amt.IsZero() {
			return true
		}
	}

	return false
}

// IsZero returns true if there are no coins or all coins are zero.
func (coins Coins) IsZero() bool {
	for _, coin := range coins {
		if !coin.IsZero() {
			return false
		}
	}
	return true
}

// IsEqual returns true if the two sets of Coins have the same value
func (coins Coins) IsEqual(coinsB Coins) bool {
	if len(coins) != len(coinsB) {
		return false
	}

	coins = coins.Sort()
	coinsB = coinsB.Sort()

	for i := 0; i < len(coins); i++ {
		if !coins[i].IsEqual(coinsB[i]) {
			return false
		}
	}

	return true
}

// Empty returns true if there are no coins and false otherwise.
func (coins Coins) Empty() bool {
	return len(coins) == 0
}

// Returns the amount of a denom from coins
func (coins Coins) AmountOf(denom string) Int {
	switch len(coins) {
	case 0:
		return ZeroInt()

	case 1:
		coin := coins[0]
		if coin.Denom == denom {
			return coin.Amount
		}
		return ZeroInt()

	default:
		midIdx := len(coins) / 2 // 2:1, 3:1, 4:2
		coin := coins[midIdx]

		if denom < coin.Denom {
			return coins[:midIdx].AmountOf(denom)
		} else if denom == coin.Denom {
			return coin.Amount
		} else {
			return coins[midIdx+1:].AmountOf(denom)
		}
	}
}

// IsAllPositive returns true if all coins have positive values.
//
// TODO: Remove once unsigned integers are used.
func (coins Coins) IsAllPositive() bool {
	if len(coins) == 0 {
		return false
	}

	for _, coin := range coins {
		if !coin.IsPositive() {
			return false
		}
	}

	return true
}

// IsAnyNegative returns true if at least one coin has negative amount.
//
// TODO: Remove once unsigned integers are used.
func (coins Coins) IsAnyNegative() bool {
	for _, coin := range coins {
		if coin.IsNegative() {
			return true
		}
	}

	return false
}

// negative returns a set of coins with all amount negative.
//
// TODO: Remove once unsigned integers are used.
func (coins Coins) negative() Coins {
	res := make([]Coin, 0, len(coins))

	for _, coin := range coins {
		res = append(res, Coin{coin.Denom, coin.Amount.Neg()})
	}

	return res
}

func (coins Coins) GetCoin(denom string) (Coin, error) {
	for _, coin := range coins {
		if coin.Denom == denom {
			return coin, nil
		}
	}
	return Coin{}, fmt.Errorf("cannot find coin with denom %s", denom)
}

// removeZeroCoins removes all zero coins from the given coin set in-place.
func removeZeroCoins(coins Coins) Coins {
	i, l := 0, len(coins)
	for i < l {
		if coins[i].IsZero() {
			// remove coin
			coins = append(coins[:i], coins[i+1:]...)
			l--
		} else {
			i++
		}
	}

	return coins[:i]
}

//-----------------------------------------------------------------------------
// Sort interface

//nolint
func (coins Coins) Len() int           { return len(coins) }
func (coins Coins) Less(i, j int) bool { return coins[i].Denom < coins[j].Denom }
func (coins Coins) Swap(i, j int)      { coins[i], coins[j] = coins[j], coins[i] }

var _ sort.Interface = Coins{}

// Sort is a helper function to sort the set of coins inplace
func (coins Coins) Sort() Coins {
	sort.Sort(coins)
	return coins
}

//-----------------------------------------------------------------------------
// Parsing & Checking

var (
	// Denominations can be 3 ~ 21 characters long.
	reABS              = `([a-z][0-9a-z]{2}[:])?`
	reCoinName         = reABS + `(([a-z][a-z0-9]{2,7}|x)\.)?([a-z][a-z0-9]{2,7})`
	reDenom            = reCoinName + `(-[a-z]{3,5})?`
	reAmount           = `[0-9]+(\.[0-9]+)?`
	reSpace            = `[[:space:]]*`
	reCoinNameCompiled = regexp.MustCompile(fmt.Sprintf(`^%s$`, reCoinName))
	reDenomCompiled    = regexp.MustCompile(fmt.Sprintf(`^%s$`, reDenom))
	reCoinCompiled     = regexp.MustCompile(fmt.Sprintf(`^(%s)%s(%s)$`, reAmount, reSpace, reDenom))
)

func ParseCoinParts(coinStr string) (denom, amount string, err error) {
	coinStr = strings.ToLower(strings.TrimSpace(coinStr))

	matches := reCoinCompiled.FindStringSubmatch(coinStr)
	if matches == nil {
		err = fmt.Errorf("invalid coin string: %s", coinStr)
		return
	}
	denom, amount = matches[3], matches[1]
	return
}

// ParseCoin parses a cli input for one coin type, returning errors if invalid.
// This returns an error on an empty string as well.
func ParseCoin(coinStr string) (coin Coin, err error) {
	coinStr = strings.ToLower(strings.TrimSpace(coinStr))

	matches := reCoinCompiled.FindStringSubmatch(coinStr)
	if matches == nil {
		return Coin{}, fmt.Errorf("invalid coin expression: %s", coinStr)
	}

	denomStr, amountStr := matches[3], matches[1]

	amount, ok := NewIntFromString(amountStr)
	if !ok {
		return Coin{}, fmt.Errorf("failed to parse coin amount: %s", amountStr)
	}

	return NewCoin(denomStr, amount), nil
}

// ParseCoins will parse out a list of coins separated by commas.
// If nothing is provided, it returns nil Coins.
// Returned coins are sorted.
func ParseCoins(coinsStr string) (coins Coins, err error) {
	if len(coinsStr) == 0 {
		return Coins{}, nil
	}

	coinStrs := strings.Split(coinsStr, ",")
	for _, coinStr := range coinStrs {
		coin, err := ParseCoin(coinStr)
		if err != nil {
			return nil, err
		}
		coins = append(coins, coin)
	}

	// Sort coins for determinism.
	coins.Sort()

	return coins, nil
}

func IsCoinNameValid(coinName string) bool {
	return reCoinNameCompiled.MatchString(coinName)
}

func IsCoinMinDenomValid(denom string) bool {
	if denom != IrisAtto && (!strings.HasSuffix(denom, MinDenomSuffix) || strings.HasPrefix(denom, Iris+"-")) {
		return false
	}
	return reDenomCompiled.MatchString(denom)
}

func (coins Coins) IsValidV0() bool {
	switch len(coins) {
	case 0:
		return true
	case 1:
		return coins[0].IsPositive()
	default:
		if !coins[0].IsPositive() {
			return false
		}

		lowDenom := coins[0].Denom

		for _, coin := range coins[1:] {
			if coin.Denom <= lowDenom {
				return false
			}
			if !coin.IsPositive() {
				return false
			}

			// we compare each coin against the last denom
			lowDenom = coin.Denom
		}

		return true
	}
}
