# iriscli service define 

## Description

Create a new service definition

## Usage

```
iriscli service define [flags]
```

## Flags

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --service-description |                         | [string] service description                                                                                                                          |          |
| --author-description  |                         | [string] service author description                                                                                                                   |          |
| --service-name        |                         | [string] service name                                                                                                                                 |   Yes    |
| --tags                |                         | [strings] service tags                                                                                                                                |          |
| --idl-content         |                         | [string] content of service interface description language                                                                                            |          |
| --file                |                         | [string] path of file which contains service interface description language                                                                           |          |
| -h, --help            |                         | help for define                                                                                                                                       |          |
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

### define a service
```shell
iriscli service define --chain-id=test  --from=node0 --fee=0.004iris --service-name=test-service --service-description=service-description --author-description=author-description --tags=tag1,tag2 --idl-content=<idl-content> --file=test.proto
```
Idl-content can be replaced by file if the file item is not empty.  [Example of IDL content](#idl-content-example).

After that, you're done with defining a new service.

```txt
Password to sign with 'node0':
Committed at block 65 (tx hash: 663B676E453F91BFCDC87B0308910501DD14DF79C88390FC15E06C4CC9612422, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:7968 Tags:[{Key:[97 99 116 105 111 110] Value:[115 101 114 118 105 99 101 45 100 101 102 105 110 101] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 49 53 57 51 54 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "service-define",
     "completeConsumedTxFee-iris-atto": "\"159360000000000\""
   }
 }
```

### IDL content example
* idl-content example

    > syntax = \\"proto3\\";\n\npackage helloworld;\n\n// The greeting service definition.\nservice Greeter {\n    //@Attribute description: sayHello\n    //@Attribute output_privacy: NoPrivacy\n    //@Attribute output_cached: NoCached\n    rpc SayHello (HelloRequest) returns (HelloReply) {}\n}\n\n// The request message containing the user's name.\nmessage HelloRequest {\n    string name = 1;\n}\n\n// The response message containing the greetings\nmessage HelloReply {\n    string message = 1;\n}\n

* IDL file example

    [test.proto](../../features/test.proto)

