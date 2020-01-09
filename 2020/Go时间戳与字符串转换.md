Go 语言里时间戳与字符串的转换可谓别具一格，直接给出使用案例。

## 举个栗子

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // 获取当地时间当前的时间戳
    timestamp := time.Now().Unix()
    fmt.Printf("timestamp: %v\n", timestamp)

    // 时间戳转字符串，这里都是转成当前时间
    tm := time.Unix(timestamp, 0)
    fmt.Println(tm.Format("2006-01-02 03:04:05"))
    fmt.Println(tm.Format("02/01/2006 15:04:05"))
    fmt.Println(tm.Format("2006-01-02 15:04:05 Mon Jan"))

    // 字符串转时间戳，第一个参数是格式，第二个是要转换的时间字符串
    // 使用这个函数时，字符里的时间是 UTC 时间
    tm2, _ := time.Parse("01/02/2006 15:04:05", "06/01/2020 00:00:00")
    // 本地时区的时间转换用下面这个方法
    tm3, _ := time.ParseInLocation("01/02/2006 15:04:05", "01/06/2020 00:00:00", time.Local)
    fmt.Println(tm2.Unix(), tm3.Unix())

    now := time.Now()
    fmt.Println(now.Unix())
    year, mon, day := now.UTC().Date()
    hour, min, sec := now.UTC().Clock()
    zone, _ := now.UTC().Zone()
    fmt.Printf("UTC 时间是 %d-%d-%d %02d:%02d:%02d %s\n",
        year, mon, day, hour, min, sec, zone) // UTC 时间是 2016-7-14 07:06:46 UTC
    fmt.Println(now.UTC().Format("2006-01-02 03:04:05"))
}
```

输出结果：
```
timestamp: 1578401853
2020-01-07 08:57:33
07/01/2020 20:57:33
2020-01-07 20:57:33 Tue Jan
1590969600 1578240000
1578401853
UTC 时间是 2020-1-7 12:57:33 UTC
2020-01-07 12:57:33
```

## 获取当前的 Unix 时间戳

go 与时间相关的方法都在 time 包里，因此要先包含 time 包。

获取当前时间的时间戳：
```go
timestamp := time.Now().Unix()
```

先通过 Now() 函数获得当前时间，Now() 返回的是 time 类型，然后通过 Unix 转为 int64 类型。

## 时间戳转字符串

Unix 时间戳转字符串是通过 Format 函数来实现的，在前文案例中，格式化为字符串时用很多数字来表示，不知道你什么感觉，反正我一次见，有点懵逼。但是**这些数字不是瞎写的，每个都是有特殊意义的**。其意义如下：

```
月份 1,01,Jan,January
日　 2,02,_2
时　 3,03,15,PM,pm,AM,am
分　 4,04
秒　 5,05
年　 06,2006
时区 -07,-0700,Z0700,Z07:00,-07:00,MST
周几 Mon,Monday
```
以月份为例进行说明：
- 1 表示月份，用数字展示，不保留前导0
- 01 表示月份，用数字展示，保留前导0
- Jan 表示月份，用英文缩写展示
- January 表示月份，用英文单词展示

这每一个数字或者字符串都是一个占位符，字符串的格式没有特殊规定，需要什么样的格式就添加什么样的字符就好了。比如，你可以输出这样奇葩的格式 `"2006-2006-2006 01-01-02 15:04:05 Mon Jan"`。

```go
    tm := time.Unix(timestamp, 0)
    fmt.Println(tm.Format("2006-01-02 03:04:05"))
    fmt.Println(tm.Format("02/01/2006 15:04:05"))
    fmt.Println(tm.Format("2006-01-02 15:04:05 Mon Jan"))
```

## 字符串转时间戳

字符串转时间戳是通过 Parse 函数来实现的，不过这里有个坑，那就 UTC 时间与时区的问题。UTC 时间是世界标准时间，所谓世界标准时间就是地理课中的本初子午线那个地方的时间。世界其他各国的时间都以那个时间为标准，世界其他地方的时间被分成了 24 个时区。

北京时间比 UTC 时间快 8 个小时，换句话说：北京时间 = UTC 时间 + 8 小时，因此北京时间也可以写作 UTC+8
夏威夷时间比 UTC 时间慢 10 个小时，换句话说：夏威夷时间 = UTC 时间 - 10 小时，因此夏威夷时间也可以写作 UTC - 10

Parse 函数解析字符串时间时，认为这个字符串的时间是 UTC 的时间，如果你想把北京时间 2020年1月1日 08:00:00 转成 Unix 时间戳要用 ParseInLocation 这个函数。

```
    // 字符串转时间戳，第一个参数是格式，第二个是要转换的时间字符串
    // 使用这个函数时，字符里的时间是 UTC 时间
    tm2, _ := time.Parse("2020-01-02 15:04:05", "2020-01-01 08:00:00")
    // 本地时区的时间转换用下面这个方法
    tm3, _ := time.ParseInLocation("2020-01-02 15:04:05", "2020-01-01 08:00:00", time.Local)
```

## 参考

[Go时间戳和日期字符串的相互转换](https://www.cnblogs.com/baiyuxiong/p/4349595.html)
[golang 格式化时间总结](https://blog.csdn.net/x356982611/article/details/87972400)