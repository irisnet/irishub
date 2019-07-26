# iriscli rand query-queue

## 描述

查询随机数请求队列，支持可选的高度

## 使用方式

```bash
iriscli rand query-queue [flags]
```

## 特有的标志

| 命令, 速记     | 类型   | 是否必须 | 默认值  | 描述                                                         |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --queue-height      | int64  | 否     |  0      | 欲查询的目标高度 |

## 示例

```bash
iriscli rand query-queue --queue-height=100000
```
