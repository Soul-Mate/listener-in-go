package main

import (
	"time"
	"fmt"
)

// 定时器
func TimerWrite() {
	go func() {
		lisSync := GetListenerSync()
		<-lisSync.SaveMarketC
		fmt.Println("start timer write")
		t := time.NewTimer(time.Hour * 6)
		for {
			select {
			case <-t.C:
				writeMarketFull()
				writeMatchFull()
				writeLeagueFull()
				writeLeague()
				fmt.Println("start timer write done")
				t.Reset(time.Hour * 6)
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
