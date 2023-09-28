package http

import (
	"context"
	"encoding/json"
	"github.com/gocarina/gocsv"
	"io"
	"microservices/pkg/log"
	"net/http"
	"net/url"
	"strings"
)

type httpRequest[T any] struct {
	domain string
}

func (h *httpRequest[T]) Get(ctx context.Context, path string, urlParams map[string]string, body *string,
	headers map[string]string) (*T, error) {
	return h.request(ctx, path, http.MethodGet, urlParams, body, headers)
}

func (h *httpRequest[T]) Post(ctx context.Context, path string, urlParams map[string]string, body *string,
	headers map[string]string) (*T, error) {
	return h.request(ctx, path, http.MethodPost, urlParams, body, headers)
}

func (h *httpRequest[T]) Put(ctx context.Context, path string, urlParams map[string]string, body *string,
	headers map[string]string) (*T, error) {
	return h.request(ctx, path, http.MethodPut, urlParams, body, headers)
}

func (h *httpRequest[T]) Delete(ctx context.Context, path string, urlParams map[string]string, body *string,
	headers map[string]string) (*T, error) {
	return h.request(ctx, path, http.MethodDelete, urlParams, body, headers)
}

func (h *httpRequest[T]) Patch(ctx context.Context, path string, urlParams map[string]string, body *string,
	headers map[string]string) (*T, error) {
	return h.request(ctx, path, http.MethodPatch, urlParams, body, headers)
}

func (h *httpRequest[T]) request(ctx context.Context, urlPath, method string, urlParams map[string]string, reqBody *string,
	headers map[string]string) (*T, error) {
	fullUrl := strings.Trim(h.domain, "/") + "/" + strings.Trim(urlPath, "/")
	if urlParams != nil {
		values := url.Values{}
		for k, v := range urlParams {
			values.Add(k, v)
		}
		fullUrl = fullUrl + "?" + values.Encode()
	}
	var reqBodyBuff io.Reader
	if reqBody != nil {
		reqBodyBuff = strings.NewReader(*reqBody)
	}
	req, err := http.NewRequestWithContext(ctx, method, urlPath, reqBodyBuff)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	log.Log(ctx, "http-request", map[string]any{
		"url":    fullUrl,
		"method": method,
		"body":   string(responseBody),
	})
	contentType := response.Header.Get("content-type")
	var res *T
	if contentType == "text/csv" {
		if err := gocsv.UnmarshalBytes(responseBody, &res); err != nil {
			return nil, err
		}
	} else {
		if err := json.Unmarshal(responseBody, &res); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func NewHttp[T any](domain string) *httpRequest[T] {
	return &httpRequest[T]{
		domain: domain,
	}
}
