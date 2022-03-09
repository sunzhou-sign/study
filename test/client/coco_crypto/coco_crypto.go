/**
 * @Author      : HanYuBin
 * @Description : nil
 * @File        : track
 * @Date        : 2021/5/24
 */

package coco_crypto

import (
	"encoding/base64"
	"errors"
	"fmt"
)

// 椰子加解密通用key
var Commonkey = [16]byte{'C', 'o', 'c', 'o', '-', 'A', 'd', 'x', ',', 'Y', 'e', 's', '!', '^', '_', '^'}

// 加密Url编码的AES数据，返回加密后的字符串数据
func CommonAesEncryptUrlEncode(value []byte) (string, error) {
	//aes_128_ecb_encrypt
	tmp, err := Aes128ECBEncrypt(Commonkey, value)
	if err != nil {
		fmt.Println("CommonAesEncryptUrlEncode err:", err)
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(tmp), nil
}

// 解密Url编码的AES数据，返回原始字符串数据
func CommonAesDecryptUrlEncode(value string) (priceStr string, err error) {
	defer func() {
		if p := recover(); p != nil {
			err = errors.New(p.(string))
		}
	}()
	var tmp []byte
	tmp, err = base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		fmt.Println("CommonAesDecryptUrlEncode base64 err:", err)
		return "", err
	}
	var price []byte
	price, err = Aes128ECBDecrypt(Commonkey, tmp)
	if err != nil {
		fmt.Println("CommonAesDecryptUrlEncode err:", err)
		return "", err
	}
	return string(price), err
}

// 功能：加密Url编码的AES数据。
// 参数1：原始数据字符串
// 参数2：16位密钥
// 返回：加密后的字符串数据
func AesEncryptUrlEncode(value []byte, key [16]byte) (string, error) {
	//aes_128_ecb_encrypt
	tmp, err := Aes128ECBEncrypt(key, value)
	if err != nil {
		fmt.Println("AesEncryptUrlEncode err:", err)
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(tmp), nil
}

// 功能：解密Url编码的AES数据。
// 参数1：加密数据字符串
// 参数2：16位密钥
// 返回：解密后的字符串数据
func AesDecryptUrlEncode(value string, key [16]byte) (string, error) {
	tmp, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		fmt.Println("AesDecryptUrlEncode base64 err:", err)
		return "", err
	}
	price, err := Aes128ECBDecrypt(key, tmp)
	if err != nil {
		fmt.Println("AesDecryptUrlEncode err:", err)
		return "", err
	}
	return string(price), nil
}
