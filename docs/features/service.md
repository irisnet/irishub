# Service User Guide

## Basic Function Description
IRIS Services (a.k.a. "iServices") intend to bridge the gap between the blockchain world and the conventional business application world, by mediating a complete lifecycle of off-chain services -- from their definition, binding (provider registration), invocation, to their governance (profiling and dispute resolution). By enhancing the IBC processing logic to support service semantics, the IRIS SDK is intended to allow distributed business services to be available across the internet of blockchains. The [Interface description language](https://en.wikipedia.org/wiki/Interface_description_language) (IDL) we introduced is
to work with the service standardized definitions to satisfy service invocations across different programming languages.
The currently supported IDL language is [protobuf](https://developers.google.com/protocol-buffers/). The main functions of this module are as follows:
1. Service Definition
2. Service Binding
3. Service Invocation
4. Dispute Resolution (TODO)
5. Service Analysis (TODO)

## Interactive process

### Service definition process

1. Any users can define a service. In service definition，use `protobuf` to standardize the definition of the service's method, its input and output parameters.

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
iriscli service define --chain-id=service-test  --from=x --fee=0.004iris --service-name=test-service --service-description=service-description --author-description=author-description --tags=tag1,tag2 --idl-content=<idl-content> --file=test.proto

# Result
Committed at block 92 (tx hash: A63241AA6666B8CFE6B1C092B707AB0FA350F108, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:8007 Tags:[{Key:[97 99 116 105 111 110] Value:[115 101 114 118 105 99 101 45 100 101 102 105 110 101]} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[49 54 48 49 52 48 48 48 48 48 48 48 48 48 48]}] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "service-define",
     "completeConsumedTxFee-iris-atto": "160140000000000"
   }
}

# Query service definition
iriscli service definition --def-chain-id=service-test --service-name=test-service

```

### Service Binding
```
# Service Binding
iriscli service bind --chain-id=service-test  --from=x --fee=0.004iris --service-name=test-service --def-chain-id=service-test --bind-type=Local  --deposit=1iris --prices=1iris --avg-rsp-time=10000 --usable-time=100

# Result
Committed at block 168 (tx hash: 02CAC60E75CD93465CAE10CE35F30B53C8A95574, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:5437 Tags:[{Key:[97 99 116 105 111 110] Value:[115 101 114 118 105 99 101 45 98 105 110 100]} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[49 48 56 55 52 48 48 48 48 48 48 48 48 48 48]}] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "service-bind",
     "completeConsumedTxFee-iris-atto": "108740000000000"
   }
}

# Query service binding
iriscli service binding --def-chain-id=service-test --service-name=test-service --bind-chain-id=service-test --provider=<your address>

# Query service binding list
iriscli service bindings --def-chain-id=service-test --service-name=test-service

# Service binding update
iriscli service update-binding --chain-id=service-test  --from=x --fee=0.004iris --service-name=test-service --def-chain-id=service-test --bind-type=Local  --deposit=1iris --prices=1iris,2iris --avg-rsp-time=10000 --usable-time=100

# Result
Committed at block 233 (tx hash: 2F5F44BAF09981D137EA667F9E872EB098A9B619, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:4989 Tags:[{Key:[97 99 116 105 111 110] Value:[115 101 114 118 105 99 101 45 98 105 110 100 105 110 103 45 117 112 100 97 116 101]} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[57 57 55 56 48 48 48 48 48 48 48 48 48 48]}] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "service-binding-update",
     "completeConsumedTxFee-iris-atto": "99780000000000"
   }
}

# Disable service binding
iriscli service disable --chain-id=service-test  --from=x --fee=0.004iris --def-chain-id=service-test --service-name=test-service

