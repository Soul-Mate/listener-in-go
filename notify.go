package main

import (
	"github.com/fsnotify/fsnotify"
	"fmt"
	"regexp"
)

func StartListener(config *Config) {
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
	ParseListener(event.Name, config)
}

func ParseListener(file string, config *Config) {
	if matched, err := regexp.MatchString(`match[0-9]+\.json`, file); matched && err == nil {
		fmt.Println("match file")
	} else {
		Check(err)
	}

	if matched, err := regexp.MatchString(`market[0-9]+\.json`, file); matched && err == nil {
		fmt.Println("market file")
	} else {
		Check(err)
	}
	file = config.Listener.RootPath + "/" + file
	if file == config.Listener.StaticFiles.League {
		fmt.Println("matched league static file")
	}

	if file == config.Listener.StaticFiles.MatchFull {
		fmt.Println("matched match_full static file")
	}

	if file == config.Listener.StaticFiles.MarketFull {
		fmt.Println("matched match_full static file")
	}

	if file == config.Listener.StaticFiles.LeaguesFull {
		fmt.Println("matched leagues_full static file")
	}
	fmt.Println("modify file: ", file)
}
