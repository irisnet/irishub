# iriscli asset

## 描述

Asset 模块允许你在 IRIS Hub 上发行/管理资产

## 使用方式

```
 iriscli asset <command>
```

## 相关命令

| Name                                                | Description          |
| --------------------------------------------------- | -------------------- |
| [create-gateway](create-gateway.md)                 | 创建网关               |
| [edit-gateway](edit-gateway.md)                     | 编辑网关信息           |
| [transfer-gateway-owner](transfer-gateway-owner.md) | 转让网关所有权         |
| [issue-token](issue-token.md)                       | 发行资产              |
| [transfer-token-owner](transfer-token-owner.md)     | 转让资产所有权         |
| [edit-token](edit-token.md)                         | 编辑资产信息           |
| [mint-token](mint-token.md)                         | 增发资产               |
| [query-token](query-token.md)                       | 查询资产详情           |
| [query-tokens](query-tokens.md)                     | 查询资产列表           |
| [query-gateway](query-gateway.md)                   | 查询网关详情           |
| [query-gateways](query-gateways.md)                 | 查询网关列表           |
| [query-fee](query-fee.md)                           | 查询资产管理相关的费用  |

## 标志

| 命令，速记   | 默认值  | 描述          | 是否必须 |
| ---------- | ------ | ------------ | -------- |
| -h, --help |        | Asset 模块帮助 |         |

## 全局标志

| 命令，速记       | 默认值          | 描述                          | 是否必须 |
| -------------- | -------------- | ---------------------------- | -------- |
| -e, --encoding | hex            | 字符串二进制编码 (hex\b64\btc ) |         |
| --home         | /root/.iriscli | 配置和数据存储目录              |         |
| -o, --output   | text           | 输出格式 (text\json)          |          |
| --tr           |                |                              |          |