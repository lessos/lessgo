package httpclient

import (
    "../../utils"
    "bytes"
    "crypto/tls"
    "io/ioutil"
    "net"
    "net/http"
    "net/url"
    "time"
)

var (
    defaultUserAgent = "lessgoHttpClient"
    defaultTimeout   = 60000 * time.Millisecond
)

type HttpClientRequest struct {
    req             *http.Request
    url             string
    timeout         time.Duration
    tlsClientConfig *tls.Config
    rsp             *http.Response
}

func NewHttpClientRequest(method, url string) *HttpClientRequest {

    var req http.Request
    req.Method = method
    req.Header = http.Header{}
    req.Header.Set("User-Agent", defaultUserAgent)

    return &HttpClientRequest{
        req:     &req,
        url:     url,
        timeout: defaultTimeout,
    }
}

// SetTimeout sets connect time out.
func (c *HttpClientRequest) SetTimeout(timeout time.Duration) *HttpClientRequest {
    c.timeout = timeout * time.Millisecond
    return c
}

// SetTLSClientConfig sets tls connection configurations if visiting https url.
func (c *HttpClientRequest) SetTLSClientConfig(config *tls.Config) *HttpClientRequest {
    c.tlsClientConfig = config
    return c
}

// Header add header item string in request.
func (c *HttpClientRequest) Header(key, value string) *HttpClientRequest {
    c.req.Header.Set(key, value)
    return c
}

// SetCookie add cookie into request.
func (c *HttpClientRequest) SetCookie(cookie *http.Cookie) *HttpClientRequest {
    c.req.Header.Add("Cookie", cookie.String())
    return c
}

// Get returns *HttpClientRequest with GET method.
func Get(url string) *HttpClientRequest {
    return NewHttpClientRequest("GET", url)
}

// Post returns *HttpClientRequest with POST method.
func Post(url string) *HttpClientRequest {
    return NewHttpClientRequest("POST", url)
}

// Put returns *HttpClientRequest with PUT method.
func Put(url string) *HttpClientRequest {
    return NewHttpClientRequest("PUT", url)
}

// Delete returns *HttpClientRequest with DELETE method.
func Delete(url string) *HttpClientRequest {
    return NewHttpClientRequest("DELETE", url)
}

// Head returns *HttpClientRequest with HEAD method.
func Head(url string) *HttpClientRequest {
    return NewHttpClientRequest("HEAD", url)
}

// Body adds request raw body.
// it supports string and []byte.
func (c *HttpClientRequest) Body(data interface{}) *HttpClientRequest {
    switch t := data.(type) {
    case string:
        bf := bytes.NewBufferString(t)
        c.req.Body = ioutil.NopCloser(bf)
        c.req.ContentLength = int64(len(t))
    case []byte:
        bf := bytes.NewBuffer(t)
        c.req.Body = ioutil.NopCloser(bf)
        c.req.ContentLength = int64(len(t))
    }
    return c
}

// Response executes request client gets response mannually.
func (c *HttpClientRequest) Response() (*http.Response, error) {

    if c.rsp != nil {
        return c.rsp, nil
    }

    url, err := url.Parse(c.url)
    if err != nil {
        return nil, err
    }
    c.req.URL = url

    client := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: c.tlsClientConfig,
            Dial:            timeoutDialer(c.timeout, c.timeout),
        },
    }
    c.rsp, err = client.Do(c.req)
    if err != nil {
        return nil, err
    }

    return c.rsp, nil
}

// timeoutDialer returns functions of connection dialer with timeout settings for http.Transport Dial field.
func timeoutDialer(cTimeout, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
    return func(netw, addr string) (net.Conn, error) {
        conn, err := net.DialTimeout(netw, addr, cTimeout)
        if err != nil {
            return nil, err
        }
        conn.SetDeadline(time.Now().Add(rwTimeout))
        return conn, nil
    }
}

// Bytes returns the body []byte in response.
// it calls Response inner.
func (c *HttpClientRequest) ReplyBytes() ([]byte, error) {
    resp, err := c.Response()
    if err != nil {
        return nil, err
    }
    if resp.Body == nil {
        return nil, nil
    }
    defer resp.Body.Close()

    data, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    return data, nil
}

// String returns the body string in response.
// it calls Response inner.
func (c *HttpClientRequest) ReplyString() (string, error) {

    data, err := c.ReplyBytes()
    if err != nil {
        return "", err
    }

    return string(data), nil
}

// ReplyJson returns the map that marshals from the body bytes as json in response .
func (c *HttpClientRequest) ReplyJson(v interface{}) error {
    data, err := c.ReplyBytes()
    if err != nil {
        return err
    }
    err = utils.JsonDecode(data, v)
    if err != nil {
        return err
    }
    return nil
}
