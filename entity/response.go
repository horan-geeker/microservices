package entity

// Response .
type Response[T any] struct {
	TraceId string `json:"traceId"`
	SpanId  string `json:"spanId"`
	Data    T      `json:"data"`
	Message string `json:"message"`
	Code    int64  `json:"code"`
}
