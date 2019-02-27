# iriscli stake unbonding-delegations-from

## 描述

基于验证者地址的所有unbonding-delegation记录查询

## 用法

```
iriscli stake unbonding-delegations-from [validator-address] [flags]
```
打印帮助信息
```
iriscli stake unbonding-delegations-from --help
```

## 示例

基于验证者地址的所有unbonding-delegation记录查询
```
iriscli stake unbonding-delegations [validator-address]
```

运行成功以后，返回的结果如下：

```json
[
  {
    "delegator_addr": "iaa13lcwnxpyn2ea3skzmek64vvnp97jsk8qmhl6vx",
    "validator_addr": "fva15grv3xg3ekxh9xrf79zd0w077krgv5xf6d6thd",
    "creation_height": "1310",
    "min_time": "2018-11-15T06:24:22.754703377Z",
    "initial_balance": {
      "denom": "iris-atto",
      "amount": "20000000000000000"
    },
    "balance": {
      "denom": "iris-atto",
      "amount": "20000000000000000"
    }
  }
]
```
