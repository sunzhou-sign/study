package coco_crypto

import (
	"encoding/base64"
	"fmt"
)

/**
 * @author     ：songdehua
 * @emall      ：200637086@qq.com
 * @date       ：Created in 2021/7/7 4:51 下午
 * @description：
 * @modified By：
 * @version    ：$
 */

// tanx加解密需要的key
var trackTanxkey = [16]byte{'1', '1', 'a', '5', 'f', '9', '0', '4', '0', 'c', 'd', 'e', '2', '9', 'a', '0'}

// tanx参数加密
func TrackTanxEncrypt(apiParam []byte) (string, error) {
	//aes_128_ecb_encrypt
	tmp, err := Aes128ECBEncrypt(trackTanxkey, apiParam)
	if err != nil {
		fmt.Println("trackApiParamEncrypt err:", err)
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(tmp), nil
}

// track参数解密
func TrackTanxDecrypt(value string) ([]byte, error) {
	tmp, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		fmt.Println("TrackApiParamDecrypt base64 err:", err)
		return nil, err
	}
	plaintext, err := Aes128ECBDecrypt(trackTanxkey, tmp)
	if err != nil {
		fmt.Println("TrackApiParamDecrypt err:", err)
		return nil, err
	}
	return plaintext, nil
}
