# iriscli asset mint-token

## 描述

资产的所有者增发一定数量的代币到指定账户地址

## 使用方式

```bash
iriscli asset mint-token <token-id> [flags]
```

## 特有的标志

| 命令，速记 | 类型    | 是否必须  | 默认值 | 描述                   |
| -------- | ------ | -------- | ----- | ---------------------- |
| --to     | string |          |       | 接收增发的账户，默认为自己 |
| --amount | uint64 | 是       | 0     | 增发数量                |

## 示例

```bash
iriscli asset mint-token kitty --amount=1000000 --from=<key-name> --chain-id=irishub --fee=0.4iris
```
