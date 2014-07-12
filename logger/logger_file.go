package logger

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	logDir    = flag.String("log_dir", "", "If non-empty, write log files in this directory")
	pid       = os.Getpid()
	program   = filepath.Base(os.Args[0])
	host      = "unknownhost"
	userName  = "unknownuser"
	flocker   sync.Mutex
	status    = 0
	out       io.Writer // destination for output // TODO
	outEnable = false   // TODO
)

func fileInit() {

	flocker.Lock()
	defer flocker.Unlock()

	if status == 1 {
		return
	}
	status = 1

	h, err := os.Hostname()
	if err == nil {
		host = shortHostname(h)
	}

	current, err := user.Current()
	if err == nil {
		userName = current.Username
	}

	flag.Parse()

	// Sanitize userName since it may contain filepath separators on Windows.
	userName = strings.Replace(userName, `\`, "_", -1)

	if len(*logDir) > 0 {

		if _, err := os.Stat(*logDir); err == nil {

			logName, _ := logName(time.Now())

			out, err = os.OpenFile(*logDir+"/"+logName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
			if err == nil {
				outEnable = true
			} else {
				fmt.Println("logDir", err)
			}
		}
	}

	go func() {

		for logEntry := range bufs {

			if outEnable {
				if _, err := out.Write([]byte(logEntry.line() + "\n")); err != nil {
					outEnable = false // TODO
				}
			}
		}
	}()
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

	name = fmt.Sprintf("%s.%s.%s_%04d%02d%02d.log",
		program, host, userName,
		t.Year(), t.Month(), t.Day())

	return name, program + ".log"
}
