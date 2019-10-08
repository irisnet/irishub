# Coin Type

## Definitions

Coin type defines the available units of a kind of token in IRISnet. The developers can specify different coin-type for  their tokens. The native token in IRIShub is `iris`, which has following available units: `iris-milli`, `iris-micro`, `iris-nano`, `iris-pico`, `iris-femto` and `iris-atto`. The conversion relationship between them are as follows:

```toml
1 iris = 10^3 iris-milli
1 iris = 10^6 iris-micro
1 iris = 10^9 iris-nano
1 iris = 10^12 iris-pico
1 iris = 10^15 iris-femto
1 iris = 10^18 iris-atto
```

All the registered types of `iris` in the system can be used with transactions.

## Structure

### CoinType

```go
type CoinType struct {
    Name    string `json:"name"`
    MinUnit Unit   `json:"min_unit"`
    Units   Units  `json:"units"`
    Origin  Origin `json:"origin"`
    Desc    string `json:"desc"`
}
```

### Unit

```go
type Unit struct {
    Denom   string `json:"denom"`
    Decimal int    `json:"decimal"`
}
```

* `Name`: The name of a token, which is also its default unit；for instance,the default unit of IRISnet is `iris`.
* `MinUnit`: The minimum unit of coin-type. The tokens in the system are all stored in the form of minimum unit, such as `iris-atto`. You could choose to use the minimum unit of the tokens when sending a transaction to the IRIShub. If you use the command line client, aka `iriscli`, you can use any system-recognized unit and the system will automatically convert to the minimum unit of this corresponding token. For example, if you execute `send`command to transfer 1iris, the command line will be processed as 10^18 iris-atto in the backend, and you will only see 10^18 `iris-atto` when searching the transaction details by transaction hash.

`Denom`: is defined as the name of this unit, and `Decimal` is defined as the precision of the unit.

For example, the precision of iris-atto is 18.

* `Unit`: defines a set of units available under coin-type.
* `Origin`: defines the source of the coin-type
  * `Native`: native tokens, such as iris and user-defined tokens
  * `External`: external system tokens, such as eth for Ethereum, etc.
  * `Gateway`: external system tokens issued by gateways
* `Desc`: Description of the coin-type.

## Query of CoinType

If you want to query the CoinType configuration of a certain token, you can use the following command:

```bash
iriscli bank coin-type <coin-name>
```

E.g. query the CoinType of `iris`:

```bash
iriscli bank coin-type iris
```

Example Output:

```bash
CoinType:
  Name:     iris
  MinUnit:  iris-atto: 18
  Units:    iris: 0,  iris-milli: 3,  iris-micro: 6,  iris-nano: 9,  iris-pico: 12,  iris-femto: 15,  iris-atto: 18
  Origin:   native
  Desc:     IRIS Network
```
