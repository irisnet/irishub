# iriscli distribution withdraw-rewards

## 介绍

取回收益

## 用法

```
iriscli distribution withdraw-rewards [flags]
```

打印帮助信息:

```
iriscli distribution withdraw-rewards --help
```

## 特有的flags

| Name, shorthand       | type   | Required | Default  | Description                                                         |
| --------------------- | -----  | -------- | -------- | ------------------------------------------------------------------- |
| --only-from-validator | string | false    | ""       | only withdraw from this validator address (in bech) |
| --is-validator        | bool   | false    | false    | Also withdraw validator's commission |

不能同时使用两个flags。

## 示例

1. 仅取回某一个委托产生的收益
    ```
    iriscli distribution withdraw-rewards --only-from-validator fva134mhjjyyc7mehvaay0f3d4hj8qx3ee3w3eq5nq --from mykey --fee=0.004iris --chain-id=irishub-test
    ```
2. 取回所有委托产生的收益，不包含验证人的佣金收益:
    ```
    iriscli distribution withdraw-rewards --from mykey --fee=0.004iris --chain-id=irishub-test
    ```
3. 取回所有委托产生的收益以及验证人的佣金收益:
    ```
    iriscli distribution withdraw-rewards --is-validator=true --from mykey --fee=0.004iris --chain-id=irishub-test
    ```