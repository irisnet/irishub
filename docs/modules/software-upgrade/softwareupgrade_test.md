# 软件升级测试

## 单节点成功升级CASE
### 第一次升级
#### 运行旧软件
```
rm -rf iris
rm -rf .iriscli
iris init gen-tx --name=x --home=iris
iris init --gen-txs --chain-id=upgrade-test -o --home=iris
iris start --home=iris
```

```
iriscli keys list
VADDR=验证人的地址
```

#### 发送升级提议
```
iriscli gov submit-proposal --name=x --proposer=$VADDR --title=ADD --description="I am crazy" --type=SoftwareUpgrade --deposit=10000000000000000000iris --chain-id=upgrade-test --fee=20000000000000000iris
```
#### 发送升级协议的YES投票
```
iriscli gov  vote --name=x --voter=$VADDR --proposalID=1 --option=Yes --chain-id=upgrade-test --fee=20000000000000000iris
```

#### 查询提议内容
```
iriscli gov query-proposal --proposalID=1        
```

#### 查询升级的版本信息
```
iriscli upgrade info
```

#### 运行新软件
```
iris1 start --home=iris
```

#### 发送消息自己已运行新软件
```
iriscli  upgrade submit-switch --name=x --from=$VADDR --proposalID=1 --chain-id=upgrade-test --fee=20000000000000000iris
```

#### 查询switch信息
```
iriscli upgrade query-switch --voter=$VADDR --proposalID=1
```

#### 使用新功能（无报错）
```
iriscli1 advanced ibc set --name=x --from=$VADDR --chain-id=upgrade-test --print-response true --fee=20000000000000000iris
iriscli1 advanced ibc get --name=x --from=$VADDR --chain-id=upgrade-test --sequence=0 --print-response true --fee=20000000000000000iris
```

### 第二次升级
#### 发送升级提议第二次升级
```
iriscli1 gov submit-proposal --name=x --proposer=$VADDR --title=ADD --description=“I am crazy” --type=SoftwareUpgrade --deposit=10000000000000000000iris --chain-id=upgrade-test --fee=20000000000000000iris
```
#### 发送升级协议的YES投票
```
iriscli gov  vote --name=x --voter=$VADDR --proposalID=2 --option=Yes --chain-id=upgrade-test --fee=20000000000000000iris
```

#### 查询提议内容
```
iriscli gov query-proposal --proposalID=2        
```

#### 查询升级的版本信息
```
iriscli upgrade info
```

#### 运行新软件
```
iris2 start --home=iris
```

#### 发送消息自己已运行新软件
```
iriscli2  upgrade submit-switch --name=x --from=$VADDR --proposalID=2 --chain-id=upgrade-test --fee=20000000000000000iris
```

#### 查询switch信息
```
iriscli2 upgrade query-switch --voter=$VADDR --proposalID=1
```

#### 使用新功能（无报错）
```
iriscli2 advanced ibc set --name=x --from=$VADDR --chain-id=upgrade-test --print-response true --fee=20000000000000000iris
iriscli2 advanced ibc get --name=x --from=$VADDR --chain-id=upgrade-test --print-response true --fee=20000000000000000iris
```

## 多节点连续升级测试

### version 0

```
iris start --home /data/iris > /data/iris/log.txt &
(run in all the nodes)

iriscli gov submit-proposal --name=silei --proposer=$VADDR --title=ADD --description="I am crazy" --type=Text --deposit=10iris --chain-id=upgrade-test --fee=20000000000000000iris --home=/data/iriscli
(run in node1)

iriscli gov query-proposal --proposalID=1

```

### version 1

```

kill iris   (run in all the nodes)

iris1 start --home /data/iris > /data/iris/log.txt &   (run in all the nodes)

iriscli1 upgrade submit-switch --name=silei --from=$VADDR --proposalID=1 --title=test --chain-id=upgrade-test --fee=20000000000000000iris --home=/data/iriscli
(run in all the nodes)

iriscli1 upgrade info

iriscli1 advanced ibc set --name=silei --from=$VADDR --chain-id=upgrade-test --sequence=0 --print-response true --fee=20000000000000000iris --home=/data/iriscli

iriscli1 advanced ibc get --name=silei --from=$VADDR --chain-id=upgrade-test --sequence=0 --print-response true --fee=20000000000000000iris --home=/data/iriscli

iriscli1 gov submit-proposal --name=silei --proposer=$VADDR --title=ADD --description="I am crazy" --type=Text --deposit=10iris --chain-id=upgrade-test --fee=20000000000000000iris --home=/data/iriscli
(run in node1)

iriscli1 gov query-proposal --proposalID=2

```

### version 2

