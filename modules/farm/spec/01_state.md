<!--
order: 1
-->

# State

## Params

`Params` defines the parameters that the farm module can manage through the `gov` module.

```go
type Params struct {
    CreatePoolFee      sdk.Coin 
    MaxRewardCategories uint32                                  
}
```

Parameters are stored in a global GlobalParams KVStore.

- `CreatePoolFee`: the cost of creating a farm pool, which will be allocated to the validator or delegator
- `MaxRewardCategories`: the farm pool can be set to reward how many types of tokens

## FarmPool

`FarmPool` records all the detailed information of the current pool, including the total amount of the pool, balance, etc.

```go
type FarmPool struct {
    Name                   string                                  
    Creator                string                                  
    Description            string                                  
    StartHeight            int64                                  
    EndHeight              int64                                  
    LastHeightDistrRewards int64                                  
    Editable               bool                                    
    TotalLpTokenLocked     sdk.Coin 
    Rules                  []RewardRule                            
}

type RewardRule struct {
    Reward          string
    TotalReward     sdk.Int
    RemainingReward sdk.Int
    RewardPerBlock  sdk.Int
    RewardPerShare  sdk.Dec
}
```

- `Name`: the name of the farm pool, globally unique.
- `Creator`: the creator of farm pool, but also the provider of rewards and fees.
- `Description`: detailed description of farm pool.
- `StartHeight`: the starting height of the farm pool activity, but the user's reward is not calculated from this height, but calculated from the moment the user staking.
- `EndHeight`: the end height of the farm pool activity. After this height, users can no longer perform stake transactions, and the reward ends after this height. The activity will be removed from the active farm pool. If there are remaining bonuses, will be refunded to the creator of the pool.
- `LastHeightDistrRewards`: `LastHeightDistrRewards` records the height of the pool that triggered the reward distribution last time. When the reward distribution is triggered next time, it will use `LastHeightDistrRewards` as the starting height and the current height as the ending height. The total rewards generated during this time period are calculated.
- `Editable`: whether the farm pool can be actively destroyed by the creator, after the farm pool is destroyed, the profit calculation ends, and the remaining money is returned to the creator.
- `TotalLpTokenLocked`: the farm pool accepts collateralized token denom, and the denom rules can be set by the users of moudle.

## RewardRule

`RewardRule` defines the rules for the pool to distribute rewards and record the remaining bonuses of the current pool.

```go
type RewardRule struct {
    Reward          string
    TotalReward     sdk.Int
    RemainingReward sdk.Int
    RewardPerBlock  sdk.Int
    RewardPerShare  sdk.Dec
}
```

- `Reward`: denom of rewards distribution.
- `TotalReward`: total amount of bonuses issued.
- `RemainingReward`: the remaining amount of the bonuses.
- `RewardPerBlock`: amount of rewards issued for each block.
- `RewardPerShare`: the current amount of rewards that each lptoken can get.

## FarmInfo

`FarmInfo` records the user's stake information.

```go
type FarmInfo struct {
    PoolName   string
    Address    string
    Locked     sdk.Int
    RewardDebt sdks.Coins
}
```

- `PoolName`: the name of farm pool.
- `Address`: the address of farmer.
- `Locked`: the total amount of user staked
- `RewardDebt`: user's total debt.

Every time the user triggers the return of earnings, the `RewardDebt` will be updated. When all `lpToken` is retrieved, `FarmInfo` is deleted.
