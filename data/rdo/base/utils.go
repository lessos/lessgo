package base

import (
    "time"
)

const (
    timeFormatDate     = "2006-01-02"
    timeFormatDateTime = "2006-01-02 15:04:05"
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
