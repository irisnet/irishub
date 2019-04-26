# 委托人

## 什么是委托人?

不能或不想运行验证人节点的用户仍然可以作为委托人参与委托过程。事实上，节点被选为验证人不是根据验证人本身自抵押额来选择的，而是根据这个验证人上所有委托的抵押额来被选择的。这是一个重要的属性，因为它使委托人能够制约验证人以防他们出现不良行为。 

如果验证人有不良行为，他的委托人将把他们的iris通证转移，从而减少他的总股份。最后，如果验证人的总委托额跌出了前100名，那么他将被退出验证人集合。

## 委托人抵押的通证

委托人与验证人有相同的状态。请注意，委托的token并不一定是绑定的。它的状态可以为"委托且绑定"、"委托且正在解除绑定"、"委托且已经解除绑定"或"释放"。 

## 委托人常用操作

* 委托

委托10iris给某个验证人：
```$xslt
iriscli stake delegate --address-delegator=<address-delegator> --address-validator=<address-validator> --chain-id=<chain-id> --from=<key_name> --fee=0.3iris --amount=10iris
```
请参阅[委托](../cli-client/stake/delegate.md)

* 查询委托
```$xslt
iriscli stake delegation --address-delegator=<address-delegator> --address-validator=<address-validator> --chain-id=<chain-id> 
```
请参阅[查询委托](../cli-client/stake/delegation.md)

* 转委托

委托人可以随时更换受委托验证人，当更换受委托验证人时，委托人抵押的通证可直接转入新的受委托人抵押池，而不无需等待3周的解绑期。
 
转委托100shares给新的验证人：
```$xslt
iriscli stake redelegate --addr-validator-dest=<addr-validator-dest>  --addr-validator-source=<addr-validator> --address-delegator=<address-delegator>  --chain-id=<chain-id>  --from=<key_name> --fee=0.3iris --shares-amount=100 
```

请参阅[转委托](../cli-client/stake/redelegate.md)

* 解绑

若委托人需要取回已委托的通证，可以通过发送[解绑交易](../cli-client/stake/unbond.md)。在IRISnet网络中，解绑期默认为**三周**(解绑期间无收益)。一旦解绑期结束，被绑定的通证将自动成为流通通证。

解绑50% shares：
```$xslt
iriscli stake unbond begin  --address-validator=<address-validator> --address-delegator=<address-delegator> --chain-id=<chain-id>  --from=<key_name> --fee=0.3iris --shares-percent=0.5
```

三周后你将会看到，账户余额增加。
