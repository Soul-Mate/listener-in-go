package main

import (
	"github.com/fsnotify/fsnotify"
	"regexp"
	"time"
	"fmt"
)

type fileRef struct {
	ref  int
	file string
}

var fileRefMap map[string]*fileRef

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
	if event.Op&fsnotify.Write == fsnotify.Write {
		// 存在文件引用
		if fr := GetFileRefMap(event.Name); fr != nil {
			// 达到计数值
			if fr.ref >= 2 {
				fmt.Println("开始写入：",fr.file)
				ParseListener(*fr)
				fmt.Println("写入完毕：",fr.file)
				// 删除
				DelFileRefMap(fr.file)
				fmt.Println("删除完毕：",fr.file)
			} else {
				fmt.Println("增加引用计数：",fr.ref)
				fr.AddFileRefValue()
				fmt.Println("计数增加完毕：",fr.ref)
			}
		} else {
			fmt.Println("设置fileRef：",fr.file)
			SetFileRefMap(event.Name)
		}
	}
}

func ParseListener(fr fileRef) {

	if matched, err := regexp.MatchString(`match[0-9]+\.json$`, fr.file); matched && err == nil {
		time.AfterFunc(time.Second*3, func() {
			ParseMatchSave(fr.file)
		})
	}

	if matched, err := regexp.MatchString(`market[0-9]+\.json`, fr.file); matched && err == nil {
		time.AfterFunc(time.Second*3, func() {
			ParseMarketSave(fr.file)
		})
	}
}


func (fr *fileRef )AddFileRefValue()  {
	fr.ref++
}

func GetFileRefMap(file string) (*fileRef) {
	if elem, ok := fileRefMap[file]; ok {
		return elem
	}
	return nil
}

func SetFileRefMap(file string) {
	fr := &fileRef{
		file: file,
		ref:  1,
	}
	fileRefMap[file] = fr
}

func DelFileRefMap(file string) {
	delete(fileRefMap, file)
}


