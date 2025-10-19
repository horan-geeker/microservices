package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/tencentyun/cos-go-sdk-v5"
	"gopkg.in/gomail.v2"
	"io"
	"microservices/entity/config"
	"microservices/entity/consts"
	"microservices/pkg/http"
	"microservices/pkg/log"
	http2 "net/http"
	"net/mail"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Tencent interface {
	SendMailTo(mailTo string, from string, title string, content string) error
	SignCosFile(ctx context.Context, name string, duration time.Duration) (string, error)
	DownloadFromCOS(ctx context.Context, url string) (resp []byte, err error)
	UploadToCOS(ctx context.Context, userId int, content []byte, filename string) (string, string, error)
}

type tencent struct {
	mailServerAddress  string
	mailServerPassword string
	cosClient          *cos.Client
	client             http.Client
}

// SendMailTo .
func (t *tencent) SendMailTo(mailTo string, from string, title string, content string) error {
	m := gomail.NewMessage()
	fromMail := mail.Address{from, t.mailServerAddress}
	m.SetHeader("From", fromMail.String())
	m.SetHeader("To", mailTo)
	m.SetHeader("Subject", title)
	m.SetBody("text/html", content)

	d := gomail.NewDialer(consts.TencentSmtpServer, consts.TencentSmtpPort, t.mailServerAddress,
		t.mailServerPassword)
	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

// UploadToCOS .
func (t *tencent) UploadToCOS(ctx context.Context, userId int, content []byte, filename string) (string, string, error) {
	env := config.GetEnvConfig()
	fileUri := "users" + "/" + strconv.Itoa(userId) + "/" + filename
	f := strings.NewReader(string(content))
	httpResponse, err := t.cosClient.Object.Put(context.Background(), fileUri, f, nil)
	if err != nil {
		return "", "", err
	}
	body, err := io.ReadAll(httpResponse.Body)
	defer httpResponse.Body.Close()
	if err != nil {
		return "", "", err
	}
	log.Info(ctx, "cos", map[string]any{
		"body":       string(body),
		"httpStatus": httpResponse.StatusCode,
	})
	return env.COSBucketDomain + "/" + fileUri, fileUri, nil
}

func (t *tencent) getFileName(filename string, fileContent []byte, suffix string) string {
	hasher := md5.New()
	hasher.Write(fileContent)
	md5String := hex.EncodeToString(hasher.Sum(nil))
	if len(suffix) == 0 {
		fileFullName := strings.Split(filename, ".")
		if len(fileFullName) >= 2 {
			return md5String + "." + fileFullName[len(fileFullName)-1]
		}
	}
	if len(filename) == 0 {
		return md5String + "." + suffix
	}
	return filename
}

// SignCosFile .
func (t *tencent) SignCosFile(ctx context.Context, name string, duration time.Duration) (string, error) {
	// 获取文件签名
	env := config.GetEnvConfig()
	urlObj, err := t.cosClient.Object.GetPresignedURL(ctx, http2.MethodGet, name, env.COSSecretId, env.COSSecretKey,
		duration, nil)
	if err != nil {
		return "", err
	}
	return urlObj.String(), nil
}

func (t *tencent) DownloadFromCOS(ctx context.Context, url string) (resp []byte, err error) {
	err = t.client.Get(ctx, "", nil, &resp, http.WithTimeout(5), http.WithRawResponse(), http.WithBaseUrl(url))
	return
}

func newTencent(tencentOpt *config.TencentOptions) Tencent {
	env := config.GetEnvConfig()
	var u, _ = url.Parse(env.COSBucketDomain)
	var b = &cos.BaseURL{BucketURL: u}
	return &tencent{
		mailServerAddress:  tencentOpt.MailServerAddress,
		mailServerPassword: tencentOpt.MailServerPassword,
		cosClient: cos.NewClient(b, &http2.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  env.COSSecretId,
				SecretKey: env.COSSecretKey,
			},
		}),
		client: http.NewClient(http.WithJsonRequest(), http.WithJsonResponse(),
			http.WithProtocol("https"), http.WithDomain("api.weixin.qq.com")),
	}
}
