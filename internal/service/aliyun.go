package service

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	log "github.com/sirupsen/logrus"
	"microservices/internal/config"
)

type templateJson struct {
	Code string `json:"code"`
}

// SendSMSCodeByAliyun .
func SendSMSCodeByAliyun(phone string, code string) error {
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou",
		config.Env.AliyunAccessKeyId, config.Env.AliyunAccessKeySecret)

	template := templateJson{code}
	templateStr, err := json.Marshal(template)

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "http"
	request.PhoneNumbers = phone
	request.SignName = config.Env.AliyunSmsSignName
	request.TemplateCode = config.Env.AliyunSmsTemplateCode
	request.TemplateParam = string(templateStr)

	response, err := client.SendSms(request)
	if err != nil {
		log.Error("send sms error", "response:", response, "error:", err)
		return err
	}
	return nil
}
