# Gov Parameters

In IRISnet, there are some special parameters that can be modified through on-chain governance.
All the IRIS holders may participate in the on-chain governance. If the community is not satisfied with certain modifiable parameters, it is available to submit a [parameter-change](../features/governance.md#usage-scenario-of-parameter-change) proposal, and the params will be changed online automatically when the proposal passes.

## Parameters in Auth

| key                           | Description                                                     | Range                     | Current |
| ----------------------------- | --------------------------------------------------------------- | ------------------------- | ------- |
| `auth/MaxMemoCharacters`      | Maximum number of characters in the memo field in a transaction | (0, 18446744073709551615] | 256     |
| `auth/TxSigLimit`             | Maximum number of signatures per transaction                    | (0, 18446744073709551615] | 7       |
| `auth/TxSizeCostPerByte`      | The amount of gas consumed per byte of the transaction          | (0, 18446744073709551615] | 10      |
| `auth/SigVerifyCostED25519`   | Gas spent on edd2519 algorithm signature verification           | (0, 18446744073709551615] | 590     |
| `auth/SigVerifyCostSecp256k1` | Gas spent on secp256k1 algorithm signature verification         | (0, 18446744073709551615] | 1000    |

## Parameters in Bank

| key                       | Description                                        | Range        | Current |
| ------------------------- | -------------------------------------------------- | ------------ | ------- |
| `bank/SendEnabled`        | Tokens that support transfer                       |              | []      |
| `bank/DefaultSendEnabled` | Whether to enable the transfer function by default | {true,false} | true    |

Details in [Bank](../features/bank.md)

## Parameters in Coinswap

| key            | Description | Range | Current              |
| -------------- | ----------- | ----- | -------------------- |
| `coinswap/Fee` | Swap Fee    | (0,1) | 0.003000000000000000 |

Details in [Coinswap](../features/coinswap.md)

## Parameters in Distribution

| key                                | Description                                       | Range        | Current |
| ---------------------------------- | ------------------------------------------------- | ------------ | ------- |
| `distribution/communitytax`        | Fees charged for withdrawal                       | [0, 1]       | 0.02    |
| `distribution/baseproposerreward`  | The base reward rate of the block proposer        | [0, 1]       | 0.01    |
| `distribution/bonusproposerreward` | Reward rate for block proposers                   | [0, 1]       | 0.04    |
| `distribution/withdrawaddrenabled` | Whether to support setting the withdrawal address | {true,false} | true    |

Details in [Distribution](../features/distribution.md)

## Parameters in Gov

| key                 | Description                                      | Range                                                    | Current                                                                                                        |
| ------------------- | ------------------------------------------------ | -------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------- |
| `gov/depositparams` | Related parameters of the deposit mortgage phase | max_deposit_period:(0, 9223372036854775807]              | {"min_deposit": [{"denom": "uiris", "amount": "1000000000"}], "max_deposit_period": "86400s" }                 |
| `gov/votingparams`  | Related parameters of the voting mortgage phase  | voting_period:(0, 9223372036854775807]                   | {"voting_period": "432000s"}                                                                                   |
| `gov/tallyparams`   | Related parameters of the voting tally phase     | quorum:[0,1]<br>threshold:(0,1]<br/>veto_threshold:(0,1] | {"quorum":"0.500000000000000000","threshold": "0.500000000000000000","veto_threshold": "0.330000000000000000"} |

Details in [Governance](../features/governance.md)

## Parameters in IBC

| key                       | Description                             | Range        | Current                            |
| ------------------------- | --------------------------------------- | ------------ | ---------------------------------- |
| `ibc/AllowedClients`      | Clients that support ibc                |              | ["06-solomachine","07-tendermint"] |
| `transfer/SendEnabled`    | Whether to enable the transfer function | {true,false} | false                              |
| `transfer/ReceiveEnabled` | Whether to enable the receive function  | {true,false} | false                              |

## Parameters in Mint

| key              | Description                 | Range    | Current |
| ---------------- | --------------------------- | -------- | ------- |
| `mint/Inflation` | Token issuance frequency    | [0, 0.2] | 0.04    |
| `mint/MintDenom` | Denom of the token mintable |          | uiris   |

Details in [Mint](../features/mint.md)

## Parameters in Service

| key                            | Description                                                 | Range                     | Current                                     |
| ------------------------------ | ----------------------------------------------------------- | ------------------------- | ------------------------------------------- |
| `service/ArbitrationTimeLimit` | Maximum time of dispute resolution                          | (0, 9223372036854775807]  | 120h0m0s                                    |
| `service/ComplaintRetrospect`  | Maximum time for submitting a dispute                       | (0, 9223372036854775807]  | 360h0m0s                                    |
| `service/MaxRequestTimeout`    | Maximum number of blocks to wait for service invocation     | (0, 9223372036854775807]  | 100                                         |
| `service/MinDepositMultiple`   | A multiple of the minimum deposit amount of service binding | (0, 9223372036854775807]  | 1000                                        |
| `service/ServiceFeeTax`        | Tax rate of service fee                                     | [0, 1)                    | 0.05                                        |
| `service/SlashFraction`        | Slash fraction                                              | [0, 1]                    | 0.001                                       |
| `service/TxSizeLimit`          | The limit of the service tx size                            | (0, 18446744073709551615] | 4000                                        |
| `service/MinDeposit`           | Minimum deposit amount                                      | amount: (0, +∞)           | [{"denom": "uiris","amount": "5000000000"}] |
| `service/BaseDenom`            | Token denom that must be used for deposits                  |                           | uiris                                       |

Details in [Service](../features/service.md)

## Parameters in Slashing

| key                                | Description                           | Range                     | Current |
| ---------------------------------- | ------------------------------------- | ------------------------- | ------- |
| `slashing/DowntimeJailDuration`    | Maximum downtime  (continuous)        | (0, 9223372036854775807]  | 10m0s   |
| `slashing/MinSignedPerWindow`      | Minimum signature rate in each window | [0, 1]                    | 0.7     |
| `slashing/SignedBlocksWindow`      | Sliding window for downtime slashing  | (0, 18446744073709551615] | 34560   |
| `slashing/SlashFractionDoubleSign` | Penalty coefficient for double sign   | [0, 1]                    | 0.01    |
| `slashing/SlashFractionDowntime`   | Penalty coefficient for downtime      | [0, 1]                    | 0.0003  |

Details in [Slashing](../features/slashing.md)

## Parameters in Staking

| key                         | Description                                                     | Range                    | Current  |
| --------------------------- | --------------------------------------------------------------- | ------------------------ | -------- |
| `staking/UnbondingTime`     | Mortgage redemption time                                        | (0, 9223372036854775807] | 1814400s |
| `staking/MaxValidators`     | Maximum number of validators                                    | (0, 4294967295]          | 100      |
| `staking/MaxEntries`        | The maximum number of unbinding/redelegation orders in progress | (0, 4294967295]          | 7        |
| `staking/BondDenom`         | Bond denom                                                      |                          | uiris    |
| `staking/HistoricalEntries` | Historical entries                                              | [0, 4294967295]          | 10000    |

Details in [Staking](../features/staking.md)

## Parameters in Token

| key                       | Description                       | Range           | Current                             |
| ------------------------- | --------------------------------- | --------------- | ----------------------------------- |
| `token/TokenTaxRate`      | Base rate for issuing/mint tokens | [0, 1]          | 0.4                                 |
| `token/IssueTokenBaseFee` | Base token for issuing tokens     | amount: (0, +∞) | {"denom": "iris","amount": "60000"} |
| `token/MintTokenFeeRatio` | Rate for mint tokens              | [0, 1]          | 0.1                                 |

Details in [Token](../features/token.md)
