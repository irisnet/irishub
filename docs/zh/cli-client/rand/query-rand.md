# iriscli rand query-rand

## 描述

使用ID查询链上生成的随机数

## 使用方式

```bash
iriscli rand query-rand [flags]
```

## 特有的标志

| 命令, 速记     | 类型    | 是否必须 | 默认值 | 描述                 |
| ------------- | -----  | ------- | ----- | ------------------- |
| --request-id  | string | 是      |       | 请求ID, 由请求交易返回 |

## 示例

```bash
iriscli rand query-rand --request-id=035a8d4cf64fcd428b5c77b1ca85bfed172d3787be9bdf0887bbe8bbeec3932c
```
