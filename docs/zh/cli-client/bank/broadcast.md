# iriscli bank broadcast

## 描述

在使用 [sign](./sign.md)生成签名的离线交易后，可以用此命令将交易广播到全网.

## Usage:

```
iriscli bank broadcast [tx-file] [flags] 
```

 

## 标志

| 命令，缩写   | 类型   | 是否必须 | 默认值                | 描述                                      |
| ------------ | ------ | -------- | --------------------- | ----------------------------------------- |
| -h, --help   |        | 否       |                       | 打印帮助信息                              |
| --chain-id   | String | 否       |                       | tendermint 节点网络ID                     |
| --height     | Int    | 否       |                       | 查询的区块高度用于获取最新的区块。        |
| --ledger     | String | 否       |                       | 使用一个联网的分账设备                    |
| --node       | String | 否       | tcp://localhost:26657 | <主机>:<端口> 链上的tendermint rpc 接口。 |
| --trust-node | String | 否       | True                  | 不验证响应的证明                          |



## 全局标志

| 命令，缩写             | 默认值         | 描述                                | 是否必须 |
| --------------------- | -------------- | ----------------------------------- | -------- |
| -e, --encoding string | hex            | 字符串二进制编码 (hex \|b64 \|btc ) | 否       |
| --home string         | $HOME/.iriscli | 配置和数据存储目录                  | 否       |
| -o, --output string   | text           | 输出格式 (text \|json)              | 否       |
| --trace               |                | 出错时打印完整栈信息                | 否       |



## 例子

### 广播离线交易

```
iriscli bank broadcast sign.json --chain-id=irishub-stage 

```

然后可以看到如下信息：

```
Committed at block 2265 (tx hash: A60224C8433487D48C8B03B51CB7A2BCB014932A97A55D946E5F30E561E1195E, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:4690 Tags:[{Key:[115 101 110 100 101 114] Value:[102 97 97 49 57 97 97 109 106 120 51 120 115 122 122 120 103 113 104 114 104 48 121 113 100 52 104 107 117 114 107 101 97 55 102 54 100 52 50 57 121 120] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[114 101 99 105 112 105 101 110 116] Value:[102 97 97 49 57 97 97 109 106 120 51 120 115 122 122 120 103 113 104 114 104 48 121 113 100 52 104 107 117 114 107 101 97 55 102 54 100 52 50 57 121 120] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 57 51 56 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "completeConsumedTxFee-iris-atto": "\"93800000000000\"",
     "recipient": "faa19aamjx3xszzxgqhrh0yqd4hkurkea7f6d429yx",
     "sender": "faa19aamjx3xszzxgqhrh0yqd4hkurkea7f6d429yx"
   }
 }
```