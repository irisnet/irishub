# Service

## Introduction

Service aims to bridge the gap between blockchain and traditional applications. It standardizes the definition and binding of off-chain services (provider registration), facilitates invocation and interaction, and mediates the service governance process (analysis and dispute resolution).

## Service Definition

### Service interface schema

The service interface, that is, input and output, needs to be specified using [JSON Schema](https://JSON-Schema.org/). Here is an example:

```json
{
  "input": {
    "$schema": "http://json-schema.org/draft-04/schema#",
    "title": "service-def-input-example",
    "description": "Schema for a service input example",
    "type": "object",
    "properties": {
      "base": {
        "description": "base token denom",
        "type": "string"
      },
      "quote": {
        "description": "quote token denom",
        "type": "string"
      }
    },
    "required": [
      "base",
      "quote"
    ]
  },
  "output": {
    "$schema": "http://json-schema.org/draft-04/schema#",
    "title": "service-def-output-example",
    "description": "Schema for a service output example",
    "type": "object",
    "properties": {
      "price": {
        "description": "price",
        "type": "number"
      }
    },
    "required": [
      "price"
    ]
  }
}
```

### Operation

`CLI`

```bash
# Create service definition
iris tx service define <service-name> <schemas-json or path/to/schemas.json> --description=<service-description> --author-description=<author-description> --tags=<tag1,tag2 ,...>

# Query service definition
iris query service definition <service-name>
```

## Service binding

By creating bindings to existing service definitions, corresponding services can be provided. Binding mainly consists of four parts: _provider address_, _pricing_, _deposit_ and _quality of service_.

### Provider address

The provider address is an endpoint used by _service provider_ (ie off-chain service/process) to listen for requests. Before the service provider can accept and process the service request, its operator or owner must create an on-chain address for it and initiate a `binding` transaction to associate this address with the relevant service definition.

To call a service, a user or consumer initiates a request to a provider address bound to a valid service by initiating a request transaction; the service provider detects and processes the request, and sends the processing result through a response transaction.

### Pricing

Pricing specifies how service providers charge for the services they provide. Pricing must conform to this [schema](./service-pricing.json). Here is an example:

```json
{
  "price": "100iris",
  "promotions_by_time": [
    {
      "start_time": "2020-01-01T00:00:00Z",
      "end_time": "2020-03-31T23:59:59Z",
      "discount": 0.7
    },
    {
      "start_time": "2020-04-01T00:00:00Z",
      "end_time": "2020-06-30T23:59:59Z",
      "discount": 0.9
    }
  ]
}
```

Price is a consideration for consumers in selecting from multiple providers that provide the same service.

### Deposit

Operating a service provider means important service responsibilities, so creating a service binding requires a certain amount of deposit. The deposit amount must be greater than the _deposit threshold_, which is the maximum value of `MinDepositMultiple * price` and `MinDeposit`. If the service provider fails to respond to the request before the timeout, a small part of its bound deposit, namely `SlashFraction * deposit`, will be fined and destroyed. If the deposit falls below the threshold, the service binding will be temporarily disabled until its owner adds enough deposit to reactivate it.

> **_Tip:_** `service/MinDepositMultiple`, `service/MinDeposit` and `service/SlashFraction` are the system parameters that can be changed.

### service quality

The quality of service commitment is based on the average number of blocks required by the provider to send the service response back to the blockchain. This is another factor that consumers consider when choosing potential providers.

### Operation

Service binding can be updated by its owner at any time to adjust pricing, increase deposits or change QoS; it can also be disabled and re-enabled. If the owner of the service provider does not want to provide services anymore, he needs to disable the binding and wait for a period of time before he can get back the deposit.

`CLI`

```bash
# Create service binding
iris tx service bind --service-name=<service-name> --provider=<provider-address> --deposit=<deposit> --qos=<qos> --pricing=<pricing-json or path/to /pricing.json>

# Update service binding
iris tx service update-binding <service-name> <provider-address> --deposit=<added-deposit> --qos=<qos> --pricing=<pricing-json or path/to/pricing.json>

# Enable an unavailable service binding
iris tx service enable <service-name> <provider-address> <added-deposit>

# Disable an available service binding
iris tx service disable <service-name> <provider-address>

# Get back the deposit bound to the service
iris tx service refund-deposit <service-name> <provider-address>

# Query all bindings of a service
iris query service bindings <service-name>

# Query all bindings owned by an account
iris query service bindings <service-name> --owner <address>

# Query the specified service binding
iris query service binding <service-name> <provider-address>

# Query pricing schema
iris query service schema pricing
```

## Service call

### Request context

Consumers specify how to call a service by creating a request context. _Request context_ The actual request is automatically generated like a smart contract. _Request context_ consists of some parameters, which can be roughly divided into the following four groups:

#### Goals and inputs

* _Service name_: The name of the target service to be called
* _Input data_: json format data conforming to the target service input schema

#### Provider filtering

* _Provider List_: A comma-separated list of addresses of candidate service providers
* _Service Fee Limit_: The maximum service fee that consumers are willing to pay for each call
* _Timeout_: The number of blocks the consumer is willing to wait for to receive a response

#### Response processing

* _Module_: the name of the module containing the callback function
* _Response threshold_: the minimum number of responses required to call the callback function

> **_Tip: _** These two parameters cannot be set from CLI and API; they are only available for other modules that use iService, such as [oracle](oracle.md) and [random](random.md).

#### Repeatability

* _Repeat_: A Boolean flag indicating whether the request context can be repeated
* _Frequency_: the number of block intervals between repeated calling batches
* _Total number_: The total number of repeated calling batches, a negative number means unlimited

### Request batch

For a repetitive request context, new requests _batch_ will be generated at the specified frequency until the specified number of batches is reached or the balance of the consumer (that is, the creator of the request context) is insufficient. For non-repetitive request contexts, only one request batch will be generated.

A request batch is composed of several _request_ objects, _request_ represents a service call initiated to a service provider that meets the selection criteria. Only those providers whose fees are not higher than the `service charge limit` and whose QoS is better than `timeout` can be selected.

### Operation

After successfully creating a request context, a _context ID_ will be returned to the consumer, and the context will be automatically started. Consumers can then update, pause, and start the context as they wish, or they can terminate it permanently.

`CLI`

```bash
# Create a repetitive request context (no callback function)
iris tx service call --service-name=<service-name> --data=<request-input> --providers=<provider-list> --service-fee-cap=1iris --timeout 50 --repeated- -frequency=50 --total=100

# Update an existing request context
iris tx service update <request-context-id> --frequency=20 --total=200

# Pause a running request context
iris tx service pause <request-context-id>

# Start a suspended request context
iris tx service start <request-context-id>

# Permanently terminate a request context
iris tx service kill <request-context-id>

# Query request context by ID
iris query service request-context <request-context-id>

# Query all requests in a request batch
iris query service requests <request-context-id> <batch-counter>

# Query all responses of a request batch
iris query service responses <request-context-id> <batch-counter>

# Query the corresponding response by request ID
iris query service response <request-id>
```

## Service response

Service providers monitor their own requests on the chain through query or event subscription. After processing a request, the service provider sends back a result of the object and optional service input