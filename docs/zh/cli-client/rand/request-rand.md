# iriscli rand request-rand

## 描述

请求一个随机数

## 使用方式

```bash
iriscli rand request-rand [flags]
```

## 特有的标志

| 命令, 速记        | 类型    | 是否必须 | 默认值 | 描述                            |
| ---------------- | -----  | ------- | ----- | ------------------------------ |
| --block-interval | uint64 |         | 10    | 请求的随机数将在指定的区块间隔后生成 |

## 示例

```bash
iriscli rand request-rand --block-interval=100 --from=<key-name> --chain-id=irishub --fee=0.4iris --commit
```
