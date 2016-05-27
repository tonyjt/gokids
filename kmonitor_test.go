package gokids

import (
  "testing"
)

func Test_Kmonitor(t *testing.T) {
  reportUrl :="http://172.172.177.23:80/monitor-web/statistic/report.do"

  productKey := "mt"

  kmonitorLog := NewLogDefault()

  KMonitorInit(reportUrl, productKey, kmonitorLog)

  key := "key1"
  source := "test"
  result := true
  duration := int64(200)
  KMonitorReport(key, source, result, duration)
}
