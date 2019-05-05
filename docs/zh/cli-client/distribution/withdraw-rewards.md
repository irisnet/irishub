# iriscli distribution withdraw-rewards

## 介绍

取回收益

## 用法

```
iriscli distribution withdraw-rewards <flags>
```

打印帮助信息:

```
iriscli distribution withdraw-rewards --help
```

## 特有标志位

| 名称                | 类型   | 是否必填 | 默认值  | 功能描述        |
| --------------------- | -----  | -------- | -------- | ------------------------------------------------------------------- |
| --only-from-validator | string | false    | ""       | 验证人地址，仅取回在这个验证人上的委托收益 |
| --is-validator        | bool   | false    | false    | 取回验证人佣金收益 |

不能同时使用两个flags。

## 示例

1. 仅取回在某一个验证人处的委托收益
    ```
    iriscli distribution withdraw-rewards --only-from-validator=<validator_address> --from=<key_name> --fee=0.3iris --chain-id=<chain-id>
    ```
2. 取回所有在外的委托收益，不包含验证人的佣金收益:
    ```
    iriscli distribution withdraw-rewards --from=<key_name> --fee=0.3iris --chain-id=<chain-id>
    ```
3. 验证人取回所有委托产生的收益以及验证人的佣金收益（仅限验证人）:
    ```
    iriscli distribution withdraw-rewards --is-validator=true --from=<key_name> --fee=0.3iris --chain-id=<chain-id>
    ```
