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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/lessos/lessgo/logger"
	"github.com/lessos/lessgo/utils"
)

var (
	i18n          = map[string]string{}
	i18nDefLocale = "en"
	i18nRegPath   = regexp.MustCompile("/+")
)

type i18nConfig struct {
	Locale string `json:"locale"`
	Data   []i18nConfigItem
}

type i18nConfigItem struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

func I18nFilter(c *Controller) {

	if v, e := c.Request.Cookie(c.service.Config.CookieKeyLocale); e == nil {
		c.Request.Locale = v.Value
	} else if len(c.Request.AcceptLanguage) > 0 {
		c.Request.Locale = c.Request.AcceptLanguage[0].Language
	} else {
		c.Request.Locale = i18nDefLocale
	}

	c.Data["LANG"] = c.Request.Locale
}

func i18nLoadMessages(file string) {

	var cfg i18nConfig

	str, err := i18nFsFileGetRead(file)
	if err != nil {
		return
	}

	if err := utils.JsonDecode(str, &cfg); err != nil {
		logger.Print("error", "Format Error: "+err.Error())
		return
	}

	cfg.Locale = strings.Replace(cfg.Locale, "_", "-", 1)

	for _, v := range cfg.Data {

		key := strings.ToLower(cfg.Locale + "." + v.Key)

		if v2, ok := i18n[key]; !ok || v2 != v.Val {
			i18n[key] = v.Val
		}
	}
}

func i18nTranslate(locale, msg string, args ...interface{}) string {

	key := strings.ToLower(locale + "." + msg)
	keydef := strings.ToLower(i18nDefLocale + "." + msg)

	if v, ok := i18n[key]; ok {
		msg = v
	} else if v, ok := i18n[keydef]; ok {
		msg = v
	}

	if len(args) > 0 {
		return fmt.Sprintf(msg, args...)
	} else {
		return msg
	}
}

func i18nFsFileGetRead(path string) (string, error) {

	path = "/" + strings.Trim(i18nRegPath.ReplaceAllString(path, "/"), "/")

	st, err := os.Stat(path)
	if err != nil || os.IsNotExist(err) {
		return "", errors.New("File Not Found")
	}

	if st.Size() > (10 * 1024 * 1024) {
		return "", errors.New("File size is too large")
	}

	fp, err := os.OpenFile(path, os.O_RDWR, 0754)
	if err != nil {
		return "", errors.New("File Can Not Open")
	}
	defer fp.Close()

	ctn, err := ioutil.ReadAll(fp)
	if err != nil {
		return "", errors.New("File Can Not Readable")
	}

	return string(ctn), nil
}
