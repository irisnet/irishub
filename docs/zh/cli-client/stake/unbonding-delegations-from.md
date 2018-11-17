# iriscli stake unbonding-delegations-from

## 描述

基于验证者地址的所有unbonding-delegation记录查询

## 用法

```
iriscli stake unbonding-delegations-from [operator-addr] [flags]
```

## 标志

| 名称, 速记           | 默认值                     | 描述                                                                 | 必需     |
| ------------------- | -------------------------- | ------------------------------------------------------------------- | -------- |
| --chain-id          |                            | [string] Tendermint节点的链ID                                        |          |
| --height            | 最新的可证明区块高度         | 查询的区块高度                                                       |          |
| --help, -h          |                            | unbonding-delegations-from命令帮助                                   |          |
| --indent            |                            | 在JSON响应中添加缩进                                                  |          |
| --ledger            |                            | 使用连接的硬件记账设备                                                |          |
| --node              | tcp://localhost:26657      | [string] Tendermint远程过程调用的接口\<主机>:\<端口>                   |          |
| --trust-node        | true                       | 关闭响应结果校验                                                      |          |

## 例子

### 基于验证者地址的所有unbonding-delegation记录查询

```shell
iriscli stake unbonding-delegations ValidatorAddress
```

运行成功以后，返回的结果如下：

```json
[
  {
    "delegator_addr": "faa13lcwnxpyn2ea3skzmek64vvnp97jsk8qmhl6vx",
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
