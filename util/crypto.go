package util


import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/config"
)

//加密
func AesCrypto(source []byte) ([]byte, error) {
	var block cipher.Block
	var err error
	if block, err = aes.NewCipher([]byte(config.Conf.Server.AesKey)); err != nil {
		return nil, fmt.Errorf("crypto err: %s", err.Error())
	}
	stream := cipher.NewCTR(block, []byte(config.Conf.Server.AesIv))
	var dest []byte
	stream.XORKeyStream(source, dest)
	return dest, nil
}

//解密
func AesDeCrypto(cryptoData []byte) ([]byte, error) {
	return AesCrypto(cryptoData)
}

