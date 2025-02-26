package log

import (
	"regexp"
	"squad/module/rcon/command"
	"squad/module/rcon/global"
	"strconv"
	"strings"
	"time"
	"yiarce/core/date"
	"yiarce/core/frame"
)

type WarStatistics struct {
	tickets    int
	chainID    string
	level      int
	faction    any
	raw        any
	action     any
	time       any
	team       any
	subFaction any
	layer      any
}
type EventBase interface {
	Match(str string) bool
}

type PlayerInfo struct {
	// 控制器ID
	ControllerID string
	// IP地址
	IP        string
	NickName  string
	Platform  map[string]string
	LoginTime int64
	// 平台账户ID
	UID              int64
	Vip              int
	PossessClassname string
	TeamID           int
	SquadID          int
	Status           bool
	BlackTime        int64
	BlackInfo        string
	Kill             int
	Die              int
	Saved            int
	Leader           bool
	Role             string
	Killed           string
	KilledTime       int64
}

type EventBaseModule struct {
	eventBase []EventBase
}

func (e *EventBaseModule) MatchAll(raw string) bool {
	if e.eventBase == nil {
		return true
	}
	for _, base := range e.eventBase {
		if base.Match(raw) {
			return true
		}
	}
	return false
}

type EventProcessor struct {
	events map[string][]EventBase
}

type AdminBroadcast struct {
	// 时间
	Time int64 `json:"time,omitempty"`
	// 链ID
	ChainID string `json:"-"`
	// 消息ID
	Message string `json:"message,omitempty"`
	// 来自
	From  string                 `json:"from,omitempty"`
	Extra map[string]interface{} `json:"extra,omitempty"`
}

type DeployableDamaged struct {
	// 时间
	Time int64 `json:"time,omitempty"`
	// 链ID
	ChainID         string                 `json:"-"`
	Deployable      string                 `json:"deployable"`
	Damage          float64                `json:"damage"`
	Weapon          string                 `json:"weapon"`
	PlayerSuffix    string                 `json:"player_suffix"`
	DamageType      string                 `json:"damage_type"`
	HealthRemaining string                 `json:"health_remaining"`
	Extra           map[string]interface{} `json:"extra,omitempty"`
}

type NewGame struct {
	// 时间
	Time int64 `json:"time,omitempty"`
	// 链ID
	ChainID        string                 `json:"-"`
	Dlc            string                 `json:"dlc"`
	MapClassname   string                 `json:"map_classname"`
	LayerClassname string                 `json:"layer_classname"`
	Extra          map[string]interface{} `json:"extra,omitempty"`
}

type PlayerConnected struct {
	// 时间
	Time int64 `json:"time,omitempty"`
	// 链ID
	ChainID          string                 `json:"-"`
	PlayerController string                 `json:"player_controller"`
	Ip               string                 `json:"ip"`
	Platform         map[string]string      `json:"platform"`
	Extra            map[string]interface{} `json:"extra,omitempty"`
}

type PlayerDisconnected struct {
	// 时间
	Time int64 `json:"time,omitempty"`
	// 链ID
	ChainID    string                 `json:"-"`
	PlayerInfo *PlayerInfo            `json:"-"`
	Ip         string                 `json:"ip"`
	EosID      string                 `json:"eos_id"`
	Extra      map[string]interface{} `json:"extra,omitempty"`
}

type PlayerDamaged struct {
	// 时间
	Time int64 `json:"time"`
	// 链ID
	ChainID      string                 `json:"-"`
	Damage       float64                `json:"damage"`
	AttackName   string                 `json:"attack_name"`
	AttackPlayer *PlayerInfo            `json:"attack_player"`
	WoundedName  string                 `json:"wounded_name"`
	Weapon       string                 `json:"weapon"`
	Extra        map[string]interface{} `json:"extra,omitempty" json:"extra,omitempty"`
}

