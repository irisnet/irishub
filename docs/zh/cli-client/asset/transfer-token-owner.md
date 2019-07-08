# iriscli asset transfer-token-owner

## 描述

转移Token的所有权

## 使用方式

```bash
iriscli asset transfer-token-owner <token-id> [flags]
```

## 标志

| 命令, 速记     | 类型   | 是否必须 | 默认值   | 描述                                                       |
| --------------------| -----  | -------- | --------|-------------------------------------------------------- |
| --to           | string | 是 | "" | 资产的新 Owner |

## 示例

```bash
iriscli asset transfer-token-owner kitty --to=<new-owner-address> --from=<key-name> --chain-id=irishub --fee=0.4iris --commit
```
