# iriscli gov

## Description

IRIShub governance module provides the basic functions as described below:
1. On-chain governance proposals on text
2. On-chain governance proposals on parameter change
3. On-chain governance proposals on software upgrade

## Usage

```shell
iriscli gov [command]
```

## Available Commands

| Name                                  | Description                                                     |
| ------------------------------------- | --------------------------------------------------------------- |
| [query-proposal](query-proposal.md)   | Query details of a single proposal                              |
| [query-proposals](query-proposals.md) | query proposals with optional filters                           |
| [query-vote](query-vote.md)           | query vote                                                      |
| [query-votes](query-votes.md)         | query votes on a proposal                                       |
| [query-deposit](query-deposit.md)     | Query details of a deposit                                      |
| [query-deposits](query-deposits.md)   | Query deposits on a proposal                                    |
| [query-tally](query-tally.md)         | Get the tally of a proposal vote                                |
| [query-params](query-params.md)       | query parameter proposal's config                               |
| [pull-params](pull-params.md)         | generate param.json file                                        |
| [submit-proposal](submit-proposal.md) | Create a new key, or import from seed                           |
| [deposit](deposit.md)                 | deposit tokens for activing proposal                            |
| [vote](vote.md)                       | vote for an active proposal, options: Yes/No/NoWithVeto/Abstain |

## Flags

| Name, shorthand | Default | Description   | Required |
| --------------- | ------- | ------------- | -------- |
| --help, -h      |         | help for gov  |          |

## Global Flags
| Name, shorthand | Default        | Description                            | Required |
| --------------- | -------------- | -------------------------------------- | -------- |
| --encoding, -e  | hex            | [string] Binary encoding (hex|b64|btc) |          |
| --home          | $HOME/.iriscli | [string] directory for config and data |          |
| --output, -o    | text           | [string] Output format (text|json)     |          |
| --trace         |                | print out full stack trace on errors   |          |

## Extended description

1. Any users can deposit some tokens to initiate a proposal. Once deposit reaches a certain value `min_deposit`, enter voting period, otherwise it will remain in the deposit period. Others can deposit the proposals on the deposit period. Once the sum of the deposit reaches `min_deposit`, enter voting period. However, if the block-time exceeds `max_deposit_period` in the deposit period, the proposal will be closed.
2. The proposals which enter voting period only can be voted by validators and delegators. The vote of a delegator who hasn't vote will be the same as his validator's vote, and the vote of a delegator who has voted will be remained. The votes wil be tallyed when reach `voting_period'.
3. More details about voting for proposals:
[CosmosSDK-Gov-spec](https://github.com/cosmos/cosmos-sdk/blob/develop/docs/spec/governance/overview.md)