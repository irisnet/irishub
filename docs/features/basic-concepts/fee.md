# Introduction

Specify the maximum fee you want to pay by `--fee`. Gas is the unit used to measure how much resource needed to execute the transaction. Specify the maximum gas limit by --gas. If maximum gas is too small, it won't be enough for executing the transaction. If the fee is too low, fee paid for each unit of gas will be low & validator won't execute the transaction too. The fee(minimum unit)/gas must be large than 2*10^10. We recommend that you set your maximum gas to 20000 and set your maximum fee to 400000000000000iris. Fee will be deduct according to the gas used and lefted fee will be returned to user.

## Fee

To secure their own validator node and maintain the avalibility of blockchain network, validators in IRISnet need a lot of equipments and works. Thus, every transactions in IRISnet should pay some fee to validators. The parameter in commands is used to specify the maximum fee the user want to pay for their transaction.

## Gas

The resource needed for every transactions are varied for different type of transactions. For example, only a few computations, queries & modifies is needed for sending some token to other peolpe. But a lot of computations, queries & modifies is needed for creating a validator.  Gas is the unit used to measure how much resource needed to execute the transaction. We list the gas needed for some typical operactions in the below:

- gas needed  for writing the tx to the database is: 10 * the size of the transaction data (in bytes)
- gas needed for reading some data from database: 10 + data length(in bytes)
- gas needed for writing some data to database: 10 + 10 * data length(in bytes)
- sign or verify a signature

The total gas needed for executing the transaction is the sum of gas needed for every operations performed during executing the transaction. User should specify the maximum gas by --gas parameter. If maximum gas is not enough for executing the transcation, the transaction won't be executed successfully and all the fee will be returned to user's account. After the transaction being executed successfully, fee will deduct according to gas used, deducted fee will be  maximum fee * gas used / maximum gas and left fee will be returned to user. Gas price is equal to maximum fee / maximum gas and stands for fee that user want to pay for each unit of gas. To keep the fee is set in a resonable range, we set a minimum limit for gas price, 2^(-8) iris/gas, transaction won't be executed if the gas price is less than this value.

Example
```
    iriscli stake unbond  --from=test --address-validator=faa1mahw6ymzvt2q3lu4pjj5pau2e8krntklgarrxy --address-delegator=faa1mahw6ymzvt2q3lu4pjj5pau2e8krntklgarrxy  --fee=0.02iris --gas=20000 --chain-id=test-irishub
```
This example is a transaction to complete the unbond operation. The maximum fee(--fee) is set to be 2000000000000000iris(2*10^15) and the maximum(--gas) gas is set to be 20000. Therefore, the gas price here is 10^11iris/gas. Suppose that 1500 gas is used to execute the transaction, then 1500000000000000 iris will be paid to validators and lefted 500000000000000 iris will be returned to user.
