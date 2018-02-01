package main

import (
	"github.com/fsnotify/fsnotify"
	"regexp"
	"time"
	"fmt"
	"sync"
)

var fileNameChan = make(chan string, 10)

var fileSyncMap = sync.Map{}

func StartListener() {
	config := GetConfig()
	done := make(chan bool)
	watcher, err := fsnotify.NewWatcher()
	Check(err)
	defer watcher.Close()

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				ReadListener(event)
			case err := <-watcher.Errors:
				Check(err)
			}
		}
	}()
	err = watcher.Add(config.Listener.RootPath)
	Check(err)

	defer close(fileNameChan)

	go func() {
		for file := range fileNameChan {
			// match
			if matched, err := regexp.MatchString(`match[0-9]+\.json$`, file); matched && err == nil {
				if v, ok := fileSyncMap.Load(file); ok {
					time.AfterFunc(time.Second*2, func() {
						fmt.Println("after parseMartchSave start....")
						ParseMatchSave(v.(string))
						fileSyncMap.Delete(v)
						fmt.Println("after parseMartchSave done....")
					})
				}
			}
			// market
			if matched, err := regexp.MatchString(`market[0-9]+\.json$`, file); matched && err == nil {
				if v, ok := fileSyncMap.Load(file); ok {
					time.AfterFunc(time.Second*2, func() {
						fmt.Println("after ParseMarketSave start....")
						ParseMarketSave(v.(string))
						fileSyncMap.Delete(v)
						fmt.Println("after ParseMarketSave done....")
					})
				}
			}

			// market full
			if file == "/home/www/marketfull.json" {
				if v, ok := fileSyncMap.Load(file); ok {
					time.AfterFunc(time.Second*15, func() {
						fmt.Println("after ParseMarketSave start....")
						ParseMarketSave(v.(string))
						fileSyncMap.Delete(v)
						fmt.Println("after ParseMarketSave done....")
					})
				}
			}
			// match full
			if file == "/home/www/matchfull.json" {
				if v, ok := fileSyncMap.Load(file); ok {
					time.AfterFunc(time.Second*10, func() {
						fmt.Println("after ParseMatchSave start....")
						ParseMatchSave(v.(string))
						fileSyncMap.Delete(v)
						fmt.Println("after ParseMatchSave done....")
					})
				}
			}
			// league full
			if file == "/home/www/leaguesfull.json" {
				if v, ok := fileSyncMap.Load(file); ok {
					time.AfterFunc(time.Second*5, func() {
						fmt.Println("after ParseLeagueSave start....")
						ParseLeagueSave(v.(string))
						fileSyncMap.Delete(v)
						fmt.Println("after ParseLeagueSave done....")
					})
				}
			}

			// league
			if file == "/home/www/league.json" {
				if v, ok := fileSyncMap.Load(file); ok {
					time.AfterFunc(time.Second*5, func() {
						fmt.Println("after ParseLeagueSave start....")
						ParseLeagueSave(v.(string))
						fileSyncMap.Delete(v)
						fmt.Println("after ParseLeagueSave done....")
					})
				}
			}
		}
	}()

	<-done
}

func ReadListener(event fsnotify.Event) {

	if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
		if _, ok := fileSyncMap.Load(event.Name); !ok {
			fmt.Println("map nount found: ", event.Name)
			fileSyncMap.Store(event.Name, event.Name)
			fileNameChan <- event.Name
		} else {
			fmt.Println("map find: ", event.Name)
		}
	}
}
