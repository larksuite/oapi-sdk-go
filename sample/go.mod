module github.com/larksuite/oapi-sdk-go/sample

go 1.17

require github.com/larksuite/oapi-sdk-go/v3 v3.0.0

replace github.com/larksuite/oapi-sdk-go/v3 => ../

require github.com/bytedance/sonic v1.5.0

require (
	github.com/chenzhuoyu/base64x v0.0.0-20211019084208-fb5309c8db06 // indirect
	github.com/klauspost/cpuid/v2 v2.0.9 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	golang.org/x/arch v0.0.0-20210923205945-b76863e36670 // indirect
)
