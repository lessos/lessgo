// Copyright 2015 lessOS.com, All rights reserved.
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

package httpsrv

import (
	"errors"

	"github.com/lessos/lessgo/deps/go.net/websocket"
)

var (
	wsErrNoConn = errors.New("No WebSocket Connection")
)

type WebSocket struct {
	conn *websocket.Conn
}

func (ws WebSocket) Receive(v interface{}) error {

	if ws.conn == nil {
		return wsErrNoConn
	}

	if err := websocket.Message.Receive(ws.conn, v); err != nil {
		ws.conn.Close()
		return err
	}

	return nil
}

func (ws WebSocket) Send(v interface{}) error {

	if ws.conn == nil {
		return wsErrNoConn
	}

	err := websocket.Message.Send(ws.conn, v)
	if err != nil {
		ws.conn.Close()
	}

	return err
}

func (ws WebSocket) SendJson(v interface{}) error {

	if ws.conn == nil {
		return wsErrNoConn
	}

	err := websocket.JSON.Send(ws.conn, v)
	if err != nil {
		ws.conn.Close()
	}

	return err
}

func (ws WebSocket) Close() {

	if ws.conn == nil {
		return
	}

	ws.conn.Close()
	ws.conn = nil
}
