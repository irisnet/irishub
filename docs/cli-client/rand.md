# iriscli rand

Rand module allows you to post a random number request to the IRIS Hub and query the random numbers or the pending random number requests.

## Available Commands

| Name                                       | Description                                                      |
| ------------------------------------------ | ---------------------------------------------------------------- |
| [request-rand](#iriscli-rand-request-rand) | Request a random number                                          |
| [query-rand](#iriscli-rand-query-rand)     | Query the generated random number by the request id              |
| [query-queue](#iriscli-rand-query-queue)   | Query the pending random number requests with an optional height |

## iriscli rand request-rand

Request a random number.

```bash
iriscli rand request-rand <flags>
```

**Flags:**

| Name, shorthand  | Type   | Required | Default | Description                                                                  |
| ---------------- | ------ | -------- | ------- | ---------------------------------------------------------------------------- |
| --block-interval | uint64 |          | 10      | The block interval after which the requested random number will be generated |

### Request a random number

Post a random number request to the IRIS Hub, the random number will be generated after `--block-interval` blocks.

```bash
iriscli rand request-rand --block-interval=100 --from=<key-name> --chain-id=irishub --fee=0.3iris --commit
```

:::tip
You will get a unique request id if the tx is committed, which can be used to query the status of the request. You can also [query the tx detail](./tendermint.md#iriscli-tendermint-tx) to get the request id.
:::

## iriscli rand query-rand

Query the generated random number by the request id.

```bash
iriscli rand query-rand <flags>
```

**Flags:**

| Name, shorthand | Type   | Required | Default | Description                               |
| --------------- | ------ | -------- | ------- | ----------------------------------------- |
| --request-id    | string | Yes      |         | The request id returned by the request tx |

## Query a random number

Query the random number after it is generated.

```bash
iriscli rand query-rand --request-id=035a8d4cf64fcd428b5c77b1ca85bfed172d3787be9bdf0887bbe8bbeec3932c
```

## iriscli rand query-queue

Query the pending random number requests with an optional block height.

```bash
iriscli rand query-queue <flags>
```

**Flags:**

| Name, shorthand | Type  | Required | Default | Description                                                |
| --------------- | ----- | -------- | ------- | ---------------------------------------------------------- |
| --queue-height  | int64 |          | 0       | The block height at which random numbers will be generated |

## Query random number request queue

Query the pending random number requests with an optional block height at which random numbers will be generated.

```bash
iriscli rand query-queue --queue-height=100000
```
