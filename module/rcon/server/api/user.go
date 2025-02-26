package api

import (
	"fmt"
	"os"
	"squad/module/rcon/command"
	"squad/module/rcon/log"
	"squad/module/rcon/orm"
	"strconv"
	"strings"
	"yiarce/core/date"
	"yiarce/core/file"
	"yiarce/core/frame"
)

func userList(data string, aId string) string {
	q := decodeInterface(data)
	times := strconv.Itoa(date.Date().Unix())
	o := orm.Table(`user u`).
		Join(`user_game_chess ugc`, `u.u_id = ugc.u_id`, `left`).
		Join(`server_role sr`, `u.sr_id = sr.sr_id`, `left`).
		Join(`user_point up`, `u.u_id = up.u_id`, `left`)
	for k, v := range q[`query`].(map[string]interface{}) {
		if k == `u_name` {
			o.Where(`locate('` + v.(string) + `',u_name) > 0`)
			continue
		}
		if k == `flag` {
			o.Where(`(u.u_vip_expire > 0 && u.u_vip_expire < ` + times + `) or (u.sr_expire > 0 && u.sr_expire < ` + times + `)`)
			continue
		}
		if k == `ban` {
			vv := v.(string)
			if vv == `1` {
				o.Where(`u_black_info != ''`)
			}
			if vv == `2` {
				o.Where(`u_black_info != '' && black_time = 0`)
			}
			continue
		}
		o.Where(k, v)
	}
	r := o.Order(`u_id`, `desc`).Page(int(q[`page`].(float64)), int(q[`num`].(float64))).
		Field(`u.u_id,u_name,u_vip_level,sr_expire,u_vip_expire,u_steam,last_time,online_time,u_eos,ugc.ugc_kill,ugc.ugc_die,ugc.ugc_rescue,up.u_points,black_time,u_black_info,u.sr_id,sr.sr_name`).
		FetchSql().
		Select()
	if r.Err() != nil {
		frame.Println(r.Sql())
		return errors(r.Err().Error())
	}
	res := r.Result()
	ra := orm.Table(`user`)
	for k, v := range q[`query`].(map[string]interface{}) {
		if k == `u_name` {
			ra.Where(`locate('` + v.(string) + `',u_name) > 0`)
			continue
		}
		if k == `flag` {
			ra.Where(`(u.u_vip_expire > 0 && u.u_vip_expire < ` + times + `) or (u.sr_expire > 0 && u.sr_expire < ` + times + `)`)
			continue
		}
		if k == `ban` {
			vv := v.(string)
			if vv == `1` {
				ra.Where(`u_black_info != ''`)
			}
			if vv == `2` {
				ra.Where(`u_black_info != '' && black_time = 0`)
			}
			continue
		}
		ra.Where(k, v)
	}
	rc := ra.Field(`count(*) num`).Find()
	result := map[string]interface{}{
		`num`:  rc.Result()[`num`],
		`list`: res,
	}
	return success(result)
}

func userBan(data string, aId string) string {
	q := decode(data)
	r := orm.Table(`user`).Where(`u_id`, q[`u_id`]).Update(map[string]string{
		`black_time`:   q[`black_time`],
		`u_black_info`: q[`u_black_info`],
	})
	if r.Err() != nil {
		return errors(r.Err().Error())
	}
	if log.PlayerInfos[q[`u_steam`]] != nil {
		t, _ := strconv.ParseInt(q[`black_time`], 10, 64)
		command.AdminKick(q[`u_steam`], q[`u_black_info`]+`-封禁至:`+date.Time(t).YMD(`-`))
		delete(log.PlayerInfos, q[`u_eos`])
	}
	if q[`u_black_info`] != `` {
		setLog(aId, `封禁了玩家 -> `+q[`u_id`]+`,理由:`+q[`u_black_info`], `封禁玩家`, `4`)
	} else {
		if q[`black_time`] == `0` {
			setLog(aId, `解封了玩家 -> `+q[`u_id`], `解封玩家`, `4`)
		} else {
			setLog(aId, `封禁了玩家 -> `+q[`u_id`]+`,理由:`+q[`u_black_info`], `封禁玩家`, `4`)
		}
	}
	return success(nil)
}

