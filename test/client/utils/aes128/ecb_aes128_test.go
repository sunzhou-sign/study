package ecb

import (
	"encoding/base64"
	"encoding/hex"
	"log"
	"testing"
)

/**
 * @author     ：songdehua
 * @emall      ：200637086@qq.com
 * @date       ：Created in 2020/4/14 4:58 下午
 * @description：
 * @modified By：
 * @version    ：$
 */

func Test(t *testing.T) {

	//str := AesEncryptECB([]byte("songdehua1234567"), []byte("2020202020202020"))
	//fmt.Printf("entry str:%v \n", string(str))
	//
	//decStr := AesDecryptECB([]byte("jT+oscMTpLvRluBldiszFlbr1/PVmJov"), []byte("2020202020202020"))
	//fmt.Printf("decStr str:%v \n", string(decStr))

	origData := []byte("Hello World") // 待加密的数据
	key := []byte("ABCDEFGHIJKLMNOP") // 加密的密钥
	encrypted := AesEncryptECB(origData, key)
	log.Println("密文(hex)：", hex.EncodeToString(encrypted))
	log.Println("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
	decrypted := AesDecryptECB(encrypted, key)
	log.Println("解密结果：", string(decrypted))

}
