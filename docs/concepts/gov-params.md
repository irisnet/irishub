# Gov Parameters

In IRISnet, there are some special parameters that can be modified through on-chain governance.
All the IRIS holders may participate in the on-chain governance. If the community is not satisfied with certain modifiable parameters, it is available to submit a [parameter-change](../features/governance.md#usage-scenario-of-parameter-change) proposal, and the params will be changed online automatically when the proposal passes.

## Parameters in Auth

| key                           | Description                  | Range   | Current |
| ----------------------------- | ---------------------------- | ------- | ------- |
| `auth/MaxMemoCharacters`      | maximum number of characters | (0, +∞) | 256     |
| `auth/TxSigLimit`             | tx signature limit           | (0, +∞) | 7       |
| `auth/TxSizeCostPerByte`      |                              | (0, +∞) | 10      |
| `auth/SigVerifyCostED25519`   |                              | (0, +∞) | 590     |
| `auth/SigVerifyCostSecp256k1` |                              | (0, +∞) | 1000    |

## Parameters in Bank

| key                       | Description | Range        | Current |
| ------------------------- | ----------- | ------------ | ------- |
| `bank/SendEnabled`        |             |              | []      |
| `bank/DefaultSendEnabled` |             | {true,false} | true    |

Details in [Bank](../features/bank.md)

## Parameters in Coinswap

| key            | Description | Range | Current              |
| -------------- | ----------- | ----- | -------------------- |
| `coinswap/Fee` |             | (0,1) | 0.003000000000000000 |

Details in [Coinswap](../features/coinswap.md)

## Parameters in Crisis

| key                  | Description | Range   | Current                                |
| -------------------- | ----------- | ------- | -------------------------------------- |
| `crisis/ConstantFee` |             | amount: (0, +∞) | {"denom": "uiris",   "amount": "1000"} |

## Parameters in Distribution

| key                                | Description                                    | Range        | Current |
| ---------------------------------- | ---------------------------------------------- | ------------ | ------- |
| `distribution/communitytax`        | proportion of contributions to community funds | (0, 0.2]     | 0.02    |
| `distribution/baseproposerreward`  | standard ratio of the block reward             | (0, 0.02]    | 0.01    |
| `distribution/bonusproposerreward` | maximum additional bonus ratio                 | (0, 0.08]    | 0.04    |
| `distribution/withdrawaddrenabled` |                                                | {true,false} | true    |

Details in [Distribution](../features/distribution.md)

## Parameters in Gov

| key                 | Description | Range                                                    | Current                                                      |
| ------------------- | ----------- | -------------------------------------------------------- | ------------------------------------------------------------ |
| `gov/depositparams` |             | max_deposit_period:(0, +∞)                               | {"min_deposit": [{"denom": "uiris", "amount": "1000000000"}], "max_deposit_period": "86400s" }, |
| `gov/votingparams`  |             | voting_period:(0, +∞)                                    | {"voting_period": "432000s"}                                 |
| `gov/tallyparams`   |             | quorum:[0,1]<br>threshold:(0,1]<br/>veto_threshold:(0,1] | {"quorum":"0.500000000000000000","threshold": "0.500000000000000000","veto_threshold": "0.330000000000000000"} |

Details in [Governance](../features/governance.md)

## Parameters in ibc

| key              | Description           | Range    | Current |
| ---------------- | --------------------- | -------- | ------- |
| `ibc/AllowedClients` |  |       | ["06-solomachine","07-tendermint"] |

## Parameters in Mint

| key              | Description           | Range    | Current |
| ---------------- | --------------------- | -------- | ------- |
| `mint/Inflation` | Inflation coefficient | [0, 0.2] | 0.04    |
| `mint/MintDenom` |                       |          | uiris   |

Details in [Mint](../features/mint.md)

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
| `service/MinDeposit`          |                                                             |                  | [{"denom": "uiris","amount": "5000000000"}] |
| `service/BaseDenom`          |                                                             |                  | uiris |

Details in [Service](../features/service.md)

## Parameters in Slashing

| key                                | Description                               | Range         | Current |
| ---------------------------------- | ----------------------------------------- | ------------- | ------- |
| `slashing/DowntimeJailDuration`    | Downtime Jail Duration                    | (0, 4week)    | 0h10m0s |
| `slashing/MinSignedPerWindow`      | Minimum voting ratio in the slash window  | [0.5, 0.9]    | 0.7     |
| `slashing/SignedBlocksWindow`      | The number of blocks for slash statistics | [100, 140000] | 34560   |
| `slashing/SlashFractionDoubleSign` | Slash ratio of double-sign                | [0.01, 0.1]   | 0.01    |
| `slashing/SlashFractionDowntime`   | Slash ratio of downtime                   | [0.005, 0.1]  | 0.0003  |

Details in [Slashing](../features/slashing.md)

## Parameters in Staking

| key                       | Description                  | Range   | Current          |
| ------------------------- | ---------------------------- | ------- | ---------------- |
| `staking/UnbondingTime`      | unbonding time               | [2week,)   | 1814400s |
| `staking/MaxValidators` | maximum number of validators | [100, 200] | 100 |
| `staking/MaxEntries` |                              | (0, +∞) | 7 |
| `staking/BondDenom` |                              |            | uiris |
| `staking/HistoricalEntries` |                              | [0, +∞) | 10000 |

Details in [Staking](../features/staking.md)

## Parameters in Token

| key                       | Description                  | Range   | Current          |
| ------------------------- | ---------------------------- | ------- | ---------------- |
| `token/TokenTaxRate`      | Asset tax rate               | [0, 1]  | 0.4              |
| `token/IssueTokenBaseFee` | Base fee for issuing tokens  | [0, +∞) | 60000000000uiris |
| `token/MintTokenFeeRatio` | Fee ratio for minting tokens | [0, 1]  | 0.1              |

Details in [Token](../features/token.md)

## Parameters in Transfer

| key                       | Description                  | Range   | Current          |
| ------------------------- | ---------------------------- | ------- | ---------------- |
| `transfer/SendEnabled`      |                | {true,false} | true |
| `transfer/ReceiveEnabled` |             | {true,false} | true |
