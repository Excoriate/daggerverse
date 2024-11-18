package main

import (
	"context"
)

func (m *AwsTagInspector) ScanS3Buckets(
	ctx context.Context,
	criteria TagCriteria,
) ([]ScanResult, error) {
}
