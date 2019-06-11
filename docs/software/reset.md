# Reset Blockchain State

## Description

IRISnet can reset blockchain state at any height. This can be accomplished by command `iris reset`.

## Usage
```		
 iris reset <flags>
```
## Flags

 | Name，shorthand     | type   | Required | Default  | Description    |		
 | ------------------- | -----  | -------- | -------- | -------------- |		
 | --height            | int    | false    | 0        | Specify the height, default value is 0 which means to export the latest state |		
 | --home              | string | false    | $HOME/.iris       | Specify the directory which stores node config and blockchain data |		
 
## Examples

1. Reset the blockchain state to block 100:
```		
 iris reset --height 100 --home=<path_to_your_home>
```