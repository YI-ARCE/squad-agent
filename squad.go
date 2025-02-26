package main

import (
	"encoding/json"
	"os"
	"squad/module/rcon"
	"squad/module/rcon/global"
	"time"
	"yiarce/core/file"
	"yiarce/core/frame"
)

func getServerConfig() {
	f, err := file.Get(`./server.json`)
	if err != nil {
		frame.Println(`配置文件读取失败,请检查当前目录下的server.json文件`)
		os.Exit(0)
	}
	var config = map[string]string{}
	err = json.Unmarshal(f.Byte(), &config)
	if err != nil {
		frame.Println(`配置解析失败,请修改当前目录下的server.json文件`)
		os.Exit(0)
	}
	if config[`listing_port`] == `` {
		frame.Println(`管理面板连接端口(listing_port)不能为空,请修改`)
		os.Exit(0)
	}
	if config[`rcon_port`] == `` {
		frame.Println(`服务器RCON端口(rcon_port)不能为空,请修改`)
		os.Exit(0)
	}
	if config[`database`] == `` {
		frame.Println(`数据库名(database)不能为空,任意填入当前主机中唯一的数据库名,注意只用英文`)
		os.Exit(0)
	}
	if config[`log_path`] == `` {
		frame.Println(`游戏日志目录(log_path)不能为空`)
		os.Exit(0)
	}
	global.ServerConfig = config
}

func main() {
	frame.SetPackageName(`squad`)
	getServerConfig()
	frame.Println("等待游戏服务器启动中...")
	time.Sleep(time.Second * 5)
	f, err := file.Get(global.ServerConfig[`log_path`])
	if err != nil {
		frame.Println(`游戏日志文件打开失败,请检查目录是否存在或未启动`)
		frame.Println(`err:`, err.Error())
		os.Exit(0)
	}
	f.Close()
	frame.Println(`游戏服务器已启动,请稍后`)
	time.Sleep(time.Second * 5)
	rcon.Link(`127.0.0.1:`+global.ServerConfig[`rcon_port`], global.ServerConfig[`rcon_password`])
	t := make(chan int, 1)
	<-t
}
