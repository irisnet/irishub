# 软件升级开发者文档

> 按照一定的规范来写新版本的软件，来通过upgrade module 来升级软件

## 增加新module或通过引入新module修改旧module逻辑

> 该部分适用于本次升级为新加module的场景。具体到代码上的限制是：
> 1. 只在 github.com/irisnet/irishub/modules 下增加一个或多个module来实现新功能或覆盖老module
> 2. 只在 github.com/irisnet/irishub/app/app.go 中增加配置新module的``keeper，router，registerWire``等代码，不能修改 ``initChainer`` 中的代码

### 如何编写新版本软件

> 使用命令 ``iriscli upgrade info`` 查询当前版本信息（假设当前内部版本号为 ``0``）

> 在 github.com/irisnet/irishub/modules 下创建新module的folder，并在其中开发实现该module的业务逻辑（假设module名字为 ``myModule``), 注意该module中返回的msgType必须为 ``当前版本号 + 1`` (在该例子中应配置为1)

```
func (msg MyModuleMsg) Type() string { return "myModule-1" }

```

> 在 github.com/irisnet/irishub/app/app.go 中注册新module，注意一定要在注册router的时候把新module的起始生效版本号配置为 ``当前版本号 + 1`` (在该例子中应配置为1):

```
app.myModuleKeeper = myModule.NewKeeper(app.accountMapper)
app.Router().AddRoute("myModule-1", []*sdk.KVStoreKey{app.keyAccount}, myModule.NewHandler(app.myModuleKeeper))

```

> 在 github.com/irisnet/irishub/app/app.go 中配置 ``func MakeCodec() *wire.Codec``:

```
myModule.RegisterWire(cdc)
```

## BugFix升级

> 该部分适用于本次升级不新加module，只修改module内部代码的场景

### 如何编写新版本软件

> 使用命令 ``iriscli upgrade info`` 查询当前版本信息（假设当前内部版本号为 ``0``）

> 修改module内部代码，指定不同代码段的生效版本号， 本次新加代码对应版本为 ``当前版本号 + 1`` (在该例子中应配置为1):

```
if uk.OnlyRunAfterVersionId(ctx, 1) {
    myKeeper.Set(ctx,msg.Addr.String()+":myModule-1")
    return sdk.Result{Log:"This is new module - myModule1 !!"}
} else {
    myKeeper.Set(ctx,msg.Addr.String()+":myModule-0")
    return sdk.Result{Log:"This is new module - myModule0 !!"}
}
```
