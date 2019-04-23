# iriscli bank send

## 描述

发送通证到指定地址, 此条命令同时包含“创建/签名/广播”交易。

## 使用方式

转账10iris
```
iriscli bank send --to=<address> --from=<key_name> --fee=0.3iris --chain-id=<chain-id> --amount=10iris
``` 

## 标志

| 命令，缩写       | 类型   | 是否必须 | 默认值                | 描述                                                         |
| ---------------- | ------ | -------- | --------------------- | ------------------------------------------------------------ |
| --amount         | String | 是       |                       | 需要发送的通证数量，比如10iris                               |
| --to             | String | 是       |                       | Bech32 编码的接收通证的地址。                                |

## 例子

### 发送通证到指定地址 

```
 iriscli bank send --to=<address> --from=<key_name> --fee=0.3iris --chain-id=<chain-id> --amount=10iris
```

命令执行完成后，返回执行的细节信息
```
Committed at block 87 (tx hash: AEA8E49C1BC9A81CAFEE8ACA3D0D96DA7B5DC43B44C06BACEC7DCA2F9C4D89FC, response:
 {
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 3839,
   "codespace": "",
   "tags": {
     "action": "send",
     "recipient": "iaa1893x4l2rdshytfzvfpduecpswz7qtpstevr742",
     "sender": "iaa106nhdckyf996q69v3qdxwe6y7408pvyvyxzhxh"
   }
 })
```
