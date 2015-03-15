package servx

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type DataRedisService struct {
	cfg  DataServiceConfig
	args []string
	err  error
}

func NewDataRedisService(cfg DataServiceConfig) (sv DataRedisService) {

	sv.cfg = cfg

	if len(sv.cfg.Addr) < 8 {
		sv.cfg.Addr = "127.0.0.1"
	}

	if _, err := os.Stat(sv.cfg.Exec); err != nil && os.IsNotExist(err) {
		sv.err = errors.New(fmt.Sprintf("Exec (%s) is not exists", sv.cfg.Exec))
		return
	}

	if sv.cfg.Dir == "" {
		sv.err = errors.New("No Dir Found")
		return
	}

	if _, err := os.Stat(sv.cfg.Dir); err != nil && os.IsNotExist(err) {
		os.MkdirAll(sv.cfg.Dir, 0755)
	}

	sv.args = []string{"--daemonize", "yes",
		"--dir", sv.cfg.Dir,
		"--dbfilename", sv.cfg.InstanceID + ".rdb"}

	if sv.cfg.Port > 0 && sv.cfg.Addr != "" {
		sv.args = append(sv.args, "--port", fmt.Sprintf("%d", sv.cfg.Port))
	}

	if sv.cfg.Socket != "" {
		sv.args = append(sv.args, "--unixsocket", sv.cfg.Socket, "--port", "7777")
	}

	return sv
}

func (s DataRedisService) Error() error {
	return s.err
}

func (s DataRedisService) Start() error {

	if n, _ := s.processID(); n > 0 {
		return nil
	}

	s.err = exec.Command(s.cfg.Exec, s.args...).Run()
	time.Sleep(5e8)

	return s.err
}

func (s DataRedisService) Stop() error {

	if n, ids := s.processID(); n > 0 {

		for _, pid := range ids {
			exec.Command("/bin/kill", "-9", pid).Run()
		}

		time.Sleep(5e8)
	}

	return nil
}

func (s DataRedisService) processID() (n int, pids []string) {

	pid, err := exec.Command("/bin/sh", "-c",
		"ps aux |grep "+s.cfg.Exec+"|grep -v grep|grep "+s.cfg.InstanceID+"|awk '{print $2}'").Output()
	if err == nil && len(pid) > 0 {
		pids = strings.Split(strings.TrimSpace(string(pid)), "\n")
	}

	return len(pids), pids
}
