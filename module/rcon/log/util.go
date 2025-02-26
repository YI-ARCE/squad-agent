package log

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

var logEvent = []string{`LogSquadTrace`, `LogSquadOnlineServices`, `LogSquad`, `LogWorld`, `LogGame`, `LogNet: Join succeeded`, `LogNet: UChannel::Close`}

// GetPlatform 获取用户平台信息
func GetPlatform(raw string) map[string]string {
	m, _ := regexp.Compile(`(\w+):\s*(\w+)`)
	arr := m.FindAllStringSubmatch(raw, -1)
	result := map[string]string{}
	for _, strs := range arr {
		result[strings.ToLower(strs[1])+`ID`] = strs[2]
	}
	return result
}

func ToTime(raw string) int64 {
	arr := strings.Split(raw, `:`)
	parse, err := time.Parse(`2006.01.02-15.04.05`, arr[0])
	if err != nil {
		panic(err.Error())
	}
	ms, err := strconv.Atoi(arr[1])
	if err != nil {
		panic(err.Error())
	}
	return parse.UnixMilli() + int64(ms)
}
func switchEventBase(eventName string) EventBaseModule {
	switch eventName {
	case `LogSquadTrace`:
		return newLogSquadTraceEventBase()
	case `LogSquad`:
		return newLogSquadEventBase()
	case `LogWorld`:
		return newLogWorldEventBase()
	case `LogGame`:
		return newLogGameEventBase()
	case `LogSquadOnlineServices`:
		return EventBaseModule{}
	case `LogNet: Join succeeded`:
		return newJoinPlayerEventBase()
	case `LogNet: UChannel::Close`:
		return newDisconnectEventBase()
	default:
		return EventBaseModule{}
	}
}
func newLogSquadTraceEventBase() EventBaseModule {
	return EventBaseModule{[]EventBase{
		&DeployableDamaged{},
		&PlayerDied{},
		&PlayerPossess{},
		&PlayerUnPossess{},
		&PlayerWounded{},
	}}
}

func newLogSquadEventBase() EventBaseModule {
	return EventBaseModule{[]EventBase{
		&AdminBroadcast{},
		&PlayerConnected{},
		&PlayerRevived{},
		&ServerTickRate{},
		&PlayerWounded{},
		&RoundTickets{},
		&PlayerDamaged{},
		&PlayerCreateSquad{},
		&GameLoadTeam{},
	}}
}

func newLogWorldEventBase() EventBaseModule {
	return EventBaseModule{[]EventBase{
		&NewGame{},
	}}
}

func newLogGameEventBase() EventBaseModule {
	return EventBaseModule{[]EventBase{
		&RoundEnded{},
	}}
}

func newJoinPlayerEventBase() EventBaseModule {
	return EventBaseModule{[]EventBase{
		&PlayerJoinSucceeded{},
	}}
}

func newDisconnectEventBase() EventBaseModule {
	return EventBaseModule{[]EventBase{
		&PlayerDisconnected{},
	}}
}
