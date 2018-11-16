# iriscli upgrade info

## 描述

查询软件版本信息和升级模块的信息

## 用法

```
iriscli upgrade info
```

## 标志

| 名称, 速记       | 默认值                     | 描述                                                        | 必需     |
| --------------- | -------------------------- | ----------------------------------------------------------------- | -------- |
| --proposal-id      |        | 软件升级提议的ID                              | 是     |
| --voter     |                            | 签名switch消息的地址                             | 是      |
| --chain-id      |                            | [string] tendermint节点的链ID                               | 是       |
| --height        | 最近可证明区块高度           | [int] 查询的区块高度                                              |          |
| --help, -h      |                            | 查询命令帮助                                                |          |
| --indent        |                            | 在JSON格式的应答中添加缩进                                   |          |
| --ledger        |                            | 使用连接的硬件记账设备                                       |          |
| --node          | tcp://localhost:26657      | [string] tendermint节点开启的远程过程调用接口\<主机>:\<端口> |          |
| --trust-node    | true                       | 关闭响应结果校验                                            |          |

## 例子

查询当前版本信息

```
iriscli upgrade info 
```

然后它会打印如下内容：

```
{
  "current_proposal_id": "0",
  "current_proposal_accept_height": "-1",
  "version": {
    "Id": "0",
    "ProposalID": "0",
    "Start": "0",
    "ModuleList": [
      {
        "Start": "0",
        "End": "9223372036854775807",
        "Handler": "bank",
        "Store": [
          "acc"
        ]
      },
      {
        "Start": "0",
        "End": "9223372036854775807",
        "Handler": "stake",
        "Store": [
          "stake",
          "acc",
          "mint",
          "distr"
        ]
      },
      .......
    ]
  }
}
```