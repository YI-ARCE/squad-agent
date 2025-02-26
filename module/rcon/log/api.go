package log

import (
	"bytes"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io"
	"os"
	command2 "squad/module/rcon/command"
	"squad/module/rcon/global"
	"strings"
	"time"
	"unsafe"
	"yiarce/core/file"
	"yiarce/core/frame"
	"yiarce/core/timing"
	"yiarce/core/yorm"
)

var skew int64 = 3
var logBytes []byte
var first = true
var db func(name string) yorm.ModelTransfer
var getUid func(info *PlayerInfo)

func Monitor(orm func(name string) yorm.ModelTransfer, uidFunc func(info *PlayerInfo), register func()) {
	db = orm
	getUid = uidFunc
	path := global.ServerConfig[`log_path`]
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(`监听服务启动失败`)
	}
	err = watcher.Add(path)
	if err != nil {
		panic(err.Error())
	}
	frame.Println(`初始化读取日志中...`)
	firstRead(path)
	first = false
	frame.Println(`首次读取完成...`, `文件偏移:`, skew)
	frame.Println(`稍后启动监听`)
	time.Sleep(time.Second * 2)
	timing.Anonymous(func() bool {
		command2.ListPlayers()
		time.Sleep(time.Second * 2)
		command2.ListSquads()
		time.Sleep(time.Second * 1)
		command2.ShowCurrentMap()
		return true
	}, time.Second*7).Start()
	frame.Println(`开始监听文件更新`)
	register()
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				readNewContent(path)
			}
			// 根据event.Name和event.Op进行相应的处理
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Println("LogError:", err)
		}
	}

}

func firstRead(filename string) {
	get, err := file.Get(filename)
	if err != nil {
		frame.Println(err.Error())
		os.Exit(0)
	}
	get.Index = skew
	get.ReadLine(func(raw string) {
		for _, v := range logEvent {
			if strings.Contains(raw, v) {
				event := switchEventBase(v)
				if event.MatchAll(raw) {
					break
				}
			}
		}
	})
	skew = get.Index
}

func readNewContent(filename string) {
	f, err := file.GetCustom(filename, os.O_RDONLY)
	defer f.Close()
	if err != nil {
		frame.Println(err.Error())
		return
	}
	fc := f.GetReader()
	fc.Seek(skew, 0)
	var buf []byte
	for {
		tmp := make([]byte, 4096)
		n, err := fc.Read(tmp)
		if err != nil && err != io.EOF {
			frame.Println(err)
		}
		if n == 0 {
			break
		} else {
			skew += int64(n)
			buf = append(buf, tmp[:n]...)
		}
	}
	parse(buf)
}

func parse(b []byte) {
	d := bytes.Split(b, []byte("\n"))
	if len(logBytes) > 0 {
		d[0] = append(logBytes, d[0]...)
		logBytes = []byte{}
	}
	//n := 0
	for _, v := range d {
		if len(v) < 2 {
			continue
		}
		//n++
		v = bytes.ReplaceAll(v, []byte("\r"), []byte{})
		str := *(*string)(unsafe.Pointer(&v))
		for _, event := range logEvent {
			if strings.Contains(str, event) {
				event := switchEventBase(event)
				if event.MatchAll(str) {
					break
				}
			}
		}
	}
	//frame.Println(`本次读到`, n, `行`)
}
