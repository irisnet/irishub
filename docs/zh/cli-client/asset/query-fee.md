# iriscli asset query-fee

## 描述

查询资产相关的费用，包括网关创建费用以及Token发行和增发费用

## 使用方式

```
iriscli asset query-fee [flags]
```

打印帮助信息:
```
iriscli asset query-fee --help
```

## 特有的标志

| 命令, 速记     | 类型   | 是否必须 | 默认值  | 描述                                          |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --subject           | string | 是     | ""      | 费用类型, 取值为gateway或者token         |
| --moniker           | string | 否    | ""       | 网关名字; 如果subject为gateway,则必须指定 |
| --id                | string | 否    | ""       | Token ID; 如果subject为token,则必须指定  |


## 示例

```
iriscli asset query-fee --subject gateway --moniker tgw
```

输出信息:
```txt
Fee: 600000iris
```

```json
{
  "exist": false,
  "fee": {
    "denom": "iris",
    "amount": "600000"
  }
}
```