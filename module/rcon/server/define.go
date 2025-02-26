package server

const (
	TypeClientConnect         = `client_connect`          // 客户端连接
	TypeHeartbeatPackage      = `heartbeat_package`       // 心跳包
	TypeLoginRequest          = `login_request`           // 登录请求
	TypeLoginResponse         = `login_response`          // 登录响应
	TypeActiveExpiredResponse = `active_expired_response` // 过期事件
	TypeActiveExpiredNotify   = `active_expired_notify`   // 过期提醒
	TypeChatEvent             = `chat_event`              // 聊天事件
	TypeLogEvent              = `log_event`               // 日志事件
	TypeGameEvent             = `game_event`              // 日志事件
	TypeGetUserList           = `get_user_list`           // 获取用户列表
	TypeGetGameInfo           = `get_game_info`           // 获取对局数据请求
	TypeGameInfo              = `game_info`               // 获取对局数据
	TypePlayerInfo            = `player_info`             // 返回用户列表
	TypeRconCommand           = `rcon_command`            // 命令提交
	TypeRconResponse          = `rcon_response`           // 事件响应
	TypeApiRequest            = `api_request`             // 接口请求
	TypeApiResponse           = `api_response`            // 接口响应
	TypeTeamInfo              = `team_info`               // 队伍信息
	TypeNextRandomPlayer      = `random_player`           // 跳边
)
