      # iriscli stake redelegate

## 介绍

把某个委托的一部分或者全部从一个验证人转移到另外一个验证人

## 用法

```
iriscli stake redelegate [flags]
```

打印帮助信息
```
iriscli stake redelegate --help
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


## 特有的flags

| 名称                       | 类型   | 是否必填 | 默认值   | 功能描述         |
| -------------------------- | -----  | -------- | -------- | ------------------------------------------------------------------- |
| --address-validator-dest   | string | true     | ""       | 目标验证人地址 |
| --address-validator-source | string | true     | ""       | 源验证人地址 |
| --shares-amount            | float  | false    | 0.0      | 转移的share数量，正数 |
| --shares-percent           | float  | false    | 0.0      | 转移的比率，0到1之间的正数 |

用户可以用`--shares-amount`或者`--shares-percent`指定装委托的token数量。记住，这两个参数不可同时使用。

## 示例

执行如下在chain-id为test的网络上重新委托10iris的命令：

```
iriscli stake redelegate --chain-id=test-irishub --from=fuxi --fee=0.004iris --address-validator-source=fva106nhdckyf996q69v3qdxwe6y7408pvyvfcwqmd --address-validator-dest=fva1xpqw0kq0ktt3we5gq43vjphh7xcjfy6sfqamll  --shares-percent=0.1
```
输出信息：
```$xslt
Committed at block 1089 (tx hash: D9A60074B1E15ED33D1C0591AF7B6AD967B13409D342980DC0C858F811021C41, response:
 {
   "code": 0,
   "data": "DAiX2MzgBRCAtK+tAQ==",
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 27011,
   "codespace": "",
   "tags": {
     "action": "begin_redelegate",
     "delegator": "faa106nhdckyf996q69v3qdxwe6y7408pvyvufy0x2",
     "destination-validator": "fva1xpqw0kq0ktt3we5gq43vjphh7xcjfy6sfqamll",
     "end-time": "\u000c\u0008\ufffd\ufffd\ufffd\ufffd\u0005\u0010\ufffd\ufffd\ufffd\ufffd\u0001",
     "source-validator": "fva106nhdckyf996q69v3qdxwe6y7408pvyvfcwqmd"
   }
 })
```