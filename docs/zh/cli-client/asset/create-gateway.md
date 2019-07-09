# iriscli asset create-gateway

## 描述

创建一个网关。网关用于映射外部资产。

## 使用方式

```bash
iriscli asset create-gateway [flags]
```

## 特有的标志

| 命令, 速记     | 类型   | 是否必须 | 默认值  | 描述                                                         |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --moniker           | string | 是    | ""       | 唯一的网关名字, 长度为3-8个英文字符 |
| --identity          | string | 否    | ""       | 可选的身份 (例如UPort or Keybase), 最大128个字符 |
| --details           | string | 否    | ""       | 可选的描述, 最大280个字符|
| --website           | string | 否    | ""       | 可选的外部网址, 最大128个字符|

## 示例

```bash
iriscli asset create-gateway --moniker=cats --identity=<pgp-id> --details="Cat Tokens" --website="www.example.com" --from=<key-name> --chain-id=irishub --fee=0.3iris
```