type PlayerDied struct {
	// 时间
	Time int64 `json:"time,omitempty"`
	// 链ID
	ChainID                  string                 `json:"-"`
	VictimName               string                 `json:"victim_name"`
	Damage                   float64                `json:"damage"`
	AttackerPlayerController string                 `json:"attacker_player_controller"`
	Weapon                   string                 `json:"weapon"`
	PlayerInfo               *PlayerInfo            `json:"-"`
	VictimPlayerInfo         *PlayerInfo            `json:"-"`
	Extra                    map[string]interface{} `json:"extra"`
}

type PlayerPossess struct {
	// 时间
	Time int64 `json:"time,omitempty"`
	// 链ID
	ChainID          string                 `json:"-"`
	PlayerSuffix     string                 `json:"player_suffix"`
	PossessClassname string                 `json:"possess_classname"`
	Player           *PlayerInfo            `json:"-"`
	Extra            map[string]interface{} `json:"extra,omitempty"`
}

type PlayerRevived struct {
	// 时间
	Time int64 `json:"time,omitempty"`
	// 链ID
	ChainID       string                 `json:"-"`
	ReviverName   string                 `json:"reviver_name"`
	VictimName    string                 `json:"victim_name"`
	RevivedPlayer *PlayerInfo            `json:"-"`
	VictimPlayer  *PlayerInfo            `json:"-"`
	Extra         map[string]interface{} `json:"extra,omitempty"`
}

type PlayerUnPossess struct {
	// 时间
	Time int64 `json:"time,omitempty"`
	// 链ID
	ChainID      string                 `json:"-"`
	PlayerSuffix string                 `json:"player_suffix"`
	PlayerInfo   *PlayerInfo            `json:"-"`
	Extra        map[string]interface{} `json:"extra,omitempty"`
}

type PlayerWounded struct {
	// 时间
	Time int64 `json:"time,omitempty"`
	// 链ID
	ChainID      string                 `json:"-"`
	PlayerSuffix string                 `json:"player_suffix"`
	Weapon       string                 `json:"weapon"`
	PlayerInfo   *PlayerInfo            `json:"-"`
	Damage       float64                `json:"damage"`
	Extra        map[string]interface{} `json:"extra,omitempty"`
}

type RoundEnded struct {
	// 时间
	Time int64 `json:"time,omitempty"`
	// 链ID
	ChainID string                 `json:"-"`
	Winner  WarStatistics          `json:"winner"`
	Loser   WarStatistics          `json:"loser"`
	Status  string                 `json:"status"`
	Extra   map[string]interface{} `json:"extra,omitempty"`
}

type RoundTickets struct {
	Time        int64                  `json:"time"`
	ChainID     string                 `json:"-"`
	Winner      string                 `json:"winner"`
	WinnerGroup string                 `json:"winner_group"`
	MapName     string                 `json:"map_name"`
	MapClass    string                 `json:"map_class"`
	Point       int64                  `json:"point"`
	Extra       map[string]interface{} `json:"extra"`
}

type RoundWinner struct {
	// 时间
	Time int64 `json:"time"`
	// 链ID
	ChainID string `json:"-"`
	// 赢家
	Winner string                 `json:"winner"`
	Extra  map[string]interface{} `json:"extra"`
}

type ServerTickRate struct {
	// 时间
	Time int64 `json:"time"`
	// 链ID
	ChainID string `json:"-"`
	// 赢家
	TickRate float64                `json:"tick_rate"`
	Status   string                 `json:"status"`
	Extra    map[string]interface{} `json:"extra,omitempty"`
}

type PlayerJoinSucceeded struct {
	// 时间
	Time int64 `json:"time,omitempty"`
	// 链ID
	ChainID string `json:"-"`
	// 踢出率
	TickRate   string                 `json:"tick_rate"`
	PlayerInfo *PlayerInfo            `json:"-"`
	Extra      map[string]interface{} `json:"extra,omitempty"`
}

