package hmac

import (
	"adx-track-v8/utils"
	"adx-track-v8/utils/aes"
	ecb "adx-track-v8/utils/aes128"
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/aosting/goTools/plog"
	"github.com/aosting/goTools/str2"
	"strconv"
)

/**********************************************
** @Des: vivohmac
** @Author: zhangxueyuan
** @Date:   2018-09-14 15:30:50
** @Last Modified by:   zhangxueyuan
** @Last Modified time: 2018-09-14 15:30:50
***********************************************/

//vivo
var (
	encSecret  = "5396b9e2ca6b400c854611f98b63ac01"
	initSecret = "6e8a0638c5f04612aed43e4c0fdc2467"

	encSecretXM  = "cdbea2618c01bd7a585a6f1273443d57"
	initSecretXM = "c57f0b2dd33417706346c20df5d8009b"

	encSecretXM_XZ  = "ceb26f94e72ad8ee5e431cece385a493"
	initSecretXM_XZ = "f24fadbd8b4a2af21b9ce12a324e4d95"

	encSecret360, _  = hex.DecodeString("4a21168c2ee4c0d27280a120acec6796dc490a1dd946e00b60998e8120941b7a")
	initSecret360, _ = hex.DecodeString("727cbe02400807b5223e1e89136b769684faa0344ff898a95d2c588b3808c747")

	encSecretBaidu, _  = hex.DecodeString("01ca67f4002428b925c84535002428b925c87de2002428b925c88015e51d958b")
	initSecretBaidu, _ = hex.DecodeString("01ca67f4002428b925c8fa85002428b925c8fc51002428b925c8fdaefe86d161")

	key = "mediav20181228ad"

	// TODO 快手提供的测试密钥 上线需要修改
	// encSecretKwai = "123456789abcdefg"  // 测试密钥
	encSecretKwai = "FD6BAB306BAFDIE7" // 生产环境密钥

)

/**
https://developers.google.com/authorized-buyers/rtb/response-guide/decrypt-price?hl=zh-CN#decryption_scheme

// Add any required base64 padding (= or ==).
final_message_valid_base64 = AddBase64Padding(final_message)

// Web-safe decode, then base64 decode.
enc_price = WebSafeBase64Decode(final_message_valid_base64)

// Message is decoded but remains encrypted.
(iv, p, sig) = enc_price // Split up according to fixed lengths.
price_pad = hmac(e_key, iv)
price = p <xor> price_pad

conf_sig = hmac(i_key, price || iv)
success = (conf_sig == sig)
*/

func DecodingVivo(price string) float64 {
	if price, e := decoding(price, encSecret, initSecret); e != nil {
		fmt.Println("DecodingVivo is err: ", e, " price is:", price)
		return 0
	} else {
		return float64(bytesToInt64(price)) / float64(1000000)
	}
}

func decodingVivo(price, encSecret, initSecret string) float64 {
	if price, e := decoding(price, encSecret, initSecret); e != nil {

		return 0
	} else {
		return float64(bytesToInt64(price)) / float64(1000000)
	}
}

func decodingXiaomi(price, encSecret, initSecret string) string {
	if price, e := decoding(price, encSecret, initSecret); e != nil {
		return "0"
	} else {
		return string(bytes.Trim(price, "\x00"))
	}
}

func DecodingXiaomi(price, buss string) (string, []byte) {

	var b []byte
	if buss == "active" || buss == "tb_deeplink" ||
		buss == "tb_new_login" || buss == "tb_first_buy" {
		return "0", b
	}
	if price1, e := decoding(price, encSecretXM, initSecretXM); e != nil {
		plog.WARN("DecodingXiaomi is err: ", e, " price is:", price, "  buss is:", buss)
		return "0", b
	} else {
		return string(bytes.TrimRight(price1, "\x00")), price1
	}
}

func DecodingXZXiaomi(price string) (string, []byte) {
	var b []byte
	if price1, e := decoding(price, encSecretXM_XZ, initSecretXM_XZ); e != nil {
		plog.WARN("DecodingXZXiaomi is err: ", e)
		return "0", b
	} else {
		return string(bytes.TrimRight(price1, "\x00")), price1
	}
}

