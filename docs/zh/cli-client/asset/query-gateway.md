# iriscli asset query-gateway

## 描述

使用网关名称查询网关信息。

## 使用方式

```bash
iriscli asset query-gateway [flags]
```

## 特有的标志

| 命令, 速记     | 类型   | 是否必须 | 默认值  | 描述                                                         |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --moniker           | string | 是     | ""       | 全局唯一的网关名称, 长度为3-8个英文字符 |

## 示例

```bash
iriscli asset query-gateway --moniker cats
```
