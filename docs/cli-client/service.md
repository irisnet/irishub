# iriscli service

Service module allows you to define, bind, invoke services on the IRIS Hub. [Read more about iService](../features/service.md).

## Available Commands

| Name                                              | Description                                           |
| ------------------------------------------------- | ----------------------------------------------------- |
| [define](#iriscli-service-define)                 | Define a new service                      |
| [definition](#iriscli-service-definition)         | Query a service definition                              |
| [bind](#iriscli-service-bind)                     | Bind a service                     |
| [binding](#iriscli-service-binding)               | Query a service binding                                 |
| [bindings](#iriscli-service-bindings)             | Query all bindings of a service definition                              |
| [set-withdraw-addr](#iriscli-service-set-withdraw-addr)             | Set a withdrawal address for a provider                              |
| [withdraw-addr](#iriscli-service-withdraw-addr)             | Query the withdrawal address of a provider                              |
| [update-binding](#iriscli-service-update-binding) | Update a service binding                              |
| [disable](#iriscli-service-disable)               | Disable an available service binding                   |
| [enable](#iriscli-service-enable)                 | Enable an unavailable service binding                 |
| [refund-deposit](#iriscli-service-refund-deposit) | Refund all deposit from a service binding             |
| [call](#iriscli-service-call)                     | Call a service                                 |
| [request](#iriscli-service-request)             | Query a request by the request ID                          |
| [requests](#iriscli-service-requests)             | Query service requests by the service binding or request context ID                          |
| [respond](#iriscli-service-respond)               | Respond to a service request                  |
| [response](#iriscli-service-response)             | Query a response by the request ID                             |
| [responses](#iriscli-service-responses)             | Query responses by the request context ID and batch counter                             |
| [request-context](#iriscli-service-request-context)             | Query a request context                             |
| [update](#iriscli-service-update)             | Update a request context                             |
| [pause](#iriscli-service-pause)             | Pause a request context                             |
| [start](#iriscli-service-start)             | Resume a paused request context                             |
| [kill](#iriscli-service-kill)             | Terminate a request context                             |
| [fees](#iriscli-service-fees)                     | Query the earned fees of a provider |
| [withdraw-fees](#iriscli-service-withdraw-fees)   | Withdraw the earned fees of a provider          |

## iriscli service define

Define a new service.

```bash
iriscli service define <flags>
```

**Flags:**

| Name, shorthand       | Default | Description                                                        | Required |
| --------------------- | ------- | ------------------------------------------------------------------ | -------- |
| --name        |         | Service name                                                       | Yes      |
| --description |         | Service description                                                |          |
| --author-description  |         | Service author description                                         |          |
| --tags                |         | Service tags                                                       |          |
| --schemas        |         | Service interface schemas content or path                  | Yes    |

### define a service

```bash
iriscli service define --chain-id=<chain-id> --from=<key-name> --fee=0.3iris 
--name=<service name> --description=<service description> --author-description=<author description>
--tags=tag1,tag2 --schemas=<schemas content or path/to/schemas.json>
```

### Schemas content example

* idl-content example

    > syntax = \\"proto3\\";\n\npackage helloworld;\n\n// The greeting service definition.\nservice Greeter {\n    //@Attribute description: sayHello\n    //@Attribute output-privacy: NoPrivacy\n    //@Attribute output-cached: NoCached\n    rpc SayHello (HelloRequest) returns (HelloReply) {}\n}\n\n// The request message containing the user's name.\nmessage HelloRequest {\n    string name = 1;\n}\n\n// The response message containing the greetings\nmessage HelloReply {\n    string message = 1;\n}\n

## iriscli service definition

Query a service definition.

```bash
iriscli service definition <service name>
```

### Query a service definition

Query the detailed info of the service definition with the specified service name.

```bash
iriscli service definition <service name>
```

## iriscli service bind

Bind a service.

```bash
iriscli service bind <flags>
```

**Flags:**

| Name, shorthand | Default | Description                                                               | Required |
| --------------- | ------- | ------------------------------------------------------------------------- | -------- |
| --service-name  |         | Service name                                                              | Yes      |
| --deposit       |         | Deposit of the binding                                                        | Yes      |
| --pricing        |        | Pricing of the binding, which is an instance of the Irishub Service Pricing schema                           | Yes          |

### Bind an existing service definition

The deposit needs to satisfy the minimum requirement, which is the maximal one between `price` * `MinDepositMultiple` and `MinDeposit`(`MinDepositMultiple` and `MinDeposit` are the system parameters, which can be modified through the governance).

The pricing must contain a `price` property which represents the base price of the service binding.

```bash
iriscli service bind --chain-id=<chain-id> --from=<key-name> --fee=0.3iris
--service-name=<service name> --deposit=10000iris --pricing=<pricing>
```

### Pricing example


## iriscli service binding

Query a service binding.

```bash
iriscli service binding <service name> <provider>
```

### Query a service binding

```bash
iriscli service binding <service name> <provider>
```

## iriscli service bindings

Query all bindings of a service definition.

```bash
iriscli service bindings <service name>
```

### Query service binding list

```bash
iriscli service bindings <service name>
```

## iriscli service update-binding

Update a service binding.

```bash
iriscli service update-binding <flags>
```

**Flags:**

| Name, shorthand | Default | Description                                                               | Required |
| --------------- | ------- | ------------------------------------------------------------------------- | -------- |
| --service-name  |         | Service name                                                              | Yes      |
| --deposit       |         | Deposit to be increased to the binding             |          |
| --pricing        |         | New pricing of the binding, which is an instance of the Irishub Service Pricing schema        |

### Update an existing service binding

The following example updates the service binding with the additional 10iris deposit

```bash
iriscli service update-binding --chain-id=<chain-id> --from=<key-name> --fee=0.3iris 
--service-name=<service-name> --deposit=10iris
```

## iriscli service disable

Disable an available service binding.

```bash
iriscli service disable <service name>
```

### Disable an available service binding

```bash
iriscli service disable <service name> --chain-id=<chain-id> --from=<key-name> --fee=0.3iris
```

## iriscli service enable

Enable an unavailable service binding.

```bash
iriscli service enable <service name> <flags>
```

**Flags:**

| Name, shorthand  | Default | Description                                                 | Required |
| ---------------- | ------- | ----------------------------------------------------------- | -------- |
| --deposit |         | Deposit to be increased to the binding |          |

### Enable an unavailable service binding

The following example enables a service binding with the additional 10iris deposit.

```bash
iriscli service enable <service name> --chain-id=<chain-id>  --from=<key-name> --fee=0.3iris --deposit=10iris
```

## iriscli service refund-deposit

Refund all deposits from a service binding.

```bash
iriscli service refund-deposit <service name>
```

### Refund all deposits from an unavailable service binding

Before refunding, you should [disable](#iriscli-service-disable) the service binding first.

```bash
iriscli service refund-deposit <service name> --chain-id=<chain-id>  --from=<key-name> --fee=0.3iris
```

## iriscli service call

Call a service.

```bash
iriscli service call <flags>
```

**Flags:**

| Name, shorthand | Default | Description                                        | Required |
| --------------- | ------- | -------------------------------------------------- | -------- |
| --service-name  |         | Service name                                       | Yes      |
| --providers     |         | Provider list                             | Yes      |
| --service-fee-cap |         | The maximum service fee to pay for a single request        | Yes      |
| --data      |         | Input data which is an instance of the Input JSON schema of the service definition | Yes      |
| --timeout | | Block interval after which the request will expire | |
| --repeated   |    false     | Indicate if the reqeust is repetitive                |          |
| --frequency   |         | Request generation frequency if repeated              |          |
| --total  |         | Repeated count if repeated    |          |

### Initiate a service invocation request

```bash
iriscli service call --chain-id=<chain-id> --from=<key name> --fee=0.3iris --service-name=<service name>
--providers=<provider list> --service-fee-cap=1iris --data=<request data> -timeout=100 --repeated --frequency=150 --total=100
```

## iriscli service requests

Query service requests by the service binding or request context ID.

```bash
iriscli service requests [<service name> <provider>] | [<request-context-id> <batch-counter>]
```

### Query active requests of a service binding

```bash
iriscli service requests <service name> <provider>
```

### Query service requests by the request context ID and batch counter

```bash
iriscli service requests <request-context-id> <batch-counter>
```

## iriscli service respond

Respond to a service request.

```bash
iriscli service respond <flags>
```

**Flags:**

| Name, shorthand    | Default | Description                                                    | Required |
| ------------------ | ------- | -------------------------------------------------------------- | -------- |
| --request-id       |         | The request ID                             | Yes      |
| --data    |         | Output data which is an instance of the Output JSON schema of the service definition             |          |
| --error    |         | Error msg which is an instance of the Error JSON schema of the service definition              |          |

### Respond to a service request

```bash
iriscli service respond --chain-id=<chain-id> --from=<key-name> --fee=0.3iris
--request-id=<request-id> --data=<response output>
```

:::tip
You can retrieve the `request-id` by querying the block result of [tendermint block](#iriscli-tendermint-block)
:::

## iriscli service response

Query a service response.

```bash
iriscli service response <request-id>
```

### Query a service response

```bash
iriscli service response <request-id>
```

:::tip
You can retrieve the `request-id` by querying the block result of [tendermint block](#iriscli-tendermint-block)
:::

## iriscli service request-context

Query a request context.

```bash
iriscli service request-context <request-context-id>
```

### Query a request context

```bash
iriscli service request-context <request-context-id>
```

:::tip
You can retrieve the `request-context-id` in the result of [service call](#iriscli-service-call)
:::

## iriscli service responses

Query responses by the request context ID and batch counter

```bash
iriscli service responses <request-context-id> <batch-counter>
```

### Query responses by the request context ID and batch counter

```bash
iriscli service responses <request-context-id> <batch-counter>
```

## iriscli service update

Update a request context

```bash
iriscli service update <request-context-id> <flags>
```

**Flags:**

| Name, shorthand | Default | Description                                        | Required |
| --------------- | ------- | -------------------------------------------------- | -------- |
| --providers     |         | New provider list, not updated if empty                            | Yes      |
| --service-fee-cap |         | New maximum service fee to pay for a single request, not updated if not provided       |      |
| --timeout | | New timeout, not updated if set to 0 | |
| --frequency   |         | New request frequency, not updated if set to 0           |          |
| --total  |         | New repeated count, not updated if set to 0    |          |

### Update a request context

```bash
iriscli service update <request-context-id> --chain-id=<chain-id> --from=<key name> --fee=0.3iris
--providers=<provider list> --service-fee-cap=1iris -timeout=0 --frequency=150 --total=100
```

## iriscli service pause

Pause a running request context.

```bash
iriscli service pause <request-context-id>
```

### Pause a running request context

```bash
iriscli service pause <request-context-id>
```

## iriscli service start

Resume a paused request context.

```bash
iriscli service start <request-context-id>
```

### Start a paused request context

```bash
iriscli service start <request-context-id>
```

## iriscli service kill

Terminate a request context.

```bash
iriscli service kill <request-context-id>
```

### Kill a request context

```bash
iriscli service kill <request-context-id>
```

## iriscli service fees

Query the earned fees of a provider.

### Query service fees

```bash
iriscli service fees <provider>
```

## iriscli service withdraw-fees

Withdraw the earned fees of a provider.

```bash
iriscli service withdraw-fees <flags>
```

### Withdraw the earned fees

```bash
iriscli service withdraw-fees --chain-id=<chain-id> --from=<key-name> --fee=0.3iris
```
