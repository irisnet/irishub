# Bech32 on IRISnet

Bech32 is a new Bitcoin address format proposed by Pieter Wuille and Greg Maxwell. Besides Bitcoin addresses, Bech32 can encode any short binary data. In the IRIS network, keys and addresses may refer to a number of different roles in the network like accounts, validators etc. The IRIS network is designed to use the Bech32 address format to provide robust integrity checks on data. The human readable part(HRP) makes it more efficient to read and the users could see error messages.


## Human Readable Part Table


| HRP        | Definition |
| ------------- |:-------------:|
|faa	|IRISnet Account Address|
|fap|	IRISnet Account Public Key|
|fva	|IRISnet Consensus Address|
|fvp|	IRISnet Consensus Public Key|

## Separator

Why include a separator in addresses? That way the human-readable part is unambiguously separated from the data part, avoiding potential collisions with other human-readable parts that share a prefix. It also allows us to avoid having character-set restrictions on the human-readable part. The separator is 1 because using a non-alphanumeric character would complicate copy-pasting of addresses (with no double-click selection in several applications). Therefore an alphanumeric character outside the normal character set was chosen.

## Encoding

Not all interfaces to users IRISnet should be exposed as bech32 interfaces. Many address are still in hex or base64 encoded form.

To covert between other binary reprsentation of addresses and keys, it is important to first apply the Amino enocoding process before bech32 encoding.


## Example

Once you create a new address, you should see the following:

`
NAME:	TYPE:	ADDRESS:						PUBKEY:
test1	local	faa18ekc4dswwrh2a6lfyev4tr25h5y76jkpqsz7kl	fap1addwnpepqgxa40ww28uy9q46gg48g6ulqdzwupyjcwfumgfjpvz7krmg5mrnw6zv8uv
`

This means you have created a new address `faa18ekc4dswwrh2a6lfyev4tr25h5y76jkpqsz7kl`, its hrp is `faa`. And its public key could be encoded into `fap1addwnpepqgxa40ww28uy9q46gg48g6ulqdzwupyjcwfumgfjpvz7krmg5mrnw6zv8uv`, its hrp is `fap`. 