package hmac

import (
	"testing"
	"unsafe"
)

/**********************************************
** @Des: 测试代码用例
** @Author: zhangxueyuan
** @Date:   2018-09-14 17:15:26
** @Last Modified by:   zhangxueyuan
** @Last Modified time: 2018-09-14 17:15:26
***********************************************/

const checkMark = "\u2713"
const ballotX = "\u2717"

func Test360Decoding(t *testing.T) {
	{

		//e := []byte{byte(0xf7), byte(0x65), byte(0x83), byte(0x44), byte(0xb6),
		//	byte(0x13), byte(0x61), byte(0x57), byte(0x1e), byte(0x2f),
		//	byte(0x76), byte(0x58), byte(0xef), byte(0x5c), byte(0x5c),
		//	byte(0x0c), byte(0x17), byte(0x0b), byte(0xbc), byte(0xde),
		//	byte(0x0f), byte(0x3b), byte(0x0c), byte(0x8c), byte(0x9e),
		//	byte(0x70), byte(0x73), byte(0x7f), byte(0x7a), byte(0x33),
		//	byte(0x44), byte(0x32),}
		//i := []byte{
		//	byte(0xaf), byte(0x29), byte(0xa0), byte(0x78), byte(0x7e),
		//	byte(0x4a), byte(0x9f), byte(0x25), byte(0xdf), byte(0x36),
		//	byte(0x03), byte(0xd2), byte(0xe6), byte(0xb7), byte(0xcd),
		//	byte(0x54), byte(0xd1), byte(0x83), byte(0x49), byte(0xf8),
		//	byte(0x7e), byte(0x7b), byte(0x6b), byte(0xd0), byte(0xa1),
		//	byte(0x23), byte(0x36), byte(0x0b), byte(0x6a), byte(0xa3),
		//	byte(0x9e), byte(0xf4),
		//}
		//encSecret = string(e)
		//initSecret = string(i)

		//tmp := decoding360("AAAAAFoXl2AAAAAAAACWt3A5aCHt6sQxMoUYWw", encSecret, initSecret)
		//t.Log(string(tmp))
		//t.Log(bytesToInt64(tmp))

		price := Decoding360("AAAAAF1pAQkAAAAAAAiyd8d73I8jO6vWQDkduQ")
		t.Logf("%f", price)
	}

}

func TestVivoDecoding(t *testing.T) {

	var Prices = []struct {
		pirce  string
		result float64
	}{
		{
			"AACnZjHF6_AAAAFd0CcGgXlvO9SiDZdE4WFjxQ",
			1000.555,
		}, {
			"AACngidzJvAAAAFd0CjbmjEZG8p-cQX3Bivejg",
			800.1,
		}, {
			"AACnlNbLZTAAAAFd0CoVF18Sq4_W8Vc6YEAnWw",
			10,
		}, {
			"AACnojEyBvsAAAFd0Cr1H0h5--oQ1vFKUpVerg",
			1,
		}, {
			"",
			0,
		},
	}

	t.Log("Use demo price to testDecoding.")
	{
		encSecret = "c546946f87ee40578d32155d86bc3f21"
		initSecret = "780f1b12d777424bb4818d43060606b1"
		for _, u := range Prices {
			tmp := decodingVivo(u.pirce, encSecret, initSecret)
			if u.result == tmp {
				t.Logf(" should hava a %f with input %s. %v", u.result, u.pirce, checkMark)
			} else {
				t.Errorf("should have a %f, but is %f  with input %s. %v", u.result, tmp, u.pirce, ballotX)
			}
		}
	}
}

func TestOneXiaomDecoding(t *testing.T) {
	//tmp, v := DecodingXiaomi("cHJpY2VlbmNvZGluZ3doZTL4isFnbCZL5J5Wzg")
	//t.Logf("%s,%q ", tmp, v)
}

func BenchmarkString(b *testing.B) {
	p := DecodingXiaomiByte("cHJpY2VlbmNvZGluZ3doZTL4isFnbCZL5J5Wzg")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		String(p)
	}
}

func BenchmarkSafe2(b *testing.B) {
	p := DecodingXiaomiByte("cHJpY2VlbmNvZGluZ3doZTL4isFnbCZL5J5Wzg")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		String2(p)
	}
}

func String(p []byte) string {
	return string(p)
}

func String2(p []byte) string {
	return *(*string)(unsafe.Pointer(&p))
}

func TestXiaomiDecoding(t *testing.T) {

	var Prices = []struct {
		pirce  string
		result string
	}{
		{
			"cHJpY2VlbmNvZGluZ3doZTzUyMb2dirzYcjALQ",
			"1.0",
		}, {
			"cHJpY2VlbmNvZGluZ3doZTjUy8b2dirzYXcL0Q",
			"5.3",
		}, {
			"cHJpY2VlbmNvZGluZ3doZTzKyOjDQx_zjiVOmw",
			"100.555",
		}, {
			"cHJpY2VlbmNvZGluZ3doZTrN1vbBdirz1G52-Q",
			"77.07",
		}, {
			"",
			"0",
		},
	}

	t.Log("Use demo price to testDecoding.")
	{

		encSecret = "encrypriceencodingwhenintegratin"
		initSecret = "integpriceencodingwhenintegratin"
		for _, u := range Prices {
			tmp := decodingXiaomi(u.pirce, encSecret, initSecret)
			if u.result == tmp {
				t.Logf(" should hava a %s with input %s. %v", u.result, u.pirce, checkMark)
			} else {
				t.Errorf("should have a %s, but is %s  with input %s. %v", u.result, tmp, u.pirce, ballotX)
			}
		}
	}
}
