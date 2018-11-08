# Introduction

The incentive mechanism of POW is widely known and explicit: once a new block is produced, the block miner will acquire a certain amount of token which is the reward for producing block. As a POS blockchain network, the Irishub also has incentive mechanism. It is much different. Strictly speaking, it is the inflation token which will be distributed to all the contributors.

As we all know, POW means proof of work. In each producing period, all miner compete to submit their proof of work and the fastest one will be the winner. Actually, all loser miners don't offer any help or collaboration to the winner miner, and they are only the competitors. So it is reasonable to grant all reward of producing block to the winner miner. However, in POS blockchain network, we can't do that. Because each block is the collaboration consequent of all validators and delegators, which means the reward should be share by all contributors.

The detailed distribution mechanism will be documented and implemented by distribution module. Here we mainly introduce how to figure out the inflation token and what is the impact to users. 


# Reward Calculation

Unlike POW network, the reward will not be paid to contributors in each block. Instead, once an hour(block time) the reward is calculated and saved in global pool. Only when contributors send transactions to withdraw reward, then will the reward tokens be transfered to users specified addresses.

## Block Time

The block time is not the machine time, because different machine may not have exactly the same time. They must have some deviation more or less which will result in non-deterministic. So here the block time is the BFT time. Please refer to this [tendermint bft-time](https://github.com/tendermint/tendermint/blob/master/docs/spec/consensus/bft-time.md) for detailed description.

## Inflation Rate

The inflation rate depends on the bonded ratio which means it always changes. The desired bonded ratio 67%. If the ratio is higher, the inflation rate will decrease. in contrast, if the bonded ratio is lower than 67%, the inflation rate will increase. Besides, the inflation rate should no more than 20% and no less than 7%. Otherwise, it will be truncated. For detailed algorithm please refer to [cosmos mint](https://github.com/cosmos/cosmos-sdk/blob/develop/docs/spec/mint/begin_block.md).

# Impact to users

The inflation calculation is an automatically process. Users have no directly interaface to this process. However, users can send delegate or unbond transactions to change the bonded ratio, therefore the inflation rate will change accordingly.

Besides, the inflation process will increate the total token amount. Users can get the total token amount by this command:
```
ubuntu@ubuntu:~$ iriscli stake pool --node=<iris node url>
Pool 
Loose Tokens: 200.1186409166
Bonded Tokens: 400.0000000000
Token Supply: 600.1186409166
Bonded Ratio: 0.6665348695
```

