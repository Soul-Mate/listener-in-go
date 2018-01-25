package main

import "fmt"

var config Config

func main() {
	// 初始化配置文件
	err := config.NewConfig()
	fmt.Println(config)

	// TODO 初始化mysql连接

	// TODO 初始化redis连接

	// 初始化监听器
	StartListener(&config)
	Check(err)
}
