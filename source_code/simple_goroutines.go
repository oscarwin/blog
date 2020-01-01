package main

import "fmt"
import "time"

func main() {
	fmt.Println("main start")
	go echo()
	fmt.Println("main end")
	time.Sleep(500 * time.Millisecond)
}

func echo() {
	fmt.Println("echo a line")
}
