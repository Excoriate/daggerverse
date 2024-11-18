package main

import (
	"github.com/your-project/config"
)

// TagCriteria defines the rules for tag validation
type TagCriteria struct {
	RequiredTags  []string          // Tags that must be present
	ForbiddenTags []string          // Tags that must not be present
	SpecificTags  map[string]string // Tags that must match exact values
}

// ScanResult represents the result of scanning a resource
type ScanResult struct {
	// ResourceType is the AWS resource type (e.g., "ec2:instance", "s3:bucket")
	ResourceType string `json:"resource_type"`

	// ResourceID is the unique identifier for the resource
	ResourceID string `json:"resource_id"`

	// ARN is the Amazon Resource Name
	ARN string `json:"arn"`

	// Region is the AWS region where the resource is located
	Region string `json:"region"`

	// Tags is a map of key-value pairs representing the resource's tags
	Tags map[string]string `json:"tags"`

	// Issues contains a list of compliance violations or problems found
	Issues []string `json:"issues"`

	// Metadata contains additional resource-specific information
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// ComplianceTag indicates the overall compliance status
	ComplianceTag string `json:"compliance_tag,omitempty"`
}

type Scanner interface {
	// Identity information
	GetResourceType() string
	GetResourceID() string
	GetARN() string
	GetRegion() string
	// Tag operations
	GetTags() map[string]string
	HasTags() bool
	HasTag(key string) bool
	GetTagValue(key string) (string, bool)

	// Scan operations
	Scan(criteria TagCriteria) ([]ScanResult, error)

	// Additional required methods
	GetMetadata() map[string]interface{}
	ValidateCompliance(criteria TagCriteria) bool
	ScanTags(criteria TagCriteria) []string

	// IsExcluded checks if the resource should be excluded from scanning
	IsExcluded() (bool, string)
}

// BaseResource provides a common implementation for AWS resources
type BaseResource struct {
	ResourceType string
	ResourceID   string
	ARN          string
	Region       string
	Tags         map[string]string
	Metadata     map[string]interface{}
}

// GetResourceType returns the type of the resource
func (r *BaseResource) GetResourceType() string {
	return r.ResourceType
}

// GetResourceID returns the resource identifier
func (r *BaseResource) GetResourceID() string {
	return r.ResourceID
}

// GetARN returns the ARN of the resource
func (r *BaseResource) GetARN() string {
	return r.ARN
}

// GetRegion returns the region of the resource
func (r *BaseResource) GetRegion() string {
	return r.Region
}

// GetTags returns all tags associated with the resource
func (r *BaseResource) GetTags() map[string]string {
	if r.Tags == nil {
		return make(map[string]string)
	}
	return r.Tags
}

// HasTags returns true if the resource has any tags
func (r *BaseResource) HasTags() bool {
	return len(r.Tags) > 0
}

// HasTag checks if a specific tag exists
func (r *BaseResource) HasTag(key string) bool {
	_, exists := r.Tags[key]
	return exists
}

// GetTagValue gets the value of a specific tag
func (r *BaseResource) GetTagValue(key string) (string, bool) {
	if r.Tags == nil {
		return "", false
	}
	value, exists := r.Tags[key]
	return value, exists
}

// ScanTags performs tag validation against criteria
func (r *BaseResource) ScanTags(criteria TagCriteria) []string {
	var issues []string

	// Validate input criteria
	if criteria.RequiredTags == nil && criteria.ForbiddenTags == nil && criteria.SpecificTags == nil {
		return []string{"Invalid criteria: no validation rules specified"}
	}

	// Check for untagged resource
	if !r.HasTags() {
		issues = append(issues, "Resource has no tags")
		return issues
	}

	// Check required tags
	for _, required := range criteria.RequiredTags {
		if required == "" {
			continue // Skip empty required tags
		}
		if !r.HasTag(required) {
			issues = append(issues, "Missing required tag: "+required)
		}
	}

	// Check forbidden tags
	for _, forbidden := range criteria.ForbiddenTags {
		if forbidden == "" {
			continue // Skip empty forbidden tags
		}
		if r.HasTag(forbidden) {
			issues = append(issues, "Contains forbidden tag: "+forbidden)
		}
	}

	// Check specific tag values
	for key, expectedValue := range criteria.SpecificTags {
		if key == "" {
			continue // Skip empty keys
		}
		if actualValue, exists := r.GetTagValue(key); !exists || actualValue != expectedValue {
			issues = append(issues, "Tag mismatch: "+key+" should be "+expectedValue)
		}
	}

	return issues
}

// ValidateCompliance checks if the resource is compliant with tag criteria
func (r *BaseResource) ValidateCompliance(criteria TagCriteria) bool {
	return len(r.ScanTags(criteria)) == 0
}

// GetMetadata returns the resource metadata
func (r *BaseResource) GetMetadata() map[string]interface{} {
	if r.Metadata == nil {
		return make(map[string]interface{})
	}
	return r.Metadata
}

// IsExcluded checks if the resource should be excluded from scanning
func (r *BaseResource) IsExcluded(configLoader *config.ConfigLoader) (bool, string) {
	return configLoader.IsResourceExcluded(r.ResourceType, r.ResourceID)
}
