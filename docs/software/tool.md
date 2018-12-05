# IRIS Command tool

## Introduction
`iristool` include [monitor](./monitor.md) and debug now.

## debug
Simple tool for simple debugging.

We try to accept both hex and base64 formats and provide a useful response.

Note we often encode bytes as hex in the logs, but as base64 in the JSON.

### Usage

* Pubkeys 
The following give the same result:

```bash
iristool debug pubkey TZTQnfqOsi89SeoXVnIw+tnFJnr4X8qVC0U8AsEmFk4=
iristool debug pubkey 4D94D09DFA8EB22F3D49EA17567230FAD9C5267AF85FCA950B453C02C126164E
```

* Txs
Pass in a hex/base64 tx and get back the full JSON

```bash
iristool debug tx [hex or base64 transaction]
```

* Hack
This is a command with boilerplate for using Go as a scripting language to hack on an existing Iris state.

If you run 
```
iristool debug hack $HOME/.iris
```
on that state, it will do a binary search on the state history to find when the state invariant was violated.
