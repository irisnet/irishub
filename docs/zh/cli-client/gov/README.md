# iriscli gov

## 描述

IRIShub Gov模块允许您提交提案，投票和查询您关注的管理信息：
1.任何用户都可以存入一些令牌来发起提案。存款达到一定值min_deposit后，进入投票期，否则将保留存款期。其他人可以在存款期内存入提案。一旦存款总额达到min_deposit，输入投票期。但是，如果冻结时间超过存款期间的max_deposit_period，则提案将被关闭。
2.进入投票期的提案只能由验证人和委托人投票。未投票的代理人的投票将与其验证人的投票相同，并且投票的代理人的投票将保留。到达“voting_period”时，票数将被计算在内。
3.关于投票建议的更多细节：[CosmosSDK-Gov-spec](https://github.com/cosmos/cosmos-sdk/blob/develop/docs/spec/governance/overview.md)

该模块提供如下所述的基本功能：
1.关于案文的连锁治理提案
2.关于参数变化的链式治理建议
3.关于软件升级的链式治理建议

## 使用方式

```shell
iriscli gov [command]
```

## 相关命令

| 命令                                  | 描述                                                             |
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

## 标志

| 名称, 速记       | 默认值   | 描述          | 是否必须  |
| --------------- | ------- | ------------- | -------- |
| --help, -h      |         | help for gov  |          |

## 全局标志

| 名称, 速记       | 默认值          | 描述                                   | 是否必须 |
| --------------- | -------------- | -------------------------------------- | -------- |
| --encoding, -e  | hex            | [string] Binary encoding (hex|b64|btc) |          |
| --home          | $HOME/.iriscli | [string] directory for config and data |          |
| --output, -o    | text           | [string] Output format (text|json)     |          |
| --trace         |                | print out full stack trace on errors   |          |
