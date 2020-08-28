# iriscli oracle

Oracle module allows you to manage the feed on IRIS Hub

## Available Commands

| Name                                       | Description                                                                            |
| ------------------------------------------ | -------------------------------------------------------------------------------------- |
| [create](#iriscli-oracle-create)           | Create a new feed, the feed will be in "paused" state                                   |
| [start](#iriscli-oracle-start)             | Start a feed in "paused" state                                                         |
| [pause](#iriscli-oracle-pause)             | Pause a feed in "running" state                                                        |
| [edit](#iriscli-oracle-edit)               | Feed creator modify the feed information and update service invocation parameters      |
| [query-feed](#iriscli-oracle-query-feed)   | Query feed definition by feed name                                                   |
| [query-feeds](#iriscli-oracle-query-feeds) | Query a group of definitions by feed state                                          |
| [query-value](#iriscli-oracle-query-value) | Query the feed's value by feed name                                                |

## iriscli oracle create

This command is used to create a "paused" feed on IRIS Hub.

```bash
iriscli oracle create [flags]
```

**Flags:**

| Name, shorthand   | Type     | Required | Default | Description                                                                                   |
| ----------------- | -------- | -------- | ------- | --------------------------------------------------------------------------------------------- |
| --feed-name       | string   | Yes      |         | The unique identifier of the feed.                                                            |
| --description     | string   |          |         | The description of the feed.                                                                  |
| --latest-history  | uint64   | Yes      |         | The maximum Number of the latest history values to be saved for the Feed, range [1, 100].                 |
| --service-name    | string   | Yes      |         | The name of the service to be invoked by the feed.                                            |
| --input           | string   | Yes      |         | The input argument (JSON) used to invoke the service.                                         |
| --providers       | []string | Yes      |         | The list of service provider addresses.                                                       |
| --service-fee-cap | string   | Yes      |         | Only providers charging a fee lower than the cap will be invoked.                             |
| --timeout         | int64    |          |         | The maximum number of blocks to wait for a response since a request is sent, beyond which the request will be ignored. |
| --frequency       | uint64   |          |         | The invocation frequency of sending repeated requests.                                                   |
| --threshold       | uint16   |          | 1       | The minimum number of responses needed for aggregation, range [1, Length(providers)].         |
| --aggregate-func  | string   | Yes      |         | The name of predefined function for processing the service responses, e.g.avg、max、min etc.  |
| --value-json-path | string   | Yes      |         | The field name or path of Service response result used to retrieve the value property of aggregate-func from response results.        |

### Create a new feed

```bash
iriscli oracle create --chain-id=irishub --from=node0 --fee=0.3iris --feed-name="test-feed" --latest-history=10 --service-name="test-service" --input=<request-data> --providers=<provide1_address>,<provider2_address> --service-fee-cap=1iris --timeout=2 --frequency=10 --total=10 --threshold=1 --aggregate-func="avg" --value-json-path="high" --commit
```

## iriscli oracle start

This command is used to start a feed in "paused" state

```bash
iriscli oracle start [feed-name] [flags]
```

### Start a "paused" feed

```bash
iriscli oracle start test-feed --chain-id=irishub --from=node0 --fee=0.3iris --commit
```

## iriscli oracle pause

This command is used to pause a feed in "running" state

```bash
iriscli oracle pause [feed-name] [flags]
```

### Pause a "running" feed

```bash
iriscli oracle pause test-feed --chain-id=irishub --from=node0 --fee=0.3iris --commit
```

## iriscli oracle edit

This command is used to edit an existing feed on IRIS Hub.

```bash
iriscli oracle edit [feed-name] [flags]
```

**Flags:**

| Name, shorthand   | Type     | Required | Default | Description                                                                                   |
| ----------------- | -------- | -------- | ------- | --------------------------------------------------------------------------------------------- |
| --feed-name       | string   | Yes      |         | The unique identifier of the feed.                                                            |
| --description     | string   |          |         | The description of the feed.                                                                  |
| --latest-history  | uint64   | Yes      |         | The maximum Number of the latest history values to be saved for the Feed, range [1, 100].                |
| --providers       | []string | Yes      |         | The list of service provider addresses.                                                       |
| --service-fee-cap | string   | Yes      |         | Only providers charging a fee lower than the cap will be invoked.                             |
| --timeout         | int64    |          |         | The maximum number of blocks to wait for a response since a request is sent, beyond which the request will be ignored. |
| --frequency       | uint64   |          |         | The invocation frequency of sending repeated requests.                                                   |
| --threshold       | uint16   |          | 1       | The minimum number of responses needed for aggregation, range [1, Length(providers)].         |

### Edit an existed feed

```bash
iriscli oracle edit test-feed --chain-id=irishub --from=node0 --fee=0.3iris --latest-history=5 --commit
```

## iriscli oracle query-feed

This command is used to query a feed 

```bash
iriscli oracle query-feed [feed-name] [flags]
```

### Query an existed feed

```bash
iriscli oracle query-feed test-feed
```

## iriscli oracle query-feeds

This command is used to query a group of feed 

```bash
iriscli oracle query-feeds [flags]
```

**Flags:**

| Name, shorthand | Type   | Required | Default | Description                                |
| --------------- | ------ | -------- | ------- | ------------------------------------------ |
| --state         | string |          |         | the state of the feed,e.g.paused、running. |

### Query a group of feed

```bash
iriscli oracle query-feeds --state=running
```

## iriscli oracle query-value

This command is used to query the result of a specified feed

```bash
iriscli oracle query-value test-feed
```

### Query the result of an existed feed

```bash
iriscli oracle query-value test-feed
```