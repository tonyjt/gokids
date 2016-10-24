package gokids

import (
    "errors"
    "strings"
)

func MessageAddMsgInfo(msgBoxUrl string, msgType uint8, sourceId string, appCode string, content string,
    uids []string) (error) {
    if len(uids) > 200 {
        return errors.New("uids num out of range")
    }
    mqData := ModelMessageNewMsgInfo(msgType, sourceId, appCode, content)
    mqData.CustomerIds = strings.Join(uids, ",")
    return MessageSendMqData(msgBoxUrl, mqData)
}

func MessageSendMqData(msgBoxUrl string, data ModelMessageMsgBoxInfo) (error) {
    key := "hzw.msgbox.msgList"//写死
    return UtilMQSend(msgBoxUrl, key, data)
}