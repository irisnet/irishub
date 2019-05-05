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
```json
[
  {
    "delegator_addr": "iaa10s0arq9khpl0cfzng3qgxcxq0ny6hmc9gtd2ft",
    "validator_src_addr": "iva1dayujdfnxjggd5ydlvvgkerp2supknth9a8qhh",
    "validator_dst_addr": "iva1h27xdw6t9l5jgvun76qdu45kgrx9lqedpg3ecs",
    "creation_height": "1130",
    "min_time": "2018-11-16T07:22:48.740311064Z",
    "initial_balance": "0.1iris",
    "balance": "0.1iris",
    "shares_src": "0.1000000000",
    "shares_dst": "0.1000000000"
  }
]
```
