package service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	config2 "microservices/entity/config"
	"os"
)

type CloudflareService interface {
	UploadToR2(ctx context.Context, bucket string, key string, content []byte, contentType string) (string, error)
	GetCustomDomain() string
}

type cloudflareService struct {
	s3Client *s3.Client
	config   *config2.CloudflareOptions
}

func NewCloudflareService(opt *config2.CloudflareOptions) CloudflareService {
	// Initialize S3 client for Cloudflare R2
	accountID := opt.AccountId
	accessKey := opt.AccessKeyId
	secretKey := opt.SecretAccessKey

	if accountID == "" || accessKey == "" || secretKey == "" {
		// Log warning or return dummy? For now returning struct with nil client will cause panic on use, but factory pattern usually assumes config exists.
		// We can return nil client and handle error in Upload.
		return &cloudflareService{config: opt}
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     accessKey,
				SecretAccessKey: secretKey,
			}, nil
		})),
		config.WithRegion("auto"), // R2 uses auto
	)
	if err != nil {
		fmt.Printf("unable to load SDK config, %v", err)
		return &cloudflareService{}
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID))
	})

	return &cloudflareService{
		s3Client: client,
	}
}

func (s *cloudflareService) UploadToR2(ctx context.Context, bucket string, key string, content []byte, contentType string) (string, error) {
	if s.s3Client == nil {
		return "", fmt.Errorf("cloudflare client not initialized (check env vars)")
	}

	_, err := s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(content),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", err
	}

	// Return URL.
	// R2 Public URL format: https://<custom-domain>/<key> or https://<bucket>.<account>.r2.cloudflarestorage.com/<key>
	// We prefer custom domain if set.
	domain := s.GetCustomDomain()
	if domain != "" {
		return fmt.Sprintf("%s/%s", domain, key), nil
	}
	return fmt.Sprintf("https://%s.r2.cloudflarestorage.com/%s", os.Getenv("CLOUDFLARE_ACCOUNT_ID"), key), nil
}

func (s *cloudflareService) GetCustomDomain() string {
	return os.Getenv("CLOUDFLARE_PUBLIC_DOMAIN") // e.g., https://assets.aicbpng.com
}
