# Docker 镜像与容器

docker 是基于 Linux 内核的 cgroup 和 namespace 等技术实现的进程隔离技术，是操作系统层面的虚拟化技术。由于隔离的进程独立于宿主和其它的隔离的进程，因此也称其为容器。

docker 容器与传统的虚拟化技术最大的区别在于，传统的虚拟化技术是在宿主机器上虚拟出一套硬件后，在其上面运行完整的操作系统，然后在该系统上运行用户进程，属于硬件层面的虚拟化技术。而 docker 是操作系统层面的虚拟化技术，不需要虚拟化硬件，同一宿主机器上的所有容器共享宿主机的操作系统内核，只是在用户态隔离进程的运行环境和存储，因此对于操作的用户而言就好像是一台独立的主机。下图是 docker 官网上的一张对比图，左边描述的是 docker 容器技术，右边描述的是虚拟机技术。

![docker与虚拟机的区别](./images/diff_container_virtual_machine.png)

这个区别使得 docker 相比与传统的虚拟化技术要轻量的多，启动一个虚拟机需要占用数 GB 的磁盘空间，需要几分钟的时间，而启动一个容器只一般只要几十 M，在秒级的时间内完成启动。

docker 中有三个重要的概念：镜像、容器和仓库。首先容器从表象上看就是一台 linux 主机，在这个主机中会启动需要的进程，比如一个 redis 的容器就会启动一个 redis 的服务器。而启动一个容器需要几个镜像，你可以想象镜像就是容器的安装包。最后仓库就是放镜像的地方，你可以想象仓库是 linux 系统的软件包库。

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
其中常见的命令选项有：
- -i 选项是让容器标准输入打开，就可以接受键盘输入了
- -t 选项是让docker分配一个伪终端，绑定到标准输入上。通过这个伪终端就可以像操作一台 linux 机器来操作这个容器了。
- --name <容器名称> 选项为容器指定一个名称
- -d 选项让容器在后台运行，什么是后台运行，下文有说明

用这个命令启动一个容器后，就可以像操作一台普通的 ubuntu 机器一样，操作这个容器了。可以尝试在这个容器中执行一些 linux 命令。
```
root@e7cd5bc73508:/# ps
  PID TTY          TIME CMD
    1 pts/0    00:00:00 bash
   10 pts/0    00:00:00 ps

root@e7cd5bc73508:/# ls
bin  boot  dev  etc  home  lib  lib64  media  mnt  opt  proc  root  run  sbin  srv  sys  tmp  usr  var
```

需要注意的是，对于这种以交互模式启动的容器，当终止交互后，容器就退出了。

`docker container start` 可以启动一个终止的容器，该命令的使用格式是：
```
docker container start [OPTIONS] CONTAINER [CONTAINER...]
```
要启动一个容器，需要知道容器的ID来作为`docker container start`的参数，接下来介绍如何查看容器。

### 查询容器状态——docker container ls

`docker container ls` 列出所有正在运行的容器

`docker container ls -a` 列出所有的容器，包括已经终止的容器

列出的容器信息包括：

- 容器ID（CONTAINER ID），这是容器的唯一标识，操作容器相关的命令都需要带上这个标识
- 依赖的镜像名称（IMAGE）
- 执行的命令（COMMAND）
- 创建的时间（CREATED）
- 容器的状态（STATUS），UP 是运行状态，Exited 是终止状态
- 暴露的端口（PORTS）
- 容器的名称（NAMES），可以在启动容器的时候通过 --name 选项指定容器的名称，如果没有指定，系统会生成一个默认名称

### 终止容器——docker container stop

`docker container stop` 命令用来终止一个正在运行的容器

终止的容器可以通过 `docker container start` 命令启动，`docker container restart` 命令会先终止容器，然后再启动容器。

### 守护态运行

大多数时候需要容器在后台运行，不需要将输出结果打印到宿主主机上。此时，可以通过 -d 参数来实现。

不使用 -d 选项启动下面这个容器，表现如下：
```
$ docker run ubuntu /bin/bash -c "while true; do echo hello world; sleep 1; done"
hello world
hello world
hello world
hello world
hello world
```

使用 -d 选项启动这个容器，表现如下：

```
$ docker run -d ubuntu /bin/bash -c "while true; do echo hello world; sleep 1; done"
574ee145e9ca501639a233601efe9574a356a8d7bdef461d240b2212b6aeaf77
```

在宿主的终端只输出了该容器的 ID。如果需要查看容器的输出可以使用 `docker container logs` 命令。本文在说容器的这些命令时都省略了最后一个参数容器 ID，这个读者在使用的时候需要自行加上，查询容器 ID 的方法前面已经介绍过了。如果想要持续获取容器的输出可以使用 `docker container logs -f` 命令。

### 进入容器

使用 -d 参数启动容器后，容器在后台运行，如果需要进入容器可以使用 `docker attach` 命令或 `docker exec` 命令。建议使用 `docker exec` 命令，原因在下面进行说明。

**docker attach**

```
$ docker run -it -d ubuntu
2b10efdeaa9f0f24f0060bf636cc8f1bae9598f5e0176b498447a7b63ea10d06

$ docker container ls
CONTAINER ID    IMAGE     COMMAND       CREATED             STATUS              PORTS          NAMES
2b10efdeaa9f    ubuntu    "/bin/bash"   13 seconds ago      Up 11 seconds                      cocky_golick

$ docker attach 2b10efdeaa9f
root@2b10efdeaa9f:/# 
```

使用 attach 命令，如果从这个终端退出后，这个容器也会被终止，这就是不推荐使用 `attach` 命令的原因。

**docker exec**

```
$ docker run -it -d ubuntu
61299a9b910b029f9a8667131a45aff92fbc35633e4f86c8fcfb72e6362c5115

$ docker container ls
CONTAINER ID        IMAGE      COMMAND         CREATED             STATUS          PORTS           NAMES
61299a9b910b        ubuntu     "/bin/bash"     6 seconds ago       Up 5 seconds                    friendly_hawking

$ docker exec -it 61299a9b910b /bin/bash
```

使用 exec 命令，可以生成新的伪终端与容器进行交互，因此退出时不会导致容器退出。

### 删除容器

`docker container rm` 命令可以删除一个容器，如果要删除一个正在运行的容器可以使用`docker container rm -f`。

如果要清除所有终止的容器可以使用 `docker container prune` 命令。

## 访问仓库

- `docker login` 登录到docker hub仓库
- `docker logout` 退出登录
- `docker search <关键词>` 搜索镜像
- `docker pull` 从仓库下载镜像
- `docker push` 将本地仓库推送到远程仓库

下面演示一下这些命令的使用：
```
$ docker login
Login with your Docker ID to push and pull images from Docker Hub. If you don't have a Docker ID, head over to https://hub.docker.com to create one.
Username: xxx
Password: xxx
WARNING! Your password will be stored unencrypted in /root/.docker/config.json.
Configure a credential helper to remove this warning. See
https://docs.docker.com/engine/reference/commandline/login/#credentials-store

Login Succeeded

$ docker pull hello-world

# 为镜像 hello-world 生成一个新的标签，而且指定仓库为自己仓库的用户名
$ docker tag hello-world:latest oscarwin/hello-world:3.0

# 将镜像推送到远程仓库，那么就可以在 docker hub 里看到自己推送的这个镜像了
$ docker push oscarwin/hello-world
``` 