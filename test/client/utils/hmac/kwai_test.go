package hmac

import (
	"fmt"
	"testing"
)

/**
 * @author     ：songdehua
 * @emall      ：200637086@qq.com
 * @date       ：Created in 2020/7/22 5:37 下午
 * @description：
 * @modified By：
 * @version    ：$
 */

func TestDecodeKwaiByECB(t *testing.T) {
	//s := "1001"
	//key := "123456789abcdefg"
	//
	//ans := EntryptKwaiECB([]byte(s), []byte(key))
	//fmt.Println(ans)

	//origData := []byte("1001")        // 待加密的数据
	//key := []byte("123456789abcdefg") // 加密的密钥
	//encrypted := ecb.AesEncryptECB(origData, key)
	//log.Println("密文(hex)：", hex.EncodeToString(encrypted))
	//log.Println("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
	//log.Println("密文(urlencode)：", util.UrlEncode(base64.StdEncoding.EncodeToString(encrypted)))
	//decrypted := ecb.AesDecryptECB(encrypted, key)
	//log.Println("解密结果：", string(decrypted))

	secret := "7WcS%2FIj3fwwiZB%2FFqX9alg%3D%3D" // 密文
	fmt.Println("原价格：", DecodingKwaiByte(secret))
}
