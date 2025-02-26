package api

import "fmt"

func Defer(f func(data string, aId string) string, data string, aId string) (raw string) {
	defer func() {
		if err := recover(); err != nil {
			raw = errors(fmt.Sprintf("%v", err))
		}
	}()
	raw = f(data, aId)
	return raw
}

var Actions = map[string]func(data string, aId string) string{
	`adminList`:            adminList,
	`createAdmin`:          createAdmin,
	`changeAdminStatus`:    changeAdminStatus,
	`deleteAdmin`:          deleteAdmin,
	`adminAuthList`:        adminAuthList,
	`adminRoleList`:        adminRoleList,
	`createAdminRole`:      createAdminRole,
	`updateAdminRoleGroup`: updateAdminRoleGroup,
	`editAdmin`:            editAdmin,
	`userList`:             userList,
	`userBan`:              userBan,
	`userBan2`:             userBan2,
	`userRemoveBan`:        userRemoveBan,
	`userVip`:              userVip,
	`dataList`:             dataList,
	`dataTag`:              dataTag,
	`userRoleList`:         userRoleList,
	`createUserRole`:       createUserRole,
	`editUserRole`:         editUserRole,
	`deleteUserRole`:       deleteUserRole,
	`getUserRoleAuth`:      getUserRoleAuth,
	`setUserRole`:          setUserRole,
	`outAdminConfig`:       outAdminConfig,
	`adminLogList`:         adminLogList,
	`randomTeamPlayer`:     randomTeamPlayer,
	`getActiveStatus`:      getActiveStatus,
	`clearHisData`:         clearHisData,
}
