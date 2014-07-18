package logger

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	locker   sync.Mutex
	levels   = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
	levelMap = map[string]int{}
	bufs     = make(chan *entry, 100000)
)

const (
	printDefault uint8 = iota
	printFormat
)

type entry struct {
	ptype      uint8
	level      string
	format     string
	fileName   string
	lineNumber int
	ltime      time.Time
	args       []interface{}
}

func LevelConfig(ls []string) {

	if len(ls) < 1 {
		return
	}

	levels = []string{}
	for _, v := range ls {
		levels = append(levels, strings.ToUpper(v))
	}

	levelInit()
}

func init() {
	levelInit()
	fileInit()
}

func levelInit() {

	locker.Lock()
	defer locker.Unlock()

	levelMap = map[string]int{}
	for _, v := range levels {

		if _, ok := levelMap[v]; ok {
			continue
		}

		levelMap[v] = len(levelMap)
	}
}

func (e *entry) line() string {

	tfmt := e.ltime.Format("2006-01-02 15:04:05.000000")
	logLine := fmt.Sprintf("%s %s:%d] %s", tfmt, e.fileName, e.lineNumber, e.level)

	if e.ptype == printDefault {
		logLine += " " + fmt.Sprint(e.args...)
	} else if e.ptype == printFormat {
		logLine += " " + fmt.Sprintf(e.format, e.args)
	}

	return logLine
}

func newEntry(ptype uint8, level, format string, a ...interface{}) {

	level = strings.ToUpper(level)
	if _, ok := levelMap[level]; !ok {
		return
	}

	// It's always the same number of frames to the user's call.
	_, fileName, lineNumber, ok := runtime.Caller(2)
	if !ok {
		fileName = "?"
		lineNumber = 1
	} else {
		slash := strings.LastIndex(fileName, "/")
		if slash >= 0 {
			fileName = fileName[slash+1:]
		}
	}
	if lineNumber < 0 {
		lineNumber = 0 // not a real line number, but acceptable to someDigits
	}

	bufs <- &entry{
		ptype:      ptype,
		level:      level,
		format:     format,
		fileName:   fileName,
		lineNumber: lineNumber,
		args:       a,
		ltime:      time.Now(),
	}
}

func Print(level string, a ...interface{}) {
	newEntry(printDefault, level, "", a...)
}

func Printf(level, format string, a ...interface{}) {
	newEntry(printFormat, level, format, a...)
}
