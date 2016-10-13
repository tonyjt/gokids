package gokids

import (
    "encoding/json"
    "github.com/go-stomp/stomp"
    "strings"
    "errors"
)

var (
    ActiveUrl map[string]string
)

func UtilMQSend(url string, key string, data interface{}) error {
    var (
        urlMqActiveArray []string
        err error
        mqUrl string
    )
    if url == "" {
        return errors.New("url is none")
    }
    if ActiveUrl == nil {
        ActiveUrl = make(map[string]string)
    }
    urlMqActiveArray = strings.Split(url, ",")
    md5Str := UtilCryptoGenerateMD5Hash(url)
    mqUrl, ok := ActiveUrl[md5Str]
    if !ok {
        mqUrl = urlMqActiveArray[0]
    }
    conn, errConn := utilMQDial(mqUrl)
    if errConn != nil {
        for _, mu := range urlMqActiveArray {

            if mu != mqUrl {
                conn, errConn = utilMQDial(mu)
                if errConn == nil {
                    ActiveUrl[md5Str] = mu
                    break
                }
            }
        }

        if errConn != nil {
            log.Error("mq url connect error,sever:%s,key:%s,data:%v,error:%s",
                url, key, data, errConn.Error())
            return errConn
        }
    } else {
        ActiveUrl[md5Str] = mqUrl
    }

    defer conn.Disconnect()

    var msg []byte
    msg, err = json.Marshal(data)
    if err != nil {
        return err
    }
    err = conn.Send(key, "text/plain", msg)

    if err != nil {
        log.Error("mq send msg error,sever:%s,key:%s,data:%v,error:%s",
            mqUrl, key, data, err.Error())
    }
    return err
}

func utilMQDial(url string) (*stomp.Conn, error) {
    conn, err := stomp.Dial("tcp", url)
    if err != nil {
        log.Error("cannot connect to server:%s,error:%s", url, err.Error())
    }
    return conn, err
}
