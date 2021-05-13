package lichv

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"unsafe"
)

func post(url string, params map[string]interface{},application_type string) (*string, error) {
	if len(params) == 0 {
		return nil, errors.New("参数错误")
	}
	if application_type == ""{
		application_type = "application/x-www-form-urlencoded"
	}
	var result []string
	for k,v := range params{
		result = append(result,fmt.Sprintf("%s=%s",k,v))
	}
	inputStr := strings.Join(result,"&")
	resp, err := http.Post(url, application_type, strings.NewReader(inputStr))
	if err != nil {
		return nil, err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	str := (*string)(unsafe.Pointer(&respBytes))
	return str, nil
}

func get(url string)(*string,error) {
	resp,err := http.Get(url)
	if err != nil {
		return nil, err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	str := (*string)(unsafe.Pointer(&respBytes))
	return str, nil
}

