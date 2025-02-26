package rcon

import (
	"encoding/binary"
	"io"
	"net"
	"os"
	"regexp"
	command2 "squad/module/rcon/command"
	"squad/module/rcon/global"
	"squad/module/rcon/log"
	"squad/module/rcon/orm"
	"squad/module/rcon/server"
	"strings"
	"sync"
	"time"
	"unsafe"
	"yiarce/core/date"
	"yiarce/core/frame"
	"yiarce/core/yorm"
	_ "yiarce/core/yorm/mysql"
)

var conn net.Conn

var buff []byte
var responseQueue []Response

var chatReg = regexp.MustCompile(`\[(ChatAll|ChatTeam|ChatSquad|ChatAdmin)] \[Online IDs:([^]]+)] (.+?) : (.*)`)
var rconConfig = map[string]string{}
var count = uint16(10)
var mutex = &sync.Mutex{}
var reload = false
var heartFirst = false
var sendFlag = false

var refreshTime int

func Link(host string, password string) {
	frame.Println(`数据连接初始化中...`)
	global.SetFunc(`getUserList`, getUserList)
	config := yorm.Config{Host: `127.0.0.1`, Username: "root", Password: password, Database: global.ServerConfig[`database`], Port: `3306`}
	mysql, err := yorm.ConnMysql(config)
	if err != nil {
		if strings.Contains(err.Error(), `Unknown database '`+global.ServerConfig[`database`]+`'`) {
			mysql = initDataBase(config)
		} else {
			frame.Println(`数据库启动异常:`, err.Error())
			return
		}
	}
	orm.Init(mysql)
	rconConfig[`host`] = host
	rconConfig[`password`] = password
	linkRcon()
	command2.Init(Send)
}

func connect() {
	count = 10
	buff = make([]byte, 0)
	c, err := net.Dial("tcp", rconConfig[`host`])
	if err != nil {
		frame.Println("RCON连接失败", err.Error())
		go func() {
			time.Sleep(time.Second * 5)
			connect()
		}()
		return
	}
	frame.Println("RCON连接成功")
	conn = c
	if !heartFirst {
		go keepAliveRcon()
		heartFirst = true
	}
	go onData()
	// 认证
	if err = authenticate(rconConfig[`password`]); err != nil {
		frame.Println(err.Error())
		return
	}
}

func linkRcon() {
	frame.Println(`数据库连接成功,启动客户端服务进程...`)
	go server.Listen(func() {
		go log.Monitor(orm.Table, orm.GetUid, logRegisterSuccess)
	}, commandClient)
	frame.Println(`连接服务器RCON`)
	connect()
}

func logRegisterSuccess() {
	sendFlag = true
}

func keepAliveRcon() {
	for {
		time.Sleep(time.Second * 60)
		Send(2, 2, `ping heart`)
	}
}

func register() {
	log.PlayerConnectedFunc = OnPlayerConnected
	log.NewGameFunc = OnNewGame
	log.PlayerDisconnectedFunc = OnPlayerDisconnected
	log.PlayerJoinSucceededFunc = OnPlayerJoinSucceeded
	log.PlayerDiedFunc = OnPlayerDied
	log.PlayerRevivedFunc = OnPlayerRevived
	log.PlayerWoundedFunc = OnPlayerWounded
	log.PlayerUnPossessFunc = OnPlayerUnPossess
	log.PlayerPossessFunc = OnPlayerPossess
	log.ServerTickRateFunc = OnServerTickRate
	log.RoundEndedFunc = OnRoundEnded
	log.RoundWinnerFunc = OnRoundWinner
	log.DeployableDamagedFunc = OnDeployableDamaged
	log.AdminBroadcastFunc = OnAdminBroadcast
	log.PlayerDamagedFunc = OnPlayerDamaged
	log.RoundTicketsFunc = OnRoundTickets
	log.PlayerSquadCreateFunc = OnSquadCreate
}

func response(r Response) {
	switch r.Type {
	case 0:
		switch r.RequestID {
		case 1:
		case 2:
			byteData := r.Payload
			str := string(byteData)
			if len(str) > 0 {
				baseCase := EventBaseCase()
				if !baseCase.Match(str) {
					//fmt.Println(`未知的RCON消息:`, str)
				}
			}
		}
	case 1:
		str := string(r.Payload)
		// 使用正则表达式解析日志条目
		matches := chatReg.FindStringSubmatch(str)
		// 检查是否找到了匹配项
		if len(matches) > 0 {
			// 输出匹配结果
			go chatEvent(ChatMessage{
				Type:     matches[1],
				Platform: log.GetPlatform(matches[2]),
				UserName: strings.Trim(matches[3], ` `),
				Message:  matches[4],
				Time:     time.Now().UnixMilli(),
			})
		} else {
			baseCase := EventBaseCase2()
			if !baseCase.Match(str) {
				//fmt.Println(`未知的RCON消息:`, str)
			}
		}
	case 2:
		if r.RequestID != 2 {
			frame.Println(`RCON密码不正确`)
			os.Exit(0)
		} else {
			if reload {
				go func() {
					time.Sleep(time.Second / 2)
					reload = false
				}()
			} else {
				frame.Println(`注册日志监听事件`)
				register()
				frame.Println(`日志监听事件注册完成`)
				if global.ActiveStatusInfo.Flag {
					go log.Monitor(orm.Table, orm.GetUid, logRegisterSuccess)
					time.Sleep(time.Second)
				}
			}

		}
	default:
		frame.Println(`未知包类型`)
		frame.Println(`包类型:`, r.Type, `包ID`, r.RequestID, string(r.Payload))
	}
}

