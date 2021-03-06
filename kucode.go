package gokids

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// KucodeVerify 图片验证码验证
func KucodeVerify(identity string, appCode string, appServiceCode string, verifyCode string, u string) (result bool, err error) {
	if identity == "" || appCode == "" || appServiceCode == "" || u == "" {
		err = errors.New("parameter issue")
		return
	}

	if verifyCode == "" {
		result = false
		return
	}
	base64Identity := base64.StdEncoding.EncodeToString([]byte(identity))

	base64AppCode := base64.StdEncoding.EncodeToString([]byte(appCode))

	base64AppServiceCode := base64.StdEncoding.EncodeToString([]byte(appServiceCode))

	base64VerifyCode := base64.StdEncoding.EncodeToString([]byte(verifyCode))

	url := fmt.Sprintf("%s?identity=%s&appCode=%s&appServiceCode=%s&verifyCode=%s", u, base64Identity, base64AppCode, base64AppServiceCode, base64VerifyCode)

	res, err := http.Get(url)

	if err != nil {
		if strings.HasSuffix(err.Error(), "EOF") {

		} else {
			log.Error("get url error, err:%s", err.Error())
		}
		return
	}

	b, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Error("read all failed, err:%s", err.Error())

	} else {
		modelRes := KucodeModelVerifyResponse{}
		err = json.Unmarshal(b, &modelRes)
		log.Info(string(b))
		if err != nil {
			log.Error("unmarshal failed, err:%s", err.Error())
		} else {
			if modelRes.Success {
				result = true
			} else {
				result = false
			}
		}
	}

	return

}
