每当提起 go 的优点时，不得不提的一点是 go 支持原生协程。对于传统的 C 或者 Java 语言，需要通过线程来实现并发，然后通过互斥锁来实现线程同步。而 go 则通过协程(goroutines)来实现并发，通过通道(channels)来实现同步。go 语言这种并发模型被称为 CSP(Com-municating Sequential Processes)。

go 并发模型里两个最重要的概念是 goroutines 和 channels：
- goroutines：一个 goroutine 就是独立运行的函数，就像是一个单独运行的线程一样，但是它并不是线程，被称为协程，协程往往被称为轻量级线程。
- channels：通道就像一个管道一样，协程之前可以通过通道来传递数据，实现同步和数据通信。

如果你是首次接触 go 语言，那么 go 的并发编程一定会让你惊讶不已，下面实现一个简单的并发程序。

```
package main

import "fmt"
import "time"

func main() {
	fmt.Println("main start")
	go echo()                          // 调用 echo 函数开启协程
	fmt.Println("main end")
	time.Sleep(500 * time.Millisecond) // 休眠500ms
}

func echo() {
	fmt.Println("echo a line")
}
```

输出：
```
$: go run simple_goroutines.go 
main start
main end
echo a line
```

这个程序是怎么运行的呢？这个程序中，首先 main 函数打印出 “main start”，然后通过一个 go 关键词就开启了一个协程，在这个协程了执行了 echo 函数，开启协程后主程序还继续往下运行打印出 “main end”，之后主程序休眠 500ms 等待协程运行结束。

## runtime.Gosched() 解析

```
package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println("outside a goroutine")
	go func() {
		fmt.Println("inside a goroutine")
	}()
	fmt.Println("outside a goroutine again")
	runtime.Gosched()
}
```

在进程中有时候想打开一个文件，或者输出一些数据，但是不想阻塞当前的函数，这时候就可以开启一个协程来执行。除了使用函数来开启协程外，还可以通过一个匿名函数来开启协程，如本例中就使用了一个匿名函数来开启协程chu