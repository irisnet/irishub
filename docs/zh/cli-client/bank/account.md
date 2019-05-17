# iriscli bank account

## 描述

查询选定账户信息

## 使用方式

```
iriscli bank account <address> <flags>
```


## 标志

| 命令，速记   | 类型   | 是否必须 | 默认值                | 描述                                      |
| ------------ | ------ | -------- | --------------------- | ----------------------------------------- |
| -h, --help   |        | 否       |                       | 打印帮助信息                              |
| --chain-id   | String | 否       |                       | tendermint 节点Chain ID                     |
| --height     | Int    | 否       |                       | 查询的区块高度用于获取最新的区块。        |
| --ledger     | String | 否       |                       | 使用ledger设备                    |
| --node       | String | 否       | tcp://localhost:26657 | <主机>:<端口> 链上的tendermint rpc 接口。 |
| --trust-node | String | 否       | True                  | 不验证响应的证明                          |



## 例子

### 查询账户信息 

```
 iriscli bank account <address>
```

执行完命令后，获得账户的详细信息如下
```
Account:
  Address:         iaa19aamjx3xszzxgqhrh0yqd4hkurkea7f646vaym
  Pubkey:          iap1addwnpepqwnsrt9m8tevhy4fdqyarunzuzzgz8e5q8jlceyf7uwpw0q0ptp2cp3lmjt
  Coins:           50iris
  Account Number:  0
  Sequence:        2
```

如果你查询一个错误的地址，将会返回如下信息
```
 iriscli bank account iaa19aamjx3xszzxgqhrh0yqd4hkurkea7f6d429zz
ERROR: decoding bech32 failed: checksum failed. Expected 46vaym, got d429zz.
```

如果查询一个在链上没有任何交易的地址，将会返回如下信息
```
iriscli bank account iaa1kenrwk5k4ng70e5s9zfsttxpnlesx5psh804vr
ERROR: {"codespace":"sdk","code":9,"message":"account iaa1kenrwk5k4ng70e5s9zfsttxpnlesx5psh804vr does not exist"}
```


## 扩展描述

查询IRISnet中的账户信息。

​    



​           
