/*
描述 :  golang  AES/ECB/PKCS5  加密解密
date : 2016-04-08
*/

package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"test/client/utils"
)

var gdt_token = "bdec711cb5c9e211"

//待测试

func Coding() {
	/*
	 *src 要加密的字符串
	 *key 用来加密的密钥 密钥长度可以是128bit、192bit、256bit中的任意一个
	 *16位key对应128bit
	 */
	src := "0.56"
	key := "0123456789abcdef"

	crypted := AesEncrypt(src, key)
	AesDecrypt(crypted, []byte(key))
	Base64URLDecode("39W7dWTd_SBOCM8UbnG6qA")

}

func Decoding(price string) (string, error) {
	if len(price) == 0 {
		err := errors.New("input params is empty!")
		return "0", err
	}
	base, err := hex.DecodeString(price)
	if err != nil {
		err = errors.New("Decoding price error" + err.Error())
		fmt.Println("Decoding price error", err)
		return "0", err
	}
	Oppo_Key := "d95fe73c6dee46e8"
	e, err := AesDecrypt([]byte(base), []byte(Oppo_Key))
	if err != nil {
		fmt.Println("AesDecrypt price error", err)
		err = errors.New("AesDecrypt price error" + err.Error())
		return "0", err
	}
	return string(e), nil
}

func DecodingGDT(price string) (int64, string) {
	if len(price) == 0 {
		return 0, "0"
	}
	var result []byte
	var err error
	if result, err = Base64URLDecode(price); err != nil {
		fmt.Println("GDT price base64 error: ", err)
		return 0, "0"
	} else {
		if result, err = AesDecrypt2(result, []byte(gdt_token)); err != nil {
			fmt.Println("GDT price aesDecrypt error: ", err)
			return 0, "0"
		}
		price := strings.TrimSpace(string(result))
		v, _ := utils.StringToInt64(price)
		return v, price
	}
}

func Base64URLDecode(data string) ([]byte, error) {
	var missing = (4 - len(data)%4) % 4
	data += strings.Repeat("=", missing)
	//res, _ := base64.URLEncoding.DecodeString(data)
	//fmt.Println("  decodebase64urlsafe is :", string(res), err)
	return base64.URLEncoding.DecodeString(data)
}

func Base64UrlSafeEncode(source []byte) string {
	// Base64 Url Safe is the same as Base64 but does not contain '/' and '+' (replaced by '_' and '-') and trailing '=' are removed.
	bytearr := base64.StdEncoding.EncodeToString(source)
	safeurl := strings.Replace(string(bytearr), "/", "_", -1)
	safeurl = strings.Replace(safeurl, "+", "-", -1)
	safeurl = strings.Replace(safeurl, "=", "", -1)
	return safeurl
}

func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("AesDecrypt err is:", err)
		return nil, err
	}
	blockMode := NewECBDecrypter(block)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	if len(origData) == 0 {
		return nil, errors.New("origData is empty")
	}
	origData, err = PKCS5UnPadding(origData)
	if err != nil {
		fmt.Println("PKCS5UnPadding err is:", err)
		return nil, err
	}
	return origData, nil
}

func AesDecrypt2(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("AesDecrypt err is:", err)
		return nil, err
	}
	blockMode := NewECBDecrypter(block)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	if len(origData) == 0 {
		return nil, errors.New("origData is empty")
	}

	return origData, nil
}

func AesEncrypt(src, key string) []byte {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Println("key error1", err)
	}
	if src == "" {
		fmt.Println("plain content empty")
	}
	ecb := NewECBEncrypter(block)
	content := []byte(src)
	content = PKCS5Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)
	// 普通base64编码加密 区别于urlsafe base64
	//fmt.Println("base64 result:", base64.StdEncoding.EncodeToString(crypted))
	//fmt.Println("base64UrlSafe result:", Base64UrlSafeEncode(crypted))
	return crypted
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	//fmt.Print("unpadding", unpadding, "length:", length)
	if (length-unpadding) > length || (length-unpadding) < 0 {
		return nil, errors.New(fmt.Sprintf("length - unpadding ;%d ; %d", (length - unpadding), unpadding))
	}
	return origData[:(length - unpadding)], nil
}

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbEncrypter ecb

// NewECBEncrypter returns a BlockMode which encrypts in electronic code book
// mode, using the given Block.
func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}
func (x *ecbEncrypter) BlockSize() int { return x.blockSize }
func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

type ecbDecrypter ecb

// NewECBDecrypter returns a BlockMode which decrypts in electronic code book
// mode, using the given Block.
func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}
func (x *ecbDecrypter) BlockSize() int { return x.blockSize }
func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}
