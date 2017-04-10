package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
)

type Info struct {
	Access_token string
	Expires_in   string
}

func main() {
	msg := "获取电脑外网IP:" + get_external() + "内网IP:" + get_internal() + "\n" + "主机名：" + getHostname()
	fmt.Println(msg)
	sendWXQY("dcy", msg)
}

//获取电脑外网IP
func get_external() string {
	response, _ := http.Get("http://myexternalip.com/raw")
	defer response.Body.Close()
	IP, _ := ioutil.ReadAll(response.Body)
	//fmt.Println(string(IP))
	return string(IP)
}

//获取电脑内网IP
func get_internal() string {
	ips := ""
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops:" + err.Error())
		//os.Exit(1)
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				//os.Stdout.WriteString(ipnet.IP.String() + "\n")
				if ips == "" {
					ips = ipnet.IP.String()
				} else {
					ips = ips + "/" + ipnet.IP.String()
				}
			}
		}
	}
	//os.Exit(0)
	return ips
}

//获取电脑主机名
func getHostname() string {
	host, err := os.Hostname()
	if err != nil {
		fmt.Printf("%s", err)
		return ""
	} else {
		return host
	}
}

//发送信息给微信用户
func sendWXQY(wxid string, content string) {
	resp1, _ := http.Get("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=wxfe9025ee4d2c1ca2&corpsecret=9CkdgSgjatw04xEf_yz_CIZl41IvoBJFCioLshBRbj_Dks6Z12W34gby44YPTxMw")
	defer resp1.Body.Close()
	body1, _ := ioutil.ReadAll(resp1.Body)
	JsonStr := string(body1)
	var myInfo Info
	json.Unmarshal([]byte(JsonStr), &myInfo)
	Access_token := myInfo.Access_token
	fmt.Println(Access_token)

	tmp := `{"touser": "` + wxid + `","msgtype": "text","agentid": 0,"text": {"content": "` + content + `"},"safe":0}`
	req := bytes.NewBuffer([]byte(tmp))

	body_type := "application/json;charset=utf-8"
	resp, _ := http.Post("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token="+Access_token, body_type, req)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
