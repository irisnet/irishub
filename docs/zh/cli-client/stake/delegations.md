# iriscli stake delegations

## 描述

查询某个委托者发起的所有委托记录

## 用法

```
iriscli stake delegations <delegator-address> [flags]
```

打印帮助信息
```
iriscli stake delegations --help
```

## 示例

### 查询某个委托者发起的所有委托记录

```
iriscli stake delegations iaa106nhdckyf996q69v3qdxwe6y7408pvyvyxzhxh
```

运行成功以后，返回的结果如下：

```
Delegation:
  Delegator:  faa1td4xnefkthfs6jg469x33shzf578fed6n7k7ua
  Validator:  fva1zkevgrasr5txhgsyqd7l9javln9et2d7k5yycy
  Shares:     1.0000000000000000000000000000
  Height:     26
```
