package ipsearch

import (
	"encoding/json"
	"fmt"
	"testing"
)

/**
 * @author xiao.luo
 * @description This is the unit test for IpSearch
 */

func TestLoad(t *testing.T) {
	fmt.Println("Test Load IP Dat ...")
	p, err := New()
	if len(p.data) <= 0 || err != nil {
		t.Fatal("the IP Dat did not loaded successfully!")
	}
}

func TestGet(t *testing.T) {
	fmt.Println("Test Get IP ...")
	p, _ := New()
	ip := "223.115.20.170"
	ipstr := p.Get(ip)
	fmt.Println(ipstr)

	if ipstr != `亚洲|中国|湖北| |潜江|联通|429005|China|CN|112.896866|30.421215` {
		t.Fatal("the IP convert by ipSearch component is not correct!")
	}
}

func TestGetInfo(t *testing.T) {
	fmt.Println("Test Get IP ...")
	p, _ := New()
	ip := "139.207.124.130"
	ipObj := p.GetInfo(ip)

	d, _ := json.Marshal(ipObj)
	fmt.Println(string(d))
}
