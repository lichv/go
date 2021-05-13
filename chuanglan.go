package lichv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"unsafe"
)

type Chuanglan struct {
	Account string
	Password string
	Sign string
}

func CreateChuanglan(account string,password string,sign string) *Chuanglan{
	return &Chuanglan{Account: account,Password: password,Sign: sign}
}

func (c *Chuanglan) Send(mobile string, content string) (*string,error) {
	params := make(map[string]interface{})
	params["account"] = c.Account
	params["password"] = c.Password
	params["phone"] = mobile
	params["msg"] = url.QueryEscape("【"+c.Sign+"】"+content)
	params["report"] = "true"
	bytesData,err := json.Marshal(params)
	if err != nil {
		fmt.Println(err.Error())
		return nil,err
	}
	reader := bytes.NewReader(bytesData)
	url := "http://smssh1.253.com/msg/send/json"
	request,err := http.NewRequest("POST",url,reader)
	if err != nil {
		fmt.Println(err.Error())
		return nil,err
	}
	client := http.Client{}
	resp,err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	respBytes,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil,err
	}
	str := (*string)(unsafe.Pointer(&respBytes))
	return str ,nil
}