# iriscli gov

## 描述

该模块提供如下所述的基本功能：
1. 参数修改提议的链上治理
2. 软件升级提议的链上治理
3. 网络终止提议的链上治理
4. Tax收入分配提议的链上治理

## 使用方式

```shell
iriscli gov <command>
```

打印子命令和参数
```
iriscli gov --help
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
| [submit-proposal](submit-proposal.md) | 提交提议并抵押初始保证金                                  |
| [deposit](deposit.md)                 | 充值保证金代币以激活提议                                            |
| [vote](vote.md)                       | 为有效的提议投票，选项：Yes/No/NoWithVeto/Abstain                   |

## 补充描述

1.任何用户都可以抵押一些代币来发起提案。抵押达到一定值(min_deposit)后，进入投票期(voting period)，否则将停留在抵押期(deposit period)。其他人可以在抵押期内充值保证金。一旦抵押总额达到min_deposit，进入投票期(voting period)。但是，如果冻结时间超过存款期间的max_deposit_period，则提案将被关闭。
2.进入投票期的提案只能由验证人投票。
3.关于投票建议的更多细节：[Governance](../../features/governance.md)
