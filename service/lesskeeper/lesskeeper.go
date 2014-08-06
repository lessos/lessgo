package lesskeeper

import (
	"../../net/httpclient"
	"../../utils"
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Reply struct {
	ErrorCode uint   `json:"errorCode,omitempty"`
	Message   string `json:"message,omitempty"`
	Action    string `json:"action,omitempty"`
	Node      Node   `json:"node,omitempty"`
}

type Node struct {
	Key       string `json:"key"`
	Value     string `json:"value,omitempty"`
	PrevValue string `json:"prevValue,omitempty"`
	Version   uint64 `json:"version,omitempty"`
	Ttl       int    `json:"ttl,omitempty"`
	Dir       bool   `json:"dir,omitempty"`
	Nodes     []Node `json:"nodes,omitempty"`
}

type Client struct {
	ApiUrl    string `json:"apiurl"`
	Bucket    string `json:"bucket"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	Timeout   int    `json:"timeout"`
}

func NewClient(url, bucket, accessKey, secretKey string) (Client, error) {

	c := Client{
		ApiUrl:    url,
		Bucket:    bucket,
		AccessKey: accessKey,
		SecretKey: secretKey,
		Timeout:   10,
	}

	return c, nil
}

func (c Client) signHandler(req *http.Request) {

	h := hmac.New(sha1.New, []byte(c.SecretKey))

	data := req.URL.Path
	if req.URL.RawQuery != "" {
		data += "?" + req.URL.RawQuery
	}
	io.WriteString(h, data+"\n")

	if req.Body != nil {

		s2, err2 := ioutil.ReadAll(req.Body)
		if err2 != nil {
			return
		}
		h.Write(s2)

		req.Body = ioutil.NopCloser(bytes.NewReader(s2))
	}

	sign := base64.URLEncoding.EncodeToString(h.Sum(nil))

	req.Header.Set("Authorization", "LOS "+c.AccessKey+":"+sign)
}

func (c Client) nodeCmd(method, url string, params map[string]string) Reply {

	hc := httpclient.NewHttpClientRequest(method, url)
	hc.SignHandler = c.signHandler
	hc.SetTimeout(time.Duration(c.Timeout) * time.Second)

	for k, v := range params {
		hc.Param(k, v)
	}

	var rsp Reply
	err := hc.ReplyJson(&rsp)
	if err != nil {
		rsp.ErrorCode = 500
	}

	return rsp
}

func (c Client) NodeSet(n Node) Reply {

	params := map[string]string{
		"value": n.Value,
	}
	if n.PrevValue != "" {
		params["prevValue"] = n.PrevValue
	}
	if n.Ttl > 0 {
		params["ttl"] = strconv.Itoa(n.Ttl)
	}

	return c.nodeCmd("PUT", c.ApiUrl+"/data/keys"+n.Key, params)
}

func (c Client) NodeGet(key string) Reply {
	return c.nodeCmd("GET", c.ApiUrl+"/data/keys"+key, map[string]string{})
}

func (c Client) NodeDel(key string) Reply {
	return c.nodeCmd("DELETE", c.ApiUrl+"/data/keys"+key, map[string]string{})
}

// Json returns the map that marshals from the reply value string as json in response .
func (n *Node) Json(v interface{}) error {
	return utils.JsonDecode(n.Value, v)
}
