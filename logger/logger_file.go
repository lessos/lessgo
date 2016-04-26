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
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	file_max_size = uint64(1024 * 1024 * 1024)
	flocker       sync.Mutex
	pid           = os.Getpid()
	program       = filepath.Base(os.Args[0])
	host          = "unknownhost"
	userName      = "unknownuser"
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
}

type logFileWriter struct {
	mu        sync.Mutex
	fp        *os.File
	nbytes    uint64
	level_tag string
}

func newLogFileWriter(level_tag string) *logFileWriter {
	return &logFileWriter{
		level_tag: level_tag,
	}
}

func (w *logFileWriter) Write(bs []byte) {

	w.mu.Lock()
	defer w.mu.Unlock()

	if w.fp != nil && w.nbytes > file_max_size {
		w.fp.Sync()
		w.fp.Close()
		w.fp = nil
		w.nbytes = 0
	}

	if w.fp == nil {
		fp, err := fileOpenInit(w.level_tag)
		if err != nil {
			return
		}
		w.fp = fp
	}

	if w.fp != nil {
		if n, err := w.fp.Write(bs); err != nil {
			w.fp = nil
		} else {
			w.nbytes += uint64(n)
		}
	}
}

func (w *logFileWriter) Close() {

	w.mu.Lock()
	defer w.mu.Unlock()

	if w.fp == nil {
		return
	}

	w.fp.Sync()
	w.fp.Close()
	w.fp = nil
	w.nbytes = 0
}

func fileOpenInit(tag string) (*os.File, error) {

	flocker.Lock()
	defer flocker.Unlock()

	if len(*logDir) < 1 {
		return nil, errors.New("No -log_dir Setup")
	}

	if _, err := os.Stat(*logDir); err != nil {
		return nil, err
	}

	logName, link := logName(tag, time.Now())

	fp, err := os.OpenFile(*logDir+"/"+logName,
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

	if err == nil {

		symlink := filepath.Join(*logDir, link)

		os.Remove(symlink)           // ignore err
		os.Symlink(logName, symlink) // ignore err

		var buf bytes.Buffer
		fmt.Fprintf(&buf, "Log file created at: %s\n",
			time.Now().Format("2006-01-02 15:04:05.000000"))
		fmt.Fprintf(&buf, "Running on machine: %s\n", host)
		fmt.Fprintf(&buf, "Binary: Built with %s %s for %s/%s\n",
			runtime.Compiler, runtime.Version(), runtime.GOOS, runtime.GOARCH)
		fmt.Fprintf(&buf, "Log line format: [%s] yyyy-mm-dd hh:ii:ss.uuuuuu file:line] msg\n", levelChar)
		fp.Write(buf.Bytes())
	}

	return fp, err
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
func logName(tag string, t time.Time) (name, link string) {

	name = fmt.Sprintf("%s.%s.%s.log.%s.%s.%d",
		program, host, userName, tag, t.Format("20060102-150405"), pid)

	return name, program + ".log." + tag
}
