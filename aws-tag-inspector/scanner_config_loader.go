package main

import (
	"context"
	"fmt"
	"regexp"

	"github.com/Excoriate/daggerverse/aws-tag-inspector/internal/dagger"
	"gopkg.in/yaml.v3"
)

// configLoader handles loading and validating configuration
type configLoader struct {
	config *inspectorConfig
}

func newCfg() *configLoader {
	return &configLoader{}
}

func (l *configLoader) loadConfig(ctx context.Context, cfg *dagger.File) (*inspectorConfig, error) {
	if cfg == nil {
		return nil, Errorf("configuration file is required")
	}

	fileContent, fileContentErr := cfg.Contents(ctx)
	if fileContentErr != nil {
		return nil, WrapError(fileContentErr, "failed to retrieve configuration file contents")
	}

	parsedCfg := &inspectorConfig{}
	cfgErr := yaml.Unmarshal([]byte(fileContent), parsedCfg)

	if cfgErr != nil {
		return nil, WrapError(cfgErr, "failed to parse configuration file contents")
	}

	if err := l.validateConfig(parsedCfg); err != nil {
		return nil, WrapError(err, "failed to validate configuration")
	}

	if err := l.compilePatternRules(parsedCfg); err != nil {
		return nil, WrapError(err, "failed to compile pattern rules")
	}

	l.config = parsedCfg
	return parsedCfg, nil
}

// validateConfig performs validation of the loaded configuration
func (l *configLoader) validateConfig(config *inspectorConfig) error {
	// Validate version
	if config.Version == "" {
		return fmt.Errorf("config version is required")
	}

	// Validate global configuration
	if err := l.validateGlobalConfig(&config.Global); err != nil {
		return fmt.Errorf("global configuration validation failed: %w", err)
	}

	// Validate resource configurations
	if err := l.validateResourceConfigs(config.Resources); err != nil {
		return fmt.Errorf("resource configuration validation failed: %w", err)
	}

	// Validate compliance levels
	if err := l.validateComplianceLevels(config.ComplianceLevels); err != nil {
		return fmt.Errorf("compliance levels validation failed: %w", err)
	}

	// Validate tag validation rules
	if err := l.validateTagValidationRules(config.TagValidation); err != nil {
		return fmt.Errorf("tag validation rules validation failed: %w", err)
	}

	// Validate notifications
	if err := l.validateNotifications(config.Notifications); err != nil {
		return fmt.Errorf("notifications validation failed: %w", err)
	}

	return nil
}

// validateGlobalConfig validates the global configuration
func (l *configLoader) validateGlobalConfig(global *globalConfig) error {
	// Validate batch size
	if global.BatchSize != nil && *global.BatchSize <= 0 {
		defaultBatchSize := 10
		global.BatchSize = &defaultBatchSize
	}

	// Validate tag criteria
	if err := l.validateTagCriteria(global.TagCriteria); err != nil {
		return fmt.Errorf("global tag criteria validation failed: %w", err)
	}

	return nil
}

