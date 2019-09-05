package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
// Coin tests

func TestCoin(t *testing.T) {
	require.Panics(t, func() { NewInt64Coin("atom-min", -1) })
	require.Panics(t, func() { NewCoin("atom-min", NewInt(-1)) })
	require.Equal(t, NewInt(5), NewInt64Coin("atom-min", 5).Amount)
	require.Equal(t, NewInt(5), NewCoin("atom-min", NewInt(5)).Amount)
}

func TestNewCoins(t *testing.T) {
	tenatom := NewInt64Coin("atom-min", 10)
	tenbtc := NewInt64Coin("btc-min", 10)
	zeroeth := NewInt64Coin("eth-min", 0)

	println(NewCoins(zeroeth, tenatom, tenbtc).String())
	tests := []struct {
		name      string
		coins     Coins
		want      Coins
		wantPanic bool
	}{
		{"empty args", []Coin{}, Coins{}, false},
		{"one coin", []Coin{tenatom}, Coins{tenatom}, false},
		{"sort after create", []Coin{tenbtc, tenatom}, Coins{tenatom, tenbtc}, false},
		{"sort and remove zeroes", []Coin{zeroeth, tenbtc, tenatom}, Coins{tenatom, tenbtc}, false},
		{"panic on dups", []Coin{tenatom, tenatom}, Coins{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				require.Panics(t, func() { NewCoins(tt.coins...) })
				return
			}
			got := NewCoins(tt.coins...)
			require.True(t, got.IsEqual(tt.want))
		})
	}
}

func TestSameDenomAsCoin(t *testing.T) {
	cases := []struct {
		inputOne Coin
		inputTwo Coin
		expected bool
	}{
		{NewInt64Coin("atom-min", 1), NewInt64Coin("atom-min", 1), true},
		{NewInt64Coin("atom-min", 1), NewInt64Coin("atom2-min", 1), false},
		{NewInt64Coin("atom-min", 1), NewInt64Coin("btc-min", 1), false},
		{NewInt64Coin("steak-min", 1), NewInt64Coin("steak-min", 10), true},
	}

	for tcIndex, tc := range cases {
		res := (tc.inputOne.Denom == tc.inputTwo.Denom)
		require.Equal(t, tc.expected, res, "coin denominations didn't match, tc #%d", tcIndex)
	}
}

func TestIsEqualCoin(t *testing.T) {
	cases := []struct {
		inputOne Coin
		inputTwo Coin
		expected bool
	}{
		{NewInt64Coin("atom-min", 1), NewInt64Coin("atom-min", 1), true},
		{NewInt64Coin("atom-min", 1), NewInt64Coin("atom2-min", 1), false},
		{NewInt64Coin("atom-min", 1), NewInt64Coin("btc-min", 1), false},
		{NewInt64Coin("steak-min", 1), NewInt64Coin("steak-min", 10), false},
	}

	for tcIndex, tc := range cases {
		res := tc.inputOne.IsEqual(tc.inputTwo)
		require.Equal(t, tc.expected, res, "coin equality relation is incorrect, tc #%d", tcIndex)
	}
}

func TestPlusCoin(t *testing.T) {
	cases := []struct {
		inputOne    Coin
		inputTwo    Coin
		expected    Coin
		shouldPanic bool
	}{
		{NewInt64Coin("atom-min", 1), NewInt64Coin("atom-min", 1), NewInt64Coin("atom-min", 2), false},
		{NewInt64Coin("atom-min", 1), NewInt64Coin("atom-min", 0), NewInt64Coin("atom-min", 1), false},
		{NewInt64Coin("atom-min", 1), NewInt64Coin("btc-min", 1), NewInt64Coin("atom-min", 1), true},
	}

	for tcIndex, tc := range cases {
		if tc.shouldPanic {
			require.Panics(t, func() { tc.inputOne.Add(tc.inputTwo) })
		} else {
			res := tc.inputOne.Add(tc.inputTwo)
			require.Equal(t, tc.expected, res, "sum of coins is incorrect, tc #%d", tcIndex)
		}
	}
}

