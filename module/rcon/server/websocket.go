package server

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/http"
	"os"
	"squad/module/rcon/global"
	"squad/module/rcon/orm"
	"squad/module/rcon/server/api"
	"strconv"
	"strings"
	"sync"
	"time"
	"yiarce/core/date"
	"yiarce/core/frame"
)

var activeFunc func()

var commandFunc func(str string)

func Listen(f func(), c func(str string)) {
	activeFunc = f
	commandFunc = c
	http.HandleFunc("/", socket)
	err := http.ListenAndServe("0.0.0.0:"+global.ServerConfig[`listing_port`], nil)
	if err != nil {
		frame.Println(`服务端启动失败: 转发通道异常`)
		time.Sleep(5)
		os.Exit(0)
	}
}

var conn sync.Map

var mutex = sync.Mutex{}

var uid = map[string]string{}

func socket(w http.ResponseWriter, r *http.Request) {
	wu := websocket.Upgrader{}
	wu.ReadBufferSize = 1024
	wu.WriteBufferSize = 1024
	wu.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	wu.Error = func(w http.ResponseWriter, r *http.Request, status int, reason error) {
		frame.Println(reason.Error())
	}
	c, err := wu.Upgrade(w, r, nil)
	if err != nil {
		frame.Println(err.Error())
		return
	}
	frame.Println(r.RemoteAddr[0:strings.IndexAny(r.RemoteAddr, ":")], "连接成功")
	id := strconv.Itoa(date.Date().Unix()) + hex.EncodeToString([]byte(r.RemoteAddr)) + hex.EncodeToString([]byte(strconv.FormatInt(rand.Int63n(99999999), 10)))
	conn.Store(id, c)
	defer func(c *websocket.Conn) {
		err := c.Close()
		if err != nil {
			frame.Println(err.Error())
		}
	}(c)
	for {
		t, p, err := c.ReadMessage()
		if err != nil {
			c = nil
			conn.Delete(id)
			delete(uid, id)
			break
		}
		switch t {
		case -1:
			conn.Delete(id)
			delete(uid, id)
		case 2:
			break
		case 1:
			onMessage(c, p, id)
		default:
			fmt.Println("类型:", t, "消息:", p)
			break
		}
	}
}
func onMessage(connect *websocket.Conn, byteRaw []byte, clientID string) {
	data := map[string]string{}
	err := json.Unmarshal(byteRaw, &data)
	if err != nil {
		frame.Println(err.Error())
		return
	}
	switch data[`type`] {
	case TypeLoginRequest:
		info := LoginRequest{}
		json.Unmarshal([]byte(data[`data`]), &info)
		lr := checkAdmin(info)
		if lr.LoginStatus {
			uid[clientID] = lr.AId
		}
		Push(connect, clientID, TypeLoginResponse, lr)
	case TypeHeartbeatPackage:
		break
	case TypeRconCommand:
		commandFunc(data[`data`])
		orm.Table(`admin_log`).Insert(map[string]string{
			`al_type`:     `2`,
			`al_content`:  data[`data`],
			`a_id`:        uid[clientID],
			`create_time`: date.Date().Timestamp(`s`),
		})
	case TypeGetUserList:
		Push(connect, clientID, TypePlayerInfo, global.ActiveFunc(`getUserList`).([]map[string]interface{}))
	case TypeGetGameInfo:
		Push(connect, clientID, TypeGameInfo, global.NowGameInfo.GameTeamInfo)
	case TypeApiRequest:
		Push(connect, clientID, TypeApiResponse, api.Actions[data[`action`]](data[`data`], uid[clientID]), data[`action`])
		orm.Table(`admin_log`).Insert(map[string]string{
			`al_type`:     `2`,
			`al_content`:  data[`data`],
			`a_id`:        uid[clientID],
			`create_time`: date.Date().Timestamp(`s`),
		})
	}
}

func Push(ctx *websocket.Conn, clientID string, types string, data interface{}, tag ...string) {
	res, errs := json.Marshal(data)
	if errs != nil {
		frame.Println(errs.Error())
	}
	var builder bytes.Buffer
	builder.WriteString(`{"type":"` + types + `","data":`)
	builder.Write(res)
	if len(tag) > 0 {
		builder.WriteString(`,"tag":"` + tag[0] + `"`)
	}
	builder.WriteString(`}`)
	err := ctx.WriteMessage(1, builder.Bytes())
	if err != nil {
		conn.Delete(clientID)
	}
}

func PushAll(types string, data interface{}, tag ...string) {
	res, errs := json.Marshal(data)
	if errs != nil {
		frame.Println(errs.Error())
	}
	var builder bytes.Buffer
	builder.WriteString(`{"type":"` + types + `","data":`)
	builder.Write(res)
	if len(tag) > 0 {
		builder.WriteString(`,"tag":"` + tag[0] + `"`)
	}
	builder.WriteString(`}`)
	byteRaw := builder.Bytes()
	conn.Range(func(key, value interface{}) bool {
		v := value.(*websocket.Conn)
		if v != nil {
			mutex.Lock()
			err := v.WriteMessage(1, byteRaw)
			if err != nil {
				frame.Println(err.Error())
				conn.Delete(key.(string))
			}
		}
		mutex.Unlock()
		return true
	})
}
