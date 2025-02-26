package api

import (
	"encoding/json"
	"math/rand"
	"squad/module/rcon/command"
	"squad/module/rcon/global"
	log2 "squad/module/rcon/log"
	"squad/module/rcon/orm"
	"time"
	"yiarce/core/date"
	"yiarce/core/frame"
)

func decode(str string) map[string]string {
	m := make(map[string]string)
	json.Unmarshal([]byte(str), &m)
	return m
}

func decodeInterface(str string) map[string]interface{} {
	m := make(map[string]interface{})
	json.Unmarshal([]byte(str), &m)
	return m
}

func encode(data interface{}) string {
	d, _ := json.Marshal(data)
	return string(d)
}

func errors(str string) string {
	return `{"code":0, "msg":"` + str + `"}`
}

func RandomTeamPlayerNow() {
	global.NowGameInfo.RandomLock = true
	var team = make([][]string, 2)
	for _, info := range log2.PlayerInfos {
		if info.TeamID == 0 {
			continue
		}
		team[info.TeamID-1] = append(team[info.TeamID-1], info.Platform[`steamID`])
	}
	// 需要跳边的玩家
	var pushTeam = make([][]string, 2)
	var teams = make([]int, 2)
	teams[0] = len(team[0])
	teams[1] = len(team[1])
	//// 需要的配平数
	poise := (teams[0] + teams[1]) / 2
	// 统计已分配的人数
	radioNum := [2]int{teams[0] / 2, teams[1] / 2}
	// 持平两边总人数
	if teams[0] != teams[1] {
		if teams[0] > teams[1] {
			radioNum[0] += 1
		} else {
			radioNum[1] += 1
		}
	}
	pushNum := [2]int{0, 0}
	for teamID, ss := range team {
		for userIndex, s := range ss {
			// 已输出人数等于该团队一半时结束循环
			if pushNum[teamID] == radioNum[teamID] {
				// 满足直接跳出
				break
			}
			if teams[teamID]-(userIndex) > (radioNum[teamID] - pushNum[teamID]) {
				randomNum := rand.Intn(100)
				if randomNum > 50 {
					pushTeam[teamID] = append(pushTeam[teamID], s)
					pushNum[teamID]++
				}
			} else {
				pushTeam[teamID] = append(pushTeam[teamID], s)
				pushNum[teamID]++
			}
		}
	}
	flag := true
	num := 0
	t1Index := 0
	t2Index := 0
	if poise > 1 {
		for {
			if t1Index == pushNum[0] && t2Index == pushNum[1] {
				break
			}
			if flag {
				if t1Index < pushNum[0] {
					command.AdminForceTeamChange(pushTeam[0][t1Index], true)
					t1Index++
					//frame.Println(`团队1->团队2`, pushTeam[0][t1Index], `团队1当前进度:`, t1Index)
					num++
				}
			} else {
				if t2Index < pushNum[1] {
					command.AdminForceTeamChange(pushTeam[1][t2Index], true)
					t2Index++
					//frame.Println(`团队2->团队1`, pushTeam[1][t2Index], `团队2当前进度:`, t2Index)
					num++
				}
			}
			time.Sleep(time.Second / 50)
			flag = !flag
		}
	}
	global.NowGameInfo.RandomLock = false
	frame.Println(`配平已执行`)
	time.Sleep(time.Second / 3)
	command.ListPlayers()
	time.Sleep(time.Second / 3)
	command.ListSquads()
}

func setLog(aId string, content string, tag string, types string) {
	orm.Table(`admin_log`).Insert(map[string]string{
		`al_type`:     types,
		`al_content`:  content,
		`a_id`:        aId,
		`al_tag`:      tag,
		`create_time`: date.Date().Timestamp(`s`),
	})
}
