# Record

Record module allows you to manage record on IRIS Hub

## Available Commands

| Name                                | Description        |
| ----------------------------------- | ------------------ |
| [create](#iris-tx-record-create)    | Create a record    |
| [record](#iris-query-record-record) | Query record by id |

## iris tx record create

Create a record

```bash
iris tx record create [digest] [digest-algo] [flags]
```

**Flags:**

| Name, shorthand | Type   | Required | Default | Description                                |
| --------------- | ------ | -------- | ------- | ------------------------------------------ |
| --uri           | string |          |         | Source uri of record, such as an ipfs link |
| --meta          | string |          |         | meta data of record                        |

## iris query record record

Query record by id

```bash
iris query record record [record-id]
```
