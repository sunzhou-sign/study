package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//判断字符串是否为空串
func IsBlank(str string) bool {
	str = strings.TrimSpace(str)
	if len(str) < 1 {
		return true
	}
	return false
}

func Str2Int32(arg string) int32 {
	num, _ := strconv.ParseInt(arg, 10, 32)
	return int32(num)
}

func Str2Float32(arg string) (result float32) {
	num, _ := strconv.ParseFloat(arg, 32)
	return float32(num)
}

func Int322String(in int32) string {
	return strconv.Itoa(int(in))
}

func StringToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func Float642String(in float64) string {
	return strconv.FormatFloat(in, 'f', 5, 64)
}

func Float322String(in float32) string {
	return strconv.FormatFloat(float64(in), 'f', 5, 32)
}

func Int642String(in int64) string {
	return strconv.FormatInt(in, 10)
}

// LocationTime 时区转换
func LocationTime() time.Time {
	local, err := time.LoadLocation("Asia/Chongqing")
	if err != nil {
		fmt.Println("util timeUtils LocationTime error:", err.Error())
		return time.Now()
	}
	return time.Now().In(local)
}

func Md5(b string) (tp string) {
	h := md5.New()
	h.Write([]byte(b))
	x := h.Sum(nil)
	y := make([]byte, 32)
	hex.Encode(y, x)

	return string(y)
}

func MyJson(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Println("marshal Err: ", err.Error())
	}
	return string(b)
}

func GetRandNumber(maxNum int32) int32 {
	rand.Seed(time.Now().UnixNano())
	return rand.Int31n(maxNum)
}

//普通的访问方式(默认连接方法)
func VisitUrl(url string) (data []byte, err error) {
	response, err := http.Get(url)
	if err != nil {
		return
	}

	defer response.Body.Close()
	data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	return
}

//json解析，去除html字符转义
func MarshalNoEscapeHTML(t interface{}) string {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	if err != nil {
		fmt.Println("marshal Err: ", err.Error())
	}
	return string(buffer.Bytes())
}

func IsNum(str string) bool {
	_, err := strconv.ParseFloat(str, 64)
	return err == nil
}

func UrlToHttp(url string) string {
	res := ""
	if strings.HasPrefix(url, "www") {
		res = "http://" + url
	} else {
		res = strings.Replace(url, "https", "http", 1)
	}

	return res
}

func Md5Sum32b(b string) string {
	h := md5.New()
	h.Write([]byte(b))
	x := h.Sum(nil)
	y := make([]byte, 32)
	hex.Encode(y, x)

	return string(y)
}

func GetNum(num int64) int32 {
	for i := 0; i < 100; i++ {
		if 1<<i == num {
			return int32(i)
		}
	}
	return -1
}

// PostHeader 发起post请求
func PostHeader(url string, client *http.Client, duration time.Duration, msg string, headers map[string]string) ([]byte, error) {
	if client == nil {
		client = &http.Client{
			Timeout: duration,
		}
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(msg))
	if err != nil {
		return nil, err
	}
	for key, header := range headers {
		req.Header.Set(key, header)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
