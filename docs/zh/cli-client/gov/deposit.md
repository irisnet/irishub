# iriscli gov deposit

## 描述
 
充值保证金以激活提议
 
## 使用方式
 
```
iriscli gov deposit [flags]
```

打印帮助信息:

```
iriscli gov deposit --help
```
## 标志

| 命令，缩写       | 类型   | 是否必须 | 默认值                | 描述                                                         |
| ---------------- | ------ | -------- | --------------------- | ------------------------------------------------------------ |
| -h, --help       |        | 否       |                       | 打印帮助                                                     |
| --chain-id       | String | 否       |                       | tendermint 节点网络ID                                        |
| --account-number | int    | 否       |                       | 账户数字用于签名通证发送                                     |
| --async          |        | 否       | True                  | 异步广播传输信息                                             |
| --dry-run        |        | 否       |                       | 忽略--gas 标志 ，执行仿真传输，但不广播。                    |
| --fee            | String | 是       |                       | 设置传输需要的手续费                                         |
| --from           | String | 是       |                       | 用于签名的私钥名称                                           |
| --gas            | String | 否       | 20000                 | 每笔交易设定的gas限额; 设置为“simulate”以自动计算所需气体    |
| --gas-adjustment | Float  | 否       | 1                     | 调整因子乘以传输模拟返回的估计值; 如果手动设置气体限制，则忽略该标志 |
| --generate-only  |        | 否       |                       | 创建一个未签名的传输并写到标准输出中。                       |
| --indent         |        | 否       |                       | 在JSON响应中增加缩进                                         |
| --json           |        | 否       |                       | 以json格式返回输出                                           |
| --memo           | String | 否       |                       | 传输中的备注信息                                             |
| --print-response |        | 否       |                       | 返回传输响应 (仅仅当 async = false时有效)                    |
| --sequence       | Int    | 否       |                       | 等待签名传输的序列号。                                       |
| --ledger         | String | 否       |                       | 使用一个联网的分账设备                                       |
| --node           | String | 否       | tcp://localhost:26657 | <主机>:<端口> 链上的tendermint rpc 接口。                    |
| --trust-node     | String | 否       | True                  | 不验证响应的证明                                             |
| --commit         | String | 否     | True                  |是否等到交易有明确的返回值，如果是True，则忽略--async的内容|


## 特殊标志
 
| 名称, 速记        | 默认值                      | 描述                                                                                                                                                 | 是否必须  |
| ---------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --deposit        |                            | [string] 发起提议的保证金                                                                                                                         | Yes      |
| --proposal-id    |                            | [string] 充值保证金的提议ID                                                                                                        | Yes      |

## 例子

### 充值保证金

```shell
iriscli gov deposit --chain-id=test --proposal-id=1 --deposit=50iris --from=node0 --fee=0.01iris
```

输入正确的密码后，你就充值了50个iris用以激活提议的投票状态。

```txt
Committed at block 7 (tx hash: C1156A7D383492AE5C2EB1BADE0080C3A36BE8AED491DC5B2331056BED5D60DC, response:
 {
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 7944,
   "codespace": "",
   "tags": {
     "action": "deposit",
     "depositor": "faa1x25y3ltr4jvp89upymegvfx7n0uduz5kmh5xuz",
     "proposal-id": "1",
     "voting-period-start": "1"
   }
 })
```

如何查询保证金充值明细？

请点击下述链接：

[query-deposit](query-deposit.md)

[query-deposits](query-deposits.md)
