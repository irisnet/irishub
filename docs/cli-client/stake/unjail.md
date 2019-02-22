# iriscli stake unjail

## Introduction

In Proof-of-Stake blockchain, validators will get block provisions by staking their token. But if they failed to keep online, they will be punished by slashing a small portion of their staked tokens. The offline validators will be removed from the validator set and put into jail, which means their voting power is zero. During the jail period, these nodes are not event validator candidates . Once the jail period is end, they can send `unjail` transactions to free themselves and become validator candidates again.

## Usage

```
iriscli stake unjail [flags]
```

Print help messages:
```
iriscli stake unjail --help
```

## Examples

### Unjail validator previously jailed for downtime

```
iriscli stake unjail --from=<key name> --fee=0.3iris --chain-id=<chain-id>
```
### Common Issue

* Check the jailing time for this validator:

```$xslt
iriscli stake signing-info fvp1zcjduepqewwc93xwvt0ym6prxx9ppfzeufs33flkcpu23n5eutjgnnqmgazsw54sfv --node=localhost:36657 --trust-node
```
If your validator is jailed, it will tell the jailing time.

```
Start height: 565, index offset: 2, jailed until: 2018-12-12 06:46:37.274910287 +0000 UTC, missed blocks counter: 2
```

If you do `unjail` before the jailing time, you will see the following error.

```$xslt
ERROR: Msg 0 failed: {"codespace":10,"code":102,"abci_code":655462,"message":"validator still jailed, cannot yet be unjailed"}
```

After that jailing period, you could submit an `unjail` transaction. 

Sample output:
```txt
Committed at block 306 (tx hash: 5A4C6E00F4F6BF795EB05D2D388CBA0E8A6E6CF17669314B1EE6A31729A22450, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:3398 Tags:[{Key:[97 99 116 105 111 110] Value:[115 101 114 118 105 99 101 45 119 105 116 104 100 114 97 119 45 102 101 101 115] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 54 55 57 54 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
```
```json
{
   "tags": {
     "action": "unjail",
     "completeConsumedTxFee-iris-atto": "\"918600000000000\"",
     "validator": "fva12zgt9hc5r5mnxegam9evjspgwhkgn4wzjxkvqy"
   }
 }
```
