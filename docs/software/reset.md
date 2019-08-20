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
 | --height            | uint   | false    | 0        | Reset state to a particular height (greater than latest height means latest height) |		
 | --home              | string | false    | $HOME/.iris       | Specify the directory which stores node config and blockchain data |		
 
## Examples

1. Reset the blockchain state to block 100:
```
 iris reset --height 100 --home=<path_to_your_home>
```