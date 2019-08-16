# iriscli bank burn

## 描述

通过减去账户余额的方式来销毁指定账户的一些token，任何人都可以使用这个交易。

## 使用方式

```bash
iriscli bank burn [flags]
```

## 标志

| 命令，缩写        | 类型    | 是否必须  | 默认值                 | 描述                                                         |
| ---------------- | ------ | -------- | --------------------- | ------------------------------------------------------------ |
| --amount         | String | 是       |                       | 销毁token数量，比如10iris                                      |

## 示例

### 销毁token

```bash
 iriscli bank burn --amount=10iris --from=<key-name> --fee=0.3iris --chain-id=irishub
```
