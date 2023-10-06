package service

import (
	"context"
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	log "github.com/sirupsen/logrus"
)

type Aliyun interface {
	SendSMSCode(ctx context.Context, phone string, code string) error
}

type aliyun struct {
	accessKeyId     string
	accessKeySecret string
	SmsSignName     string
	SmsTemplateCode string
}

// SendSMSCode .
func (a *aliyun) SendSMSCode(ctx context.Context, phone string, code string) error {
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", a.accessKeyId, a.accessKeySecret)

	type templateJson struct {
		Code string `json:"code"`
	}

	template := templateJson{code}
	templateStr, err := json.Marshal(template)

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "http"
	request.PhoneNumbers = phone
	request.SignName = a.SmsSignName
	request.TemplateCode = a.SmsTemplateCode
	request.TemplateParam = string(templateStr)

	response, err := client.SendSms(request)
	if err != nil {
		log.Error("send sms error", "response:", response, "error:", err)
		return err
	}
	return nil
}

func newAliyun(srv *serviceInstance) Aliyun {
	return &aliyun{
		accessKeyId:     srv.aliyunOptions.AccessKeyId,
		accessKeySecret: srv.aliyunOptions.AccessKeySecret,
		SmsSignName:     srv.aliyunOptions.SmsSignName,
		SmsTemplateCode: srv.aliyunOptions.SmsTemplateCode,
	}
}
