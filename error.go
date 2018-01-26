package main

import (
	"log"
)
// 打印错误
func Check(err error)  {
	if err != nil {
		log.Fatal(err)
	}
}