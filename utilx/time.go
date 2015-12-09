package utilx

import (
	"strconv"
	"time"
)

const (
	timeFormatDate     = "2006-01-02"
	timeFormatDateTime = "2006-01-02 15:04:05"
	timeFormatMetaTime = "20060102150405.000"
	timeFormatAtom     = time.RFC3339
)

var (
	TimeZone = time.UTC
)

func TimeFormat(format string) string {

	switch format {
	case "datetime":
		return timeFormatDateTime
	case "date":
		return timeFormatDate
	case "atom":
		return timeFormatAtom
	}

	return timeFormatDateTime
}

func TimeNow(format string) string {
	return time.Now().In(TimeZone).Format(TimeFormat(format))
}

func TimeNowAdd(format, add string) string {
	taf, err := time.ParseDuration(add)
	if err != nil {
		taf = 0
	}
	return time.Now().Add(taf).In(TimeZone).Format(TimeFormat(format))
}

func TimeZoneFormat(t time.Time, tz, format string) string {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return t.Format(TimeFormat(format))
	}
	return t.In(loc).Format(TimeFormat(format))
}

func TimeParse(timeString, format string) time.Time {

	tp, err := time.ParseInLocation(TimeFormat(format), timeString, TimeZone)
	if err != nil {
		return time.Now().In(TimeZone)
	}

	return tp
}

func TimeMetaNow() uint64 {

	t := time.Now().UTC().Format(timeFormatMetaTime)

	if u64, err := strconv.ParseUint(t[:14]+t[15:], 10, 64); err == nil {
		return u64
	}

	return 0
}