type PlayerCreateSquad struct {
	// 时间
	Time int64 `json:"time,omitempty"`
	// 链ID
	ChainID string `json:"-"`
	// 踢出率
	SquadID    int                    `json:"squad_id"`
	PlayerInfo *PlayerInfo            `json:"-"`
	SquadName  string                 `json:"squad_name"`
	TeamID     int                    `json:"team_id"`
	Extra      map[string]interface{} `json:"extra,omitempty"`
}

type GameLoadTeam struct {
	// 时间
	Time int64 `json:"time,omitempty"`
	// 链ID
	ChainID string `json:"-"`
	// 踢出率
	TeamName string `json:"team_name"`
}

func (c *GameLoadTeam) Match(raw string) bool {
	m, _ := regexp.Compile(`^\[([0-9.:-]+)]\[([ 0-9]*)]LogSquad: Loaded Faction : (.+) `)
	arr := m.FindStringSubmatch(raw)
	if len(arr) < 2 {
		return false
	}
	c.TeamName = global.TeamAbstract[arr[3]]
	if len(global.NowGameInfo.GameTeamDesc) < 2 {
		global.NowGameInfo.GameTeamDesc = append(global.NowGameInfo.GameTeamDesc, c.TeamName)
		//frame.Println(`团队ID:`, len(global.NowGameInfo.GameTeamDesc), `团队名:`, c.TeamName, arr[3], len(c.TeamName))
	}
	return true
}

func (c *PlayerCreateSquad) Match(raw string) bool {
	m, _ := regexp.Compile(`^\[([0-9.:-]+)]\[([ 0-9]*)]LogSquad: (.+) \(Online IDs:([^)]+)\) has created Squad (\d+) \(Squad Name: (.+)\) on (.+)`)
	arr := m.FindStringSubmatch(raw)
	if len(arr) < 2 {
		return false
	}
	c.Time = ToTime(arr[1])
	c.PlayerInfo = PlayerInfos[GetPlatform(arr[4])[`eosID`]]
	c.SquadID, _ = strconv.Atoi(arr[5])
	c.SquadName = arr[6]
	team := strings.ReplaceAll(arr[7], "\r", ``)
	//name := ``
	if team == `Team 1` || team == `Team 2` {
		id, _ := strconv.Atoi(team[len(team)-1:])
		//name = global.NowGameInfo.GameTeamDesc[id-1]
		c.TeamID = id
	} else {
		//name = global.TeamDesc[team]
		if global.NowGameInfo.GameTeamDesc[0] == team {
			c.TeamID = 1
		}
		if global.NowGameInfo.GameTeamDesc[1] == team {
			c.TeamID = 2
		}
	}
	if c.TeamID == 0 {

	} else {
		global.NowGameInfo.GameTeamSquad[c.TeamID-1][arr[5]] = c.Time
		if !first {
			if len(global.NowGameInfo.GameTeamDesc) > 1 {
				command.AdminBroadcast(`[` + date.TimeMill(c.Time).Custom(`H:i:s`) + `]` + `[` + c.PlayerInfo.NickName + `] 创建了队伍 ` + c.SquadName + ` 在 ` + global.NowGameInfo.GameTeamDesc[c.TeamID-1])
			}
			time.Sleep(time.Second / 2)
			command.ListSquads()
			go PlayerSquadCreateFunc(c)
		}
	}
	return true
}

