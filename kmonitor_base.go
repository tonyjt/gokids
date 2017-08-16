package gokids

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	kMonitorHostNanme string
	kmonitor          *KMonitor
	kmonitorReportUrl string
	kmonitorKey       string
)

//Base Kmonitor methods
//KMonitorBaseInit init when process start
func KMonitorBaseInit(reportUrl string, kmonitorLog ILog) {
	kMonitorHostNanme, _ = os.Hostname()
	kmonitorReportUrl = reportUrl
	log = kmonitorLog
	kmonitor = newKMonitor()
	kmonitor.start()
}

//KMonitorBaseReport report after calling
func KMonitorBaseReport(productKey string, key string, source string, result bool, duration int64) error {
	if kmonitor == nil {
		return errors.New("kmonitor not init")
	}

	source = strings.Trim(source, " ")

	key = strings.Trim(key, " ")

	if key == "" {
		return errors.New("parameter key is empty")
	}

	if source != "" {
		key = key + "_" + source
	}

	reqStatus := KMonitorReqStatus{}

	reqStatus.ProductKey = productKey
	reqStatus.Key = key
	if result {
		reqStatus.Status = 1
	} else {
		reqStatus.Status = 0
	}
	reqStatus.Duration = duration

	kmonitor.C <- reqStatus

	return nil
}

func newKMonitor() *KMonitor {
	m := &KMonitor{}
	m.data = make(map[string]*KMonitorKeyData)
	//m.C = make(chan string, 100000)
	m.C = make(chan KMonitorReqStatus)
	m.ticker = time.NewTicker(time.Second * 60)
	return m
}

// 启动
func (m *KMonitor) start() {
	go m.startChan()
	go m.startTicker()
}

func (m *KMonitor) add(v KMonitorReqStatus) {
	m.Lock()
	defer m.Unlock()

	//status := v[0:1]
	//url := v[2:]
	url := strings.ToLower(strings.TrimSpace(v.Key))
	dataKey := fmt.Sprintf("%s%s", v.ProductKey, url)
	//fmt.Println("monitor recevier", status, url)

	value, exists := m.data[dataKey]
	if exists {
		//if status == 1 {
		//atomic.AddInt64(&value.sucessCount, 1)
		//}
		//atomic.AddInt64(&value.totalCount, 1)
		if v.Status == 1 {
			value.SucessCount++
		}
		value.TotalCount++
		value.Duration = value.Duration + v.Duration
	} else {
		value = &KMonitorKeyData{}
		value.Url = url
		value.ProductKey = v.ProductKey
		if v.Status == 1 {
			value.SucessCount = 1
		}
		value.Duration = v.Duration
		// if status == 1 {
		// 	value.SucessCount = 1
		// }
		value.TotalCount = 1
		m.data[dataKey] = value
	}
}

// 启动ticker
func (m *KMonitor) startTicker() {
	for range m.ticker.C {
		//fmt.Println("kmonitor ticket", now)
		m.report()
	}
}

func (m *KMonitor) report() {
	m.Lock()
	defer m.Unlock()

	reqMaps := m.getMonitorData()
	m.data = make(map[string]*KMonitorKeyData)
	for _, req := range reqMaps {
		go m.postToMonitor(req)
	}

}

// 启动接收
func (m *KMonitor) startChan() {
	for v := range m.C {
		m.add(v)
	}
}

//上报格式优化
func (m *KMonitor) postToMonitor(req *KMonitorReportReq) {
	b, _ := json.Marshal(req)
	//form := neturl.Values{}
	data := url.Values{}
	data.Add("c", "0")
	data.Add("d", "["+string(b)+"]")

	//fmt.Println(kmonitorReportUrl, string(b))

	//res, err := postForms(kmonitorReportUrl, form)
	_, err := http.PostForm(kmonitorReportUrl, data)
	//fmt.Println("post to monito", err, string(res))
	if err != nil {
		if strings.HasSuffix(err.Error(), "EOF") {
			//暂时不记录
			//log.Error("monitor error,error:%s,res:%v,data:%v",err,res,string(b))
		} else {
			//srvlog.Error("monitor error", "err", err, "res", res, "data", string(b))
			//log.Error("monitor error,error:%s,res:%v,data:%v",err,res,string(b))
		}
	}
}

// 获取监控数据
func (m *KMonitor) getMonitorData() map[string]*KMonitorReportReq {
	reqMaps := make(map[string]*KMonitorReportReq)
	data := m.data
	if data != nil {
		for _, v := range m.data {
			var req *KMonitorReportReq
			key := v.ProductKey
			if _, ok := reqMaps[key]; ok {
				req = reqMaps[key]
			} else {
				req = &KMonitorReportReq{}
				req.T = time.Now().Format("06-01-02 15:04:05")
				req.P = v.ProductKey
				req.R = make([]KMonitorReportReqData, 0)
			}

			d := KMonitorReportReqData{}
			d.Code = "web" + "@@@@" + kMonitorHostNanme
			d.SKey = v.Url
			d.TotalCount = v.TotalCount
			d.SuccessCount = v.SucessCount
			d.Duration = v.Duration
			req.R = append(req.R, d)
			reqMaps[key] = req
		}
	}
	return reqMaps
}
