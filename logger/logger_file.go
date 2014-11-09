package logger

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	flocker      sync.Mutex
	onceFileInit sync.Once
	out          *os.File // destination for output // TODO
	logDir       = flag.String("log_dir", "", "If non-empty, write log files in this directory")

	// file name format args
	pid      = os.Getpid()
	program  = filepath.Base(os.Args[0])
	host     = "unknownhost"
	userName = "unknownuser"
)

func init() {

	h, err := os.Hostname()
	if err == nil {
		host = shortHostname(h)
	}

	current, err := user.Current()
	if err == nil {
		userName = current.Username
	}

	// Sanitize userName since it may contain filepath separators on Windows.
	userName = strings.Replace(userName, `\`, "_", -1)

	//
	onceFileInit.Do(fileInit)
}

func fileInit() {

	go func() {

		for logEntry := range bufs {

			if out == nil && len(*logDir) > 0 {
				fileOpenInit()
			}

			if out != nil {
				if _, err := out.Write([]byte(logEntry.line() + "\n")); err != nil {
					out = nil
				}
			}
		}
	}()
}

func fileOpenInit() {

	flocker.Lock()
	defer flocker.Unlock()

	if len(*logDir) < 1 {
		return
	}

	if _, err := os.Stat(*logDir); err == nil {

		logName, _ := logName(time.Now())

		if out, err = os.OpenFile(*logDir+"/"+logName,
			os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644); err != nil {
			out = nil
		}
	}
}

// shortHostname returns its argument, truncating at the first period.
// For instance, given "www.google.com" it returns "www".
func shortHostname(hostname string) string {
	if i := strings.Index(hostname, "."); i >= 0 {
		return hostname[:i]
	}
	return hostname
}

// logName returns a new log file name containing tag, with start time t, and
// the name for the symlink for tag.
//  func logName(tag string, t time.Time) (name, link string) {
func logName(t time.Time) (name, link string) {

	name = fmt.Sprintf("%s.%s.%s.%04d%02d%02d.%d.log",
		program, host, userName,
		t.Year(), t.Month(), t.Day(), pid)

	return name, program + ".log"
}
