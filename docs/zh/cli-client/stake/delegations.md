# iriscli stake delegations

## 描述

查询某个委托者发起的所有委托记录

## 用法

```
iriscli stake delegations [delegator-addr] [flags]
```

## 标志

| 名称, 速记             | 默认值                     | 描述                                                                  | 必需     |
| --------------------- | -------------------------- | -------------------------------------------------------------------- | -------- |
| --chain-id            |                            | [string] Tendermint节点的链ID                                         |          |
| --height              | 最新的可证明区块高度         | 查询的区块高度                                                        |          |
| --help, -h            |                            | delegations命令帮助                                                  |          |
| --indent              |                            | 在JSON响应中添加缩进                                                  |          |
| --ledger              |                            | 使用连接的硬件记账设备                                                 |          |
| --node                | tcp://localhost:26657      | [string] Tendermint远程过程调用的接口\<主机>:\<端口>                   |          |
| --trust-node          | true                       | 关闭响应结果校验                                                      |          |

## Examples

### 查询某个委托者发起的所有委托记录

```shell
iriscli stake delegations DelegatorAddress
```

运行成功以后，返回的结果如下：

```json
[
  {
    "delegator_addr": "faa13lcwnxpyn2ea3skzmek64vvnp97jsk8qmhl6vx",
    "validator_addr": "fva15grv3xg3ekxh9xrf79zd0w077krgv5xf6d6thd",
    "shares": "0.2000000000",
    "height": "290"
  }
]
```
