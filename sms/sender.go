package sms

import (
	"context"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

const (
	RedisSmsKey     = "sms:code:"
	CodeLen         = 6
	ExpiredDuration = 120
)

type Sender interface {
	SendShortMessage(context context.Context, mobile string, secretCode string) error
}

type AliyunSender struct {
	client       *dysmsapi20170525.Client
	accessKey    string
	accessSecret string
}

func NewAliyunMessageSender(accessKey string, accessSecret string) (Sender, error) {
	client, err := CreateClient(tea.String(accessSecret), tea.String(accessSecret))
	if err != nil {
		return nil, fmt.Errorf("cannot create aliyun dysmsapi Client")
	}
	return &AliyunSender{
		accessKey:    accessKey,
		accessSecret: accessSecret,
		client:       client,
	}, err
}

func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *dysmsapi20170525.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result = &dysmsapi20170525.Client{}
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}

func (sender *AliyunSender) SendShortMessage(context context.Context, mobile string, secretCode string) error {
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		SignName:      tea.String("阿里云短信测试"),
		TemplateCode:  tea.String("SMS_154950909"),
		PhoneNumbers:  tea.String(mobile),
		TemplateParam: tea.String("{\"code\":\"" + secretCode + "\"}"),
	}
	runtime := &util.RuntimeOptions{}

	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		_, _err := sender.client.SendSmsWithOptions(sendSmsRequest, runtime)
		if _err != nil {
			return _err
		}

		return nil
	}()

	if tryErr != nil {
		return fmt.Errorf("fail to send request : %s", tryErr)
	}
	return nil
}
