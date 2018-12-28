# iriscli tendermint txs

## 介绍

搜索查询所有符合指定匹配条件的交易

## 用法

```
iriscli tendermint txs [flags]
```

## 标志

| 名称，速记 | 默认值              |功能介绍                                                             | 是否必填     |
| --------------- | -------------------- | --------------------------------------------------------- | -------- |
| --chain-id      | ""                   | 区块链网络ID   | yes     |
| --node string   | tcp://localhost:26657| 节点查询rpc接口|
| --help, -h      |                      | 帮助信息 |    |
| --trust-node    | true                 | 是否信任查询节点     |          |
| --tags          | ""                   | 匹配条件     |          |
| --page          | 0                    | 分页的页码     |          |
| --size          | 100                  | 分页的大小     |          |

## 示例

### 查询交易

```shell
iriscli tendermint txs --tags action=send&sender=faa1c6al0vufl8efggzsvw34hszua9pr4qqyry37jn --chain-id=fuxi-4000 --trust-node=true
```

示例结果：

```
{
  "hash": "CD117378EC1CE0BA4ED0E0EBCED01AF09DA8F6B7",
  "height": "100722",
  "tx": {
    "type": "auth/StdTx",
    "value": {
      "msg": [
        {
          "type": "cosmos-sdk/Send",
          "value": {
            "inputs": [
              {
                "address": "faa1c6al0vufl8efggzsvw34hszua9pr4qqyry37jn",
                "coins": [
                  {
                    "denom": "iris-atto",
                    "amount": "3650000000000000000"
                  }
                ]
              }
            ],
            "outputs": [
              {
                "address": "faa1v2ezk7yvkgjq87ey54etfuxc87353ulrvq28z9",
                "coins": [
                  {
                    "denom": "iris-atto",
                    "amount": "3650000000000000000"
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
            "amount": "4787310000000000"
          }
        ],
        "gas": "6631"
      },
      "signatures": [
        {
          "pub_key": {
            "type": "tendermint/PubKeySecp256k1",
            "value": "A/ZQqJkDnqiN7maj4N9we8u8hE1dUpFD72+bD2PZgH+V"
          },
          "signature": "MEQCIEiNg0y3Xp9YgpY00cuYV6yoRIIXS1/Z7rOJeRwK8WipAiABfHZAS/yDMqPnBEPud1eJX8cZ6hhex1C7CGq286oclw==",
          "account_number": "162",
          "sequence": "3"
        }
      ],
      "memo": ""
    }
  },
  "result": {
    "log": "Msg 0: ",
    "gas_wanted": "6631",
    "gas_used": "4361",
    "tags": [
      {
        "key": "c2VuZGVy",
        "value": "ZmFhMWM2YWwwdnVmbDhlZmdnenN2dzM0aHN6dWE5cHI0cXF5cnkzN2pu"
      },
      {
        "key": "cmVjaXBpZW50",
        "value": "ZmFhMXYyZXprN3l2a2dqcTg3ZXk1NGV0ZnV4Yzg3MzUzdWxydnEyOHo5"
      },
      {
        "key": "Y29tcGxldGVDb25zdW1lZFR4RmVlLWlyaXMtYXR0bw==",
        "value": "MzE0ODQ2MzExNDE2MDc2MA=="
      }
    ]
  }
}

```
