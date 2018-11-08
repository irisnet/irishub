# Introduction

In POW blockchain network, once a new block is produced, the block miner will acquire a certain amount of token which is the producing block reward. Irishub is a POS blockchain network, it does have reward, but we call it inflation token instead of producing block reward.

As we all know, POW means proof of work. Each miner will compete to submit its proof of work and the fastest one is the winner. All the loser miners at a given height don't offer any help or collaboration to the winner miner, and they are only the competitors. So it is reasonable to grant all producing block reward to the winner miner. However, in POS blockchain network, each block is the collaboration consequent of all validators and delegators, which means the inflation token should be share by all contributors.

The detailed distribution mechanism will be documented and implemented by distribution module. Here we mainly introduce how to figure out the inflation token and what is the impact to users. 


# Reward Calculation

Unlike POW network, we won't calculate inflation in each block, instead, we do it once an hour(block time).

## Block Time

The block time is not the machine time, because different machine may not have the same time. They must have some deviation more or less which will result in non-deterministic. So here the block time is the BFT time. Please refer to this [tenermint bft-time](https://github.com/tendermint/tendermint/blob/master/docs/spec/consensus/bft-time.md) for detailed description.

## Inflation Rate

The inflation rate depends on the bonded ratio which mean it always changes. The desired bonded ratio 67%. If the ratio is higher, the inflation rate will decrease. in contrast, if the bonded ratio is lower than 67%, the inflation rate will increase. Besides, the inflation rate should no more than 20% and no less than 7%. For detailed algorithm please refer to [cosmos mint](https://github.com/cosmos/cosmos-sdk/blob/develop/docs/spec/mint/begin_block.md).

# Impact to users
