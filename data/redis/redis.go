package redis

import (
	"runtime"
	"sync"
	"time"
)

type Connector struct {
	ctype    string
	clink    string
	ctimeout time.Duration
	conns    chan *conn
	locker   sync.Mutex     // TODO
	stats    map[string]int // TODO
	config   Config
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

	c := Connector{
		ctype:    "tcp",
		clink:    cfg.Host + ":" + cfg.Port,
		ctimeout: time.Duration(cfg.Timeout) * time.Second,
		conns:    make(chan *conn, cfg.MaxConn),
		config:   cfg,
	}

	if c.ctimeout < 1*time.Second {
		c.ctimeout = 1 * time.Second
	}

	if len(cfg.Socket) > 1 {
		c.ctype, c.clink = "unix", cfg.Socket
	}

	for i := 0; i < cfg.MaxConn; i++ {

		cn, err := DialTimeout(c.ctype, c.clink, c.ctimeout)
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

	rs := cn.client.Cmd(cmd, args...)

	if rs.Err != nil && rs.Err.Error() == "use of closed network connection" {

		time.Sleep(1e9)

		if cnc, err := DialTimeout(c.ctype, c.clink, c.ctimeout); err == nil {
			cn.client = cnc
		}
	}

	return rs
}

func (cr *Connector) Close() {

	for i := 0; i < cr.config.MaxConn; i++ {
		cn, _ := cr.pull()
		cn.client.Close()
	}
}

func (c *Connector) push(cn *conn) {
	c.conns <- cn
}

func (c *Connector) pull() (cn *conn, err error) {
	return <-c.conns, nil
}
