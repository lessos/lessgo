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
	InstanceID        string
	HttpPort          uint16 // e.g. 8080
	HttpAddr          string // e.g. "", "127.0.0.1", "unix:/tmp/app.socket (if the port is zero)"
	UrlBasePath       string
	CookieKeyLocale   string
	CookieKeySession  string
	LessIdsServiceUrl string
}

func (c *Config) TemplateFuncRegister(name string, fn interface{}) {

	if _, ok := templateFuncs[name]; ok {
		return
	}

	templateFuncs[name] = fn
}

func (c *Config) I18n(file string) {
	i18nLoadMessages(file)
}
