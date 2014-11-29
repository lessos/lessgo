package pagelet

import (
	"errors"
	"fmt"
	"github.com/lessos/lessgo/utils"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var (
	i18n      map[string]string = map[string]string{}
	sysLocale string            = "en"
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

	if v, e := c.Request.Cookie(Config.LocaleCookieKey); e == nil {
		c.Request.Locale = v.Value
	} else if len(c.Request.AcceptLanguage) > 0 {
		c.Request.Locale = c.Request.AcceptLanguage[0].Language
	} else {
		c.Request.Locale = sysLocale
	}

	c.ViewData["LANG"] = c.Request.Locale
}

func i18nLoadMessages(file string) {

	var cfg i18nConfig

	str, err := i18nFsFileGetRead(file)
	if err != nil {
		return
	}

	if err := utils.JsonDecode(str, &cfg); err != nil {
		fmt.Println("Format Error:", err)
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
	keydef := strings.ToLower(sysLocale + "." + msg)

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

	reg, _ := regexp.Compile("/+")
	path = "/" + strings.Trim(reg.ReplaceAllString(path, "/"), "/")

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
