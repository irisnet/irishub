# Service

## Basic Function Description

IRIS Services (a.k.a. "iService") intend to bridge the gap between the blockchain world and the conventional business application world, by mediating a complete lifecycle of off-chain services -- from their definition, binding (provider registration), invocation, to their governance (profiling and dispute resolution). By enhancing the IBC processing logic to support service semantics, the IRIS SDK is intended to allow distributed business services to be available across the internet of blockchains. The [Interface description language](https://en.wikipedia.org/wiki/Interface_description_language) (IDL) we introduced is
to work with the service standardized definitions to satisfy service invocations across different programming languages.
The currently supported IDL language is [protobuf](https://developers.google.com/protocol-buffers/). The main functions of this module are as follows:

1. Service Definition
2. Service Binding
3. Service Invocation
4. Dispute Resolution (TODO)
5. Service Analysis (TODO)

### System parameters

The following parameters can be modified by [governance](governance.md)

* `MinDepositMultiple`    a multiple of the minimum deposit amount of service binding
* `MaxRequestTimeout`     maximum number of waiting blocks for service invocation
* `ServiceFeeTax`         tax rate of service fee
* `SlashFraction`         slash fraction
* `ComplaintRetrospect`   maximum time for submit a dispute
* `ArbitrationTimeLimit`  maximum time of dispute resolution

## Interactive process

### Service definition

Any users can define a service. In service definition, use `protobuf` to standardize the definition of the service's method, its input and output parameters. In order to support attributes of iService better, IRISnet has made some extensions to `protobuf`, please refer to [IDL extension](#idl-extension) for details.

```bash
# create a new service definition
iriscli service define --chain-id=<chain-id>  --from=<key-name> --fee=0.6iris --gas=100000 --service-name=<service-name> --service-description=<service-description> --author-description=<author-description> --tags=<tag1>,<tag2> --idl-content=<idl-content> --file=</***/***.proto>

# query service definition
iriscli service definition --def-chain-id=<def-chain-id> --service-name=<service-name>
```

### Service Binding

The minimum deposit amount for Service Binding is `MinDepositMultiple * Service fee`. The service provider can update his service binding and adjust the Service fee at any time, disable and enable the service binding. If the provider want to refund the deposit, he needs to disable service binding first and wait for a period of `ComplaintRetrospectParameter` + `ArbitrationTimelimitParameter`.

```bash
# create a new service binding
iriscli service bind --chain-id=<chain-id>  --from=<key-name> --fee=0.3iris --service-name=<service-name> --def-chain-id=<def-chain-id> --bind-type=Local  --deposit=1000iris --prices=1iris --avg-rsp-time=10000 --usable-time=9999

# query service binding
iriscli service binding --def-chain-id=<def-chain-id> --service-name=<service-name> --bind-chain-id=<bind-chain-id> --provider=<provider-account-address>

# query service bindings
iriscli service bindings --def-chain-id=<def-chain-id> --service-name=<service-name>

# update a service binding
iriscli service update-binding --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --service-name=<service-name> --def-chain-id=<def-chain-id> --bind-type=Local  --deposit=1iris --prices=1iris,2iris --avg-rsp-time=10000 --usable-time=100

# disable a available service binding
iriscli service disable --chain-id=<chain-id>  --from=<key-name> --fee=0.3iris --def-chain-id=<def-chain-id> --service-name=<service-name>

# enable an unavailable service binding
iriscli service enable --chain-id=<chain-id>  --from=<key-name> --fee=0.3iris --def-chain-id=<def-chain-id> --service-name=<service-name> --deposit=100iris

# refund all deposit from a service binding
iriscli service refund-deposit --chain-id=<chain-id>  --from=<key-name> --fee=0.3iris --def-chain-id=<def-chain-id> --service-name=<service-name>
```

### Service Invocation

If the service consumer needs to initiate a service invocation request, the service fee specified by the service provider needs to be paid. The service provider needs to respond to the service request within the block height defined by `MaxRequestTimeout`. If the service provider does not respond in time, the deposit of the 'SlashFraction' ratio will be deducted from the service provider's service binding deposit and the service fee will be refunded to the service consumer's return pool. If the service call is responded normally, the system will deduct the `ServiceFeeTax` ratio from the service fee, and add the remaining service fee to the service provider's incoming pool. The service provider/consumer can initiate the `withdraw-fees`/`refund-fees` transaction to retrieve all of the tokens in the incoming/return pool.

```bash
# initiate service invocation
iriscli service call --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --def-chain-id=<def-chain-id> --service-name=<service-name> --method-id=1 --bind-chain-id=<bind-chain-id> --provider=<provider-account-address> --service-fee=1iris --request-data=<request-data>

# query service requests
iriscli service requests --def-chain-id=<def-chain-id> --service-name=<service-name> --bind-chain-id=<bind-chain-id> --provider=<provider-account-address>

# respond a service invocation
iriscli service respond --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --request-chain-id=<request-chain-id> --request-id=<request-id (e.g.230-130-0)> --response-data=<response-data>

# query a service response
iriscli service response --request-chain-id=<request-chain-id> --request-id=<request-id (e.g.230-130-0)>

# query return and incoming fee of a particular address
iriscli service fees <account-address>

# refund all fees from service return fees
iriscli service refund-fees --chain-id=<chain-id> --from=<key-name> --fee=0.3iris

# withdraw all fees from service incoming fees
iriscli service withdraw-fees --chain-id=<chain-id> --from=<key-name> --fee=0.3iris
```

## IDL extension

When using proto file to standardize the definition of the service's method, its input and output parameters, the method attributes can be added through annotations.

### Annotation standard

* If `//@Attribute attribute:  value` wrote on top of `rpc method`, it will be added to the method attributes. Eg.

    > //@Attribute description: sayHello

### Currently supported attributes

* `description` The name of this method in the service
* `output_privacy` Whether the output of the method is encrypted, {`NoPrivacy`,`PubKeyEncryption`}
* `output_cached` Whether the output of the method is cached, {`OffChainCached`, `NoCached`}

### IDL content example

* idl-content example

    > syntax = \"proto3\";\n\npackage helloworld;\n\n// The greeting service definition.\nservice Greeter {\n    //@Attribute description: sayHello\n    //@Attribute output_privacy: NoPrivacy\n    //@Attribute output_cached: NoCached\n    rpc SayHello (HelloRequest) returns (HelloReply) {}\n}\n\n// The request message containing the user's name.\nmessage HelloRequest {\n    string name = 1;\n}\n\n// The response message containing the greetings\nmessage HelloReply {\n    string message = 1;\n}\n

* IDL file example

    [test.proto](https://github.com/irisnet/irishub/blob/master/docs/features/test.proto)
