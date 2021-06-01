package lichv

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func IsMatch(text string, filter string) bool {
	reg := regexp.MustCompile(filter)
	result := reg.FindAllString(text, -1)
	if len(result) > 0 {
		return true
	}else{
		return false
	}
}
func In(haystack interface{}, needle interface{}) (bool, error) {
	sVal := reflect.ValueOf(haystack)
	kind := sVal.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		for i := 0; i < sVal.Len(); i++ {
			if sVal.Index(i).Interface() == needle {
				return true, nil
			}
		}

		return false, nil
	}

	return false, errors.New("ErrUnSupportHaystack")
}

func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))
}

func GeneTimeUUID() string {
	now := time.Now().UnixNano()/1000
	return strconv.FormatUint(uint64(now),36)+strconv.Itoa(rand.New(rand.NewSource(now)).Intn(90)+10)
}

func URLAppendParams(uri string, key ,value string) (string,error) {
	l, err := url.Parse(uri)
	if err != nil {
		return uri,err
	}

	query := l.Query()
	query.Set(key,value)
	encodeurl := l.Scheme + "://" + l.Host + "?" + query.Encode()
	return encodeurl,nil
}

func Strval(value interface{}) string {
	// interface 转 string
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}

func BoolVal(value interface{}) bool {
	var key = false
	if value == nil {
		return key
	}
	switch value.(type) {
	case float64:
		ft := value.(float64)
		it := strconv.FormatFloat(ft, 'f', -1, 64)
		key = it == "0.0"
	case float32:
		ft := value.(float32)
		it := strconv.FormatFloat(float64(ft), 'f', -1, 64)
		key = it == "0.0"
	case int:
		it := value.(int)
		key = it != 0
	case uint:
		it := value.(uint)
		key = it != 0
	case int8:
		it := value.(int8)
		key = it != 0
	case uint8:
		it := value.(uint8)
		key = it != 0
	case int16:
		it := value.(int16)
		key = it != 0
	case uint16:
		it := value.(uint16)
		key = it != 0
	case int32:
		it := value.(int32)
		key = it != 0
	case uint32:
		it := value.(uint32)
		key = it != 0
	case int64:
		it := value.(int64)
		key = it != 0
	case uint64:
		it := value.(uint64)
		key = it != 0
	case string:
		it := value.(string)
		key = it != "" && it != "false"
	case []byte:
		it := value.([]byte)
		key = len(it) > 0
	default:
		newValue, _ := json.Marshal(value)
		it := string(newValue)
		key = len(it) > 0
	}
	return key
}

func HttpRequest(url, method, postdata string,headers map[string]interface{}) (interface{}, error) {
	client := &http.Client{}
	req, _ := http.NewRequest(strings.ToUpper(method), url, strings.NewReader(postdata))
	if len(headers) > 0 {
		for k, v := range headers {
			value := Strval(v)
			req.Header.Set(k,value)
		}
	}else{
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.129 Safari/537.36")
	}
	for k, v := range headers {
		value := Strval(v)
		req.Header.Set(k,value)
	}
	if strings.ToLower(method) == "post" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("http "+strings.ToLower(method)+" error", err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	var result interface{}
	err = json.Unmarshal(body, &result)
	if err == nil {
		return result, nil
	}
	err = xml.Unmarshal(body, &result)
	if err == nil {
		return result, nil
	}
	return body, nil
}

func urlJoin(href, base string) string {
	uri, err := url.Parse(href)
	if err != nil {
		return ""
	}
	baseUrl, err := url.Parse(base)
	if err != nil {
		return ""
	}
	return baseUrl.ResolveReference(uri).String()
}