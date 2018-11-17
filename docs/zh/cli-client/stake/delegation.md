# iriscli stake delegation

## 描述

基于委托者和验证者地址查询委托交易

## 用法

```
iriscli stake delegation [flags]
```

## 标志

| 名称, 速记             | 默认值                      | 描述                                                                 | 必需     |
| --------------------- | -------------------------- | -------------------------------------------------------------------- | -------- |
| --address-delegator   |                            | [string] 委托者bech地址                                               | Yes      |
| --address-validator   |                            | [string] 验证者bech地址                                               | Yes      |
| --chain-id            |                            | [string] Tendermint节点的链ID                                         |          |
| --height              | 最新的可证明区块高度         | 查询的区块高度                                                        |          |
| --help, -h            |                            | delegation命令帮助                                                    |          |
| --indent              |                            | 在JSON响应中添加缩进                                                   |          |
| --ledger              |                            | 使用连接的硬件记账设备                                                 |          |
| --node                | tcp://localhost:26657      | [string] Tendermint远程过程调用的接口\<主机>:\<端口>                    |          |
| --trust-node          | true                       | 关闭响应结果校验                                                       |          |

## 例子

### 查询验证者

```shell
iriscli stake delegation --address-validator=ValidatorAddress --address-delegator=DelegatorAddress

```

运行成功以后，返回的结果如下：

```txt
Delegation
Delegator: faa13lcwnxpyn2ea3skzmek64vvnp97jsk8qmhl6vx
Validator: fva15grv3xg3ekxh9xrf79zd0w077krgv5xf6d6thd
Shares: 0.2000000000Height: 290
```
