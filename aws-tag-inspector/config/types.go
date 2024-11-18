package config

import (
	"regexp"
)

type Config struct {
	Version          string                     `yaml:"version"`
	Global           GlobalConfig               `yaml:"global"`
	Resources        map[string]ResourceConfig  `yaml:"resources"`
	ComplianceLevels map[string]ComplianceLevel `yaml:"compliance_levels"`
	TagValidation    TagValidation              `yaml:"tag_validation"`
	Notifications    NotificationConfig         `yaml:"notifications"`
}

type GlobalConfig struct {
	BatchSize     int               `yaml:"batch_size"`
	RequiredTags  []string          `yaml:"required_tags"`
	ForbiddenTags []string          `yaml:"forbidden_tags"`
	SpecificTags  map[string]string `yaml:"specific_tags"`
}

type ResourceConfig struct {
	Enabled           bool               `yaml:"enabled"`
	BatchSize         *int               `yaml:"batch_size,omitempty"`
	RequiredTags      []string           `yaml:"required_tags"`
	ForbiddenTags     []string           `yaml:"forbidden_tags"`
	SpecificTags      map[string]string  `yaml:"specific_tags"`
	ExcludedResources []ExcludedResource `yaml:"excluded_resources"`
}

type ExcludedResource struct {
	Pattern string `yaml:"pattern"`
	Reason  string `yaml:"reason"`
}

type ComplianceLevel struct {
	RequiredTags []string          `yaml:"required_tags"`
	SpecificTags map[string]string `yaml:"specific_tags"`
}

type TagValidation struct {
	AllowedValues map[string][]string `yaml:"allowed_values"`
	PatternRules  map[string]string   `yaml:"pattern_rules"`
	compiledRules map[string]*regexp.Regexp
}

type NotificationConfig struct {
	Slack SlackConfig `yaml:"slack"`
	Email EmailConfig `yaml:"email"`
}

type SlackConfig struct {
	Enabled  bool              `yaml:"enabled"`
	Channels map[string]string `yaml:"channels"`
}

type EmailConfig struct {
	Enabled    bool     `yaml:"enabled"`
	Recipients []string `yaml:"recipients"`
}
