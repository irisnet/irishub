# 介绍

irislcd是一个能连接到任何全节点并且提供rest API接口的服务。通过这些API接口，用户可以发起交易或者查询区块链上的数据。irislcd可以验证全节点返回数据的完整性和有效性，因此irislcd可以以最小的带宽资源，最小的计算需求和最小的存储消耗来获取可全节点一样的数据安全。另外，irislcd还提供swagger-ui页面，这个页面列出来所以得API接口并且附有详细的说明文档。

## Irislcd options

irislcd有如下这些参数可以配置

| 参数名           | 类型      | 默认值                   | 是否必须 | 功能介绍                                          |
| --------------- | --------- | ----------------------- | -------- | ---------------------------------------------------- |
| chain-id        | 字符串    | null                    | true     | 所连接的tendermint区块链网络的ID |
| home            | 字符串    | "$HOME/.irislcd"        | false    | irislcd的home目录，用来存储秘钥和历史验证信息 |
| node            | 字符串    | "tcp://localhost:26657" | false    | 所要连接的全节点的url |
| laddr           | 字符串    | "tcp://localhost:1317"  | false    | 侦听的网络端口 |
| trust-node      | 布尔型    | false                   | false    | 是否信任全节点 |
| max-open        | 整型      | 1000                    | false    | 最大支持的连接数 |
| cors            | 字符串    | ""                      | false    | 是否支持跨域请求 |

## Start Irislcd

启动irislcd:
```
irislcd start --chain-id=<chain-id>
```
在浏览器中访问以下的url就可以打开Swagger-UI:
```
http://localhost:1317/swagger-ui/
```
打印版本号.
```
irislcd version
```
如果所连接的全节点是可信的，可以加上`trust-node`:
```
irislcd start --chain-id=<chain-id> --trust-node
```
