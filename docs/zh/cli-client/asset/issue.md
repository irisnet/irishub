# iriscli asset issue

## 描述

实现中IRIS Hub链上发行资产。

## 使用方式

发行10000000000 kitty
```
iriscli asset issue [flags]
```


## 标志

| 命令，缩写         | 类型       | 是否必须  | 默认值              | 描述                                                                                                                                             |
| ---------------- | --------- | ------- | ------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------ |
| --family          | string    | 否       |  fungible          | 取值可为fungible 或 non-fungible                                                                                                                 |
| --name            | string    | 是       |                    | 新发行资产的名字，比如IRISnet。限定为最长32位的中英文字母、数字和_的组合                                                                                     |
| --symbol          | string    | 是       |                    | 新发行资产的唯一标志，限定为3-8位的中英文字母、数字和_的组合                                                                                                                                                          |
| --symbol-at-source| string    | 否       |                    | 对于外部资产和网关资产，可用于指定值源链的symbol                                                                                                                                                                                  |
| --symbol-min-alias| string    | 否       |                    | 资产最小单位别名，alphanumeric，首字符必须为字母，长度[3,10]，toLowerCase                                                                                                                                                                                   |
| --source          | string    | 否       |  native            | 可选：native; external; gateway                                                                                                             |
| --initial-supply  | uint64    | 是       |                    | 该资产的初始发行量，该值不能超过1000亿                                                                                                                 |
| --max-supply      | uint64    | 否       |  1000000000000     | 该资产的最大发行量，资产的总供应量不能超过该值                                                                                                                 |
| --decimal         | uint8     | 否       |  0                 | 该资产值允许的最高小数位，最大为18位                                                                                                                   |
| --mintable        | boolean   | 否       |  false             | 初始发行后，该资产是否允许增发                                                                                                                        |
| --operators       | []Address | 否       |  []                | 资产的操纵者拥有除转让资产所有者的所有权限。但当资产所有者丢失私钥且存在两个及两个以上的资产操纵者时，操纵者们可以通过发送由超过2/3的资产操纵者签名的转让资产交易来转让资产所有者。|




## 例子

### 发行本地资产


```
iriscli asset issue-token --family=fungible --name=kittyToken --symbol=kitty --source=native --initial-supply=100000000000 --max-supply=1000000000000 --decimal=0 --mintable=<true/false>    --fee=1iris
```


### 发行网关资产

```
iriscli asset issue-token --family=fungible --symbol-at-source=cat --name=kittyToken --symbol=kitty --source=gateway --gateway=gtty --initial-supply=100000000000 --max-supply=1000000000000 --decimal=0 --mintable=<true/false>  --fee=1iris
```

输出信息:
```txt
Committed at block 306 (tx hash: 5A4C6E00F4F6BF795EB05D2D388CBA0E8A6E6CF17669314B1EE6A31729A22450, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:3398 Tags:[{Key:[97 99 116 105 111 110] Value:[115 101 114 118 105 99 101 45 119 105 116 104 100 114 97 119 45 102 101 101 115] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 54 55 57 54 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
```

```json
  {
    "code": 0,
    "data": null,
    "log": "Msg 0: ",
    "info": "",
    "gas_wanted": 50000,
    "gas_used": 7008,
    "codespace": "",
    "tags": [
      {
        "key": "action",
        "value": "issue_token"
      },
      {
        "key": "token-id",
        "value": "kitty"
      },
      {
        "key": "token-denom",
        "value": "kitty-min"
      },
      {
        "key": "token-source",
        "value": "native"
      },
      {
        "key": "token-gateway",
        "value": ""
      },
      {
        "key": "token-owner",
        "value": "faa1j8mlkem7s9a0jjhkd39zd24xh9gdj6zus77v7q"
      }
    ]
  })
 

```
