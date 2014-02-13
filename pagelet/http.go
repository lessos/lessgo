package pagelet

import (
    "net/http"
)

type Request struct {
    *http.Request
    ContentType string
}

type Response struct {
    Status      int
    ContentType string
    Out         http.ResponseWriter
}

func NewResponse(w http.ResponseWriter) *Response {
    return &Response{Out: w}
}

func NewRequest(r *http.Request) *Request {
    return &Request{
        Request:     r,
        ContentType: "text/html",
    }
}

func (resp *Response) WriteHeader(status int, ctype string) {

    if resp.Status == 0 {
        resp.Status = status
        resp.Out.WriteHeader(resp.Status)
    }

    if resp.ContentType == "" {
        resp.ContentType = ctype
        resp.Out.Header().Set("Content-Type", resp.ContentType)
    }
}
