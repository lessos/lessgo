package redis

import (
	radix "github.com/fzzy/radix/redis"
	"runtime"
	"sync"
	"time"
)

type Client struct {
	cfg    Config
	conns  chan *conn
	locker sync.Mutex     // TODO
	stats  map[string]int // TODO
}

type conn struct {
	radix  *radix.Client
	locker sync.Mutex     // TODO
	stats  map[string]int // TODO
}

func NewClient(cfg Config) (*Client, error) {

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

	c := Client{
		cfg:   cfg,
		conns: make(chan *conn, cfg.MaxConn),
	}

	for i := 0; i < cfg.MaxConn; i++ {

		ctype, clink := "tcp", cfg.Host+":"+cfg.Port
		if len(cfg.Socket) > 1 {
			ctype, clink = "unix", cfg.Socket
		}

		cn, err := radix.DialTimeout(ctype, clink, timeout)
		if err != nil {
			return &c, err
		}
		c.conns <- &conn{radix: cn}
	}

	return &c, nil
}

func (c *Client) Cmd(cmd string, args ...interface{}) *radix.Reply {

	cn, _ := c.pull()
	defer c.push(cn)

	return cn.radix.Cmd(cmd, args...)
}

func (c *Client) push(cn *conn) {
	c.conns <- cn
}

func (c *Client) pull() (cn *conn, err error) {
	return <-c.conns, nil
}
