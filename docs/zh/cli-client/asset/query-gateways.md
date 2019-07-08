# iriscli asset query-gateways

## 描述

查询所有网关信息，支持可选的owner参数

## 使用方式

```bash
iriscli asset query-gateways [flags]
```

## 特定的标志

| 命令, 速记     | 类型   | 是否必须 | 默认值  | 描述                                                         |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --owner           | Address | 否     |        | 要查询的目标所有者地址 |

## 示例

```bash
iriscli asset query-gateways --owner=<owner-address>
```
