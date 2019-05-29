# iriscli asset issue

## 描述

实现中IRIS Hub链上发行资产。

## 使用方式

发行100000000000iris
```
iriscli asset issue [flags]
```


## 标志

| 命令，缩写         | 类型       | 是否必须  | 默认值              | 描述                                                                                                                                             |
| ---------------- | --------- | ------- | ------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------ |
| --family         | byte      | 否       |  00                | 00 - fungible;01 - non-fungible                                                                                                                 |
| --name           | string    | 是       |                    | 新发行资产的名字，比如IRISnet。限定为最长32位的中英文字母、数字和_的组合                                                                                     |
| --symbol         | string    | 是       |                    | 新发行资产的symbol，比如iris。限定为3-6位的中英文字母、数字的组合                                                                                          |
| --source         | string    | 否       |  00                | 保留值 - 00(native); 01 (外部资产); 网关IDs                                                                                                         |
| --initial-supply | uint64    | 是       |                    | 该资产的初始发行量，该值不能超过1000亿                                                                                                                 |
| --max-supply     | uint64    | 是       |  1000000000000     | 该资产的最大发行量，资产的总供应量不能超过该值                                                                                                                 |
| --decimal        | uint8     | 否       |  0                 | 该资产值允许的最高小数位，最大为18位                                                                                                                   |
| --mintable       | boolean   | 否       |  false             | 初始发行后，该资产是否允许增发                                                                                                                        |
| --operators      | []Address | 否       |  []                | 资产的操纵者拥有除转让资产所有者的所有权限。但当资产所有者丢失私钥且存在两个及两个以上的资产操纵者时，操纵者们可以通过发送由超过2/3的资产操纵者签名的转让资产交易来转让资产所有者。|




## 例子

### 发行资产

```
iriscli asset issue iris --family=00 --name=IRISnet --symbol=iris --source=00 --initial-supply=100000000000 --max-supply=1000000000000 --decimal=0 --mintalbe=false  --operators=<account address A>,<account address B>

```

#### TODO:运行成功以后，返回的结果格式
