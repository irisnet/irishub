# Params

Params module allows you to query the system parameters which can be governed (except the gov params) by the [gov module](./gov.md).

```bash
iris query params subspace [subspace] [key] [flags]
```

Subspace currently supports the following:`auth`, `bank`, `staking`, `mint`, `distribution`, `slashing`, `gov`, `crisis`, `token`, `record`, `htlc`, `coinswap`, `service`.

Among them, the parameters available for query for each subspace are as follows:

## auth

| key                      | description                                                     | default |
| ------------------------ | --------------------------------------------------------------- | ------- |
| `MaxMemoCharacters`      | Maximum number of characters in the memo field in a transaction | 256     |
| `TxSigLimit`             | Maximum number of signatures per transaction                    | 7       |
| `TxSizeCostPerByte`      | The amount of gas consumed per byte of the transaction          | 10      |
| `SigVerifyCostED25519`   | Gas spent on edd2519 algorithm signature verification           | 590     |
| `SigVerifyCostSecp256k1` | Gas spent on secp256k1 algorithm signature verification         | 1000    |

## bank

| key                  | description                                        | default |
| -------------------- | -------------------------------------------------- | ------- |
| `SendEnabled`        | Tokens that support transfer                       | {}      |
| `DefaultSendEnabled` | Whether to enable the transfer function by default | true    |

## staking

| key                 | description                                                     | default   |
| ------------------- | --------------------------------------------------------------- | --------- |
| `UnbondingTime`     | Mortgage redemption time                                        | 3w(weeks) |
| `MaxValidators`     | Maximum number of validators                                    | 100       |
| `MaxEntries`        | The maximum number of unbinding/redelegation orders in progress | 7         |
| `BondDenom`         | Bond denom                                                      | uiris     |
| `HistoricalEntries` |                                                                 | 100       |

## mint

| key         | description                 | default |
| ----------- | --------------------------- | ------- |
| `Inflation` | Token issuance frequency    | 0.04    |
| `MintDenom` | Denom of the token mintable | uiris   |

## distribution

| key                   | description                                       | default |
| --------------------- | ------------------------------------------------- | ------- |
| `communitytax`        | Fees charged for withdrawal                       | 0.02    |
| `baseproposerreward`  | The base reward rate of the block proposer        | 0.01    |
| `bonusproposerreward` | Reward rate for block proposers                   | 0.04    |
| `withdrawaddrenabled` | Whether to support setting the withdrawal address | true    |

## slashing

| key                       | description                           | default |
| ------------------------- | ------------------------------------- | ------- |
| `SignedBlocksWindow`      | Sliding window for downtime slashing  | 100     |
| `MinSignedPerWindow`      | Minimum signature rate in each window | 0.5     |
| `DowntimeJailDuration`    | Maximum downtime  (continuous)        | 10m     |
| `SlashFractionDoubleSign` | Penalty coefficient for double sign   | 1/20    |
| `SlashFractionDowntime`   | Penalty coefficient for downtime      | 1/100   |

## gov

| key             | description                                      | default                                                         |
| --------------- | ------------------------------------------------ | --------------------------------------------------------------- |
| `depositparams` | Related parameters of the deposit mortgage phase | `min_deposit`:    10000000uiris; `max_deposit_period`: 2d(days) |
| `votingparams`  | Related parameters of the voting mortgage phase  | `voting_period`: 2d(days)                                       |
| `tallyparams`   | Related parameters of the voting tally phase     | `quorum`: 0.334; `threshold`: 0.5; `veto_threshold`: 0.334      |

## crisis

| key           | description  | default   |
| ------------- | ------------ | --------- |
| `ConstantFee` | Constant Fee | 1000uiris |

## token

| key                 | description                       | default            |
| ------------------- | --------------------------------- | ------------------ |
| `TokenTaxRate`      | Base rate for issuing/mint tokens | 0.4                |
| `IssueTokenBaseFee` | Base token for issuing tokens     | 60000 * 10^6 uiris |
| `MintTokenFeeRatio` | Rate for mint tokens              | 0.1                |

## coinswap

| key             | description                   | default |
| --------------- | ----------------------------- | ------- |
| `StandardDenom` | The name of the token charged | uiris   |

## htlc

| key           | description                              | default |
| ------------- | ---------------------------------------- | ------- |
| `AssetParams` | Array of supported assets,`[]AssetParam` | None    |

AssetParam参数如下：

| key                        | description                       | Example                                      |
| -------------------------- | --------------------------------- | -------------------------------------------- |
| `AssetParam.Denom`         | The name  of asset                | "htltbcbnb"                                  |
| `AssetParam.SupplyLimit`   | The supply limit of  asset        | 100000                                       |
| `AssetParam.Active`        |                                   | Asset state: live or paused                  | true |
| `AssetParam.DeputyAddress` | Deputy's IRISHub address          | “iaa18n3x722r4jpwmshlxnw3ehlpfzywupzefthcz5” |
| `AssetParam.FixedFee`      | Deputy's fixed fee on other chain | 1000                                         |
| `AssetParam.MinSwapAmount` | Minimum swap amount               | 1                                            |
| `AssetParam.MaxSwapAmount` | Maximum swap amount               | 10000                                        |
| `AssetParam.MinBlockLock`  | Minimum swap expire height        | 50                                           |
| `AssetParam.MaxBlockLock`  | Maximum swap expire height        | 25480                                        |

## service

| key                         | description                                         | default    |
| --------------------------- | --------------------------------------------------- | ---------- |
| `MaxRequestTimeout`         | Maximum service request timeout                     | 100(block) |
| `MinDepositMultiple`        | Minimum deposit multiple                            | 200        |
| `MinDeposit`                | Minimum deposit amount                              | 6000uiris  |
| `ServiceFeeTax`             | Service rate                                        | 0.1        |
| `SlashFraction`             | Slash fraction                                      | 0.001      |
| `ComplaintRetrospect`       | Complaint retrospect                                | 15d        |
| `ArbitrationTimeLimit`      | Arbitration period                                  | 5d         |
| `TxSizeLimit`               | The maximum number of bytes per service transaction | 4000       |
| `BaseDenom`                 | Tokens supported by service fees                    | uiris      |
| `RestrictedServiceFeeDenom` | Restricted service fee denom                        | false      |
