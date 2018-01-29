package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"regexp"
	"time"
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
	fmt.Println("modify file: ", file)

	if matched, err := regexp.MatchString(`match[0-9]+\.json$`, file); matched && err == nil {
		time.AfterFunc(time.Second*3, func() {
			ParseMatchSave(file)
		})
	}

	if matched, err := regexp.MatchString(`market[0-9]+\.json`, file); matched && err == nil {
		time.AfterFunc(time.Second*3, func() {
			ParseMarketSave(file)
		})
	}
}
