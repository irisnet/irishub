# iriscli service disable 

## Description

Disable a available service binding

## Usage

```
iriscli service disable [flags]
```

## Flags

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --def-chain-id        |                         | [string] the ID of the blockchain defined of the service                                                                                              |  Yes     |
| --service-name        |                         | [string] service name                                                                                                                                 |  Yes     |
| -h, --help            |                         | help for disable                                                                                                                                      |          |
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

### Disable a available service binding
```shell
iriscli service disable --chain-id=test  --from=node0 --fee=0.004iris --def-chain-id=test --service-name=test-service
```

After that, you're done with disabling a available service binding.

```txt
Password to sign with 'node0':
Committed at block 537 (tx hash: 66C5EE634339D168A07C073C6BF209D80081762EB8451974ABC33A41914A7158, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:3510 Tags:[{Key:[97 99 116 105 111 110] Value:[115 101 114 118 105 99 101 45 100 105 115 97 98 108 101] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 55 48 50 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "service-disable",
     "completeConsumedTxFee-iris-atto": "\"70200000000000\""
   }
 }
```