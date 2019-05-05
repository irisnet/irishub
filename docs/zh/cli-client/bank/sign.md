# iriscli bank sign

## 描述

签名生成的离线交易文件。该文件由`generate-only`标志生成。

## 使用方式

```
iriscli bank sign <file> <flags>
```

 

## 标志

| 命令，速记       | 类型   | 是否必须 | 默认值                | 描述                                                         |
| ---------------- | ------ | -------- | --------------------- | ------------------------------------------------------------ |
| -h, --help       |        | 否       |                       | 打印帮助                                                     |
| --append         | Bool  | 是       | True                  | 将签名附加到现有签名。 如果禁用，旧签名将被覆盖              |
| --name           | String | 是       |                       | 用于签名的私钥名称                                           |
| --offline        | Boole | 是       | False                 | 离线模式. 不查询本地缓存                                     |
| --print-sigs     | Bool  | 是       | False                 | 打印必须签署交易的地址和已签名的地址，然后退出               |


## 例子

### 签名一个离线的转账交易

首先你必须使用`iriscli bank send`命令和`generate-only`标志来生成一个转账交易，如下

```  
iriscli bank send --to=<address> --from=<key_name> --fee=0.3iris --chain-id=<chain-id> --amount=10iris --generate-only

{"type":"auth/StdTx","value":{"msg":[{"type":"cosmos-sdk/Send","value":{"inputs":[{"address":"iaa19aamjx3xszzxgqhrh0yqd4hkurkea7f646vaym","coins":[{"denom":"iris-atto","amount":"10000000000000000000"}]}],"outputs":[{"address":"iaa19aamjx3xszzxgqhrh0yqd4hkurkea7f646vaym","coins":[{"denom":"iris-atto","amount":"10000000000000000000"}]}]}}],"fee":{"amount":[{"denom":"iris-atto","amount":"4000000000000000"}],"gas":"200000"},"signatures":null,"memo":""}}
```

保存输出到文件中，如`tx.json`
```
iriscli bank send --to=<address> --from=<key_name> --fee=0.3iris --chain-id=<chain-id> --amount=10iris --generate-only > tx.json
```

接着来签名这个离线文件
```
iriscli bank sign tx.json --name=<key_name> 

{"type":"auth/StdTx","value":{"msg":[{"type":"cosmos-sdk/Send","value":{"inputs":[{"address":"iaa106nhdckyf996q69v3qdxwe6y7408pvyvyxzhxh","coins":[{"denom":"iris-atto","amount":"10000000000000000000"}]}],"outputs":[{"address":"iaa1893x4l2rdshytfzvfpduecpswz7qtpstevr742","coins":[{"denom":"iris-atto","amount":"10000000000000000000"}]}]}}],"fee":{"amount":[{"denom":"iris-atto","amount":"40000000000000000"}],"gas":"200000"},"signatures":[{"pub_key":{"type":"tendermint/PubKeySecp256k1","value":"Auouudrg0P86v2kq2lykdr97AJYGHyD6BJXAQtjR1gzd"},"signature":"sJewd6lKjma49rAiGVfdT+V0YYerKNx6ZksdumVCvuItqGm24bEN9msh7IJ12Sil1lYjqQjdAcjVCX/77FKlIQ==","account_number":"0","sequence":"3"}],"memo":"test"}}
```


随后得到签名详细信息，如下输出中你会看到签名信息。 

**sJewd6lKjma49rAiGVfdT+V0YYerKNx6ZksdumVCvuItqGm24bEN9msh7IJ12Sil1lYjqQjdAcjVCX/77FKlIQ==**

签完名后的交易可以在网络中被广播 [broadcast command](./broadcast.md)