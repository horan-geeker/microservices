package middleware

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"math/rand"
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

	requestId := c.Request.Header.Get("requestId")
	if requestId == "" {
		// 服务端生成的 request_id 增加 auto 前缀
		requestId = "auto" + strconv.Itoa(rand.Intn(999999-100000)+100000)
	}

	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "requestId", requestId))

	log.WithFields(log.Fields{
		"method":    c.Request.Method,
		"url":       c.Request.URL.Host + c.Request.URL.String(),
		"ip":        c.Request.Header.Get("x-forwarded-for"),
		"requestId": requestId,
		"body":      string(buf),
		"event":     "api",
		"header":    c.Request.Header,
	}).Info()

	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw

	c.Next()

	body := blw.body.String()
	if len(body) > 4096 {
		body = body[:4096]
	}
	// Log response body
	log.WithFields(log.Fields{
		"method":     c.Request.Method,
		"url":        c.Request.URL.Host + c.Request.URL.String(),
		"requestId":  requestId,
		"httpStatus": c.Writer.Status(),
		"timeCost":   time.Since(begin).Milliseconds(),
		"event":      "response",
		"body":       body,
	}).Info()
}
