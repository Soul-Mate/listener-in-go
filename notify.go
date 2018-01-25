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
	fmt.Println("event: ", event.Op.String())
	switch {
	case event.Op&fsnotify.Create == fsnotify.Create:
		ParseListener(event.Name, config)
	case event.Op&fsnotify.Chmod == fsnotify.Chmod:
		ParseListener(event.Name, config)
	case event.Op&fsnotify.Remove == fsnotify.Remove:
		ParseListener(event.Name, config)
	case event.Op&fsnotify.Write == fsnotify.Write:
		ParseListener(event.Name, config)
	}
}

func ParseListener(file string, config *Config) {
	var (
		matched bool
		err     error
	)
	matched, err = regexp.MatchString(`match[0-9]+\.json`, file)
	Check(err)
	if matched {
		fmt.Println("matched match file")
	}

	matched, err = regexp.MatchString(`market[0-9]+\.json`, file)
	Check(err)
	if matched {
		fmt.Println("matched market file")
	}

	switch file {
	case config.Listener.StaticFiles.League:
		fmt.Println("matched league static file")
	case config.Listener.StaticFiles.MatchFull:
		fmt.Println("matched match_full static file")
	case config.Listener.StaticFiles.MarketFull:
		fmt.Println("matched match_full static file")
	case config.Listener.StaticFiles.LeaguesFull:
		fmt.Println("matched leagues_full static file")
	}
}