func (c *PlayerJoinSucceeded) Match(raw string) bool {
	m, _ := regexp.Compile(`^\[([0-9.:-]+)]\[([ 0-9]*)]LogNet: Join succeeded: (.+)`)
	arr := m.FindStringSubmatch(raw)
	if len(arr) < 2 {
		return false
	}
	c.Time = ToTime(arr[1])
	c.ChainID = arr[2]
	c.PlayerInfo = JoinRequest[c.ChainID]
	c.PlayerInfo.NickName = strings.ReplaceAll(arr[3], "\r", ``)
	c.PlayerInfo.LoginTime = date.Date().Source().UnixMilli()
	c.Extra = make(map[string]interface{})
	BaseControllers[c.PlayerInfo.ControllerID] = c.PlayerInfo
	NickNameInfo[c.PlayerInfo.NickName] = c.PlayerInfo
	PlayerInfos[c.PlayerInfo.Platform[`eosID`]] = c.PlayerInfo
	getUid(c.PlayerInfo)
	JoinRequest[c.ChainID] = nil
	if !first {
		if c.PlayerInfo.BlackTime > 0 {
			command.AdminKick(c.PlayerInfo.Platform[`steamID`], c.PlayerInfo.BlackInfo+` 封禁至 `+date.Time(c.PlayerInfo.BlackTime).YMD(`-`))
		}
		db(`user`).Where(`u_id`, c.PlayerInfo.UID).Update(map[string]interface{}{
			`last_time`: date.Date().Unix(),
		})
		go PlayerJoinSucceededFunc(c)
	}
	return true
}

func (c *ServerTickRate) Match(raw string) bool {
	m, _ := regexp.Compile(`^\[([0-9.:-]+)]\[([ 0-9]*)]LogSquad: USQGameState: Server Tick Rate: ([0-9.]+)`)
	arr := m.FindStringSubmatch(raw)
	if len(arr) < 2 {
		return false
	}
	c.Time = ToTime(arr[1])
	c.ChainID = arr[2]
	c.Extra = make(map[string]interface{})
	floatNum, _ := strconv.ParseFloat(arr[3], 64)
	c.TickRate = floatNum
	switch {
	case c.TickRate >= 50:
		c.Status = "非常好"
	case c.TickRate >= 40:
		c.Status = "良好"
	case c.TickRate >= 30:
		c.Status = "稍差"
	case c.TickRate >= 20:
		c.Status = "非常差"
	case c.TickRate >= 0:
		c.Status = "无法游玩"
	default:
		c.Status = "崩溃辣,救命"
	}
	if !first {
		go ServerTickRateFunc(c)
	}
	return true
}

func (c *RoundWinner) Match(raw string) bool {
	m, _ := regexp.Compile(`^\[([0-9.:-]+)]\[([ 0-9]*)]LogSquadTrace: \[DedicatedServer](?:ASQGameMode::)?DetermineMatchWinner\(\): (.+) won on (.+)`)
	arr := m.FindStringSubmatch(raw)
	if len(arr) < 2 {
		return false
	}
	c.Extra = make(map[string]interface{})
	//frame.Println(raw)
	//frame.Println(`RoundWinner`, arr)
	return true
}

func (c *RoundTickets) Match(raw string) bool {
	m, _ := regexp.Compile(`^\[([0-9.:-]+)]\[([ 0-9]*)]LogSquadGameEvents: Display: Team ([0-9]), (.*) \( ?(.*?) ?\) has (won|lost) the match with ([0-9]+) Tickets on layer (.*) \(level (.*)\)!`)
	arr := m.FindStringSubmatch(raw)
	if len(arr) < 2 {
		return false
	}
	c.Extra = make(map[string]interface{})
	c.Time = ToTime(arr[1])
	c.ChainID = arr[2]
	c.Winner = arr[5]
	c.WinnerGroup = arr[4]
	c.MapName = arr[9]
	num, _ := strconv.ParseInt(arr[7], 10, 64)
	c.Point = num
	c.MapClass = arr[8]
	if !first {
		db(`game_bill`).Insert(map[string]interface{}{
			`gb_type`:     2,
			`gb_msg`:      c.Winner + `(` + c.WinnerGroup + `) 结算点数:` + arr[7],
			`log_time`:    c.Time / 1000,
			`gt_id`:       global.NowGameInfo.Tag,
			`create_time`: date.Date().Unix(),
		})
		go RoundTicketsFunc(c)
	}
	return true
}

