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
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type Request struct {
	*http.Request
	ContentType    string
	AcceptLanguage []AcceptLanguage
	Locale         string
	RequestPath    string
	RawBody        []byte
	WebSocket      WebSocket
}

func NewRequest(r *http.Request) *Request {

	req := &Request{
		Request:        r,
		ContentType:    resolveContentType(r),
		AcceptLanguage: resolveAcceptLanguage(r),
		Locale:         "",
		RawBody:        []byte{},
	}

	if req.Body != nil {

		if body, err := ioutil.ReadAll(req.Body); err == nil {
			req.RawBody = body
			req.Body = ioutil.NopCloser(bytes.NewReader(req.RawBody))
		}
	}

	return req
}

func (req *Request) RawAbsUrl() string {

	scheme := "http"

	if len(req.URL.Scheme) > 0 {
		scheme = req.URL.Scheme
	}

	return fmt.Sprintf("%s://%s%s", scheme, req.Host, req.RequestURI)
}

// Get the content type.
// e.g. From "multipart/form-data; boundary=--" to "multipart/form-data"
// If none is specified, returns "text/html" by default.
func resolveContentType(r *http.Request) string {

	contentType := r.Header.Get("Content-Type")

	if contentType == "" {
		return "text/html"
	}

	return strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
}

// Resolve the Accept-Language header value.
//
// The results are sorted using the quality defined in the header for each language range with the
// most qualified language range as the first element in the slice.
//
// See the HTTP header fields specification
// (http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.4) for more details.
func resolveAcceptLanguage(r *http.Request) acceptLanguages {

	var (
		rals = strings.Split(r.Header.Get("Accept-Language"), ",")
		als  = make(acceptLanguages, len(rals))
	)

	for i, v := range rals {

		if v2 := strings.Split(v, ";q="); len(v2) == 2 {
			quality, err := strconv.ParseFloat(v2[1], 32)
			if err != nil {
				quality = 1
			}
			als[i] = AcceptLanguage{v2[0], float32(quality)}
		} else if len(v2) == 1 {
			als[i] = AcceptLanguage{v, 1}
		}
	}

	sort.Sort(als)

	return als
}

// A single language from the Accept-Language HTTP header.
type AcceptLanguage struct {
	Language string
	Quality  float32
}

// A collection of sortable AcceptLanguage instances.
type acceptLanguages []AcceptLanguage

func (al acceptLanguages) Len() int           { return len(al) }
func (al acceptLanguages) Swap(i, j int)      { al[i], al[j] = al[j], al[i] }
func (al acceptLanguages) Less(i, j int) bool { return al[i].Quality > al[j].Quality }
func (al acceptLanguages) String() string {
	output := bytes.NewBufferString("")
	for i, language := range al {
		output.WriteString(fmt.Sprintf("%s (%1.1f)", language.Language, language.Quality))
		if i != len(al)-1 {
			output.WriteString(", ")
		}
	}
	return output.String()
}
