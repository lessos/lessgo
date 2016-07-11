// Copyright 2013-2016 lessgo Author, All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logger

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	locker    sync.Mutex
	levels    = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
	levelChar = "DIWEF"
	levelMap  = map[string]int{}
	levelOut  = map[string]int{}
	bufs      = make(chan *entry, 100000)
	onceCmd   sync.Once
	log_ws    = map[int]*logFileWriter{}

	logDir      = flag.String("log_dir", "", "If non-empty, write log files in this directory")
	logToStderr = flag.Bool("logtostderr", false, "log to standard error instead of files")
	minLogLevel = flag.Int("minloglevel", 1, "Messages logged at a lower level than this"+
		" don't actually get logged anywhere")
	logToLevels = flag.Bool("logtolevels", false, "Write log to multi level files")
)

const (
	printDefault uint8 = iota
	printFormat
)

type entry struct {
	ptype      uint8
	level      int
	format     string
	fileName   string
	lineNumber int
	ltime      time.Time
	args       []interface{}
}

func init() {

	if !flag.Parsed() && len(os.Args) > 1 {

		for i := 1; i < len(os.Args); i++ {

			if os.Args[i][0] != '-' {
				continue
			}

			key, val, ok := args_parse(i)
			if !ok {
				continue
			}

			switch key {

			case "log_dir":
				*logDir = val

			case "logtostderr":
				if val == "true" {
					*logToStderr = true
				}

			case "minLogLevel":
				if v, err := strconv.Atoi(val); err == nil {
					*minLogLevel = v
				}

			case "logtolevels":
				if val == "true" {
					*logToLevels = true
				}
			}
		}
	}

	levelInit()
}

func args_parse(i int) (string, string, bool) {

	key := strings.TrimLeft(os.Args[i], "-")

	if si := strings.Index(key, "="); si > 0 {

		if key[si+1:] != "" {
			return key[:si], key[si+1:], true
		}

		key = key[:si]
	}

	if (i+2) > len(os.Args) || os.Args[i+1][0] == '-' {
		return key, "", false
	}

	if val := strings.TrimLeft(os.Args[i+1], "="); val != "" {
		return key, val, true
	}

	if (i+3) > len(os.Args) || os.Args[i+2][0] == '-' {
		return key, "", false
	}

	return key, os.Args[i+2], true
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

func levelInit() {

	locker.Lock()
	defer locker.Unlock()

	//
	for _, wr := range log_ws {
		wr.Close()
	}
	log_ws = map[int]*logFileWriter{}

	//
	levelMap = map[string]int{}
	levelChar = ""
	for _, tag := range levels {

		if _, ok := levelMap[tag]; !ok {
			levelMap[tag] = len(levelMap)
			levelChar += tag[0:1]
		}
	}

	if *minLogLevel < 0 {
		*minLogLevel = 0
	} else if *minLogLevel >= len(levelMap) {
		*minLogLevel = len(levelMap) - 1
	}

	//
	for tag, level := range levelMap {

		//
		if (*logToLevels == true && level >= *minLogLevel) ||
			(*logToLevels == false && level == *minLogLevel) {

			levelOut[tag] = level
		}
	}

	onceCmd.Do(outputAction)
}

func (e *entry) line() string {

	logLine := fmt.Sprintf("%s %s %s:%d] ", string(levelChar[e.level]),
		e.ltime.Format("2006-01-02 15:04:05.000000"), e.fileName, e.lineNumber)

	if e.ptype == printDefault {
		logLine += fmt.Sprint(e.args...)
	} else if e.ptype == printFormat {
		logLine += fmt.Sprintf(e.format, e.args...)
	}

	return logLine + "\n"
}

func newEntry(ptype uint8, level_tag, format string, a ...interface{}) {

	level_tag = strings.ToUpper(level_tag)

	level, ok := levelMap[level_tag]
	if !ok || level < *minLogLevel {
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

func outputAction() {

	go func() {

		for logEntry := range bufs {

			bs := []byte(logEntry.line())

			if *logToStderr {
				os.Stderr.Write(bs)
			}

			if len(*logDir) > 0 {

				for level_tag, level := range levelOut {

					if logEntry.level < level {
						continue
					}

					locker.Lock()
					wr, _ := log_ws[level]

					if wr == nil {
						wr = newLogFileWriter(level_tag)
						log_ws[level] = wr
					}

					locker.Unlock()

					wr.Write(bs)
				}
			}
		}
	}()
}
