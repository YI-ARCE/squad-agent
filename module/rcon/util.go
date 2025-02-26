package rcon

import (
	"sort"
	"squad/module/rcon/global"
	"squad/module/rcon/log"
	"strconv"
)

func EventBaseCase() EventBaseModule {
	return EventBaseModule{eventBase: []EventBase{
		EventAdminCamera{},
		EventUnAdminCamera{},
		EventAdminWarn{},
		EventKick{},
		EventBanPlayer{},
		EventCurrentMap{},
		EventListPlayers{},
		EventListSquads{},
		EventNextMap{},
	}}
}

func EventBaseCase2() EventBaseModule {
	return EventBaseModule{eventBase: []EventBase{
		EventAdminCamera{},
		EventUnAdminCamera{},
	}}
}

func getUserList(...interface{}) interface{} {
	var players []map[string]interface{}
	var keys []string
	for k := range log.PlayerInfos {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		info := log.PlayerInfos[k]
		squadId := strconv.Itoa(info.SquadID)
		name := ``
		if info.TeamID > 0 {
			name = global.NowGameInfo.GameTeamInfo[info.TeamID-1].Squad[squadId].Name
		}
		player := map[string]interface{}{
			`nick_name`:  info.NickName,
			`platform`:   info.Platform,
			`login_time`: info.LoginTime,
			`uid`:        info.UID,
			`vip`:        info.Vip,
			`arm_type`:   info.PossessClassname,
			`team_id`:    info.TeamID,
			`squad_id`:   info.SquadID,
			`squad_name`: name,
			`role`:       info.Role,
		}
		players = append(players, player)
	}
	return players
}

func getPanelSquadList() map[string]interface{} {
	team := make([]map[string]global.SquadInfo, 2)
	team[0] = global.NowGameInfo.GameTeamInfo[0].Squad
	team[1] = global.NowGameInfo.GameTeamInfo[1].Squad
	var keys []string
	for k := range log.PlayerInfos {
		keys = append(keys, k)
	}
	teamNum := []int{0, 0}
	sort.Strings(keys)
	for _, v := range keys {
		info := log.PlayerInfos[v]
		if info.TeamID == 0 {
			continue
		}
		user := map[string]interface{}{
			`nick_name`:  info.NickName,
			`platform`:   info.Platform,
			`login_time`: info.LoginTime,
			`uid`:        info.UID,
			`vip`:        info.Vip,
			`kill`:       info.Kill,
			`die`:        info.Die,
			`saved`:      info.Saved,
			`leader`:     info.Leader,
			`arm_type`:   info.PossessClassname,
			`sign`:       global.NowGameInfo.GameUserSign[info.Platform[`eosID`]],
			`role`:       info.Role,
		}
		teamNum[info.TeamID-1]++
		team[info.TeamID-1][strconv.Itoa(info.SquadID)] = team[info.TeamID-1][strconv.Itoa(info.SquadID)].Append(user, info.Leader)
	}
	return map[string]interface{}{
		`teamInfo`: []string{global.NowGameInfo.GameTeamInfo[0].Name, global.NowGameInfo.GameTeamInfo[1].Name},
		`squad`:    team,
		`team_num`: teamNum,
	}
}

func command(raw string) {
	Send(2, 2, raw)
}

func commandClient(raw string) {
	Send2(2, 2, raw)
}
