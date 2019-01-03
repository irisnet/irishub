# iriscli stake delegations-to

## 介绍

查询在某个验证人上的所有委托

## 用法

```
iriscli stake delegations-to [validator-address] [flags]
```
打印帮助信息
```
iriscli stake delegations-to --help
```

## 示例

查询在某个验证人上的所有委托
```
iriscli stake delegations-to fva1yclscskdtqu9rgufgws293wxp3njsesx7s40m2
```

示例结果

```json
[
  {
    "delegator_addr": "faa13lcwnxpyn2ea3skzmek64vvnp97jsk8qmhl6vx",
    "validator_addr": "fva1yclscskdtqu9rgufgws293wxp3njsesx7s40m2",
    "shares": "0.2000000000",
    "height": "290"
  }
]
```
