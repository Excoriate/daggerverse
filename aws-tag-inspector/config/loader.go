package config

import (
	"fmt"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

// ConfigLoader handles loading and validating configuration
type ConfigLoader struct {
	config *Config
}

// NewConfigLoader creates a new configuration loader
func NewConfigLoader() *ConfigLoader {
	return &ConfigLoader{}
}

// LoadConfig loads configuration from a YAML file
func (l *ConfigLoader) LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := l.validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	if err := l.compilePatternRules(&config); err != nil {
		return nil, fmt.Errorf("failed to compile pattern rules: %w", err)
	}

	l.config = &config
	return &config, nil
}

// validateConfig performs validation of the loaded configuration
func (l *ConfigLoader) validateConfig(config *Config) error {
	if config.Version == "" {
		return fmt.Errorf("config version is required")
	}

	// Validate global configuration
	if config.Global.BatchSize <= 0 {
		config.Global.BatchSize = 10 // Default batch size
	}

	// Validate resource configurations
	for resourceType, resourceConfig := range config.Resources {
		if resourceConfig.BatchSize != nil && *resourceConfig.BatchSize <= 0 {
			return fmt.Errorf("invalid batch size for resource type %s", resourceType)
		}

		// Validate excluded resource patterns
		for _, excluded := range resourceConfig.ExcludedResources {
			if _, err := regexp.Compile(excluded.Pattern); err != nil {
				return fmt.Errorf("invalid exclusion pattern %s for resource type %s: %w",
					excluded.Pattern, resourceType, err)
			}
		}
	}

	return nil
}

// compilePatternRules pre-compiles regex patterns for tag validation
func (l *ConfigLoader) compilePatternRules(config *Config) error {
	config.TagValidation.compiledRules = make(map[string]*regexp.Regexp)

	for tagName, pattern := range config.TagValidation.PatternRules {
		compiled, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("invalid pattern for tag %s: %w", tagName, err)
		}
		config.TagValidation.compiledRules[tagName] = compiled
	}

	return nil
}

// GetTagCriteria generates TagCriteria for a specific resource type
func (l *ConfigLoader) GetTagCriteria(resourceType string) TagCriteria {
	if l.config == nil {
		return TagCriteria{}
	}

	// Start with global criteria
	criteria := TagCriteria{
		RequiredTags:  append([]string{}, l.config.Global.RequiredTags...),
		ForbiddenTags: append([]string{}, l.config.Global.ForbiddenTags...),
		SpecificTags:  make(map[string]string),
	}

	// Copy global specific tags
	for k, v := range l.config.Global.SpecificTags {
		criteria.SpecificTags[k] = v
	}

	// Merge resource-specific criteria if it exists
	if resourceConfig, exists := l.config.Resources[resourceType]; exists && resourceConfig.Enabled {
		criteria.RequiredTags = append(criteria.RequiredTags, resourceConfig.RequiredTags...)
		criteria.ForbiddenTags = append(criteria.ForbiddenTags, resourceConfig.ForbiddenTags...)

		// Merge specific tags (resource-specific tags override global tags)
		for k, v := range resourceConfig.SpecificTags {
			criteria.SpecificTags[k] = v
		}
	}

	return criteria
}

// IsResourceExcluded checks if a resource should be excluded from scanning
func (l *ConfigLoader) IsResourceExcluded(resourceType, resourceID string) (bool, string) {
	if l.config == nil {
		return false, ""
	}

	resourceConfig, exists := l.config.Resources[resourceType]
	if !exists || !resourceConfig.Enabled {
		return false, ""
	}

	for _, excluded := range resourceConfig.ExcludedResources {
		if matched, _ := regexp.MatchString(excluded.Pattern, resourceID); matched {
			return true, excluded.Reason
		}
	}

	return false, ""
}
