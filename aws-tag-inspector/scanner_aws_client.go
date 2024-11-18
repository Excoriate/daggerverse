package main

import (
	"context"
	"fmt"

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

// AWSClient manages AWS service clients
type AWSClient struct {
	cfg aws.Config
	// serviceClients stores initialized service clients
	serviceClients map[string]interface{}
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

// NewAWSClient creates a new AWS client with the provided configuration
func NewAWSClient(ctx context.Context, cfg AWSClientConfig) (*AWSClient, error) {
	var awsCfg aws.Config
	var err error

	// Configure AWS SDK options
	opts := []func(*config.LoadOptions) error{
		config.WithRegion(cfg.Region),
		config.WithRetryMaxAttempts(cfg.MaxRetries),
	}

	// Handle credentials configuration
	if cfg.AccessKeyID != "" && cfg.SecretAccessKey != "" {
		opts = append(opts, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			cfg.SessionToken,
		)))
	} else if cfg.Profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(cfg.Profile))
	}

	// Handle custom endpoint if specified
	if cfg.Endpoint != "" {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL: cfg.Endpoint,
			}, nil
		})
		opts = append(opts, config.WithEndpointResolverWithOptions(customResolver))
	}

	// Load the AWS configuration
	awsCfg, err = config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config: %w", err)
	}

	return &AWSClient{
		cfg:            awsCfg,
		serviceClients: make(map[string]interface{}),
	}, nil
}

// GetServiceClient returns a client for the specified AWS service
func (c *AWSClient) GetServiceClient(serviceName string) (interface{}, error) {
	// Check if client is already initialized
	if client, exists := c.serviceClients[serviceName]; exists {
		return client, nil
	}

	// Get factory for service
	factory, exists := serviceFactories[serviceName]
	if !exists {
		return nil, fmt.Errorf("unsupported AWS service: %s", serviceName)
	}

	// Create new client
	client := factory.CreateClient(c.cfg)
	c.serviceClients[serviceName] = client

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

// GetS3Client returns an S3 client
func (c *AWSClient) GetS3Client() (*s3.Client, error) {
	client, err := c.GetServiceClient("s3")
	if err != nil {
		return nil, err
	}
	return client.(*s3.Client), nil
}

// Example of how to add a new service client:
/*
// Step 1: Add the AWS SDK import
import "github.com/aws/aws-sdk-go-v2/service/rds"

// Step 2: Create a factory for the new service
type RDSClientFactory struct{}

func (f *RDSClientFactory) CreateClient(cfg aws.Config) interface{} {
    return rds.NewFromConfig(cfg)
}

// Step 3: Register the factory in init()
func init() {
    serviceFactories["rds"] = &RDSClientFactory{}
}

// Step 4: Add a convenience method
func (c *AWSClient) GetRDSClient() (*rds.Client, error) {
    client, err := c.GetServiceClient("rds")
    if err != nil {
        return nil, err
    }
    return client.(*rds.Client), nil
}
*/
