package dateparse

import (
	"fmt"
	"strconv"
	"time"
)

const (
	DefaultTime0000 = "0000-00-00 00:00:00"
)

// ParseYmdHisTZ 解析日期 2021-12-12T12:12:12Z
func ParseYmdHisTZ(ymdhis string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05Z", ymdhis)
}

// ParseYmdHisCST 解析日期 2021-12-12T12:12:12+08:00
func ParseYmdHisCST(ymdhis string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05+08:00", ymdhis)
}

// ParseYmdHis 解析日期 2021-12-12 12:12:12
func ParseYmdHis(ymdhis string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", ymdhis)
}

// ParseYmdHis2Ts 解析日期 2021-12-12 12:12:12 转成时间戳
func ParseYmdHis2Ts(ymdhis string) (int64, error) {
	t, err := time.Parse("2006-01-02 15:04:05", ymdhis)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

// ParseYmdHis2TsMs 解析日期 2021-12-12 12:12:12 转成时间戳[毫秒]
func ParseYmdHis2TsMs(ymdhis string) (int64, error) {
	t, err := time.Parse("2006-01-02 15:04:05", ymdhis)
	if err != nil {
		return 0, err
	}
	return t.UnixMilli(), nil
}

// ParseYmd 解析日期 2021-12-12
func ParseYmd(ymd string) (time.Time, error) {
	return time.Parse("2006-01-02", ymd)
}

// TsToDate 时间戳转日期
func TsToDate(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02 15:04:05")
}

// TimestampToDateYmd 时间戳转日期
func TimestampToDateYmd(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02")
}

// TimestampToDateYm 时间戳转日期
func TimestampToDateYm(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01")
}

// TimeFormat 时间格式化
func TimeFormat(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// TsToDateMilli 毫秒级时间戳转日期
func TsToDateMilli(tsMilli int64) string {
	return time.UnixMilli(tsMilli).Format("2006-01-02 15:04:05")
}

// TsToDateYmd 时间戳转日期 年-月-日
func TsToDateYmd(ts time.Time) string {
	return ts.Format("2006-01-02")
}

// TsToDateYm 时间戳转日期 年-月
func TsToDateYm(ts time.Time) string {
	return ts.Format("2006-01")
}

// TsToDateYW 时间戳转日期 年 周
func TsToDateYW(ts time.Time) string {
	year, week := ts.ISOWeek()
	return fmt.Sprintf("%d %d", year, week)
}

// TsToDateYmdH 时间戳转日期
func TsToDateYmdH(ts time.Time) string {
	return ts.Format("2006-01-02 15")
}

// TsToDateYmdHi 时间戳转日期
func TsToDateYmdHi(ts time.Time) string {
	return ts.Format("2006-01-02 15:04")
}

// GetFirstDateOfMonth 获取某个月的第一天0点
func GetFirstDateOfMonth(d time.Time) time.Time {
	return GetZeroTime(d.AddDate(0, 0, -d.Day()+1))
}

// GetFirstDateOfNextMonth 获取下个月第一天0点
func GetFirstDateOfNextMonth(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, 0)
}

// GetZeroTime 获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

// GetZeroTimeOfNextDay 获取下一天0点时间
func GetZeroTimeOfNextDay(d time.Time) time.Time {
	return GetZeroTime(d.AddDate(0, 0, 1))
}

// GetTheDayBeforeZeroTime 获取前一天0点时间
func GetTheDayBeforeZeroTime(d time.Time) time.Time {
	return GetZeroTime(d.AddDate(0, 0, -1))
}

// GetFirstDateOfWeek 获取某周的第一天0点（定义一周的第一天为 周1）
func GetFirstDateOfWeek(d time.Time) time.Time {
	if d.Weekday() == time.Sunday {
		return GetZeroTime(d.AddDate(0, 0, -7).AddDate(0, 0, 1))
	} else {
		return GetZeroTime(d.AddDate(0, 0, -int(d.Weekday())).AddDate(0, 0, 1))
	}
}

// GetFirstDateOfNextWeek 获取下周第一天0点时间（定义一周的第一天为 周1）
func GetFirstDateOfNextWeek(d time.Time) time.Time {
	return GetFirstDateOfWeek(d).AddDate(0, 0, 7)
}

// SecondFriendly 友好化显示秒数
func SecondFriendly(ts int64) string {
	hour := ts / 3600
	min := (ts / 60) % 60
	second := ts % 60

	var str string

	if hour < 10 {
		str = fmt.Sprintf("0%d", hour)
	} else {
		str = fmt.Sprintf("%d", hour)
	}

	if min < 10 {
		str = str + ":" + fmt.Sprintf("0%d", min)
	} else {
		str = str + ":" + fmt.Sprintf("%d", min)
	}

	if second < 10 {
		str = str + ":" + fmt.Sprintf("0%d", second)
	} else {
		str = str + ":" + fmt.Sprintf("%d", second)
	}

	return str
}

func DayFriendly(second int64) string {
	if second <= 0 {
		return ""
	}
	tmp := second / 86400
	if second%86400 > 0 {
		tmp = tmp + 1
	}
	return strconv.FormatInt(tmp, 10)
}
