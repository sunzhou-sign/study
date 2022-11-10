package utils

import (
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

const (
	URL_CODE_MACH = `%([0-9a-fA-F]{2})`
)

func UrlDecode(from string) (to string) {
	//1.匹配出所需转码的字符集
	re, err := regexp.Compile(URL_CODE_MACH)
	if nil != err {
		//log.Println("正则匹配失败!")
		to = from
		return to
	}
	findStr := re.FindAllString(from, -1) //正则匹配url转义编码

	//2.采用map进行去重
	mapStr := make(map[string]int) //生成集合
	for _, s := range findStr {
		mapStr[s] += 1
		//log.Println("正则匹配:", s)
	}

	//3.遍历去重后的集合,并替换编码
	for k, _ := range mapStr {
		sk, err := url2s(k)
		if nil != err {
			continue
		} else {
			from = strings.Replace(from, k, sk, -1)
		}
		//log.Println(k, "转码为:", sk)
	}
	to = from

	return to
}
func UrlEncode(from string) (to string) {
	to = url.QueryEscape(from)

	return to
}
func url2s(from string) (to string, err error) {
	//将百分号转码为正常字符
	bs := strings.Replace(from, `%`, ``, -1)
	is, err := strconv.ParseInt(bs, 16, 64)
	if nil != err {
		log.Println("util.url2s:strconv.ParseInt(bs, 16, 64) error!")
	}
	to = string(is)

	return
}

func UrlTest() {
	//%7B%22down_x%22%3ADOWN_X%2C%22down_y%22%3ADOWN_Y%2C%22up_x%22%3AUP_X%2C%22up_y%22%3AUP_Y%7D
	//`{"down_x":${DOWN_X},"down_y":${DOWN_Y},"up_x":${UP_X},"up_y":${UP_Y}}`
	//%7B%22down_x%22%3A${DOWN_X}%2C%22down_y%22%3A${DOWN_Y}%2C%22up_x%22%3A${UP_X}%2C%22up_y%22%3A${UP_Y}%7D
	ss := `{"down_x":DOWN_X,"down_y":DOWN_Y,"up_x":UP_X,"up_y":UP_Y}`
	rr := UrlEncode(ss)
	log.Println(rr)

	s := "%7B%22down_x%22%3A${DOWN_X}%2C%22down_y%22%3A${DOWN_Y}%2C%22up_x%22%3A${UP_X}%2C%22up_y%22%3A${UP_Y}%7D"
	r := UrlDecode(s)
	log.Println(r)
}
