# Gov Parameters

In IRISnet, there are some special parameters that can be modified through on-chain governance.
All the IRIS holders may participate in the on-chain governance. If the community is not satisfied with certain modifiable parameters, it is available to submit a [parameter-change](../features/governance.md#usage-scenario-of-parameter-change) proposal, and the params will be changed online automatically when the proposal passes.

## Parameters in Auth

| key                      | Description                          | Range               | Current       |
| ------------------------ | ------------------------------------ | ------------------- | ------------- |
| `auth/gasPriceThreshold` | minimum of gas price                 | (0, 10^18iris-atto] | 6000000000000 |
| `auth/txSizeLimit`       | the limitation of the normal tx size | [500, 1500]         | 1000          |

## Parameters in Stake

| key                   | Description                  | Range      | Current  |
| --------------------- | ---------------------------- | ---------- | -------- |
| `stake/MaxValidators` | maximum number of validators | [100, 200] | 100      |
| `stake/UnbondingTime` | unbonding time               | [2week,)   | 504h0m0s |

Details in [Stake](../features/stake.md)

## Parameters in Distribution

| key                         | Description                                    | Range     | Current |
| --------------------------- | ---------------------------------------------- | --------- | ------- |
| `distr/BaseProposerReward`  | standard ratio of the block reward             | (0, 0.02] | 0.01    |
| `distr/BonusProposerReward` | maximum additional bonus ratio                 | (0, 0.08] | 0.04    |
| `distr/CommunityTax`        | proportion of contributions to community funds | (0, 0.2]  | 0.02    |

Details in [Distribution](../features/distribution.md)

## Parameters in Mint

| key              | Description           | Range    | Current |
| ---------------- | --------------------- | -------- | ------- |
| `mint/Inflation` | Inflation coefficient | [0, 0.2] | 0.04    |

Details in [Mint](../features/mint.md)

## Parameters in Slashing

| key                                | Description                               | Range         | Current |
| ---------------------------------- | ----------------------------------------- | ------------- | ------- |
| `slashing/CensorshipJailDuration`  | Censorship Jail Duration                  | (0, 4week)    | 48h0m0s |
| `slashing/DoubleSignJailDuration`  | DoubleSign Jail Duration                  | (0, 4week)    | 48h0m0s |
| `slashing/DowntimeJailDuration`    | Downtime Jail Duration                    | (0, 4week)    | 36h0m0s |
| `slashing/MaxEvidenceAge`          | Acceptable earliest time of the evidence  | [1day, +∞)    | 51840   |
| `slashing/MinSignedPerWindow`      | Minimum voting ratio in the slash window  | [0.5, 0.9]    | 0.7     |
| `slashing/SignedBlocksWindow`      | The number of blocks for slash statistics | [100, 140000] | 34560   |
| `slashing/SlashFractionCensorship` | Slash ratio of censorship                 | [0.005, 0.1]  | 0       |
| `slashing/SlashFractionDoubleSign` | Slash ratio of double-sign                | [0.01, 0.1]   | 0.01    |
| `slashing/SlashFractionDowntime`   | Slash ratio of downtime                   | [0.005, 0.1]  | 0.0003  |

Details in [Slashing](../features/slashing.md)

## Parameters in Asset

| key                          | Description                                      | Range   | Current                           |
| ---------------------------- | ------------------------------------------------ | ------- | --------------------------------- |
| `asset/AssetTaxRate`         | Asset tax rate                                   | [0, 1]  | 0.4                               |
| `asset/IssueTokenBaseFee`    | Base fee for issuing tokens                      | [0, +∞) | 60000000000000000000000iris-atto  |
| `asset/MintTokenFeeRatio`    | Fee ratio for minting tokens                     | [0, 1]  | 0.1                               |
| `asset/CreateGatewayBaseFee` | Base fee for creating gateway tokens             | [0, +∞) | 120000000000000000000000iris-atto |
| `asset/GatewayAssetFeeRatio` | Fee ratio for issuing and minting gateway tokens | [0, 1]  | 0.1                               |

Details in [Asset](../features/asset.md)

## Parameters in Service

| key                            | Description                                                 | Range            | Current  |
| ------------------------------ | ----------------------------------------------------------- | ---------------- | -------- |
| `service/ArbitrationTimeLimit` | maximum time of dispute resolution                          | [5days, 10days]  | 120h0m0s |
| `service/ComplaintRetrospect`  | maximum time for submitting a dispute                       | [15days, 30days] | 360h0m0s |
| `service/MaxRequestTimeout`    | maximum number of blocks to wait for service invocation     | [20, +∞)         | 100      |
| `service/MinDepositMultiple`   | a multiple of the minimum deposit amount of service binding | [500, 5000]      | 1000     |
| `service/ServiceFeeTax`        | tax rate of service fee                                     | (0, 0.2]         | 0.01     |
| `service/SlashFraction`        | slash fraction                                              | (0, 0.01]        | 0.001    |
| `service/TxSizeLimit`          | the limit of the service tx size                            | [2000, 6000]     | 4000     |

Details in [Service](../features/service.md)
