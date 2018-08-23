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

## 多节点连续升级测试 (以4个节点为例)

### version 0 （运行0号版本，并提交软件升级提议）

在4个节点分别运行iris：
```
iris start --home /data/iris > /data/iris/log.txt &
```
1号节点发起软件升级提议：
```
iriscli gov submit-proposal --name=silei --proposer=$VADDR --title=ADD --description="I am crazy" --type=Text --deposit=10iris --chain-id=upgrade-test --fee=20000000000000000iris --home=/data/iriscli
```
各节点投票通过该提议：
```
iriscli gov vote --name=silei --voter=$VADDR --proposalID=1 --option=Yes --chain-id=upgrade-test --fee=20000000000000000iris --home=/data/iriscli
```
查询软件升级提议：
```
iriscli gov query-proposal --proposalID=1
```


### version 1 (升级到1号版本)

各节点退出iris，下载并启动iris1
```
kill iris

iris1 start --home /data/iris > /data/iris/log.txt & 
```
各节点发送切换到新版本消息：
```
iriscli1 upgrade submit-switch --name=silei --from=$VADDR --proposalID=1 --title=test --chain-id=upgrade-test --fee=20000000000000000iris --home=/data/iriscli
```
查询软件升级信息：
```
iriscli1 upgrade info
```

运行新版本特有命令，检查升级结果：
```
iriscli1 advanced ibc set --name=silei --from=$VADDR --chain-id=upgrade-test --sequence=0 --print-response true --fee=20000000000000000iris --home=/data/iriscli

iriscli1 advanced ibc get --name=silei --from=$VADDR --chain-id=upgrade-test --sequence=0 --print-response true --fee=20000000000000000iris --home=/data/iriscli
```

1号节点发起软件升级提议, 提议升级到2号版本：
```
iriscli1 gov submit-proposal --name=silei --proposer=$VADDR --title=ADD --description="I am crazy" --type=Text --deposit=10iris --chain-id=upgrade-test --fee=20000000000000000iris --home=/data/iriscli

iriscli1 gov query-proposal --proposalID=2

```

各节点投票通过该提议：
```
iriscli1 gov vote --name=silei --voter=$VADDR --proposalID=2 --option=Yes --chain-id=upgrade-test --fee=20000000000000000iris --home=/data/iriscli
```

### version 2  (从1号版本升级到2号版本)

各节点退出iris1，下载并启动iris2
```
kill iris1 

iris2 start --home /data/iris > /data/iris/log.txt & 
```
各节点发送切换到新版本消息：
```
iriscli2 upgrade submit-switch --name=silei --from=$VADDR --proposalID=1 --title=test --chain-id=upgrade-test --fee=20000000000000000iris --home=/data/iriscli

iriscli2 upgrade info
```
运行新版本特有命令，检查升级结果：
```
iriscli2 advanced ibc set --name=silei --from=$VADDR --chain-id=upgrade-test --sequence=0 --print-response true --fee=20000000000000000iris --home=/data/iriscli

iriscli2 advanced ibc get --name=silei --from=$VADDR --chain-id=upgrade-test --sequence=0 --print-response true --fee=20000000000000000iris --home=/data/iriscli

```

### bug fix upgrade from version 1 (基于1号版本进行bug fix的软件升级)

1号节点发起软件升级提议, 提议进行bug-fix升级：
```
iriscli1 gov submit-proposal --name=silei --proposer=$VADDR --title=ADD --description="I am crazy" --type=SoftwareUpgrade --deposit=10000000000000000000iris --chain-id=upgrade-test --fee=20000000000000000iris --home=/data/iriscli
```

各节点投票通过该提议：
```
iriscli1 gov vote --name=silei --voter=$VADDR --proposalID=2 --option=Yes --chain-id=upgrade-test --fee=20000000000000000iris --home=/data/iriscli

iriscli1 gov query-proposal --proposalID=2
```

各节点退出iris1，下载并启动iris2-bugfix:
```
iris2-bugfix start --home /data/iris
```
各节点发送切换到新版本消息：
```
iriscli2-bugfix upgrade submit-switch --name=silei --from=$VADDR --proposalID=2 --title=test --chain-id=upgrade-test --fee=20000000000000000iris --home=/data/iriscli

iriscli2-bugfix upgrade query-switch --voter=$VADDR --proposalID=3 --home=/data/iriscli

iriscli2-bugfix upgrade info
```
运行新版本特有命令，检查升级结果：
```
iriscli2-bugfix advanced ibc set --name=silei --from=$VADDR --chain-id=upgrade-test --sequence=0 --print-response true --fee=20000000000000000iris --home=/data/iriscli
```

### bug fix upgrade from version iris2-bugfix (基于iris2-bugfix版本继续进行bug-fix的软件升级)

1号节点发起软件升级提议, 提议进行bug-fix升级：
```
iriscli2-bugfix gov submit-proposal --name=silei --proposer=$VADDR --title=ADD --description="I am crazy" --type=SoftwareUpgrade --deposit=10000000000000000000iris --chain-id=upgrade-test --fee=20000000000000000iris --home=/data/iriscli
```

各节点投票通过该提议：
```
iriscli2-bugfix gov vote --name=silei --voter=$VADDR --proposalID=3 --option=Yes --chain-id=upgrade-test --fee=20000000000000000iris --home=/data/iriscli

iriscli2-bugfix gov query-proposal --proposalID=3
```

各节点退出iris2-bugfix，下载并启动iris3-bugfix:
```
iris3-bugfix start --home /data/iris
```
各节点发送切换到新版本消息：
```
iriscli3-bugfix upgrade submit-switch --name=silei --from=$VADDR --proposalID=3 --title=test --chain-id=upgrade-test --fee=20000000000000000iris --home=/data/iriscli

iriscli3-bugfix upgrade query-switch --voter=$VADDR --proposalID=3 --home=/data/iriscli

iriscli3-bugfix upgrade info
```
运行新版本特有命令，检查升级结果：
```
iriscli3-bugfix advanced ibc set --name=silei --from=$VADDR --chain-id=upgrade-test --sequence=0 --print-response true --fee=20000000000000000iris --home=/data/iriscli
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