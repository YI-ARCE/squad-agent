package api

import (
	"squad/module/rcon/global"
	"squad/module/rcon/orm"
	"strconv"
	"yiarce/core/date"
	"yiarce/core/frame"
)

func dataList(data string, aId string) string {
	d := decodeInterface(data)
	query := d[`query`].(map[string]interface{})
	q := orm.Table(`game_bill`)
	for key, value := range query {
		if key == `search` {
			s, _ := value.(string)
			if s != `` {
				q.Where(`locate('` + s + `',gb_msg) > 0`)
			}
			continue
		}
		if key == `gb_type` {
			q.Where(key + ` in (` + value.(string) + `)`)
			continue
		}
	}
	rs := q.Field(`gb_type,gb_msg,log_time`).
		Order(`gb_id`, `desc`).
		Order(`log_time`, `desc`).Page(int(d[`page`].(float64)), int(d[`num`].(float64))).
		Select()
	qc := orm.Table(`game_bill`)
	for key, value := range query {
		qc.Where(key, value)
	}
	rc := qc.Field(`count(*) num`).Find()
	if rs.Err() != nil {
		frame.Println(rs.Err().Error())
	}
	result := map[string]interface{}{
		`num`:  rc.Result()[`num`],
		`list`: rs.Result(),
	}
	return success(result)
}

func dataTag(data string, aId string) string {
	rs := orm.Table(`game_tag`).Field(`gt_id,gt_map,gt_layer,log_time`).Order(`gt_id`, `desc`).Select()
	return success(rs.Result())
}

func getActiveStatus(data string, aId string) string {
	return success(map[string]string{
		`expire`: date.Time(global.ActiveStatusInfo.Expire).Timestamp(`s`),
	})
}

func clearHisData(data string, aId string) string {
	times := date.ParseDate(date.Date().YMD(`-`) + ` 00:00:00`).Unix()
	times -= 86400 * 3
	r := orm.Table(`game_tag`).Where(`create_time < ` + strconv.Itoa(times)).Field(`group_concat(gt_id) ids`).FetchSql().Find()
	if r.Err() != nil {
		return errors(r.Err().Error())
	}
	ids := r.Result()[`ids`]
	if ids != `` {
		orm.Table(`game_tag`).Where(`gt_id in (` + ids + `)`).Delete()
		orm.Table(`game_bill`).Where(`gt_id in (` + ids + `)`).Delete()
		setLog(aId, `清除了三天前的历史对局数据`, `清除数据`, `9`)
	}
	return success(nil)
}
