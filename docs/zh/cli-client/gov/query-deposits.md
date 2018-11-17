# iriscli gov query-deposits

## 描述

查询指定提议的保证金详细情况

## 使用方式

```
iriscli gov query-deposits [flags]
```

## 标志

| 名称, 速记       | 默认值                      | 描述                                                                                                                                                 | 是否必须  |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --chain-id      |                            | [string] tendermint节点的链ID                                                                                                                 | Yes      |
| --height        |                            | [int] 查询的区块高度                                                                                  |          |
| --help, -h      |                            | 查询命令帮助                                                                                                                              |          |
| --indent        |                            | 在JSON响应中添加缩进                                                                                                                          |          |
| --ledger        |                            | 使用连接的硬件记账设备                                                                                                                        |          |
| --node          | tcp://localhost:26657      | [string] tendermint节点开启的远程过程调用接口\<主机>:\<端口>                                                                                  |          |
| --proposal-id   |                            | [string] 提议ID                                                                                                        | Yes      |
| --trust-node    | true                       | 关闭响应结果校验                                                                                                                    |          |

## 例子

### 查询所有保证金

```shell
iriscli gov query-deposits --chain-id=test --proposal-id=1
```

你可以查询到指定提议的所有保证金代币，包括每个存款人的充值详情。

```txt
[
  {
    "depositor": "faa1c4kjt586r3t353ek9jtzwxum9x9fcgwetyca07",
    "proposal_id": "1",
    "amount": [
      {
        "denom": "iris-atto",
        "amount": "35000000000000000000"
      }
    ]
  }
]
```