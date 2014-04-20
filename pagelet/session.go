package pagelet

import (
    "../service/lessids"
)

type Session struct {
    AccessToken string
}

func SessionFilter(c *Controller) {

    c.Session = Session{}

    if lessids.ServiceUrl != Config.LessIdsServiceUrl {
        lessids.ServiceUrl = Config.LessIdsServiceUrl
    }

    /* if token := c.Params.Get("setcookie"); token != "" {

        ck := &http.Cookie{
            Name:     "access_token",
            Value:    token,
            Path:     "/",
            HttpOnly: true,
            Expires:  session.Expired.UTC(),
        }
        http.SetCookie(r.Response.Out, ck)
    } */

    if c.Session.AccessToken = c.Params.Get("access_token"); c.Session.AccessToken == "" {

        if token_cookie, err := c.Request.Cookie("access_token"); err == nil {
            c.Session.AccessToken = token_cookie.Value
        }
    }
}

func (s *Session) SessionFetch() (session lessids.SessionEntry, err error) {
    return lessids.SessionFetch(s.AccessToken)
}

func (s *Session) IsLogin() bool {

    println(s)

    session, err := s.SessionFetch()
    if err != nil || session.Uid == 0 {
        return false
    }

    return true
}

func (s *Session) IsAllowed(privilege string) bool {
    return lessids.IsAllowed(privilege, Config.InstanceId, s.AccessToken)
}