func DecodingXiaomiByte(price string) []byte {
	price1, _ := decoding(price, encSecretXM, initSecretXM)
	return price1
}

func DecodingKwaiByte(price string) string {
	if len(price) == 0 {
		return "0"
	}
	ivs := utils.UrlDecode(price)
	b, err := base64.StdEncoding.DecodeString(ivs)
	if err != nil {
		fmt.Println("base64.StdEncoding.DecodeString() err: ", err)
		return "0"
	}
	price1 := ecb.AesDecryptECB(b, []byte(encSecretKwai))
	return string(price1)
}

func decoding360(price, encSecret, initSecret string) []byte {
	var p []byte
	if price, e := decoding(price, encSecret, initSecret); e != nil {
		return p
	} else {
		return price
	}
}

func Decoding360(price string) float64 {

	if price, e := decoding(price, string(encSecret360), string(initSecret360)); e != nil {
		plog.WARN("DecodingVivo is err: ", e)
		return 0
	} else {
		//微分--> 分
		return float64(bytesToInt64(price)) / 10000
	}
}

func DecodingBaidu(price string) float64 {
	if str2.IsBlank(price) {
		return 0
	}
	if price, e := decoding(price, string(encSecretBaidu), string(initSecretBaidu)); e != nil {
		plog.WARN("DecodingBaidu is err: ", e)
		return 0
	} else {
		// 转换为分
		return float64(bytesToInt64(price))
	}
}

func DecodingTencent(price string) float64 {
	//if price, e := decoding(price, string(encSecretBaidu), string(initSecretBaidu)); e != nil {
	//	plog.WARN("DecodingTencent is err: ", e)
	//	return 0
	//} else {
	//	// 转换为分
	//	return float64(bytesToInt64(price))
	//}

	str := ecb.DecryptDESECB([]byte(price), []byte(key))
	if str == "" {
		return 0
	}
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return f
}

func decoding(price, encSecret, initSecret string) (pricebyte []byte, err error) {

	if len(price) == 0 {
		return pricebyte, errors.New("price is null")
	}

	var ivs []byte
	ivs, err = aes.Base64URLDecode(price)
	if err != nil {
		plog.ERROR("Decoding price error", err)
		return pricebyte, err
	}
	if len(ivs) != 28 {
		plog.ERROR("Decoding price error size!=28")
		return pricebyte, errors.New("Decoding price error size!=28")
	}

	iv := ivs[0:16]
	enc_price := ivs[16:24]
	signature := ivs[24:28]

	pad := hmac.New(sha1.New, []byte(encSecret))
	_, err = pad.Write(iv)
	if err != nil {
		plog.ERROR("hmac.New.Write error", err)
		return pricebyte, err
	}
	price_pad := pad.Sum(nil)
	tmpprice, _ := safeXORBytes(price_pad, enc_price)
	conf_tmp := hmac.New(sha1.New, []byte(initSecret))
	_, err = conf_tmp.Write(tmpprice)

	if err != nil {
		plog.ERROR("hmac.New.Write error", err)
		return pricebyte, err
	}
	_, err = conf_tmp.Write(iv)
	if err != nil {
		plog.ERROR("hmac.New.Write error", err)
		return pricebyte, err
	}
	conf_sig := conf_tmp.Sum(nil)
	if campareSign(signature, conf_sig) {
		//return float64(bytesToInt64(tmpprice)) / float64(1000000)
		return tmpprice, nil
	} else {
		plog.ERROR("signature error")
		err = errors.New("signature error")
	}
	return pricebyte, err
}

func safeXORBytes(a, b []byte) ([]byte, int) {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	dst := make([]byte, n)
	for i := 0; i < n; i++ {
		dst[i] = a[i] ^ b[i]
	}
	return dst, n
}

func campareSign(decodesige []byte, confsign []byte) bool {
	if len(decodesige) > len(confsign) {
		return false
	} else {
		for i, _ := range decodesige {
			if decodesige[i] != confsign[i] {
				return false
			}
		}
		return true
	}
}

func bytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

func pbindex(bs []byte) {
	for i, _ := range bs {
		fmt.Print(i, ":", bs[i], " ")
	}
	fmt.Println()
}
