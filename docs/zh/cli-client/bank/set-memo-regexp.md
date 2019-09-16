# iriscli bank set-memo-regexp

## 描述

设置校验 memo 的正则表达式

## 使用方式

```bash
iriscli bank set-memo-regexp --regexp=<regular-expression> --from=<key-name> --fee=<native-fee> --chain-id=<chain-id>
```

## 标志

| 命令，速记 | 类型    | 是否必须  | 默认值 | 描述                                      |
| -------- | ------ | -------- | ----- | ---------------------------------------- |
| --regexp | string | 是       |       | 正则表达式，最大长度50， 例如 ^[A-Za-z0-9]+$ |

## 示例

### 交易发送者为自己账户设置 memo 正则表达式

```bash
iriscli bank set-memo-regexp --regexp=^[A-Za-z0-9]+$ --from=<key-name> --fee=0.3iris --chain-id=irishub
```
