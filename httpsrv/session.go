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
	"github.com/lessos/lessgo/service/lessids"
)

type Session struct {
	InstanceID  string
	AccessToken string
}

func SessionFilter(c *Controller) {

	c.Session = Session{
		InstanceID: c.service.Config.InstanceID,
	}

	if lessids.ServiceUrl != c.service.Config.LessIdsServiceUrl {
		lessids.ServiceUrl = c.service.Config.LessIdsServiceUrl
	}

	/* if token := c.Params.Get("setcookie"); token != "" {

	    ck := &http.Cookie{
	        Name:     c.service.Config.CookieKeySession,
	        Value:    token,
	        Path:     "/",
	        HttpOnly: true,
	        Expires:  session.Expired.UTC(),
	    }
	    http.SetCookie(r.Response.Out, ck)
	} */

	if c.Session.AccessToken = c.Params.Get(c.service.Config.CookieKeySession); c.Session.AccessToken == "" {

		if token_cookie, err := c.Request.Cookie(c.service.Config.CookieKeySession); err == nil {
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
	return lessids.AccessAllowed(privilege, s.InstanceID, s.AccessToken)
}
