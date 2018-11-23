# Mint User Guide

## Introduction

The incentive mechanism of POW is widely known and explicit: once a new block is produced, the block miner will acquire a certain amount of token as well as the accumulation of transaction fee in the block. As a POS blockchain network, the IRISHUB incentive mechanism is much different. 

As we all know, POW means proof of work. In each block producing period, all miners compete to calculate their work proof and the fastest one will be the winner. Actually, all loser miners don't offer any positive help or collaboration to the winner miner, and they are only the competitors. So it is reasonable to grant all reward to the winner miner. However, in POS blockchain network, we can't do that. Because each block producing process is the collaboration of all validators and delegators, which means the benefit should be share by all these contributors. There are two sources of revenue, one is the transaction fee of the packaged transaction in the block. The other is regular inflation, which will produce new tokens.

As for how to distribute inflation token to contributors, we will document and implement it in distribution module. Here we mainly introduce how to figure out the inflation token and what is the impact to users. 

## Inflation Calculation

Unlike POW network, the inflation token will not be paid to contributors in each block. Only when contributors explicitly send transactions to withdraw reward, then will the inflation token be transfered to users specified addresses. Besides, the token inflation is triggered once an hour, and the new produced token will be saved in global pool. 

### Block Time

The block time is not the machine time, because different machines may not have exactly the same time. They must have some deviation more or less which will result in non-deterministic. So here the block time is the BFT time. Please refer to this [tendermint bft-time](https://github.com/tendermint/tendermint/blob/master/docs/spec/consensus/bft-time.md) for detailed description.

### Inflation Rate

The inflation rate depends on the bonded ratio which means it always changes. The desired bonded ratio 67%. If the ratio is higher, the inflation rate will decrease. In contrast, if the bonded ratio is lower, the inflation rate will increase. Besides, the inflation rate should no more than 20% and no less than 7%. Otherwise, it will be truncated.

Suppose the inflation rate is 10%, and total token amount is 10000iris, then the inflation token will be 0.114iris(10000iris*10%/8766, one year has 8766 hours). After this inflation, the total token amount will be 10000.114iris.

## Impact to users

The inflation calculation is an automatically process. Users have no directly interface to this process. However, users can send delegation or unboud transactions to change the bonded ratio, therefore the inflation rate will change accordingly.

Besides, the inflation process will increate the total token amount. Users can get the total token amount by this command:
```
ubuntu@ubuntu:~$ iriscli stake pool --node=<iris node url>
Pool 
Loose Tokens: 200.1186409166
Bonded Tokens: 400.0000000000
Token Supply: 600.1186409166
Bonded Ratio: 0.6665348695
```

