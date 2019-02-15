# iriscli service

## 描述
Service 允许你通过区块链定义、绑定、调用服务。

## 用法

```shell
iriscli service [command]
```

## 相关命令

| Name                                  | Description            |
| ------------------------------------  | ---------------------- |
| [define](define.md)                   | 创建一个新的服务定义       |
| [definition](definition.md)           | 查询服务定义              |
| [bind](bind.md)                       | 创建一个新的服务绑定       |
| [binding](binding.md)                 | 查询服务绑定              |
| [bindings](bindings.md)               | 查询服务绑定列表           |
| [update-binding](update-binding.md)   | 更新一个存在的服务绑定      |
| [disable](disable.md)                 | 禁用一个可用的服务绑定      |
| [enable](enable.md)                   | 启用一个不可用的服务绑定     
| [refund-deposit](refund-deposit.md)   | 取回所有押金               |
| [call](call.md)                       | 调用服务方法                     |
| [requests](requests.md)               | 查询服务请求列表                     |
| [respond](respond.md)                 | 响应服务调用       |
| [response](response.md)               | 查询服务响应       |
| [fees](fees.md)                       | 查询指定地址的服务费退款和收入      |
| [refund-fees](refund-fees.md)         | 从服务费退款中退还所有费用  |
| [withdraw-fees](withdraw-fees.md)     | 从服务费收入中取回所有费用 |


## 标志

| 名称, 速记       | 默认值   | 描述            | 必需     |
| --------------- | ------- | ---------------- | -------- |
| --help, -h      |         | 服务命令帮助       |          |