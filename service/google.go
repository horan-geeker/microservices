package service

import (
	"bytes"
	"context"
	"fmt"
	"google.golang.org/genai"
	"io"
	"microservices/entity/config"
	entity "microservices/entity/model"
	"microservices/pkg/log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Google interface {
	ChatWithGemini(ctx context.Context, apiKey string, prompt string, images []*entity.Image) (*genai.GenerateContentResponse, error)
}

type google struct {
	opts *config.GoogleOptions
}

type GeminiResponse struct {
	Content string `json:"content"`
}

func newGoogle(opts *config.GoogleOptions) Google {
	return &google{
		opts: opts,
	}
}

type GoogleTransport struct {
	rt    http.RoundTripper
	proxy string
}

func NewGoogleTransport(proxyURL string) *GoogleTransport {
	transport := &http.Transport{}
	if proxyURL != "" {
		if proxy, err := url.Parse(proxyURL); err == nil {
			transport.Proxy = http.ProxyURL(proxy)
		}
	}
	return &GoogleTransport{
		rt:    transport,
		proxy: proxyURL,
	}
}

func (t *GoogleTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	fullUrl := req.URL.String()
	method := req.Method
	reqBody, _ := io.ReadAll(req.Body)
	req.Body = io.NopCloser(bytes.NewBuffer(reqBody))
	reqHeaders := make(map[string]string)
	for k, v := range req.Header {
		reqHeaders[k] = strings.Join(v, ",")
	}
	var httpStatus int
	begin := time.Now()

	resp, err := t.rt.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	respBody, _ := io.ReadAll(resp.Body)
	resp.Body = io.NopCloser(bytes.NewBuffer(respBody))
	reqBodyStr := string(reqBody)
	respBodyStr := string(respBody)
	if len(respBodyStr) > 1000 {
		respBodyStr = respBodyStr[:1000]
	}
	httpStatus = resp.StatusCode
	defer func() {
		traceId, _ := req.Context().Value("traceId").(string)
		spanId, _ := req.Context().Value("spanId").(string)
		if len(reqBodyStr) > 1000 {
			reqBodyStr = reqBodyStr[:1000]
		}
		log.Info(req.Context(), "google-client", map[string]any{
			"url":            fullUrl,
			"reqBody":        reqBodyStr,
			"method":         method,
			"responseStatus": httpStatus,
			"responseBody":   respBodyStr,
			"headers":        reqHeaders,
			"timeCost":       time.Since(begin).Milliseconds(),
			"traceId":        traceId,
			"spanId":         spanId,
		})
	}()
	return resp, nil
}

func (g *google) ChatWithGemini(ctx context.Context, apiKey string, prompt string, images []*entity.Image) (*genai.GenerateContentResponse, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("gemini API key is not configured")
	}

	// 使用配置中的代理URL，如果没有配置则使用默认值
	proxyURL := g.opts.ProxyURL
	if proxyURL == "" {
		proxyURL = "http://clash:7890" // 默认代理
	}

	httpClient := &http.Client{
		Transport: NewGoogleTransport(proxyURL),
		Timeout:   60 * time.Second,
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:     apiKey,
		HTTPClient: httpClient,
	})
	if err != nil {
		return nil, err
	}

	generateConfig := &genai.GenerateContentConfig{}
	parts := []*genai.Part{
		genai.NewPartFromText(prompt),
	}
	for _, image := range images {
		parts = append(parts, &genai.Part{
			InlineData: &genai.Blob{
				MIMEType: image.MimeType,
				Data:     image.Content,
			},
		})
	}
	response, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash-image",
		[]*genai.Content{
			genai.NewContentFromParts(parts, genai.RoleUser),
		},
		generateConfig,
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}