func TestMinusCoin(t *testing.T) {
	cases := []struct {
		inputOne    Coin
		inputTwo    Coin
		expected    Coin
		shouldPanic bool
	}{
		{NewInt64Coin("atom-min", 1), NewInt64Coin("btc-min", 1), NewInt64Coin("atom-min", 1), true},
		{NewInt64Coin("atom-min", 10), NewInt64Coin("atom-min", 1), NewInt64Coin("atom-min", 9), false},
		{NewInt64Coin("atom-min", 5), NewInt64Coin("atom-min", 3), NewInt64Coin("atom-min", 2), false},
		{NewInt64Coin("atom-min", 5), NewInt64Coin("atom-min", 0), NewInt64Coin("atom-min", 5), false},
		{NewInt64Coin("atom-min", 1), NewInt64Coin("atom-min", 5), Coin{}, true},
	}

	for tcIndex, tc := range cases {
		if tc.shouldPanic {
			require.Panics(t, func() { tc.inputOne.Sub(tc.inputTwo) })
		} else {
			res := tc.inputOne.Sub(tc.inputTwo)
			require.Equal(t, tc.expected, res, "difference of coins is incorrect, tc #%d", tcIndex)
		}
	}

	tc := struct {
		inputOne Coin
		inputTwo Coin
		expected int64
	}{NewInt64Coin("atom-min", 1), NewInt64Coin("atom-min", 1), 0}
	res := tc.inputOne.Sub(tc.inputTwo)
	require.Equal(t, tc.expected, res.Amount.Int64())
}

func TestIsGTECoin(t *testing.T) {
	cases := []struct {
		inputOne Coin
		inputTwo Coin
		expected bool
	}{
		{NewInt64Coin("atom-min", 1), NewInt64Coin("atom-min", 1), true},
		{NewInt64Coin("atom-min", 2), NewInt64Coin("atom-min", 1), true},
		{NewInt64Coin("atom-min", 1), NewInt64Coin("btc-min", 1), false},
	}

	for tcIndex, tc := range cases {
		res := tc.inputOne.IsGTE(tc.inputTwo)
		require.Equal(t, tc.expected, res, "coin GTE relation is incorrect, tc #%d", tcIndex)
	}
}

func TestIsLTCoin(t *testing.T) {
	cases := []struct {
		inputOne Coin
		inputTwo Coin
		expected bool
	}{
		{NewInt64Coin("atom-min", 1), NewInt64Coin("atom-min", 1), false},
		{NewInt64Coin("atom-min", 2), NewInt64Coin("atom-min", 1), false},
		{NewInt64Coin("atom-min", 0), NewInt64Coin("btc-min", 1), false},
		{NewInt64Coin("atom-min", 1), NewInt64Coin("btc-min", 1), false},
		{NewInt64Coin("atom-min", 1), NewInt64Coin("atom-min", 1), false},
		{NewInt64Coin("atom-min", 1), NewInt64Coin("atom-min", 2), true},
	}

	for tcIndex, tc := range cases {
		res := tc.inputOne.IsLT(tc.inputTwo)
		require.Equal(t, tc.expected, res, "coin LT relation is incorrect, tc #%d", tcIndex)
	}
}

func TestCoinIsZero(t *testing.T) {
	coin := NewInt64Coin("atom-min", 0)
	res := coin.IsZero()
	require.True(t, res)

	coin = NewInt64Coin("atom-min", 1)
	res = coin.IsZero()
	require.False(t, res)
}

// ----------------------------------------------------------------------------
// Coins tests

func TestIsZeroCoins(t *testing.T) {
	cases := []struct {
		inputOne Coins
		expected bool
	}{
		{Coins{}, true},
		{Coins{NewInt64Coin("atom-min", 0)}, true},
		{Coins{NewInt64Coin("atom-min", 0), NewInt64Coin("btc-min", 0)}, true},
		{Coins{NewInt64Coin("atom-min", 1)}, false},
		{Coins{NewInt64Coin("atom-min", 0), NewInt64Coin("btc-min", 1)}, false},
	}

	for _, tc := range cases {
		res := tc.inputOne.IsZero()
		require.Equal(t, tc.expected, res)
	}
}

