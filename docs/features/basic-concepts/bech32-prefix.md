# Bech32 on IRISnet

Bech32 is a new Bitcoin address format proposed by Pieter Wuille and Greg Maxwell. Besides Bitcoin addresses, Bech32 can encode any short binary data. In the IRIS network, keys and addresses may refer to a number of different roles in the network like accounts, validators etc. The IRIS network is designed to use the Bech32 address format to provide robust integrity checks on data. The human readable part(HRP) makes it more efficient to read and the users could see error messages. More details in [bip-0173](https://github.com/bitcoin/bips/blob/master/bip-0173.mediawiki)


## Human Readable Part Table


| HRP        | Definition |
| -----------|:-------------|
|iaa|   IRISnet Account Address|
|fap|    IRISnet Account Public Key|
|fva|   IRISnet Validator's Operator Address|
|fvp|   IRISnet Validator's Operator Public Key|
|fca|   Tendermint Consensus Address|
|fcp|    Tendermint Consensus Public Key|

## Encoding

Not all interfaces to users IRISnet should be exposed as bech32 interfaces. Many address are still in hex or base64 encoded form.

To covert between other binary reprsentation of addresses and keys, it is important to first apply the Amino enocoding process before bech32 encoding.


## Account Key Example

Once you create a new address, you should see the following:

```
NAME:    TYPE:           ADDRESS:                                PUBKEY:
test1    local    iaa18ekc4dswwrh2a6lfyev4tr25h5y76jkpqsz7kl    fap1addwnpepqgxa40ww28uy9q46gg48g6ulqdzwupyjcwfumgfjpvz7krmg5mrnw6zv8uv
```

This means you have created a new address `iaa18ekc4dswwrh2a6lfyev4tr25h5y76jkpqsz7kl`, its hrp is `iaa`. And its public key could be encoded into `fap1addwnpepqgxa40ww28uy9q46gg48g6ulqdzwupyjcwfumgfjpvz7krmg5mrnw6zv8uv`, its hrp is `fap`. 

## Validator Key Example

A Tendermint Consensus Public key is generated when the node is created with  `iris init`.
You can get this value with   
```
iris tendermint show-validator
```

Example output:
```
fcp1zcjduepqwh0tqpqrewe9lrr87ywgjq50gd3m82mgz0qwsmu62s83pukrqsfs5lv2kw
```
