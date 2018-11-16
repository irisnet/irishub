# iriscli service definition

## Description

Query service definition

## Usage

```
iriscli service definition [flags]
```

## Flags

| Name, shorthand | Default                    | Description                                                         | Required |
| --------------- | -------------------------- | ------------------------------------------------------------------- | -------- |
| --def-chain-id  |                            | [string] the ID of the blockchain defined of the service            | Yes      |
| --service-name  |                            | [string] service name                                               | Yes      |
| --help, -h      |                            | help for definition                                                 |          |
| --chain-id      |                            | [string] Chain ID of tendermint node                                |          |
| --height        | most recent provable block | [int] block height to query                                         |          |
| --indent        |                            | Add indent to JSON response                                         |          |
| --ledger        |                            | Use a connected Ledger device                                       |          |
| --node          | tcp://localhost:26657      | [string] \<host>:\<port> to tendermint rpc interface for this chain |          |
| --trust-node    | true                       | Don't verify proofs for responses                                   |          |


## Examples

### Query a service definition

```shell
iriscli service definition --def-chain-id=test --service-name=test-service
```

After that, you will get detail info for the service definition which has the specfied define chain id and service name.

```json
{
  "SvcDef": {
    "name": "test-service",
    "chain_id": "test",
    "description": "service-description",
    "tags": [
      "tag1",
      "tag2"
    ],
    "author": "faa1ydhmma8l4m9dygsh7l08fgrwka6yczs0gkfnvd",
    "author_description": "author-description",
    "idl_content": "syntax = \"proto3\";\n\npackage helloworld;\n\n// The greeting service definition.\nservice Greeter {\n    //@Attribute description: sayHello\n    //@Attribute output_privacy: NoPrivacy\n    //@Attribute output_cached: NoCached\n    rpc SayHello (HelloRequest) returns (HelloReply) {}\n}\n\n// The request message containing the user's name.\nmessage HelloRequest {\n    string name = 1;\n}\n\n// The response message containing the greetings\nmessage HelloReply {\n    string message = 1;\n}\n"
  },
  "methods": [
    {
      "id": "1",
      "name": "SayHello",
      "description": "sayHello",
      "output_privacy": "NoPrivacy",
      "output_cached": "NoCached"
    }
  ]
}
```