func TestEqualCoins(t *testing.T) {
	cases := []struct {
		inputOne Coins
		inputTwo Coins
		expected bool
	}{
		{Coins{}, Coins{}, true},
		{Coins{NewInt64Coin("atom-min", 0)}, Coins{NewInt64Coin("atom-min", 0)}, true},
		{Coins{NewInt64Coin("atom-min", 0), NewInt64Coin("btc-min", 1)}, Coins{NewInt64Coin("atom-min", 0), NewInt64Coin("btc-min", 1)}, true},
		{Coins{NewInt64Coin("atom-min", 0)}, Coins{NewInt64Coin("btc-min", 0)}, false},
		{Coins{NewInt64Coin("atom-min", 0)}, Coins{NewInt64Coin("atom-min", 1)}, false},
		{Coins{NewInt64Coin("atom-min", 0)}, Coins{NewInt64Coin("atom-min", 0), NewInt64Coin("btc-min", 1)}, false},
		{Coins{NewInt64Coin("atom-min", 0), NewInt64Coin("btc-min", 1)}, Coins{NewInt64Coin("btc-min", 1), NewInt64Coin("atom-min", 0)}, true},
	}

	for tcnum, tc := range cases {
		res := tc.inputOne.IsEqual(tc.inputTwo)
		require.Equal(t, tc.expected, res, "Equality is differed from expected. tc #%d, expected %b, actual %b.", tcnum, tc.expected, res)
	}
}

func TestPlusCoins(t *testing.T) {
	zero := NewInt(0)
	one := NewInt(1)
	two := NewInt(2)

	cases := []struct {
		inputOne Coins
		inputTwo Coins
		expected Coins
	}{
		{Coins{{"atom-min", one}, {"btc-min", one}}, Coins{{"atom-min", one}, {"btc-min", one}}, Coins{{"atom-min", two}, {"btc-min", two}}},
		{Coins{{"atom-min", zero}, {"btc-min", one}}, Coins{{"atom-min", zero}, {"btc-min", zero}}, Coins{{"btc-min", one}}},
		{Coins{{"atom-min", two}}, Coins{{"btc-min", zero}}, Coins{{"atom-min", two}}},
		{Coins{{"atom-min", one}}, Coins{{"atom-min", one}, {"btc-min", two}}, Coins{{"atom-min", two}, {"btc-min", two}}},
		{Coins{{"atom-min", zero}, {"btc-min", zero}}, Coins{{"atom-min", zero}, {"btc-min", zero}}, Coins(nil)},
	}

	for tcIndex, tc := range cases {
		res := tc.inputOne.Add(tc.inputTwo)
		assert.True(t, res.IsValid())
		require.Equal(t, tc.expected, res, "sum of coins is incorrect, tc #%d", tcIndex)
	}
}

func TestMinusCoins(t *testing.T) {
	zero := NewInt(0)
	one := NewInt(1)
	two := NewInt(2)

	testCases := []struct {
		inputOne    Coins
		inputTwo    Coins
		expected    Coins
		shouldPanic bool
	}{
		{Coins{{"atom-min", two}}, Coins{{"atom-min", one}, {"btc-min", two}}, Coins{{"atom-min", one}, {"btc-min", two}}, true},
		{Coins{{"atom-min", two}}, Coins{{"btc-min", zero}}, Coins{{"atom-min", two}}, false},
		{Coins{{"atom-min", one}}, Coins{{"btc-min", zero}}, Coins{{"atom-min", one}}, false},
		{Coins{{"atom-min", one}, {"btc-min", one}}, Coins{{"atom-min", one}}, Coins{{"btc-min", one}}, false},
		{Coins{{"atom-min", one}, {"btc-min", one}}, Coins{{"atom-min", two}}, Coins{}, true},
	}

	for i, tc := range testCases {
		if tc.shouldPanic {
			require.Panics(t, func() { tc.inputOne.Sub(tc.inputTwo) })
		} else {
			res := tc.inputOne.Sub(tc.inputTwo)
			assert.True(t, res.IsValid())
			require.Equal(t, tc.expected, res, "sum of coins is incorrect, tc #%d", i)
		}
	}
}

