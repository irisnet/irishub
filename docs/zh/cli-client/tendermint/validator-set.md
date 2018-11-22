# iriscli tendermint validator-set

## 描述
根据指定高度在验证器上查询

## 用法

```
  iriscli tendermint validator-set [height] [flags]

```

## 标志

| 名称, 速记 | 默认值                    | 描述                                                             | 必需      |
| --------------- | -------------------------- | --------------------------------------------------------- | -------- |
| --chain-id    | 无 | [string] tendermint节点的链ID   | 是       |
| --node string     |   tcp://localhost:26657                         | 要连接的节点  |                                     
| --help, -h      |           无| 	下载命令帮助|
| --trust-node    | true                       | 信任连接的完整节点，关闭响应结果校验                                            |          |

## 例子 

```shell
iriscli tendermint validator-set 114360 --chain-id=fuxi-4000 
```
之后你会在验证器上查询到该高度的信息
