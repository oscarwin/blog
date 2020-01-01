## 虚拟内存
每个进程都有自己独立的地址空间，进程的数据存储在这个地址空间中，使得进程之间相互不影响。以32位x86的Linux系统为例，其地址空间范围是0x00000000-0xFFFFFFFF，也就是4G。这4G并不是真实的4G物理内存，设想一下，如果每个进程都占用4G的物理内存，那即使再大的内存条也抗住不这样的消耗。因此，现代操作系统都采用了虚拟内存的设计。虚拟内存和物理内存都按照一定的大小分成很多页，比如每页是4k的话，4G的虚拟内存就可以分成1048576页。虚拟内存的页通过一张映射关系表与实际物理内存的页建立对应关系，如下图所示。

![ 图1 ](https://user-gold-cdn.xitu.io/2018/9/25/16610d7213771ae3?w=630&h=444&f=png&s=19910)

4G的虚拟地址空间并非都被使用，很少有程序会用满虚拟地址空间，对于使用了的虚拟地址也并不是全部都映射到物理内存。实际上只有一部分映射到物理内存，一部分映射到磁盘的交换区(swap分区)，这样就大大减少了单个程序对物理内存的占用。

进程访问虚拟地址空间的某个地址的过程是这样的：进程访问某个地址，首先根据这个地址计算得到分页，然后查询映射关系表，如果该分页映射到了物理内存帧上，则从物理内存中的数据取出返回；如果该分页没有映射到物理内存帧上，则将该分页的数据从磁盘加载到内存中，并更新映射关系表。之所以可以采用这种方式来设计有两个前提原因：
    
（1）刚被访问过的数据，很有可能被再次访问；

（2）二八原则——80%的时间都在访问20%的数据；

## Linux内存的分布

Linux的内存管理将这4GB的地址分为两个部分，一部分是内核空间的内存，一部分是用户空间的内存，内核空间占据3G-4G范围的地址，用户空间占据0G-3G范围的地址。内核空间是操作系统内核使用内存，用户空间就是运行在操作系统上的程序可以使用的内存。

![](https://user-gold-cdn.xitu.io/2018/9/25/1661105a95da211a?w=1065&h=614&f=jpeg&s=81078)

### Linux内核空间

Linux内核空间是Linux内核使用的地址空间，Linux内核总是驻留在内存中，用户级进程不能访问内核的地址空间。之前所提到的虚拟内存的页表与物理内存的映射关系表就保存在内核的地址空间里。内核空间的内存分配非常复杂，本文主要讨论用户空间的内存分配。

### Linux用户空间
Linux用户地址空间从低位到高位的顺序可以分为：文本段(Text Segment)、初始化数据段(Data Segment)、未初始化数据段(Bss Segment)、堆(Heap)、栈(Stack)和环境变量区(Environment variables)

**文本段**

用户空间的最低位是文本段，包含了程序运行的机器码。文本段具有只读属性，防止进程意外修改了指令导致程序出错。而且对于多进程的程序，可以共享同一份程序代码，这样减少了对物理内存的占用。

但是文本段并不是从0x0000 0000开始的，而是从0x0804 8000开始。0x0804 8000以下的地址是保留区，进程是不能去访问该地址段的数据，因此C语言中将为空的指针指向0。

**初始化数据段**

文本段上面就是初始化的数据段，数据段包含显示初始化的全局变量和静态变量。当程序被加载到内存中时，从可执行文件中读取这些数据的值，并加载到内存。因此，可执行文件中需要保存这些变量的值。

**Bss**

Bss段包含未初始化的全局变量和静态变量，还包含显示初始化为0的全局变量(根据编译器的实现)。当程序被加载到内存中时，这一段内存就会被初始化为0。可执行文件中只需要保存这一段内存的起始地址就行，因此减小了可执行文件的大小。

**堆**

堆从下自上增长(根据实现)，用于动态分配内存。堆的顶端成为program break，可以通过brk和sbrk函数调整堆顶的位置。c语言通过malloc函数实现动态内存分配，通过free释放分配的内存，后面会详细描述这两个函数的实现。堆上的内存通过一个双向链表进行维护，链表的每个节点保存这块内存的大小是否可用等信息。在堆上分配内存可能会导致以下问题：
（1）分配的内存，没有释放，就会导致内存泄漏；
（2）频繁的分配小块的内存有可能导致堆上都是剩余的小块的内存，这称为内存碎片；

**栈**

栈是一个动态增长和收缩的段，栈是自顶向下增长。栈由栈帧组成，每调用一个函数，系统会为每个当前调用的函数分配一个栈帧，栈帧从存储了参数的实参，以及函数中使用的局部变量，当函数返回时，该函数的栈帧就会弹出，函数中的局部变量因此也就被销毁了。

**环境变量**

在栈上面还有一小段空间，这段空间里保存的是环境变量和命令行参数，环境变量和命令行参数都是指向字符串的数组argv和environ。

### malloc和free实现

动态内存的分配是通过维护一个双向链表来实现，每个节点保存该内存块的大小的使用情况。malloc的分配有多种算法，比如首次适配原则，最优适配原则等。我们这里采用首次适配原则。实际上free函数，当堆顶有大块的内存时，会通过sbrk函数降低堆顶的地址，我们这里并不做处理。

malloc和free函数
```cpp
/* my_malloc.h */
#ifndef _MY_MALLOC_H_
#define _MY_MALLOC_H_

#include <unistd.h>
#include <stdlib.h>
#include <stdbool.h>

//保存每个内存块的信息
typedef struct _MEM_CONTROL_BLOCK_
{
    unsigned int uiBlockSize;      //当前块的大小
    unsigned int uiPrevBlockSize;  //前一个内存块的大小
    bool bIsAvalible;              //该内存块是否已经被分配内存
} MemControlBlock;

#define INIT_BLOCK_SIZE        (0x21000)          //初始化堆的大小
#define MEM_CONTROL_BLOCK_SIZE (sizeof(MemControlBlock))

void* g_pMallocStartAddr;     //维护堆底地址
void* g_pMallocEndAddr;       //维护堆顶地址

//初始化堆段
void malloc_init()
{
    g_pMallocStartAddr = sbrk(INIT_BLOCK_SIZE);
    g_pMallocEndAddr = g_pMallocStartAddr + INIT_BLOCK_SIZE;

    //初始化时堆只有一个内存块
    MemControlBlock* pFirstBlock;
    pFirstBlock = (MemControlBlock*)g_pMallocStartAddr;
    pFirstBlock->bIsAvalible = 1;
    pFirstBlock->uiBlockSize = INIT_BLOCK_SIZE;
    pFirstBlock->uiPrevBlockSize = 0;
}

void* my_malloc(unsigned int uiMallocSize)
{
    static bool bIsInit = false;
    if(!bIsInit)
    {
        malloc_init();
        bIsInit = true;
    }

    void* pCurAddr = g_pMallocStartAddr;
    MemControlBlock* pCurBlock = NULL;
    MemControlBlock* pLeaveBlock = NULL;
    void* pRetAddr = NULL;

    //判断是否到了堆顶
    while (pCurAddr < g_pMallocEndAddr)
    {
        pCurBlock = (MemControlBlock*)pCurAddr;
        if (pCurBlock->bIsAvalible)
        {
            //判断该块可用的内存大小是否满足分配的需求
            if (pCurBlock->uiBlockSize - MEM_CONTROL_BLOCK_SIZE >= uiMallocSize)
            {
                //该块分配空间后剩余的空间是否足够分配一个控制块，如果不能则把该块全部分配了
                if ((pCurBlock->uiBlockSize - MEM_CONTROL_BLOCK_SIZE) <= (uiMallocSize + MEM_CONTROL_BLOCK_SIZE))
                {
                    pCurBlock->bIsAvalible = 0;
                    pRetAddr = pCurAddr;
                    break;
                }
                else
                {
                    //分配内存，并将剩余的空间独立成一个块
                    pLeaveBlock = (MemControlBlock*)(pCurAddr + MEM_CONTROL_BLOCK_SIZE + uiMallocSize);
                    pLeaveBlock->bIsAvalible = 1;
                    pLeaveBlock->uiBlockSize = pCurBlock->uiBlockSize - MEM_CONTROL_BLOCK_SIZE - uiMallocSize;
                    pLeaveBlock->uiPrevBlockSize = MEM_CONTROL_BLOCK_SIZE + uiMallocSize;

                    pCurBlock->bIsAvalible = 0;
                    pCurBlock->uiBlockSize = MEM_CONTROL_BLOCK_SIZE + uiMallocSize;

                    pRetAddr = pCurAddr;
                    break;
                }
            }
            else
            {
                pCurAddr += pCurBlock->uiBlockSize;
                continue;
            }
        }
        else
        {
            pCurAddr += pCurBlock->uiBlockSize;
            continue;
        }
    }

    //已有的块中找不到合适的块，则通过sbrk函数增加堆顶地址
    if (!pRetAddr)
    {
        unsigned int uiAppendMemSize = uiMallocSize + MEM_CONTROL_BLOCK_SIZE;
        unsigned int uiPrevBlockSize = pCurBlock->uiBlockSize;
        if(*((int*)sbrk(uiAppendMemSize)) == -1)
        {
            return NULL;
        }
        g_pMallocEndAddr = g_pMallocEndAddr + uiAppendMemSize;
        pCurBlock = (MemControlBlock*)pCurAddr;
        pCurBlock->bIsAvalible = 0;
        pCurBlock->uiBlockSize = uiAppendMemSize;
        pCurBlock->uiPrevBlockSize = uiPrevBlockSize;

        pRetAddr = pCurAddr;
    }

    return pRetAddr + MEM_CONTROL_BLOCK_SIZE;
}

void my_free(void* pFreeAddr)
{
    if (pFreeAddr == NULL)
    {
        return;
    }

    MemControlBlock* pCurBlock = (MemControlBlock*)(pFreeAddr - MEM_CONTROL_BLOCK_SIZE);
    MemControlBlock* pPrevBlock = (MemControlBlock*)(pFreeAddr - MEM_CONTROL_BLOCK_SIZE - pCurBlock->uiPrevBlockSize);
    MemControlBlock* pNextBlock = (MemControlBlock*)(pFreeAddr - MEM_CONTROL_BLOCK_SIZE + pCurBlock->uiBlockSize);
    if (pCurBlock->bIsAvalible == 0)
    {
        pCurBlock->bIsAvalible = 1;

        //判断前一个内存块是否可用
        if (pCurBlock->uiPrevBlockSize != 0 && pPrevBlock->bIsAvalible)
        {
            pPrevBlock->uiBlockSize += pCurBlock->uiBlockSize;

            if((void*)pNextBlock < g_pMallocEndAddr)
            {
                pNextBlock->uiPrevBlockSize = pPrevBlock->uiBlockSize;
            }
        }
    }

    return;
}

#endif //_MY_MALLOC_H_
```

**测试程序**

这个测试程序就是循环在堆上动态分配内存，然后释放内存，可以选择释放起始块的位置，也可以选择间隔的块数量。

```
/*test_malloc.c*/
#include <unistd.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include "my_malloc.h"

#define MAX_ALLOCS 1000

int main(int argc, char const *argv[])
{
    if(argc < 3 || strcmp(argv[1], "--help") == 0)
    {
        printf("Usage: %s num_allocs block_size [step [min [max]]]\n", argv[1]);
        return -1;
    }

    int iNumAllocs = atoi(argv[1]);
    int iBlockSize = atoi(argv[2]);
    if(iBlockSize > MAX_ALLOCS)
    {
        printf("Out range of the max mallocs.\n");
    }

    int iStep, iMin, iMax;
    iStep = (argc > 3) ? atoi(argv[3]) : 1;
    iMin = (argc > 4) ? atoi(argv[4]) : 1;
    iMax = (argc > 5) ? atoi(argv[5]) : iNumAllocs;

    printf("Initial program break: %10p\n", sbrk(0));

    void* pArr[iNumAllocs];
    memset(pArr, 0, sizeof(pArr));
    for(int i = 0; i < iNumAllocs; ++i)
    {
        pArr[i] = my_malloc(iBlockSize);
        if(pArr[i] == NULL)
        {
            printf("malloc failed.\n");
            return -1;
        }
        printf("After malloc, program break: %10p\n", sbrk(0));
    }

    printf("After alloc, program break: %10p\n", sbrk(0));

    for (int i = iMin; i <= iMax; i += iStep)
    {
        my_free(pArr[i - 1]);
    }

    printf("After free, program break: %10p\n", sbrk(0));

    return 0;
}
```

#### 测试结果

	$：./a.out 10 100 2 
	Initial program break:   0xa13000
	After malloc, program break:   0xa34000
	After malloc, program break:   0xa34000
	After malloc, program break:   0xa34000
	After malloc, program break:   0xa34000
	After malloc, program break:   0xa34000
	After malloc, program break:   0xa34000
	After malloc, program break:   0xa34000
	After malloc, program break:   0xa34000
	After malloc, program break:   0xa34000
	After malloc, program break:   0xa34000
	After alloc, program break:   0xa34000
	After free, program break:   0xa34000