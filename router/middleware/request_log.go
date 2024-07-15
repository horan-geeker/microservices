package middleware

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"math/rand"
	"microservices/pkg/consts"
	"strconv"
	"time"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write .
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// WriteString .
func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// RequestLogger .
func RequestLogger(c *gin.Context) {
	begin := time.Now()

	buf, _ := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(buf))

	traceId := c.Request.Header.Get("traceId")
	spanId := c.Request.Header.Get("spanId")
	if spanId == "" {
		// 服务端生成的 spanId 增加 auto 前缀
		spanId = "auto" + strconv.Itoa(rand.Intn(999999-100000)+100000)
	}
	ctx := context.WithValue(c.Request.Context(), "spanId", spanId)
	ctx = context.WithValue(ctx, "traceId", traceId)
	c.Request = c.Request.WithContext(ctx)

	log.WithFields(log.Fields{
		"method":  c.Request.Method,
		"url":     c.Request.URL.Host + c.Request.URL.String(),
		"ip":      c.Request.Header.Get("x-forwarded-for"),
		"traceId": traceId,
		"spanId":  spanId,
		"body":    string(buf),
		"event":   "request",
		"header":  c.Request.Header,
	}).Info()

	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw

	c.Next()

	body := blw.body.String()
	if len(body) > consts.MaxResponseLogLength {
		body = body[:consts.MaxResponseLogLength] + "..."
	}
	// Log response body
	log.WithFields(log.Fields{
		"method":     c.Request.Method,
		"url":        c.Request.URL.Host + c.Request.URL.String(),
		"traceId":    traceId,
		"spanId":     spanId,
		"httpStatus": c.Writer.Status(),
		"timeCost":   time.Since(begin).Milliseconds(),
		"event":      "response",
		"body":       body,
	}).Info()
}
