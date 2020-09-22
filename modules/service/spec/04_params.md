<!--
order: 4
-->

# Parameters

The service module contains the following parameters:

| Key                  | Type             | Example                              |
| -------------------- | ---------------- | ------------------------------------ |
| MaxRequestTimeout    | int64            | 100                                  |
| MinDepositMultiple   | int64            | 1000                                 |
| MinDeposit           | array (coins)    | [{"denom":"stake","amount":"10000"}] |
| ServiceFeeTax        | string (dec)     | "0.01"                               |
| SlashFraction        | string (dec)     | "0.001"                              |
| ComplaintRetrospect  | string (time ns) | "172800000000000"                    |
| ArbitrationTimeLimit | string (time ns) | "172800000000000"                    |
| TxSizeLimit          | uint64           | 4000                                 |
| BaseDenom            | string           | "stake"                              |
