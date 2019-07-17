# iriscli gov submit-proposal

## 描述

提交区块链治理提议以及发起提议所涉及的初始保证金，其中提议的类型包括PlainText/ParameterChange/SoftwareUpgrade/AddToken这几种类型。

## 使用方式

```
iriscli gov submit-proposal <flags>
```
打印帮助信息:

```
iriscli gov submit-proposal --help
```
## 标志

| 名称, 速记        | 默认值                      | 描述                                                                                                                                                 | 是否必须  |
| ---------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --deposit        |                            | 提议的保证金（至少30%minDeposit）                                                                                                                         |          |
| --description    |                            | 提议的描述                                                                                                           | Yes      |
| --param          |                            | 提议参数,例如: mint/Inflation=0.050                                                                                |          |
| --title          |                            | 提议标题                                                                                                                           | Yes      |
| --type           |                            | 提议类型,例如:PlainText/ParameterChange/SoftwareUpgrade/SoftwareHalt/CommunityTaxUsage/AddToken                                                                  | Yes      |
| --version           |            0                | 新协议的版本信息                                                                           |       |
| --software           |           " "                 |  新协议的软件地址                                                                       |       |
| --switch-height           |       0                     |  新版本协议升级的高度                                                     |       |
| --threshold        | "0.8"   |  软件升级的阈值                                              |               |
| token-canonical-symbol   |        | 外部代币名称                                                 | |
| --token-symbol |  | 代币符号 | |
| --token-name |  | 代币名称 | |
| --token-decimal |  | 代币最大精度 | |
| --token-min-unit-alias |  | 代币最小单温别名 | |
| --token-initial-supply |  | 代币初始总量 | |

## 例子

提议者必须抵押至少30%的`MinDeposit`，详情见 [Gov](../../features/governance.md)

### 提交一个`ParameterChange`类型的提议

修改Inflation参数的提议：

```shell
iriscli gov submit-proposal --chain-id=<chain-id> --title=<proposal_title> --param='mint/Inflation=0.050' --type=ParameterChange --description=<proposal_description> --from=<key_name> --fee=0.3iris --deposit="3000iris" 
```

param的值可以通过 `iriscli params`查询(Gov模块本身参数不可以修改),详细[请参考](../params/README.md)

### 提交一个`SoftwareUpgrade`类型的提议

```shell
iriscli gov submit-proposal --chain-id=<chain-id> --title=<proposal_title> --description=<proposal_description>  --type=SoftwareUpgrade --description=<proposal_description> --from=<key_name> --fee=0.3iris --software=https://github.com/irisnet/irishub/tree/v0.9.0 --version=2 --switch-height=80 --threshold=0.9 --deposit="3000iris" 
```

### 提交一个`AddToken`类型的提议

```shell
iriscli gov submit-proposal --chain-id=irishub-test --from=node0 --fee=4iris --type=AddToken --description=test --title=test-proposal --deposit=50000iris --commit --home=$iris_root_path --token-symbol=btc --token-canonical-symbol=btc --token-name=btcToken --token-decimal=18 --token-min-unit-alias=atto --token-initial-supply=200000
```

###  如何查询提议详情？

[query-proposal](query-proposal.md)

[query-proposals](query-proposals.md)