func (c *RoundEnded) Match(raw string) bool {
	m, _ := regexp.Compile(`^\[([0-9.:-]+)]\[([ 0-9]*)]LogGameState: Match State Changed from InProgress to WaitingPostMatch`)
	arr := m.FindStringSubmatch(raw)
	if len(arr) < 2 {
		return false
	}
	c.Extra = make(map[string]interface{})
	c.Time = ToTime(arr[1])
	c.ChainID = arr[2]
	if strings.Contains(raw, `to WaitingPostMatch`) {
		c.Status = `已结束`
	} else {
		c.Status = `已开始`
	}
	if !first {
		go RoundEndedFunc(c)
	}
	return true
}

func (c *PlayerWounded) Match(raw string) bool {
	m, _ := regexp.Compile(`^\[([0-9.:-]+)]\[([ 0-9]*)]LogSquadTrace: \[DedicatedServer](?:ASQSoldier::)?Wound\(\): Player:(.+) KillingDamage=(?:-)*([0-9.]+) from ([A-z_0-9]+) \(Online IDs:([^)|]+)\| Controller ID: ([\w\d]+)\) caused by ([A-z_0-9-]+)_C`)
	arr := m.FindStringSubmatch(raw)
	if len(arr) < 2 {
		return false
	}
	c.Extra = make(map[string]interface{})
	c.Time = ToTime(arr[1])
	c.ChainID = arr[2]
	c.PlayerSuffix = arr[3]
	c.Weapon = arr[8]
	floatNum, _ := strconv.ParseFloat(arr[4], 64)
	c.Damage = floatNum
	c.PlayerInfo = PlayerInfos[GetPlatform(arr[6])[`eosID`]]
	if !first {
		go PlayerWoundedFunc(c)
	}
	return true
}

func (c *PlayerUnPossess) Match(raw string) bool {
	m, _ := regexp.Compile(`^\[([0-9.:-]+)]\[([ 0-9]*)]LogSquadTrace: \[DedicatedServer](?:ASQPlayerController::)?OnUnPossess\(\): PC=(.+) \(Online IDs:([^)]+)\)`)
	arr := m.FindStringSubmatch(raw)
	if len(arr) < 2 {
		return false
	}
	c.Extra = make(map[string]interface{})
	c.Time = ToTime(arr[1])
	c.ChainID = arr[2]
	c.PlayerSuffix = arr[3]
	c.PlayerInfo = PlayerInfos[GetPlatform(arr[4])[`eosID`]]
	if !first {
		go PlayerUnPossessFunc(c)
	}
	return true
}

func (c *PlayerRevived) Match(raw string) bool {
	m, _ := regexp.Compile(`^\[([0-9.:-]+)]\[([ 0-9]*)]LogSquad: (.+) \(Online IDs:([^)]+)\) has revived (.+) \(Online IDs:([^)]+)\)\.`)
	arr := m.FindStringSubmatch(raw)
	if len(arr) < 2 {
		return false
	}
	c.Extra = make(map[string]interface{})
	c.Time = ToTime(arr[1])
	c.ChainID = arr[2]
	c.ReviverName = arr[3]
	c.VictimName = arr[5]
	c.RevivedPlayer = PlayerInfos[GetPlatform(arr[4])[`eosID`]]
	c.RevivedPlayer.Saved++
	c.VictimPlayer = PlayerInfos[GetPlatform(arr[6])[`eosID`]]
	if !first {
		go PlayerRevivedFunc(c)
	}
	return true
}

func (c *PlayerPossess) Match(raw string) bool {
	m, _ := regexp.Compile(`^\[([0-9.:-]+)]\[([ 0-9]*)]LogSquadTrace: \[DedicatedServer](?:ASQPlayerController::)?OnPossess\(\): PC=(.+) \(Online IDs:([^)]+)\) Pawn=([A-z0-9_]+)_C`)
	arr := m.FindStringSubmatch(raw)
	if len(arr) < 2 {
		return false
	}
	p := PlayerInfos[GetPlatform(arr[4])[`eosID`]]
	if p == nil {
		return false
	}
	c.Extra = make(map[string]interface{})
	c.Time = ToTime(arr[1])
	c.ChainID = arr[2]
	c.PlayerSuffix = arr[3]
	c.PossessClassname = arr[5]
	c.Player = p
	c.Player.PossessClassname = c.PossessClassname
	if !first {
		go PlayerPossessFunc(c)
	}
	return true
}

