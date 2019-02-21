  
Linux系统时间有两种。

（1）日历时间。该值是自协调世界时(UTC)1970年1月1日00:00:00这个特定时间以来所经过的秒数累计值。基本数据类型用`time_t`保存。最后通过转换才能得到我们平时所看到的24小时制或者12小时间制的时间。

（2）进程时间。也被称为CPU时间，用以度量进程使用的中央处理器资源。进程时间以时钟滴答计算。

## 获取时间戳
### time()
```
#include <time.h>
time_t time(time_t *calptr)
```
- `time`返回当前时间的时间戳，也就是从世界时到现在的秒数;
- `time_t`实际就是一个`uint64_t`；
- `calptr`不为空时，时间戳也会写入到该指针中；

调用示例：
```
#include <time.h>
#include <iostream>
#include <stdlib.h>

using namespace std;

int main()
{
	time_t curTime;
	curTime = time(NULL);
	cout << curTime << endl;
	return 0;
}
```
结果：
返回一串数值，如`1533287924`

### gettimeofday()和clock_gettime()
`time`函数只能得到秒精度的时间，为了获得更高精度的时间戳，需要其他函数。`gettimeofday`函数可以获得微秒精度的时间戳，用结构体`timeval`来保存；`clock_gettime`函数可以获得纳秒精度的时间戳，用结构体`timespec`来保存。

```
#include <sys/time.h>

int gettimeofday(struct timeval *tp, void *tzp);
因为历史原因tzp的唯一合法值是NULL，因此调用时写入NULL即可。

int clock_gettime(clockid_t clock_id, strcut timespec *tsp);
clock_id有多个选择，当选择为CLOCK_REALTIME时与time的功能相似，但是时间精度更高。
```

两个函数使用的结构体定义如下：
```
struct timeval
{
    long tv_sec; /*秒*/
    long tv_usec; /*微秒*/
};

struct timespec
{
	time_t tv_sec;  //秒
	long tv_nsec;   //纳秒
};
```

调用示例：
```
#include <time.h>
#include <sys/time.h>
#include <iostream>
#include <stdlib.h>

using namespace std;

int main()
{
	time_t dwCurTime1;
	dwCurTime1 = time(NULL);

    struct timeval stCurTime2;
    gettimeofday(&stCurTime2, NULL);

    struct timespec stCurTime3;
    clock_gettime(CLOCK_REALTIME, &stCurTime3);

    cout << "Time1: " << dwCurTime1 << "s" << endl;
    cout << "Time2: " << stCurTime2.tv_sec << "s, " << stCurTime2.tv_usec << "us" << endl;
    cout << "Time3: " << stCurTime3.tv_sec << "s, " << stCurTime3.tv_nsec << "ns" << endl;

	return 0;
}
```

结果：

	编译时要在编译命令最后加上-lrt链接Real Time动态库，如
	g++ -o time2 test_time_linux_2.cpp -lrt
	
	Time1: 1533289490s
	Time2: 1533289490s, 133547us
	Time3: 1533289490s, 133550060ns

### 可视化时间

#### tm结构体
得到的时间戳不能直观的展示现在的时间，为此需要使用`tm`结构体来表示成我们日常所见的时间，该结构体定义如下：
```
struct tm
{
    int tm_sec;  /*秒，正常范围0-59， 但允许至61*/
    int tm_min;  /*分钟，0-59*/
    int tm_hour; /*小时， 0-23*/
    int tm_mday; /*日，即一个月中的第几天，1-31*/
    int tm_mon;  /*月， 从一月算起，0-11*/  1+p->tm_mon;
    int tm_year;  /*年， 从1900至今已经多少年*/  1900＋ p->tm_year;
    int tm_wday; /*星期，一周中的第几天， 从星期日算起，0-6*/
    int tm_yday; /*从今年1月1日到目前的天数，范围0-365*/
    int tm_isdst; /*日光节约时间的旗标*/
};
```

#### time_t转成tm
`gmtime` 和`localtime`可以将`time_t`类型的时间戳转为`tm`结构体，用法如下：

```
struct tm* gmtime(const time_t *timep);
//将time_t表示的时间转换为没有经过时区转换的UTC时间，是一个struct tm结构指针

stuct tm* localtime(const time_t *timep);
//和gmtime功能类似，但是它是经过时区转换的时间，也就是可以转化为北京时间。
```

#### 固定格式打印时间
得到`tm`结构体后，可以将其转为字符串格式的日常使用的时间，或者直接从`time_t`进行转换，分别可以使用以下两个函数达到目的。不过这两个函数只能打印固定格式的时间。
```
//这两个函数已经被标记为弃用，尽量使用后面将要介绍的函数
char *asctime(const struct tm* timeptr);
char *ctime(const time_t *timep);
```
调用示例：
```
#include <time.h>
#include <sys/time.h>
#include <iostream>
#include <stdlib.h>

using namespace std;

int main()
{
	time_t dwCurTime1;
	dwCurTime1 = time(NULL);

    struct tm* pTime;
    pTime = localtime(&dwCurTime1);

    char* strTime1;
    char* strTime2;
    strTime1 = asctime(pTime);
    strTime2 = ctime(&dwCurTime1);

    cout << strTime1 << endl;
    cout << strTime2 << endl;

	return 0;
}
```

结果：
```
Fri Aug  3 18:24:29 2018
Fri Aug  3 18:24:29 2018
```

#### 灵活安全的时间转换函数strftime()
上述两个函数因为可能出现缓冲区溢出的问题而被标记为弃用，因此更加安全的方法是采用`strftime`方法。
```
/*
** @buf：存储输出的时间
** @maxsize：缓存区的最大字节长度
** @format：指定输出时间的格式
** @tmptr：指向结构体tm的指针
*/
size_t strftime(char* buf, size_t maxsize, const char *format, const struct tm *tmptr);
```

