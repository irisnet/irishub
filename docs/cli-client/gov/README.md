# iriscli gov

## Description

This module provides the basic functions as described below:
1. On-chain governance proposals on parameter change
2. On-chain governance proposals on software upgrade 
3. On-chain governance proposals on software halt
4. On-chain governance proposals on tax usage

## Usage

```shell
iriscli gov <command>
```

Print all supported subcommands and flags:
```
iriscli gov --help
```

## Available Commands

| Name                                  | Description                                                     |
| ------------------------------------- | --------------------------------------------------------------- |
| [query-proposal](query-proposal.md)   | Query details of a single proposal                              |
| [query-proposals](query-proposals.md) | Query proposals with optional filters                           |
| [query-vote](query-vote.md)           | Query vote                                                      |
| [query-votes](query-votes.md)         | Query votes on a proposal                                       |
| [query-deposit](query-deposit.md)     | Query details of a deposit                                      |
| [query-deposits](query-deposits.md)   | Query deposits on a proposal                                    |
| [query-tally](query-tally.md)         | Get the tally of a proposal vote                                |
| [query-params](query-params.md)       | Query parameter proposal's config                               |
| [submit-proposal](submit-proposal.md) | Submit a proposal along with an initial deposit                          |
| [deposit](deposit.md)                 | Deposit tokens for active proposal                            |
| [vote](vote.md)                       | vote for an active proposal, options: Yes/No/NoWithVeto/Abstain |


## Extended description

1. Any user can deposit some tokens to submit a proposal. Once the deposit reaches a certain value `min_deposit`, the proposal enter voting period, otherwise it will remain in the deposit period. Others can deposit the proposals in the deposit period. Once the sum of the deposit reaches `min_deposit`, the proposal enter voting period. However, if the block-time exceeds `max_deposit_period` in the deposit period, the proposal will be closed.
2. The proposals which enter voting period only can be voted by validators. 
3. More details about voting for proposals: [Governance](../../features/governance.md)