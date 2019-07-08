# iriscli bank send

## 描述

发送通证到指定地址, 此条命令同时包含`创建/签名/广播`交易。

## 使用方式

```bash
iriscli bank send [flags]
```

## 标志

| 命令，缩写       | 类型   | 是否必须 | 默认值                | 描述                                                         |
| ---------------- | ------ | -------- | --------------------- | ------------------------------------------------------------ |
| --amount         | String | 是       |                       | 需要发送的通证数量，比如10iris                               |
| --to             | String | 是       |                       | Bech32 编码的接收通证的地址。                                |

## 示例

### 发送通证到指定地址

```bash
 iriscli bank send --from=<key-name> --to=<address> --amount=10iris --fee=0.3iris --chain-id=irishub
```
