# iriscli gov

## 描述

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
| [query-proposal](query-proposal.md)   | 查询单个提议的详细信息                                             |
| [query-proposals](query-proposals.md) | 通过可选过滤器查询提议                                             |
| [query-vote](query-vote.md)           | 查询投票信息                                                      |
| [query-votes](query-votes.md)         | 查询提议的投票情况                                                 |
| [query-deposit](query-deposit.md)     | 查询保证金详情                                                    |
| [query-deposits](query-deposits.md)   | 查询提议的保证金                                                  |
| [query-tally](query-tally.md)         | 查询提议投票的统计                                                 |
| [query-params](query-params.md)       | 查询参数提议的配置                                                 |
| [pull-params](pull-params.md)         | 生成param.json文件                                               |
| [submit-proposal](submit-proposal.md) | 创建新密钥，或者通过助记词导入恢复                                   |
| [deposit](deposit.md)                 | 充值保证金代币以激活提议                                            |
| [vote](vote.md)                       | 为有效的提议投票，选项：Yes/No/NoWithVeto/Abstain                   |

## 标志

| 名称, 速记       | 默认值   | 描述          | 是否必须  |
| --------------- | ------- | ------------- | -------- |
| --help, -h      |         | help for gov  |          |

## 全局标志

| 名称, 速记       | 默认值          | 描述                                   | 是否必须  |
| --------------- | -------------- | -------------------------------------- | -------- |
| --encoding, -e  | hex            | [string] Binary encoding (hex &brvbar; b64 &brvbar; btc) |          |
| --home          | $HOME/.iriscli | [string] directory for config and data |          |
| --output, -o    | text           | [string] Output format (text &brvbar; json)     |          |
| --trace         |                | print out full stack trace on errors   |          |

## 补充描述

1.任何用户都可以存入一些令牌来发起提案。存款达到一定值min_deposit后，进入投票期，否则将保留存款期。其他人可以在存款期内存入提案。一旦存款总额达到min_deposit，输入投票期。但是，如果冻结时间超过存款期间的max_deposit_period，则提案将被关闭。
2.进入投票期的提案只能由验证人和委托人投票。未投票的代理人的投票将与其验证人的投票相同，并且投票的代理人的投票将保留。到达“voting_period”时，票数将被计算在内。
3.关于投票建议的更多细节：[CosmosSDK-Gov-spec](https://github.com/cosmos/cosmos-sdk/blob/develop/docs/spec/governance/overview.md)
