# iriscli keys export

## Description

Export single keystore to a json file

## Usage

```shell
iriscli keys export Mykey --output-file=<path_to_backup_keystore>
```

## Flags

| Name, shorthand | Default   | Description          | Required |
| --------------- | --------- | -------------------- | -------- |
| --output-file   |           | The path of keystore |          |
| --help, -h      |           | Help for export      |          |


## Example

### Export

A password is required before exporting to a json file
```shell
iriscli keys export Mykey --output-file=<path_to_backup_keystore>
```

### Import

Use the specified password to import the keystore
```shell
iriscli keys add Mykey --recover --keystore=<path_to_backup_keystore>
```
