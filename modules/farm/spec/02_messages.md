<!--
order: 2
-->

# Messages

## MsgCreatePool

The farm pool can be created by any user via a `MsgCreatePool` message.

```go
type MsgCreatePool struct {
    Name           string
    Description    string
    LpTokenDenom   string
    StartHeight    int64
    RewardPerBlock sdk.Coins
    TotalReward    sdk.Coins
    Editable       bool
    Creator        string
}
```

This message is expected to fail if:

- `LpTokenDenom` does not comply with the rules specified on the chain.
- the name of farm pool has exist.
- `StartHeight` is less than the current block height.
- `TotalReward` is less than `RewardPerBlock`.
- the balance of creator is not enough to pay `CreatePoolFee+TotalReward`.
- The length of `TotalReward` is greater than `MaxRewardCategoryN`.

In addition, the `Endheight = StartHeight + TotalReward/RewardPerBlock`. Because there may be multiple tokens for event rewards, the end heights may be inconsistent. In order to reduce the complexity of the system, take the smallest value among all heights as the final end height. After the event ends, the remaining bonuses will be refunded to creator's account.

At the beginning of the activity, because there was no user participating, so `RewardPerShare=0`, which means that the user has no income from the beginning of the activity to the user's first stake, and every time the user's `stake`、 `unstake` 、`harvest` will trigger
in the calculation of `RewardPerShare` (calculate the income that each lptoken can obtain before), the user's previous income is equal to `lastTotalLocked * RewardPerShare - lastDebt`, after the user gets back the income, record the user's current total debt (total income that has been withdrawn, lastDebt) is equal to `currentTotalLocked * RewardPerShare`.

## MsgDestroyPool

The farm pool can be destroyed by creator via a `MsgDestroyPool` message.

```go
type MsgDestroyPool struct {
    PoolName string
    Creator  string
}
```

This message is expected to fail if:

- the farm pool is not exist.
- the `Creator` is not the creator of the farm pool.
- the farm pool activity has ended.

When creator actively destroys the pool, it will also trigger the recalculation of `RewardPerShare`, calculate the current value of each `lpToken` (relative to the reward), and calculate the current remaining bonus in the pool, and return the bonus to creator to remove the pool from the active queue.

When the pool becomes inactive (expired or ended), the user can only unstake his staking `lptoken` from the pool and get rewards, or it can be retrieved multiple times, but the reward obtained are the same as all at once. note that `RewardPerShare` will no longer be updated at this time

## MsgAdjustPool

Creators can use `MsgAdjustPool` to add reward to the pool before the end of the activity to achieve the purpose of extending the activity.

```go
type MsgAdjustPool struct {
    PoolName string
    AdditionalReward sdk.Coins
    RewardPerBlock   sdk.Coins
    Creator  string
}
```

This message is expected to fail if:

- the farm pool is not exist.
- the `Creator` is not the creator of the farm pool.
- the farm pool activity has ended.
- additional reward types are not within the scope of the pool definition

When the creator adds bonuses to the pool, it is equivalent to extending the end height of the pool and does not affect the user's previous earnings.

## MsgStake

Any user can retrieve the staking `lpToken` through `MsgStake` and trigger the return of reward.

```go
type MsgStake struct {
    PoolName string
    Amount   sdk.Coin
    Sender   string
}
```

This message is expected to fail if:

- the farm pool is not exist.
- the farm activity has not yet started.
- the farm activity has ended.
- the `lpToken` staked by the user is not specified by the pool

Each stake operation of the user will trigger the update of the pool (including `RemainingReward`, `RewardPerShare`, `TotalLpTokenLocked`), but only the rewards of the current user will be retrieved. The rewards of other users are collected by accumulating `RewardPerShare`. When the rewards are issued, Will update the current user’s debt for the next calculation of revenue.

## MsgUnstake

Any user can retrieve their staked `lpToken` through `MsgUnstake`.

```go
type MsgUnstake struct {
    PoolName string
    Amount   sdk.Coin
    Sender   string
}
```

This message is expected to fail if:

- the farm pool is not exist.
- the farmer information is not exist.
- the amount of `lpToken` retrieved by users is greater the amount he has staked.

When the user `Unstake`, there may be two situations, one is that the current farm activity has ended, because the activity has ended, all rewards have been solidified, so `RewardPerShare` will not be updated again, only `TotalLpTokenLocked` and `RemainingReward` will be updated, and the reward will be calculated When the activity ends, use the `RewardPerShare` at the end of the activity to calculate the revenue; the other is that the current farm activity is in progress. In this case, it has been described above, and normal update of the pool is enough.

## MsgHarvest

Any user can get back the rewards through `MsgHarvest`. The only difference from `MsgUnstake` is that you don’t need to only get back the revenue and not get back the `lptoken`.

```go
type MsgHarvest struct {
    PoolName string
    Sender   string
}
```

This message is expected to fail if:

- the farm pool is not exist.
- the farmer information is not exist.
- the farm activity has ended.
