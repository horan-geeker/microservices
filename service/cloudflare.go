package service

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"microservices/entity/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type CloudflareService interface {
	UploadToR2(ctx context.Context, key string, content []byte, contentType string) (string, error)
	GetCustomDomain() string
	SignUrl(ctx context.Context, key string, expired int64) (string, error) // Added
}

type cloudflareService struct {
	s3Client      *s3.Client
	presignClient *s3.PresignClient // Added
	config        *config.CloudflareOptions
}

func NewCloudflareService(opt *config.CloudflareOptions) CloudflareService {
	// Initialize S3 client for Cloudflare R2
	accountID := opt.AccountId
	accessKey := opt.AccessKeyId
	secretKey := opt.SecretAccessKey

	if accountID == "" || accessKey == "" || secretKey == "" {
		// Log warning or return dummy? For now returning struct with nil client will cause panic on use, but factory pattern usually assumes config exists.
		// We can return nil client and handle error in Upload.
		return &cloudflareService{config: opt}
	}

	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(), // This is AWS config loader, not our entity/config package. Need to alias our import if clash.
		// Oh wait, `github.com/aws/aws-sdk-go-v2/config` is imported as `config`.
		// Our entity config is usually imported as `microservices/entity/config`.
		// Let's check imports.
		// existing: 	"github.com/aws/aws-sdk-go-v2/service/s3"
		// I need to import "microservices/entity/config" as well. I should alias AWS config or Entity config.
		// I'll alias AWS config as `awsconfig`.
		awsconfig.WithCredentialsProvider(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     accessKey,
				SecretAccessKey: secretKey,
			}, nil
		})),
		awsconfig.WithRegion("auto"), // R2 uses auto
	)
	if err != nil {
		fmt.Printf("unable to load SDK config, %v", err)
		return &cloudflareService{config: opt}
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID))
	})

	return &cloudflareService{
		s3Client:      client,
		presignClient: s3.NewPresignClient(client), // Added
		config:        opt,
	}
}

func (s *cloudflareService) UploadToR2(ctx context.Context, key string, content []byte, contentType string) (string, error) {
	if s.s3Client == nil {
		return "", fmt.Errorf("cloudflare client not initialized (check config)")
	}

	_, err := s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.config.Bucket),
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
	return fmt.Sprintf("https://%s.r2.cloudflarestorage.com/%s", s.config.AccountId, key), nil
}

func (s *cloudflareService) GetCustomDomain() string {
	return s.config.PublicDomain
}

func (s *cloudflareService) SignUrl(ctx context.Context, key string, expired int64) (string, error) {
	if s.presignClient == nil {
		return "", fmt.Errorf("presign client not initialized")
	}

	request, err := s.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(expired) * time.Second
	})
	if err != nil {
		return "", err
	}

	return request.URL, nil
}
