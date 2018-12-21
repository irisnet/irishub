# iriscli record submit

## 描述

向链上提交一个存证

## 用法

```
iriscli record submit [flags]
```

## 标志

| 名称, 速记        | 默认值                     | 描述                                                                               | 必需      |
| ---------------  | -------------------------- | ---------------------------------------------------------------------------------- | -------- |
| --account-number |                            | [int] 用来签名交易的AccountNumber                                                   |          |
| --async          |                            | 异步广播交易                                                                        |          |
| --chain-id       |                            | [string] tendermint节点的链ID                                                       | 是       |
| --description    | description                | [string] 上传文件的描述信息                                                          |          |
| --dry-run        |                            | 忽略--gas标志并进行本地的交易仿真                                                     |          |
| --fee            |                            | [string] 支付的交易费用                                                              | 是       |
| --from           |                            | [string] 用来签名的私钥名                                                            | 是       |
| --from-addr      |                            | [string] 指定generate-only模式下的from address                                      |          |
| --gas string     | 200000                     | 单次交易的gas上限; 设置为"simulate"将自动计算相应的阈值                                |          |
| --gas-adjustment | 1                          | [float] 这个调整因子将乘以交易仿真返回的估计值; 如果手动设置了gas的上限，这个标志将被忽略 |          |
| --generate-only  |                            | 构建一个未签名交易并将其打印到标准输出                                                 |          |
| -h, --help       |                            | 提交命令帮助                                                                         |          |
| --indent         |                            | 在JSON响应中添加缩进                                                                 |          |
| --json           |                            | 输出将以json格式返回                                                                 |          |
| --ledger         |                            | 使用连接的硬件记账设备                                                                |          |
| --memo           |                            | [string] 发送交易的备忘录                                                            |          |
| --node           | tcp://localhost:26657      | [string] [string] tendermint节点开启的远程过程调用接口\<主机>:\<端口>                  |          |
| --onchain-data   |                            | [string] 上链数据                                                                    | 是      |
| --print-response |                            | 返回交易响应 (当且仅当同步模式下使用)                                                 |          |
| --sequence       |                            | [int] 用来签名交易的sequence number                                                  |          |
| --trust-node     | true                       | 关闭响应结果校验                                                                     |          |

## 例子

### 提交存证

```shell
iriscli record submit --chain-id="test" --onchain-data="this is my on chain data" --from=node0 --fee=0.1iris
```

运行成功以后，返回的结果如下：

```txt
Committed at block 486 (tx hash: 8AB91BF0E61AD2C860402B88579EE83167506E7C3A8597E873976915D82D4F1B, response:
 {
   "code": 0,
   "data": "cmVjb3JkOmFiNTYwMmJhYzEzZjExNzM3ZTg3OThkZDU3ODY5YzQ2ODE5NGVmYWQyZGIzNzYyNTc5NWYxZWZkOGQ5ZDYzYzY=",
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 3764,
   "codespace": "",
   "tags": {
     "action": "submit_record",
     "ownerAddress": "faa1x25y3ltr4jvp89upymegvfx7n0uduz5kmh5xuz",
     "record-id": "record:ab5602bac13f11737e8798dd57869c468194efad2db37625795f1efd8d9d63c6"
   }
 })
```

本次存证操作的record-id如下:

```txt
"record-id": "record:ab5602bac13f11737e8798dd57869c468194efad2db37625795f1efd8d9d63c6"
```

请务必备份record-id，以备将来查询本次存证。若丢失record-id，本次存证再也无法查询到。
