package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gocarina/gocsv"
	"io"
	"microservices/pkg/consts"
	"microservices/pkg/log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client interface {
	Post(ctx context.Context, path string, urlParams map[string]string, reqBody any, respBody any, opts ...Option) error
	Get(ctx context.Context, path string, urlParams map[string]string, respBody any, opts ...Option) error
	Put(ctx context.Context, path string, urlParams map[string]string, reqBody any, respBody any, opts ...Option) error
	Delete(ctx context.Context, path string, urlParams map[string]string, reqBody any, respBody any,
		opts ...Option) error
	Patch(ctx context.Context, path string, urlParams map[string]string, reqBody any, respBody any,
		opts ...Option) error
}

type client struct {
	client *http.Client
	opts   []Option
}

// 创建 一个 http 客户端，通过传入参数 WithBaseUrl 来设置基础的 url
func NewClient(opts ...Option) Client {
	c := &client{}
	c.opts = make([]Option, 0, len(opts)+1)
	c.opts = append(c.opts, WithProtocol("http"), WithJsonRequest(), WithJsonResponse(), WithTimeout(time.Second*5))
	c.opts = append(c.opts, opts...)
	c.client = &http.Client{}
	return c
}

// form/data 上传文件示例
//
//	 reqBody := &bytes.Buffer{}
//		writer := multipart.NewWriter(reqBody)
//		part, err := writer.CreateFormFile("file", filepath.Base(filePath))
//		if err != nil {
//			return err
//		}
//		_, err = io.Copy(part, file)
//		if err != nil {
//			return err
//		}
//		writer.Close()
//		var resp any
//
// 需要注意使用 http.WithFormDataRequest(writer.FormDataContentType()) 来附加请求头 content-type 和 boundary
//
//	if err := s.client.Post(ctx, "", nil, reqBody, &resp, http.WithTimeout(time.Second*60), http.WithFormDataRequest(writer.FormDataContentType())); err != nil {
//		return err
//	}
func (c *client) Post(ctx context.Context, path string, urlParams map[string]string, reqBody any, respBody any, opts ...Option) error {
	return c.send(ctx, path, http.MethodPost, urlParams, reqBody, respBody, opts...)
}

func (c *client) Get(ctx context.Context, path string, urlParams map[string]string, respBody any, opts ...Option) error {
	return c.send(ctx, path, http.MethodGet, urlParams, nil, respBody, opts...)
}

func (c *client) Put(ctx context.Context, path string, urlParams map[string]string, reqBody any, respBody any, opts ...Option) error {
	return c.send(ctx, path, http.MethodPut, urlParams, reqBody, respBody, opts...)
}

func (c *client) Delete(ctx context.Context, path string, urlParams map[string]string, reqBody any, respBody any,
	opts ...Option) error {
	return c.send(ctx, path, http.MethodDelete, urlParams, reqBody, respBody, opts...)
}

func (c *client) Patch(ctx context.Context, path string, urlParams map[string]string, reqBody any, respBody any,
	opts ...Option) error {
	return c.send(ctx, path, http.MethodPatch, urlParams, reqBody, respBody, opts...)
}
func (c *client) send(ctx context.Context, path string, method string, urlParams map[string]string, reqBody any, respBody any, opts ...Option) error {
	options, err := c.getOptions(opts)
	if err != nil {
		return err
	}
	fullUrl := options.getBaseUrl() + path
	if urlParams != nil {
		values := url.Values{}
		for k, v := range urlParams {
			values.Add(k, v)
		}
		fullUrl = fullUrl + "?" + values.Encode()
	}
	var body io.Reader
	if reqBody != nil {
		if options.RequestHeader.Get(ContentType) == ContentTypeApplicationJson {
			reqBody, err := json.Marshal(reqBody)
			if err != nil {
				return err
			}
			body = bytes.NewReader(reqBody)
		} else if strings.Contains(options.RequestHeader.Get(ContentType), ContentTypeFormData) {
			body = reqBody.(*bytes.Buffer)
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, fullUrl, body)
	if err != nil {
		return err
	}
	req.Header = options.RequestHeader
	var responseReturn *string
	resp, err := c.client.Do(req)
	defer func() {
		var httpStatus *int
		if err == nil {
			resp.Body.Close()
			httpStatus = &resp.StatusCode
		}
		log.Info(ctx, "http-request", map[string]any{
			"url":            fullUrl,
			"reqBody":        reqBody,
			"method":         method,
			"responseStatus": httpStatus,
			"responseBody":   responseReturn,
		})
	}()
	if err != nil {
		return err
	}
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	responseStr := string(responseBody)
	if len(responseStr) > consts.MaxResponseLogLength {
		responseStr = responseStr[:consts.MaxResponseLogLength] + "..."
	}
	responseReturn = &responseStr
	//if resp.StatusCode != http.StatusOK {
	//	return ecode.ErrHttpRequestError
	//}
	contentType := resp.Header.Get(ContentType)
	if options.ResponseHeader.Get(contentType) == ContentTypeApplicationJson || contentType == ContentTypeApplicationJson {
		if err := json.Unmarshal(responseBody, respBody); err != nil {
			return err
		}
	} else if options.ResponseHeader.Get(contentType) == "text/csv" || contentType == "text/csv" {
		if err := gocsv.UnmarshalBytes(responseBody, respBody); err != nil {
			return err
		}
	} else {
		respBody = responseBody
	}
	return nil
}

// getOptions returns Options
func (c *client) getOptions(opts []Option) (*Options, error) {
	options := &Options{}
	// 附加全局参数
	for _, opt := range c.opts {
		opt(options)
	}
	// 附加每个单独请求的参数
	for _, opt := range opts {
		opt(options)
	}
	if options.Timeout > 0 {
		c.client.Timeout = options.Timeout
	}
	if err := options.parseURL(); err != nil {
		return nil, err
	}
	if options.Protocol == "" {
		return nil, errors.New("protocol is empty")
	}
	if options.Domain == "" && options.IP == "" {
		return nil, errors.New("domain and ip is empty")
	}
	return options, nil
}
