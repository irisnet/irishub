# iriscli stake

## 介绍

Stake模块提供了一系列查询staking状态和发送staking交易的命令

## 用法

```
iriscli stake [subcommand]
```

打印子命令和flags：
```
iriscli stake --help
```

## 子命令

| 子命令                                                          | 功能                                                          |
| ------------------------------------------------------------- | --------------------------------------------------------------|
| [validator](validator.md)                                     | 查询某个验证者 |
| [validators](validators.md)                                   | 查询所有的验证者 |
| [delegation](delegation.md)                                   | 基于委托者地址和验证者地址的委托查询 |
| [delegations](delegations.md)                                 | 基于委托者地址的所有委托查询 |
| [delegations-to](delegations-to.md)                           | 查询在某个验证人上的所有委托      |
| [unbonding-delegation](unbonding-delegation.md)               | 基于委托者地址和验证者地址的unbonding-delegation记录查询 |
| [unbonding-delegations](unbonding-delegations.md)             | 基于委托者地址的所有unbonding-delegation记录查询 |
| [unbonding-delegations-from](unbonding-delegations-from.md)   | 基于验证者地址的所有unbonding-delegation记录查询|
| [redelegations-from](redelegations-from.md)                   | 基于某一验证者的所有重新委托查询  |
| [redelegation](redelegation.md)                               | 基于委托者地址，原源验证者地址和目标验证者地址的重新委托记录查询 |
| [redelegations](redelegations.md)                             | 基于委托者地址的所有重新委托记录查询 |
| [pool](pool.md)                                               | 查询最新的权益池 |
| [parameters](parameters.md)                                   | 查询最新的权益参数信息 |
| [signing-info](signing-info.md)                               | 查询验证者签名信息 |
| [create-validator](create-validator.md)                       | 以自委托的方式创建一个新的验证者 |
| [edit-validator](edit-validator.md)                           | 编辑已存在的验证者信息 |
| [delegate](delegate.md)                                       | 委托一定代币到某个验证者 |
| [unbond](unbond.md)                                           | 从指定的验证者解绑一定的股份 |
| [redelegate](redelegate.md)                                   | 重新委托一定的token从一个验证者到另一个验证者 |
| [unjail](unjail.md)                                           | 恢复之前由于宕机被惩罚的验证者的身份 |

