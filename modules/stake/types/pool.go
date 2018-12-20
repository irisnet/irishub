package types

import (
	"bytes"
	"fmt"

	"github.com/irisnet/irishub/modules/bank"
	sdk "github.com/irisnet/irishub/types"
)

type Pool struct {
	BondedTokens sdk.Dec `json:"bonded_tokens"` // reserve of bonded tokens
}

// nolint
func (p Pool) Equal(p2 PoolMgr) bool {
	bz1 := MsgCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := MsgCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

// initial pool for testing
func InitialPool() Pool {
	return Pool{
		BondedTokens: sdk.ZeroDec(),
	}
}

//_______________________________________________________________________

// Pool - dynamic parameters of the current state
type PoolMgr struct {
	BankKeeper   bank.Keeper
	Pool 		 Pool
}

func (p PoolMgr) increaseBondedToken(ctx sdk.Context, bondedTokens sdk.Dec) PoolMgr {
	round := bondedTokens.TruncateInt()
	change := bondedTokens.Sub(sdk.NewDecFromInt(round))
	if !change.IsZero() {
		panic("token is not integer")
	}

	p.Pool.BondedTokens = p.Pool.BondedTokens.Add(bondedTokens)
	balance := sdk.NewCoin(StakeDenom, round)
	p.BankKeeper.DecreaseLoosenToken(ctx, sdk.Coins{balance})
	return p
}

func (p PoolMgr) decreaseBondedTokens(ctx sdk.Context, bondedTokens sdk.Dec) PoolMgr {
	round := bondedTokens.TruncateInt()
	change := bondedTokens.Sub(sdk.NewDecFromInt(round))
	if !change.IsZero() {
		panic("token is not integer")
	}

	p.Pool.BondedTokens = p.Pool.BondedTokens.Sub(bondedTokens)
	balance := sdk.NewCoin(StakeDenom, round)
	p.BankKeeper.IncreaseLoosenToken(ctx, sdk.Coins{balance})
	return p
}

type PoolStatus struct {
	LooseTokens  sdk.Dec `json:"loose_tokens"`  // tokens which are not bonded in a validator
	BondedTokens sdk.Dec `json:"bonded_tokens"` // reserve of bonded tokens
}

// Sum total of all staking tokens in the pool
func (p PoolStatus) TokenSupply() sdk.Dec {
	return p.LooseTokens.Add(p.BondedTokens)
}

// get the bond ratio of the global state
func (p PoolStatus) BondedRatio() sdk.Dec {
	supply := p.TokenSupply()
	if supply.GT(sdk.ZeroDec()) {
		return p.BondedTokens.Quo(supply)
	}
	return sdk.ZeroDec()
}

// HumanReadableString returns a human readable string representation of a
// pool.
func (p PoolStatus) HumanReadableString() string {

	resp := "Pool \n"
	resp += fmt.Sprintf("Loose Tokens: %s\n", p.LooseTokens)
	resp += fmt.Sprintf("Bonded Tokens: %s\n", p.BondedTokens)
	resp += fmt.Sprintf("Token Supply: %s\n", p.TokenSupply())
	resp += fmt.Sprintf("Bonded Ratio: %v\n", p.BondedRatio())
	return resp
}
