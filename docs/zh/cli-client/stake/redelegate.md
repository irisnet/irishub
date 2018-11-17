# iriscli stake redelegate

## 介绍

把某个委托的一部分或者全部从一个验证人转移到另外一个验证人

## 用法

```
iriscli stake redelegate [flags]
```

打印帮助信息
```
iriscli stake redelegate --help
```

## 特有flags

| 名称                       | 类型   | 是否必填 | 默认值   | 功能描述         |
| -------------------------- | -----  | -------- | -------- | ------------------------------------------------------------------- |
| --address-validator-dest   | string | true     | ""       | 目标验证人地址 |
| --address-validator-source | string | true     | ""       | 源验证人地址 |
| --shares-amount            | float  | false    | 0.0      | 转移的share数量，正数 |
| --shares-percent           | float  | false    | 0.0      | 转移的比率，0到1之间的正数 |

用户可以用`--shares-amount`或者`--shares-percent`指定装委托的token数量。记住，这两个参数不可同时使用。

## 示例

```
iriscli stake redelegate --chain-id=<chain-id> --from=<key name> --fee=0.004iris --address-validator-source=<SourceValidatorAddress> --address-validator-dest=<DestinationValidatorAddress> --shares-percent=0.1
```

