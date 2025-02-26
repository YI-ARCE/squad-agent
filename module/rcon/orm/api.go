package orm

import (
	"errors"
	"squad/module/rcon/global"
	"squad/module/rcon/log"
	"strconv"
	"yiarce/core/date"
	log2 "yiarce/core/log"
	"yiarce/core/yorm"
)

var db *yorm.Db

func Init(dbs *yorm.Db) {
	db = dbs
}

func GetUid(info *log.PlayerInfo) {
	r := db.Table(`user`).Where(`u_steam`, info.Platform[`steamID`]).Field(`u_id,u_black_info,u_vip_level,u_vip_expire,black_time`).Find()
	data := r.Result()
	if data[`u_id`] == `` {
		id, _ := createUserInfo(info.NickName, info.Platform[`eosID`], info.Platform[`steamID`], info.LoginTime)
		if id > 0 {
			info.UID = id
		}
	} else {
		flag := db.Table(`user_game_chess`).Where(`u_id`, data[`u_id`]).Field(`u_id`).Find().Result()[`u_id`]
		if flag == `` {
			db.Table(`user_game_chess`).Insert(map[string]interface{}{
				`u_id`: data[`u_id`],
			})
		}
		parseInt, _ := strconv.ParseInt(data[`u_id`], 10, 64)
		blackTimeInt, _ := strconv.ParseInt(data[`black_time`], 10, 64)
		info.UID = parseInt
		info.BlackTime = blackTimeInt
		info.BlackInfo = data[`u_black_info`]
		info.Vip, _ = strconv.Atoi(data[`u_vip_level`])
		times, _ := strconv.Atoi(data[`u_vip_expire`])
		if times > 0 {
			if times < date.Date().Unix() {
				info.Vip = 0
			}
		}
	}
}

func CreateUser(info *log.PlayerInfo) error {
	result := db.Table(`user`).Where(`u_steam = ` + info.Platform[`steamID`]).Field(`u_id`).Find()
	if result.Result()[`u_id`] != `` {
		return errors.New(`您已创建账户`)
	} else {
		data := map[string]interface{}{
			`u_name`:    info.NickName,
			`u_eos`:     info.Platform[`eosID`],
			`u_steam`:   info.Platform[`steamID`],
			`last_time`: date.Date().Unix(),
		}
		r := db.Table(`user`).Insert(data)
		if r.Id() < 1 {
			return errors.New(`创建失败,请重试`)
		}
		return nil
	}
}

func createUserInfo(name string, eos string, steam string, loginTime int64) (int64, error) {
	r := db.Table(`user`).FetchSql().Insert(map[string]interface{}{
		`u_name`:      name,
		`u_eos`:       eos,
		`u_steam`:     steam,
		`last_time`:   date.TimeMill(loginTime).Unix(),
		`create_time`: date.Date().Unix(),
	})
	if r.Err() != nil {
		log2.Default(r.Err().Error())
		log2.Default(r.Sql())
		return 0, r.Err()
	}
	id := r.Id()
	db.Table(`user_game_chess`).Insert(map[string]interface{}{
		`u_id`: id,
	})
	return id, nil
}

func BanUser(steamId string, times int64, message string) error {
	r := db.Table(`user`).Where(`u_steam = ` + steamId).Update(map[string]interface{}{
		`black_time`: times,
		`black_info`: message,
	})
	if r.Num() < 1 {
		return errors.New(`操作失败,请重试`)
	}
	return nil
}

// SignIn 签到
func SignIn(platform map[string]string) (int, error) {
	user := db.Table(`user`).Where(`u_steam`, platform[`steamID`]).Field(`u_id,u_name`).Find()
	if user.Err() != nil {
		log2.Default(user.Err().Error())
		return 0, errors.New(`签到失败`)
	}
	userInfo := user.Result()
	if userInfo[`u_id`] == `` {
		return 0, errors.New(`玩家信息未建立,无法签到`)
	}
	r := db.Table(`user_point`).Where(`u_id`, userInfo[`u_id`]).Field(`u_id,u_points,points_time`).Find()
	if r.Err() != nil {
		log2.Default(r.Err().Error())
		return 0, errors.New(`签到失败`)
	}
	result := r.Result()
	d := date.Date()
	point, _ := strconv.Atoi(result[`u_points`])
	if result[`u_id`] == `` {
		r := db.Table(`user_point`).Insert(map[string]interface{}{`u_id`: userInfo[`u_id`], `u_points`: `1`, `points_time`: date.Date().Unix()})
		if r.Err() != nil {
			log2.Default(r.Err().Error())
			return 0, errors.New(`签到失败`)
		}
		point = 1
	} else {
		r := db.Table(`user_point`).Where(`u_id`, userInfo[`u_id`]).Update(map[string]interface{}{
			`u_points`:    yorm.Raw(`u_points + 1`),
			`points_time`: d.Unix(),
		})
		if r.Err() != nil {
			return 0, errors.New(`签到失败`)
		}
		point += 1
	}
	db.Table(`game_bill`).Insert(map[string]interface{}{
		`gb_type`:     8,
		`atk_u_id`:    userInfo[`u_id`],
		`victim_u_id`: 0,
		`gb_msg`:      userInfo[`u_name`] + ` 本局签到成功`,
		`log_time`:    d.Unix(),
		`gt_id`:       global.NowGameInfo.Tag,
		`create_time`: d.Unix(),
	})
	return point, nil
}

// 用户离开服务器,更新他的游玩时长
func SettleGameTime(info *log.PlayerInfo) {
	db.Table(`user`).Where(`u_id`, info.UID).Update(map[string]interface{}{
		`online_time`: yorm.Raw(`online_time + ` + strconv.Itoa(date.Date().Unix()-date.TimeMill(info.LoginTime).Unix())),
	})
}

func Table(name string) yorm.ModelTransfer {
	return db.Table(name)
}

func UpdatePoint(steamId string, data map[string]string) {

}
