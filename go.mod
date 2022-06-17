module github.com/feishu/oapi-sdk-go

go 1.13

replace (
	github.com/feishu/oapi-sdk-go/core => ./core
	github.com/feishu/oapi-sdk-go/event => ./event
	github.com/feishu/oapi-sdk-go/service/docx => ./service/docx
	github.com/feishu/oapi-sdk-go/service/im => ./service/im
)
