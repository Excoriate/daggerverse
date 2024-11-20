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
	if config.Version == "" {
		return fmt.Errorf("config version is required")
	}

	// Validate global configuration
	if config.Global.BatchSize != nil && *config.Global.BatchSize <= 0 {
		defaultBatchSize := 10
		config.Global.BatchSize = &defaultBatchSize
	}

	// Validate resource configurations
	for resourceType, resourceConfig := range config.Resources {
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

func (l *configLoader) validateComplianceLevels(levels map[string]ComplianceLevel) error {
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
