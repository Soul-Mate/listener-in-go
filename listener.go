package main

import (
	"path/filepath"
	"os"
	"fmt"
	"regexp"
	"time"
)

type ListenerSync struct {
	SaveMarketC chan bool
	SaveMatchC  chan int
}

var listenerSync *ListenerSync

func main() {
	// 启动时将根目录文件写入
	go WalkRootPath()
	// 定时器
	TimerWrite()
	// 初始化监听器
	StartListener()
}

func NewListenerSync() *ListenerSync {
	listenerSync := new(ListenerSync)
	listenerSync.SaveMarketC = make(chan bool)
	listenerSync.SaveMatchC = make(chan int)
	return listenerSync
}

func GetListenerSync() *ListenerSync {
	if listenerSync == nil {
		listenerSync = NewListenerSync()
	}
	return listenerSync
}


// 遍历监听根目录
func WalkRootPath() {
	sync := GetListenerSync()
	fmt.Println("walk root path...")
	config := GetConfig()
	rootPath := config.Listener.RootPath
	filepath.Walk(config.Listener.RootPath, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// match
			matched, err := regexp.MatchString(`match[0-9]+\.json$`, file)
			if err == nil {
				if matched {
					ParseMatchSave(file)
				}
			}
			// market
			matched, err = regexp.MatchString(`market[0-9]+\.json$`, file)
			if err == nil {
				if matched {
					ParseMarketSave(file)
				}
			}
			if file == rootPath+"/"+config.Listener.StaticFiles.League {
				fmt.Println(file)
				time.AfterFunc(time.Second, func() {
					ParseLeagueSave(file)
				})
			}

			if file == rootPath+"/"+config.Listener.StaticFiles.LeaguesFull {
				time.AfterFunc(time.Second, func() {
					ParseLeagueSave(file)
					fmt.Println(file)
				})
			}
			//
			if file == rootPath+"/"+config.Listener.StaticFiles.MarketFull {
				time.AfterFunc(time.Second, func() {
					ParseMarketSave(file)
				})
			}
			if file == rootPath+"/"+config.Listener.StaticFiles.MatchFull {
				time.AfterFunc(time.Second, func() {
					ParseMatchSave(file)
				})
			}
		}
		return nil
	})
	fmt.Println("walk root path done...")
	sync.SaveMarketC <- true
}
