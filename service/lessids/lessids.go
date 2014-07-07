package lessids

import (
	"../../net/httpclient"
	"../../utils"
	"errors"
	"sync"
	"time"
)

type ResponseJson struct {
	Status     int    `json:"status"`
	Message    string `json:"message"`
	ApiVersion string `json:"apiVersion"`
}

type SessionEntry struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Uid          uint32    `json:"uid"`
	Uuid         string    `json:"uuid"`
	Uname        string    `json:"uname"`
	Name         string    `json:"name"`
	Data         string    `json:"data"`
	Roles        string    `json:"roles"`
	Expired      time.Time `json:"expired"`
	InnerExpired time.Time
}

var (
	locker            sync.Mutex
	ServiceUrl        = ""
	sessions          = map[string]SessionEntry{}
	nextClean         = time.Now()
	innerExpiredRange = time.Second * 1800
)

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

func SessionFetch(token string) (session SessionEntry, err error) {

	if ServiceUrl == "" || token == "" {
		return session, errors.New("Unauthorized")
	}

	if session, ok := sessions[token]; ok {
		return session, nil
	}

	hc := httpclient.Get(ServiceUrl + "/service/auth?access_token=" + token)

	var rsjson struct {
		ResponseJson
		Data SessionEntry
	}

	err = hc.ReplyJson(&rsjson)
	if err != nil || rsjson.Status != 200 || rsjson.Data.Uid == 0 {
		return session, errors.New("Unauthorized")
	}

	rsjson.Data.InnerExpired = time.Now().Add(innerExpiredRange)

	if rsjson.Data.InnerExpired.After(rsjson.Data.Expired) {
		rsjson.Data.InnerExpired = rsjson.Data.Expired
	}

	locker.Lock()
	sessions[token] = rsjson.Data // TODO Cache API
	locker.Unlock()

	return rsjson.Data, nil
}

func AccessAllowed(privilege, instanceid, token string) bool {

	//fmt.Println("lessids.AccessAllowed", ServiceUrl, privilege, instanceid, token)

	if !IsLogin(token) {
		return false
	}

	var req struct {
		AccessToken string `json:"access_token"`
		Data        struct {
			InstanceId string `json:"instanceid"`
			Privilege  string `json:"privilege"`
		} `json:"data"`
	}
	req.AccessToken = token
	req.Data.InstanceId = instanceid
	req.Data.Privilege = privilege

	reqjson, _ := utils.JsonEncode(req)

	hc := httpclient.Post(ServiceUrl + "/service/access-allowed")
	hc.Header("contentType", "application/json; charset=utf-8")
	hc.Body(reqjson)

	var rsjson ResponseJson
	err := hc.ReplyJson(&rsjson)
	if err != nil || rsjson.Status != 200 {
		return false
	}

	return true
}
