package main

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/Excoriate/daggerverse/aws-tag-inspector/internal/dagger"
)

// resourceScanner defines the interface for scanning different AWS resource types
type resourceScanner interface {
	Scan(criteria TagCriteria) ([]ScanResult, error)
}

// Scan performs a comprehensive AWS resource tag inspection and validation.
//
// This method dynamically scans AWS resources based on the configuration,
// supporting extensibility for multiple resource types.
//
// Parameters:
//   - ctx: Optional context for controlling the scan operation's lifecycle and timeout.
//
// Returns:
//   - A Dagger file containing the scan results in JSON format.
//   - An error if the scan process encounters any critical failures.
func (m *AwsTagInspector) Scan(
	ctx context.Context,
) (*dagger.File, error) {
	if m.Cfg == nil {
		return nil, Errorf("configuration is required for scanning")
	}

	if !m.Cfg.Global.Enabled {
		return nil, Errorf("scanning is globally disabled in configuration")
	}

	var allResults []ScanResult

	// Dynamically scan configured resources
	for resourceType, resourceConfig := range m.Cfg.Resources {
		// Skip if resource is not enabled
		if !resourceConfig.Enabled {
			continue
		}

		// Determine tag criteria (resource-specific or global)
		tagCriteria := resourceConfig.TagCriteria
		if reflect.DeepEqual(tagCriteria, TagCriteria{}) {
			tagCriteria = m.Cfg.Global.TagCriteria
		}

		// Ensure compliance level is set, prioritizing resource-specific configuration
		if tagCriteria.ComplianceLevel == "" {
			// If resource-specific compliance level is empty, use global compliance level
			if m.Cfg.Global.TagCriteria.ComplianceLevel != "" {
				tagCriteria.ComplianceLevel = m.Cfg.Global.TagCriteria.ComplianceLevel
			}
		}

		// Scan based on resource type
		results, err := m.scanResourceByType(ctx, resourceType, tagCriteria)
		if err != nil {
			return nil, WrapError(err, fmt.Sprintf("failed to scan %s resources", resourceType))
		}

		allResults = append(allResults, results...)
	}

	// Convert results to JSON
	jsonContent, err := m.formatResults(allResults)
	if err != nil {
		return nil, WrapError(err, "failed to format scan results")
	}

	// Get a container from the AWS client
	container := m.Ctr
	if container == nil {
		return nil, Errorf("failed to get container for results")
	}

	// Create a file in the ocntainer
	resultsFile := container.Directory("mnt/").
		WithNewFile("scan-results.json", jsonContent)

	// return the file only, extracted from the container.
	return resultsFile.File("scan-results.json"), nil
}

// scanResourceByType dynamically scans a specific resource type
func (m *AwsTagInspector) scanResourceByType(
	ctx context.Context,
	resourceType string,
	tagCriteria TagCriteria,
) ([]ScanResult, error) {
	if m.AWSClient == nil {
		return nil, Errorf("AWS client is not initialized")
	}

	switch resourceType {
	case "s3":
		s3Scanner, err := NewS3Scanner(ctx, m.AWSClient, m.Cfg)
		if err != nil {
			return nil, WrapError(err, "failed to create S3 scanner")
		}

		results, err := s3Scanner.Scan(tagCriteria)
		if err != nil {
			return nil, WrapError(err, "failed to scan S3 buckets")
		}
		return results, nil
	// Add more resource types here as they are implemented
	// case "ec2":
	// 	return m.scanEC2InstancesAPI(ctx, tagCriteria)
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}
}

// formatResults converts scan results to a JSON string
func (m *AwsTagInspector) formatResults(results []ScanResult) (string, error) {
	if len(results) == 0 {
		return "[]", nil
	}

	// Group results by resource type
	groupedResults := make(map[string][]ScanResult)
	for _, result := range results {
		groupedResults[result.ResourceType] = append(groupedResults[result.ResourceType], result)
	}

	// Format summary
	summary := struct {
		TotalResources int                     `json:"total_resources"`
		ResourceTypes  map[string][]ScanResult `json:"resource_types"`
		Compliance     struct {
			Compliant    int `json:"compliant"`
			NonCompliant int `json:"non_compliant"`
		} `json:"compliance"`
	}{
		TotalResources: len(results),
		ResourceTypes:  groupedResults,
	}

	// Calculate compliance stats
	for _, result := range results {
		if result.ComplianceTag == "compliant" {
			summary.Compliance.Compliant++
		} else {
			summary.Compliance.NonCompliant++
		}
	}

	// Marshal to JSON with indentation for readability
	jsonBytes, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal results to JSON: %w", err)
	}

	return string(jsonBytes), nil
}

// ValidateConfig validates the configuration file for the AWS Tag Inspector.
//
// This method performs comprehensive validation of the provided configuration file,
// ensuring its structural integrity, syntax, and compatibility with the AWS Tag Inspector.
//
// Parameters:
//   - ctx: Optional context for controlling the validation operation's lifecycle and timeout.
//   - config: A Dagger file representing the configuration to be validated.
//
// Returns:
//   - A string containing the validated configuration file's contents.
//   - An error if the configuration is invalid or cannot be processed.
func (m *AwsTagInspector) ValidateConfig(
	ctx context.Context,
	config *dagger.File,
) (string, error) {
	cfgLoader := newCfg()

	_, cfgErr := cfgLoader.loadConfig(ctx, config)
	if cfgErr != nil {
		return "", WrapError(cfgErr, "configuration file is invalid, failed to load it")
	}

	if err := cfgLoader.validateConfig(cfgLoader.config); err != nil {
		return "", WrapError(err, "configuration file is invalid, failed to validate it")
	}

	cfgContent, cfgContentErr := config.Contents(ctx)
	if cfgContentErr != nil {
		return "", WrapError(cfgContentErr, "failed to get configuration file's content")
	}

	fmt.Println("configuration file is valid ✅")

	return cfgContent, nil
}
