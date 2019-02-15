# iriscli stake

## Introduction

Stake module provides a set of subcommands to query staking state and send staking transactions.

## Usage

```
iriscli stake [subcommand] [flags]
```

Print all supported subcommands and flags:
```
iriscli stake --help
```

## Available Commands

| Name                            | Description                                                   |
| --------------------------------| --------------------------------------------------------------|
| [validator](validator.md)       | Query a validator                                             |
| [validators](validators.md)     | Query for all validators                                      |
| [delegation](delegation.md)     | Query a delegation based on address and validator address     |
| [delegations](delegations.md)   | Query all delegations made from one delegator                 |
| [delegations-to](delegations-to.md)   | Query all delegations to one validator                 |
| [unbonding-delegation](unbonding-delegation.md)               | Query an unbonding-delegation record based on delegator and validator address                 |
| [unbonding-delegations](unbonding-delegations.md)             | Query all unbonding-delegations records for one delegator                                     |
| [unbonding-delegations-from](unbonding-delegations-from.md)   | Query all unbonding delegatations from a validator                                            |
| [redelegations-from](redelegations-from.md)                   | Query all outgoing redelegatations from a validator                                           |
| [redelegation](redelegation.md)                               | Query a redelegation record based on delegator and a source and destination validator address |
| [redelegations](redelegations.md)                             | Query all redelegations records for one delegator                                             |
| [pool](pool.md)                                               | Query the current staking pool values                                                         |
| [parameters](parameters.md)                                   | Query the current staking parameters information                                              |
| [signing-info](signing-info.md)                               | Query a validator's signing information                                                       |
| [create-validator](create-validator.md)                       | Create new validator initialized with a self-delegation to it                                 |
| [edit-validator](edit-validator.md)                           | Edit and existing validator account                                                           |
| [delegate](delegate.md)                                       | Delegate liquid tokens to an validator                                                        |
| [unbond](unbond.md)                                           | Unbond shares from a validator                                                                |
| [redelegate](redelegate.md)                                   | Redelegate illiquid tokens from one validator to another                                      |
| [unjail](unjail.md)                                           | Unjail validator previously jailed for downtime                                               |

