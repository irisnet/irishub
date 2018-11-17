# iriscli distribution withdraw-address

## 介绍

查收取回收益时的收款地址

## 用法

```
iriscli distribution withdraw-address [delegator address] [flags]
```

打印帮助信息:

```
iriscli distribution withdraw-address --help
```

## 特有的flags

这个命令没有特有的flag，它有一个输入参数：委托人地址


## 示例

```
iriscli distribution withdraw-address <delegator address>
```
执行结果示例
```
faa1ezzh0humhy3329xg4avhcjtay985nll0zswc5j
```
If the given delegator didn't specify other withdraw address, the query result will be empty.