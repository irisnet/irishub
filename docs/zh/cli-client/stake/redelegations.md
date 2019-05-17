# iriscli stake redelegations

## 描述

基于委托者地址的所有转委托记录查询

## 用法

```
iriscli stake redelegations <delegator-address> <flags>
```

打印帮助信息
```
iriscli stake redelegations --help
```

## 示例

基于委托者地址的所有转委托记录查询
```
iriscli stake redelegations <delegator-address> 
```

运行成功以后，返回的结果如下：
```
Redelegation:
Delegator: iaa10s0arq9khpl0cfzng3qgxcxq0ny6hmc9gtd2ft
Source Validator: iva1dayujdfnxjggd5ydlvvgkerp2supknth9a8qhh
Destination Validator: iva1h27xdw6t9l5jgvun76qdu45kgrx9lqedpg3ecs
Creation height: 1130
Min time to unbond (unix): 2018-11-16 07:22:48.740311064 +0000 UTC
Source shares: 0.1000000000
Destination shares: 0.1000000000
```
