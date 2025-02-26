package rcon

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"regexp"
	"squad/module/rcon/global"
	"squad/module/rcon/log"
	"squad/module/rcon/orm"
	"squad/module/rcon/server"
	"strconv"
	"strings"
	"time"
	"yiarce/core/date"
	"yiarce/core/frame"
)

type Request struct {
	RequestID uint16
	Type      uint32
	Count     uint16
	Payload   []byte
}

type Buffers struct {
	incomingData *bytes.Buffer
}

type Response struct {
	RequestID int
	Count     int
	Type      int
	Payload   []byte
}

type ChatMessage struct {
	Type     string            `json:"type"`
	Platform map[string]string `json:"platform"`
	UserName string            `json:"nick_name"`
	Message  string            `json:"message"`
	Time     int64             `json:"time"`
}

var reg = map[string]*regexp.Regexp{}

func init() {
	squads, _ := regexp.Compile(`ID: (\d+) \| Name: (.+) \| Size: (\d+) \| Locked: (True|False) \| Creator Name: (.+) \| Creator Online IDs:([^|]+)`)
	reg[`squads`] = squads
}

func (r Request) Serialize() []byte {
	data := make([]byte, len(r.Payload)+14)
	binary.LittleEndian.PutUint32(data[0:4], uint32(len(r.Payload)+10))
	binary.LittleEndian.PutUint16(data[4:6], r.RequestID)
	binary.LittleEndian.PutUint16(data[6:8], r.Count)
	binary.LittleEndian.PutUint32(data[8:12], r.Type)
	copy(data[12:], r.Payload)
	data[len(data)-2] = 0x00 // padding
	data[len(data)-1] = 0x00 // padding
	return data
}

func (r *Response) Deserialize(data []byte) error {
	l := len(data)
	if l < 9 {
		return fmt.Errorf("data too short")
	}
	r.RequestID = int(binary.LittleEndian.Uint16(data[4:6]))
	r.Count = int(binary.LittleEndian.Uint16(data[6:8]))
	r.Type = int(binary.LittleEndian.Uint32(data[8:12]))
	r.Payload = data[12 : l-2]
	return nil
}

type EventBase interface {
	Match(str string) bool
}

type EventBaseModule struct {
	eventBase []EventBase
}

func (c *EventBaseModule) Match(str string) bool {
	for _, base := range c.eventBase {
		if base.Match(str) {
			return true
		}
	}
	return false
}

type EventAdminCamera struct {
	Name string `json:"name"`
	Time int64  `json:"time"`
}

type EventUnAdminCamera struct {
	Name string `json:"name"`
	Time int64  `json:"time"`
}

type EventAdminWarn struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
	Time   int64  `json:"time"`
}

type EventKick struct {
	Player string `json:"player"`
	Reason string `json:"reason"`
	Time   int64  `json:"time"`
}

type EventSquadCreated struct {
	SquadID   string `json:"squad_id"`
	SquadName string `json:"squad_name"`
	TeamName  string `json:"team_name"`
	Time      int64  `json:"time"`
}

type EventBanPlayer struct {
	PlayerID string `json:"player_id"`
	Name     string `json:"name"`
	Interval string `json:"interval"`
	Time     int64  `json:"time"`
}

type EventCurrentMap struct {
	Level string `json:"level"`
	Layer string `json:"layer"`
	Time  int64  `json:"time"`
}

type EventNextMap struct {
	Level string `json:"level"`
	Layer string `json:"layer"`
	Time  int64  `json:"time"`
}

type EventListPlayers struct {
}

type EventListSquads struct {
	Team []string
}

