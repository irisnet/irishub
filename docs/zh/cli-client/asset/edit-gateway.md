# iriscli asset edit-gateway

## 描述

编辑指定名字的网关信息

## 使用方式

```bash
iriscli asset edit-gateway [flags]
```

## 特有的标志

| 命令, 速记     | 类型   | 是否必须 | 默认值  | 描述                                                         |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --moniker           | string | 是    | ""       | 唯一的网关名字, 长度为3-8; 以英文字符开始，后接英文字符与数字|
| --identity          | string | 否    | ""       | 可选的身份 (例如UPort or Keybase), 最大128个字符 |
| --details           | string | 否    | ""       | 可选的描述, 最大280个字符|
| --website           | string | 否    | ""       | 可选的外部网址, 最大128个字符|

## 示例

```bash
iriscli asset edit-gateway --moniker=cats --identity=<pgp-id> --details="Cat Tokens" --website="http://www.example.com" --from=<key-name> --chain-id=irishub --fee=0.4iris --commit
```
