# iriscli tendermint tx

## Description

Search for the transaction which has the same hash in all existing blocks

## Usage

```bash
iriscli tendermint tx <hash> <flags>
```

## Flags

| Name, shorthand | Default               | Description                                                   | Required |
| --------------- | --------------------- | ------------------------------------------------------------- | -------- |
| --chain-id      |                       | Chain ID of Tendermint node                                   | yes      |
| --node string   | tcp://localhost:26657 | Node to connect to                                            |
| --help, -h      |                       | help for tx                                                   |          |
| --trust-node    | true                  | Trust connected full node (don't verify proofs for responses) |          |

## Examples

### tx

```shell
iriscli tendermint tx CD117378EC1CE0BA4ED0E0EBCED01AF09DA8F6B7 --chain-id=<chain-id> --trust-node
```

You will get the following result.

```json
{
  "hash": "50F8D75FC1F0C2643A0D09189B7FB44246AB00AF89779215FFBC0740E6C59F3A",
  "height": "3411",
  "tx": {
    "type": "irishub/bank/StdTx",
    "value": {
      "msg": [
        {
          "type": "irishub/bank/Send",
          "value": {
            "inputs": [
              {
                "address": "faa10t6tn0ntgrzetmzwlr9x8fj4j29qrcax0p52dm",
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
                "address": "faa1m9m9t8paa48xgmaxg7gxzq3a5rcl4neecm4f94",
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
            "amount": "400000000000000000"
          }
        ],
        "gas": "50000"
      },
      "signatures": [
        {
          "pub_key": {
            "type": "tendermint/PubKeySecp256k1",
            "value": "AtRMRPAKXstV6/NN8cizi55lMeqrtyzkR6UmSMcujYpG"
          },
          "signature": "ctjaNgszonLxoVd2weWe1TleCxg8vmSoYNuJNI1OEE5Ll/+NY0PEnDHeUsTkq71t8HgYkFkM636EssP9TAmttQ==",
          "account_number": "2",
          "sequence": "1"
        }
      ],
      "memo": ""
    }
  },
  "result": {
    "Code": 0,
    "Data": null,
    "Log": "Msg 0: ",
    "Info": "",
    "GasWanted": "50000",
    "GasUsed": "6678",
    "Tags": [
      {
        "key": "action",
        "value": "send"
      },
      {
        "key": "sender",
        "value": "faa10t6tn0ntgrzetmzwlr9x8fj4j29qrcax0p52dm"
      },
      {
        "key": "recipient",
        "value": "faa1m9m9t8paa48xgmaxg7gxzq3a5rcl4neecm4f94"
      }
    ],
    "Codespace": "",
    "XXX_NoUnkeyedLiteral": {},
    "XXX_unrecognized": null,
    "XXX_sizecache": 0
  },
  "timestamp": "2019-07-01T07:40:05Z",
  "coin_flow": null
}
```
