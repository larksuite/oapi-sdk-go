package response

type Response struct {
	Error
	Data interface{} `json:"data"`
}
