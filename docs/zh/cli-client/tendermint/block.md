# iriscli tendermint block

## 描述

在给定高度获取区块的验证数据。如果未指定高度，则将使用最新高度作为默认高度。

## 用法

```
  iriscli tendermint block <height> <flags>
```
或者
```
  iriscli tendermint block [flags]
```

## 标志

| 名称, 速记 | 默认值                    | 描述                                                             | 必需      |
| --------------- | -------------------------- | --------------------------------------------------------- | -------- |
| --chain-id    | 无 | tendermint节点的Chain ID   | 是       |
| --node string     |   tcp://localhost:26657                         | 要连接的节点  |                                     
| --help, -h      |           无| 	命令帮助|
| --trust-node    | true                       | 信任连接的完整节点，关闭响应结果校验                                            |          |

## 例子 

### 获得114263高度的区块信息

```shell
iriscli tendermint block 114263 --chain-id=<chain-id>
```
### 获得最新区块

```shell
iriscli tendermint block --chain-id=<chain-id> --trust-node
```