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
