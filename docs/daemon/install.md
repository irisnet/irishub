---
order: 2
---

# Install

## Latest Version

The Latest version of IRIShub is [v0.15.2](https://github.com/irisnet/irishub/releases/latest)

::: tip
Please replace <latest_iris_version> below with v0.15.2
:::

## Install `go`

::: tip
**Go 1.12.5+** is required for the IRIShub.
:::

Install `go` by following the [official docs](https://golang.org/doc/install).

Remember to set your `$GOPATH`, `$GOBIN`, and `$PATH` environment variables, for example:

```bash
mkdir -p $HOME/go/bin
echo "export GOPATH=$HOME/go" >> ~/.bashrc
echo "export GOBIN=$GOPATH/bin" >> ~/.bashrc
echo "export PATH=$PATH:$GOBIN" >> ~/.bashrc
source ~/.bashrc
```

Verify that `go` has been installed successfully

```bash
go version
```

## Install `iris`

After setting up `go` correctly, you should be able to compile and run `iris`.

Make sure that your server can access to google.com because our project depends on some libraries provided by google. (If you are not able to access google.com, you can also try to add a proxy: `export GOPROXY=https://goproxy.io`)

```bash
git clone --branch <latest_iris_version> https://github.com/irisnet/irishub
cd irishub
# source scripts/setTestEnv.sh # to build or install the testnet version
make get_tools install
```

If your environment variables have set up correctly, you should not get any errors by running the above commands.
Now check your `iris` version.

```bash
iris version
<latest_iris_version>

iriscli version
<latest_iris_version>
```
