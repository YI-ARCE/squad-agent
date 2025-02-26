package command

import (
	"time"
	"yiarce/core/frame"
)

var send func(id uint16, types uint32, commend string, flag ...bool) error

func Init(f func(id uint16, types uint32, commend string, flag ...bool) error) {
	send = f
}

func AdminKick(steamId string, msg string) {
	send(2, 2, `AdminKick `+steamId+` `+msg)
}

func ListPlayers() {
	send(2, 2, `ListPlayers`)
}

func ListSquads() {
	send(2, 2, `ListSquads`)
}

func AdminWarn(steamId string, msg string) {
	send(2, 2, `AdminWarn `+steamId+` `+msg)
}

func AdminBroadcast(msg string) {
	send(2, 2, `AdminBroadcast `+msg)
}

func AdminForceTeamChange(steamID string, flag ...bool) {
	for {
		err := send(2, 2, `AdminForceTeamChange `+steamID, flag...)
		if err != nil {
			frame.Println(err.Error())
			time.Sleep(time.Second)
		} else {
			break
		}
	}
}

func ShowCurrentMap() {
	send(2, 2, `ShowCurrentMap`)
}
