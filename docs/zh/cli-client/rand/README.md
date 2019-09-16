# iriscli rand

## Description

此模块允许你向 IRIS Hub 发送随机数请求，查询随机数或待处理的随机数请求队列

## Usage

```bash
iriscli rand <command>
```

打印子命令和参数：

```bash
iriscli rand --help
```

## 相关命令

| 命令                             | 描述                          |
| ------------------------------- | ----------------------------- |
| [request-rand](request-rand.md) | 请求一个随机数                  |
| [query-rand](query-rand.md)     | 使用ID查询链上生成的随机数        |
| [query-queue](query-queue.md)   | 查询随机数请求队列，支持可选的高度 |
