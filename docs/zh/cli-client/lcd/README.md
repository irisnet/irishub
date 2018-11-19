# 介绍

IRISLCD节点是一个提供REST APIs接口的服务器，它可以连接到任何IRISHUB的全节点。通过这些REST APIs，用户可以发送交易和查询区块链数据。对于从链上返回的数据，IRISLCD可以独立验证数据的有效性。因此，它可以提供与全节点相同的安全性，同时对带宽，计算和存储资源的要求最低。此外，它还提供了swagger-ui，它详细描述了它提供的API以及如何使用它们。

## IRISLCD的用法

IRISLCD有两个子命令:

| 子命令      | 功能                 | 示例命令 |
| --------------- | --------------------------- | --------------- |
| version         | 打印版本信息   | irislcd version |
| start           | 启动一个IRISLCD节点  | irislcd start --chain-id=<chain-id> |

`start`子命令有如下参数可配置

| 参数名称        | 类型      | 默认值                 | 是否必填 | 功能描述                                          |
| --------------- | --------- | ----------------------- | -------- | ---------------------------------------------------- |
| chain-id        | string    | null                    | true     | Tendermint节点的chain ID |
| home            | string    | "$HOME/.irislcd"        | false    | 配置home目录，key和proof相关的信息都存于此 |
| node            | string    | "tcp://localhost:26657" | false    | 全节点的rpc地址 |
| laddr           | string    | "tcp://localhost:1317"  | false    | 侦听的地址和端口 |
| trust-node      | bool      | false                   | false    | 是否信任全节点 |
| max-open        | int       | 1000                    | false    | 最大连接数 |
| cors            | string    | ""                      | false    | 允许跨域访问的地址 |

## 示例命令

1. 如果所连接的全节点是可信的，那么可以使能信任模式；如果不可行，那么要关闭信任模式。默认是非信任模式，打开信任模式的方法是加上`--trust-node`：
```bash
irislcd start --chain-id=<chain-id> --trust-node
```

2. 如果需要在其他机器上访问此IRISLCD节点，还需要配置`--laddr`参数，例如：
```bash
irislcd start --chain-id=<chain-id> --laddr=tcp://0.0.0.0:1317
```