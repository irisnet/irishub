# IService User Guide

## Basic Function Description
IRIS Services (a.k.a. "iServices") intend to bridge the gap between the blockchain world and the conventional business application world, 
by mediating a complete lifecycle of off-chain services -- from their definition, binding (provider registration), invocation, to their 
governance (profiling and dispute resolution). By enhancing the IBC processing logic to support service semantics, the IRIS SDK is intended 
to allow distributed business services to be available across the internet of blockchains. We introduced the [Interface description language](https://en.wikipedia.org/wiki/Interface_description_language) 
(IDL) to work with the service standardized definitions to satisfy service invocations across different programming languages.
The currently supported IDL language is [protobuf](https://developers.google.com/protocol-buffers/). The main function points of this module are as follows:
1. Service Definition
2. Service Binding
3. Service Invocation (TODO)
4. Dispute Resolution (TODO)
5. Service Analysis (TODO)

## Interactive process

### Service definition process

1. Any users can define a service. In service definition，Use `protobuf` to standardize the definition of the service's method, its input and output parameters.

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
iriscli iservice define --chain-id=service-test  --from=x --fee=0.004iris --service-name=test-service --service-description=service-description --author-description=author-description --tags tag1,tag2 --messaging=Unicast --idl-content=<idl-content> --file=test.proto

# Result
Committed at block 92 (tx hash: A63241AA6666B8CFE6B1C092B707AB0FA350F108, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:8007 Tags:[{Key:[97 99 116 105 111 110] Value:[115 101 114 118 105 99 101 45 100 101 102 105 110 101]} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[49 54 48 49 52 48 48 48 48 48 48 48 48 48 48]}] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "service-define",
     "completeConsumedTxFee-iris-atto": "160140000000000"
   }
}

# Query service definition
iriscli iservice definition --def-chain-id=service-test --service-name=test-service

```

### Service Binding
```
# Service binding
iriscli iservice bind --chain-id=service-test  --from=x --fee=0.004iris --service-name=test-service --def-chain-id=service-test --bind-type=Local  --deposit=1iris --prices=1iris --avg-rsp-time=10000 --usable-time=100 --expiration=-1

# Result
Committed at block 168 (tx hash: 02CAC60E75CD93465CAE10CE35F30B53C8A95574, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:5437 Tags:[{Key:[97 99 116 105 111 110] Value:[115 101 114 118 105 99 101 45 98 105 110 100]} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[49 48 56 55 52 48 48 48 48 48 48 48 48 48 48]}] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "service-bind",
     "completeConsumedTxFee-iris-atto": "108740000000000"
   }
}

# Query service binding
iriscli iservice binding --def-chain-id=service-test --service-name=test-service --bind-chain-id=service-test --provider=<your address>

# Query service bindings
iriscli iservice bindings --def-chain-id=service-test --service-name=test-service

# Service binding update
iriscli iservice update-binding --chain-id=service-test  --from=x --fee=0.004iris --service-name=test-service --def-chain-id=service-test --bind-type=Local  --deposit=1iris --prices=1iris,2iris --avg-rsp-time=10000 --usable-time=100 --expiration=-1

# Result
Committed at block 233 (tx hash: 2F5F44BAF09981D137EA667F9E872EB098A9B619, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:4989 Tags:[{Key:[97 99 116 105 111 110] Value:[115 101 114 118 105 99 101 45 98 105 110 100 105 110 103 45 117 112 100 97 116 101]} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[57 57 55 56 48 48 48 48 48 48 48 48 48 48]}] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "service-binding-update",
     "completeConsumedTxFee-iris-atto": "99780000000000"
   }
}

# Refund Deposit
iriscli iservice refund-deposit --chain-id=service-test  --from=x --fee=0.004iris --def-chain-id=service-test --service-name=test-service

# Result
Committed at block 1563 (tx hash: 748CEA6EA9DEFB384FFCFBE68A3CB6D8B643361B, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:5116 Tags:[{Key:[97 99 116 105 111 110] Value:[115 101 114 118 105 99 101 45 114 101 102 117 110 100 45 100 101 112 111 115 105 116]} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[49 48 50 51 50 48 48 48 48 48 48 48 48 48 48]}] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "service-refund-deposit",
     "completeConsumedTxFee-iris-atto": "102320000000000"
   }
}
```

## CLI Command Details

```
iriscli iservice define --chain-id=service-test  --from=x --fee=0.004iris --service-name=test-service --service-description=service-description --author-description=author-description --tags "tag1 tag2" --messaging=Unicast --idl-content=<idl-content> --file=test.proto
```
* `--service-name`  The name of iService
* `--service-description`  The description of this iService
* `--author-description`  The self-description of the iService creator which is optional
* `--tags`  The keywords of this iService
* `--messaging`  The messaging type of this iService{`Unicast`,`Multicast`}
* `--idl-content`  The standardized definition of the methods for this iService
* `--file`  Idl-content can be replaced by files,if the item is not empty.

```
iriscli iservice definition --def-chain-id=service-test --service-name=test-service
```
* `--def-chain-id`  The ID of the blockchain defined of the iService
* `--service-name`  The name of iService

```
iriscli iservice bind --chain-id=service-test  --from=x --fee=0.004iris --service-name=test-service --def-chain-id=service-test --bind-type=Local  --deposit=1iris --prices=1iris --avg-rsp-time=10000 --usable-time=100 --expiration=-1
```
* `--def-chain-id` The ID of the blockchain defined of the iService
* `--service-name` The name of iService
* `--bind-type` Set whether the service is local or global
* `--deposit` The deposit of service provider
* `--prices` Service prices, a list sorted by service method
* `--avg-rsp-time` The average service response time in milliseconds
* `--usable-time` An integer represents the number of usable service invocations per 10,000
* `--expiration` The blockchain height where this binding expires; a negative number means "never expire"

```
iriscli iservice binding --def-chain-id=service-test --service-name=test-service --bind-chain-id=service-test --provider=<your address>
```
* `--def-chain-id` The ID of the blockchain defined of the iService
* `--service-name` The name of iService
* `--bind-chain-id`  The ID of the blockchain bound of the iService
* `--provider` The bech32 encoded account created the iService binding

```
iriscli iservice bindings --def-chain-id=service-test --service-name=test-service
```
* Refer to iriscli iservice binding

```
iriscli iservice update-binding --chain-id=service-test  --from=x --fee=0.004iris --service-name=test-service --def-chain-id=service-test --bind-type=Local  --deposit=1iris --prices=1iris,2iris --avg-rsp-time=10000 --usable-time=100 --expiration=-1
```
* Refer to iriscli iservice bind

```
iriscli iservice refund-deposit --chain-id=service-test  --from=x --fee=0.004iris --def-chain-id=service-test --service-name=test-service
```
* `--def-chain-id` The ID of the blockchain defined of the iService
* `--service-name` The name of iService

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
