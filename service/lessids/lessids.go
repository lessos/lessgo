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

package lessids

import (
	"errors"
	"sync"
	"time"

	"github.com/lessos/lessgo/net/httpclient"
	"github.com/lessos/lessgo/types"
	"github.com/lessos/lessgo/utils"
	"github.com/lessos/lessgo/utilx"
)

var (
	locker            sync.Mutex
	ServiceUrl        = "http://127.0.0.1:50101/ids"
	sessions          = map[string]UserSession{}
	nextClean         = time.Now()
	innerExpiredRange = time.Second * 1800
)

type UserSession struct {
	types.TypeMeta `json:",inline"`
	AccessToken    string    `json:"access_token"`
	RefreshToken   string    `json:"refresh_token"`
	UserID         string    `json:"userid"`
	UserName       string    `json:"username"`
	ClientAddr     string    `json:"client_addr,omitempty"`
	Name           string    `json:"name"`
	Data           string    `json:"data"`
	Roles          string    `json:"roles"`
	Expired        string    `json:"expired"`
	InnerExpired   time.Time `json:"inner_expired,omitempty"`
	Timezone       string    `json:"timezone"`
}

type UserAccessEntry struct {
	types.TypeMeta `json:",inline"`
	AccessToken    string `json:"access_token"`
	InstanceID     string `json:"instanceid"`
	Privilege      string `json:"privilege"`
}

func innerExpiredClean() {

	if nextClean.After(time.Now()) {
		return
	}

	locker.Lock()
	defer locker.Unlock()

	for k, v := range sessions {

		if v.InnerExpired.Before(time.Now()) {
			continue
		}

		delete(sessions, k)
	}

	nextClean = time.Now().Add(time.Second * 60)
}

func LoginUrl(backurl string) string {
	return ServiceUrl + "/service/login?continue=" + backurl
}

func IsLogin(token string) bool {

	if _, err := SessionFetch(token); err != nil {
		return false
	}

	return true
}

func SessionFetch(token string) (session UserSession, err error) {

	if ServiceUrl == "" || token == "" {
		return session, errors.New("Unauthorized")
	}

	if session, ok := sessions[token]; ok {
		return session, nil
	}

	hc := httpclient.Get(ServiceUrl + "/v1/service/auth?access_token=" + token)

	var us UserSession

	err = hc.ReplyJson(&us)
	if err != nil || us.Error != nil || us.Kind != "UserSession" {
		return session, errors.New("Unauthorized")
	}

	us.InnerExpired = time.Now().Add(innerExpiredRange)

	exp := utilx.TimeParse(us.Expired, "atom")
	if us.InnerExpired.After(exp) {
		us.InnerExpired = exp
	}

	locker.Lock()
	sessions[token] = us // TODO Cache API
	locker.Unlock()

	return us, nil
}

func AccessAllowed(privilege, instanceid, token string) bool {

	if !IsLogin(token) {
		return false
	}

	req := UserAccessEntry{
		AccessToken: token,
		InstanceID:  instanceid,
		Privilege:   privilege,
	}

	js, _ := utils.JsonEncode(req)
	hc := httpclient.Post(ServiceUrl + "/v1/service/access-allowed")
	hc.Header("contentType", "application/json; charset=utf-8")
	hc.Body(js)

	var us UserAccessEntry
	if err := hc.ReplyJson(&us); err != nil || us.Kind != "UserAccessEntry" {
		return false
	}

	return true
}
