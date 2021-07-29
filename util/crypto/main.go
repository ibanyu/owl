package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/util"
	"log"
)

const (
	aesKey = "12345678abcdefgh"
	aesIv  = "abcdabcd12345678"
)

//加密
func AesEny(plaintext []byte) []byte {
	var (
		block cipher.Block
		err   error
	)
	//创建aes
	if block, err = aes.NewCipher([]byte(aesKey)); err != nil {
		log.Fatal(err)
	}
	//创建ctr
	stream := cipher.NewCTR(block, []byte(aesIv))
	//加密, src,dst 可以为同一个内存地址
	stream.XORKeyStream(plaintext, plaintext)
	return plaintext
}

//解密
func AesDec2(ciptext []byte) []byte {
	//对密文再进行一次按位异或就可以得到明文
	//例如：3的二进制是0011和8的二进制1000按位异或(相同为0,不同为1)后得到1011，
	//对1011和8的二进制1000再进行按位异或得到0011即是3
	return AesEny(ciptext)
}

func main() {
	plaintext := []byte("aaaaaa")
	fmt.Println("明文", string(plaintext))
	ciptext := AesEny(plaintext)
	xstr := util.StringifyByteDirectly(ciptext)
	fmt.Println("xstr aes:", xstr)
	fmt.Println("xstr destring aes:", util.ParseStringedByte(xstr))
	fmt.Println("加密", ciptext)
	platext2 := AesDec2(util.ParseStringedByte(xstr))
	fmt.Println("解密", string(platext2))

}
