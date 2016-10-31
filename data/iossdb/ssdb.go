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

package iossdb

import (
	"fmt"
	"net"
	"runtime"
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

	if cfg.Timeout < 3 {
		cfg.Timeout = 3
	} else if cfg.Timeout > 600 {
		cfg.Timeout = 600
	}

	cr := &Connector{
		ctype:    "tcp",
		clink:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		ctimeout: time.Duration(cfg.Timeout) * time.Second,
		conns:    make(chan *Client, cfg.MaxConn),
		config:   cfg,
	}

	for i := 0; i < cfg.MaxConn; i++ {

		cn, err := dialTimeout(cr.ctype, cr.clink)
		if err != nil {
			return cr, err
		}

		if cr.config.Auth != "" {
			cn.Cmd("auth", cr.config.Auth)
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

	cn.sock.SetDeadline(time.Now().Add(cr.ctimeout))

	var rpl *Reply

	for try := 1; try <= 3; try++ {

		rpl = cn.Cmd(args...)
		if rpl.State != ReplyFail {
			break
		}

		time.Sleep(time.Duration(try) * time.Second)

		if cn0, err := dialTimeout(cr.ctype, cr.clink); err == nil {

			cn = cn0
			cn.sock.SetDeadline(time.Now().Add(cr.ctimeout))

			if cr.config.Auth != "" {
				cn.Cmd("auth", cr.config.Auth)
			}
		}
	}

	cr.push(cn)

	return rpl
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