我们可以根据format指向字符串中格式，将timeptr中存储的时间信息按照format指定的形式输出到buf中，最多向缓冲区buf中存放maxsize个字符。该函数返回向buf指向的字符串中放置的字符数。

函数strftime()的操作有些类似于sprintf()：识别以百分号(%)开始的格式命令集合，格式化输出结果放在一个字符串中。格式化命令说明串 strDest中各种日期和时间信息的确切表示方法。格式串中的其他字符原样放进串中。格式命令列在下面，它们是区分大小写的。

	%a 星期几的简写
	%A 星期几的全称
	%b 月分的简写
	%B 月份的全称
	%c 标准的日期的时间串
	%C 年份的后两位数字
	%d 十进制表示的每月的第几天
	%D 月/天/年
	%e 在两字符域中，十进制表示的每月的第几天
	%F 年-月-日
	%g 年份的后两位数字，使用基于周的年
	%G 年分，使用基于周的年
	%h 简写的月份名
	%H 24小时制的小时
	%I 12小时制的小时
	%j 十进制表示的每年的第几天
	%m 十进制表示的月份
	%M 十时制表示的分钟数
	%n 新行符
	%p 本地的AM或PM的等价显示
	%r 12小时的时间
	%R 显示小时和分钟：hh:mm
	%S 十进制的秒数
	%t 水平制表符
	%T 显示时分秒：hh:mm:ss
	%u 每周的第几天，星期一为第一天 （值从0到6，星期一为0）
	%U 第年的第几周，把星期日做为第一天（值从0到53）
	%V 每年的第几周，使用基于周的年
	%w 十进制表示的星期几（值从0到6，星期天为0）
	%W 每年的第几周，把星期一做为第一天（值从0到53）
	%x 标准的日期串
	%X 标准的时间串
	%y 不带世纪的十进制年份（值从0到99）
	%Y 带世纪部分的十制年份
	%z，%Z 时区名称，如果不能得到时区名称则返回空字符。
	%% 百分号

调用示例：
```
#include <time.h>
#include <sys/time.h>
#include <iostream>
#include <stdlib.h>

using namespace std;

int main()
{
	time_t dwCurTime1;
	dwCurTime1 = time(NULL);

    struct tm* pTime;
    pTime = localtime(&dwCurTime1);

    char buf[100];

    strftime(buf, 100, "time: %r, %a %b %d, %Y", pTime);

    cout << buf << endl;

	return 0;
}
```
结果：

	time: 08:18:12 PM, Fri Aug 03, 2018

###	时间函数之间的关系图

![](https://user-gold-cdn.xitu.io/2018/11/10/166fc7017db7a06b?w=1279&h=1079&f=jpeg&s=151410)

## 进程时间

进程时间是进程被创建后使用CPU的时间	，进程时间被分为以下两个部分：

- 用户CPU时间：在用户态模式下使用CPU的时间
- 内核CPU时间：在内核态模式下使用CPU的时间。这是执行内核调用或其他特殊任务所需要的时间。

### clock函数
clock函数提供了一个简单的接口用于取得进程时间，它返回一个值描述进程使用的总的CPU时间（包括用户时间和内核时间），该函数定义如下：

	#include <time.h>
	clock_t clock(void)
	//if error, return -1

clock函数返回值得计量单位是CLOCKS_PER_SEC，将返回值除以这个计量单位就得到了进程时间的秒数

### times函数
times函数也是一个进程时间函数，有更加具体的进程时间表示，函数定义如下：

	#include <sys/times.h>
	clock_t times(struct tms* buf);

	struct tms{
		clock_t tms_utime;
		clock_t tms_stime;
		clock_t tms_cutime;
		clock_t tms_cstime;
	};
times函数虽然返回类型还是clock_t，但是与clock函数返回值的计量单位不同。times函数的返回值得计量单位要通过sysconf(_SC_CLK_TCK_)来获得。

Linux系统编程手册上一个完整的使用案例如下：

	#include <time.h>
	#include <sys/times.h>
	#include <unistd.h>
	#include <stdio.h>
	
	static void displayProcessTime(const char* msg)
	{
		struct tms t;
		clock_t clockTime;
		static long clockTick = 0;
	
		if (msg != NULL) 
		{
			printf("%s\n", msg);
		}
	
		if (clockTick == 0)
		{
			clockTick = sysconf(_SC_CLK_TCK);
			if (clockTick < 0) return;
		}
	
		clockTime = clock();
		printf("clock return %ld CLOCKS_PER_SEC (%.2f seconds)\n", (long)clockTime, (double)clockTime/CLOCKS_PER_SEC);
		
		times(&t);
		printf("times return user CPU = %.2f; system CPU = %.2f\n", (double)t.tms_utime / clockTick, (double)t.tms_stime / clockTick);
	}
	
	int main()
	{
		printf("CLOCKS_PER_SEC = %ld, sysconf(_SC_CLK_TCK) = %ld\n", (long)CLOCKS_PER_SEC, sysconf(_SC_CLK_TCK));
	
		displayProcessTime("start:");
		for (int i = 0; i < 1000000000; ++i)
		{
			getpid();
		}
		printf("\n");
		displayProcessTime("end:");
	
		return 0;
	}

## 参考
[1] [http://www.runoob.com/w3cnote/cpp-time_t.html](http://www.runoob.com/w3cnote/cpp-time_t.html) 

[2] Unix高级环境编程(第三版)

[3] Unix系统编程手册
