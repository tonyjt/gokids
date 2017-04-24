package gokids

import (
    "time"
)

type ModelPushApiTaskInfo struct {
    Name string `json:"name"`
    StartTime string `json:"startTime"`
    FailureTime string `json:"failureTime"`
    Content string `json:"content"`
    Memo string `json:"memo"`
    SourceSystem string `json:"sourceSystem"`
    AimSystem string `json:"aimSystem"`
    TaskType string `json:"taskType"`
}

type ModelPushApiResponseInfo struct {
    Content uint64 `json:"content"`
    ErrorCode string `json:"errorCode"`
    Meg string `json:"msg"`
    Success bool `json:"success"`
}

type ModelPushMqOcmTaskInfo struct {
    TaskCode int `json:"taskCode"`
    Count int `json:"count"`
    DeviceType int `json:"deviceType"` //设备类型（0表示ios，1表示android，-1表示全部）
    Rows []ModelPushMqOcmTaskRowsInfo `json:"rows"`
}

type ModelPushMqOcmTaskRowsInfo struct {
    CustomerId uint64 `json:"customerId"`
    Badge int `json:"badge"`
    Sound string `json:"sound"`
    Params string `json:"params"`
}

func ModelPushGetApiTaskInfo(taskType string, startTime time.Time, sourceSystem string,
    aimSystem string, content string) (ModelPushApiTaskInfo) {
    info := ModelPushApiTaskInfo{}
    info.Name = "message"
    info.StartTime = UtilTimeGetCommonDateYdmHis(startTime)
    info.FailureTime = UtilTimeGetCommonDateYdmHis(UtilTimeFewDaysLater(startTime, 1))
    info.Content = content
    info.SourceSystem = sourceSystem
    info.AimSystem = aimSystem
    info.TaskType = taskType

    return info
}