package aes

import (
	"testing"
)

/**********************************************
** @Des: coding_test
** @Author: zhangxueyuan
** @4.Date:   2019-03-25 19:19:49
** @Last Modified by:   zhangxueyuan
** @Last Modified time: 2019-03-25 19:19:49
***********************************************/

func TestOppo(t *testing.T) {
	src := "0.56"
	key := "0123456789abcdef"

	crypted := AesEncrypt(src, key)
	a := Base64UrlSafeEncode(crypted)
	t.Log(a)
	AesDecrypt(crypted, []byte(key))
	Base64URLDecode("39W7dWTd_SBOCM8UbnG6qA")
}

func TestAesDecrypt2(t *testing.T) {
	src := "O-0mVdLfTGAnt3TClMitSg=="
	key := "NDAzMzY3LDE0MjI4"
	p, _ := Base64URLDecode(src)
	s, _ := AesDecrypt2(p, []byte(key))
	t.Log(string(s))

	gdt_token = "NDAzMzY3LDE0MjI4"
	t.Log(DecodingGDT(src))
	t.Log(DecodingGDT("VwSSFElLxs3wWy0LMsTy5Q=="))
	gdt_token = "bdec711cb5c9e211"
	t.Log(DecodingGDT("8N4CFH1xMH_j6yTfny4u0w=="))
	t.Log(DecodingGDT("svQRI9HfWpb5U0BB2oQ6Lg=="))
	t.Log(DecodingGDT("E6NbwpBaixr__es85CKLtQ=="))
}
