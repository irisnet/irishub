# iriscli service

Service module allows you to define, bind, invoke services on the IRIS Hub. [Read more about iService](#TODO)

## Available Commands

| Name                                                | Description                                           |
| --------------------------------------------------  | ----------------------------------------------------- |
| [define](#iriscli-service-define)                   | Create a new service definition                       |
| [definition](#iriscli-service-definition)           | Query service definition                              |
| [bind](#iriscli-service-bind)                       | Create a new service binding                          |
| [binding](#iriscli-service-binding)                 | Query service binding                                 |
| [bindings](#iriscli-service-bindings)               | Query service bindings                                |
| [update-binding](#iriscli-service-update-binding)   | Update a service binding                              |
| [disable](#iriscli-service-disable)                 | Disable a available service binding                   |
| [enable](#iriscli-service-enable)                   | Enable an unavailable service binding                 |
| [refund-deposit](#iriscli-service-refund-deposit)   | Refund all deposit from a service binding             |
| [call](#iriscli-service-call)                       | Call a service method                                 |
| [requests](#iriscli-service-requests)               | Query service requests                                |
| [respond](#iriscli-service-respond)                 | Respond a service method invocation                   |
| [response](#iriscli-service-response)               | Query a service response                              |
| [fees](#iriscli-service-fees)                       | Query return and incoming fee of a particular address |
| [refund-fees](#iriscli-service-refund-fees)         | Refund all fees from service return fees              |
| [withdraw-fees](#iriscli-service-withdraw-fees)     | Withdraw all fees from service incoming fees          |

## iriscli service define

Create a new service definition

```bash
iriscli service define <flags>
```

**Unique Flags:**

| Name, shorthand       | Default | Description                                                        | Required |
| --------------------- | ------- | ------------------------------------------------------------------ | -------- |
| --service-description |         | Service description                                                |          |
| --author-description  |         | Service author description                                         |          |
| --service-name        |         | Service name                                                       | true     |
| --tags                |         | Service tags                                                       |          |
| --idl-content         |         | Content of service interface description language                  |          |
| --file                |         | Path of file which contains service interface description language |          |

### define a service

```bash
iriscli service define --chain-id=<chain-id>  --from=<key-name> --fee=0.3iris --service-name=<service-name> --service-description=<service-description> --author-description=<author-description> --tags=tag1,tag2 --idl-content=<idl-content> --file=test.proto
```

Idl-content can be replaced by file if the file item is not empty.  [Example of IDL content](#idl-content-example).

### IDL content example

* idl-content example

    > syntax = \\"proto3\\";\n\npackage helloworld;\n\n// The greeting service definition.\nservice Greeter {\n    //@Attribute description: sayHello\n    //@Attribute output-privacy: NoPrivacy\n    //@Attribute output-cached: NoCached\n    rpc SayHello (HelloRequest) returns (HelloReply) {}\n}\n\n// The request message containing the user's name.\nmessage HelloRequest {\n    string name = 1;\n}\n\n// The response message containing the greetings\nmessage HelloReply {\n    string message = 1;\n}\n

* IDL file example

    [test.proto](https://github.com/irisnet/irishub/blob/master/docs/features/test.proto)

## iriscli service definition

Query service definition

```bash
iriscli service definition <flags>
```

**Unique Flags:**

| Name, shorthand | Default | Description                                     | Required |
| --------------- | ------- | ----------------------------------------------- | -------- |
| --def-chain-id  |         | The ID of the blockchain defined of the service | true     |
| --service-name  |         | Service name                                    | true     |

### Query a service definition

Query the detail info of the service definition which has the specified define chain id and service name.

```bash
iriscli service definition --def-chain-id=<service-define-chain-id> --service-name=<service-name>
```

## iriscli service bind

Create a new service binding

```bash
iriscli service bind <flags>
```

**Unique Flags:**

| Name, shorthand       | Default | Description                                                               | Required |
| --------------------- | ------- | ------------------------------------------------------------------------- | -------- |
| --avg-rsp-time        |         | The average service response time in milliseconds                         | true     |
| --bind-type           |         | Type of binding, valid values can be Local and Global                     | true     |
| --def-chain-id        |         | The ID of the blockchain defined of the service                           | true     |
| --deposit             |         | Deposit of binding                                                        | true     |
| --prices              |         | Prices of binding, will contains all method                               |          |
| --service-name        |         | Service name                                                              | true     |
| --usable-time         |         | An integer represents the number of usable service invocations per 10,000 | true     |

### Add a binding to an existing service definition

In service binding, you need to define a `deposit`, the minimum mortgage amount of this `deposit` is `price` * `MinDepositMultiple` (defined by system parameters, can be modified through governance)

```bash
iriscli service bind --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --service-name=<service-name> --def-chain-id=<service-define-chain-id> --bind-type=Local --deposit=1000iris --prices=1iris --avg-rsp-time=10000 --usable-time=9999
```

## iriscli service binding

Query service binding

```bash
iriscli service binding <flags>
```

**Unique Flags:**

| Name, shorthand | Default                    | Description                                        | Required |
| --------------- | -------------------------- | -------------------------------------------------- | -------- |
| --bind-chain-id |                            | The ID of the blockchain bond of the service       | true     |
| --def-chain-id  |                            | The ID of the blockchain defined of the service    | true     |
| --provider      |                            | Bech32 encoded account created the service binding | true     |
| --service-name  |                            | Service name                                       | true     |

### Query a service binding

```bash
iriscli service binding --def-chain-id=<service-define-chain-id> --service-name=<service-name> --bind-chain-id=<service-bind-chain-id> --provider=<provider-address>
```

## iriscli service bindings

Query service bindings

```bash
iriscli service bindings <flags>
```

**Unique Flags:**

| Name, shorthand | Default | Description                                     | Required |
| --------------- | ------- | ----------------------------------------------- | -------- |
| --def-chain-id  |         | The ID of the blockchain defined of the service | true     |
| --service-name  |         | Service name                                    | true     |

### Query service binding list

```bash
iriscli service bindings --def-chain-id=<chain-id> --service-name=<service-name>
```

## iriscli service update-binding

Update a service binding

```bash
iriscli service update-binding <flags>
```

**Unique Flags:**
| Name, shorthand | Default | Description                                                               | Required |
| --------------- | ------- | ------------------------------------------------------------------------- | -------- |
| --avg-rsp-time  |         | The average service response time in milliseconds                         |          |
| --bind-type     |         | Type of binding, valid values can be Local and Global                     |          |
| --def-chain-id  |         | The ID of the blockchain defined of the service                           | true     |
| --deposit       |         | Deposit of binding, will add to the current deposit balance               |          |
| --prices        |         | Prices of binding, will contains all method                               |          |
| --service-name  |         | Service name                                                              | true     |
| --usable-time   |         | An integer represents the number of usable service invocations per 10,000 |          |

### Update an existing service binding

The following example updates a service binding alone with 10iris additional deposit

```bash
iriscli service update-binding --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --service-name=<service-name> --def-chain-id=<service-define-chain-id> --bind-type=Local --deposit=10iris --prices=1iris --avg-rsp-time=10000 --usable-time=9999
```

## iriscli service disable

Disable an active service binding

```bash
iriscli service disable <flags>
```

**Unique Flags:**

| Name, shorthand       | Default | Description                                     | Required |
| --------------------- | ------- | ----------------------------------------------- | -------- |
| --def-chain-id        |         | The ID of the blockchain defined of the service | true     |
| --service-name        |         | Service name                                    | true     |

### Disable an active service binding

```bash
iriscli service disable --chain-id=<chain-id>  --from=<key-name> --fee=0.3iris --def-chain-id=<service-define-chain-id> --service-name=<service>
```

## iriscli service enable

Enable an inactive service binding

```bash
iriscli service enable <flags>
```

**Unique Flags:**

| Name, shorthand       | Default | Description                                                 | Required |
| --------------------- | ------- | ----------------------------------------------------------- | -------- |
| --def-chain-id        |         | The ID of the blockchain defined of the service             | true     |
| --deposit string      |         | Deposit of binding, will add to the current deposit balance |          |
| --service-name        |         | Service name                                                | true     |

### Enable an inactive service binding

The following example activates a service binding alone with 10iris additional deposit

```bash
iriscli service enable --chain-id=<chain-id>  --from=<key-name> --fee=0.3iris --def-chain-id=<service-define-chain-id> --service-name=<service-name> --deposit=10iris
```

## iriscli service refund-deposit

Refund all deposits from a service binding

```bash
iriscli service refund-deposit <flags>
```

**Unique Flags:**

| Name, shorthand       | Default | Description                                     | Required |
| --------------------- | ------- | ----------------------------------------------- | -------- |
| --def-chain-id        |         | The ID of the blockchain defined of the service | true     |
| --service-name        |         | Service name                                    | true     |

### Refund all deposits from an inactive service binding

Before refunding, you should [disable](disable) the service binding first.

```bash
iriscli service refund-deposit --chain-id=<chain-id>  --from=<key-name> --fee=0.3iris --def-chain-id=<service-define-chain-id> --service-name=<service-name>
```

## iriscli service call

Invoke a service method

```bash
iriscli service call <flags>
```

**Unique Flags:**

| Name, shorthand       | Default | Description                                        | Required |
| --------------------- | ------- | -------------------------------------------------- | -------- |
| --def-chain-id        |         | The ID of the blockchain defined of the service    | true     |
| --service-name        |         | Service name                                       | true     |
| --method-id           |         | The method id called                               | true     |
| --bind-chain-id       |         | The ID of the blockchain bond of the service       | true     |
| --provider            |         | Bech32 encoded account created the service binding | true     |
| --service-fee         |         | Fee to pay for a service invocation                |          |
| --request-data        |         | Hex encoded request data of a service invocation   |          |

### Initiate a service invocation request

```bash
iriscli service call --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --def-chain-id=<service-define-chain-id> --service-name=<service-name> --method-id=1 --bind-chain-id=<service-bind-chain-id> --provider=<provider-address> --service-fee=1iris --request-data=<request-data>
```

## iriscli service requests

Query service requests

```bash
iriscli service requests <flags>
```

**Unique Flags:**

| Name, shorthand       | Default | Description                                        | Required |
| --------------------- | ------- | -------------------------------------------------- | -------- |
| --def-chain-id        |         | The ID of the blockchain defined of the service    | true     |
| --service-name        |         | Service name                                       | true     |
| --bind-chain-id       |         | The ID of the blockchain bond of the service       | true     |
| --provider            |         | Bech32 encoded account created the service binding | true     |

### Query service request list

```bash
iriscli service requests --def-chain-id=<service-define-chain-id> --service-name=<service-name> --bind-chain-id=<service-bind-chain-id> --provider=<provider-address>
```

## iriscli service respond

Respond a service method invocation

```bash
iriscli service respond <flags>
```

**Unique Flags:**

| Name, shorthand       | Default | Description                                                    | Required |
| --------------------- | ------- | -------------------------------------------------------------- | -------- |
| --request-chain-id    |         | The ID of the blockchain that the service invocation initiated | true     |
| --request-id          |         | The ID of the service invocation                               | true     |
| --response-data       |         | Hex encoded response data of a service invocation              |          |

### Respond to a service invocation

```bash
iriscli service respond --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --request-chain-id=<request-chain-id> --request-id=<request-id> --response-data=<response-data>
```

:::tip
You can figure out the `request-id` in the result of [service call](#iriscli-service-call)
:::

## iriscli service response

Query a service response

```bash
iriscli service response <flags>
```

**Unique Flags:**

| Name, shorthand       | Default | Description                                                    | Required |
| --------------------- | ------- | -------------------------------------------------------------- | -------- |
| --request-chain-id    |         | The ID of the blockchain that the service invocation initiated | true     |
| --request-id          |         | The ID of the service invocation                               | true     |

### Query a service response

```bash
iriscli service response --request-chain-id=<request-chain-id> --request-id=<request-id>
```

:::tip
You can figure out the `request-id` in the result of [service call](#iriscli-service-call)
:::

## iriscli service fees

Query return and incoming fee of a service provider address

### Query service fees

```bash
iriscli service fees <service-provider-address>
```

Example Output:

```json
{
  "returned-fee": [
    {
      "denom": "iris-atto",
      "amount": "10000000000000000"
    }
  ],
  "incoming-fee": [
    {
      "denom": "iris-atto",
      "amount": "10000000000000000"
    }
  ]
}
```

## iriscli service refund-fees

Refund all fees from service return fees

```bash
iriscli service refund-fees <flags>
```

### Refund fees from service return fees

```bash
iriscli service refund-fees --chain-id=<chain-id> --from=<key-name> --fee=0.3iris
```

## iriscli service withdraw-fees

Withdraw service incoming fees

```bash
iriscli service withdraw-fees <flags>
```

### Withdraw service incoming fees

```bash
iriscli service withdraw-fees --chain-id=<chain-id> --from=<key-name> --fee=0.3iris
```
