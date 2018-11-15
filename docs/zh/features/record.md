# Record User Guide

## 基本功能描述

1. 文本数据的链上存证
2. 文本文件的链上存证 (TODO)
3. 存证参数修改提议的链上治理 (TODO)

## 交互流程

### 存证流程

1. 任何用户可以发起存证请求，存证过程会花费用户一部分token，如果该数据的存证之前在链上并不存在,则该请求可以成功完成，并且链上会记录相关元数据，并返还给用户一个存证ID以确认该用户对这份数据的所有权。
2. 其他人若对完全相同的数据发起存证请求，则该请求会被直接否决，并提示相关存证数据已经存在。
3. 任何用户都可以根据存证ID在链上进行检索/下载操作。
4. 目前每次存证数据最大不超过1K，未来将结合治理模块实现参数的动态调整。

## 使用场景
### 创建使用环境

```
rm -rf iris
rm -rf .iriscli
iris init gen-tx --name=x --home=iris
iris init --gen-txs --chain-id=record-test -o --home=iris
iris start --home=iris
```

### 链上存证的使用场景

场景一：通过命令行对相关数据进行存证

```
# 根据--onchain-data指定需要存证的文本数据
iriscli record submit --description="test" --onchain-data=x --from=x --fee=0.04iris

# 结果
Committed at block 4 (tx hash: F649D5465A28842B50CAE1EE5950890E33379C45, response: {Code:0 Data:[114 101 99 111 114 100 58 97 98 57 100 57 57 100 48 99 102 102 54 53 51 100 99 54 101 56 53 52 53 99 56 99 99 55 50 101 53 53 51 51 100 101 97 97 97 49 50 53 53 50 53 52 97 100 102 100 54 98 48 48 55 52 101 50 56 54 57 54 54 49 98] Log:Msg 0:  Info: GasWanted:200000 GasUsed:3857 Tags:[{Key:[97 99 116 105 111 110] Value:[115 117 98 109 105 116 45 114 101 99 111 114 100]} {Key:[111 119 110 101 114 65 100 100 114 101 115 115] Value:[102 97 97 49 109 57 51 99 103 57 54 51 56 121 115 104 116 116 100 109 119 54 57 97 121 118 51 103 106 53 102 116 116 109 108 120 51 99 102 121 107 109]} {Key:[114 101 99 111 114 100 45 105 100] Value:[114 101 99 111 114 100 58 97 98 57 100 57 57 100 48 99 102 102 54 53 51 100 99 54 101 56 53 52 53 99 56 99 99 55 50 101 53 53 51 51 100 101 97 97 97 49 50 53 53 50 53 52 97 100 102 100 54 98 48 48 55 52 101 50 56 54 57 54 54 49 98]} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[2 189 149 142 250 208 0]}] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "submit-record",
     "completeConsumedTxFee-iris-atto": "\u0002\ufffd\ufffd\ufffd\ufffd\ufffd\u0000",
     "ownerAddress": "faa1m93cg9638yshttdmw69ayv3gj5fttmlx3cfykm",
     "record-id": "record:ab9d99d0cff653dc6e8545c8cc72e5533deaaa1255254adfd6b0074e2869661b"
   }
 }

# 查询存证情况
iriscli record query --record-id=x

# 下载存证数据
iriscli record download --record-id=x --file-name="download"

```

场景二：通过命令行查询链上包含存证数据的交易

```
# 查询存证上链情况
iriscli tendermint txs --tag "action='submit-record'"
```

## 命令详情

```
iriscli record submit --description="test" --onchain-data=x --chain-id="record-test" --from=x --fee=0.04iris
```

* `--onchain-data`  需要存证的数据


```
iriscli record query --record-id=x --chain-id="record-test"
```

* `--record-id` 待查询存证的ID


```
iriscli record download --record-id=x --file-name="download" --chain-id="record-test"
```

* `--file-name` 用来存放存证数据的文件名，其位于`--home`指定的目录中