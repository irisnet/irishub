# iriscli stake unbonding-delegation

## 描述

基于委托者地址和验证者地址的unbonding-delegation记录查询

## 用法

```
iriscli stake unbonding-delegation [flags]
```
打印帮助信息
```
iriscli stake unbonding-delegation --help
```

## 特有的flags

| 名称, 速记           | 默认值                     | 描述                                                                 | 必需     |
| ------------------- | -------------------------- | ------------------------------------------------------------------- | -------- |
| --address-delegator |                            | [string] 委托者bech地址                                              | Yes      |
| --address-validator |                            | [string] 验证者bech地址                                             | Yes      |

## 示例

查询unbonding-delegation
```
iriscli stake unbonding-delegation --address-delegator=DelegatorAddress --address-validator=ValidatorAddress
```

运行成功以后，返回的结果如下：

```txt
Unbonding Delegation
Delegator: iaa13lcwnxpyn2ea3skzmek64vvnp97jsk8qrcezvm
Validator: fva15grv3xg3ekxh9xrf79zd0w077krgv5xf6d6thd
Creation height: 1310
Min time to unbond (unix): 2018-11-15 06:24:22.754703377 +0000 UTC
Expected balance: 0.02iris
```
