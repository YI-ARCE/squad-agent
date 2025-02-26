package server

import (
	"squad/module/rcon/orm"
	"yiarce/core/date"
	"yiarce/core/frame"
)

// 检查管理员账户
func checkAdmin(loginRequest LoginRequest) LoginResponse {
	r := orm.Table(`admin`).Where(`a_username`, loginRequest.User).Where(`a_password`, loginRequest.Password).FetchSql().Field(`a_id,a_status,ar_id`).Find()
	if r.Err() != nil {
		frame.Println(r.Err().Error())
		frame.Println(r.Sql())
	}
	res := r.Result()
	lr := LoginResponse{LoginStatus: false}
	if res[`a_id`] == `` {
		return lr
	}
	if res[`ar_id`] == `0` {
		lr.Msg = `账户暂未启用`
		return lr
	}
	if res[`a_status`] == `4` {
		lr.Msg = `账户禁止使用`
		return lr
	}
	lr.LoginStatus = true
	lr.ExpireTime = date.Date().Source().UnixMilli() + 86400
	rt := orm.Table(`admin_role`).Where(`ar_id`, res[`ar_id`]).Find()
	lr.Role = rt.Result()
	lr.AId = res[`a_id`]
	orm.Table(`admin`).Where(`a_username`, loginRequest.User).Where(`a_password`, loginRequest.Password).Update(map[string]string{
		`last_time`: date.Date().Timestamp(`s`),
	})
	return lr
}
