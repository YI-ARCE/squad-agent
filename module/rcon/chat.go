package rcon

import (
	"squad/module/rcon/global"
	log2 "squad/module/rcon/log"
	"squad/module/rcon/orm"
	"squad/module/rcon/server"
	"strconv"
	"strings"
	"yiarce/core/date"
	"yiarce/core/log"
)

func chatEvent(message ChatMessage) {
	server.PushAll(server.TypeChatEvent, message)
	switch message.Type {
	case `ChatAdmin`:
		return
	default:
		switch strings.ToLower(message.Message) {
		case `qd`, `签到`:
			if global.NowGameInfo.GameUserSign[message.Platform[`eosID`]] {
				Send(2, 2, `AdminWarn `+message.Platform[`eosID`]+` 本局已签到`)
				return
			}
			point, err := orm.SignIn(message.Platform)
			if err != nil {
				Send(2, 2, `AdminWarn `+message.Platform[`eosID`]+` `+err.Error())
				return
			} else {
				global.NowGameInfo.GameUserSign[message.Platform[`eosID`]] = true
				Send(2, 2, `AdminWarn `+message.Platform[`eosID`]+` 签到成功,积分剩余: `+strconv.Itoa(point))
			}
		case `tb`, `跳边`:
			if global.NowGameInfo.GameAllJumpStatus >= 10 || global.NowGameInfo.GameUserJumpStatus[message.Platform[`eosID`]] {
				Send(2, 2, `AdminWarn `+message.Platform[`eosID`]+` 本局跳边次数已用尽`)
			} else {
				global.NowGameInfo.GameAllJumpStatus++
				global.NowGameInfo.GameUserJumpStatus[message.Platform[`eosID`]] = true
				Send(2, 2, `AdminForceTeamChange `+message.Platform[`eosID`])
				Send(2, 2, `AdminWarn `+message.Platform[`eosID`]+` 跳边成功`)
			}
		case `killed`:
			user := log2.PlayerInfos[message.Platform[`eosID`]]
			if user != nil {
				if user.Killed != `` && (date.Date().Source().Unix()-user.KilledTime) > 10 {
					atkUser := log2.PlayerInfos[user.Killed]
					if atkUser != nil {
						Send(2, 2, `AdminWarn `+message.Platform[`eosID`]+` [`+atkUser.NickName+`]K:`+strconv.Itoa(atkUser.Kill)+`-D:`+strconv.Itoa(atkUser.Die)+`-A:`+strconv.Itoa(atkUser.Saved))
					}
				}
			}
		case `op`:
			if date.Date().Unix()-global.NowGameInfo.OpBroadcastTime > 5 {
				str := `当前在线OP:` + "\n"
				index := 0
				for _, i2 := range log2.PlayerInfos {
					if i2.Vip == 6 {
						str += strconv.Itoa(index+1) + `. ` + i2.NickName + "\n"
						index++
					}
				}
				if index == 0 {
					command(`AdminBroadcast 当前暂无OP在线!`)
				} else {
					command(`AdminBroadcast ` + str)
				}
			}
		case `jftb`, `积分跳边`:
			r := orm.Table(`user u`).Join(`user_point up`, `u.u_id = up.u_id`).
				Where(`u.u_steam`, message.Platform[`steamID`]).
				Field(`u.u_id,up.u_points`).
				FetchSql().
				Find()
			if r.Err() != nil {
				log.Default(r.Sql())
				Send(2, 2, `AdminWarn `+message.Platform[`eosID`]+` 跳边失败`)
			}
			result := r.Result()
			if result[`u_id`] == `` {
				Send(2, 2, `AdminWarn `+message.Platform[`eosID`]+` 跳边失败,账户未建立`)
			} else {
				num, _ := strconv.Atoi(result[`u_points`])
				if num < 5 {
					Send(2, 2, `AdminWarn `+message.Platform[`eosID`]+` 积分不足,跳边失败`)
				} else {
					ru := orm.Table(`user_point`).Where(`u_id`, result[`u_id`]).Update(map[string]interface{}{
						`u_points`: num - 5,
					})
					if ru.Err() != nil {
						Send(2, 2, `AdminWarn `+message.Platform[`eosID`]+` 跳边失败,无法操作`)
					} else {
						Send(2, 2, `AdminForceTeamChange `+message.Platform[`steamID`])
						Send(2, 2, `AdminWarn `+message.Platform[`eosID`]+` 跳边成功 剩余积分:`+strconv.Itoa(num-5))
					}
				}
			}
		case `sor`, `sr`, `sorry`, `sry`, `对不起`, `抱歉`:
			if global.NowGameInfo.GameUserTK[message.Platform[`eosID`]] > 0 {
				global.NowGameInfo.GameUserTK[message.Platform[`eosID`]] -= 1
				if global.NowGameInfo.GameUserTK[message.Platform[`eosID`]] < 1 {
					Send(2, 2, `AdminWarn `+message.Platform[`eosID`]+` 感谢配合`)
				} else {
					Send(2, 2, `AdminWarn `+message.Platform[`eosID`]+` 还需道歉`+strconv.Itoa(global.NowGameInfo.GameUserTK[message.Platform[`eosID`]])+`次`)
				}
			}
		case `jfcx`, `cxjf`, `jf`, `积分查询`:
			r := orm.Table(`user_point up`).Join(`user u`, `u.u_id = up.u_id`).
				Where(`u.u_steam`, message.Platform[`steamID`]).
				Field(`u.u_id,up.u_points`).
				Find()
			if r.Err() != nil {
				Send(2, 2, `AdminWarn `+message.Platform[`eosID`]+` 积分查询失败`)
				return
			}
			result := r.Result()
			if result[`u_id`] == `` {
				Send(2, 2, `AdminWarn `+message.Platform[`eosID`]+` 无积分`)
				return
			}
			Send(2, 2, `AdminWarn `+message.Platform[`eosID`]+` 当前剩余积分: `+result[`u_points`])
		}
	}
}
