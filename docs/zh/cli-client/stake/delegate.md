# iriscli stake delegate

## 介绍

发送委托交易

## 用法

```
iriscli stake delegate [flags]
```

打印帮助信息
```
iriscli stake delegate --help
```

## 特有的flags

| 名称                | 类型   | 是否必填 | 默认值   | 功能描述         |
| --------------------| -----  | -------- | -------- | ---------------- |
| --address-delegator | string | true     | ""       | 验证人地址 |
| --amount            | string | true     | ""       | 委托的token数量 |

## 示例

```
iriscli stake delegate --chain-id=ChainID --from=KeyName --fee=Fee --amount=CoinstoBond --address-validator=ValidatorOwnerAddress
```