```

kill iris1   (run in all the nodes)

iris2 start --home /data/iris > /data/iris/log.txt &   (run in all the nodes)

iriscli2 upgrade submit-switch --name=silei --from=$VADDR --proposalID=1 --title=test --chain-id=upgrade-test --fee=20000000000000000iris --home=/data/iriscli
(run in all the nodes)

iriscli2 upgrade info

iriscli2 advanced ibc set --name=silei --from=$VADDR --chain-id=upgrade-test --sequence=0 --print-response true --fee=20000000000000000iris --home=/data/iriscli

iriscli2 advanced ibc get --name=silei --from=$VADDR --chain-id=upgrade-test --sequence=0 --print-response true --fee=20000000000000000iris --home=/data/iriscli

```

### bug fix upgrade from version 1

```

iriscli1 gov submit-proposal --name=silei --proposer=$VADDR --title=ADD --description="I am crazy" --type=SoftwareUpgrade --deposit=10000000000000000000iris --chain-id=upgrade-test --fee=20000000000000000iris --home=/data/iriscli

iriscli1 gov vote --name=silei --voter=$VADDR --proposalID=2 --option=Yes --chain-id=upgrade-test --fee=20000000000000000iris --home=/data/iriscli

iriscli1 gov query-proposal --proposalID=2

iris2-bugfix start --home /data/iris


iriscli2-bugfix upgrade submit-switch --name=silei --from=$VADDR --proposalID=2 --title=test --chain-id=upgrade-test --fee=20000000000000000iris --home=/data/iriscli

iriscli2-bugfix upgrade query-switch --voter=$VADDR --proposalID=3 --home=/data/iriscli

iriscli2-bugfix upgrade info

iriscli2-bugfix advanced ibc set --name=silei --from=$VADDR --chain-id=upgrade-test --sequence=0 --print-response true --fee=20000000000000000iris --home=/data/iriscli

```

## 多节点升级测试(非docker)
```
rm -rf iris1
rm -rf iris2
rm -rf iris3
rm -rf iris4
rm -rf .iriscli

iris init gen-tx --name=x1 --home=iris1
iris init gen-tx --name=x2 --home=iris2
iris init gen-tx --name=x3 --home=iris3
iris init gen-tx --name=x4 --home=iris4

cp iris2/config/gentx/gentx-XXXX.json iris1/config/gentx/
cp iris3/config/gentx/gentx-XXXX.json iris1/config/gentx/
cp iris4/config/gentx/gentx-XXXX.json iris1/config/gentx/

iris init --gen-txs --chain-id=upgrade-test -o --home=iris1

cp ~/iris1/config/genesis.json ~/iris2/config/
cp ~/iris1/config/genesis.json ~/iris3/config/
cp ~/iris1/config/genesis.json ~/iris4/config/

vi iris2/config/config.toml
vi iris3/config/config.toml
vi iris4/config/config.toml
6628995f6eae0c7d810867e467f23530c55b1232@localhost:26656

iris start --home=iris1
iris start --home=iris2
iris start --home=iris3
iris start --home=iris4
```

### 地址赋值
```
VADDR1=
VADDR2=
VADDR3=
VADDR4=
```

### 发送升级提议
```
iriscli gov submit-proposal --name=x1 --proposer=$VADDR1 --title=ADD --description=“I am crazy” --type=SoftwareUpgrade --deposit=10000000000000000000iris --chain-id=upgrade-test --fee=20000000000000000iris
```
### 投票
```
iriscli gov  vote --name=x1 --voter=$VADDR1 --proposalID=1 --option=Yes --chain-id=upgrade-test --fee=20000000000000000iris
iriscli gov  vote --name=x2 --voter=$VADDR2 --proposalID=1 --option=Yes --chain-id=upgrade-test --fee=20000000000000000iris
iriscli gov  vote --name=x3 --voter=$VADDR3 --proposalID=1 --option=Yes --chain-id=upgrade-test --fee=20000000000000000iris
iriscli gov  vote --name=x4 --voter=$VADDR4 --proposalID=1 --option=Yes --chain-id=upgrade-test --fee=20000000000000000iris

```
### 查询提议内容
```
iriscli gov query-proposal --proposalID=1        
```

### 查询升级的版本信息
```
iriscli upgrade info
```

#### 运行新软件
```
iris1 start --home=iris1
iris1 start --home=iris2
iris1 start --home=iris3
iris1 start --home=iris4
```

