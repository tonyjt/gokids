package gokids

//KMonitorInit init when process start
func KMonitorInit(reportUrl string,productKey string,kmonitorLog ILog){
	kmonitorKey = productKey
	KMonitorBaseInit(reportUrl, kmonitorLog)
}

func KMonitorReport(key string,source string,result bool,duration int64) error {
	return KMonitorBaseReport(kmonitorKey, key, source, result, duration)
}
