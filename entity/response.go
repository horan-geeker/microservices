package entity

// Response .
type Response struct {
	TraceId string         `json:"traceId"`
	SpanId  string         `json:"spanId"`
	Data    map[string]any `json:"data"`
	Message string         `json:"message"`
	Code    int64          `json:"code"`
}
