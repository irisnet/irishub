# Random

Rand module allows you to post a random number request to the IRIS Hub and query the random numbers or the pending random number requests.

## Available Commands

| Name                                       | Description                                                  |
| ------------------------------------------ | ------------------------------------------------------------ |
| [request-rand](#iris-tx-rand-request-rand) | Request a random number                                      |
| [query-rand](#iris-q-rand-rand)            | Query the generated random number by the request id          |
| [query-queue](#iris-q-rand-queue)          | Query the pending random number requests with an optional height |

## iris tx rand request-rand

Request a random number.

```bash
iris tx rand request-rand [flags]
```

**Flags:**

| Name, shorthand   | Type   | Required | Default | Description                                                                  |
| ----------------- | ------ | -------- | ------- | ---------------------------------------------------------------------------- |
| --block-interval  | uint64 | true     | 10      | The block interval after which the requested random number will be generated |
| --oracle          | bool   |          | false   | Whether to use the oracle method                                             |
| --service-fee-cap | string |          | ""      | Max service fee, required if "oracle" is true                                |

### Request a random number

Post a random number request to the IRIS Hub, the random number will be generated after `--block-interval` blocks.

```bash
# without oracle
iris tx rand request-rand --block-interval=100 --from=<key-name> --chain-id=irishub --fee=0.3iris --commit

# with oracle
iris tx rand request-rand --block-interval=100 --oracle=true --service-fee-cap=1iris --from=<key-name> --chain-id=irishub --fee=0.3iris --commit
```

:::tip
You will get a unique request id if the tx is committed, which can be used to query the status of the request. You can also [query the tx detail](./tendermint.md#iriscli-tendermint-tx) to get the request id.
:::

## iris q rand rand

Query the generated random number by the request id.

```bash
iris query rand rand <request-id> [flags]
```

### Query a random number

Query the random number after it is generated.

```bash
iris q rand rand <request-id>
```

## iris q rand queue

Query the pending random number requests with an optional block height.

```bash
iris query rand queue <gen-height> [flags]
```

### Query random number request queue

Query the pending random number requests with an optional block height at which random numbers will be generated or request service.

```bash
iris query rand queue 100000
```
