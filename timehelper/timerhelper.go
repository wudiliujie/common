package timehelper

import (
	"github.com/wudiliujie/common/log"
	"time"
)

func GetNowDateInt() int32 {
	now := time.Now()
	return int32(now.Year()*10000 + int(now.Month())*100 + now.Day())
}
func GetNowTimeInt() int32 {
	now := time.Now()
	return int32(now.Hour()*10000 + int(now.Minute())*100 + now.Second())
}
func GetDateInt(time time.Time) int32 {
	return int32(time.Year()*10000 + int(time.Month())*100 + time.Day())
}
func GetDateIntTimeStamp(timeStamp int64) int32 {
	_time := time.Unix(timeStamp, 0)
	return int32(_time.Year()*10000 + int(_time.Month())*100 + _time.Day())
}
func AddDateInt(date int32, addday int32) int32 {
	t := GetTimeByDateInt(date)
	t = t.Add(time.Hour * 24 * time.Duration(addday))
	return int32(t.Year()*10000 + int(t.Month())*100 + t.Day())
}
func SubDaysInt(date1 int32, data2 int32) int32 {
	t1 := GetTimeByDateInt(date1)
	t2 := GetTimeByDateInt(data2)
	return SubDays(t1, t2)
}
func GetTimeByDateInt(date int32) time.Time {
	year := int(date / 10000)
	month := int(date / 100 % 100)
	day := int(date % 100)
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	return t
}
func StrToTime(str string) time.Time {
	local, _ := time.LoadLocation("Local")
	t, err := time.ParseInLocation("2006-01-02 15:04:05", str, local)
	if err != nil {
		log.Error("StrToTime:%v", err)
		return time.Date(2010, time.January, 1, 0, 0, 0, 0, local)
	}

	return t
}
func TimeToStr(t time.Time) string {
	return t.Format("2006-01-02")
}

//时间戳专业字符串
func TimeStampToStr(timestamp int64) string {
	tm := time.Unix(timestamp, 0)
	return tm.Format("2006-01-02 15:04:05")
}
func SubDays(t1, t2 time.Time) int32 {
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, time.Local)
	t2 = time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, time.Local)
	return int32(t1.Sub(t2).Hours() / 24)
}

//获取当前秒数
func GetNowSceond(timestamp int64) int64 {
	now := time.Unix(timestamp, 0)
	return int64(now.Hour()*3600 + now.Minute()*60 + now.Second())
}

//获取当前秒数
func GetNowDayTimestamp(timestamp int64) int64 {
	return timestamp/86400*86400 - 28800
}

//获取星期一的秒数
func GetMondaySecond(timestamp int64) int64 {
	now := time.Unix(timestamp, 0)
	week := int32(now.Weekday())
	if week == 0 {
		week = 7
	}
	a := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).Add(-time.Duration(week-1) * time.Hour * 24)

	return a.Unix()
}