// validateResourceConfigs validates resource-specific configurations
func (l *configLoader) validateResourceConfigs(resources map[string]resourceConfig) error {
	for resourceType, resourceConfig := range resources {
		// Validate batch size
		if resourceConfig.BatchSize != nil && *resourceConfig.BatchSize <= 0 {
			return fmt.Errorf("invalid batch size for resource type %s", resourceType)
		}

		// Validate tag criteria
		if err := l.validateTagCriteria(resourceConfig.TagCriteria); err != nil {
			return fmt.Errorf("invalid tag criteria for resource type %s: %w", resourceType, err)
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

// validateTagValidationRules validates the tag validation configuration
func (l *configLoader) validateTagValidationRules(tagValidation tagValidation) error {
	// Validate allowed values
	for tagName, values := range tagValidation.AllowedValues {
		if len(values) == 0 {
			return fmt.Errorf("no allowed values specified for tag %s", tagName)
		}

		// Validate each value is not empty
		for _, value := range values {
			if value == "" {
				return fmt.Errorf("empty value found in allowed values for tag %s", tagName)
			}
		}
	}

	// Validate pattern rules
	for tagName, pattern := range tagValidation.PatternRules {
		if pattern == "" {
			return fmt.Errorf("empty pattern rule for tag %s", tagName)
		}
		if _, err := regexp.Compile(pattern); err != nil {
			return fmt.Errorf("invalid regex pattern for tag %s: %w", tagName, err)
		}
	}

	return nil
}

// validateNotifications validates the notification configuration
func (l *configLoader) validateNotifications(notifications notificationConfig) error {
	// Validate Slack notifications
	if notifications.Slack.Enabled {
		if len(notifications.Slack.Channels) == 0 {
			return fmt.Errorf("Slack notifications enabled but no channels specified")
		}

		// Validate each channel name
		for _, channel := range notifications.Slack.Channels {
			if channel == "" {
				return fmt.Errorf("empty Slack channel name found")
			}
		}
	}

	// Validate Email notifications
	if notifications.Email.Enabled {
		if len(notifications.Email.Recipients) == 0 {
			return fmt.Errorf("email notifications enabled but no recipients specified")
		}

		// Validate email format for each recipient
		emailRegex := regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
		for _, email := range notifications.Email.Recipients {
			if !emailRegex.MatchString(email) {
				return fmt.Errorf("invalid email format: %s", email)
			}
		}

		// Validate frequency
		validFrequencies := []string{"daily", "hourly", "weekly"}
		if notifications.Email.Frequency != "" {
			found := false
			for _, freq := range validFrequencies {
				if notifications.Email.Frequency == freq {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("invalid email notification frequency: %s", notifications.Email.Frequency)
			}
		}
	}

	return nil
}

// compilePatternRules pre-compiles regex patterns for tag validation
func (l *configLoader) compilePatternRules(config *inspectorConfig) error {
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

// getTagCriteria generates TagCriteria for a specific resource type
func (l *configLoader) getTagCriteria(resourceType string) TagCriteria {
	if l.config == nil {
		return TagCriteria{}
	}

	// Start with global criteria
	criteria := l.config.Global.TagCriteria

	// Merge resource-specific criteria if it exists
	if resourceConfig, exists := l.config.Resources[resourceType]; exists && resourceConfig.Enabled {
		// Merge required tags
		if len(resourceConfig.TagCriteria.RequiredTags) > 0 {
			criteria.RequiredTags = append(criteria.RequiredTags, resourceConfig.TagCriteria.RequiredTags...)
		}

		// Merge forbidden tags
		if len(resourceConfig.TagCriteria.ForbiddenTags) > 0 {
			criteria.ForbiddenTags = append(criteria.ForbiddenTags, resourceConfig.TagCriteria.ForbiddenTags...)
		}

		// Merge specific tags (resource-specific tags override global tags)
		if len(resourceConfig.TagCriteria.SpecificTags) > 0 {
			if criteria.SpecificTags == nil {
				criteria.SpecificTags = make(map[string]string)
			}
			for k, v := range resourceConfig.TagCriteria.SpecificTags {
				criteria.SpecificTags[k] = v
			}
		}
	}

	return criteria
}

// isResourceExcluded checks if a resource should be excluded from scanning
func (l *configLoader) isResourceExcluded(resourceType, resourceID string) (bool, string) {
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

// validateTagCriteria validates the tag criteria for a resource
func (l *configLoader) validateTagCriteria(criteria TagCriteria) error {
	if err := l.validateMinimumRequiredTags(criteria); err != nil {
		return err
	}

	// Validate required tags
	for _, tag := range criteria.RequiredTags {
		if tag == "" {
			return fmt.Errorf("empty required tag found")
		}
	}

	// Validate forbidden tags
	for _, tag := range criteria.ForbiddenTags {
		if tag == "" {
			return fmt.Errorf("empty forbidden tag found")
		}
	}

	// Validate specific tags
	for k, v := range criteria.SpecificTags {
		if k == "" {
			return fmt.Errorf("empty specific tag key found")
		}
		if v == "" {
			return fmt.Errorf("empty specific tag value found for key: %s", k)
		}
	}

	return nil
}

func (l *configLoader) validateComplianceLevels(levels map[string]complianceLevel) error {
	for levelName, level := range levels {
		if levelName == "" {
			return Errorf("compliance level name cannot be empty")
		}
		// Validate required tags
		for _, tag := range level.RequiredTags {
			if tag == "" {
				return fmt.Errorf("empty required tag in compliance level '%s'", levelName)
			}
		}
		// Validate specific tags
		for key, value := range level.SpecificTags {
			if key == "" || value == "" {
				return fmt.Errorf("empty key or value in specific tags of compliance level '%s'", levelName)
			}
		}
	}
	return nil
}

// Add validateMinimumRequiredTags function
func (l *configLoader) validateMinimumRequiredTags(criteria TagCriteria) error {
	if criteria.MinimumRequiredTags < 0 {
		return fmt.Errorf("minimum required tags cannot be negative")
	}

	if criteria.MinimumRequiredTags > len(criteria.RequiredTags) {
		return fmt.Errorf("minimum required tags (%d) cannot be greater than the number of required tags (%d)",
			criteria.MinimumRequiredTags, len(criteria.RequiredTags))
	}

	return nil
}
