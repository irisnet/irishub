# iriscli record submit

## Description

Submit a record on the chain

## Usage

```
iriscli record submit [flags]
```

## Flags

| Name, shorthand  | Default                    | Description                                                                                 | Required |
| ---------------  | -------------------------- | ------------------------------------------------------------------------------------------- | -------- |
| --account-number |                            | [int] AccountNumber number to sign the tx                                                   |          |
| --async          |                            | Broadcast transactions asynchronously                                                       |          |
| --chain-id       |                            | [string] Chain ID of tendermint node                                                        | Yes      |
| --description    | description                | [string] Uploaded file description                                                          |          |
| --dry-run        |                            | Ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it     |          |
| --fee            |                            | [string] Fee to pay along with transaction                                                  | Yes      |
| --from           |                            | [string] Name of private key with which to sign                                             | Yes      |
| --from-addr      |                            | [string] Specify from address in generate-only mode                                         |          |
| --gas string     | 200000                     | Gas limit to set per-transaction; set to "simulate" to calculate required gas automatically |          |
| --gas-adjustment | 1                          | [float] Adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored |          |
| --generate-only  |                            | Build an unsigned transaction and write it to STDOUT                                        |          |
| -h, --help       |                            | help for submit                                                                             |          |
| --indent         |                            | Add indent to JSON response                                                                 |          |
| --json           |                            | return output in json format                                                                |          |
| --ledger         |                            | Use a connected Ledger device                                                               |          |
| --memo           |                            | [string] Memo to send along with transaction                                                |          |
| --node           | tcp://localhost:26657      | [string] \<host>:\<port> to tendermint rpc interface for this chain                           |          |
| --onchain-data   |                            | [string] on chain data source                                                               | Yes      |
| --print-response |                            | return tx response (only works with async = false)                                          |          |
| --sequence       |                            | [int] Sequence number to sign the tx                                                        |          |
| --trust-node     | true                       | Don't verify proofs for responses                                                           |          |

## Examples

### Submit a record

```shell
iriscli record submit --chain-id="test" --onchain-data="this is my on chain data" --from=node0 --fee=0.1iris
```

After that, you're done with submitting a new record, but remember to back up your record id, it's the only way to retrieve your record.

```txt
Committed at block 486 (tx hash: 8AB91BF0E61AD2C860402B88579EE83167506E7C3A8597E873976915D82D4F1B, response:
 {
   "code": 0,
   "data": "cmVjb3JkOmFiNTYwMmJhYzEzZjExNzM3ZTg3OThkZDU3ODY5YzQ2ODE5NGVmYWQyZGIzNzYyNTc5NWYxZWZkOGQ5ZDYzYzY=",
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 3764,
   "codespace": "",
   "tags": {
     "action": "submit_record",
     "ownerAddress": "faa1x25y3ltr4jvp89upymegvfx7n0uduz5kmh5xuz",
     "record-id": "record:ab5602bac13f11737e8798dd57869c468194efad2db37625795f1efd8d9d63c6"
   }
 })
```

