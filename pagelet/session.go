package pagelet

import (
	"github.com/lessos/lessgo/service/lessids"
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
	        Name:     Config.SessionCookieKey,
	        Value:    token,
	        Path:     "/",
	        HttpOnly: true,
	        Expires:  session.Expired.UTC(),
	    }
	    http.SetCookie(r.Response.Out, ck)
	} */

	if c.Session.AccessToken = c.Params.Get(Config.SessionCookieKey); c.Session.AccessToken == "" {

		if token_cookie, err := c.Request.Cookie(Config.SessionCookieKey); err == nil {
			c.Session.AccessToken = token_cookie.Value
		}
	}
}

func (s *Session) SessionFetch() (session lessids.SessionEntry, err error) {
	return lessids.SessionFetch(s.AccessToken)
}

func (s *Session) IsLogin() bool {

	session, err := s.SessionFetch()
	if err != nil || session.Uid == 0 {
		return false
	}

	return true
}

func (s *Session) AccessAllowed(privilege string) bool {
	return lessids.AccessAllowed(privilege, Config.InstanceId, s.AccessToken)
}
