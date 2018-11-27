# iriscli upgrade query-switch

## 描述

查询switch信息来知道某人对某升级提议是否发送了switch消息。

## 用法

```
iriscli upgrade query-switch --proposalID <proposalID> --voter <voter address>
```
打印帮助信息:

```
iriscli upgrade query-switch  --help
```

## 标志

| 名称, 速记       | 默认值                     | 描述                                                        | 必需     |
| --------------- | -------------------------- | ----------------------------------------------------------------- | -------- |
| --proposalID      |        | 软件升级提议的ID                              | 是     |
| --voter     |                            | 签名switch消息的地址                             | 是      |

## 例子

查询用户`faa1qvt2r6hh9vyg3kh4tnwgx8wh0kpa7q2lsk03fe`是否对升级提议（ID为5）发送了switch消息

```
iriscli upgrade query-switch --proposalID=5 --voter=faa1qvt2r6hh9vyg3kh4tnwgx8wh0kpa7q2lsk03fe
```
