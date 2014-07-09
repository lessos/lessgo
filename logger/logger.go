package logger

import (
	"log"
	"os"
)

var (
	lgrpath = "/var/log/lessgo.log"
	lgrfile *os.File
	lgr     *log.Logger
)

func Initialize(path string) error {

	lgrpath = path

	return nil
}

func Logf(format string, a ...interface{}) {

	if lgrfile == nil || lgr == nil {

		lgrfile, err := os.OpenFile(lgrpath,
			os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return
		}

		lgr = log.New(lgrfile, "", log.Ldate|log.Ltime)
	}

	lgr.Printf(format, a...)
}
