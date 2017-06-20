package gokids

import "testing"

func TestKucodeVerify(t *testing.T) {
	Init(NewLogDefault())

	u := "http://test.verifycode.haiziwang.com/ucode-web/ucode/check.do"

	identity := "1234"
	appCode := "xxr-membergo"
	appServiceCode := "kidslogin"
	vc := "3952"

	result, err := KucodeVerify(identity, appCode, appServiceCode, vc, u)

	if err != nil {
		t.Errorf("err is :%s", err.Error())
	} else if !result {
		t.Error("result is false")
	}

}
