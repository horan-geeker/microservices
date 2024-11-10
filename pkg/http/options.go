package http

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const ContentType = "Content-Type"
const ContentTypeApplicationJson = "application/json"
const ContentTypeFormData = "multipart/form-data"

// Option sets client options.
type Option func(*Options)

type Options struct {
	IP             string
	Domain         string
	Timeout        time.Duration
	Protocol       string
	ResponseHeader http.Header
	RequestHeader  http.Header
	Port           int
	BaseUrl        string
	PrefixUri      string
}

func (o *Options) getBaseUrl() string {
	target := o.Domain
	if target == "" {
		target = o.IP
	}
	baseUrl := fmt.Sprintf("%s://%s", o.Protocol, target)
	if o.Port != 0 {
		return fmt.Sprintf("%s:%d", baseUrl, o.Port)
	}
	if o.PrefixUri != "" {
		return fmt.Sprintf("%s/%s", baseUrl, strings.Trim(o.PrefixUri, "/"))
	}
	return baseUrl
}

func (o *Options) parseURL() error {
	if o.BaseUrl != "" {
		parsedUrl, err := url.Parse(o.BaseUrl)
		if err != nil {
			return err
		}
		o.Protocol = parsedUrl.Scheme // 协议
		if strings.Contains(parsedUrl.Host, ":") {
			parts := strings.Split(parsedUrl.Host, ":")
			o.Domain = parts[0] // 域名
		}
		o.Port, err = strconv.Atoi(parsedUrl.Port()) // 端口，如果未指定则返回空字符串
		if err != nil {
			return err
		}
		o.PrefixUri = parsedUrl.Path // 路径
	}
	return nil
}

func WithTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.Timeout = timeout
	}
}

func WithDomain(domain string) Option {
	return func(o *Options) {
		o.Domain = domain
	}
}

func WithProtocol(Protocol string) Option {
	return func(o *Options) {
		o.Protocol = Protocol
	}
}

func WithIP(ip string) Option {
	return func(o *Options) {
		o.IP = ip
	}
}

// WithBaseUrl 识别 http://127.0.0.1:8080 格式的 url，解析到 protocol， ip 或 domain 和 port 中
func WithBaseUrl(baseUrl string) Option {
	return func(o *Options) {
		o.BaseUrl = baseUrl
	}
}

func WithRequestHeader(header http.Header) Option {
	return func(o *Options) {
		o.RequestHeader = header
	}
}

func WithJsonResponse() Option {
	return func(o *Options) {
		if o.ResponseHeader == nil {
			o.ResponseHeader = make(http.Header)
		}
		o.ResponseHeader.Set(ContentType, ContentTypeApplicationJson)
	}
}

func WithRawResponse() Option {
	return func(o *Options) {
		if o.ResponseHeader == nil {
			o.ResponseHeader = make(http.Header)
		}
		o.ResponseHeader.Set(ContentType, "application/octet-stream")
	}
}

func WithJsonRequest() Option {
	return func(o *Options) {
		if o.RequestHeader == nil {
			o.RequestHeader = make(http.Header)
		}
		o.RequestHeader.Set(ContentType, ContentTypeApplicationJson)
	}
}
func WithFormDataRequest(contentTypeWithBoundary string) Option {
	return func(o *Options) {
		if o.RequestHeader == nil {
			o.RequestHeader = make(http.Header)
		}
		o.RequestHeader.Set(ContentType, contentTypeWithBoundary)
	}
}
