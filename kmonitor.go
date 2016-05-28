package gokids

import (
	"os"
	"strings"
	"encoding/json"
	"net/url"
	"time"
  "errors"
  "net/http"
)



var kMonitorHostNanme string

var kmonitor *KMonitor

var kmonitorReportUrl string
var kmonitorKey string

//KMonitorInit init when process start
func KMonitorInit(reportUrl string,productKey string,kmonitorLog ILog){
    kMonitorHostNanme, _ = os.Hostname()
    kmonitorReportUrl = reportUrl
    kmonitorKey = productKey
    log = kmonitorLog
    kmonitor = newKMonitor()

    kmonitor.start()
}

//KMonitorReport report after calling
func KMonitorReport(key string,source string,result bool,duration int64) error{
    if kmonitor ==nil{
      return errors.New("kmonitor not init")
    }

    source = strings.Trim(source, " ")

    key =strings.Trim(key, " ")

    if key == "" {
      return errors.New("parameter key is empty")
    }

    if source !=""{
      key =key + "_"+source
    }

    reqStatus := KMonitorReqStatus{}

    reqStatus.Key = key
    if result{
      reqStatus.Status =1
    }else{
      reqStatus.Status =0
    }
    reqStatus.Duration = duration

    kmonitor.C <- reqStatus

    return nil
}




func newKMonitor() *KMonitor{
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

func (m *KMonitor) add(v KMonitorReqStatus){
  m.Lock()
  defer m.Unlock()

  //status := v[0:1]
  //url := v[2:]
  url := strings.ToLower(strings.TrimSpace(v.Key))
  //fmt.Println("monitor recevier", status, url)

  value, exists := m.data[url]
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
    if v.Status == 1 {
      value.SucessCount = 1
    }
    value.Duration = v.Duration
    // if status == 1 {
    // 	value.SucessCount = 1
    // }
    value.TotalCount = 1
    m.data[url] = value
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

	req := m.getMonitorData()
	m.data = make(map[string]*KMonitorKeyData)
	go m.postToMonitor(req)

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

	//fmt.Println("km post", string(b))

	//res, err := postForms(kmonitorReportUrl, form)
  res, err := http.PostForm(kmonitorReportUrl, data)
	//fmt.Println("post to monito", err, string(res))
	if err != nil {
		if strings.HasSuffix(err.Error(), "EOF") {
			//暂时不记录
			log.Error("monitor error,error:%s,res:%v,data:%v",err,res,string(b))
		} else {
			//srvlog.Error("monitor error", "err", err, "res", res, "data", string(b))
      log.Error("monitor error,error:%s,res:%v,data:%v",err,res,string(b))
		}
	}
}

// 获取监控数据
func (m *KMonitor) getMonitorData() *KMonitorReportReq {
	req := &KMonitorReportReq{}
	//req.T = time.Now().Format(time.RFC3339)
	req.T = time.Now().Format("06-01-02 15:04:05")
	req.P = kmonitorKey
	req.R = make([]KMonitorReportReqData, 0)

	data := m.data
	if data != nil {
		for _, v := range m.data {
			d := KMonitorReportReqData{}
			d.Code = "web" + "@@@@" + kMonitorHostNanme
			d.SKey = v.Url
			d.TotalCount = v.TotalCount
			d.SuccessCount = v.SucessCount
			d.Duration = v.Duration
			req.R = append(req.R, d)
		}
	}
	return req
}
