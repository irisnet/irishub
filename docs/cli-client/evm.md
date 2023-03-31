# EVM

## Available Commands

| Name                           | Description                                              |
| ------------------------------ | -------------------------------------------------------- |
| [raw](#iris-tx-evm-raw)        | Build cosmos transaction from raw ethereum transaction.  |
| [code](#iris-q-evm-code)       | Gets code from an account.                               |
| [params](#iris-q-evm-params)   | Get the evm params.                                      |
| [storage](#iris-q-evm-storage) | Gets storage for an account with a given key and height. |


## iris tx evm raw

Build cosmos transaction from raw ethereum transaction.

```bash
iris tx evm raw <TX_HEX>
```

## iris q evm code

Gets code from an account.

```bash
iris q evm code <ADDRESS>
```

## iris q evm params

Get the evm params.

```bash
iris q evm params
```

## iris q evm storage

Gets storage for an account with a given key and height.

```bash
iris q evm storage <ADDRESS> <KEY>
```