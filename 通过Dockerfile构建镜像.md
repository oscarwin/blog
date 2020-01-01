高效mac使用



HomeBrew 是 mac 的包管理系统，类似于 centos 的 yum 或者 ubuntu 的 apt-get。因此是首先要安装的。安装的方法也很简单，打开终端就运行下面的命令就可以安装了：

```
/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
```

然后就可以用 brew 来安装大部分第三方的应用了：

```
brew install wget // 安装 wget 工具
```

Mac 自带的终端虽然也能满足一定的需求，但是有更好用，更强大的终端 iTerm2，安装的方法是在终端执行以下命令：
```
brew cask install iterm2
```

iTerm2 中默认没有快捷键来实现光标向左向右移动一个单词，因此需要手动添加来实现，添加快捷键 option + → 来表示向右移动一个单词。iTerm2 -> preferences 进入设置窗口。

![添加option ]()

iTerm2 相关的快捷键

1. Cmd + T：新建一个标签页
2. Cmd + D：垂直切分窗口，切分为左右两个
2. Shift + Cmd + D：水平切分窗口，切分成上下两个窗口

1. 更强大的shell——zsh

1. Markdown编辑利器——Typora

快捷键