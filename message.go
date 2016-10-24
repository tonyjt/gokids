package gokids

import (
    "errors"
    "strings"
)

var (
    MessageMsgBoxChannel chan ModelMessageMsgBoxChannel
    MessageMsgBoxLimitNum uint8 //失败重试次数
)

//异步使用必须启动时运行
func RunMessageMsgBox() {
    MessageMsgBoxChannel = make(chan ModelMessageMsgBoxChannel)
    MessageMsgBoxLimitNum = 5
    go func() {
        for {
            select {
            case msg := <-MessageMsgBoxChannel:
                go dealMessageSendMqData(msg)
            }
        }
    }()
}

func dealMessageSendMqData(msg ModelMessageMsgBoxChannel) {
    err := MessageSendMqData(msg.MsgBoxUrl, msg.MqData)
    if err != nil {
        if msg.LimitNum < MessageMsgBoxLimitNum {
            msg.LimitNum++
            MessageMsgBoxChannel <- msg
        } else {
            log.Error("message send msg box info error :%s, data: %v", err.Error(), msg)
        }
    }
    return
}

func MessageAddMsgInfo(msgBoxUrl string, msgType uint8, sourceId string, appCode string, content string,
    uids []string) (error) {
    if len(uids) > 200 {
        return errors.New("uids num out of range")
    }
    mqData := ModelMessageNewMsgInfo(msgType, sourceId, appCode, content)
    mqData.CustomerIds = strings.Join(uids, ",")
    return MessageSendMqData(msgBoxUrl, mqData)
}

func MessageAddMsgInfoAsync(msgBoxUrl string, msgType uint8, sourceId string, appCode string, content string,
    uids []string) {
    if len(uids) > 200 {
        uids = uids[:200]
    }
    mqData := ModelMessageNewMsgInfo(msgType, sourceId, appCode, content)
    mqData.CustomerIds = strings.Join(uids, ",")
    data := ModelMessageMsgBoxChannel{msgBoxUrl, mqData, 0}
    go func(msgBoxChannel chan ModelMessageMsgBoxChannel, mData ModelMessageMsgBoxChannel) {
        msgBoxChannel <- mData
    }(MessageMsgBoxChannel, data)
}

func MessageSendMqData(msgBoxUrl string, data ModelMessageMsgBoxInfo) (error) {
    key := "hzw.msgbox.msgList"//写死
    return UtilMQSend(msgBoxUrl, key, data)
}