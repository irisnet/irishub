# iriscli bank send

## 描述

发送通证到指定地址 

## 使用方式

```
iriscli bank send --to=<account address> --from <key name> --fee=0.004iris --chain-id=<chain-id> --amount=10iris
```

 

## 标志

| 命令，速记       | 类型   | 是否必须 | 默认值                | 描述                                                         |
| ---------------- | ------ | -------- | --------------------- | ------------------------------------------------------------ |
| -h, --help       |        | 否       |                       | 打印帮助                                                     |
| --chain-id       | String | 否       |                       | tendermint 节点网络ID                                        |
| --account-number | int    | 否       |                       | 账户数字用于签名通证发送                                     |
| --amount         | String | 是       |                       | 需要发送的通证数量，比如10iris                               |
| --async          |        | 否       | True                  | 异步广播传输信息                                             |
| --dry-run        |        | 否       |                       | 忽略--gas 标志 ，执行仿真传输，但不广播。                    |
| --fee            | String | 是       |                       | 设置传输需要的手续费                                         |
| --from           | String | 是       |                       | 用于签名的私钥名称                                           |
| --from-addr      | string | 否       |                       | 在generate-only模式下指定的源地址                            |
| --gas            | String | 否       | 20000                 | 每笔交易设定的gas限额; 设置为“simulate”以自动计算所需气体    |
| --gas-adjustment | Float  | 否       | 1                     | 调整因子乘以传输模拟返回的估计值; 如果手动设置气体限制，则忽略该标志 |
| --generate-only  |        | 否       |                       | 创建一个未签名的传输并写到标准输出中。                       |
| --indent         |        | 否       |                       | 在JSON响应中增加缩进                                         |
| --json           |        | 否       |                       | 以json格式返回输出                                           |
| --memo           | String | 否       |                       | 传输中的备注信息                                             |
| --print-response |        | 否       |                       | 返回传输响应 (仅仅当 async = false时有效)                    |
| --sequence       | Int    | 否       |                       | 等待签名传输的序列号。                                       |
| --to             | String | 是       |                       | Bech32 编码的接收通证的地址。                                |
| --ledger         | String | 否       |                       | 使用一个联网的分账设备                                       |
| --node           | String | 否       | tcp://localhost:26657 | <主机>:<端口> 链上的tendermint rpc 接口。                    |
| --trust-node     | String | 否       | True                  | 不验证响应的证明                                             |



## 全局标志

| 命令，速记            | 默认值         | 描述                                | 是否必须 |
| --------------------- | -------------- | ----------------------------------- | -------- |
| -e, --encoding string | hex            | 字符串二进制编码 (hex \|b64 \|btc ) | 否       |
| --home string         | /root/.iriscli | 配置和数据存储目录                  | 否       |
| -o, --output string   | text           | 输出格式 (text \|json)              | 否       |
| --trace               |                | 出错时打印完整栈信息                | 否       |



## 

## 例子

### 发送通证到指定地址 

```
 iriscli bank send --to=faa19aamjx3xszzxgqhrh0yqd4hkurkea7f6d429yx  --from=test  --fee=0.004iris --chain-id=irishub-test --amount=10iris
```

命令执行完成后，返回执行的细节信息

```
[root@ce7da33d46c3 iriscli]# iriscli bank send --to=faa19aamjx3xszzxgqhrh0yqd4hkurkea7f6d429yx  --from=test  --fee=0.004iris --chain-id=irishub-test --amount=10iris
Password to sign with 'test':
Committed at block 2265 (tx hash: A60224C8433487D48C8B03B51CB7A2BCB014932A97A55D946E5F30E561E1195E, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:4690 Tags:[{Key:[115 101 110 100 101 114] Value:[102 97 97 49 57 97 97 109 106 120 51 120 115 122 122 120 103 113 104 114 104 48 121 113 100 52 104 107 117 114 107 101 97 55 102 54 100 52 50 57 121 120] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[114 101 99 105 112 105 101 110 116] Value:[102 97 97 49 57 97 97 109 106 120 51 120 115 122 122 120 103 113 104 114 104 48 121 113 100 52 104 107 117 114 107 101 97 55 102 54 100 52 50 57 121 120] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 57 51 56 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "completeConsumedTxFee-iris-atto": "\"93800000000000\"",
     "recipient": "faa19aamjx3xszzxgqhrh0yqd4hkurkea7f6d429yx",
     "sender": "faa19aamjx3xszzxgqhrh0yqd4hkurkea7f6d429yx"
   }
 }

```
