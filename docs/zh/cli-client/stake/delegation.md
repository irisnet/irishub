# iriscli stake delegation

## 描述

基于委托者和验证者地址查询委托交易

## 用法

```
iriscli stake delegation [flags]
```
打印帮助信息
```
iriscli stake delegation --help
```
## 特有的flags

| 名称, 速记             | 默认值                      | 描述                                                                 | 必需     |
| --------------------- | -------------------------- | -------------------------------------------------------------------- | -------- |
| --address-delegator   |                            | [string] 委托者bech地址                                               | Yes      |
| --address-validator   |                            | [string] 验证者bech地址                                               | Yes      |

## 示例

### 查询验证者

```
iriscli stake delegation --address-validator=iva106nhdckyf996q69v3qdxwe6y7408pvyv3hgcms --address-delegator=iaa106nhdckyf996q69v3qdxwe6y7408pvyvyxzhxh

```

运行成功以后，返回的结果如下：

```txt
Delegation
Delegator: iaa13lcwnxpyn2ea3skzmek64vvnp97jsk8qrcezvm
Validator: iva15grv3xg3ekxh9xrf79zd0w077krgv5xfzzunhs
Shares: 200.0000000
Height: 290
```
