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
    FixedFee      sdk.Int
    MinSwapAmount sdk.Int
    MaxSwapAmount sdk.Int
    MinBlockLock  uint64
    MaxBlockLock  uint64
}

type SupplyLimit struct {
    Limit          sdk.Int
    TimeLimited    bool
    TimePeriod     time.Duration
    TimeBasedLimit sdk.Int
}
```
