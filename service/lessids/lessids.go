package lessids

import (
    "../../net/httpclient"
    "../../pagelet"
    "sync"
    "time"
)

type ResponseJson struct {
    Status     int    `json:"status"`
    Message    string `json:"message"`
    ApiVersion string `json:"apiVersion"`
}

type Session struct {
    AccessToken  string    `json:"access_token"`
    RefreshToken string    `json:"refresh_token"`
    Uid          uint32    `json:"uid"`
    Uname        string    `json:"uname"`
    Name         string    `json:"name"`
    Data         string    `json:"data"`
    Expired      time.Time `json:"expired"`
    InnerExpired time.Time
}

var (
    locker            sync.Mutex
    ServiceUrl        = ""
    sessions          = map[string]Session{}
    nextClean         = time.Now()
    innerExpiredRange = time.Second * 3600
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

func IsLogin(r *pagelet.Request) bool {

    if ServiceUrl == "" {
        return false
    }

    cookie, err := r.Request.Cookie("access_token")
    if err == nil {

        if _, ok := sessions[cookie.Value]; ok {

            innerExpiredClean()
            return true
        }
    }

    hc := httpclient.Get(ServiceUrl + "/service/auth?access_token=" + cookie.Value)

    var rsjson struct {
        ResponseJson
        Data Session
    }

    err = hc.ReplyJson(&rsjson)
    if err != nil || rsjson.Status != 200 || rsjson.Data.Uid == 0 {
        return false
    }

    rsjson.Data.InnerExpired = time.Now().Add(innerExpiredRange)

    if rsjson.Data.InnerExpired.After(rsjson.Data.Expired) {
        rsjson.Data.InnerExpired = rsjson.Data.Expired
    }

    locker.Lock()
    // TODO Cache API
    sessions[cookie.Value] = rsjson.Data
    locker.Unlock()

    return true
}