func TestCoins(t *testing.T) {
	good := Coins{
		{"gas-min", NewInt(1)},
		{"mineral-min", NewInt(1)},
		{"tree-min", NewInt(1)},
	}
	empty := Coins{
		{"gold-min", NewInt(0)},
	}
	null := Coins{}
	badSort1 := Coins{
		{"tree-min", NewInt(1)},
		{"gas-min", NewInt(1)},
		{"mineral-min", NewInt(1)},
	}

	// both are after the first one, but the second and third are in the wrong order
	badSort2 := Coins{
		{"gas-min", NewInt(1)},
		{"tree-min", NewInt(1)},
		{"mineral-min", NewInt(1)},
	}
	badAmt := Coins{
		{"gas-min", NewInt(1)},
		{"mineral-min", NewInt(1)},
		{"tree-min", NewInt(0)},
	}
	dup := Coins{
		{"gas-min", NewInt(1)},
		{"gas-min", NewInt(1)},
		{"mineral-min", NewInt(1)},
	}
	neg := Coins{
		{"gas-min", NewInt(-1)},
		{"mineral-min", NewInt(1)},
	}

	assert.True(t, good.IsValid(), "Coins are valid")
	assert.True(t, good.IsAllPositive(), "Expected coins to be positive: %v", good)
	assert.False(t, null.IsAllPositive(), "Expected coins to not be positive: %v", null)
	assert.True(t, good.IsAllGTE(empty), "Expected %v to be >= %v", good, empty)
	assert.False(t, good.IsAllLT(empty), "Expected %v to be < %v", good, empty)
	assert.True(t, empty.IsAllLT(good), "Expected %v to be < %v", empty, good)
	assert.False(t, badSort1.IsValid(), "Coins are not sorted")
	assert.False(t, badSort2.IsValid(), "Coins are not sorted")
	assert.True(t, badAmt.IsValid(), "Coins can include 0 amounts")
	assert.False(t, dup.IsValid(), "Duplicate coin")
	assert.False(t, neg.IsValid(), "Negative first-denom coin")
}

func TestCoinsGT(t *testing.T) {
	one := NewInt(1)
	two := NewInt(2)

	assert.False(t, Coins{}.IsAllGT(Coins{}))
	assert.False(t, Coins{}.IsAllGT(Coins{Coin{"atom-min", ZeroInt()}}))
	assert.True(t, Coins{{"atom-min", one}}.IsAllGT(Coins{}))
	assert.False(t, Coins{{"atom-min", one}}.IsAllGT(Coins{{"atom-min", one}}))
	assert.False(t, Coins{{"atom-min", one}}.IsAllGT(Coins{{"btc-min", one}}))
	assert.True(t, Coins{{"atom-min", one}, {"btc-min", two}}.IsAllGT(Coins{{"btc-min", one}}))
	assert.False(t, Coins{{"atom-min", one}, {"btc-min", one}}.IsAllGT(Coins{{"btc-min", two}}))
}

func TestCoinsGTE(t *testing.T) {
	one := NewInt(1)
	two := NewInt(2)

	assert.True(t, Coins{}.IsAllGTE(Coins{}))
	assert.True(t, Coins{}.IsAllGTE(Coins{Coin{"atom-min", ZeroInt()}}))
	assert.True(t, Coins{{"atom-min", one}}.IsAllGTE(Coins{}))
	assert.True(t, Coins{{"atom-min", one}}.IsAllGTE(Coins{{"atom-min", one}}))
	assert.False(t, Coins{{"atom-min", one}}.IsAllGTE(Coins{{"btc-min", one}}))
	assert.True(t, Coins{{"atom-min", one}, {"btc-min", one}}.IsAllGTE(Coins{{"btc-min", one}}))
	assert.False(t, Coins{{"atom-min", one}, {"btc-min", one}}.IsAllGTE(Coins{{"btc-min", two}}))
}

func TestCoinsLT(t *testing.T) {
	one := NewInt(1)
	two := NewInt(2)

	assert.False(t, Coins{}.IsAllLT(Coins{}))
	assert.False(t, Coins{}.IsAllLT(Coins{Coin{"atom-min", ZeroInt()}}))
	assert.False(t, Coins{{"atom-min", one}}.IsAllLT(Coins{}))
	assert.False(t, Coins{{"atom-min", one}}.IsAllLT(Coins{{"atom-min", one}}))
	assert.False(t, Coins{{"atom-min", one}}.IsAllLT(Coins{{"btc-min", one}}))
	assert.False(t, Coins{{"atom-min", one}, {"btc-min", one}}.IsAllLT(Coins{{"btc-min", one}}))
	assert.False(t, Coins{{"atom-min", one}, {"btc-min", one}}.IsAllLT(Coins{{"btc-min", two}}))
	assert.False(t, Coins{{"atom-min", one}, {"btc-min", one}}.IsAllLT(Coins{{"atom-min", one}, {"btc-min", one}}))
	assert.True(t, Coins{{"atom-min", one}, {"btc-min", one}}.IsAllLT(Coins{{"atom-min", two}, {"btc-min", two}}))
	assert.True(t, Coins{}.IsAllLT(Coins{{"atom-min", one}}))
}

