# iriscli record download

## 描述

下载存证ID指定的存证

## 用法

```
iriscli record download [record ID] [flags]
```

## 标志

| 名称, 速记 | 默认值                    | 描述                                                             | 必需      |
| --------------- | -------------------------- | --------------------------------------------------------- | -------- |
| --path            |     $HOME/.iriscli          | [string]存储下载的目录                             | 是     |
| --chain-id      |                            | [string] tendermint节点的链ID                              | 是       |
| --file-name     |                            | [string] 下载文件名                                        | 是       |
| --height        | 最近可证明区块高度           | [int] 查询的区块高度                                       |          |
| --help, -h      |                            | 下载命令帮助                                               |          |
| --indent        |                            | 在JSON响应中添加缩进                                       |          |
| --ledger        |                            | 使用连接的硬件记账设备                                      |          |
| --node          | tcp://localhost:26657      | [string] tendermint节点开启的远程过程调用接口\<主机>:\<端口> |          |
| --record-id     |                            | [string] 存证ID                                            | 是       |
| --trust-node    | true                       | 关闭响应结果校验                                            |          |

## 例子

### 下载存证

```shell
iriscli record download --chain-id=test --record-id=MyRecordID --file-name="download.txt"
```

iriscli会先在home目录(default: ~/.irislci)中创建一个用户指定的文件(download.txt)，然后把从链上下载的数据保存在此文件中。

```txt
[ONCHAIN] Downloading ~/.iriscli/download.txt from blockchain directly...
[ONCHAIN] Download file from blockchain complete.
```