func (c *PlayerDied) Match(raw string) bool {
	m, _ := regexp.Compile(`^\[([0-9.:-]+)]\[([ 0-9]*)]LogSquadTrace: \[DedicatedServer](?:ASQSoldier::)?Die\(\): Player:(.+) KillingDamage=(?:-)*([0-9.]+) from ([A-z_0-9]+) \(Online IDs:([^)|]+)\| Contoller ID: ([\w\d]+)\) caused by ([A-z_0-9-]+)_C`)
	arr := m.FindStringSubmatch(raw)
	if len(arr) < 2 {
		return false
	}
	if strings.Contains(arr[6], `INVALID`) {
		return false
	}
	c.Extra = make(map[string]interface{})
	p := GetPlatform(arr[6])
	c.Time = ToTime(arr[1])
	c.ChainID = arr[2]
	c.VictimName = arr[3]
	c.PlayerInfo = PlayerInfos[p[`eosID`]]
	c.Weapon = arr[8]
	floatNum, _ := strconv.ParseFloat(arr[4], 64)
	c.Damage = floatNum
	c.AttackerPlayerController = arr[5]
	c.VictimName = arr[3]
	//分辨是否加了前缀的名字
	if c.VictimName[0] != ' ' {
		c.VictimName = c.VictimName[strings.Index(c.VictimName, ` `)+1:]
	} else {
		c.VictimName = c.VictimName[1:]
	}
	victimUser := NickNameInfo[c.VictimName]
	if victimUser != nil {
		victimUser.Die++
	}
	if c.PlayerInfo != nil {
		c.PlayerInfo.Kill++
	}
	if !first {
		if c.PlayerInfo != nil && victimUser != nil {
			c.VictimPlayerInfo = victimUser
			msg := ``
			if c.PlayerInfo.Platform[`eosID`] != c.VictimPlayerInfo.Platform[`eosID`] {
				msg = c.PlayerInfo.NickName + ` 使用武器: ` + c.Weapon + ` 击杀了 ` + victimUser.NickName
				if c.PlayerInfo.TeamID == c.VictimPlayerInfo.TeamID {
					msg += ` (友军击杀)`
				}
			} else {
				msg = c.PlayerInfo.NickName + ` 自杀了`
			}
			cr := db(`game_bill`).FetchSql().Insert(map[string]interface{}{
				`gb_type`:     1,
				`gb_msg`:      msg,
				`atk_u_id`:    c.PlayerInfo.UID,
				`victim_u_id`: victimUser.UID,
				`gt_id`:       global.NowGameInfo.Tag,
				`log_time`:    c.Time / 1000,
				`create_time`: date.Date().Unix(),
			})
			if cr.Err() != nil {
				frame.Println(cr.Sql())
				frame.Println(cr.Err().Error())
			}
			if c.PlayerInfo.TeamID > 0 && c.VictimPlayerInfo.TeamID > 0 && c.PlayerInfo.TeamID == c.VictimPlayerInfo.TeamID && c.PlayerInfo.Platform[`eosID`] != c.VictimPlayerInfo.Platform[`eosID`] {
				if global.NowGameInfo.GameUserTK[c.PlayerInfo.Platform[`eosID`]] < 1 {
					global.NowGameInfo.GameUserTK[c.PlayerInfo.Platform[`eosID`]] += 1
					command.AdminWarn(c.PlayerInfo.Platform[`steamID`], `你击杀了友军,请及时发送道歉字样`)
					command.AdminBroadcast(c.PlayerInfo.NickName + ` 击杀了友军,请及时发送道歉字样`)
					go func() {
						time.Sleep(time.Second * 120)
						if global.NowGameInfo.GameUserTK[c.PlayerInfo.Platform[`eosID`]] > 0 {
							global.NowGameInfo.GameUserTK[c.PlayerInfo.Platform[`eosID`]] = 0
							command.AdminKick(c.PlayerInfo.Platform[`steamID`], `踢出原因:攻击队友未道歉!`)
							//command.AdminWarn(c.PlayerInfo.Platform[`steamID`], `踢出原因:攻击队友未道歉!`)
						}
					}()
				} else {
					global.NowGameInfo.GameUserTK[c.PlayerInfo.Platform[`eosID`]] += 1
				}
			}
		}
		go PlayerDiedFunc(c)
	}
	return true
}

