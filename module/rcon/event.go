package rcon

import (
	"fmt"
	"squad/module/rcon/global"
	"squad/module/rcon/log"
	"squad/module/rcon/orm"
	"squad/module/rcon/server"
	"squad/module/rcon/server/api"
	"strconv"
	"yiarce/core/date"
	"yiarce/core/frame"
	"yiarce/core/yorm"
)

func OnAdminBroadcast(c *log.AdminBroadcast) {
	server.PushAll(server.TypeGameEvent, c, `admin_broadcast`)
}

func OnDeployableDamaged(c *log.DeployableDamaged) {
	//fmt.Println(c.PlayerSuffix, `使用武器:`, c.Weapon, `打中`, c.Deployable, `造成`, fmt.Sprintf(`%.2f`, c.Damage), `点伤害`, `剩余血量:`, c.HealthRemaining, `伤害类型:`, c.DamageType)
	//server.PushAll(server.TypeGameEvent, c, `deployable_damaged`)
}

func OnSquadCreate(c *log.PlayerCreateSquad) {
	orm.Table(`game_bill`).Insert(map[string]interface{}{
		`gb_type`:     5,
		`atk_u_id`:    c.PlayerInfo.UID,
		`victim_u_id`: 0,
		`gb_msg`:      c.PlayerInfo.NickName + `创建了队伍 ` + c.SquadName,
		`log_time`:    c.Time / 1000,
		`gt_id`:       global.NowGameInfo.Tag,
		`create_time`: date.Date().Unix(),
	})
}

func OnNewGame(c *log.NewGame) {
	global.NowGameInfo.Restart()
	r := orm.Table(`game_tag`).Insert(map[string]string{
		`gt_map`:      c.MapClassname,
		`gt_layer`:    c.LayerClassname,
		`log_time`:    strconv.FormatInt(c.Time/1000, 10),
		`create_time`: strconv.Itoa(date.Date().Unix()),
	})
	orm.Table(`game_bill`).Insert(map[string]interface{}{
		`gb_type`:     3,
		`gb_msg`:      `当前对局地图: ` + c.LayerClassname,
		`gt_id`:       global.NowGameInfo.Tag,
		`log_time`:    c.Time / 1000,
		`create_time`: date.Date().Unix(),
	})
	if r.Err() != nil {
		frame.Println(`新对局记录操作失败`, r.Err().Error())
	}
	global.NowGameInfo.Tag = r.Id()
	go func() {
		command(`ShowCurrentMap`)
	}()
	//server.PushAll(server.TypeGameEvent, c, `new_game`)
}

func OnPlayerConnected(c *log.PlayerConnected) {
	//fmt.Println(`玩家:`, c.Platform[`eosID`], `连接服务器中`, `IP:`, c.Ip, `steamID:`, c.Platform[`steamID`])
	server.PushAll(server.TypeGameEvent, c, `player_connected`)
}

func OnPlayerDisconnected(c *log.PlayerDisconnected) {
	//fmt.Println(`玩家:`, c.PlayerInfo.NickName, `伤心欲绝的退出了游戏`)
	c.Extra[`nickname`] = c.PlayerInfo.NickName
	server.PushAll(server.TypeGameEvent, c, `player_disconnected`)
	orm.SettleGameTime(c.PlayerInfo)
	orm.Table(`game_bill`).Insert(map[string]interface{}{
		`gb_type`:     7,
		`atk_u_id`:    c.PlayerInfo.UID,
		`victim_u_id`: 0,
		`gb_msg`:      c.PlayerInfo.NickName + `退出了游戏`,
		`log_time`:    c.Time / 1000,
		`gt_id`:       global.NowGameInfo.Tag,
		`create_time`: date.Date().Unix(),
	})
	// 使用后清空用户信息回收内存
	c.PlayerInfo = nil
}

func OnPlayerDamaged(c *log.PlayerDamaged) {
	//fmt.Println(c.AttackName, `使用`, c.Weapon, `武器对`, c.WoundedName, `造成`, fmt.Sprintf(`%.2f`, c.Damage), `点伤害`)
	orm.Table(`game_bill`).Insert(map[string]interface{}{
		`gb_type`:     9,
		`atk_u_id`:    0,
		`victim_u_id`: 0,
		`gb_msg`:      c.AttackName + `使用` + c.Weapon + `对` + c.WoundedName + `造成` + fmt.Sprintf(`%.2f`, c.Damage) + `点伤害`,
		`log_time`:    c.Time / 1000,
		`gt_id`:       global.NowGameInfo.Tag,
		`create_time`: date.Date().Unix(),
	})
	server.PushAll(server.TypeGameEvent, c, `player_damaged`)
}

