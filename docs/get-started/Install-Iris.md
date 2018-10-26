# Install Iris

### Step 1: Configure Your Server

All the blockchains in IRISnet is based on Cosmos-SDK, which is a framework for building blockchain applications in Golang. It is being used to build [Cosmos Hub](https://cosmos.network/). It's recommended to run a validator node on Linux server.

**Recommanded Configurations:**

1. 2 CPU
2. Memory: 4GB
3. Disk: 60GB SSD
4. OS: Ubuntu 16.04 LTS
5. Allow all incoming connections on TCP port 26656 and 26657

### Step 2: Install Iris

There are two ways to get Iris running on your server. You can download the binary files from our release pages, or you can download the source code and compile it locally.

#### Download Binary Directly

Go to the download page: https://github.com/irisnet/irishub/releases/  then get the release v0.5.0 on your computer.
`unzip -C /usr/local/bin  iris$VERSION.$OS-$ARCH.zip` to `/usr/local/bin/ ` 

You can verify you have the right version installed by running the following commands:

```
$ iris version
v0.6.0

$ iriscli version
v0.6.0
```

#### Compile Source Code

- Install Go 1.10+

```
$ sudo add-apt-repository ppa:gophers/archive
$ sudo apt-get update
$ sudo apt-get install golang-1.10-go
```

> Note that golang-1.10-go puts binaries in /usr/lib/go-1.10/bin. If you want them on your PATH, you need to make that change yourself.

Using snaps also works quite well:

```
This will give you the latest version of go
$ sudo snap install --classic go
```

> A restart is required for the command to be recognized.

Then you need to verify the versions of Go:

```
$ go version
go version go1.10.3 darwin/amd64
```

Then, you need to add `GOPATH` to system `PATH` , then your system could correctly compile the code.

Open your `.profile` in your home directory. Add the following lines at the end of file:

```
GOPATH=$HOME/go
PATH=$GOPATH/bin:$PATH
```

Save the file and exit the editor. Then run the following to make your bash reload your profile configurations.

```
$ source $HOME/.profile
```

Now you should see something like this if you echo your\$GOPATH and \$PATH variables

```
$ echo $GOPATH
/home/iris/go
$ echo $PATH
/home/isir/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
```

- Get the code and compile Iris

After setup Go correctly, you should be able to compile and run **Iris**.
Make sure that you can access to google.com for that our project used some libraries provided by google.
```
mkdir -p $GOPATH/src/github.com/irisnet
cd $GOPATH/src/github.com/irisnet
git clone https://github.com/irisnet/irishub
cd irishub && git checkout v0.5.0
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
make get_vendor_deps && make install
```

If your environment variables have set up correctly, you should not get any errors by running the above commands.
Now check your **Iris** version.

```
$ iris version
v0.6.0
$ iriscli version
v0.6.0
```

### How to Update

Get latest code (you can also `git fetch` only the version desired), ensure the dependencies are up to date, then recompile.

```
cd $GOPATH/src/github.com/irisnet/irishub
git fetch -a origin
rm Gopkg.lock
git checkout v0.6.0
make get_vendor_deps
make install
```

## Use BaaS 

* https://irisnet.wancloud.cloud/#/
* https://baas.zhigui.com/IRISnet