### 发送消息自己已运行新软件
```
iriscli1  upgrade submit-switch --name=x1 --from=$VADDR1 --proposalID=1 --chain-id=upgrade-test --fee=20000000000000000iris
iriscli1  upgrade submit-switch --name=x2 --from=$VADDR2 --proposalID=1 --chain-id=upgrade-test --fee=20000000000000000iris
iriscli1  upgrade submit-switch --name=x3 --from=$VADDR3 --proposalID=1 --chain-id=upgrade-test --fee=20000000000000000iris
iriscli1  upgrade submit-switch --name=x4 --from=$VADDR4 --proposalID=1 --chain-id=upgrade-test --fee=20000000000000000iris
```

### 查询switch信息
```
iriscli1 upgrade query-switch --voter=$VADDR1 --proposalID=1
iriscli1 upgrade query-switch --voter=$VADDR2 --proposalID=1
iriscli1 upgrade query-switch --voter=$VADDR3 --proposalID=1
iriscli1 upgrade query-switch --voter=$VADDR4 --proposalID=1
```

### 测试新功能
```
iriscli1 advanced ibc set --name=x --from=$VADDR --chain-id=upgrade-test --sequence=0 --print-response true
iriscli1 advanced ibc get --name=x --from=$VADDR --chain-id=upgrade-test --sequence=0 --print-response true
```

### 发送升级提议
```
iriscli1 gov submit-proposal --name=x1 --proposer=$VADDR1 --title=ADD --description=“I am crazy” --type=SoftwareUpgrade --deposit=10000000000000000000iris --chain-id=upgrade-test --fee=20000000000000000iris
```
### 投票
```
iriscli1 gov  vote --name=x1 --voter=$VADDR1 --proposalID=2 --option=Yes --chain-id=upgrade-test --fee=20000000000000000iris
iriscli1 gov  vote --name=x2 --voter=$VADDR2 --proposalID=2 --option=Yes --chain-id=upgrade-test --fee=20000000000000000iris
iriscli1 gov  vote --name=x3 --voter=$VADDR3 --proposalID=2 --option=Yes --chain-id=upgrade-test --fee=20000000000000000iris
iriscli1 gov  vote --name=x4 --voter=$VADDR4 --proposalID=2 --option=Yes --chain-id=upgrade-test --fee=20000000000000000iris

```
### 查询提议内容
```
iriscli gov query-proposal --proposalID=2      
```

### 查询升级的版本信息
```
iriscli upgrade info
```

#### 运行新软件
```
iris2 start --home=iris1
iris2 start --home=iris2
iris2 start --home=iris3
iris2 start --home=iris4
```

### 发送消息自己已运行新软件
```
iriscli2  upgrade submit-switch --name=x1 --from=$VADDR1 --proposalID=2 --chain-id=upgrade-test --fee=20000000000000000iris
iriscli2  upgrade submit-switch --name=x2 --from=$VADDR2 --proposalID=2 --chain-id=upgrade-test --fee=20000000000000000iris
iriscli2  upgrade submit-switch --name=x3 --from=$VADDR3 --proposalID=2 --chain-id=upgrade-test --fee=20000000000000000iris
iriscli2  upgrade submit-switch --name=x4 --from=$VADDR4 --proposalID=2 --chain-id=upgrade-test --fee=20000000000000000iris
```

### 查询switch信息
```
iriscli2 upgrade query-switch --voter=$VADDR1 --proposalID=1
iriscli2 upgrade query-switch --voter=$VADDR2 --proposalID=1
iriscli2 upgrade query-switch --voter=$VADDR3 --proposalID=1
iriscli2 upgrade query-switch --voter=$VADDR4 --proposalID=1
```

### 测试新功能
```
iriscli2 advanced ibc set --name=x --from=$VADDR --chain-id=upgrade-test --sequence=0 --print-response true
iriscli2 advanced ibc get --name=x --from=$VADDR --chain-id=upgrade-test --sequence=0 --print-response true
```

## iris测试命令 
```
rm -rf .iris*
rm -rf iris*
iris init gen-tx --name=x --home=iris1
iris init --gen-txs --chain-id=upgrade-test -o --home=iris1
iris start --home=iris1
iriscli keys list
iriscli gov  submit-proposal --name=x --proposer=$VADDR --title=ADD --description=“I am crazy” --type=Text --deposit=10iris --chain-id=upgrade-test 
iriscli gov deposit --depositer=$VADDR  --deposit=1iris --name=x --proposalID=1 --chain-id=upgrade-test
iriscli gov query-proposal --proposalID=1
iriscli upgrade info
iriscli  upgrade submit-switch --name=x --from=$VADDR --proposalID=1 --chain-id=upgrade-test
iriscli upgrade query-switch --voter=$VADDR --proposalID=1
```

```
iris unsafe_reset_all --home=iris1
iris unsafe_reset_all --home=iris2
iris unsafe_reset_all --home=iris3
iris unsafe_reset_all --home=iris4
```