func OnPlayerDied(c *log.PlayerDied) {
	atkUser := c.PlayerInfo
	victimUser := c.VictimPlayerInfo
	pushData := map[string]interface{}{}
	if atkUser != nil && victimUser != nil {
		pushData[`atk_user`] = atkUser.NickName
		pushData[`vic_user`] = victimUser.NickName
		pushData[`time`] = c.Time
		pushData[`flag`] = false
		self := atkUser.NickName != victimUser.NickName
		flag := atkUser.TeamID != victimUser.TeamID
		if self && flag {
			orm.Table(`user_game_chess`).Where(`u_id`, atkUser.UID).Update(map[string]interface{}{
				`ugc_kill`: yorm.Raw(`ugc_kill+1`),
			})
			orm.Table(`user_game_chess`).Where(`u_id`, victimUser.UID).Update(map[string]interface{}{
				`ugc_die`: yorm.Raw(`ugc_die+1`),
			})
		}
		if self && !flag {
			pushData[`flag`] = true
		}
		server.PushAll(server.TypeGameEvent, pushData, `player_died`)
	}
}

func OnPlayerPossess(c *log.PlayerPossess) {
	//fmt.Println(c.Player.NickName, `选择了`, c.PossessClassname, `角色`, `加入战局`)
	c.Player.PossessClassname = c.PossessClassname
	//server.PushAll(server.TypeGameEvent, c, `player_possess`)
}

func OnPlayerRevived(c *log.PlayerRevived) {
	//fmt.Println(c.ReviverName, `妙手回春拯救了`, c.VictimName, `请击杀更多敌对玩家以表谢意~`)
	orm.Table(`user_game_chess`).Where(`u_id`, c.RevivedPlayer.UID).Update(map[string]interface{}{
		`ugc_rescue`: yorm.Raw(`ugc_rescue+1`),
	})
	server.PushAll(server.TypeGameEvent, c, `player_revived`)
}

func OnPlayerUnPossess(c *log.PlayerUnPossess) {
	//fmt.Println(c.PlayerSuffix, `放弃等待救援,复活后祝您干翻对面~`)
	//server.PushAll(server.TypeGameEvent, c, `player_un_possess`)
}

func OnPlayerWounded(c *log.PlayerWounded) {
	//fmt.Println(c.PlayerSuffix, `被武器`, c.Weapon, `命中`, `损失`, fmt.Sprintf(`%.2f`, c.Damage), `生命值`)
	//server.PushAll(server.TypeGameEvent, c, `player_wounded`)
}

func OnRoundWinner(c *log.RoundWinner) {
}

func OnRoundEnded(c *log.RoundEnded) {
	if c.Status == `已结束` {
		if global.NowGameInfo.RandomPlayer {
			frame.Println(`设置了自动配平,执行中...`)
			go func() {
				api.RandomTeamPlayerNow()
				server.PushAll(server.TypeApiResponse, map[string]bool{
					`flag`: false,
				}, `randomTeamPlayer`)
			}()
		}
	}
	//fmt.Println(`当前战局` + c.Status)
	server.PushAll(server.TypeGameEvent, c, `round_ended`)
}

func OnRoundTickets(c *log.RoundTickets) {
	//fmt.Println(`当前地图`, c.MapName+`(`, c.MapClass+`) 中`, c.Winner+`(`+c.WinnerGroup+`)`, `以`, c.Point, `点数取得了胜利,他们是战无不胜哒!`)
	server.PushAll(server.TypeGameEvent, c, `round_tickets`)
}

func OnServerTickRate(c *log.ServerTickRate) {
	//fmt.Println(`当前服务器处理效率为`, c.TickRate, `状态代表`, c.Status)
	if c.TickRate < 30 {
		server.PushAll(server.TypeGameEvent, c, `server_tick_rate`)
	}
}

func OnPlayerJoinSucceeded(c *log.PlayerJoinSucceeded) {
	// 存入用户平台ID
	c.Extra[`nick_name`] = c.PlayerInfo.NickName
	c.Extra[`platform`] = c.PlayerInfo.Platform
	c.Extra[`uid`] = c.PlayerInfo.UID
	c.Extra[`vip`] = c.PlayerInfo.Vip
	//frame.Println(`玩家:`, c.PlayerInfo.NickName+`(`+c.PlayerInfo.Platform[`eosID`]+`)`, `加入游戏`, `VIP`, c.PlayerInfo.Vip)
	orm.Table(`game_bill`).Insert(map[string]interface{}{
		`gb_type`:     4,
		`atk_u_id`:    c.PlayerInfo.UID,
		`victim_u_id`: 0,
		`gb_msg`:      c.PlayerInfo.NickName + `加入了游戏`,
		`gt_id`:       global.NowGameInfo.Tag,
		`log_time`:    c.Time / 1000,
		`create_time`: date.Date().Unix(),
	})
	if c.PlayerInfo.Vip > 0 {
		if c.PlayerInfo.Vip == 6 {
			command(`AdminBroadcast 欢迎管理员用户: ` + c.PlayerInfo.NickName + ` 加入游戏!`)
		} else {
			command(`AdminBroadcast 欢迎尊贵的会员用户: ` + c.PlayerInfo.NickName + ` 加入游戏!`)
		}
	}
	//fmt.Println(`玩家:`, c.PlayerInfo.NickName+`(`+c.PlayerInfo.Platform[`eosID`]+`)`, `加入游戏`, `UID`, c.PlayerInfo.UID)
	server.PushAll(server.TypeGameEvent, c, `player_join_succeeded`)
}
