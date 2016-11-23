package gokids

import (
    "net/http"
    "strings"
    "io/ioutil"
)

func UtilCurlGet(url string, header []string) (ret []byte, err error) {
    client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return ret, err
    }
    for _, v := range header {
        t := strings.Split(v, ":")
        length := len(t)
        if length == 2 {
            req.Header.Add(t[0], t[1])
        } else if length == 1 {
            req.Header.Add(t[0], "")
        }
    }

    resp, err := client.Do(req)
    if err != nil {
        return ret, err
    }
    defer resp.Body.Close()
    ret, err = ioutil.ReadAll(resp.Body)

    return ret, err
}
