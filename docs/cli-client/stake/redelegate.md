# iriscli stake redelegate

## Description

Redelegate illiquid tokens from one validator to another

## Usage

```
iriscli stake redelegate [flags]
```

## Flags

| Name, shorthand              | Default               | Description                                                         | Required |
| ---------------------------- | --------------------- | ------------------------------------------------------------------- | -------- |
| --account-number             |                       | [int] AccountNumber number to sign the tx                           |          |
| --address-validator-dest     |                       | [string] Bech address of the destination validator                  | Yes      |
| --address-validator-source   |                       | [string] Bech address of the source validator                       | Yes      |
| --async                      |                       | Broadcast transactions asynchronously                               |          |
| --chain-id                   |                       | [string] Chain ID of tendermint node                                | Yes      |
| --dry-run                    |                       | Ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it ||
| --fee                        |                       | [string] Fee to pay along with transaction                          | Yes      |
| --from                       |                       | [string] Name of private key with which to sign                     | Yes      |
| --from-addr                  |                       | [string] Specify from address in generate-only mode                 |          |
| --gas                        | 200000                | [string] Gas limit to set per-transaction; set to "simulate" to calculate required gas automatically ||
| --gas-adjustment             | 1                     | [float] Adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored ||
| --generate-only              |                       | build an unsigned transaction and write it to STDOUT                |          |
| --help, -h                   |                       | help for redelegate                                                 |          |
| --indent                     |                       | Add indent to JSON response                                         |          |
| --json                       |                       | return output in json format                                        |          |
| --ledger                     |                       | Use a connected Ledger device                                       |          |
| --memo                       |                       | [string] Memo to send along with transaction                                 |          |
| --node                       | tcp://localhost:26657 | [string] \<host>:\<port> to tendermint rpc interface for this chain          |          |
| --print-response             |                       | return tx response (only works with async = false)                  |          |
| --sequence                   |                       | [int] Sequence number to sign the tx                                      |          |
| --shares-amount              |                       | [string] Amount of source-shares to either unbond or redelegate as a positive integer or decimal ||
| --shares-percent             |                       | [string] Percent of source-shares to either unbond or redelegate as a positive integer or decimal >0 and <=1 ||
| --trust-node                 | true                  | Don't verify proofs for responses                                   |          |

## Examples

### Redelegate illiquid tokens from one validator to another

```shell
iriscli stake redelegate --chain-id=ChainID --from=KeyName --fee=Fee --address-validator-source=SourceValidatorAddress --address-validator-dest=DestinationValidatorAddress --shares-percent=SharesPercent
```

After that, you're done with redelegating specified liquid tokens from one validator to another validator.

```txt
Committed at block 648 (tx hash: E59EE3C8F04D62DA0F5CFD89AC96402A92A56728692AEA47E8A126CDDA58E44B, response: {Code:0 Data:[11 8 185 204 185 223 5 16 247 169 147 42] Log:Msg 0:  Info: GasWanted:200000 GasUsed:29085 Tags:[{Key:[97 99 116 105 111 110] Value:[98 101 103 105 110 45 114 101 100 101 108 101 103 97 116 105 111 110] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[100 101 108 101 103 97 116 111 114] Value:[102 97 97 49 48 115 48 97 114 113 57 107 104 112 108 48 99 102 122 110 103 51 113 103 120 99 120 113 48 110 121 54 104 109 99 57 115 121 116 106 102 107] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[115 111 117 114 99 101 45 118 97 108 105 100 97 116 111 114] Value:[102 118 97 49 100 97 121 117 106 100 102 110 120 106 103 103 100 53 121 100 108 118 118 103 107 101 114 112 50 115 117 112 107 110 116 104 97 106 112 99 104 50] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[100 101 115 116 105 110 97 116 105 111 110 45 118 97 108 105 100 97 116 111 114] Value:[102 118 97 49 104 50 55 120 100 119 54 116 57 108 53 106 103 118 117 110 55 54 113 100 117 52 53 107 103 114 120 57 108 113 101 100 101 56 104 112 99 100] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[101 110 100 45 116 105 109 101] Value:[11 8 185 204 185 223 5 16 247 169 147 42] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 53 56 49 55 48 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "begin-redelegation",
     "completeConsumedTxFee-iris-atto": "\"5817000000000000\"",
     "delegator": "faa10s0arq9khpl0cfzng3qgxcxq0ny6hmc9sytjfk",
     "destination-validator": "fva1h27xdw6t9l5jgvun76qdu45kgrx9lqede8hpcd",
     "end-time": "\u000b\u0008\ufffd̹\ufffd\u0005\u0010\ufffd\ufffd\ufffd*",
     "source-validator": "fva1dayujdfnxjggd5ydlvvgkerp2supknthajpch2"
   }
}
```
