package main

import (
	"time"
	"fmt"
)

// 定时器
func TimerWrite() {
	go func() {
		sync := GetListenerSync()
		fmt.Println(<-sync.SaveMarketC)
		fmt.Println("start timer write")
		t := time.NewTimer(nowAdd6HourUnix())
		for {
			select {
			case <-t.C:
				time.AfterFunc(time.Second, func() {
					writeMarketFull()
					writeMatchFull()
					writeLeagueFull()
					writeLeague()
					fmt.Println("start timer write done")
					t.Reset(nowAdd6HourUnix())
				})
			}
		}
	}()
}

// 当前时间增加6小时后的unix时间戳
func nowAdd6HourUnix() time.Duration {
	return time.Duration(time.Now().Add(time.Hour * 6).Unix())
}

func writeMarketFull() {
	config := GetConfig()
	ParseMarketSave(config.Listener.RootPath + "/" + config.Listener.StaticFiles.MarketFull)
}

func writeMatchFull() {
	config := GetConfig()
	ParseMatchSave(config.Listener.RootPath + "/" + config.Listener.StaticFiles.MatchFull)
}

func writeLeagueFull() {
	config := GetConfig()
	ParseLeagueSave(config.Listener.RootPath + "/" + config.Listener.StaticFiles.LeaguesFull)
}

func writeLeague() {
	config := GetConfig()
	ParseLeagueSave(config.Listener.RootPath + "/" + config.Listener.StaticFiles.League)
}
