# iriscli bank sign

## 描述

签名生成的离线传输文件。该文件由 --generate-only 标志生成。

## 使用方式

```
iriscli bank sign <file> [flags]
```

 

## 标志

| 命令，缩写      | 类型   | 是否必须 | 默认值                | 描述                                                         |
| ---------------- | ------ | -------- | --------------------- | ------------------------------------------------------------ |
| -h, --help       |        | 否       |                       | 打印帮助                                                     |
| --append         | Boole  | 是       | True                  | 将签名附加到现有签名。 如果禁用，旧签名将被覆盖              |
| --name           | String | 是       |                       | 用于签名的私钥名称                                           |
| --offline        | Boole  | 是       | False                 | 离线模式. 不查询本地缓存                                     |
| --print-sigs     | Boole  | 是       | False                 | 打印必须签署交易的地址和已签名的地址，然后退出               |
| --chain-id       | String | 否       |                       | tendermint 节点网络ID                                        |
| --account-number | Int    | 否       |                       | 账户数字用于签名通证发送                                     |
| --amount         | String | 是       |                       | 需要发送的通证数量，比如10iris                               |
| --async          |        | 否       | True                  | 异步广播传输信息                                             |
| --dry-run        |        | 否       |                       | 忽略--gas 标志 ，执行仿真传输，但不广播。                    |
| --fee            | String | 是       |                       | 设置传输需要的手续费                                         |
| --from           | String | 是       |                       | 用于签名的私钥名称                                           |
| --from-addr      | String | 否       |                       | 在generate-only模式下指定的源地址                            |
| --gas            | String | 否       | 20000                 | 每笔交易设定的gas限额; 设置为“simulate”以自动计算所需气体    |
| --gas-adjustment | Float  | 否       | 1                     | 调整因子乘以传输模拟返回的估计值; 如果手动设置气体限制，则忽略该标志 |
| --generate-only  |        | 否       |                       | 创建一个未签名的传输并写到标准输出中。                       |
| --indent         |        | 否       |                       | 在JSON响应中增加缩进                                         |
| --json           |        | 否       |                       | 以json格式返回输出                                           |
| --memo           | String | 否       |                       | 传输中的备注信息                                             |
| --print-response |        | 否       |                       | 返回传输响应 (仅仅当 async = false时有效)                    |
| --sequence       | Int    | 否       |                       | 等待签名传输的序列号。                                       |
| --to             | String | 否       |                       | Bech32 编码的接收通证的地址。                                |
| --ledger         | String | 否       |                       | 使用一个联网的分账设备                                       |
| --node           | String | 否       | tcp://localhost:26657 | <主机>:<端口> 链上的tendermint rpc 接口。                    |
| --trust-node     | String | 否       | True                  | 不验证响应的证明                                             |



## 全局标志

| 命令，缩写            | 默认值         | 描述                                | 是否必须 | 类型   |
| --------------------- | -------------- | ----------------------------------- | -------- | ------ |
| -e, --encoding string | hex            | 字符串二进制编码 (hex \|b64 \|btc ) | False    | String |
| --home string         | /root/.iriscli | 配置和数据存储目录                  | False    | String |
| -o, --output string   | text           | 输出格式(text \|json)               | False    | String |
| --trace               |                | 出错时打印完整栈信息                | False    |        |

## 例子

### 对一个离线发送文件签名

首先你必须使用 **iriscli bank send**  命令和标志 **--generate-only** 来生成一个发送记录，如下

```  
iriscli bank send --to=faa19aamjx3xszzxgqhrh0yqd4hkurkea7f6d429yx  --from=test  --fee=0.004iris --chain-id=irishub-test --amount=10iris --generate-only

{"type":"auth/StdTx","value":{"msg":[{"type":"cosmos-sdk/Send","value":{"inputs":[{"address":"faa19aamjx3xszzxgqhrh0yqd4hkurkea7f6d429yx","coins":[{"denom":"iris-atto","amount":"10000000000000000000"}]}],"outputs":[{"address":"faa19aamjx3xszzxgqhrh0yqd4hkurkea7f6d429yx","coins":[{"denom":"iris-atto","amount":"10000000000000000000"}]}]}}],"fee":{"amount":[{"denom":"iris-atto","amount":"4000000000000000"}],"gas":"200000"},"signatures":null,"memo":""}}
```



保存输出到文件中，如  /root/output/output/node0/test_send_10iris.txt.

接着来签名这个离线文件.

```
iriscli bank sign /root/output/output/node0/test_send_10iris.txt --name=test  --offline=false --print-sigs=false --append=true
```

随后得到签名详细信息，如下输出中你会看到签名信息。 

**ci+5QuYUVcsARBQWyPGDgmTKYu/SRj6TpCGvrC7AE3REMVdqFGFK3hzlgIphzOocGmOIa/wicXGlMK2G89tPJg==**

```
iriscli bank sign /root/output/output/node0/test_send_10iris.txt --name=test  --offline=false --print-sigs=false --append=true
Password to sign with 'test':
{"type":"auth/StdTx","value":{"msg":[{"type":"cosmos-sdk/Send","value":{"inputs":[{"address":"faa19aamjx3xszzxgqhrh0yqd4hkurkea7f6d429yx","coins":[{"denom":"iris-atto","amount":"10000000000000000000"}]}],"outputs":[{"address":"faa19aamjx3xszzxgqhrh0yqd4hkurkea7f6d429yx","coins":[{"denom":"iris-atto","amount":"10000000000000000000"}]}]}}],"fee":{"amount":[{"denom":"iris-atto","amount":"4000000000000000"}],"gas":"200000"},"signatures":[{"pub_key":{"type":"tendermint/PubKeySecp256k1","value":"AzlCwiA5Tvxwi7lMB/Hihfp2qnaks5Wrrgkg/Jy7sEkF"},"signature":"ci+5QuYUVcsARBQWyPGDgmTKYu/SRj6TpCGvrC7AE3REMVdqFGFK3hzlgIphzOocGmOIa/wicXGlMK2G89tPJg==","account_number":"0","sequence":"2"}],"memo":""}}
```

