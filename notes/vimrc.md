# vim 配置

Vim 的全局配置一般在/etc/vim/vimrc或者/etc/vimrc，对所有用户生效。用户个人的配置在~/.vimrc。

如果只对单次编辑启用某个配置项，可以在命令模式下，先输入一个冒号，再输入配置。举例来说，set number 这个配置可以写在 .vimrc 里面，也可以在命令模式输入 :set number。

配置项一般都有"打开"和"关闭"两个设置。"关闭"就是在"打开"前面加上前缀"no"。

    " 打开
    set number

    " 关闭
    set nonumber

总结一个自己常用的 vim 配置

```
set autochdir                   " 自动切换当前目录为当前文件所在的目录
set cursorline                  " 突出显示当前行
set equalalways                 " 分割窗口时保持相等的宽/高
set nocompatible                " 关闭 vi 兼容模式
set number                      " 显示行号
set nobackup                    " 覆盖文件时不备份
set noswapfile                  " 编辑时不产生交换文件
set noexpandtab                 " 插入 tab 符号不以空格替换
set history=1000                " 设置冒号命令和搜索命令的命令历史列表的长度
"set autoindent                 " 开启自动缩进
"set smartindent                " 开启新行时使用智能自动缩进
set smarttab                    " 开启新行时使用智能 tab 缩进
set tabstop=4                   " 设定 tab 长度为 4
set shiftwidth=4                " 设定 << 和 >> 命令移动时的宽度为 4
set showmatch                   " 插入括号时，短暂地跳转到匹配的对应括号
set backspace=indent,eol,start  "不设定在插入状态无法用退格键和 Delete 键删除回车符
set guioptions=t                " 隐藏菜单栏、工具栏、滚动条
set ruler                       " 打开状态栏标尺
set incsearch                   " 输入搜索内容时就显示搜索结果
set hlsearch                    " 搜索时高亮显示被找到的文本
set ignorecase                  " 搜索时忽略大小写
set fileencodings=ucs-bom,utf-8,cp936,gb18030,gb2312,big5,euc-jp,euc-kr,latin1
set background=dark
syntax on                       " 自动语法高亮
set mouse=a                     " 允许使用鼠标

iab xdate <c-r>=strftime("%Y-%m-%d")<CR>
iab xfile <c-r>=expand("%:t")<CR>
```