func (c *PlayerDamaged) Match(raw string) bool {
	m, _ := regexp.Compile(`^\[([0-9.:-]+)]\[([ 0-9]*)]LogSquad: Player:(.+) ActualDamage=([0-9.]+) from (.+) \(Online IDs:([^|]+)\| Player Controller ID: ([^ ]+)\)caused by ([A-z_0-9-]+)_C`)
	arr := m.FindStringSubmatch(raw)
	if len(arr) < 2 {
		return false
	}
	c.Extra = make(map[string]interface{})
	//frame.Println(arr)
	//frame.Println(`PlayerDamaged`, arr[1:])
	c.Time = ToTime(arr[1])
	c.ChainID = arr[2]
	floatNum, _ := strconv.ParseFloat(arr[4], 64)
	c.Damage = floatNum
	c.Weapon = arr[8]
	c.AttackPlayer = PlayerInfos[GetPlatform(arr[6])[`eosID`]]
	c.AttackName = arr[5]
	c.WoundedName = arr[3]
	if !first {
		go PlayerDamagedFunc(c)
	}
	return true
}

func (c *PlayerDisconnected) Match(raw string) bool {
	m, _ := regexp.Compile(`^\[([0-9.:-]+)]\[([ 0-9]*)]LogNet: UChannel::Close: Sending CloseBunch\. ChIndex == [0-9]+\. Name: \[UChannel\] ChIndex: [0-9]+, Closing: [0-9]+ \[UNetConnection\] RemoteAddr: ([\d.]+):[\d]+, Name: EOSIpNetConnection_[0-9]+, Driver: GameNetDriver EOSNetDriver_[0-9]+, IsServer: YES, PC: ([^ ]+PlayerController_C_[0-9]+), Owner: [^ ]+PlayerController_C_[0-9]+, UniqueId: RedpointEOS:([\d\w]+)`)
	arr := m.FindStringSubmatch(raw)
	if len(arr) < 2 {
		return false
	}
	c.Extra = make(map[string]interface{})
	c.Ip = arr[3]
	c.Time = ToTime(arr[1])
	c.EosID = arr[5]
	c.ChainID = arr[2]
	p := PlayerInfos[c.EosID]
	if p != nil {
		c.PlayerInfo = p
		c.PlayerInfo.Status = false
		delete(PlayerInfos, c.EosID)
		delete(BaseControllers, c.PlayerInfo.ControllerID)
		delete(NickNameInfo, c.PlayerInfo.NickName)
		if !first {
			go PlayerDisconnectedFunc(c)
		} else {
			c.PlayerInfo = nil
		}
	}
	return true
}

func (c *PlayerConnected) Match(raw string) bool {
	m, _ := regexp.Compile(`^\[([0-9.:-]+)]\[([ 0-9]*)]LogSquad: PostLogin: NewPlayer: BP_PlayerController_C .+PersistentLevel\.([^\s]+) \(IP: ([\d.]+) \| Online IDs:([^)|]+)\)`)
	arr := m.FindStringSubmatch(raw)
	if len(arr) < 2 {
		return false
	}
	c.Extra = make(map[string]interface{})
	c.Ip = arr[4]
	c.PlayerController = arr[3]
	c.Time = ToTime(arr[1])
	c.Platform = GetPlatform(arr[5])
	c.ChainID = arr[2]
	JoinRequest[c.ChainID] = &PlayerInfo{
		arr[3], arr[4], ``, c.Platform, 0, 0, 0, ``, 0, 0, true, 0, ``, 0, 0, 0, false, ``, ``, 0,
	}
	if !first {
		go PlayerConnectedFunc(c)
	}
	return true
}

