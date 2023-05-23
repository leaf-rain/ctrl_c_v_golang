package timer

import (
	"fmt"
	"strconv"
	"time"
)

var TimeTemplates = map[string]string{
	"YYYYMM":              "200601",
	"YYYYWW":              "2006",
	"YYYYMMDD":            "20060102",
	"YYYY-MM-DD":          "2006-01-02",
	"YYYY/MM/DD":          "2006/01/02",
	"YYYY#MM#DD":          "2006#01#02",
	"YYYY-MM-DD hh:mm:ss": "2006-01-02 15:04:05",
	"YYYYMMDD hh:mm:ss":   "20060102 15:04:05",
	"DDHHMM":              "021504",
}

func ToYYYYMMDD(t time.Time) string {
	return t.Format(TimeTemplates["YYYYMMDD"])
}

func ToYYYYMM(t time.Time) string {
	return t.Format(TimeTemplates["YYYYMM"])
}

func ToDDHHMM(t time.Time) string {
	return t.Format(TimeTemplates["DDHHMM"])
}

func ToYYYYWW(t time.Time) string {
	return t.Format(TimeTemplates["YYYYWW"]) + strconv.FormatInt(WeekByDate(t), 10)
}

func ToYYYYMMDDhhmmss(t time.Time) string {
	return t.Format(TimeTemplates["YYYYMMDD hh:mm:ss"])
}

func ToYYYY_MM_DDhhmmss(t time.Time) string {
	return t.Format(TimeTemplates["YYYY-MM-DD hh:mm:ss"])
}

func ToYYYYDay(t time.Time) string {
	return fmt.Sprintf("%d%d", t.Year(), t.YearDay())
}

// 判断时间是当年的第几周
func WeekByDate(t time.Time) int64 {
	_, w := t.ISOWeek()
	return int64(w)
}

func GetMaxTimeForDay(t time.Time) int64 {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endUnix := start.AddDate(0, 0, 1).Unix() - 1
	return endUnix
}
