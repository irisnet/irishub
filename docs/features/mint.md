# Mint User Guide

## Introduction

The incentive mechanism of POW is widely known and explicit: once a new block is produced, the block miner will acquire a certain amount of token as well as the accumulation of transaction fee in the block. As a POS blockchain network, the IRISHUB has similar way to produce reward token but much different mechanism to distribute the reward to each contributor.

In each block producing period, all POW miners compete to calculate their work proof and the fastest one will be the winner. Actually, all loser miners don't offer any positive help or collaboration to the winner miner, and they are only the competitors. So it is reasonable to grant all reward to the winner miner. However, in POS blockchain network, we can't do that. Because each block producing process is the collaboration of all validators and delegators, which means the benefit should be share by all these contributors. As for how to distribute reward token to contributors, we will document and implement it in distribution module.

The reward is composed by two parts, one is the collected transaction fee from the transactions in each block. Another one is regular inflation in each block, which will produce new tokens. The mint module is in charge of calculating the inflated token amount and add the inflated token to reward pool.

## Inflation Calculation

### Block Time

The block time is not the machine time, because different machines may not have exactly the same time. They must have some deviation more or less which will result in non-deterministic. So here the block time is the BFT time. Please refer to this [tendermint bft-time](https://github.com/tendermint/tendermint/blob/master/docs/spec/consensus/bft-time.md) for detailed description.

### Inflation Rate

The inflation rate is assigned to 4% per year in genesis file. This value can be modified by governance. As for how to change the value by governance, please refer to [governance](governance.md).

### Calculation

This is the calculation equation:
```
 blockCostTime  = (current block BFT time) - (last block BFT time)
 AnnualInflationAmount = inflationBasement * inflationRate
 blockInflationAmount = AnnualInflationAmount * blockCostTime / (year)
```
The value of `inflationBasement` is specified in genesis file. By default its value `2000000000iris`(2 billion iris, `1 iris` equal `1*10^18 iris-atto`), and its value will never be changed.
Suppose `blockCostTime` is 5000 milisecond, and `inflationRate` is `4%`, then the inflation amount will be `12675235125611580094iris-atto` (`12.675235125611580094iris`)

## Impact to users

The inflation calculation is automatically trigged by each block. So once a new block is produces, new tokens will be created and the loosen token will increase accordingly. Users have no directly interface to affect this process. 

Here we provide two interfaces to query the total loosen token.

1. `iriscli stake pool`
```
ubuntu@ubuntu:~$ iriscli stake pool --node=<iris node url>
Pool
Loose Tokens: 1846663.900384156921391687
Bonded Tokens: 425182.329615843078608313
Token Supply: 2271846.230000000000000000
Bonded Ratio: 0.187152776500000000
```

2. `iriscli bank token-stats`
```
ubuntu@ubuntu:~$ iriscli bank token-stats --trust-node=false --chain-id [chain-id] --node=[iris node url]
{
  "loosen_token": [
    "1864477.596384156921391687iris"
  ],
  "burned_token": null,
  "bonded_token": "425182.329615843078608313iris"
}
```
