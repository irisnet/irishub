# Fee Model

Specify the maximum fee you want to pay by `--fee`. Gas is the unit used to measure how many resources needed to execute the transaction. Specify the maximum gas limit by `-—gas`. If the maximum gas is too small, it won't be enough for executing the transaction. If the fee is too low, fee paid for each unit of gas will be less than the minimum gas-price and validators won't execute the transaction neither. The fee(minimum unit)/gas must be greater than 6*10^12 iris-atto. We recommend that you set your maximum gas to 50000 and set your maximum fee to 0.3iris. Fee will be consumed according to actual gas used and spare fee will be reimbursed to users.

:::tip
Certain types of transactions may consume more gas (e.g. creating validators), so it is recommended to set `--gas=100000 --fee=0.6iris` to ensure the smooth execution of transactions.
:::

## Fee

To secure their own validator node and maintain the availability of blockchain network, validators in IRISnet need a lot of equipments and resources. Thus, every transactions in IRISnet should pay some fee to validators. The parameter in commands is used to specify the maximum fee the user wants to pay for their transaction.

## Gas

The resources needed for every transaction vary by different type of transactions. For example, only a few computations, queries & modifies are needed for the transfer transaction. But more computations, queries & modifies are needed for creating a validator. Gas is the unit used to measure how many resources needed to execute the transaction. We list the gas needed for some typical operations in the below:

| operations                          | gas needed                                                   |
| ----------------------------------- | ------------------------------------------------------------ |
| writing the tx to the database      | 10 * the size of the transaction data (in bytes)             |
| for reading some data from database | 10 + ln(data_length)/ln(1.02), data_length is the length of the read data(in bytes) |
| writing some data to database       | 10 + 10 * ln(data_length)/ln(1.02), data_length is the length of the write data(in bytes) |
| signature verification              | 100                                                          |

The total gas needed for executing the transaction is the sum of gas needed for every operation performed during executing the transaction. The user should specify the maximum gas by `--gas` parameter. If the maximum gas is not enough for executing the transaction, the transaction won't be executed successfully and all the fee will be returned to the user's account. After the transaction being executed successfully, fee will deduct according to gas used. The deducted fee will be `maximum-fee * gas-used / maximum-gas` and left fee will be returned to the user. Gas price is equal to `maximum-fee / maximum-gas` and stands for fee that the user wants to pay for each unit of gas. To keep the fee set in a reasonable range, we set a minimum limit for gas price: `6*10^(-6) iris/gas`, the transaction won't be executed if the gas price is less than this value.

Example

```bash
iriscli bank send --amount=1iris --fee=0.3iris --gas=50000 --chain-id=<chain-id> --from=<key-name> --to=<account-address>
```

This example is a transfer transaction. The maximum fee `--fee` is set to be 0.3iris and the maximum gas `--gas` is set to be 50000. Therefore, the gas price here is 6000iris-nano/Gas. Suppose that 10000 gas is used to execute the transaction, then 0.06iris will be paid to validators and left 0.24iris will be refunded to user.
