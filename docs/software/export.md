# Export Blockchain State

## Description

IRISnet can export blockchain state at any height and output json format string. Save the out json string to a json file and the json file can be used as genesis file of a new blockchain. This can be accomplished by command `iris export`

## Usage

```
iris export [flags]
```

## Flags

| Name，shorthand     | type   | Required | Default  | Description    |
| ------------------- | -----  | -------- | -------- | -------------- |
| --for-zero-height   | bool   | false    | false    | Do some clean up work before exporting state. If you want use the exported state to start a new blockchain, please add this flag. Otherwise, just leave out it |
| --height            | int    | false    | 0        | Specify the height, defalut value is 0 which means to export the latest state |
| --home              | string | false    | $HOME/.iris       | Specify the directory which stores node config and blockchain data |

## Examples

1. Export the latest blockchain state:
```
iris export
```
2. Export blockchain state at height 10000
```
iris export --height 10000
```
3. If you want to export the blockchain state at height 105000 and use the exported state as genesis state of another blockchain, you can try this command:
```
iris export --height=105000 --for-zero-height --home=[your_home]
```