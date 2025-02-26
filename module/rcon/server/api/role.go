package api

import (
	"squad/module/rcon/orm"
	"yiarce/core/date"
)

func userRoleList(data string, aId string) string {
	r := orm.Table(`server_role`).Select()
	return success(r.Result())
}

func createUserRole(data string, aId string) string {
	insert := decode(data)
	insert[`create_time`] = date.Date().Timestamp(`s`)
	q := orm.Table(`server_role`).Where(`sr_value`, insert[`sr_value`]).Field(`sr_id`).Find()
	if q.Result()[`sr_id`] != `` {
		return errors(`角色已存在`)
	}
	r := orm.Table(`server_role`).Insert(insert)
	if r.Id() > 0 {
		setLog(aId, `创建了玩家角色组 ->`+insert[`sr_name`], `创建玩家角色组`, `4`)
		return success(nil)
	}
	if r.Err() != nil {
		return errors(r.Err().Error())
	}
	return errors(`添加失败`)
}

func editUserRole(data string, aId string) string {
	up := decode(data)
	r := orm.Table(`server_role`).Where(`sr_id`, up[`sr_id`]).Update(up)
	if r.Num() < 1 {
		if r.Err() != nil {
			return errors(r.Err().Error())
		}
		setLog(aId, `更新了玩家角色组 ->`+up[`sr_id`], `更新玩家角色组`, `4`)
		return errors(`更新失败`)
	}
	return success(nil)
}

func deleteUserRole(data string, aId string) string {
	q := decode(data)
	r := orm.Table(`user`).Where(`sr_id`, q[`sr_id`]).Field(`u_id`).Find()
	if r.Result()[`u_id`] != `` {
		return errors(`权限组已被分配给玩家,无法删除,请先移除/更换玩家权限`)
	}
	d := orm.Table(`server_role`).Where(`sr_id`, q[`sr_id`]).Delete()
	if d.Num() < 1 {
		if d.Err() != nil {
			return errors(d.Err().Error())
		}
		return errors(`删除失败`)
	}
	setLog(aId, `删除了玩家角色组 ->`+q[`sr_id`], `删除玩家角色组`, `4`)
	return success(nil)
}

func getUserRoleAuth(data string, aId string) string {
	r := orm.Table(`server_auth`).Select()
	return success(r.Result())
}
