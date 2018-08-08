#软件升级测试

## 单节点成功升级CASE
### 运行旧软件
```
rm -rf basecoin1
rm -rf .basecli
basecoind init gen-tx --name=x --home=basecoin1
basecoind init --gen-txs --chain-id=upgrade-test -o --home=basecoin1
basecoind start --home=basecoin1
```

```
basecli keys list
VADDR=验证人的地址
```

### 发送升级提议（直接通过，2个区块高度后）
```
basecli gov submit-proposal --name=x --proposer=$VADDR --title=ADD --description=“I am crazy” --type=Text --deposit=10iris --chain-id=upgrade-test 
```

### 查询提议内容
```
basecli gov query-proposal --proposalID=1        
```

### 查询升级的版本信息
```
basecli upgrade info
```

#### 运行新软件
```
basecoind1 start --home=basecoin1
```

### 发送消息自己已运行新软件
```
basecli  upgrade submit-switch --name=x --from=$VADDR --proposalID=1 --chain
-id=upgrade-test
```

### 查询switch信息
```
basecli upgrade query-switch --voter=$VADDR --proposalID=1
```

### 使用新功能（无报错）
```
basecli1 advanced ibc set --name=x --from=$VADDR --chain-id=upgrade-test --sequence=0 --print-response true
basecli1 advanced ibc get --name=x --from=$VADDR --chain-id=upgrade-test --sequence=0 --print-response true
```
## 多节点成功升级测试
```
rm -rf basecoin1
rm -rf .basecli
rm -rf basecoin2
rm -rf .basecli
rm -rf basecoin3
rm -rf .basecli
rm -rf basecoin4
rm -rf .basecli

basecoind init gen-tx --name=x1 --home=basecoin1
basecoind init gen-tx --name=x2 --home=basecoin2
basecoind init gen-tx --name=x3 --home=basecoin3
basecoind init gen-tx --name=x4 --home=basecoin4

cp basecoin2/config/gentx/gentx-XXXX.json basecoin1/config/gentx/
cp basecoin3/config/gentx/gentx-XXXX.json basecoin1/config/gentx/
cp basecoin4/config/gentx/gentx-XXXX.json basecoin1/config/gentx/

basecoind init --gen-txs --chain-id=upgrade-test -o --home=basecoin1

cp ~/basecoin1/config/genesis.json ~/basecoin2/config/
cp ~/basecoin1/config/genesis.json ~/basecoin3/config/
cp ~/basecoin1/config/genesis.json ~/basecoin4/config/

vi basecoin2/config/config.toml
vi basecoin3/config/config.toml
vi basecoin4/config/config.toml
6628995f6eae0c7d810867e467f23530c55b1232@localhost:26656

basecoind start --home=basecoin1
basecoind start --home=basecoin2
basecoind start --home=basecoin3
basecoind start --home=basecoin4
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
basecli gov submit-proposal --name=x1 --proposer=$VADDR1 --title=ADD --description=“I am crazy” --type=Text --deposit=10iris --chain-id=upgrade-test 
```

### 查询提议内容
```
basecli gov query-proposal --proposalID=1        
```

### 查询升级的版本信息
```
basecli upgrade info
```

#### 运行新软件
```
basecoind1 start --home=basecoin1
basecoind1 start --home=basecoin2
basecoind1 start --home=basecoin3
basecoind1 start --home=basecoin4
```

### 发送消息自己已运行新软件
```
basecli  upgrade submit-switch --name=x1 --from=$VADDR1 --proposalID=1 --chain-id=upgrade-test
basecli  upgrade submit-switch --name=x2 --from=$VADDR2 --proposalID=1 --chain-id=upgrade-test
basecli  upgrade submit-switch --name=x3 --from=$VADDR3 --proposalID=1 --chain-id=upgrade-test
basecli  upgrade submit-switch --name=x4 --from=$VADDR4 --proposalID=1 --chain-id=upgrade-test
```

### 查询switch信息
```
basecli upgrade query-switch --voter=$VADDR1 --proposalID=1
basecli upgrade query-switch --voter=$VADDR2 --proposalID=1
basecli upgrade query-switch --voter=$VADDR3 --proposalID=1
basecli upgrade query-switch --voter=$VADDR4 --proposalID=1
```

### 测试新功能
```
basecli1 advanced ibc set --name=x --from=$VADDR --chain-id=upgrade-test --sequence=0 --print-response true
basecli1 advanced ibc get --name=x --from=$VADDR --chain-id=upgrade-test --sequence=0 --print-response true
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
##
```
basecoind unsafe_reset_all --home=basecoin1
basecoind unsafe_reset_all --home=basecoin2
basecoind unsafe_reset_all --home=basecoin3
basecoind unsafe_reset_all --home=basecoin4
```