package main

import (
	"github.com/fsnotify/fsnotify"
	"fmt"
	"regexp"
)

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
				ReadListener(event, config)
			case err := <-watcher.Errors:
				// TODO 出现错误
				Check(err)
			}
		}
	}()

	err = watcher.Add(config.Listener.RootPath)
	Check(err)
	<-done
}

func ReadListener(event fsnotify.Event, config *Config) {
	if event.Op&fsnotify.Create == fsnotify.Create ||
		event.Op&fsnotify.Write == fsnotify.Write {
		ParseListener(event.Name, config)
	}
}

func ParseListener(file string, config *Config) {
	if matched, err := regexp.MatchString(`match[0-9]+\.json`, file); matched && err == nil {
		ParseMatchSave(file)
	} else {
		Check(err)
	}

	if matched, err := regexp.MatchString(`market[0-9]+\.json`, file); matched && err == nil {
		fmt.Println("market file")
	} else {
		Check(err)
	}
	filePrefix := config.Listener.RootPath + "/"
	if file == filePrefix + config.Listener.StaticFiles.League {
		fmt.Println("matched league static file")
	}

	if file == filePrefix + config.Listener.StaticFiles.MatchFull {
		ParseMatchSave(file)
	}

	if file == filePrefix + config.Listener.StaticFiles.MarketFull {
		fmt.Println("matched match_full static file")
	}

	if file == filePrefix + config.Listener.StaticFiles.LeaguesFull {
		fmt.Println("matched leagues_full static file")
	}
	fmt.Println("modify file: ", file)
}