func decodeData(data []byte) {
	buff = append(buff, data...)
	for len(buff) >= 4 {
		if len(buff) == 7 {
			buff = []byte{}
			l := len(responseQueue)
			for i := 0; i < l; i++ {
				response(responseQueue[i])
			}
			responseQueue = []Response{}
			break
		}
		size := int(binary.LittleEndian.Uint32(buff[:4]))
		packetSize := size + 4
		buffLen := len(buff)
		if buffLen < packetSize {
			if packetSize > 10000 {
				if !reload {
					reload = true
					// 出现异常了,终止连接后重连
					frame.Println(`接收数据包异常,初始化rcon`)
					conn.Close()
					conn = nil
					go connect()
					return
				}
				return
			}
			//frame.Println(`包接收不全,等待更多包接收,当前包大小:`, buffLen, `所需包大小:`, packetSize)
			break
		}
		r := Response{}
		err := r.Deserialize(buff[:packetSize])
		if err != nil {
			frame.Println(`解包失败:`, err.Error())
			buff = buff[packetSize:]
			continue
		}
		buff = buff[packetSize:]
		if r.Type == 2 || r.Count == 0 {
			response(r)
		} else {
			if len(r.Payload) > 4 {
				l := len(responseQueue)
				flag := true
				for i := 0; i < l; i++ {
					if responseQueue[i].Count == r.Count {
						responseQueue[i].Payload = append(responseQueue[i].Payload, r.Payload...)
						flag = false
						break
					}
				}
				if flag {
					responseQueue = append(responseQueue, r)
				}
				//frame.Println(`包响应:`, string(r.Payload))
				//frame.Println(`包数据:`, r.Payload[len(r.Payload)-4:])
			}
		}
	}
}

func onData() {
	for {
		if conn == nil {
			return
		}
		buffer := make([]byte, 4096)
		n, err := conn.Read(buffer)
		if err != nil {
			if !reload {
				if err != io.EOF {
					frame.Println(err.Error())
				} else {
					continue
				}
			}
			return
		}
		decodeData(buffer[:n])
	}
}

func Send2(id uint16, types uint32, commend string, flag ...bool) error {
	if !sendFlag {
		return nil
	}
	if conn == nil {
		return nil
	}
	if reload && types != 3 {
		return nil
	} else {
		if global.NowGameInfo.RandomLock && len(flag) < 1 {
			return nil
		}
	}
	req := buildRequest(id, types, commend, count)
	mutex.Lock()
	_, err := conn.Write(req.Serialize())
	mutex.Unlock()
	if err != nil {
		return err
	}
	count++
	if count >= 65500 {
		count = 10
	}
	mutex.Lock()
	conn.Write(buildRequest(2, 2, ``, count).Serialize())
	mutex.Unlock()
	if date.Date().Unix()-refreshTime > 1 {
		if strings.Contains(commend, `AdminDisbandSquad`) || strings.Contains(commend, `AdminKick`) ||
			strings.Contains(commend, `AdminForceTeamChange`) || strings.Contains(commend, `AdminRemovePlayerFromSquad`) {
			refreshTime = date.Date().Unix()
			go func() {
				time.Sleep(time.Second / 10)
				Send(2, 2, `ListPlayers`)
				time.Sleep(time.Second / 30)
				Send(2, 2, `ListSquads`)
			}()
		}
	}
	return nil
}

func Send(id uint16, types uint32, commend string, flag ...bool) error {
	if !sendFlag && types != 3 {
		return nil
	}
	if conn == nil {
		return nil
	}
	if reload && types != 3 {
		return nil
	} else {
		if global.NowGameInfo.RandomLock && len(flag) < 1 {
			return nil
		}
	}
	req := buildRequest(id, types, commend, count)
	mutex.Lock()
	_, err := conn.Write(req.Serialize())
	mutex.Unlock()
	if err != nil {
		return err
	}
	if commend != `ping heart` && types == 2 {
		//frame.Println(`发送命令: `, commend)
		count++
		if count >= 65500 {
			count = 10
		}
		mutex.Lock()
		conn.Write(buildRequest(2, 2, ``, count).Serialize())
		mutex.Unlock()
		if len(flag) < 1 {
			if date.Date().Unix()-refreshTime > 1 {
				if strings.Contains(commend, `AdminDisbandSquad`) || strings.Contains(commend, `AdminKick`) ||
					strings.Contains(commend, `AdminForceTeamChange`) || strings.Contains(commend, `AdminRemovePlayerFromSquad`) {
					refreshTime = date.Date().Unix()
					go func() {
						time.Sleep(time.Second / 10)
						Send(2, 2, `ListPlayers`)
						time.Sleep(time.Second / 30)
						Send(2, 2, `ListSquads`)
					}()
				}
			}
		}

	}
	return nil
}

func buildRequest(id uint16, types uint32, body string, count uint16) Request {
	c := Request{}
	c.RequestID = id
	c.Type = types
	c.Count = count
	c.Payload = *(*[]byte)(unsafe.Pointer(&body))
	return c
}

func authenticate(password string) error {
	err := Send(2, 3, password)
	if err != nil {
		return err
	}
	return nil
}
