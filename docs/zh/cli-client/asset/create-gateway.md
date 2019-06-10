# iriscli asset create-gateway

## 描述

创建一个网关。网关用于映射外部资产。

## 使用方式

```
iriscli asset create-gateway [flags]
```

打印帮助信息:
```
iriscli asset create-gateway --help
```

## 特有的标志

| 命令, 速记     | 类型   | 是否必须 | 默认值  | 描述                                                         |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --moniker           | string | 是    | ""       | 唯一的网关名字, 长度为3-8个英文字符|
| --identity          | string | 否    | ""       | 可选的身份 (例如UPort or Keybase) |
| --details           | string | 否    | ""       | 可选的描述, 最大280个字符|
| --website           | string | 否    | ""       | 可选的外部网址, 最大128个字符|


## 示例

```
iriscli asset create-gateway --moniker=<name> --identity=<identity> --details=<description> 
--website=<website> --chain-id=<chain-id> --from=<key_name> --fee=0.3iris
```

输出信息:
```txt
Committed at block 306 (tx hash: 5A4C6E00F4F6BF795EB05D2D388CBA0E8A6E6CF17669314B1EE6A31729A22450, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:3398 Tags:[{Key:[97 99 116 105 111 110] Value:[115 101 114 118 105 99 101 45 119 105 116 104 100 114 97 119 45 102 101 101 115] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 54 55 57 54 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
```

```json
 {
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 16462,
   "codespace": "",
   "tags": {
     "moniker": "testgw"
   }
 })
```