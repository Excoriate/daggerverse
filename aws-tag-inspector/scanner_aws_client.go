package main

import (
	"context"
	"fmt"
	"regexp"

	"github.com/Excoriate/daggerverse/aws-tag-inspector/internal/dagger"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// AWSClientConfig holds the configuration for AWS clients
type AWSClientConfig struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	SessionToken    string
	Profile         string
	Endpoint        string
	MaxRetries      int
}

// AWSClient represents an AWS client configuration
type AWSClient struct {
	cfg        aws.Config
	container  *dagger.Container
	serviceMap map[string]interface{}
}

// ServiceClientFactory defines the interface for creating service clients
type ServiceClientFactory interface {
	CreateClient(cfg aws.Config) interface{}
}

// EC2ClientFactory implements ServiceClientFactory for EC2
type EC2ClientFactory struct{}

func (f *EC2ClientFactory) CreateClient(cfg aws.Config) interface{} {
	return ec2.NewFromConfig(cfg)
}

// S3ClientFactory implements ServiceClientFactory for S3
type S3ClientFactory struct{}

func (f *S3ClientFactory) CreateClient(cfg aws.Config) interface{} {
	return s3.NewFromConfig(cfg)
}

// serviceFactories maps service names to their factories
var serviceFactories = map[string]ServiceClientFactory{
	"ec2": &EC2ClientFactory{},
	"s3":  &S3ClientFactory{},
}

// NewAWSClient creates a new AWS client with the given configuration
func NewAWSClient(ctx context.Context, cfg AWSClientConfig) (*AWSClient, error) {
	if cfg.Region == "" {
		return nil, fmt.Errorf("AWS region is required")
	}

	// Validate region format
	if !isValidAWSRegion(cfg.Region) {
		return nil, fmt.Errorf("invalid AWS region format: %s", cfg.Region)
	}

	if cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" {
		return nil, fmt.Errorf("AWS credentials are required")
	}

	// Load AWS configuration with custom options
	opts := []func(*config.LoadOptions) error{
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			cfg.SessionToken, // Include session token if provided
		)),
		config.WithRetryMaxAttempts(cfg.MaxRetries),
	}

	// Add custom endpoint if specified
	if cfg.Endpoint != "" {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           cfg.Endpoint,
				SigningRegion: cfg.Region,
				// Add support for different services
		})
			}, nil,
		})
		opts = append(opts, config.WithEndpointResolverWithOptions(customResolver))
	}
 err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &AWSClient{
		cfg:        awsCfg,
		serviceMap: make(map[string]interface{}),
	}, nil
}

// isValidAWSRegion validates the AWS region format
func isValidAWSRegion(region string) bool {
	// AWS region format: [a-z]{2}-[a-z]+-\d{1}
	// Examples: us-east-1, eu-west-2, ap-southeast-1
	regionPattern := `^[a-z]{2}-[a-z]+-\d{1}$`
	match, _ := regexp.MatchString(regionPattern, region)
	return match
}

// GetServiceClient returns a cached service client or creates a new one
func (c *AWSClient) GetServiceClient(service string) (interface{}, error) {
	if c.serviceMap == nil {
		c.serviceMap = make(map[string]interface{})
	}

	if client, exists := c.serviceMap[service]; exists {
		return client, nil
	}

	var client interface{}
	switch service {
	case "s3":
		client = s3.NewFromConfig(c.cfg)
	// Add other services as needed
	default:
		return nil, fmt.Errorf("unsupported service: %s", service)
	}

	c.serviceMap[service] = client
	return client, nil
}

// GetEC2Client returns an EC2 client
func (c *AWSClient) GetEC2Client() (*ec2.Client, error) {
	client, err := c.GetServiceClient("ec2")
	if err != nil {
		return nil, err
	}
	return client.(*ec2.Client), nil
}

// GetS3Client returns an S3 client with proper configuration
func (c *AWSClient) GetS3Client() (*s3.Client, error) {
	if c.serviceMap == nil {
		c.serviceMap = make(map[string]interface{})
	}

	if client, exists := c.serviceMap["s3"]; exists {
		return client.(*s3.Client), nil
	}

	// Configure S3 specific options
	s3Opts := s3.Options{
		Region:       c.cfg.Region,
		UsePathStyle: true, // Use path-style addressing for better compatibility
	}

	client := s3.NewFromConfig(c.cfg, func(o *s3.Options) {
		*o = s3Opts
	})

	c.serviceMap["s3"] = client
	return client, nil
}

// Container returns the Dagger container associated with this client
func (c *AWSClient) Container() *dagger.Container {
	return c.container
}

// SetContainer sets the Dagger container for this client
func (c *AWSClient) SetContainer(container *dagger.Container) {
	c.container = container
}