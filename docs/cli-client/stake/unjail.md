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
iriscli stake unjail --from=<key name> --fee=0.004iris --chain-id=<chain-id>
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
```json
{
   "tags": {
     "action": "unjail",
     "completeConsumedTxFee-iris-atto": "\"918600000000000\"",
     "validator": "fva12zgt9hc5r5mnxegam9evjspgwhkgn4wzjxkvqy"
   }
 }
```
