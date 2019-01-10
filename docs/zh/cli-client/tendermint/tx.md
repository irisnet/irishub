# iriscli tendermint tx

## 描述

在所有提交的区块上寻找匹配此txhash的交易

## 用法

```
iriscli tendermint tx [hash] [flags]
```

## 标志

| 名称, 速记 | 默认值                    | 描述                                                             | 必需      |
| --------------- | -------------------------- | --------------------------------------------------------- | -------- |
| --chain-id    | 无 | [string] tendermint节点的链ID   | 是       |
| --node string     |   tcp://localhost:26657                         | 要连接的节点  |                                     
| --help, -h      |           无| 	下载命令帮助|
| --trust-node    | true                       | 信任连接的完整节点，关闭响应结果校验                                            |          |

## 例子

### 

```shell
iriscli tendermint tx 6D361F4BE14B6B9596104A9411DB7962501CA264 --chain-id=irishub-test --trust-node=true
```
得到得结果如下：
```
{
  "hash": "6D361F4BE14B6B9596104A9411DB7962501CA264",
  "height": "111293",
  "tx": {
    "type": "auth/StdTx",
    "value": {
      "msg": [
        {
          "type": "cosmos-sdk/Send",
          "value": {
            "inputs": [
              {
                "address": "faa14d59j363uha3xxk2xqqhav8xaztym2rkdjdf8v",
                "coins": [
                  {
                    "denom": "iris-atto",
                    "amount": "10000000000000000000"
                  }
                ]
              }
            ],
            "outputs": [
              {
                "address": "faa18q343g6wje2plmhfekmwyw83jhznqltreqyvs3",
                "coins": [
                  {
                    "denom": "iris-atto",
                    "amount": "10000000000000000000"
                  }
                ]
              }
            ]
          }
        }
      ],
      "fee": {
        "amount": [
          {
            "denom": "iris-atto",
            "amount": "4000000000000000"
          }
        ],
        "gas": "10000"
      },
      "signatures": [
        {
          "pub_key": {
            "type": "tendermint/PubKeySecp256k1",
            "value": "AwuoO1hR2F/p6F7fVqMwZ3L7ugcxbQJp3jSYiFrucsq6"
          },
          "signature": "MEQCIHs8ZuwVpBRQr/LmGihyeuuDW2or3/LEJJNao6KmZMnpAiA2dhfgfLwskqq4M5QOKXqELT6TqXHEA7f/SR4ghDQUMQ==",
          "account_number": "26",
          "sequence": "716"
        }
      ],
      "memo": ""
    }
  },
  "result": {
    "log": "Msg 0: ",
    "gas_wanted": "10000",
    "gas_used": "4496",
    "tags": [
      {
        "key": "c2VuZGVy",
        "value": "ZmFhMTRkNTlqMzYzdWhhM3h4azJ4cXFoYXY4eGF6dHltMnJrZGpkZjh2"
      },
      {
        "key": "cmVjaXBpZW50",
        "value": "ZmFhMThxMzQzZzZ3amUycGxtaGZla213eXc4M2poem5xbHRyZXF5dnMz"
      },
      {
        "key": "Y29tcGxldGVDb25zdW1lZFR4RmVlLWlyaXMtYXR0bw==",
        "value": "MTc5ODQwMDAwMDAwMDAwMA=="
      }
    ]
  }
}
```
