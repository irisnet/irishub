# 介绍 

这里主要介绍distribution模块对于的命令行接口。

# 查询接口

在默认情况下，所以查询命令处于信任模式，也就是在查询的时候不要求全节点返回查询结果对于的proof，并且即使全节点返回proof也不会对其进行验证。如果用户不信任所连接的全节点，可以通过在查询命令后天就`--trust-node=false`来使能非信任模式。

1. 查询撤回地址

    示例命令:
    ```bash
    iriscli distribution withdraw-address faa1vm068fnjx28zv7k9kd9j85wrwhjn8vfsxfmcrz
    ```
    如果委托人没有指定过撤回地址，那么查询结果为空。

2. 查询委托(delegation)的收益分配记录

    示例命令:
    ```bash
    iriscli distribution delegation-distr-info --address-delegator=faa1ezzh0humhy3329xg4avhcjtay985nll0zswc5j \
    --address-validator=fva1ezzh0humhy3329xg4avhcjtay985nll0hpyhf4
    ```
    示例查询结果：
    ```json
    {
      "delegator_addr": "faa1ezzh0humhy3329xg4avhcjtay985nll0zswc5j",
      "val_operator_addr": "fva1ezzh0humhy3329xg4avhcjtay985nll0hpyhf4",
      "del_pool_withdrawal_height": "4044"
    }
    ```
    从这个查询结果可知，这个委托是在4044高度创建的，或者在4044高度发起过撤回交易。

2. 查询委托人所有的委托(delegation)的收益分配记录

    示例命令: 
    ```bash
    iriscli distribution delegator-distr-info faa1ezzh0humhy3329xg4avhcjtay985nll0zswc5j
    ```
    示例查询结果：
    ```json
    [
      {
        "delegator_addr": "faa1ezzh0humhy3329xg4avhcjtay985nll0zswc5j",
        "val_operator_addr": "fva14a70gzu0v2w8dlfx462c9sldvja24qaz6vv4sg",
        "del_pool_withdrawal_height": "10859"
      },
      {
        "delegator_addr": "faa1ezzh0humhy3329xg4avhcjtay985nll0zswc5j",
        "val_operator_addr": "fva1ezzh0humhy3329xg4avhcjtay985nll0hpyhf4",
        "del_pool_withdrawal_height": "4044"
      }
    ]
    ```

4. 查询验证人收益分配记录

    示例命令: 
    ```bash
    iriscli distribution delegator-distr-info faa1ezzh0humhy3329xg4avhcjtay985nll0zswc5j
    ```
    示例查询结果：
    ```json
    {
      "operator_addr": "faa1ezzh0humhy3329xg4avhcjtay985nll0zswc5j",
      "fee_pool_withdrawal_height": "416",
      "del_accum": {
        "update_height": "416",
        "accum": "10000.0000000000"
      },
      "del_pool": "10.1300000000iris",
      "val_commission": "1.2345000000iris"
    }
    ```
	从这个查询结果可知，这个委托是在4044高度创建的，或者在4044高度发起过撤回交易。

# 发交易接口

1. 设置收益收款地址

    示例命令: 
    ```bash
    iriscli distribution set-withdraw-addr faa1syva9fvh8m6dc6wjnjedah64mmpq7rwwz6nj0k --from mykey --fee=0.004iris --chain-id=irishub-test
    ```
    
2. withdraw rewards 

    1. 仅撤回某一个委托产生的收益
    ```bash
    iriscli distribution withdraw-rewards --only-from-validator fva134mhjjyyc7mehvaay0f3d4hj8qx3ee3w3eq5nq --from mykey --fee=0.004iris --chain-id=irishub-test
    ```
    2. 仅撤回某所有委托产生的收益
    ```bash
    iriscli distribution withdraw-rewards --from mykey --fee=0.004iris --chain-id=irishub-test
    ```
    3. 仅撤回某所有委托产生的收益和验证人的佣金收益
    ```bash
    iriscli distribution withdraw-rewards --is-validator=true --from mykey --fee=0.004iris --chain-id=irishub-test
    ```