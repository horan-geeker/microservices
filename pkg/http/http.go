package http

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gocarina/gocsv"
	"io"
	"microservices/pkg/log"
	"net/http"
	"net/url"
)

type HttpRequest[T any] struct {
	timeout int
}

func (h *HttpRequest[T]) Get(ctx context.Context, path string, urlParams map[string]string, body map[string]any,
	headers map[string]string) (*T, error) {
	return h.request(ctx, path, http.MethodGet, urlParams, body, headers)
}

func (h *HttpRequest[T]) GetRawResponse(ctx context.Context, path string, urlParams map[string]string, body map[string]any,
	headers map[string]string) ([]byte, error) {
	response, err := h.sendRequest(ctx, path, http.MethodGet, urlParams, body, headers)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(response.Body)
}

func (h *HttpRequest[T]) Post(ctx context.Context, path string, urlParams map[string]string, body map[string]any,
	headers map[string]string) (*T, error) {
	return h.request(ctx, path, http.MethodPost, urlParams, body, headers)
}

func (h *HttpRequest[T]) Put(ctx context.Context, path string, urlParams map[string]string, body map[string]any,
	headers map[string]string) (*T, error) {
	return h.request(ctx, path, http.MethodPut, urlParams, body, headers)
}

func (h *HttpRequest[T]) Delete(ctx context.Context, path string, urlParams map[string]string, body map[string]any,
	headers map[string]string) (*T, error) {
	return h.request(ctx, path, http.MethodDelete, urlParams, body, headers)
}

func (h *HttpRequest[T]) Patch(ctx context.Context, path string, urlParams map[string]string, body map[string]any,
	headers map[string]string) (*T, error) {
	return h.request(ctx, path, http.MethodPatch, urlParams, body, headers)
}

func (h *HttpRequest[T]) request(ctx context.Context, urlPath, method string, urlParams map[string]string, reqBody map[string]any,
	headers map[string]string) (*T, error) {
	response, err := h.sendRequest(ctx, urlPath, method, urlParams, reqBody, headers)
	if err != nil {
		return nil, err
	}
	contentType := response.Header.Get("content-type")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
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

func (h *HttpRequest[T]) sendRequest(ctx context.Context, urlPath, method string, urlParams map[string]string, reqBody map[string]any,
	headers map[string]string) (*http.Response, error) {
	if urlParams != nil {
		values := url.Values{}
		for k, v := range urlParams {
			values.Add(k, v)
		}
		urlPath = urlPath + "?" + values.Encode()
	}
	var reqBodyBuff []byte
	if reqBody != nil {
		var err error
		reqBodyBuff, err = json.Marshal(reqBody)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, urlPath, bytes.NewReader(reqBodyBuff))
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
	responseBody, err := io.ReadAll(response.Body)
	response.Body = io.NopCloser(bytes.NewBuffer(responseBody))
	log.Info(ctx, "http-request", map[string]any{
		"url":            urlPath,
		"params":         urlParams,
		"reqBody":        reqBody,
		"method":         method,
		"responseStatus": response.StatusCode,
		"responseBody":   string(responseBody),
	})
	return response, err
}

func NewHttp[T any](opts *Options) *HttpRequest[T] {
	return &HttpRequest[T]{
		timeout: opts.Timeout,
	}
}
