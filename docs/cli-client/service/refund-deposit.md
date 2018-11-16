# iriscli service refund-deposit 

## Description

Refund all deposit from a service binding

## Usage

```
iriscli service refund-deposit [flags]
```

## Flags

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --def-chain-id        |                         | [string] the ID of the blockchain defined of the service                                                                                              |  Yes     |
| --service-name        |                         |  [string] service name                                                                                                                                |  Yes     |
| -h, --help            |                         |  help for refund-deposit                                                                                                                              |          |
| --account-number      |                         | [int] AccountNumber number to sign the tx                                                                                                             |          |
| --async               |                         | broadcast transactions asynchronously                                                                                                                 |          |
| --chain-id            |                         | [string] Chain ID of tendermint node                                                                                                                  |   Yes    |
| --dry-run             |                         | ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it                                                               |          |
| --fee                 |                         | [string] Fee to pay along with transaction                                                                                                            |   Yes    |
| --from                |                         | [string] Name of private key with which to sign                                                                                                       |   Yes    |
| --from-addr           |                         | [string] Specify from address in generate-only mode                                                                                                   |          |
| --gas                 |  200000                 | [string] gas limit to set per-transaction; set to "simulate" to calculate required gas automatically                                                  |          |
| --gas-adjustment      |  1                      | [float] adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  |          |
| --generate-only       |                         | build an unsigned transaction and write it to STDOUT                                                                                                  |          |
| --indent              |                         | Add indent to JSON response                                                                                                                           |          |
| --json                |                         | return output in json format                                                                                                                          |          |
| --ledger              |                         | Use a connected Ledger device                                                                                                                         |          |
| --memo                |                         | [string] Memo to send along with transaction                                                                                                          |          |
| --node                |  tcp://localhost:26657  | [string] <host>:<port> to tendermint rpc interface for this chain                                                                                     |          |
| --print-response      |                         | return tx response (only works with async = false)                                                                                                    |          |
| --sequence            |                         | [int] Sequence number to sign the tx                                                                                                                  |          |
| --trust-node          |  true                   | Don't verify proofs for responses                                                                                                                     |          |

## Examples

### Refund all deposit from an unavailable service binding
```shell
iriscli service refund-deposit --chain-id=test  --from=node0 --fee=0.004iris --def-chain-id=test --service-name=test-service
```

After that, you're done with refunding all deposit from a service binding.

```txt
Password to sign with 'node0':
Committed at block 991 (tx hash: 8A7F0EA61AB73A8B241945C8942EC8593774346B36BB70E36E138A23E7A473EE, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:4614 Tags:[{Key:[97 99 116 105 111 110] Value:[115 101 114 118 105 99 101 45 114 101 102 117 110 100 45 100 101 112 111 115 105 116] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 57 50 50 56 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "service-refund-deposit",
     "completeConsumedTxFee-iris-atto": "\"92280000000000\""
   }
 }
```