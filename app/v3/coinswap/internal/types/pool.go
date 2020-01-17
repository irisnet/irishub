package types

import (
	"fmt"

	sdk "github.com/irisnet/irishub/types"
)

type Pool struct {
	Coins sdk.Coins `json:"coins"`
	Name  string    `json:"name"`
}

func NewPool(name string, balance sdk.Coins) Pool {
	return Pool{
		Coins: balance,
		Name:  name,
	}
}

func (p Pool) String() string {
	return fmt.Sprintf(`
  Pool:  %s
  Coins: %s`, p.Name, p.Coins.String())
}

func (p Pool) BalanceOf(denom string) sdk.Int {
	return p.Coins.AmountOf(denom)
}

func (p Pool) Balance() sdk.Coins {
	return p.Coins
}

func (p *Pool) Add(coin sdk.Coins) *Pool {
	p.Coins = p.Coins.Add(coin)
	return p
}

func (p *Pool) Sub(coin sdk.Coins) *Pool {
	p.Coins = p.Coins.Sub(coin)
	return p
}
