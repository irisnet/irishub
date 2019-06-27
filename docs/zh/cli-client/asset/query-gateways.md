# iriscli asset query-gateways

## 描述

查询所有网关信息，支持可选的owner参数

## 使用方式

```
iriscli asset query-gateways [flags]
```

打印帮助信息:
```
iriscli asset query-gateways --help
```

## 特定的标志

| 命令, 速记     | 类型   | 是否必须 | 默认值  | 描述                                                         |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --owner           | Address | 否     |        | 要查询的目标所有者地址 |


## 示例

```
iriscli asset query-gateways --owner=faa1an4wfvsnxrp97lug5fngct6melhgcuvdv2qye3
```

输出信息:
```txt
Gateways for owner faa1an4wfvsnxrp97lug5fngct6melhgcuvdv2qye3:
  Moniker: tgw, Identity: exchange, Details: testgateway, Website: http://testgateway.io
  Moniker: tgwx, Identity: exchange, Details: testgateway2, Website: http://testgateway2.io
```

```json
[
  {
    "owner": "faa1an4wfvsnxrp97lug5fngct6melhgcuvdv2qye3",
    "moniker": "tgw",
    "identity": "exchange",
    "details": "testgateway",
    "website": "http://testgateway.io"
  },
  {
    "owner": "faa1an4wfvsnxrp97lug5fngct6melhgcuvdv2qye3",
    "moniker": "tgwx",
    "identity": "exchange",
    "details": "testgateway2",
    "website": "http://testgateway2.io"
  }
]
```