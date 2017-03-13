package gokids

import (
	"net/http"
	"strings"
	"io/ioutil"
	"encoding/json"
	"time"
	"fmt"
	"math/rand"
	"sync"
)

var (
	esfAddr map[string][]EsfCmdAddr
	esfAddrM sync.RWMutex
	esfUrl string
	esfTicker *time.Ticker
)
//EsfCmdInit 初始化
func EsfCmdInit(esfGetUrl string,cmdList []string,esfCmdLog ILog, updateDuration time.Duration){

	if log ==nil{
		log = esfCmdLog
	}

	if esfAddr == nil{
		esfAddr = make(map[string][]EsfCmdAddr)
	}

	esfUrl = esfGetUrl

	for _,cmd := range cmdList{
		addrs := esfCmdGetAddressByCmd(cmd)
		if len(addrs) >0{
			esfAddr[cmd] = addrs
		}
	}



	if esfTicker ==nil{
		esfTicker = time.NewTicker(updateDuration)
		go esfStartTicker()
	}

}


//EsfCmdGetAddr 获取命令字服务器地址
func EsfCmdGetAddr(cmd32 uint32, defaultAddr string) (strAddr []string){
	cmd := fmt.Sprintf("%x",cmd32 / 0xffff)

	ars,ok := esfAddr[cmd]

	if !ok || len(ars) ==0{
		ars = esfCmdGetAddressByCmd(cmd)
		esfAddrM.Lock()
		esfAddr[cmd] = ars
		esfAddrM.Unlock()
	}

	for _,a := range ars{
		if a.Enablestatus ==1{
			strA := fmt.Sprintf("%s:%d",a.Address,a.Port)
			strAddr = append(strAddr,strA)
		}
	}

	if len(strAddr) ==0 {
		strAddr = append(strAddr,defaultAddr)
	}
	return strAddr
}
//EsfCmdGetAddrRandom 获取命令字服务器地址，随机取一个
func EsfCmdGetAddrRandom(cmd uint32,defaultAddr string)(strAddr string){
	a := EsfCmdGetAddr(cmd,defaultAddr)

	if len(a) ==0{
		strAddr = defaultAddr
	}else if len(a) ==1{
		strAddr = a[0]
	}else{
		rand.Seed(time.Now().UnixNano())

		i := rand.Intn(len(a))

		strAddr = a[i]
	}
	return
}

func esfCmdGetAddressByCmd(cmd string)(addrs []EsfCmdAddr){

	url :=fmt.Sprintf("%s?reportload=0&aoData=[{\"name\":\"theme\",\"value\":\"%s\"}]",esfUrl,cmd)

	res,err:=http.Get(url)

	if err!=nil{
		if strings.HasSuffix(err.Error(), "EOF") {

		}else{
			log.Error("get url error, err:%s",err.Error())
		}
		return
	}

	b,errRead:= ioutil.ReadAll(res.Body)

	if errRead!=nil{
		log.Error("read all failed, err:%s",errRead.Error())

	}else{
		modelRes:= EsfCmdResponse{}
		errJson:= json.Unmarshal(b, &modelRes)

		if errJson!=nil{
			log.Error("unmarshal failed, err:%s",errJson.Error())
		}else{
			if modelRes.Result && modelRes.Data.ITotalRecords >0 {
				for _,a:= range modelRes.Data.AaData{
					if a.Enablestatus ==1{
						addrs = append(addrs,a)
					}
				}
			}else{
				log.Error("get url error, data :%v",modelRes)
			}
		}
	}

	return addrs
}


func esfStartTicker(){
	for range esfTicker.C{
		for cmd,_ := range esfAddr{
			addrs := esfCmdGetAddressByCmd(cmd)
			if len(addrs) >0{
				esfAddrM.Lock()
				esfAddr[cmd] = addrs
				esfAddrM.Unlock()
			}
		}
	}
}

