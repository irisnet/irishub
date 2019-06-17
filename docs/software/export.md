# Export Blockchain State

## Description

IRISnet can export blockchain state and output json format string. Save the out json string to a json file and the json file can be used as genesis file of a new blockchain. This can be accomplished by command `iris export`.

By default, IRISnet stores snapshots of every 10,000 blocks and the last 100 blocks. You can export the blockchain state of any existing snapshot height.

If you want to export the state of a nonexistent snapshot height, you need to [reset](reset.md) the blockchain state to the specified height first. 

## Usage
```		
 iris export <flags>
```
## Flags

 | Name，shorthand     | type   | Required | Default  | Description    |		
 | ------------------- | -----  | -------- | -------- | -------------- |		
 | --for-zero-height   | bool   | false    | false    | Do some clean up work before exporting state. If you want use the exported state to start a new blockchain, please add this flag. Otherwise, just leave out it |
 | --height            | uint   | false    | 0        | Export state from a particular height, default value is 0 which means to export the latest state |	
 | --home              | string | false    | $HOME/.iris  | Specify the directory which stores node config and blockchain data |		
 | --output-file       | string | false    | genesis.json | Target file to save exported state |
  
## Examples

1. Export the current blockchain state 

```		
 iris export --home=<path_to_your_home>
```

2. Export blockchain state from a particular height, the height must be an existing snapshot height
```
iris export --height 10000 --home=<path_to_your_home>
```

3. If you want to export the blockchain state from a particular height and use the exported state as genesis state of another blockchain
```		
iris export --height 10000 --for-zero-height --home=<path_to_your_home>	
```