func (c EventListSquads) Match(raw string) bool {
	flag := false
	if strings.Contains(raw, `----- Active Squads -----`) {
		flag = true
		if len(raw) == 25 {
			frame.Println(`跳过了`)
			return true
		}
		raw = raw[26:]
		arr := strings.Split(raw, "Team ID:")[1:]
		for _, s := range arr {
			index := strings.Index(s, "\n")
			if index == -1 {
				index = len(s)
				s += "\n"
			}
			id, _ := strconv.Atoi(s[1:2])
			if id == 0 {
				continue
			}
			teamName := s[4 : index-1]
			s = s[index+1:]
			arrSquad := strings.Split(s, "\n")
			global.NowGameInfo.GameTeamInfo[id-1].Name = teamName
			squad := map[string]global.SquadInfo{}
			squad[`0`] = global.SquadInfo{
				Name:     `无队伍`,
				Size:     0,
				Lock:     false,
				UserList: make([]map[string]interface{}, 0),
				Time:     0,
			}
			for _, s := range arrSquad {
				if s == `` {
					break
				}
				match := reg[`squads`].FindStringSubmatch(s)
				size, _ := strconv.Atoi(match[3])
				lock := match[4] == `True`
				squad[match[1]] = global.SquadInfo{
					Id:       match[1],
					Name:     match[2],
					Size:     size,
					Lock:     lock,
					UserList: make([]map[string]interface{}, 0),
					Time:     global.NowGameInfo.GameTeamSquad[id-1][match[1]],
				}
			}
			global.NowGameInfo.GameTeamInfo[id-1].Squad = squad
		}
		server.PushAll(server.TypeTeamInfo, getPanelSquadList())
	}
	return flag
}

func (c EventListPlayers) Match(raw string) bool {
	flag := false
	strs := strings.Split(raw, "\n")
	m, _ := regexp.Compile(`^ID: (\d+) \| Online IDs:([^|]+)\| Name: (.+) \| Team ID: (\d|N/A) \| Squad ID: (\d+|N/A) \| Is Leader: (True|False) \| Role: (.+)$`)
	for _, str := range strs {
		if strings.Contains(str, `----- Active Players -----`) {
			flag = true
			continue
		}
		arr := m.FindStringSubmatch(str)
		if len(arr) > 1 {
			flag = true
			teamID, _ := strconv.ParseInt(arr[4], 10, 64)
			squadID, _ := strconv.ParseInt(arr[5], 10, 64)
			p := log.GetPlatform(arr[2])
			if p[`steamID`] != `` && log.PlayerInfos[p[`eosID`]] != nil {
				log.PlayerInfos[p[`eosID`]].TeamID = int(teamID)
				log.PlayerInfos[p[`eosID`]].SquadID = int(squadID)
				log.PlayerInfos[p[`eosID`]].Leader = arr[6] == `True`
				log.PlayerInfos[p[`eosID`]].Role = arr[7]
			}
		}
	}
	if flag {
		server.PushAll(server.TypeRconResponse, getUserList(), "user_list")
		return true
	}
	return flag
}

func (c EventNextMap) Match(raw string) bool {
	m, _ := regexp.Compile(`^Next level is (.*), layer is (.*)`)
	arr := m.FindStringSubmatch(raw)
	if len(arr) < 2 {
		return false
	}
	c.Level = arr[1]
	c.Layer = arr[2]
	c.Time = date.Date().Source().UnixMilli()
	server.PushAll(server.TypeRconResponse, c, "rcon_next_map")
	return true
}

func (c EventCurrentMap) Match(raw string) bool {
	m, _ := regexp.Compile(`^Current level is (.*), layer is (.*)`)
	arr := m.FindStringSubmatch(raw)
	if len(arr) < 2 {
		return false
	}
	c.Level = arr[1]
	c.Layer = arr[2]
	c.Time = date.Date().Source().UnixMilli()
	if global.NowGameInfo.First {
		//teams := strings.Split(strings.Split(c.Layer, `,`)[1], ` `)
		global.NowGameInfo.First = false
		return true
	}
	server.PushAll(server.TypeRconResponse, c, "rcon_current_map")
	go func() {
		time.Sleep(time.Second / 2)
		server.PushAll(server.TypeApiResponse, map[string]bool{
			`flag`: global.NowGameInfo.RandomPlayer,
		}, `randomTeamPlayer`)
	}()
	return true
}

