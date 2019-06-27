# iriscli asset query-gateway

## 描述

查询由指定的名字所标识的网关信息。

## 使用方式

```
iriscli asset query-gateway [flags]
```

打印帮助信息:
```
iriscli asset query-gateway --help
```

## 特有的标志

| 命令, 速记     | 类型   | 是否必须 | 默认值  | 描述                                                         |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --moniker           | string | 是     | ""       | 唯一的网关名字, 长度为3-8个英文字符 |


## 示例

```
iriscli asset query-gateway --moniker tgw
```

输出信息:
```txt
Gateway:
  Owner:             faa1an4wfvsnxrp97lug5fngct6melhgcuvdv2qye3
  Moniker:           tgw
  Identity:          exchange
  Details:           testgateway
  Website:           http://testgateway.io
```

```json
{
  "type": "irishub/asset/Gateway",
  "value": {
    "owner": "faa1an4wfvsnxrp97lug5fngct6melhgcuvdv2qye3",
    "moniker": "tgw",
    "identity": "exchange",
    "details": "testgateway",
    "website": "http://testgateway.io"
  }
}
```