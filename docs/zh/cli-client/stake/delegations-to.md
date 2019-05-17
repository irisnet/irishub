# iriscli stake delegations-to

## 介绍

查询在某个验证人上的所有委托

## 用法

```
iriscli stake delegations-to <validator-address> <flags>
```

打印帮助信息
```
iriscli stake delegations-to --help
```

## 示例

查询在某个验证人上的所有委托
```
iriscli stake delegations-to iva1yclscskdtqu9rgufgws293wxp3njsesxxlnhmh
```

示例结果
```
Delegation:
  Delegator:  iaa13lcwnxpyn2ea3skzmek64vvnp97jsk8qrcezvm
  Validator:  iva1yclscskdtqu9rgufgws293wxp3njsesxxlnhmh
  Shares:     100.0000000000000000000000000000
  Height:     0
Delegation:
  Delegator:  iaa1td4xnefkthfs6jg469x33shzf578fed6n7k7ua
  Validator:  iva1yclscskdtqu9rgufgws293wxp3njsesxxlnhmh
  Shares:     1.0000000000000000000000000000
  Height:     26
```
