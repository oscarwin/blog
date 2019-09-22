# libevent 源码解析-事件循环

最近阅读了 libevent 的源码，写一篇文章来总结自己学习到知识。使用libevent应该优先选用最新的稳定版本，而阅读源码为了降低难度，我选择了1.4的版本，也就是patches-1.4分支。读这篇文章需要 Unix 网络编程的基础，知道 reactor 模式，如果对此还有疑问可以看我这篇文章[典型服务器模式原理分析与实践]()

## libevent 的文件结构

关于 libevent 的文件结构这篇文章[The Libevent Reference Manual: Preliminaries](http://www.wangafu.net/~nickm/libevent-book/Ref0_meta.html)说明的比较清楚了，这里简要说明一下。

**event 和 event_base**

event 和 event_base 是 libevent 的核心，也是我们要探讨的核心，主要围绕两个结构体类型 event 和 event_base 展开，event 定义了事件的结构，event_base 则是事件循环的框架，这两个结构体分别定义在 event.h 和 event-internal.h 文件中。在 event.c 中定义了事件初始化，事件注册，事件删除等 API，还包含了事件循环框架 event base 相关的 API。

**evbuffer 和 bufferevent**

evbuffer 和 bufferevent 则处理了 libevent 中关于读写缓冲的问题，这两个结构体也定义在 event.h 头文件中，而相关 API 则分别在 buffer.c 文件和 evbuffer.c 文件中定义相关API。bufferevent 是一个缓冲区管理结构体，在其中包含了两个 evbuffer 指针，一个是读缓存区，一个是写缓存区。evbuffer 则是与底层 IO 打交道的。另外不得不提的是 bufferevent 中为读缓存区和写缓存区都设定了一个高水位和低水位，高水位是为了避免单个缓存区占用过多的内存，低水位是为了减少回调函数调用的次数，提高效率。

**evutil**

这个模块主要就是不同平台网络通信的实现，也就是 IO 多路复用在不同平台下的实现，已经套接字编程在不同平台下的实现。另外就是使用的一些公共方法了。

## 事件循环

libevent 将事件进一步抽象化了，除了读和写事件，还包括定时事件，甚至将信号也转化成了事件来处理。

首先看一下 event 的结构体。libevent 用链表来保存注册事件和激活的事件，分别存在

```
struct event {
    /*
    ** libevent 用双向链表来保存注册的所有事件，包括IO事件，信号事件。
    ** ev_next 存储了该事件在事件链表中的位置
    ** 另外，libevent 还用另一个链表来存储激活的事件，通过遍历激活的事件链表来分发任务
    ** ev_active_next 存储了该事件在激活事件链表中的位置
    ** 类似，ev_signal_next 就是该事件在信号事件链表中的位置
    */
	TAILQ_ENTRY (event) ev_next;
	TAILQ_ENTRY (event) ev_active_next;
	TAILQ_ENTRY (event) ev_signal_next;
    /* libevent 用最小堆来管理超时时间，min_heap_idx 保存堆顶的 index */
	unsigned int min_heap_idx;	/* for managing timeouts */

    /* event_base 是整个事件循环的核心，每个 event 都处在一个 event_base 中，ev_base 保存这个结构体的指针 */
	struct event_base *ev_base;
    /* 对于 IO 事件，ev_fd 是绑定的文件描述符，对于 signal 事件，ev_fd 是绑定的信号 */
	int ev_fd;
    /* 要处理的事件类型， */
	short ev_events;
    /* 事件就绪执行时，调用ev_callback的次数，通常为1 */
	short ev_ncalls;
	short *ev_pncalls;	/* Allows deletes in callback */
    /* 事件超时的时间长度 */
	struct timeval ev_timeout;
    /* 优先级 */
	int ev_pri;		/* smaller numbers are higher priority */
    /* 响应事件时调用的callback函数 */
	void (*ev_callback)(int, short, void *arg);
	void *ev_arg;

	int ev_res;		/* result passed to event callback */
    /* 表示事件所处的状态 */
	int ev_flags;
};
```