func TestCoinsLTE(t *testing.T) {
	one := NewInt(1)
	two := NewInt(2)

	assert.True(t, Coins{}.IsAllLTE(Coins{}))
	assert.True(t, Coins{}.IsAllLTE(Coins{Coin{"atom-min", ZeroInt()}}))
	assert.False(t, Coins{{"atom-min", one}}.IsAllLTE(Coins{}))
	assert.True(t, Coins{{"atom-min", one}}.IsAllLTE(Coins{{"atom-min", one}}))
	assert.False(t, Coins{{"atom-min", one}}.IsAllLTE(Coins{{"btc-min", one}}))
	assert.False(t, Coins{{"atom-min", one}, {"btc-min", one}}.IsAllLTE(Coins{{"btc-min", one}}))
	assert.False(t, Coins{{"atom-min", one}, {"btc-min", one}}.IsAllLTE(Coins{{"btc-min", two}}))
	assert.True(t, Coins{{"atom-min", one}, {"btc-min", one}}.IsAllLTE(Coins{{"atom-min", one}, {"btc-min", one}}))
	assert.True(t, Coins{{"atom-min", one}, {"btc-min", one}}.IsAllLTE(Coins{{"atom-min", one}, {"btc-min", two}}))
	assert.True(t, Coins{}.IsAllLTE(Coins{{"atom-min", one}}))
}

func TestParse(t *testing.T) {
	one := NewInt(1)

	cases := []struct {
		input    string
		valid    bool  // if false, we expect an error on parse
		expected Coins // if valid is true, make sure this is returned
	}{
		{"", true, Coins{}},
		{"1foo", true, Coins{{"foo", one}}},
		{"10bar", true, Coins{{"bar", NewInt(10)}}},
		{"99bar,1foo", true, Coins{{"bar", NewInt(99)}, {"foo", one}}},
		{"98 bar , 1 foo  ", true, Coins{{"bar", NewInt(98)}, {"foo", one}}},
		{"  55\t \t bling\n", true, Coins{{"bling", NewInt(55)}}},
		{"2foo, 97 bar", true, Coins{{"bar", NewInt(97)}, {"foo", NewInt(2)}}},
		{"5 mycoin,", false, nil},                                                  // no empty coins in a list
		{"2 3foo, 97 bar", false, Coins{{"3foo", NewInt(2)}, {"bar", NewInt(97)}}}, // 3foo is invalid coin name
		{"11me coin, 12you coin", false, nil},                                      // no spaces in coin names
		{"1.2btc", false, nil},                                                     // amount must be integer
		{"5foo-bar", true, Coins{{"foo-bar", NewInt(5)}}},
		{"5uni:gdex.eth-min", true, Coins{{"uni:gdex.eth-min", NewInt(5)}}},
		{"1un:x.btc", false, Coins{{"un:x.btc", NewInt(1)}}},
		{"1uni:x.btc", true, Coins{{"uni:x.btc", NewInt(1)}}},
		{"1 uni:x.btc", true, Coins{{"uni:x.btc", NewInt(1)}}},
		{"1uni:x.btc-min", true, Coins{{"uni:x.btc-min", NewInt(1)}}},
		{"1uni:iris-atto", true, Coins{{"uni:iris-atto", NewInt(1)}}},
	}

	for tcIndex, tc := range cases {
		res, err := ParseCoins(tc.input)
		if !tc.valid {
			require.NotNil(t, err, "%s: %#v. tc #%d", tc.input, res, tcIndex)
		} else if assert.Nil(t, err, "%s: %+v", tc.input, err) {
			require.Equal(t, tc.expected, res, "coin parsing was incorrect, tc #%d", tcIndex)
		}
	}
}

