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
	"errors"
	"time"
)

const (
	// ssdb-server total amount of the commands and arguments must be less than 10MB.
	batch_cmds_max int64 = 4194304 // 4 MB
)

type Batch struct {
	cmds [][]interface{}
	cr   *Connector
}

type batch_cmds struct {
	num int
	req []byte
}

// Batch is in DEVELOPMENT PREVIEW, DO NOT USE IT IN PRODUCTION
//
// This feature is implemented in client side, the ssdb-server does not support
// batch command execution. The server will execute each command as if they are
// separated.
func (cr *Connector) Batch() *Batch {

	return &Batch{
		cmds: [][]interface{}{},
		cr:   cr,
	}
}

func (b *Batch) Cmd(args ...interface{}) {
	b.cmds = append(b.cmds, args)
}

func (b *Batch) Exec() ([]*Reply, error) {

	var (
		rpls            = []*Reply{}
		cmds            = []batch_cmds{}
		cmds_size int64 = 0
	)

	if len(b.cmds) < 1 || b.cr == nil {
		return rpls, errors.New("client error")
	}

	for _, args := range b.cmds {

		if buf, err := send_buf(args); err == nil {

			n := int64(len(buf))

			if n > batch_cmds_max {
				return rpls, errors.New("client error")
			}

			if len(cmds) < 1 || (cmds_size+n) > batch_cmds_max {
				cmds = append(cmds, batch_cmds{
					num: 0,
					req: []byte{},
				})
				cmds_size = 0
			}

			cmds[len(cmds)-1].num++
			cmds[len(cmds)-1].req = append(cmds[len(cmds)-1].req, buf...)

			cmds_size += n

		} else {
			return rpls, errors.New("client error")
		}
	}

	if len(cmds) < 1 {
		return rpls, errors.New("client error")
	}

	cn, _ := b.cr.pull()

	tto := b.cr.ctimeout * 10
	if tto > (600 * time.Second) {
		tto = 600 * time.Second
	}

	for _, bcmd := range cmds {

		tton := time.Now().Add(tto)
		cn.sock.SetReadDeadline(tton)
		cn.sock.SetWriteDeadline(tton)

		if _, err := cn.sock.Write(bcmd.req); err != nil {

			b.cr.push(cn)
			b.cmds = [][]interface{}{}
			b.cr = nil

			return rpls, err
		}

		for i := 0; i < bcmd.num; i++ {

			r := &Reply{
				State: ReplyError,
			}

			resp, err := cn.recv()

			if err != nil {

				r.State = err.Error()

			} else if len(resp) < 1 {

				r.State = ReplyFail

			} else {

				switch resp[0].String() {
				case ReplyOK, ReplyNotFound, ReplyError, ReplyFail, ReplyClientError:
					r.State = resp[0].String()
				}

				if r.State == ReplyOK {

					for k, v := range resp {

						if k > 0 {
							r.Data = append(r.Data, v)
						}
					}
				}
			}

			rpls = append(rpls, r)
		}
	}

	b.cr.push(cn)
	b.cmds = [][]interface{}{}
	b.cr = nil

	return rpls, nil
}
