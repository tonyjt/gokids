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

var (
    kMonitorHostNanme string
    kmonitor *KMonitor
    kmonitorReportUrl string
    kmonitorKey string
    KMonitorStartTimeKey string
)

func init() {
    KMonitorStartTimeKey = "kMonitorStartTime"
}

//KMonitorInit init when process start
func KMonitorInit(reportUrl string, productKey string, kmonitorLog ILog) {
    kMonitorHostNanme, _ = os.Hostname()
    kmonitorReportUrl = reportUrl
    kmonitorKey = productKey
    log = kmonitorLog
    kmonitor = newKMonitor()

    kmonitor.start()
}

//KMonitorReport report after calling
func KMonitorReport(key string, source string, result bool, duration int64) error {
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
    urlPath := strings.ToLower(strings.TrimSpace(v.Key))
    //fmt.Println("monitor recevier", status, url)

    value, exists := m.data[urlPath]
    if !exists {
        value = &KMonitorKeyData{}
    }
    if v.Status == 1 {
        value.SuccessCount++
    }
    value.TotalCount++
    value.Duration += v.Duration
    m.data[urlPath] = value
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
    //fmt.Println(m.data)
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
    //fmt.Println("post to req", req)
    if len(req.R) == 0 {
        return
    }
    b, _ := json.Marshal(req)
    //form := neturl.Values{}
    data := url.Values{}
    data.Add("c", "0")
    data.Add("d", "[" + string(b) + "]")

    //fmt.Println("km post", string(b))

    res, err := http.PostForm(kmonitorReportUrl, data)
    //fmt.Println("post to monito", err, res)
    if err != nil {
        if strings.HasSuffix(err.Error(), "EOF") {
            //暂时不记录
            //log.Error("monitor error,error:%s,res:%v,data:%v",err,res,string(b))
        } else {
            //srvlog.Error("monitor error", "err", err, "res", res, "data", string(b))
            log.Error("monitor error,error:%s,res:%v,data:%v", err, res, string(b))
        }
    }
}

// 获取监控数据
func (m *KMonitor) getMonitorData() *KMonitorReportReq {
    req := &KMonitorReportReq{}
    //req.T = time.Now().Format(time.RFC3339)
    req.T = time.Now().Format("06-01-02 15:04:05")
    req.P = kmonitorKey

    for k, v := range m.data {
        d := KMonitorReportReqData{}
        d.Code = "web" + "@@@@" + kMonitorHostNanme
        d.SKey = k
        d.TotalCount = v.TotalCount
        d.SuccessCount = v.SuccessCount
        d.Duration = v.Duration
        delete(m.data, k)
        req.R = append(req.R, d)
    }

    return req
}
