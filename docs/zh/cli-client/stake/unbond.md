# iriscli stake unbond

## 介绍

Unbond shares from a validator

## 用法

```
iriscli stake unbond [flags]
```

打印帮助信息

```
iriscli stake unbond --help
```

## 特有flags

| 名称                | 类型   | 是否必填 | 默认值   | 功能描述         |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --address-validator | string | true     | ""       | 验证人地址 |
| --shares-amount     | float  | false    | 0.0      | 解绑的share数量，正数 |
| --shares-percent    | float  | false    | 0.0      | 解绑的比率，0到1之间的正数 |

用户可以用`--shares-amount`或者`--shares-percent`指定解绑定的token数量，这两个参数不可同时使用。

## 示例

```
iriscli stake unbond --address-validator=<ValidatorAddress> --shares-percent=0.1 --from=<key name> --chain-id=<chain-id> --fee=0.004iris
```
