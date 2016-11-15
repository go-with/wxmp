package server

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"sort"
	"strings"

	"github.com/go-with/util"
)

// AES解密
func aesDecrypt(ciphertext []byte, aesKey []byte) (plaintext []byte, err error) {
	if len(ciphertext)%len(aesKey) != 0 {
		err = errors.New("ciphertext is not a multiple of the block size")
		return
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return
	}
	iv := make([]byte, aes.BlockSize)
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}
	plaintext = make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	return
}

// AES加密
func aesEncrypt(plaintext []byte, aesKey []byte) (ciphertext []byte, err error) {
	if len(plaintext)%len(aesKey) != 0 {
		plaintext = pkcs7Pad(plaintext, len(aesKey))
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return
	}
	iv := make([]byte, aes.BlockSize)
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	ciphertext = make([]byte, len(plaintext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)

	return
}

// PKCS#7填充
func pkcs7Pad(plaintext []byte, blockSize int) []byte {
	size := blockSize - len(plaintext)%blockSize
	pads := bytes.Repeat([]byte{byte(size)}, size)
	return append(plaintext, pads...)
}

// AESKey解码
func decodeAESKey(encodingAESKey string) []byte {
	return base64Decode(encodingAESKey + "=")
}

// Base64解码
func base64Decode(str string) []byte {
	data, _ := base64.StdEncoding.DecodeString(str)
	return data
}

// Base64编码
func base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// 对消息签名
func signMsg(token, timestamp, nonce string, encrypt ...string) string {
	ss := sort.StringSlice{
		token,
		timestamp,
		nonce,
	}
	if len(encrypt) > 0 {
		ss = append(ss, encrypt[0])
	}
	ss.Sort()
	return util.SHA1Str(strings.Join(ss, ""))
}
