# 特殊权益用户

## 简介

IRISnet引入了两种由基金会控制且具有一定特殊权益的系统用户：profiler和trustee。

* Profiler的权益
    1. 通过治理提交软件升级/停止提议。
    2. 使用`profiling`模式发起服务调用，`profiling`模式会免除服务费。

* Trustee的权益
    1. 通过`CommunityTaxUsage`治理取回交易税费时，只能使用Trustee address作为取回地址。
    2. 发起`withdraw-tax`交易可以从`iService`服务费税池中提取代币到指定账户。

* Genesis Profiler/Genesis Trustee的权益（在创世的genesis.json中定义）
    1. 只有Genesis Profiler可以 添加/删除 普通Profiler账户
    2. 只有Genesis Trustee可以 添加/删除 Trustee账户

## 使用场景

1. 添加profiler和trustee

    添加profiler （仅限Genesis Profiler）

    ```bash
    iriscli guardian add-profiler --address=<profiler-address> --description=<profiler-description> --chain-id=<chain-id> --from=<key-name> --fee=0.3iris
    ```

    添加trustee（仅限Genesis Trustee）

    ```bash
    iriscli guardian add-trustee --address=<trustee-address> --description=<trustee-description> --chain-id=<chain-id> --from=<key-name> --fee=0.3iris 
    ```

2. 查询profiler和trustee列表

    查询profiler列表

    ```bash
    iriscli guardian profilers
    ```

    查询trustee列表

    ```bash
    iriscli guardian trustees
    ```

3. Profiler提交软件升级/停止提议

    详细参考[upgrade](upgrade.md)

4. Profiler使用`profiling`模式发起服务调用

    该模式免除服务费

    ```bash
    iriscli service call --def-chain-id=<def-chain-id> --service-name=<service-name> --method-id=<method-id> --bind-chain-id=<bind-chain-id> --provider=<provider-address> --service-fee=1iris --request-data=<request-data> --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --profiling=true
    ```

5. 通过`CommunityTaxUsage`治理取回交易税费

    详细参考[governance](governance.md#社区基金使用提议)

6. Trustee从`iService`服务费税池中提取代币到指定账户

    ```bash
    iriscli service withdraw-tax --dest-address=<destination-address> --withdraw-amount=1iris --chain-id=<chain-id> --from=<key-name> --fee=0.3iris
    ```

7. 删除profiler和trustee

    删除profiler：（仅限Genesis Profiler）

    ```bash
    iriscli guardian delete-profiler --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --address=<profiler-address>
    ```

    删除trustee：（仅限Genesis Trustee）

    ```bash
    iriscli guardian delete-trustee --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --address=<trustee-address>
    ```
