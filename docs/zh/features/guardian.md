# Guardian User Guide

## 简介
IRISnet引入了两种由基金会控制的特权系统用户，profiler和trustee。 

* Profiler的特权[TODO]
    1. 通过治理提交软件升级/停止提议。
    2. 使用profiling模式发起服务调用，profiling模式会免除服务费。
    
* Trustee的特权[TODO]
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