module github.com/larksuite/oapi-sdk-go/sample

go 1.13

require (
	github.com/go-redis/redis/v8 v8.11.4
	github.com/larksuite/oapi-sdk-go/v2 v2.0.10-0.20220302071802-a789bd255ddb
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.6.1 // indirect
)

replace github.com/larksuite/oapi-sdk-go/v2 => ../v2
