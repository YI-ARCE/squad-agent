package rcon

import (
	"os"
	"squad/module/rcon/global"
	"strings"
	"yiarce/core/frame"
	"yiarce/core/yorm"
)

func initDataBase(config yorm.Config) *yorm.Db {
	config.Database = ``
	mysql, err := yorm.ConnMysql(config)
	if err != nil {
		frame.Println(err.Error())
		os.Exit(0)
	}
	frame.Println(`初始化数据库中...`)
	mysql.Execute(`use mysql`)
	mysql.Execute(`UPDATE user SET host = '%' WHERE User = 'root'`)
	mysql.Execute(`flush privileges`)
	mysql.Execute(`create schema ` + global.ServerConfig[`database`])
	mysql.Execute(`use ` + global.ServerConfig[`database`])
	for _, s := range databaseSql {
		arrs := strings.Split(s, `;`)
		for _, sqlRaw := range arrs {
			if strings.ReplaceAll(sqlRaw, "\n", ``) == `` {
				continue
			}
			_, _, err := mysql.Execute(sqlRaw)
			if err != nil {
				frame.Println(err.Error())
				frame.Println(sqlRaw)
				os.Exit(0)
			}
		}
	}
	return mysql
}
