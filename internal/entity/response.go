package entity

// Response .
type Response struct {
	RequestId string         `json:"request_id"`
	Data      map[string]any `json:"data"`
	Message   string         `json:"message"`
	Code      int            `json:"code"`
}
