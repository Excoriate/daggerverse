package main

import (
	"regexp"
)

// inspectorConfig represents the overall configuration structure for the AWS tag inspector.
// It contains global settings, resource-specific configurations, compliance levels,
// tag validation rules, and notification settings.
type inspectorConfig struct {
	Version          string                     `yaml:"version"`
	Global           globalConfig               `yaml:"global"`
	Resources        map[string]resourceConfig  `yaml:"resources"`
	ComplianceLevels map[string]complianceLevel `yaml:"compliance_levels"`
	TagValidation    tagValidation              `yaml:"tag_validation"`
	Notifications    notificationConfig         `yaml:"notifications"`
}

// globalConfig defines the default configuration settings that apply across all resources.
// It includes batch processing size, required and forbidden tags, and specific tag requirements.
type globalConfig struct {
	Enabled     bool        `yaml:"enabled"`
	BatchSize   *int        `yaml:"batch_size,omitempty"`
	TagCriteria TagCriteria `yaml:"tag_criteria"`
}

// resourceConfig provides configuration specific to individual resource types.
// It allows for more granular control over tag requirements, exclusions, and processing.
type resourceConfig struct {
	Enabled           bool               `yaml:"enabled"`
	BatchSize         *int               `yaml:"batch_size,omitempty"`
	TagCriteria       TagCriteria        `yaml:"tag_criteria"`
	ExcludedResources []excludedResource `yaml:"excluded_resources"`
}

// excludedResource defines a specific resource to be excluded from tag inspection,
// with a pattern to match and a reason for exclusion.
type excludedResource struct {
	Pattern string `yaml:"pattern"`
	Reason  string `yaml:"reason"`
}

// complianceLevel specifies the tag requirements for achieving a particular
// compliance status or level within the tag inspection process.
type complianceLevel struct {
	RequiredTags []string          `yaml:"required_tags"`
	SpecificTags map[string]string `yaml:"specific_tags"`
}

// tagValidation contains rules for validating tags across resources.
// It includes allowed values for specific tags and pattern-based validation rules.
type tagValidation struct {
	AllowedValues map[string][]string `yaml:"allowed_values"`
	PatternRules  map[string]string   `yaml:"pattern_rules"`
	compiledRules map[string]*regexp.Regexp
}

// notificationConfig manages the notification settings for reporting
// tag inspection results through different channels.
type notificationConfig struct {
	Slack slackConfig `yaml:"slack"`
	Email emailConfig `yaml:"email"`
}

// slackConfig defines the configuration for Slack notifications,
// including whether they are enabled and which channels to use.
type slackConfig struct {
	Enabled  bool              `yaml:"enabled"`
	Channels map[string]string `yaml:"channels"`
}

// emailConfig specifies the email notification settings,
// including whether email notifications are enabled and the list of recipients.
type emailConfig struct {
	Enabled    bool     `yaml:"enabled"`
	Recipients []string `yaml:"recipients"`
}
