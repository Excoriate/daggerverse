package main

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

// S3Scanner implements the Scanner interface for S3 buckets
type S3Scanner struct {
	BaseResource
	client       *s3.Client
	ctx          context.Context
	batchSize    int
	configLoader *configLoader
}

// NewS3Scanner creates a new S3 scanner instance
func NewS3Scanner(ctx context.Context, awsClient *AWSClient, cfg *configLoader, opts ...S3ScannerOption) (*S3Scanner, error) {
	client, err := awsClient.GetS3Client()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize S3 client: %w", err)
	}

	scanner := &S3Scanner{
		BaseResource: BaseResource{
			ResourceType: "s3:bucket",
			Region:       awsClient.cfg.Region,
		},
		client:       client,
		ctx:          ctx,
		batchSize:    10, // Default batch size
		configLoader: cfg,
	}

	// Apply options
	for _, opt := range opts {
		opt(scanner)
	}

	// Apply configuration from loader if available
	if cfg != nil && cfg.config != nil {
		if resourceConfig, exists := cfg.config.Resources["s3"]; exists {
			if resourceConfig.BatchSize != nil {
				scanner.batchSize = *resourceConfig.BatchSize
			}
		}
	}

	return scanner, nil
}

// S3ScannerOption defines functional options for S3Scanner
type S3ScannerOption func(*S3Scanner)

// WithBatchSize sets the number of buckets to process in parallel
func WithBatchSize(size int) S3ScannerOption {
	return func(s *S3Scanner) {
		if size > 0 {
			s.batchSize = size
		}
	}
}

// Scan implements the Scanner interface for S3 buckets
func (s *S3Scanner) Scan(criteria TagCriteria) ([]ScanResult, error) {
	buckets, err := s.listBuckets()
	if err != nil {
		return nil, fmt.Errorf("failed to list buckets: %w", err)
	}

	results := make([]ScanResult, 0, len(buckets))
	errChan := make(chan error)
	resultChan := make(chan ScanResult)

	sem := make(chan struct{}, s.batchSize)
	var wg sync.WaitGroup

	// Process buckets concurrently
	go func() {
		wg.Wait()
		close(resultChan)
		close(errChan)
	}()

	for _, bucket := range buckets {
		wg.Add(1)
		sem <- struct{}{} // Acquire semaphore

		go func(bucket types.Bucket) {
			defer wg.Done()
			defer func() { <-sem }() // Release semaphore

			result, err := s.scanBucket(*bucket.Name, criteria)
			if err != nil {
				errChan <- fmt.Errorf("failed to scan bucket %s: %w", *bucket.Name, err)
				return
			}
			resultChan <- result
		}(bucket)
	}

	// Collect results and errors
	var scanErrors []error
	go func() {
		for err := range errChan {
			scanErrors = append(scanErrors, err)
		}
	}()

	for result := range resultChan {
		results = append(results, result)
	}

	if len(scanErrors) > 0 {
		return results, fmt.Errorf("encountered %d errors during scan: %v", len(scanErrors), scanErrors[0])
	}

	return results, nil
}

// listBuckets retrieves all S3 buckets
func (s *S3Scanner) listBuckets() ([]types.Bucket, error) {
	result, err := s.client.ListBuckets(s.ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}
	return result.Buckets, nil
}

// scanBucket scans a single bucket for tag compliance
func (s *S3Scanner) scanBucket(bucketName string, criteria TagCriteria) (ScanResult, error) {
	// Get bucket location
	location, err := s.client.GetBucketLocation(s.ctx, &s3.GetBucketLocationInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return ScanResult{}, err
	}

	// Initialize result
	result := ScanResult{
		ResourceType: s.ResourceType,
		ResourceID:   bucketName,
		ARN:          fmt.Sprintf("arn:aws:s3:::%s", bucketName),
		Region:       string(location.LocationConstraint),
		Tags:         make(map[string]string),
		Metadata: map[string]interface{}{
			"CreationDate": "", // Will be populated if available
		},
	}

	// Get bucket tags
	tags, err := s.client.GetBucketTagging(s.ctx, &s3.GetBucketTaggingInput{
		Bucket: aws.String(bucketName),
	})

	// Handle no tags case (don't treat as error)
	if err != nil {
		// Check if the error is because there are no tags
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "NoSuchTagSet" {
			// This is fine - bucket just has no tags
			result.Tags = make(map[string]string)
		} else {
			return result, fmt.Errorf("failed to get bucket tags: %w", err)
		}
	} else if tags != nil {
		// Convert tags to map
		for _, tag := range tags.TagSet {
			if tag.Key != nil && tag.Value != nil {
				result.Tags[*tag.Key] = *tag.Value
			}
		}
	}

	// Determine compliance level from tags or default
	if level, exists := s.Tags["ComplianceLevel"]; exists {
		criteria.ComplianceLevel = level
	} else {
		criteria.ComplianceLevel = s.configLoader.config.Global.TagCriteria.ComplianceLevel
	}

	// Scan tags using base implementation, passing compliance levels
	result.Issues = s.ScanTags(criteria, s.configLoader.config.ComplianceLevels)

	// Set compliance tag based on issues
	if len(result.Issues) == 0 {
		result.ComplianceTag = "compliant"
	} else {
		result.ComplianceTag = "non-compliant"
	}

	return result, nil
}

// GetBucketMetadata retrieves additional metadata for a bucket
func (s *S3Scanner) GetBucketMetadata(bucketName string) (map[string]interface{}, error) {
	metadata := make(map[string]interface{})

	// Get bucket versioning
	versioning, err := s.client.GetBucketVersioning(s.ctx, &s3.GetBucketVersioningInput{
		Bucket: aws.String(bucketName),
	})
	if err == nil {
		metadata["Versioning"] = versioning.Status
	}

	// Get bucket encryption
	encryption, err := s.client.GetBucketEncryption(s.ctx, &s3.GetBucketEncryptionInput{
		Bucket: aws.String(bucketName),
	})
	if err == nil && encryption.ServerSideEncryptionConfiguration != nil {
		metadata["Encryption"] = encryption.ServerSideEncryptionConfiguration.Rules
	}

	// Get bucket policy status
	policyStatus, err := s.client.GetBucketPolicyStatus(s.ctx, &s3.GetBucketPolicyStatusInput{
		Bucket: aws.String(bucketName),
	})
	if err == nil && policyStatus.PolicyStatus != nil {
		metadata["PublicAccessAllowed"] = policyStatus.PolicyStatus.IsPublic
	}

	return metadata, nil
}

// Example usage for Dagger integration:
/*
func (m *AwsTagInspector) ScanS3Buckets(
	ctx context.Context,
	awsConfig AWSClientConfig,
	criteria TagCriteria,
) ([]ScanResult, error) {
	awsClient, err := NewAWSClient(ctx, awsConfig)
	if err != nil {
		return nil, err
	}

	scanner, err := NewS3Scanner(ctx, awsClient, WithBatchSize(20))
	if err != nil {
		return nil, err
	}

	return scanner.Scan(criteria)
}
*/
