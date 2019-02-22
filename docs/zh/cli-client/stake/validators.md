# iriscli stake validators

## 描述

查询所有验证者

## 用法

```
iriscli stake validators [flags]
```
打印帮助信息
```
iriscli stake validators --help
```

## 示例

查询验证者
```
iriscli stake validators
```

运行成功以后，返回的结果如下：

```txt
Validator
Operator Address: fva15grv3xg3ekxh9xrf79zd0w077krgv5xf6d6thd
Validator Consensus Pubkey: fcp1zcjduepq8fnuxnceuy4n0fzfc6rvf0spx56waw67lqkrhxwsxgnf8zgk0nus2r55he
Jailed: false
Status: Bonded
Tokens: 100.0000000000
Delegator Shares: 100.0000000000
Description: {node0   }
Bond Height: 0
Unbonding Height: 0
Minimum Unbonding Time: 1970-01-01 00:00:00 +0000 UTC
Commission: {{0.0000000000 0.0000000000 0.0000000000 0001-01-01 00:00:00 +0000 UTC}}
```
