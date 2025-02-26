package global

type SquadInfo struct {
	Id       string                   `json:"id"`
	Name     string                   `json:"name"`
	Size     int                      `json:"size"`
	Lock     bool                     `json:"lock"`
	UserList []map[string]interface{} `json:"user_list"`
	Time     int64                    `json:"time"`
}

func (receiver SquadInfo) Append(data map[string]interface{}, flag bool) SquadInfo {
	if flag {
		r := []map[string]interface{}{data}
		receiver.UserList = append(r, receiver.UserList...)
	} else {
		receiver.UserList = append(receiver.UserList, data)
	}
	return receiver
}

type TeamInfo struct {
	Name  string
	Squad map[string]SquadInfo
}

type NowGame struct {
	// 当前对局中的用户跳边次数,大于5次则禁止跳边
	GameAllJumpStatus int
	// 当前对局用户跳边次数,已跳过则无法在使用命令跳
	GameUserJumpStatus map[string]bool
	// 游戏开始时间
	GameStartTime int64
	GameTeamInfo  []*TeamInfo
	GameTeamDesc  []string
	GameTeamSquad []map[string]int64
	GameUserTK    map[string]int
	GameUserSign  map[string]bool
	RandomLock    bool
	// 开局初始化标识
	First        bool
	RandomPlayer bool
	Tag          int64
	// OP在线广播
	OpBroadcastTime int
}

var ServerConfig = map[string]string{}

var TeamAbstract = map[string]string{
	`RGF`:    `俄罗斯陆军`,
	`VDV`:    `俄罗斯空降军`,
	`PLA`:    `中国人民解放军`,
	`PLAAGF`: `PLA Amphibious Ground Forces`,
	`PLANMC`: `中国人民解放军海军陆战队`,
	`USA`:    `美国陆军`,
	`USMC`:   `美国海军陆战队`,
	`ADF`:    `澳大利亚国防军`,
	`BAF`:    `英国军队`,
	`CAF`:    `加拿大军队`,
	`IMF`:    `非正规民兵`,
	`INS`:    `叛乱军队`,
	`MEA`:    `中东联军`,
	`TLF`:    `Turkish Land Forces`,
}

var TeamDesc = map[string]string{
	`俄罗斯陆军`:                        `RGF`,
	`俄罗斯空降军`:                       `VDV`,
	`中国人民解放军`:                      `PLA`,
	`PLA Amphibious Ground Forces`: `PLAAGF`,
	`中国人民解放军海军陆战队`:                 `PLANMC`,
	`美国陆军`:                         `USA`,
	`美国海军陆战队`:                      `USMC`,
	`澳大利亚国防军`:                      `ADF`,
	`英国军队`:                         `BAF`,
	`加拿大军队`:                        `CAF`,
	`非正规民兵`:                        `IMF`,
	`叛乱军队`:                         `INS`,
	`中东联军`:                         `MEA`,
	`Turkish Land Forces`:          `TLF`,
}

func (c *NowGame) Restart() {
	c.GameStartTime = 0
	c.GameUserJumpStatus = make(map[string]bool)
	c.GameAllJumpStatus = 0
	c.GameTeamInfo = make([]*TeamInfo, 2)
	c.GameTeamInfo[0] = &TeamInfo{
		Name:  ``,
		Squad: make(map[string]SquadInfo),
	}
	c.GameTeamInfo[1] = &TeamInfo{
		Name:  ``,
		Squad: make(map[string]SquadInfo),
	}
	c.GameTeamDesc = make([]string, 0)
	c.GameTeamSquad = make([]map[string]int64, 2)
	c.GameTeamSquad[0] = make(map[string]int64)
	c.GameTeamSquad[1] = make(map[string]int64)
	c.GameUserTK = make(map[string]int)
	c.GameUserSign = make(map[string]bool)
	c.First = true
	c.RandomPlayer = false
	c.Tag = 0
	c.OpBroadcastTime = 0
}

type ActiveStatus struct {
	Flag       bool
	ServerTime int64
	Diff       int64
	Expire     int64
}

var ActiveStatusInfo = ActiveStatus{}

var NowGameInfo = NowGame{}

func init() {
	NowGameInfo.Restart()
}

var funcPool = make(map[string]func(...interface{}) interface{})

func SetFunc(name string, f func(...interface{}) interface{}) {
	funcPool[name] = f
}

func ActiveFunc(name string, args ...interface{}) interface{} {
	return funcPool[name](args)
}
