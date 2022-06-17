package event

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/larksuite/oapi-sdk-go/core"
)

type EventHandler interface {
	Event() interface{}                        // 用于返回事件消息结构体（即承载回调消息内容的结构体）
	Handle(context.Context, interface{}) error // 用于处理事件
}

type IReqHandler interface {
	ParseReq(ctx context.Context, req *EventReq) (string, error)
	DecryptEvent(ctx context.Context, cipherEventJsonStr string) (string, error)
	VerifyUrl(ctx context.Context, plainEventJsonStr string) (*EventResp, error)
	VerifySign(ctx context.Context, req *EventReq) error
	DoHandle(ctx context.Context, plainEventJsonStr string) (*EventResp, error)
}

type ReqHandler struct {
	IReqHandler
	*core.Config
}

func (h *ReqHandler) Handle(ctx context.Context, req *EventReq) (*EventResp, error) {

	// 解析请求
	cipherEventJsonStr, err := h.IReqHandler.ParseReq(ctx, req)
	if err != nil {
		return nil, err
	}

	// 消息解密
	plainEventJsonStr, err := h.IReqHandler.DecryptEvent(ctx, cipherEventJsonStr)
	if err != nil {
		return nil, err
	}

	// url验证逻辑处理
	resp, err := h.IReqHandler.VerifyUrl(ctx, plainEventJsonStr)
	if err != nil {
		return nil, err
	}
	if resp != nil {
		return resp, nil
	}

	// 签名验证
	err = h.IReqHandler.VerifySign(ctx, req)
	if err != nil {
		return nil, err
	}

	// 执行逻辑
	eventResp, err := h.IReqHandler.DoHandle(ctx, plainEventJsonStr)
	if err != nil {
		return nil, err
	}

	return eventResp, nil
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

type OptionFunc func(config *core.Config)

func WithLogger(logger core.Logger) OptionFunc {
	return func(config *core.Config) {
		config.Logger = logger
	}
}

func WithLogLevel(logLevel core.LogLevel) OptionFunc {
	return func(config *core.Config) {
		config.LogLevel = logLevel
	}
}
