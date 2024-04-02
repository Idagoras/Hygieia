package oss

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"time"
)

var BucketName = "hygieia-bucket-0001"

type Service interface {
	Upload(bucketName string, objectName string, timeout time.Duration, filePath string, retryCount int) error
	Download(bucketName string, objectName string, timeout time.Duration, savePath string, retryCount int) error
	Delete(bucketName string, objectName string, timeout time.Duration, retryCount int) error
	GetDownloadURL(bucketName string, objectName string, timeout time.Duration, retryCount int, validTime time.Duration) (string, error)
	ExistObject(bucketName string, objectName string) (bool, error)
}

type AliyunService struct {
	client   *oss.Client
	endPoint string
}

// LTAI5tDqg7uWWau237s5mfcn  CR3BTjjGEy0cpicElqF91JAnVsNXtT

func NewAliyunService() (*AliyunService, error) {
	provider, err := oss.NewEnvironmentVariableCredentialsProvider()
	if err != nil {
		return nil, fmt.Errorf("failed to new aliyun oss provider : %v", err)
	}
	client, err := oss.New("", "", "", oss.SetCredentialsProvider(&provider))
	if err != nil {
		return nil, fmt.Errorf("failed to new aliyun oss client : %v", err)
	}
	return &AliyunService{client: client, endPoint: "oss-cn-guangzhou.aliyuncs.com"}, nil
}

func (service *AliyunService) CreateBucket(bucketName string) error {
	err := service.client.CreateBucket(bucketName)
	if err != nil {
		return nil
	}
	return nil
}

func (service *AliyunService) Upload(bucketName string, objectName string, timeout time.Duration, filePath string, retryCount int) error {
	bucket, err := service.client.Bucket(bucketName)
	if err != nil {
		return err
	}

	err = bucket.PutObjectFromFile(objectName, filePath)
	if err != nil {
		return err
	}
	return nil
}

func (service *AliyunService) Download(bucketName string, objectName string, timeout time.Duration, savePath string, retryCount int) error {
	bucket, err := service.client.Bucket(bucketName)
	if err != nil {
		return err
	}

	err = bucket.GetObjectToFile(objectName, savePath)
	if err != nil {
		return err
	}
	return nil
}

func (service *AliyunService) GetDownloadURL(bucketName string, objectName string, timeout time.Duration, retryCount int, validTime time.Duration) (string, error) {
	bucket, err := service.client.Bucket(bucketName)
	if err != nil {
		return "", err
	}
	signedURL, err := bucket.SignURL(objectName, oss.HTTPGet, int64(validTime.Seconds()))
	if err != nil {
		return "", err
	}
	return signedURL, nil
}

func (service *AliyunService) Delete(bucketName string, objectName string, timeout time.Duration, retryCount int) error {
	bucket, err := service.client.Bucket(bucketName)
	if err != nil {
		return err
	}
	err = bucket.DeleteObject(objectName)
	if err != nil {
		return err
	}
	return nil
}
