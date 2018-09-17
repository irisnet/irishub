# Fee Token specification

## Previous implementation in Cosmos-SDK

* There is no any constrain to fee token, just ensure gas is enough. Gas has no explicit connection with fee token.
* All fee token will be deducted

## Current implementation in Cosmos-SDK

* In global parameter store, the native fee token name has been specified. In iris network, the native fee token is iris. Other tokens will be ignored. The native fee token can't be modified.
* Transactions should specify the max gas (gasLimit) to be consumed in transaction execution. If the gas is exhausted, the transaction execution will be terminated and return error. No fee will be charged.
* gasPrice = feeToken / gasLimit. gasPrice should be no less than gasPriceThreshold (default value: 2*10^(-8)iris, iris precision: 10^18). Otherwise the transactions will be considered as invalid. gasPriceThreshold is a global parameter which can be modified by governance.
* Normally, the fee token will not be deducted totally. The deducted quantity depends on how many the gas is consumed. If half of gas is consumed, then only half of fee token will be deducted.
 
## How to specify fee for transactions

Now we must specify fee for each transaction properly. Otherwise, your transaction will be rejected. Here we take sending token for example. Now we can send token by command line and rest API:

For command Line:
```
iriscli send --to=faa1ze87hkku6vgnqp7gqkn528d60ttr06xuudx0k9 --name=moniker --chain-id=test-chain-UieBV2 --amount=10iris --node=tcp://localhost:26657 \
--gas=10000 --fee=0.0005iris
```
In command line, we specify gas and fee token by two option: "--gas" and "--fee". The "--gas" option can be omitted. Because gas has default value: 200000. However, the "--fee" option can't be omitted. Because its default value is empty.
   
For rest API:
```
{
	"amount":[{"denom":"iris","amount":"10"}],
	"name":"moniker",
	"password":"1234567890",
	"chain_id":"test-chain-UieBV2",
	"sequence":"9",
	"account_number":"0",
	"gas":"10000",
	"fee":"0.0005iris"
}
```
In rest API, the gas and fee token are specified by gas field and fee field in json body. Both of them can't be omitted. 

The transaction senders should ensure the gas is enough for transaction execution and the gasPrice is no less than gasPriceThreshold. In addition, only specify iris as fee token. Other token will be ignored. 

## Suggestions for test

For both command line and rest API:

* Try to specify no fee.
* Try to specify other token instead of iris token as fee.
* Try to specify few iris token.
* Try to specify more gas so that the gas price is very low.
* Try different kinds of transaction.

## Future improvement

* Fee token whitelist will be brought in. All tokens in the whitelist can be used as fee token. The whitelist can be modified by governance.
* Currently, the transaction senders have no motivation to pay more fee, just ensure the gasPrice is no less than the gasPriceThreshold. Next new mechanisms will be brought in to encourage users to pay more fee. For example, proposals will tender to include transactions with higher gasPrice. Transactions with lower gasPrice have to wait for more time.
* When a proposal builds a block, it should check whether the sum of all transaction gas is greater than blockGasLimit. If so, remove some transactions the accommodate this rule. With blockGasLimit, transactions which will consumed too much resource, such as memory , disk space and execution time, will be rejected. These special transaction may cause crash or consensus failure.
* Consider to charge fee even if the transactions encounter execution failure. This is helpful for defencing DDos attack.
cfgtyhdryd
