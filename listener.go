package main

import (
	"path/filepath"
	"os"
	"fmt"
	"regexp"
)

func main() {
	// 启动时将根目录文件写入
	//WalkRootPath()
	// 定时器
	TimerWrite()
	// 初始化监听器
	StartListener()
}

// 遍历监听根目录
func WalkRootPath() {
	conf := GetConfig()
	filepath.Walk(conf.Listener.RootPath, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// market
			matched ,err := regexp.MatchString(`match\d+\.json$`, file)
			if err == nil {
				if matched {
					fmt.Println(file)
				}
			}
			// match
			matched, err = regexp.MatchString(`market\d+\.json$`,file)
			if err == nil {
				if matched {
					fmt.Println(file)
				}
			}
			// league
			// league full
			// market full
			// match full
		}
		return nil
	})
}
