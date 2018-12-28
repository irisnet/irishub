# Guardian User Guide

## 简介
IRISnet引入了两种由基金会控制的特权系统用户，profiler和trustee。 

* Profiler的特权
    1. 通过治理提交软件升级/停止提议。
    2. 使用profiling模式发起服务调用，profiling模式会免除服务费。
    
* Trustee的特权
    1. 如果`TxTaxUsage`提议的使用类型是`Distribution`或`Grant`，作为目的地地址。
    2. 发送交易从系统服务费税池中提取代币到账户。
    
## 使用场景
1. 添加profiler

    只有profiler能添加新的profiler
    ```shell
    iriscli guardian add-profiler --profiler-address=[profiler address] --profiler-name=[name] --chain-id=[chain-id] --from=[key name] --fee=0.004iris 
    ```
    
2. 查询profiler和trustee列表

    查询profiler列表
    ```shell
    iriscli guardian profilers
    ```
    查询trustee列表
    ```shell
    iriscli guardian trustees
    ```
    
3. Profiler提交软件升级/停止提议

    详细参考[upgrade](upgrade.md)

4. Profiler使用profiling模式发起服务调用
    ```shell
    iriscli service call --def-chain-id=[def-chain-id] --service-name=[service-name] --method-id=[method-id] --bind-chain-id=[bind-chain-id] --provider=[provider address] --service-fee=1iris --request-data=[request-data] --chain-id=[chain-id] --from=[key name] --fee=0.004iris
    ```
    
5. Trustee作为`TxTaxUsage`提议的目的地地址

    详细参考[governance](governance.md#proposals-on-transaction-fee-community-tax-usage)
    
6. Trustee取回服务费税金
    ```shell
    iriscli service withdraw-tax --dest-address=[destination address] --withdraw-amount=1iris --chain-id=[chain-id] --from=[key name] --fee=0.004iris 
    ```