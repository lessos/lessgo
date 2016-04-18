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

type Config struct {
	HttpAddr         string `json:"http_addr,omitempty"` // e.g. "127.0.0.1", "unix:/tmp/app.sock"
	HttpPort         uint16 `json:"http_port,omitempty"` // e.g. 8080
	HttpTimeout      uint16 `json:"http_timeout,omitempty"`
	UrlBasePath      string `json:"url_base_path,omitempty"`
	CookieKeyLocale  string `json:"cookie_key_locale,omitempty"`
	CookieKeySession string `json:"cookie_key_session,omitempty"`
}

func (c *Config) TemplateFuncRegister(name string, fn interface{}) {
	TemplateFuncs[name] = fn
}

func (c *Config) I18n(file string) {
	i18nLoadMessages(file)
}
