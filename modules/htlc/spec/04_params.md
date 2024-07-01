<!--
order: 4
-->

# Parameters

The htlc module contains the following parameters:

| Key         | Type         | Example |
| :---------- | :----------- | :------ |
| AssetParams | []AssetParam |         |

```go
type Params struct {
    AssetParams []AssetParam
}

type AssetParam struct {
    Denom         string
    SupplyLimit   SupplyLimit
    Active        bool
    DeputyAddress string
    FixedFee      math.Int
    MinSwapAmount math.Int
    MaxSwapAmount math.Int
    MinBlockLock  uint64
    MaxBlockLock  uint64
}

type SupplyLimit struct {
    Limit          math.Int
    TimeLimited    bool
    TimePeriod     time.Duration
    TimeBasedLimit math.Int
}
```
