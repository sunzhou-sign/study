package hmac

import (
	"adx-track-v8/utils/aes"
	"encoding/hex"
	"testing"
)

/**********************************************
** @Des: baidu_test
** @Author: zhangxueyuan
** @4.Date:   2019-11-20 14:23:35
** @Last Modified by:   zhangxueyuan
** @Last Modified time: 2019-11-20 14:23:35
***********************************************/

func TestBase64Decode(t *testing.T) {
	ivs, err := aes.Base64URLDecode("Uja0xQADFz97jEpgW5IA8g0f455XNIjPRj8IqA")
	if err != nil {
		t.Error(err)
	} else {
		if "5236b4c50003173f7b8c4a605b9200f20d1fe39e573488cf463f08a8" == string(ivs) {
			t.Log("baidu Base64Decode succ")
		}
	}
}

func TestDecoding(t *testing.T) {
	baiduS, _ := hex.DecodeString("005e38cb01abc22fd5acd77001abc22fd5acfe8001abc22fd5ad06509b1bb875")
	baidui, _ := hex.DecodeString("005e38cb01abc22fd5ad6be001abc22fd5ad73b001abc22fd5ad7b802e88c532")

	t.Log(baiduS)
	t.Log(baidui)
	r, err := decoding("Uja0xQADFz97jEpgW5IA8g0f455XNIjPRj8IqA", string(baiduS), String(baidui))
	if err != nil {
		t.Error(err)
	} else {
		t.Log(bytesToInt64(r))
		t.Log(string(r))
	}
}

func TestBaidu(t *testing.T) {
	p := DecodingBaidu("XlTSsQACLYh7jEpgW5IA8oEjPWaBGv8usef3tg")
	t.Log(p)

}