func (c *NewGame) Match(raw string) bool {
	m, _ := regexp.Compile(`^\[([0-9.:-]+)]\[([ 0-9]*)]LogWorld: Bringing World /([A-z]+)/(?:Maps/)?([A-z0-9-]+)/(?:.+/)?([A-z0-9-]+)\.[A-z0-9-]+`)
	arr := m.FindStringSubmatch(raw)
	if len(arr) < 2 {
		return false
	}
	c.Extra = make(map[string]interface{})
	c.Time = ToTime(arr[1])
	c.ChainID = arr[2]
	c.Dlc = arr[3]
	c.MapClassname = arr[4]
	c.LayerClassname = arr[5]
	for _, info := range PlayerInfos {
		info.Kill = 0
		info.Saved = 0
		info.Die = 0
	}
	if c.MapClassname == `Maps` || c.LayerClassname == `TransitionMap` {
		return true
	}
	if !first {
		go NewGameFunc(c)
	} else {
		global.NowGameInfo.Restart()
		r := db(`game_tag`).Where(`log_time`, c.Time/1000).Where(`gt_map`, c.MapClassname).Where(`gt_layer`, c.LayerClassname).Field(`gt_id`).Find()
		res := r.Result()
		if res[`gt_id`] == `` {
			rt := db(`game_tag`).Insert(map[string]string{
				`gt_map`:      c.MapClassname,
				`gt_layer`:    c.LayerClassname,
				`log_time`:    strconv.FormatInt(c.Time/1000, 10),
				`create_time`: strconv.Itoa(date.Date().Unix()),
			})
			global.NowGameInfo.Tag = rt.Id()
		} else {
			tag, _ := strconv.ParseInt(res[`gt_id`], 10, 64)
			global.NowGameInfo.Tag = tag
		}
	}
	return true
}

func (c *DeployableDamaged) Match(raw string) bool {
	m, _ := regexp.Compile(`^\[([0-9.:-]+)]\[([ 0-9]*)]LogSquadTrace: \[DedicatedServer](?:ASQDeployable::)?TakeDamage\(\): ([A-z0-9_]+)_C_[0-9]+: ([0-9.]+) damage attempt by causer ([A-z0-9_]+)_C_[0-9]+ instigator (.+) with damage type ([A-z0-9_]+)_C health remaining ([0-9.]+)`)
	arr := m.FindStringSubmatch(raw)
	if len(arr) < 2 {
		return false
	}
	c.Extra = make(map[string]interface{})
	c.Deployable = arr[3]
	c.Time = ToTime(arr[1])
	c.ChainID = arr[2]
	floatNum, _ := strconv.ParseFloat(arr[4], 64)
	c.Damage = floatNum
	c.Weapon = arr[5]
	c.PlayerSuffix = arr[6]
	c.DamageType = arr[7]
	c.HealthRemaining = arr[8]
	if !first {
		go DeployableDamagedFunc(c)
	}
	return true
}

func (c *AdminBroadcast) Match(raw string) bool {
	m, _ := regexp.Compile(`^\[([0-9.:-]+)]\[([ 0-9]*)]LogSquad: ADMIN COMMAND: Message broadcasted <(.+)> from (.+)`)
	arr := m.FindStringSubmatch(raw)
	if len(arr) < 2 {
		return false
	}
	c.Extra = make(map[string]interface{})
	c.Time = ToTime(arr[1])
	c.Message = arr[3]
	c.From = arr[4]
	if !first {
		go AdminBroadcastFunc(c)
	}
	return true
}
