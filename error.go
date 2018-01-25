package main

import (
	"log"
	"runtime/debug"
)
// 打印错误
// 打印调用栈
func Check(err error)  {
	if err != nil {
		log.Fatal(err)
		debug.PrintStack()
	}
}