# Result
Committed at block 241 (tx hash: 0EF936E1228F9838D0343D0FB3613F5E938602B7, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:4861 Tags:[{Key:[97 99 116 105 111 110] Value:[115 101 114 118 105 99 101 45 100 105 115 97 98 108 101] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 57 55 50 50 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "service-disable",
     "completeConsumedTxFee-iris-atto": "\"97220000000000\""
   }
}

# Enable service binding
iriscli service enable --chain-id=service-test  --from=x --fee=0.004iris --def-chain-id=service-test --service-name=test-service --deposit=1iris

# Result
Committed at block 176 (tx hash: 74AE647B8A311501CA82DACE90AA28CDB4695803, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:6330 Tags:[{Key:[97 99 116 105 111 110] Value:[115 101 114 118 105 99 101 45 101 110 97 98 108 101] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 49 50 54 54 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "service-enable",
     "completeConsumedTxFee-iris-atto": "\"126600000000000\""
   }
}

# Refund Deposit
iriscli service refund-deposit --chain-id=service-test  --from=x --fee=0.004iris --def-chain-id=service-test --service-name=test-service

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
iriscli service define --chain-id=service-test  --from=x --fee=0.004iris --service-name=test-service --service-description=service-description --author-description=author-description --tags=tag1,tag2 --idl-content=<idl-content> --file=test.proto
```
* `--service-name`  The name of service
* `--service-description`  The description of this service
* `--author-description`  The self-description of the service creator which is optional
* `--tags`  The keywords of this service
* `--idl-content`  The standardized definition of the methods for this service
* `--file`  Idl-content can be replaced by files,if the item is not empty.

```
iriscli service definition --def-chain-id=service-test --service-name=test-service
```
* `--def-chain-id`  The ID of the blockchain defined of the service
* `--service-name`  The name of service

```
iriscli service bind --chain-id=service-test  --from=x --fee=0.004iris --service-name=test-service --def-chain-id=service-test --bind-type=Local  --deposit=1iris --prices=1iris --avg-rsp-time=10000 --usable-time=100
```
* `--def-chain-id` The ID of the blockchain defined of the service
* `--service-name` The name of service
* `--bind-type` Set whether the service is local or global
* `--deposit` The deposit of service provider
* `--prices` Service prices, a list sorted by service method
* `--avg-rsp-time` The average service response time in milliseconds
* `--usable-time` An integer represents the number of usable service invocations per 10,000

```
iriscli service binding --def-chain-id=service-test --service-name=test-service --bind-chain-id=service-test --provider=<your address>
```
* `--def-chain-id` The ID of the blockchain defined of the service
* `--service-name` The name of service
* `--bind-chain-id`  The ID of the blockchain bound of the service
* `--provider` The blockchain address of bech32 encoded account 

```
iriscli service bindings --def-chain-id=service-test --service-name=test-service
```
* Refer to iriscli service binding

```
iriscli service update-binding --chain-id=service-test  --from=x --fee=0.004iris --service-name=test-service --def-chain-id=service-test --bind-type=Local  --deposit=1iris --prices=1iris,2iris --avg-rsp-time=10000 --usable-time=100
```
* `--def-chain-id` The ID of the blockchain defined of the service
* `--service-name` The name of service
* `--bind-type` Set whether the service is local or global
* `--deposit` Add to the current deposit balance of service provider
* `--prices` Service prices, a list sorted by service method
* `--avg-rsp-time` The average service response time in milliseconds
* `--usable-time` An integer represents the number of usable service invocations per 10,000

```
iriscli service disable --chain-id=service-test  --from=x --fee=0.004iris --def-chain-id=service-test --service-name=test-service
```
* `--def-chain-id` The ID of the blockchain defined of the service
* `--service-name` The name of service

```
iriscli service enable --chain-id=service-test  --from=x --fee=0.004iris --def-chain-id=service-test --service-name=test-service --deposit=1iris
```
* `--def-chain-id` The ID of the blockchain defined of the service
* `--service-name` The name of service
* `--deposit` Add to the current deposit balance of service provider

```
iriscli service refund-deposit --chain-id=service-test  --from=x --fee=0.004iris --def-chain-id=service-test --service-name=test-service
```
* `--def-chain-id` The ID of the blockchain defined of the service
* `--service-name` The name of service

## IDL extension
When using proto file to standardize the definition of the service's method, its input and output parameters, the method attributes can be added through annotations.

### Annotation standard
* If `//@Attribute attribute： value` wrote on top of `rpc method`，it will be added to the method attributes. Eg.
> //@Attribute description: sayHello

### Currently supported attributes
* `description` The name of this method in the service
* `output_privacy` Whether the output of the method is encrypted，{`NoPrivacy`,`PubKeyEncryption`}
* `output_cached` Whether the output of the method is cached，{`OffChainCached`，`NoCached`}

### IDL content example
* idl-content example

    > syntax = \"proto3\";\n\npackage helloworld;\n\n// The greeting service definition.\nservice Greeter {\n    //@Attribute description: sayHello\n    //@Attribute output_privacy: NoPrivacy\n    //@Attribute output_cached: NoCached\n    rpc SayHello (HelloRequest) returns (HelloReply) {}\n}\n\n// The request message containing the user's name.\nmessage HelloRequest {\n    string name = 1;\n}\n\n// The response message containing the greetings\nmessage HelloReply {\n    string message = 1;\n}\n

* IDL file example

    [test.proto](https://github.com/irisnet/irishub/blob/master/docs/features/test.proto)
