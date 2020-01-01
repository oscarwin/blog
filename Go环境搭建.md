
## GO 环境搭建

Go 安装整体来说都是非常简单的，下面分别说明在不同操作系统上的安装步骤

### Linux 上安装

**下载安装包**





访问[下载网址](https://golang.google.cn/dl/)下载 linux 系统的最新版本安装包，或者复制下载链接（安装包的下载链接上右键-复制链接地址）后使用 wget 命令进行下载，我写这篇文章时最新稳定版本为 1.13.4。

```
wget https://dl.google.com/go/go1.13.4.linux-amd64.tar.gz
```

**解压到安装目录**

假设我们想把 Go 安装在 /usr/local 目录下，那么就将下载下来的安装解压到这个目录，在下载的安装包目录下执行下面的命令，在 /usr/local 目录下就会出现一个 go 文件夹。

```
tar -zxvf go1.13.4.linux-amd64.tar.gz -C /usr/local
```

如果你想安装在其他目录，那么这一步解压操作和后面的操作都相应替换为你的安装目录。

**设置环境变量**

执行以下命令将 Go 的 bin 目录添加到环境变量：

```
export PATH=$PATH:/usr/local/go/bin
```

**验证安装是否成功**

命令行运行 go，得到以下输出说明安装成功。

```
$ go
Go is a tool for managing Go source code.

Usage:

        go <command> [arguments]

The commands are:

        bug         start a bug report
        build       compile packages and dependencies
        clean       remove object files and cached files
        doc         show documentation for package or symbol
        env         print Go environment information
        fix         update packages to use new APIs
        fmt         gofmt (reformat) package sources
        generate    generate Go files by processing source
        get         add dependencies to current module and install them
        install     compile and install packages and dependencies
        list        list packages or modules
        mod         module maintenance
        run         compile and run Go program
        test        test packages
        tool        run specified go tool
        version     print Go version
        vet         report likely mistakes in packages

Use "go help <command>" for more information about a command.

Additional help topics:

        buildmode   build modes
        c           calling between Go and C
        cache       build and test caching
        environment environment variables
        filetype    file types
        go.mod      the go.mod file
        gopath      GOPATH environment variable
        gopath-get  legacy GOPATH go get
        goproxy     module proxy protocol
        importpath  import path syntax
        modules     modules, module versions, and more
        module-get  module-aware go get
        module-auth module authentication using go.sum
        module-private module configuration for non-public modules
        packages    package lists and patterns
        testflag    testing flags
        testfunc    testing functions

Use "go help <topic>" for more information about that topic.
```

### Mac 上安装

## 工作空间

### Go 的环境变量

Go 环境搭建主要会涉及到3个环境变量：GOROOT、GOPATH、GOBIN。

- GOROOT：Go 的安装目录。我的安装目录是 /usr/local，所以该变量的值设为 /usr/local/go。
- GOPATH：Go 的工作目录。所谓工作目录就是存放你自己 Go 项目的地方，我打算将工作目录放在 /home/mygo，因此将该变量的值设定为 /home/mygo。
- GOBIN：Go 可执行文件的目录。go install 命令会将 go 项目编译成可执行文件，如果设置了这个环境变量，那么可执行文件就会发到这个目录下。该环境变量不是必须的，如果你还不是很清楚，那么就先不要设定好了。

Go 的工作目录的格式是固定的，在 Go 的工作目录里分为三个文件夹 src、pkg、bin。你需要自己创建这些文件夹，每个文件夹的作用如下：

- src：存放项目的源代码
- pkg：存放编译后的文件，编译后的包存放在这个目录供其他包进行使用
- bin：编译后生成的可执行文件。如果没有设定环境变量 GOBIN，那么 go install 命令的得到的可执行文件就会放入这个文件夹，如果设定了 GOBIN，那么编译得到的可执行文件就放入 GOBIN 的目录

之前安装 Go 时直接在 shell 中通过 export 命令设定环境变量，这种方式设定的环境变量是临时的，机器重启后需要重新设定，现在我们用修改文件的方式来设定环境变量。

执行下面的命令，用 vim 打开 /etc/profile，在该文件的最后添加下面几行，保存后执行 source /etc/profile。

```
export GOROOT=/usr/local/go
export GOPATH=/home/mygo
#export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```

其中我注释了 GOBIN 的环境变量，你可以启用也可以注释。另外在 PATH 里添加了 $GOROOT/bin 和 $GOPATH/bin，这样这两个目录里的可执行文件就可以在任意目录下使用命令行运行了。

### 一个最简单的 go 项目

前面提到在 Go 的工作目录中 src 文件用来存放源代码，一般没新建一个项目就在 src 目录下建一个文件夹。现在我们在 src 目录下建立一个 myhello 的文件，然后进入该文件夹新建一个名为 hello.go 的文件。

```
cd /home/mygo/src
mkdir myhello
touch hello.go
```

向 hello.go 中输入以下代码，并保存文件：

```
package myhello

import (
        "fmt"
)

func Hello() {
        fmt.Println("hello!")
}
```

然后在 myhello 目录执行里 go install 命令，这时候就会在 pkg 目录下生成 linux_amd64 文件夹，该文件夹中有一个 myhello.a 文件。linux_amd64 的名称是与平台相关的。因为 myhello 项目中没有 main 函数，所以 go install 命令不会将其编译成可执行文件，只是将其编译成供其他包使用的包文件，并存放在 pkg 目录。这时候工作目录结构如下：

```
/home/mygo$ tree
.
├── bin
├── pkg
│   └── linux_amd64
│       └── myhello.a
└── src
    └── myhello
        └── hello.go
```

那么包已经编译好了，如何提供给其他项目使用呢？这时候在 src 目录再建立一个 myapp 文件夹，在该文件夹下建一个 main.go 文件，并输入以下代码。

```
package main

import myhello

func main() {
    myhello.Hello()                                                                                                                                           
}
```

在 myapp 目录下执行 go install 命令，就编译并安装了该项目，会在 bin 目录里生成一个可执行文件 myapp。这时候的工作目录结构如下：

```
/home/mygo$ tree
.
├── bin
│   └── myapp
├── pkg
│   └── linux_amd64
│       └── myhello.a
└── src
    ├── myapp
    │   └── main.go
    └── myhello
        └── hello.go
```

安装该项目后，就可以在任意目录里执行 myapp 程序，因为之前已经将 $GOPATH/bin 目录加入 PATH 环境变量。

```
$ myapp
hello!
```

### 获取远程包

除了使用自己编写的包以外，项目中还会经常依赖到开源的一些工具包。通过 go get 命令就可以获取远程包，而且 go get 命令不仅仅下载了该包，而且还会自动进行安装，目前 go get 支持多个开源社区（例如：github、googlecode、bitbucket、Launchpad）。

以 github 上 kafka 的 go 版本客户端为例，该项目的地址为：[https://github.com/segmentio/kafka-go](https://github.com/segmentio/kafka-go)。那么运行如下命令就能将其下载并安装到本机（安装过程有点慢，稍微等一下）。

```
go get -u github.com/segmentio/kafka-go
```

> go get -u 参数可以自动更新包，而且 go get 的时候会自动获取该包依赖的其他第三方包

该命令执行完成后，在 src 目录里会多出一个 github.com 目录，所有从 github 上获取的远程包的源码都会放在这个目录下，在该目录里就有刚刚从 github 上下载下来的 segmentio/kafka-go。除此之外，在 pkg/linux_amd64 目录下也会多出一个 github.com 目录，这个目录中就存放了远程包的编译文件。

这个时候的包结构如下所示：

```
.
├── bin
│   └── myapp
├── pkg
│   └── linux_amd64
│       ├── github.com
│       │   └── segmentio
│       │       └── kafka-go.a
│       └── myhello.a
└── src
    ├── github.com
    │   └── segmentio
    │       └── kafka-go
    │           ├── balancer.go
    │           ├── balancer_test.go
    │           ├── batch.go
    │           ├── batch_test.go
    │           ├── buffer.go
    │           ├── client.go
    |           ......
    ├── myapp
    │   └── main.go
    └── myhello
        └── hello.go
```