func userBan2(data string, aId string) string {
	q := decodeInterface(data)
	day := int64(q[`time`].(float64))
	tt := int64(date.Date().Unix()) + day*86400
	if day == -1 {
		tt = 0
	}
	blackInfo := q[`u_black_info`].(string)
	r := orm.Table(`user`).Where(`u_id`, q[`u_id`]).FetchSql().Update(map[string]interface{}{
		`black_time`:   tt,
		`u_black_info`: blackInfo,
	})
	if r.Err() != nil {
		return errors(r.Err().Error())
	}
	eosId := q[`eos`].(string)
	steamId := q[`steam`].(string)
	if log.PlayerInfos[eosId] != nil {
		if tt == 0 {
			command.AdminKick(steamId, blackInfo+`-永久封禁`)
		} else {
			command.AdminKick(steamId, blackInfo+`-封禁至:`+date.Time(tt).YMD(`-`))
		}
		delete(log.PlayerInfos, eosId)
	}
	setLog(aId, `封禁了玩家 ->`+fmt.Sprintf(`%v`, q[`u_id`]), `封禁玩家`, `4`)
	return success(nil)
}

func userRemoveBan(data string, aId string) string {
	q := decode(data)
	r := orm.Table(`user`).Where(`u_id`, q[`u_id`]).Update(map[string]string{
		`black_time`:   `0`,
		`u_black_info`: ``,
	})
	if r.Err() != nil {
		return errors(r.Err().Error())
	}
	setLog(aId, `解封了玩家 ->`+fmt.Sprintf(`%v`, q[`u_id`]), `解封玩家`, `4`)
	return success(nil)
}

func userVip(data string, aId string) string {
	q := decode(data)
	r := orm.Table(`user`).Where(`u_id`, q[`u_id`]).Update(map[string]string{
		`u_vip_level`:  q[`u_vip_level`],
		`u_vip_expire`: q[`u_vip_expire`],
	})
	if r.Err() != nil {
		return errors(r.Err().Error())
	}
	setLog(aId, `设置玩家VIP角色 ->`+q[`u_id`], `设置VIP`, `4`)
	return success(nil)
}

func setUserRole(data string, aId string) string {
	d := decode(data)
	r := orm.Table(`user`).Where(`u_id`, d[`u_id`]).FetchSql().Update(d)
	if r.Num() < 0 {
		if r.Err() != nil {
			return errors(r.Err().Error())
		}
		return errors(`更新失败`)
	}
	setLog(aId, `设置VIP角色权限 ->`+d[`u_id`], `设置VIP角色权限`, `4`)
	return success(nil)
}

func outAdminConfig(data string, aId string) string {
	users := orm.Table(`user`).Where(`sr_id > 0`).Field(`u_steam,u_name,sr_id`).Select().Result()
	result := orm.Table(`server_auth`).Field(`sa_id,sa_value`).Select().Result()
	auths := make(map[string]string)
	for _, v := range result {
		auths[v[`sa_id`]] = v[`sa_value`]
	}
	result = nil
	optData := `// 刷新时间` + date.Date().Custom(`Y年M月D日 H点I分S秒`) + "\n\n"
	roles := orm.Table(`server_role`).Where(`sr_auth != ''`).Field(`sr_id,sr_name,sr_value,sr_auth`).Select().Result()
	role := make(map[string]map[string]string)
	for _, v := range roles {
		role[v[`sr_id`]] = map[string]string{`name`: v[`sr_name`], `value`: v[`sr_value`], `auth`: v[`sr_auth`]}
		arr := strings.Split(v[`sr_auth`], `,`)
		group := `Group=` + v[`sr_value`] + `:`
		for _, s := range arr {
			group += auths[s] + `,`
		}
		optData += `// ` + v[`sr_name`] + "\n" + group[:len(group)-1] + "\n"
	}
	optData += "\n\n"
	roles = nil
	adminsArr := map[string][]string{}
	for _, v := range users {
		str := `Admin=` + v[`u_steam`] + `:` + role[v[`sr_id`]][`value`] + ` // ` + v[`u_name`]
		adminsArr[v[`sr_id`]] = append(adminsArr[v[`sr_id`]], str)
	}
	for s, i := range adminsArr {
		if len(i) < 1 {
			continue
		}
		optData += `// ` + role[s][`name`] + "\n"
		for _, s2 := range i {
			optData += s2 + "\n"
		}
		optData += "\n\n"
	}
	file.Set(`./squad_server/SquadGame/ServerConfig`, `Admins.cfg`, []byte(optData), os.O_TRUNC|os.O_CREATE, 0777)
	setLog(aId, `刷新了服务器配置文件Admins.cfg`, `配置文件`, `4`)
	return success(map[string]string{
		`raw`: optData,
	})
}
