package tools

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

type DecryptErr struct {
	Message string
}

func newDecryptErr(message string) *DecryptErr {
	return &DecryptErr{Message: message}
}

func (e DecryptErr) Error() string {
	return e.Message
}

func Decrypt(encryptBs []byte, keyStr string) ([]byte, error) {
	type AESMsg struct {
		Encrypt string `json:"encrypt"`
	}
	var encrypt AESMsg
	err := json.Unmarshal(encryptBs, &encrypt)
	if err != nil {
		return nil, newDecryptErr(fmt.Sprintf("dataDecrypter jsonUnmarshalError[%v]", err))
	}
	buf, err := base64.StdEncoding.DecodeString(encrypt.Encrypt)
	if err != nil {
		return nil, newDecryptErr(fmt.Sprintf("base64StdEncode Error[%v]", err))
	}
	if len(buf) < aes.BlockSize {
		return nil, newDecryptErr("cipher  too short")
	}
	key := sha256.Sum256([]byte(keyStr))
	block, err := aes.NewCipher(key[:sha256.Size])
	if err != nil {
		return nil, newDecryptErr(fmt.Sprintf("AESNewCipher Error[%v]", err))
	}
	iv := buf[:aes.BlockSize]
	buf = buf[aes.BlockSize:]
	// CBC mode always works in whole blocks.
	if len(buf)%aes.BlockSize != 0 {
		return nil, newDecryptErr("ciphertext is not a multiple of the block size")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(buf, buf)
	n := strings.Index(string(buf), "{")
	if n == -1 {
		n = 0
	}
	m := strings.LastIndex(string(buf), "}")
	if m == -1 {
		m = len(buf) - 1
	}
	return buf[n : m+1], nil
}
