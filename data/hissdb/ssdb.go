package hissdb

import (
	"net"
	"runtime"
	"fmt"
	"time"
)

type Connector struct {
	ctype    string
	clink    string
	ctimeout time.Duration
	conns    chan *Client
	config   Config
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

	cr := &Connector{
		ctype:    "tcp",
		clink:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		ctimeout: time.Duration(cfg.Timeout) * time.Second,
		conns:    make(chan *Client, cfg.MaxConn),
		config:   cfg,
	}

	if cr.ctimeout < 1*time.Second {
		cr.ctimeout = 10 * time.Second
	}

	for i := 0; i < cfg.MaxConn; i++ {

		cn, err := dialTimeout(cr.ctype, cr.clink)
		if err != nil {
			return cr, err
		}
		cr.conns <- cn
	}

	return cr, nil
}

func dialTimeout(network, addr string) (*Client, error) {

	raddr, err := net.ResolveTCPAddr(network, addr)
	if err != nil {
		return nil, err
	}
	sock, err := net.DialTCP(network, nil, raddr)
	if err != nil {
		return nil, err
	}

	return &Client{sock: sock}, nil
}

func (cr *Connector) Cmd(args ...interface{}) *Reply {

	cn, _ := cr.pull()
	defer cr.push(cn)

	cn.sock.SetReadDeadline(time.Now().Add(cr.ctimeout))
	cn.sock.SetWriteDeadline(time.Now().Add(cr.ctimeout))

	return cn.Cmd(args...)
}

func (cr *Connector) Close() {

	for i := 0; i < cr.config.MaxConn; i++ {
		cn, _ := cr.pull()
		cn.Close()
	}
}

func (cr *Connector) push(cn *Client) {
	cr.conns <- cn
}

func (cr *Connector) pull() (cn *Client, err error) {
	return <-cr.conns, nil
}
