# iriscli asset transfer-gateway-owner

## 描述

转让指定Gateway的所有权到另一个地址。

## 使用方式

```bash
iriscli asset transfer-gateway-owner [flags]
```

## 特定的标志

| 命令, 速记 | 类型     | 是否必须 | 默认值 | 描述               |
| --------- | ------- | ------- | ----- | ----------------- |
| --moniker | string  | 是      |       | 被转让的网关唯一标识 |
| --to      | Address | 是      |       | 网关的新 Owner     |

## 示例

```bash
iriscli asset transfer-gateway-owner --moniker=cats --to=<new-owner-address> --from=<key-name> --chain-id=irishub --fee=0.3iris --commit
```
