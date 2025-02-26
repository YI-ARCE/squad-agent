package api

import (
	command2 "squad/module/rcon/command"
	"squad/module/rcon/global"
	"squad/module/rcon/orm"
	"strconv"
	"yiarce/core/date"
	"yiarce/core/frame"
	"yiarce/core/log"
)

func success(data interface{}) string {
	if data == nil {
		return `{"code":1, "msg":"success","data":{}}`
	}
	return `{"code":1, "msg":"success","data":` + encode(data) + `}`
}

func adminList(data string, aId string) string {
	r := orm.Table("admin a").Join(`admin_role ar`, `a.ar_id = ar.ar_id`, `left`).Field(`a.*,ar.ar_name`).Select()
	if r.Err() != nil {
		log.Default(r.Err())
		return errors(r.Err().Error())
	}
	return success(r.Result())
}

func adminRoleList(data string, aId string) string {
	r := orm.Table(`admin_role`).Where(`ar_id > 1`).FetchSql().Select()
	if r.Err() != nil {
		frame.Println(r.Err().Error())
		return errors(r.Err().Error())
	}
	return success(r.Result())
}

func adminAuthList(data string, aId string) string {
	q := decode(data)
	o := orm.Table(`admin_auth`)
	if q[`role`] != `*` {
		o.Where(`aa_id in (` + q[`role`] + `)`)
	}
	r := o.Select()
	return success(r.Result())
}

func createAdmin(data string, aId string) string {
	req := decode(data)
	r := orm.Table(`admin`).Where(`a_username`, req[`a_username`]).Field(`a_id`).Find()
	if r.Err() != nil {
		return errors(r.Err().Error())
	}
	if r.Result()[`a_id`] != `` {
		return errors(`用户名已存在,请更换`)
	}
	req[`create_time`] = strconv.Itoa(date.Date().Unix())
	rs := orm.Table(`admin`).Insert(req)
	if rs.Err() != nil {
		return errors(rs.Err().Error())
	}
	setLog(aId, `创建了面板账户 ->`+req[`a_username`], `创建面板账户`, `4`)
	return success(nil)
}

func createAdminRole(data string, aId string) string {
	req := decode(data)
	req[`ar_group`] = ``
	rs := orm.Table(`admin_role`).Insert(req)
	if rs.Err() != nil {
		return errors(rs.Err().Error())
	}
	setLog(aId, `创建了面板角色 ->`+req[`ar_name`], `创建面板角色`, `4`)
	return success(nil)
}

func updateAdminRoleGroup(data string, aId string) string {
	req := decode(data)
	rs := orm.Table(`admin_role`).Where(`ar_id`, req[`ar_id`]).Update(map[string]string{
		`ar_group`: req[`ar_group`],
	})
	if rs.Err() != nil {
		return errors(rs.Err().Error())
	}
	setLog(aId, `更新了面板角色 ->`+req[`ar_id`], `更新面板角色`, `4`)
	return success(nil)
}

func editAdmin(data string, aId string) string {
	req := decode(data)
	r := orm.Table(`admin`).Where(`a_id`, req[`a_id`]).FetchSql().Update(req)
	if r.Err() != nil {
		frame.Println(r.Err().Error())
		frame.Println(r.Sql())
	}
	setLog(aId, `修改了面板账户 ->`+req[`ar_id`], `修改面板账户`, `4`)
	return success(nil)
}

func changeAdminStatus(data string, aId string) string {
	q := decode(data)
	r := orm.Table(`admin`).Where(`a_id`, q[`a_id`]).Update(map[string]string{`a_status`: q[`status`]})
	if r.Err() != nil {
		return errors(r.Err().Error())
	}
	setLog(aId, `修改了面板账户状态 ->`+q[`a_id`]+q[`a_status`], `修改面板账户状态`, `4`)
	return success(nil)
}

func deleteAdmin(data string, aId string) string {
	q := decode(data)
	r := orm.Table(`admin`).Where(`a_id`, q[`a_id`]).Delete()
	if r.Err() != nil {
		return errors(r.Err().Error())
	}
	setLog(aId, `删除了面板账户 ->`+q[`a_id`]+q[`a_status`], `删除面板账户`, `4`)
	return success(nil)
}

func adminLogList(data string, aId string) string {
	d := decodeInterface(data)
	query := d[`query`].(map[string]interface{})
	q := orm.Table(`admin_log al`).Join(`admin a`, `al.a_id = a.a_id`)
	for key, value := range query {
		if key == `search` {
			s, _ := value.(string)
			if s != `` {
				q.Where(`locate('` + s + `',al_content) > 0`)
			}
		} else {
			q.Where(key, value)
		}
	}
	r := q.Order(`al_id`, `desc`).Page(int(d[`page`].(float64)), int(d[`num`].(float64))).Field(`al.*,a.a_nickname`).Select()
	rc := orm.Table(`admin_log`).Field(`count(*) num`).Find()
	result := map[string]interface{}{
		`list`: r.Result(),
		`num`:  rc.Result()[`num`],
	}
	return success(result)
}

func randomTeamPlayer(data string, aId string) string {
	q := decode(data)
	if q[`query`] != `1` {
		if q[`next`] != `` {
			global.NowGameInfo.RandomPlayer = q[`next`] == `1`
			setLog(aId, `设置了下局自动配平`, `配平`, `4`)
		} else {
			setLog(aId, `操作了立即配平`, `配平`, `4`)
			RandomTeamPlayerNow()
			command2.ListPlayers()
		}
	}
	return success(map[string]bool{
		`flag`: global.NowGameInfo.RandomPlayer,
	})
}
