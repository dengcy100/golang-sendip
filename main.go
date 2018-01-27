package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
)

type Info struct {
	Access_token string
	Expires_in   string
}

func main() {
	msg := "获取电脑外网IP:" + get_external() + "\n" + "内网IP:" + get_internal() + "\n" + "主机名：" + getHostname()
	fmt.Println(msg)
	sendWXQY("dcy", msg)
	paiming := getgongzu()
	if strings.Contains(paiming, "83932") == false {
		sendWXQY("dcy", paiming)
	}
}

//获取电脑外网IP
func get_external() string {
	response, _ := http.Get("http://www.ip.cn/")
	defer response.Body.Close()
	IP, _ := ioutil.ReadAll(response.Body)
	fmt.Println(strings.Split(strings.Split(string(IP), "<code>")[1], "</code>")[0])
	return string(strings.Split(strings.Split(string(IP), "<code>")[1], "</code>")[0])
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
	corpid := "wxfe9025ee4d2c1ca2"
	corpsecret := "el-2BwXsf2dHt9kxucDOlfruFJk7uVpHEcmx2lZY8jSIeksLuMvVMUL1b49KiAM2"
	resp1, _ := http.Get("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=" + corpid + "&corpsecret=" + corpsecret)
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

//获取公租房排名
func getgongzu() string {
	response, _ := http.Get("http://bzflh.szjs.gov.cn/TylhW/lhmcAction.do?method=queryLhmcInfo1&pageNumber=1&pageSize=10&waitTpye=2&bahzh=BHR00088615&xingm=%25E9%2582%2593%25E5%25B4%2587%25E6%2584%25BF&sfz=441427198405160813")
	defer response.Body.Close()
	content, _ := ioutil.ReadAll(response.Body)
	jsonStr := string(content)
	//fmt.Println(string(content))

	AREA_PAIX := strings.Split(jsonStr, "AREA_PAIX")[1]
	PAIX := strings.Split(jsonStr, "\"PAIX")[1]
	//fmt.Println(PAIX)
	qupaiming := "公租房区排名：" + strings.Split(strings.Split(AREA_PAIX, "\"")[1], ",")[0]
	shipaiming := "市排名：" + strings.Split(strings.Split(PAIX, "\"")[1], ",")[0]
	fmt.Println(qupaiming)
	fmt.Println(shipaiming)
	return qupaiming + "\n" + shipaiming
}
