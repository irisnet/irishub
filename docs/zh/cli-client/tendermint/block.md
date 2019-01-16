# iriscli tendermint block

## 描述

根据填写得区块高度去匹配此项

## 用法

```
  iriscli tendermint block [height] [flags]

```

## 标志

| 名称, 速记 | 默认值                    | 描述                                                             | 必需      |
| --------------- | -------------------------- | --------------------------------------------------------- | -------- |
| --chain-id    | 无 | [string] tendermint节点的链ID   | 是       |
| --node string     |   tcp://localhost:26657                         | 要连接的节点  |                                     
| --help, -h      |           无| 	下载命令帮助|
| --trust-node    | true                       | 信任连接的完整节点，关闭响应结果校验                                            |          |

## 例子 

### 获得114263高度的区块信息

```shell
iriscli tendermint block 114263  --chain-id=irishub-test
```
### 获得最新区块

```shell
iriscli tendermint block  --chain-id=irishub-test --trust-node=true

```