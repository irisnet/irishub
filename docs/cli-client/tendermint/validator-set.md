# iriscli tendermint validator-set

## Description

Get the full tendermint validator set at given height. If no height is specified, the latest height will be used as default.


## Usage

```
  iriscli tendermint validator-set [height] [flags]
```
or
```
  iriscli tendermint validator-set
```

## Flags

| Name, shorthand | Default                    |Description                                                             | Required     |
| --------------- | -------------------------- | --------------------------------------------------------- | -------- |
| --chain-id    |     | Chain ID of Tendermint node   | yes     |
| --node string     |   tcp://localhost:26657                         | Node to connect to (default "tcp://localhost:26657")  |                                     
| --help, -h      |       | 	help for block|    |
| --trust-node    |              true         | Trust connected full node (don't verify proofs for responses)     |          |

## Examples

### Get validator-set at height 114360

```shell
 iriscli tendermint validator-set 114360 --chain-id=irishub-test --trust-node=true

```

### Get the latest validator-set

```shell
 iriscli tendermint validator-set --chain-id=fuxi-4000 --trust-node=true

```
You will get the following result.

```apple js
{"block_height":"500","validators":[{"address":"fva1znj2x9p7jaww8x0a4ptcse4t68yezytkdjs6my","pub_key":"fvp1zcjduepqwh2pdw3gqstxjye9n2p9gp072e28qyrmpcegu2jg250r7k8y6naqw9epgu","proposer_priority":"-577","voting_power":"4100"},{"address":"fva1rhr9jqskza9und06mfhpdgdlm8q999yuzyhu70","pub_key":"fvp1zcjduepqjmq8r0zrqqpp2d99vlyuld0ga4qfju4uccaxrjwyqv5ykjx0p38sj5xpsm","proposer_priority":"-189","voting_power":"1000"},{"address":"fva19mqr37y2fq57xcpkxq95xrr37yjk7rchsktw73","pub_key":"fvp1zcjduepq5uqrykdrkg7tsr57kk58mjg530jf80zalujgc75y4a6g0uqzk85qra6rnq","proposer_priority":"-77","voting_power":"200"},{"address":"fva127nsk0yxt6843huqwr5sngsse3qdkehd3d0wzz","pub_key":"fvp1zcjduepqf0zsfyzvfreujl996g658tu59l8hy4x73epqnj53r7g7c3lqhazsw6tl5z","proposer_priority":"-1977","voting_power":"100"},{"address":"fva1exfvfnj63vesc5l9xzt6xg9ezfy593zq3m9cyz","pub_key":"fvp1zcjduepq02wagnhd8zcw2x5v68evlpvrthtc5ynkm6rtp8q3e6axw5ns7ylqk0497d","proposer_priority":"2823","voting_power":"300"}]}
```





