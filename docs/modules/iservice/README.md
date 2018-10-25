# IService User Guide

## Basic Function Description

1. Service Definition
2. Service Binding (TODO)
3. Service Invocation (TODO)
4. Dispute Resolution (TODO)
5. Service Analysis (TODO)

## Interactive process

### Service definition process

1. Any users can define a service. In service definition，Use [protobuf](https://developers.google.com/protocol-buffers/) to standardize the definition of the service's method, its input and output parameters.

## Usage Scenario
### Create an environment

```
rm -rf iris
rm -rf .iriscli
iris init gen-tx --name=x --home=iris
iris init --gen-txs --chain-id=service-test -o --home=iris
iris start --home=iris
```

### Service Definition

```
# Service definition
iriscli iservice define --chain-id=service-test  --from=x --fee=0.004iris --name=test-service --service-description=service-description --author-description=author-description --tags "tag1 tag2" --broadcast=Broadcast --idl-content=<idl-content> --file=test.proto

# Result
Committed at block 1040 (tx hash: 58FD40B739F592F5BD9B904A661B8D7B19C63FA9, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:13601 Tags:[{Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[247 102 151 120 200 0]}] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "completeConsumedTxFee-iris-atto": "159740000000000"
   }
}

# Query service definition
iriscli iservice definition --name=test-service --chain-id=service-test

```

## CLI Command Details

```
iriscli iservice define --chain-id=service-test  --from=x --fee=0.004iris --name=test-service --service-description=service-description --author-description=author-description --tags "tag1 tag2" --broadcast=Broadcast --idl-content=<idl-content> --file=test.proto
```

* `--name`  The name of iService
* `--service-description`  The description of this iService
* `--author-description`  The self-description of the iService creator which is optional
* `--tags`  The keywords of this iService
* `--broadcast`  The Broadcast type of this iService{`Broadcast`,`Unicast`}
* `--idl-content`  The standardized definition of the methods for this iService
* `--file`  Idl-content can be replaced by files,if the item is not empty.

```
iriscli iservice definition --name=test-service --chain-id=service-test
```

* `--chain-id`  The ID of the blockchain defined of the iService
* `--name`  The name of iService

## IDL extension
When using proto file to standardize the definition of the service's method, its input and output parameters, the method attributes can be added through annotations.

### Annotation standard
* If `//@Attribute attribute： value` wrote on top of `rpc method`，it will be added to the method attributes. Eg.
> //@Attribute description: sayHello

### Currently supported attributes
* `description` The name of this method in the iService
* `output_privacy` Whether the output of the method is encrypted，{`NoPrivacy`,`PubKeyEncryption`}
* `output_cached` Whether the output of the method is cached，{`OffChainCached`，`NoCached`}

### IDL content example
* idl-content example
> syntax = \"proto3\";\n\npackage helloworld;\n\n// The greeting service definition.\nservice Greeter {\n    //@Attribute description: sayHello\n    //@Attribute output_privacy: NoPrivacy\n    //@Attribute output_cached: NoCached\n    rpc SayHello (HelloRequest) returns (HelloReply) {}\n}\n\n// The request message containing the user's name.\nmessage HelloRequest {\n    string name = 1;\n}\n\n// The response message containing the greetings\nmessage HelloReply {\n    string message = 1;\n}\n

* IDL file example
[test.proto](./test.proto)
