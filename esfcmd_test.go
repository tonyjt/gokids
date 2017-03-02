package gokids

import (
	"testing"
	"time"
)


func Test_EsfCmdInit(t *testing.T) {
	url := "http://esf.haiziwang.com/esf-web/cpp/cppServiceList.do"
	//url := "http://test.esf.haiziwang.com:9090/esf-web/cpp/cppServiceList.do"

	var c []string

	c = append(c,"4006")

	EsfCmdInit(url,c,NewLogDefault(),5 * time.Second)
}

func Test_EsfCmdGetAddressByCmd(t *testing.T){
	Test_EsfCmdInit(t)

	a := EsfCmdGetAddr(0x40061802,"172.172.200.28:53101")

	if len(a) == 0{
		t.Errorf("addr is empty")
	}
	log.Info("%v",a)
}