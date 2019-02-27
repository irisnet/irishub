# iriscli bank account

## 描述

查询选定账户信息

## 使用方式

```
iriscli bank account [address] [flags] 
```

 

## 标志

| 命令，速记   | 类型   | 是否必须 | 默认值                | 描述                                      |
| ------------ | ------ | -------- | --------------------- | ----------------------------------------- |
| -h, --help   |        | 否       |                       | 打印帮助信息                              |
| --chain-id   | String | 否       |                       | tendermint 节点网络ID                     |
| --height     | Int    | 否       |                       | 查询的区块高度用于获取最新的区块。        |
| --ledger     | String | 否       |                       | 使用一个联网的分账设备                    |
| --node       | String | 否       | tcp://localhost:26657 | <主机>:<端口> 链上的tendermint rpc 接口。 |
| --trust-node | String | 否       | True                  | 不验证响应的证明                          |



## 全局标志

| 命令，速记            | 默认值         | 描述                                | 是否必须 |
| --------------------- | -------------- | ----------------------------------- | -------- |
| -e, --encoding string | hex            | 字符串二进制编码 (hex \|b64 \|btc ) | 否       |
| --home string         | $HOME/.iriscli | 配置和数据存储目录                  | 否       |
| -o, --output string   | text           | 输出格式 (text \|json)              | 否       |
| --trace               |                | 出错时打印完整栈信息                | 否       |



## 例子

### 查询账户信息 

```
 iriscli bank account iaa19aamjx3xszzxgqhrh0yqd4hkurkea7f646vaym
```

执行完命令后，获得账户的详细信息如下

```
{

  "address": "iaa19aamjx3xszzxgqhrh0yqd4hkurkea7f646vaym",

  "coins": [

    "50iris"

  ],

  "public_key": {

    "type": "tendermint/PubKeySecp256k1",

    "value": "AzlCwiA5Tvxwi7lMB/Hihfp2qnaks5Wrrgkg/Jy7sEkF"

  },

  "account_number": "0",

  "sequence": "1"

}



```
如果你查询一个错误的地址，将会返回如下信息
```
 iriscli bank account iaa19aamjx3xszzxgqhrh0yqd4hkurkea7f6d429zz
ERROR: decoding bech32 failed: checksum failed. Expected 46vaym, got d429zz.
```
如果查询一个空地址，，将会返回如下信息。
```
iriscli bank account iaa1kenrwk5k4ng70e5s9zfsttxpnlesx5psh804vr
ERROR: No account with address iaa1kenrwk5k4ng70e5s9zfsttxpnlesx5psh804vr was found in the state.
Are you sure there has been a transaction involving it?
```


## 扩展描述

查询iris网络中的账户信息。

​    



​           
