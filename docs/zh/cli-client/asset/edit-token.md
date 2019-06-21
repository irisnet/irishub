# iriscli asset edit-token

## 描述

编辑指定ID的资产信息

## 使用方式

```
iriscli asset edit-token <token-id> [flags]
```

打印帮助信息:
```
iriscli asset edit-token <token-id> --help
```

## 特有的标志

| 命令, 速记     | 类型   | 是否必须 | 默认值  | 描述                                                         |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --name           | string | 是    |        | 资产的新名字，限定为32位的Unicode字符，如 "IRISnet"|
| --symbol         | string | 是    |        | 资产设置的新的唯一标志，限定为3-8位的中英文字母、数字和_的组合 |


## 示例

```
iriscli asset edit-token kitty  --name=kittyToken --symbol=cat    --fee=1iris
```

输出信息:
```txt
Committed at block 14 (tx hash: D24319D696F16C42E4F0508B28889A3EC3CC371EE92786F1945BA97BA1F6223D, response:
```

```json
 {
      "code": 0,
      "data": null,
      "log": "Msg 0: ",
      "info": "",
      "gas_wanted": 50000,
      "gas_used": 3696,
      "codespace": "",
      "tags": [
        {
          "key": "action",
          "value": "edit_token"
        },
        {
          "key": "token",
          "value": "cat"
        }
      ]
    }
```