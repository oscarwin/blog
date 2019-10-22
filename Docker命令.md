
## 使用镜像

### 下载镜像——docker pull

下载镜像使用 docker pull 命令来实现，该命令格式如下：
```
docker pull [OPTIONS] NAME[:TAG|@DIGEST]
```

- 可以通过 docker pull --help 查询帮助信息
- OPTIONS 中 -a 选项表示拉取所有标签的镜像，-q 选项表示静默下载，不打印输出信息到屏幕
- NAME 由`仓库名/镜像名:标签名`的格式组成，仓库名没有则从默认仓库 `library` 拉取，也就是官方镜像，镜像名必须有，标签名没有则默认拉取 `latest` 标签的镜像

比如：
```
$ docker pull ubuntu:18.04
18.04: Pulling from library/ubuntu
bf5d46315322: Pull complete
9f13e0ac480c: Pull complete
e8988b5b3097: Pull complete
40af181810e7: Pull complete
e6f7c7e5c03e: Pull complete
Digest: sha256:147913621d9cdea08853f6ba9116c2e27a3ceffecf3b492983ae97c3d643fbbe
Status: Downloaded newer image for ubuntu:18.04
```

### 列出镜像——docker image ls

`docker image ls` 命令可以列出系统中已经下载了的所有镜像。docker 中很多命令与 linux 的命令相似，因此在使用的时候可以类比，docker 中有诸如 rm、ls、ps等命令。

```
$ docker image ls         
REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
zookeeper           latest              e41846a619f5        17 hours ago        224MB
ubuntu              18.04               cf0f3ca922e0        3 days ago          64.2MB
ubuntu              latest              cf0f3ca922e0        3 days ago          64.2MB
redis               latest              f7302e4ab3a8        2 months ago        98.2MB
hello-world         latest              fce289e99eb9        9 months ago        1.84kB
```
列表分别包含了`仓库名`、`标签`、`镜像ID`、`创建时间`、`镜像大小`。`docker image ls`命令展示的是一个完整镜像的大小，但是由于镜像是分层存储的，不同的镜像如果使用了相同的层，这个层只会存储一份数据，所有镜像公用，因此镜像实际占用的空间大小要比显示的小。

### 删除镜像——docker image rm

`docker image rm` 命令可以删除镜像，一般配合`docker image ls`命令来使用。`docker image rm`命令的使用格式是：
```
docker image rm [OPTIONS] IMAGE [IMAGE...]
```

- 使用 `docker image rm --help` 可以查看帮助
- OPTIONS: --force/-f 强制删除
- 镜像可以是镜像名，或者镜像名加标签，或者镜像的ID

比如：

```
$ docker image rm ubuntu:latest
Untagged: ubuntu:18.04
Untagged: ubuntu@sha256:a7b8b7b33e44b123d7f997bd4d3d0a59fafc63e203d17efedf09ff3f6f516152
Deleted: sha256:cf0f3ca922e08045795f67138b394c7287fbc0f4842ee39244a1a1aaca8c5e1c
Deleted: sha256:c808877c0adcf4ff8dcd2917c5c517dcfc76e9e8a035728fd8f0eae195d11908
Deleted: sha256:cdf75cc6b4d28e72a9931be2a88c6c421ad03cbf984b099916a74f107e6708ff
Deleted: sha256:b9997ded97a1c277d55be0d803cf76ee6e7b2e8235d610de0020a7c84c837b93
Deleted: sha256:a090697502b8d19fbc83afb24d8fb59b01e48bf87763a00ca55cfff42423ad36
```

**Untagged 和 Deleted**
在执行 `docker image rm` 时，出现了两种情况，Untagged 和 Deleted。因为一个镜像可以对应多个标签，当执行 `rm` 命令时可能只是移除这个标签，只有当移除最后一个标签的时候才会删除这个镜像。除此之外，如果有别的镜像依赖这个镜像作为基础层，那么这是这个镜像也不会被删除。另外，如果有其他容器依赖于这个镜像，那么也不能删除这个镜像，因为删除这个镜像必然会导致容器出现故障。

## 操作容器

### 运行容器——docker run

运行容器有两种方式，一种是新建一个容器并启动，另一种是启动一个已终止的容器。

`docker run` 命令用来新建一个容器并启动，该命令的使用格式是：
```
docker run [OPTIONS] IMAGE [COMMAND] [ARG...]
``` 

比如下面这个命令会通过 ubuntu 镜像创建并运行一个容器，然后执行 `/bin/bash` 命令。
```
$ docker run -it ubuntu /bin/bash
```
其中：
- -i 选项是让容器标准输入打开，就可以接受键盘输入了
- -t 选项是让docker分配一个伪终端，绑定到标准输入上
- /bin/bash 就是要在容器中执行的命令

用这个命令启动一个容器后，就可以像操作一台普通的 ubuntu 机器一样，操作这个容器了。可以尝试在这个容器中执行一些 linux 命令。
```
root@e7cd5bc73508:/# ps
  PID TTY          TIME CMD
    1 pts/0    00:00:00 bash
   10 pts/0    00:00:00 ps

root@e7cd5bc73508:/# ls
bin  boot  dev  etc  home  lib  lib64  media  mnt  opt  proc  root  run  sbin  srv  sys  tmp  usr  var
```