# Slashing
Slashing module can unjail validator previously jailed for downtime

## Available Commands

| Name                                                | Description                                     |
| --------------------------------------------------- | ----------------------------------------------- |
| [unjail](#iris-tx-slashing-unjail)                  | Unjail validator previously jailed for downtime |
| [params](#iris-query-slashing-params)               | Query the current slashing parameters           |
| [signing-info](#iris-query-slashing-signing-info)   | Query a validator's signing information         |
| [signing-infos](#iris-query-slashing-signing-infos) | Query signing information of all validators     |

## iris tx slashing unjail

Unjail validator previously jailed for downtime.

```bash
iris tx slashing unjail [flags]
```

## iris query slashing params

Query the current slashing parameters .

```bash
iris query bank total [flags]
```

## iris query slashing signing-info

Query a validator's signing information

```bash
iris query slashing signing-info [validator-conspub] [flags]
```

## iris query slashing signing-infos

Query signing information of all validators

```bash
iris query slashing signing-infos [flags]
```