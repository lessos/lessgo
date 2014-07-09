package redis

import (
	"runtime"
	"sync"
	"time"
)

type Connector struct {
	cfg    Config
	conns  chan *conn
	locker sync.Mutex     // TODO
	stats  map[string]int // TODO
}

type conn struct {
	client *Client
	locker sync.Mutex     // TODO
	stats  map[string]int // TODO
}

func NewConnector(cfg Config) (*Connector, error) {

	if cfg.MaxConn < 1 {

		cfg.MaxConn = 1

	} else {

		maxconn := runtime.NumCPU() * 2
		if maxconn > 100 {
			maxconn = 100
		}

		if cfg.MaxConn > maxconn {
			cfg.MaxConn = maxconn
		}
	}

	timeout := time.Duration(cfg.Timeout) * time.Second
	if timeout < 1*time.Second {
		timeout = 1 * time.Second
	}

	c := Connector{
		cfg:   cfg,
		conns: make(chan *conn, cfg.MaxConn),
	}

	for i := 0; i < cfg.MaxConn; i++ {

		ctype, clink := "tcp", cfg.Host+":"+cfg.Port
		if len(cfg.Socket) > 1 {
			ctype, clink = "unix", cfg.Socket
		}

		cn, err := DialTimeout(ctype, clink, timeout)
		if err != nil {
			return &c, err
		}
		c.conns <- &conn{client: cn}
	}

	return &c, nil
}

func (c *Connector) Cmd(cmd string, args ...interface{}) *Reply {

	cn, _ := c.pull()
	defer c.push(cn)

	return cn.client.Cmd(cmd, args...)
}

func (c *Connector) push(cn *conn) {
	c.conns <- cn
}

func (c *Connector) pull() (cn *conn, err error) {
	return <-c.conns, nil
}
