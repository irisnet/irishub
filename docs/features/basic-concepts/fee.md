# Introduction

Specify the maximum fee you want to pay by `--fee`. Gas is the unit used to measure how much resource needed to execute the transaction. 
Specify the maximum gas limit by `--gas`. 
If maximum gas is too small, it won't be enough for executing the transaction. 
If the fee is too low, fee paid for each unit of gas will be less than gaslimit and 
validators won't execute the transaction neither. 
The fee(minimum unit)/gas must be large than 6*10^12. 
We recommend that you set your maximum gas to 50000 and set your maximum fee to 0.3iris. 
Fee will be consumed according to actual gas used and spare fee will be reimbursed to users.

## Fee

To secure their own validator node and maintain the avalibility of blockchain network, 
validators in IRISnet need a lot of equipments and resources. 
Thus, every transactions in IRISnet should pay some fee to validators. 
The parameter in commands is used to specify the maximum fee the user want to pay for their transaction.

## Gas

The resource needed for every transactions are varied for different type of transactions. For example, only a few computations, queries & modifies is needed for sending some token to other peolpe. But a lot of computations, queries & modifies is needed for creating a validator.  Gas is the unit used to measure how much resource needed to execute the transaction. We list the gas needed for some typical operactions in the below:

- gas needed  for writing the tx to the database is: 10 * the size of the transaction data (in bytes)
- gas needed for reading some data from database: 10 + data length(in bytes)
- gas needed for writing some data to database: 10 + 10 * data length(in bytes)
- sign or verify a signature

The total gas needed for executing the transaction is the sum of gas needed for every operations performed during executing the transaction. User should specify the maximum gas by --gas parameter. If maximum gas is not enough for executing the transcation, the transaction won't be executed successfully and all the fee will be returned to user's account. After the transaction being executed successfully, fee will deduct according to gas used, deducted fee will be  maximum fee * gas used / maximum gas and left fee will be returned to user. Gas price is equal to maximum fee / maximum gas and stands for fee that user want to pay for each unit of gas. To keep the fee is set in a resonable range, we set a minimum limit for gas price, 2^(-8) iris/gas, transaction won't be executed if the gas price is less than this value.

Example
```
    iriscli stake unbond  --from=test --shares-amount=10 --address-validator=iaa1mahw6ymzvt2q3lu4pjj5pau2e8krntklsj9mxe  --fee=0.3iris --gas=50000 --chain-id=<chain-id>
```
This example is a transaction to complete the unbond operation. The maximum fee(--fee) is set to be 0.3iris and the maximum(--gas) gas is set to be 50000. Therefore, the gas price here is 6000iris-nano/Gas. Suppose that 10000 gas is used to execute the transaction, then 0.06iris will be paid to validators and lefted 0.24iris will be returned to user.
