# Service

> **_NOTE:_** Commands shown in this document are for illustration purpose only.  For accurate syntax of commands, please refer to [cli docs](../cli-client/service.md).

## Summary

IRIS Service (a.k.a. "iService") is intended to bridge the gap between the blockchain world and the conventional application world.  It formalizes off-chain service definition and binding (provider registration), facilitates invocation and interaction with those services, and mediates service governance process (profiling and dispute resolution).

## Service Definition

### Service interface schema
Any user can define services on the blockchain. The interface of a service must be specified in terms of its _input_ and _output_ using the standard language of [JSON Schema](https://json-schema.org/).  Here is an example:

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

### Service result schema
Service providers respond to a user request (containing an input object) by sending back a response that consists of a _result_ object and an optional output object, the latter of which is required only when the result code equals 200.  The result object must conform to this [schema](service-result.json), and here is an example of a compliant instance:

```json
{
  "result" : {
    "code": 400,
    "message": "user input out of range"
  }
}
```

Once you have a definition ready, you can publish it to the blockchain by executing the following command:

```bash
# create a new service definition
iriscli service define <service-name> <schemas-json or path/to/schemas.json> --description=<service-description> --author-description=<author-description> --tags=<tag1,tag2,...>

# query service definition
iriscli service definition <service-name>
```

## Service Binding

Whoever is willing to provide a service as specified by an existing definition can do so by creating a _binding_ to that definition.  A binding essentially consists of three components: _provider address_ (address of whoever executes the `bind` transaction), _pricing_ and _deposit_.  

### Provider address
A consumer should be able to publish a service request (input) destined to the provider address, and see a response (output) transaction coming back from this address.

### Pricing
The pricing object must conform to this [schema](service-pricing.json), and the following is a compliant instance:  

```json
{
  "price": "0.1iris",
  "promotions_by_time": [
    {
      "start_time": "2020-01-01T00:00:00Z",
      "end_time": "2020-03-31T23:59:59Z",
      "discount": 0.7
    },
    {
      "start_time": "2020-04-01T00:00:00Z",
      "end_time": "2019-06-30T23:59:59Z",
      "discount": 0.9
    }
  ]
}
```

### Deposit
Operating a service provider signifies serious commitment, therefore, a deposit is required for creating a binding.  The deposit amount must be larger than the _deposit threshold_, derived as `max(DepositMultiple*price, MinDeposit)`.  If a provider fails to respond to a request before it times out, a small portion of its binding deposit, i.e., `SlashFraction*deposit`, will be slashed.  Should the deposit drop below the threshold, the binding would be disabled temporarily until its owner re-activates it by adding more deposit.

> **_NOTE:_** `service/DepositMultiple`, `service/MinDeposit` and `service/SlashFraction` are system parameters that can be changed through on-chain [governance](governance.md).

### Lifecycle
Service bindings can be updated at any time by their owners to adjust pricing or increase deposit; they can be disabled and re-enabled as well.  If a binding owner no longer wants to run the service provider, she needs to disable the binding and wait for a certain period of time before she can claim back her deposit.

```bash
# create a new service binding
iriscli service bind <service-name> <deposit> <pricing-json or path/to/pricing.json>

# update a service binding
iriscli service update-binding <service-name> --deposit=<added-deposit> --pricing=<pricing-json or path/to/pricing.json>

# set withdrawal account
iriscli service set-withdraw-addr <withdrawal-address>

# withdraw earned fees into withdrawal account
iriscli service withdraw-fees

# enable an inactive service binding
iriscli service enable <service-name> <added-deposit>

# disable an active service binding
iriscli service disable <service-name>

# request refund of service binding deposit
iriscli service refund-deposit <service-name>

# a trustee withdraws service tax into given account
iriscli service withdraw-tax <destination-address>

# query service binding
iriscli service binding <service-name> <provider-address>

# query service bindings
iriscli service bindings <service-name>

# query a provider's withdrawal address
iriscli service withdraw-addr <provider-address>

# query a provider's earned fees
iriscli service fees <provider-address>

# query system schemas (valid names: pricing, result)
iriscli service schema <schema-name>
```

## Service Invocation

If the service consumer needs to initiate a service invocation request, the service fee specified by the service provider needs to be paid. The service provider needs to respond to the service request within the block height defined by `MaxRequestTimeout`. If the service provider does not respond in time, the deposit of the 'SlashFraction' ratio will be deducted from the service provider's service binding deposit and the service fee will be refunded to the service consumer's return pool. If the service call is responded normally, the system will deduct the `ServiceFeeTax` ratio from the service fee, and add the remaining service fee to the service provider's incoming pool. The service provider/consumer can initiate the `withdraw-fees`/`refund-fees` transaction to retrieve all of the tokens in the incoming/return pool.

```bash
```
