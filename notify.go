package main

import (
	"github.com/fsnotify/fsnotify"
	"regexp"
	"time"
)

type fileRef struct {
	ref  int
	file string
}

var fileRefMap = make(map[string]*fileRef)

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
				Check(err)
			}
		}
	}()

	err = watcher.Add(config.Listener.RootPath)
	Check(err)
	<-done
}

func ReadListener(event fsnotify.Event, config *Config) {

	if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
		if event.Name == "/home/www/matchfull.json" {
			time.AfterFunc(time.Second*10, func() {
				ParseMatchSave(event.Name)
			})
		}

		if matched, err := regexp.MatchString(`match[0-9]+\.json$`, event.Name); matched && err == nil {
			CallWrite(event, ParseMatchSave)
		}

		if matched, err := regexp.MatchString(`market[0-9]+\.json`, event.Name); matched && err == nil {
			CallWrite(event, ParseMarketSave)
		}
	}
}

func CallWrite(e fsnotify.Event, f func(file string)) {
	time.AfterFunc(time.Second*3, func() {
		f(e.Name)
	})
}

func (fr *fileRef) AddFileRefValue() {
	fr.ref++
}

func GetFileRefMap(file string) (*fileRef) {
	if elem, ok := fileRefMap[file]; ok {
		return elem
	}
	return nil
}

func SetFileRefMap(file string) {
	fileRefMap[file] = new(fileRef)
	fileRefMap[file].file = file
	fileRefMap[file].ref = 1
}

func DelFileRefMap(file string) {
	delete(fileRefMap, file)
}
