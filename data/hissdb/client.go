package hissdb

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
)

type Client struct {
	sock     *net.TCPConn
	recv_buf bytes.Buffer
}

func Connect(ip string, port int) (*Client, error) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return nil, err
	}
	sock, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}
	var c Client
	c.sock = sock
	return &c, nil
}

func (c *Client) Cmd(args ...interface{}) *Reply {

	r := &Reply{
		State: ReplyError,
		Data:  []string{},
	}

	if err := c.send(args); err != nil {
		r.State = ReplyFail
		return r
	}

	resp, err := c.recv()
	if err != nil || len(resp) < 1 {
		return r
	}

	switch resp[0] {
	case ReplyOK, ReplyNotFound, ReplyError, ReplyFail, ReplyClientError:
		r.State = resp[0]
	}

	if r.State == ReplyOK {
		for k, v := range resp {

			if k == 0 {
				continue
			}

			r.Data = append(r.Data, v)
		}
	}

	return r
}

func (c *Client) send(args []interface{}) error {
	var buf bytes.Buffer
	for _, arg := range args {
		var s string
		switch arg := arg.(type) {
		case string:
			s = arg
		case []byte:
			s = string(arg)
		case []string:
			for _, s := range arg {
				buf.WriteString(fmt.Sprintf("%d", len(s)))
				buf.WriteByte('\n')
				buf.WriteString(s)
				buf.WriteByte('\n')
			}
			continue
		case uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64:
			s = fmt.Sprintf("%d", arg)
		case float64:
			s = fmt.Sprintf("%f", arg)
		case bool:
			if arg {
				s = "1"
			} else {
				s = "0"
			}
		case nil:
			s = ""
		default:
			return fmt.Errorf("bad arguments")
		}
		buf.WriteString(fmt.Sprintf("%d", len(s)))
		buf.WriteByte('\n')
		buf.WriteString(s)
		buf.WriteByte('\n')
	}
	buf.WriteByte('\n')
	_, err := c.sock.Write(buf.Bytes())
	if err != nil {
		// TODO
	}
	return err
}

func (c *Client) recv() ([]string, error) {
	var tmp [8192]byte
	for {
		n, err := c.sock.Read(tmp[0:])
		if err != nil {
			return nil, err
		}
		c.recv_buf.Write(tmp[0:n])
		resp := c.parse()
		if resp == nil || len(resp) > 0 {
			return resp, nil
		}
	}
}

func (c *Client) parse() []string {
	resp := []string{}
	buf := c.recv_buf.Bytes()
	var idx, offset int
	idx = 0
	offset = 0

	for {
		idx = bytes.IndexByte(buf[offset:], '\n')
		if idx == -1 {
			break
		}
		p := buf[offset : offset+idx]
		offset += idx + 1
		//fmt.Printf("> [%s]\n", p);
		if len(p) == 0 || (len(p) == 1 && p[0] == '\r') {
			if len(resp) == 0 {
				continue
			} else {
				c.recv_buf.Next(offset)
				return resp
			}
		}

		size, err := strconv.Atoi(string(p))
		if err != nil || size < 0 {
			return nil
		}
		if offset+size >= c.recv_buf.Len() {
			break
		}

		v := buf[offset : offset+size]
		resp = append(resp, string(v))
		offset += size + 1
	}

	return []string{}
}

// Close The Client Connection
func (c *Client) Close() error {
	return c.sock.Close()
}