func (c EventBanPlayer) Match(raw string) bool {
	m, _ := regexp.Compile(`Banned player ([0-9]+)\. \[Online IDs=([^\]]+)\] (.*) for interval (.*)`)
	arr := m.FindStringSubmatch(raw)
	if len(arr) < 2 {
		return false
	}
	c.PlayerID = arr[1]
	c.Name = arr[3]
	c.Interval = arr[4]
	c.Time = date.Date().Source().UnixMilli()
	server.PushAll(server.TypeRconResponse, c, "rcon_ban_player")
	return true
}

func (c EventSquadCreated) Match(raw string) bool {
	m, _ := regexp.Compile(`(.+) \(Online IDs:([^)]+)\) has created Squad (\d+) \(Squad Name: (.+)\) on (.+)`)
	arr := m.FindStringSubmatch(raw)
	if len(arr) < 2 {
		return false
	}
	c.SquadID = arr[3]
	c.SquadName = arr[4]
	c.TeamName = arr[5]
	c.Time = date.Date().Source().UnixMilli()
	server.PushAll(server.TypeRconResponse, c, "rcon_squad_created")
	return true
}

func (c EventKick) Match(raw string) bool {
	m, _ := regexp.Compile(`Kicked player ([0-9]+)\. \[Online IDs=([^\]]+)\] (.*)`)
	arr := m.FindStringSubmatch(raw)
	if len(arr) < 2 {
		return false
	}
	c.Player = arr[3]
	c.Reason = arr[1]
	c.Time = date.Date().Source().UnixMilli()
	server.PushAll(server.TypeRconResponse, c, "rcon_kick_player")
	return true
}

func (c EventAdminWarn) Match(raw string) bool {
	m, _ := regexp.Compile(`Remote admin has warned player (.*)\. Message was "(.*)"`)
	arr := m.FindStringSubmatch(raw)
	if len(arr) < 2 {
		return false
	}
	c.Name = arr[1]
	c.Reason = arr[2]
	c.Time = date.Date().Source().UnixMilli()
	server.PushAll(server.TypeRconResponse, c, "rcon_admin_warn")
	return true
}

func (c EventUnAdminCamera) Match(raw string) bool {
	m, _ := regexp.Compile(`\[Online IDs:([^\]]+)\] (.+) has unpossessed admin camera\.`)
	arr := m.FindStringSubmatch(raw)
	if len(arr) < 2 {
		return false
	}
	c.Name = arr[2]
	c.Time = date.Date().Source().UnixMilli()
	server.PushAll(server.TypeRconResponse, c, "rcon_un_admin_camera")
	orm.Table(`game_bill`).Insert(map[string]interface{}{
		`gb_type`:     6,
		`atk_u_id`:    0,
		`victim_u_id`: 0,
		`gt_id`:       global.NowGameInfo.Tag,
		`gb_msg`:      c.Name + ` 退出上帝视角`,
		`log_time`:    c.Time / 1000,
		`create_time`: date.Date().Unix(),
	})
	return true
}

func (c EventAdminCamera) Match(raw string) bool {
	m, _ := regexp.Compile(`\[Online Ids:([^\]]+)\] (.+) has possessed admin camera\.`)
	arr := m.FindStringSubmatch(raw)
	if len(arr) < 2 {
		return false
	}
	c.Name = arr[2]
	c.Time = date.Date().Source().UnixMilli()
	orm.Table(`game_bill`).Insert(map[string]interface{}{
		`gb_type`:     6,
		`atk_u_id`:    0,
		`victim_u_id`: 0,
		`gt_id`:       global.NowGameInfo.Tag,
		`gb_msg`:      c.Name + ` 进入了上帝视角`,
		`log_time`:    c.Time / 1000,
		`create_time`: date.Date().Unix(),
	})
	server.PushAll(server.TypeRconResponse, c, "rcon_admin_camera")
	return true
}
