package gokids

import (
    "time"
)

type ModelMessageMsgBoxInfo struct {
    CustomerIds string `json:"customerIds"`
    MsgType uint8 `json:"msgType"`
    SourceId string `json:"sourceId"`
    CreateTime string `json:"createTime"`
    Content string `json:"content"`
    AppCode string `json:"appCode"`
}

func ModelMessageNewMsgInfo(msgType uint8, sourceId string, appCode string, content string) (ModelMessageMsgBoxInfo) {
    data := ModelMessageMsgBoxInfo{}
    data.MsgType = msgType
    data.SourceId = sourceId
    data.CreateTime = UtilTimeGetCommonDateYdmHis(time.Now())
    data.AppCode = appCode
    data.Content = content

    return data
}

func ModelMessageNewSqOperateMsgInfo(content string) (ModelMessageMsgBoxInfo) {
    data := ModelMessageMsgBoxInfo{}
    data.MsgType = 6
    data.SourceId = "SNS"
    data.CreateTime = UtilTimeGetCommonDateYdmHis(time.Now())
    data.AppCode = "HZW_MALL"
    data.Content = content

    return data
}

func ModelMessageNewSqQuestionMsgInfo(content string) (ModelMessageMsgBoxInfo) {
    data := ModelMessageMsgBoxInfo{}
    data.MsgType = 7
    data.SourceId = "SNS"
    data.CreateTime = UtilTimeGetCommonDateYdmHis(time.Now())
    data.AppCode = "rkhy"
    data.Content = content

    return data
}