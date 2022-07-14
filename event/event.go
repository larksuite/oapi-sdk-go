package larkevent

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/larksuite/oapi-sdk-go/v3/core"
)

type EventHandler interface {
	Event() interface{}                        // 用于返回事件消息结构体（即承载回调消息内容的结构体）
	Handle(context.Context, interface{}) error // 用于处理事件
}

type EventHandlerModel interface {
	RawReq(req *EventReq)
}

type IReqHandler interface {
	Handle(ctx context.Context, req *EventReq) *EventResp
	Logger() larkcore.Logger
}

type DecryptErr struct {
	Message string
}

func newDecryptErr(message string) *DecryptErr {
	return &DecryptErr{Message: message}
}
func (e DecryptErr) Error() string {
	return e.Message
}

// eventDecrypt returns decrypt bytes
func EventDecrypt(encrypt string, secret string) ([]byte, error) {
	buf, err := base64.StdEncoding.DecodeString(encrypt)
	if err != nil {
		return nil, newDecryptErr(fmt.Sprintf("base64 decode error: %v", err))
	}
	if len(buf) < aes.BlockSize {
		return nil, newDecryptErr("cipher too short")
	}
	key := sha256.Sum256([]byte(secret))
	block, err := aes.NewCipher(key[:sha256.Size])
	if err != nil {
		return nil, newDecryptErr(fmt.Sprintf("AES new cipher error %v", err))
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

func Signature(timestamp string, nonce string, eventEncryptKey string, body string) string {
	var b strings.Builder
	b.WriteString(timestamp)
	b.WriteString(nonce)
	b.WriteString(eventEncryptKey)
	b.WriteString(body)
	bs := []byte(b.String())
	h := sha256.New()
	_, _ = h.Write(bs)
	bs = h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

type OptionFunc func(config *larkcore.Config)

func WithLogger(logger larkcore.Logger) OptionFunc {
	return func(config *larkcore.Config) {
		config.Logger = logger
	}
}

func WithLogLevel(logLevel larkcore.LogLevel) OptionFunc {
	return func(config *larkcore.Config) {
		config.LogLevel = logLevel
	}
}
