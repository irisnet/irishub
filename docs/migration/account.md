## Introduction

IRISHub v2.0 changed the generative structure of the entire account system from `auth.BaseAccountProto` to `ethermint.ProtoAccount` , resulting in incompatibility of newly created accounts;

**Note:** **Accounts before the v2.0 upgrade will not be affected and can still be used normally**

**The following content takes nyancat testnet as an example：**

For accounts before the v2.0 upgrade, the query structure through LCD is as follows: https://lcd.nyancat.irisnet.org/cosmos/auth/v1beta1/accounts/iaa1e0rx87mdj79zejewuc4jg7ql9ud2286g2us8f2

```json
{
     "account": {
         "@type": "/cosmos.auth.v1beta1.BaseAccount",
         "address":"iaa1e0rx87mdj79zejewuc4jg7ql9ud2286g2us8f2",
         "pub_key": {
             "@type": "/cosmos.crypto.secp256k1.PubKey",
             "key":"AiOFJ3Jclq/8y3xV85ALNFuA7FJo1IMoTxYoB3ddMrMr"
         },
         "account_number": "1251",
         "sequence": "12983"
     }
}
```

For accounts upgraded in v2.0, the query structure through LCD is as follows: https://lcd.nyancat.irisnet.org/cosmos/auth/v1beta1/accounts/iaa1g4uak38a8fhkg5v5qky3fc9g6h50yrdcn7waug

```json
{
     "account": {
         "@type": "/ethermint.types.v1.EthAccount",
         "base_account": {
             "address":"iaa1g4uak38a8fhkg5v5qky3fc9g6h50yrdcn7waug",
             "pub_key": {
                 "@type": "/ethermint.crypto.v1.ethsecp256k1.PubKey",
                 "key":"AhHKT0xpnrOmpikkd1lEPxiEHG4ngItq06KLhwU2UQHO"
             },
             "account_number": "5564",
             "sequence": "6"
         },
         "code_hash": "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"
     }
}
```


### Using Proto

A new type needs to be registered, and the type structure is here: https://github.com/bianjieai/ethermint/blob/v0.20.0-irishub-1/proto/ethermint/types/v1/account.proto

When the query structure is `/ethermint.types.v1.EthAccount`, use the new structure analysis;

When the query structure is `/cosmos.auth.v1beta1.BaseAccount`, use the old structure analysis;

You can refer to the core-sdk-go written by the irisnet team. The relevant hash is at: https://github.com/irisnet/core-sdk-go/commit/68ed671727e057edb185935c42710f8777dab62f

### Useing LCD

When parsing, please make compatibility according to the required type;

If the type is `/cosmos.auth.v1beta1.BaseAccount`, use the old structure analysis;

If the type is `/ethermint.types.v1.EthAccount`, use the new structure analysis;