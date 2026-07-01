// Package newrelic — subscriptions.go
//
// CloudWatch log subscription filter management.
// Ported from: newrelic_lambda_cli/subscriptions.py
package newrelic

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"

	"github.com/vectorsight/serverless-tool/lambda/internal/awsclient"
)

const (
	// SubscriptionFilterName is the name used for NR log subscription filters.
	SubscriptionFilterName = "NewRelicLogStreaming"
)

// getLogGroupName builds the CloudWatch log group name from a function name or ARN.
// Ported from: subscriptions.py :: _get_log_group_name()
func getLogGroupName(functionName string) string {
	if strings.Contains(functionName, ":") {
		parts := strings.Split(functionName, ":")
		if len(parts) >= 7 {
			return "/aws/lambda/" + parts[6]
		}
	}
	return "/aws/lambda/" + functionName
}

// CreateLogSubscription creates a CloudWatch log subscription filter
// for the specified function, pointing to the NR log ingestion Lambda.
//
// Ported from: subscriptions.py :: create_log_subscription()
func CreateLogSubscription(
	ctx context.Context,
	clients *awsclient.Factory,
	functionName string,
	destinationARN string,
) error {
	cwClient, err := clients.CloudWatchLogs(ctx)
	if err != nil {
		return fmt.Errorf("failed to create CloudWatch Logs client: %w", err)
	}

	logGroupName := getLogGroupName(functionName)

	_, err = cwClient.PutSubscriptionFilter(ctx, &cloudwatchlogs.PutSubscriptionFilterInput{
		LogGroupName:   aws.String(logGroupName),
		FilterName:     aws.String(SubscriptionFilterName),
		FilterPattern:  aws.String(""),
		DestinationArn: aws.String(destinationARN),
	})
	if err != nil {
		return fmt.Errorf("failed to create log subscription for %s: %w", functionName, err)
	}

	return nil
}

// RemoveLogSubscription removes the NR log subscription filter.
//
// Ported from: subscriptions.py :: remove_log_subscription()
func RemoveLogSubscription(
	ctx context.Context,
	clients *awsclient.Factory,
	functionName string,
) error {
	cwClient, err := clients.CloudWatchLogs(ctx)
	if err != nil {
		return fmt.Errorf("failed to create CloudWatch Logs client: %w", err)
	}

	logGroupName := getLogGroupName(functionName)

	_, err = cwClient.DeleteSubscriptionFilter(ctx, &cloudwatchlogs.DeleteSubscriptionFilterInput{
		LogGroupName: aws.String(logGroupName),
		FilterName:   aws.String(SubscriptionFilterName),
	})
	if err != nil {
		// Ignore "not found" errors — the filter might not exist
		if strings.Contains(err.Error(), "ResourceNotFoundException") {
			return nil
		}
		return fmt.Errorf("failed to remove log subscription for %s: %w", functionName, err)
	}

	return nil
}
