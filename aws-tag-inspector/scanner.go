package main

import (
	"context"

	"github.com/Excoriate/daggerverse/aws-tag-inspector/internal/dagger"
)

func (m *AwsTagInspector) scanS3BucketsAPI(
	// ctx is the context for the scan function
	// +optional
	ctx context.Context,
	// tagCriteria is the criteria for the scan function
	tagCriteria TagCriteria,
) ([]ScanResult, error) {
	return nil, nil
}

func (m *AwsTagInspector) Scan(
	// ctx is the context for the scan function
	// +optional
	ctx context.Context,
	// dryRun is a boolean flag to indicate if the scan should be a dry run
	// +optional
	dryRun bool,
) (*dagger.File, error) {
	return nil, nil
}
