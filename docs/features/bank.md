# Bank User Guide

## Introduction 
This module is mainly used to transfer coins between accounts、query account balances, and provide a common offline transaction signing and broadcasting method. In addition, the available units of tokens in the IRIShub system are defined using [coin_type](./basic-concepts/coin-type.md).

## Usage Scenario

1. Query the coin_type configuration of a certain token:
    ```bash
    iriscli bank coin-type [coin-name]
    ```
    For example, coin_type of iris will be returned if the coin-name is iris:
    ```json
    {
     "name": "iris",
     "min_unit": {
       "denom": "iris-atto",
       "decimal": "18"
     },
     "units": [
       {
         "denom": "iris",
         "decimal": "0"
       },
       {
         "denom": "iris-milli",
         "decimal": "3"
       },
       {
         "denom": "iris-micro",
         "decimal": "6"
       },
       {
         "denom": "iris-nano",
         "decimal": "9"
       },
       {
         "denom": "iris-pico",
         "decimal": "12"
       },
       {
         "denom": "iris-femto",
         "decimal": "15"
       },
       {
         "denom": "iris-atto",
         "decimal": "18"
       }
     ],
     "origin": 1,
     "desc": "IRIS Network"
    }
    ```

2. Query account

    Query the account information of a certain account address, including the balance, the public key, the account number and the transaction number.
    ```bash
    iriscli bank account [account address]
    ```

3. Transfer between accounts

    For example, transfer from account A to account B10iris:
    ```bash
    iriscli bank send --to [address of wallet B] --amount=10iris --fee=0.004iris --from=[key name of wallet A] --chain-id=[chain-id]
    ```
    IRISnet supports multiple tokens in circulation, and in the future IRISnet will be able to include multiple tokens in one transaction -- tokens can be any coin_type registered in IRISnet. 

4. Sign transactions generated offline

    To improve account security, IRISnet supports offline signing of transactions to protect the account's private key. In any transaction, you can build an unsigned transaction using the flag --generate-only=true. Use transfer transactions as an example:
    ```bash
    iriscli bank send --to [address of wallet B] --amount=10iris --fee=0.004iris --from=[key name of wallet A] --generate-only=true
    ```
    Return the built transaction with empty signatures:
    ```json
    {
      "type": "auth/StdTx",
      "value": {
        "msg": [
          {
            "type": "cosmos-sdk/Send",
            "value": {
              "inputs": [
                {
                  "address": "faa1ydhmma8l4m9dygsh7l08fgrwka6yczs0gkfnvd",
                  "coins": [
                    {
                      "denom": "iris-atto",
                      "amount": "100000000000000000000"
                    }
                  ]
                }
              ],
              "outputs": [
                {
                  "address": "faa1ut8aues05kq0nkcj3lzkyhk7eyfasrdfnf7wph",
                  "coins": [
                    {
                      "denom": "iris-atto",
                      "amount": "100000000000000000000"
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
              "amount": "40000000000000000"
            }
          ],
          "gas": "200000"
        },
        "signatures": null,
        "memo": ""
      }
    }
    ```
    Save the result to a file.
    
    Send signature transaction:
    ```bash
    iriscli bank sign [file] --chain-id=[chain-id]  --name [key name] 
    ```
    Return signed transactions:
    ```json
    {
      "type": "auth/StdTx",
      "value": {
        "msg": [
          {
            "type": "cosmos-sdk/Send",
            "value": {
              "inputs": [
                {
                  "address": "faa1ydhmma8l4m9dygsh7l08fgrwka6yczs0gkfnvd",
                  "coins": [
                    {
                      "denom": "iris-atto",
                      "amount": "100000000000000000000"
                    }
                  ]
                }
              ],
              "outputs": [
                {
                  "address": "faa1ut8aues05kq0nkcj3lzkyhk7eyfasrdfnf7wph",
                  "coins": [
                    {
                      "denom": "iris-atto",
                      "amount": "100000000000000000000"
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
              "amount": "40000000000000000"
            }
          ],
          "gas": "200000"
        },
        "signatures": [
          {
            "pub_key": {
              "type": "tendermint/PubKeySecp256k1",
              "value": "A+qXW5isQDb7blT/KwEgQHepji8RfpzIstkHpKoZq0kr"
            },
            "signature": "5hxk/R81SWmKAGi4kTW2OapezQZpp6zEnaJbVcyDiWRfgBm4Uejq8+CDk6uzk0aFSgAZzz06E014UkgGpelU7w==",
            "account_number": "0",
            "sequence": "11"
          }
        ],
        "memo": ""
      }
    }
    ```
    Save the result to a file.
    
5. Broadcast transactions

    Broadcast offline signed transactions. Here you just use the transaction generated by above sign command. Of course, you can generate your signed transaction by any methods, eg. [irisnet-crypto](https://github.com/irisnet/irisnet-crypto).
    ```bash
    iriscli bank broadcast [file]
    ```
    The transaction will be broadcast and executed in IRISnet.
     