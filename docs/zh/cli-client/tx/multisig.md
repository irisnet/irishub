# iriscli tx multisign

## 描述

多重签名，多个用户对同一个交易进行数字签名,且达到最低签名人数后，该交易方可被正常广播。

## 使用方式

```
iriscli tx multisign <file> <key name> <[signature]...> <flags>
```


## 标志

| 命令，速记   | 类型 | 是否必须 | 默认值 | 描述    |
| ---------- | ---- | ------ | ----- | ------ |
| -h, --help |      |        |       | 打印帮助 |

## 例子

### 创建多签账户

首先你必须先创建一个多签账户, 请参阅 [add](../keys/add.md)

创建一个包含3个子账户的多签账户，且指定签名人数必须大于等于2人，该交易才能被正常广播。

```  
iriscli keys add <multi_account_keyname> --multisig-threshold=2 --multisig=<signer_keyname_1>,<signer_keyname_2>,<signer_keyname_3>...
```

::: tips
<signer_keyname> 可以为 local/offline/ledger 类型， 但不允许为multi类型。

其中， offline类型的账户可以在add时通过指定 --pubkey 生成。
:::

### 多签账户构造交易

以转账交易为例, 指定--generate-only，构造交易并生成Tx-generate.json文件：
```  
iriscli bank send --amount=1iris --fee=0.3iris --chain-id=<chain-id> --from=<multi_account_keyname> --to=<address> --generate-only > Tx-generate.json
```

### 多账户签名交易

先用 `iriscli keys show <multi_account_keyname>` 获取`<multi_account_address>`

由于指定的threshold=2， 则最少需要2人签名就能完成交易签名。分别对Tx-generate.json进行签名， 并生成对应的Tx-sign.json文件。

用signer_1对交易进行签名:
```  
iriscli tx sign Tx-generate.json --name=<signer_keyname_1> --chain-id=<chain-id> --multisig=<multi_account_address> --signature-only >Tx-sign-user_1.json
```

用signer_2对交易进行签名:
```  
iriscli tx sign Tx-generate.json --name=<signer_keyname_2> --chain-id=<chain-id> --multisig=<multi_account_address> --signature-only >Tx-sign-user_2.json
```

### 把多个签名合并，并生成签名后的交易

合并多个签名，并生成签名后的交易文件Tx-signed.json：

```  
iriscli tx multisign --chain-id=<chain-id> Tx-generate.json <multi_account_keyname> Tx-sign-user_1.json Tx-sign-user_2.json > Tx-signed.json
```


### 广播签名后交易

签完名后的交易可以在网络中被广播 [broadcast command](broadcast.md)

```  
iriscli tx broadcast Tx-signed.json --commit
```
