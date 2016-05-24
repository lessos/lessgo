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

type Batch struct {
	cmds [][]interface{}
	cr   *Connector
}

// Batch is in DEVELOPMENT PREVIEW, DO NOT USE IT IN PRODUCTION
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
		rpls = []*Reply{}
		cmdn = 0
		cmds = []byte{}
	)

	if len(b.cmds) < 1 || b.cr == nil {
		return rpls, errors.New("client error")
	}

	for _, args := range b.cmds {

		if buf, err := send_buf(args); err == nil {
			cmds = append(cmds, buf...)
			cmdn++
		}
	}

	if cmdn < 1 {
		return rpls, errors.New("client error")
	}

	cn, _ := b.cr.pull()

	cn.sock.SetReadDeadline(time.Now().Add(b.cr.ctimeout))
	cn.sock.SetWriteDeadline(time.Now().Add(b.cr.ctimeout))

	if _, err := cn.sock.Write(cmds); err != nil {

		b.cr.push(cn)
		b.cmds = [][]interface{}{}
		b.cr = nil

		return rpls, err
	}

	for i := 0; i < cmdn; i++ {

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

	b.cr.push(cn)
	b.cmds = [][]interface{}{}
	b.cr = nil

	return rpls, nil
}
