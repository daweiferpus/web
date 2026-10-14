package public

import (
	"time"
)

type TradeDateStamp struct {
	OneDayAgo   int64
	ThreeDayAgo int64
	AWeekAgo    int64
	AMonthAgo   int64
}

func GetTradeTimeStamp(date time.Time) int64 {
	return time.Date(date.Year(), date.Month(), date.Day(), 15, 0, 0, 0, date.Location()).UnixNano() / 1e6
}

func GetZeroClockTimeStamp(date time.Time) int64 {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location()).UnixNano() / 1e6
}

func GetNowZeroClockTimeStamp() time.Time {
	t := time.Now()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func GetZeroClockTimeStr(format string) string {
	t := time.Now()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Format(format)
}

//计算基于当前时间的上一日 上三日 上一周 上一个月的时间挫
func CalcDate() *TradeDateStamp {

	date := time.Now()
	weekday := int(date.Weekday())
	var diff int
	//按交易日计算的逻辑需要完善 需要引入交易日历
	if weekday > 1 && weekday <= 6 {
		diff = 1
	} else if weekday == 0 { //周日
		diff = 2
	} else if weekday == 1 {
		diff = 3
	}

	date.AddDate(0, 0, -1)
	trdDateStp := TradeDateStamp{}
	trdDateStp.OneDayAgo = GetTradeTimeStamp(date.AddDate(0, 0, -diff))
	trdDateStp.ThreeDayAgo = GetTradeTimeStamp(date.AddDate(0, 0, -(diff + 2)))
	trdDateStp.AWeekAgo = GetTradeTimeStamp(date.AddDate(0, 0, -7))
	trdDateStp.AMonthAgo = GetTradeTimeStamp(date.AddDate(0, -1, 0))

	return &trdDateStp
}
