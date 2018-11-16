# iriscli stake delegate

## Description

Delegate liquid tokens to an validator

## Usage

```
iriscli stake delegate [flags]
```

## Flags

| Name, shorthand              | Default               | Description                                                         | Required |
| ---------------------------- | --------------------- | ------------------------------------------------------------------- | -------- |
| --account-number             |                       | [int] AccountNumber number to sign the tx                           |          |
| --address-validator          |                       | [string] Bech address of the validator                                       | Yes      |
| --amount                     |                       | [string] Amount of coins to bond                                             | Yes      |
| --async                      |                       | broadcast transactions asynchronously                               |          |
| --chain-id                   |                       | [string] Chain ID of tendermint node                                         | Yes      |
| --dry-run                    |                       | Ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it | |
| --fee                        |                       | [string] Fee to pay along with transaction                                   | Yes      |
| --from                       |                       | [string] Name of private key with which to sign                              | Yes      |
| --from-addr                  |                       | [string] Specify from address in generate-only mode                          |          |
| --gas                        | 200000                | [string] Gas limit to set per-transaction; set to "simulate" to calculate required gas automatically |          |
| --gas-adjustment             | 1                     | [float] Adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set |manually this flag is ignored |          |
| --generate-only              |                       | Build an unsigned transaction and write it to STDOUT                |          |
| --help, -h                   |                       | Help for delegate                                                   |          |
| --indent                     |                       | Add indent to JSON response                                         |          |
| --json                       |                       | Return output in json format                                        |          |
| --ledger                     |                       | Use a connected Ledger device                                       |          |
| --memo                       |                       | [string] Memo to send along with transaction                                 |          |
| --node                       | tcp://localhost:26657 | [string] \<host>:\<port> to tendermint rpc interface for this chain          |          |
| --print-response             |                       | return tx response (only works with async = false)                  |          |
| --sequence int               |                       | Sequence number to sign the tx                                      |          |
| --trust-node                 | true                  | Don't verify proofs for responses                                   |          |

## Examples

### Delegate liquid tokens to an validator

```shell
iriscli stake delegate --chain-id=ChainID --from=KeyName --fee=Fee --amount=CoinstoBond --address-validator=ValidatorOwnerAddress
```

After that, you're done with delegating liquid tokens to specified validator.

```txt
Committed at block 2352 (tx hash: 2069F0453619637EE67EFB0CFC53713AF28A0BB89137DEB4574D8B13E723999B, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:15993 Tags:[{Key:[97 99 116 105 111 110] Value:[100 101 108 101 103 97 116 101] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[100 101 108 101 103 97 116 111 114] Value:[102 97 97 49 51 108 99 119 110 120 112 121 110 50 101 97 51 115 107 122 109 101 107 54 52 118 118 110 112 57 55 106 115 107 56 113 109 104 108 54 118 120] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[100 101 115 116 105 110 97 116 105 111 110 45 118 97 108 105 100 97 116 111 114] Value:[102 118 97 49 53 103 114 118 51 120 103 51 101 107 120 104 57 120 114 102 55 57 122 100 48 119 48 55 55 107 114 103 118 53 120 102 54 100 54 116 104 100] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 55 57 57 54 53 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "delegate",
     "completeConsumedTxFee-iris-atto": "\"7996500000000000\"",
     "delegator": "faa13lcwnxpyn2ea3skzmek64vvnp97jsk8qmhl6vx",
     "destination-validator": "fva15grv3xg3ekxh9xrf79zd0w077krgv5xf6d6thd"
   }
 }
```
