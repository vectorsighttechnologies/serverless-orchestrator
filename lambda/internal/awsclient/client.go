// Package awsclient provides a factory for AWS SDK v2 service clients.
//
// Following the DRY principle, all client creation is centralised here
// so handlers never directly call aws.Config. This also makes testing
// easier — you can swap the factory for mocks.
package awsclient

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// Factory creates and caches AWS service clients.
// Using sync.Once ensures clients are initialised exactly once (thread-safe).
type Factory struct {
	cfg     aws.Config
	once    sync.Once
	initErr error
}

// New creates a new Factory. AWS config is lazily loaded on first client request.
func New() *Factory {
	return &Factory{}
}

// init loads the AWS config once.
func (f *Factory) init(ctx context.Context) error {
	f.once.Do(func() {
		f.cfg, f.initErr = awsconfig.LoadDefaultConfig(ctx)
	})
	return f.initErr
}

// Lambda returns an AWS Lambda service client.
func (f *Factory) Lambda(ctx context.Context) (*lambda.Client, error) {
	if err := f.init(ctx); err != nil {
		return nil, err
	}
	return lambda.NewFromConfig(f.cfg), nil
}

// IAM returns an AWS IAM service client.
func (f *Factory) IAM(ctx context.Context) (*iam.Client, error) {
	if err := f.init(ctx); err != nil {
		return nil, err
	}
	return iam.NewFromConfig(f.cfg), nil
}

// CloudFormation returns an AWS CloudFormation service client.
func (f *Factory) CloudFormation(ctx context.Context) (*cloudformation.Client, error) {
	if err := f.init(ctx); err != nil {
		return nil, err
	}
	return cloudformation.NewFromConfig(f.cfg), nil
}

// CloudWatchLogs returns an AWS CloudWatch Logs service client.
func (f *Factory) CloudWatchLogs(ctx context.Context) (*cloudwatchlogs.Client, error) {
	if err := f.init(ctx); err != nil {
		return nil, err
	}
	return cloudwatchlogs.NewFromConfig(f.cfg), nil
}

// STS returns an AWS STS service client.
func (f *Factory) STS(ctx context.Context) (*sts.Client, error) {
	if err := f.init(ctx); err != nil {
		return nil, err
	}
	return sts.NewFromConfig(f.cfg), nil
}

// Region returns the configured AWS region.
func (f *Factory) Region(ctx context.Context) (string, error) {
	if err := f.init(ctx); err != nil {
		return "", err
	}
	return f.cfg.Region, nil
}
