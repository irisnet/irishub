# iriscli bank burn

## 描述

通过减去账户余额的方式来销毁指定账户的一些token，任何人都可以使用这个交易。

## 使用方式

销毁10iris
```
iriscli bank burn --from=<key_name> --fee=0.3iris --chain-id=<chain-id> --amount=10iris
```


## 标志

| 命令，缩写       | 类型   | 是否必须 | 默认值                | 描述                                                         |
| ---------------- | ------ | -------- | --------------------- | ------------------------------------------------------------ |
| --amount         | String | 是       |                       | 销毁token数量，比如10iris                               |


## 例子

### 销毁token

```
 iriscli bank burn --from=<key_name> --fee=0.3iris --chain-id=<chain-id> --amount=10iris
```

命令执行完成后，返回执行的细节信息
```
[Committed at block 87 (tx hash: AEA8E49C1BC9A81CAFEE8ACA3D0D96DA7B5DC43B44C06BACEC7DCA2F9C4D89FC, response:
  {
    "code": 0,
    "data": null,
    "log": "Msg 0: ",
    "info": "",
    "gas_wanted": 200000,
    "gas_used": 3839,
    "codespace": "",
    "tags": {
      "action": "burn",
      "burnFrom": "iaa106nhdckyf996q69v3qdxwe6y7408pvyvyxzhxh"
      "burnAmount": "10iris"
    }
  })
```
