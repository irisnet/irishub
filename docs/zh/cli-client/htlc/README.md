# iriscli htlc

## 描述

管理 HTCL，同其他链上的资产进行原子交换

## 使用方式

```bash
 iriscli htlc [command]
```

## 相关命令

| 命令                         | 描述                                           |
| --------------------------- | --------------------------------------------- |
| [create](create.md)         | 创建 HTLC                                      |
| [claim](claim.md)           | 将一个 OPEN 状态的 HTLC 中锁定的资金发放到收款人地址 |
| [refund](refund.md)         | 从过期的 HTLC 中取回退款                         |
| [query-htlc](query-htlc.md) | 查询一个 HTLC 的详细信息                         |

## 全局标志

| 命令，速记       | 默认值          | 描述                         | 是否必须 | 类型   |
| -------------- | -------------- | --------------------------- | ------ | ------ |
| -e, --encoding | hex            | 字符串二进制编码 (hex\b64\btc) |        | string |
| --home         | /root/.iriscli | 配置和数据存储目录             |        | string |
| -o, --output   | text           | 输出格式 (text\json)         |        | string |
| --trace        |                | 出错时打印完整栈信息           |        |        |

## 标志

| 命令，速记   | 默认值  | 描述        | 是否必须 |
| ---------- | ------ | ----------- | ------ |
| -h, --help |        | HTLC模块帮助 |        |
