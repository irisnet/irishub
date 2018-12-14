# iriscli tendermint txs

## Description

Search all transactions which match the given tag list

## Usage

```
iriscli tendermint txs [flags]

```

## Flags

| Name, shorthand | Default              |Description                                                             | Required     |
| --------------- | -------------------- | --------------------------------------------------------- | -------- |
| --chain-id      | ""                   | Chain ID of Tendermint node   | yes     |
| --node string   | tcp://localhost:26657| Node to connect to (default "tcp://localhost:26657")  |
| --help, -h      |                      | 	help for txs|    |
| --trust-node    | true                 | Trust connected full node (don't verify proofs for responses)     |          |
| --tags          | ""                   | tag:value list of tags that must match     |          |
| --page          | 0                    | Pagination page     |          |
| --size          | 100                  | Pagination size     |          |

## Examples

### Search transactions

```shell
iriscli tendermint txs --tags action=send&sender=faa1c6al0vufl8efggzsvw34hszua9pr4qqyry37jn --chain-id=fuxi-4000 --trust-node=true
```

You will get the following result.

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
