package servx

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"
)

const ssdb_config = `work_dir = {{.work_dir}}
pidfile = {{.pidfile}}

server:
	ip: {{.server_ip}}
	port: {{.server_port}}
	{{.auth}}

logger:
	level: error
	output: error.log
	rotate:
		size: 100000000

leveldb:
	# in MB
	cache_size: 500
	# in KB
	block_size: 32
	# in MB
	write_buffer_size: 64
	# in MB
	compaction_speed: 1000
	# yes|no
	compression: yes

`

type DataSSDBService struct {
	cfg     DataServiceConfig
	cfgsets map[string]string
	err     error
}

func NewDataSSDBService(cfg DataServiceConfig) (sv DataSSDBService) {

	sv.cfg = cfg

	if sv.cfg.Port < 1 {
		sv.err = errors.New(fmt.Sprintf("Server Port Not Set"))
		return
	}

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

	if len(sv.cfg.Pass) > 0 && len(sv.cfg.Pass) < 32 {
		sv.err = errors.New("Server Auth-Password must be at least 32 characters")
		return
	}

	if _, err := os.Stat(sv.cfg.Dir); err != nil && os.IsNotExist(err) {
		os.MkdirAll(sv.cfg.Dir, 0755)
	}

	sv.cfgsets = map[string]string{
		"work_dir":    sv.cfg.Dir,
		"pidfile":     sv.cfg.Dir + "/ssdb.pid",
		"server_ip":   sv.cfg.Addr,
		"server_port": fmt.Sprintf("%d", sv.cfg.Port),
		"auth":        "",
	}

	if len(sv.cfg.Pass) > 31 {
		sv.cfgsets["auth"] = "auth: " + sv.cfg.Pass
	}

	tpl, err := template.New("def").Parse(ssdb_config)
	if err != nil {
		sv.err = err
		return
	}

	var buf bytes.Buffer

	if sv.err = tpl.Execute(&buf, sv.cfgsets); sv.err != nil {
		return
	}

	fp, err := os.OpenFile(sv.cfg.Dir+"/"+sv.cfg.InstanceID+".conf", os.O_RDWR|os.O_CREATE, 0754)
	if err != nil {
		sv.err = err
		return
	}
	defer fp.Close()

	fp.Seek(0, 0)
	fp.Truncate(int64(len(buf.String())))

	if _, sv.err = fp.WriteString(buf.String()); sv.err != nil {
		return
	}

	return sv
}

func (s DataSSDBService) Error() error {
	return s.err
}

func (s DataSSDBService) Start() error {

	if n, _ := s.processID(); n > 0 {
		return nil
	}

	s.err = exec.Command(s.cfg.Exec, "-d", s.cfg.Dir+"/"+s.cfg.InstanceID+".conf", "-s", "restart").Run()
	time.Sleep(5e8)

	return s.err
}

func (s DataSSDBService) Stop() error {

	if n, ids := s.processID(); n > 0 {

		for _, pid := range ids {
			exec.Command("/bin/kill", "-9", pid).Run()
		}

		time.Sleep(5e8)
	}

	return nil
}

func (s DataSSDBService) processID() (n int, pids []string) {

	pid, err := exec.Command("/bin/sh", "-c",
		"ps aux |grep "+s.cfg.Exec+"|grep -v grep|grep "+s.cfg.InstanceID+"|awk '{print $2}'").Output()
	if err == nil && len(pid) > 0 {
		pids = strings.Split(strings.TrimSpace(string(pid)), "\n")
	}

	return len(pids), pids
}