func TestSortCoins(t *testing.T) {
	good := Coins{
		NewInt64Coin("gas-min", 1),
		NewInt64Coin("mineral-min", 1),
		NewInt64Coin("tree-min", 1),
	}
	empty := Coins{
		NewInt64Coin("gold-min", 0),
	}
	badSort1 := Coins{
		NewInt64Coin("tree-min", 1),
		NewInt64Coin("gas-min", 1),
		NewInt64Coin("mineral-min", 1),
	}
	badSort2 := Coins{ // both are after the first one, but the second and third are in the wrong order
		NewInt64Coin("gas-min", 1),
		NewInt64Coin("tree-min", 1),
		NewInt64Coin("mineral-min", 1),
	}
	badAmt := Coins{
		NewInt64Coin("gas-min", 1),
		NewInt64Coin("tree-min", 0),
		NewInt64Coin("mineral-min", 1),
	}
	dup := Coins{
		NewInt64Coin("gas-min", 1),
		NewInt64Coin("gas-min", 1),
		NewInt64Coin("mineral-min", 1),
	}

	cases := []struct {
		coins         Coins
		before, after bool // valid before/after sort
	}{
		{good, true, true},
		{empty, true, true},
		{badSort1, false, true},
		{badSort2, false, true},
		{badAmt, false, true},
		{dup, false, false},
	}

	for tcIndex, tc := range cases {
		require.Equal(t, tc.before, tc.coins.IsValid(), "coin validity is incorrect before sorting, tc #%d", tcIndex)
		tc.coins.Sort()
		require.Equal(t, tc.after, tc.coins.IsValid(), "coin validity is incorrect after sorting, tc #%d", tcIndex)
	}
}

func TestAmountOf(t *testing.T) {
	case0 := Coins{}
	case3 := Coins{
		NewInt64Coin("gold-min", 0),
	}
	case4 := Coins{
		NewInt64Coin("gas-min", 1),
		NewInt64Coin("mineral-min", 1),
		NewInt64Coin("tree-min", 1),
	}
	case5 := Coins{
		NewInt64Coin("mineral-min", 1),
		NewInt64Coin("tree-min", 1),
	}
	case8 := Coins{
		NewInt64Coin("gas-min", 8),
	}

	cases := []struct {
		coins           Coins
		amountOf        int64
		amountOfGAS     int64
		amountOfMINERAL int64
		amountOfTREE    int64
	}{
		{case0, 0, 0, 0, 0},
		{case3, 0, 0, 0, 0},
		{case4, 0, 1, 1, 1},
		{case5, 0, 0, 1, 1},
		{case8, 0, 8, 0, 0},
	}

	for _, tc := range cases {
		assert.Equal(t, NewInt(tc.amountOfGAS), tc.coins.AmountOf("gas-min"))
		assert.Equal(t, NewInt(tc.amountOfMINERAL), tc.coins.AmountOf("mineral-min"))
		assert.Equal(t, NewInt(tc.amountOfTREE), tc.coins.AmountOf("tree-min"))
	}
}

func TestIsCoinNameValid(t *testing.T) {

	cases := []struct {
		name         string
		coinName     string
		expectedPass bool
	}{
		{"standard iris", "iris", true},
		{"standard native", "abc", true},
		{"standard external", "x.abc", true},
		{"standard gateway", "gdex.abc", true},
		{"standard uni", "uni:gdex.abc", true},
		{"with-suffix", "iris-atto", false},
		{"with-suffix-min", "abc-min", false},
	}

	for _, tc := range cases {
		res := IsCoinNameValid(tc.coinName)
		if tc.expectedPass {
			require.True(t, res)
		} else {
			require.False(t, res)
		}
	}
}

func TestIsCoinMinDenomValid(t *testing.T) {

	cases := []struct {
		name         string
		denom        string
		expectedPass bool
	}{
		{"standard iris", "iris-atto", true},
		{"standard native", "abc-min", true},
		{"standard external", "x.abc-min", true},
		{"standard gateway", "gdex.abc-min", true},
		{"standard uni", "uni:gdex.abc-min", true},
		{"non-suffix", "iris", false},
		{"non-suffix !iris", "abc", false},
	}

	for _, tc := range cases {
		res := IsCoinMinDenomValid(tc.denom)
		if tc.expectedPass {
			require.True(t, res)
		} else {
			require.False(t, res)
		}
	}
}
