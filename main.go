package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Info struct {
	Access_token string
	Expires_in   string
}

func main() {

	response, _ := http.Get("http://myexternalip.com/raw")
	defer response.Body.Close()
	IP, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(IP))
	sendWXQY("获取电脑IP:" + string(IP))

}

func sendWXQY(content string) {
	resp1, _ := http.Get("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=wxfe9025ee4d2c1ca2&corpsecret=_AX6ICTZegQADd9JaQWp3xXhTGw4cDBZYF-7Q7yAt574jD3rgWHft24oK3hG1p_U")

	defer resp1.Body.Close()
	body1, _ := ioutil.ReadAll(resp1.Body)
	JsonStr := string(body1)
	var myInfo Info
	json.Unmarshal([]byte(JsonStr), &myInfo)
	Access_token := myInfo.Access_token
	fmt.Println(Access_token)

	tmp := `{"touser": "dcy","msgtype": "text","agentid": 0,"text": {"content": "` + content + `"},"safe":0}`
	req := bytes.NewBuffer([]byte(tmp))

	body_type := "application/json;charset=utf-8"
	resp, _ := http.Post("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token="+Access_token, body_type, req)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
