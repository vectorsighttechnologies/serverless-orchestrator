// Package newrelic — integrations.go
//
// AWS Integration management via CloudFormation stacks.
//
// Supports two methods:
//   - Metric Streams (recommended) — deploys NR's official nested CF templates
//   - API Polling (legacy) — creates an IAM role + NerdGraph link
package newrelic

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cftypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	"github.com/vectorsighttechnologies/serverless-orchestrator/lambda/internal/awsclient"
	"github.com/vectorsighttechnologies/serverless-orchestrator/lambda/internal/config"
	"github.com/vectorsighttechnologies/serverless-orchestrator/lambda/internal/types"
)

const (
	// MetricStreamsStackName is the CF stack name for Metric Streams integration.
	MetricStreamsStackName = "NewRelicMetricStreams"

	// APIPollingStackName is the CF stack name for API Polling integration.
	APIPollingStackName = "NewRelicLambdaIntegrationRole"

	// MetricStreamsTemplateURL is NR's official nested CF template for Metric Streams.
	MetricStreamsTemplateURL = "https://nr-downloads-main.s3.amazonaws.com/cloud_integrations/aws/cloudformation/newrelic-cloudformation-mstreams.yml"

	// APIPollingTemplateURL is NR's official CF template for API Polling integration.
	APIPollingTemplateURL = "https://nr-downloads-main.s3.amazonaws.com/cloud_integrations/aws/cloudformation/newrelic-cloudformation-polling.yml"
)

// GetIntegrationStatus checks the CF stack status for the NR integration.
//
func GetIntegrationStatus(
	ctx context.Context,
	clients *awsclient.Factory,
	cfg *config.Config,
) (*types.IntegrationStatusResponse, error) {
	cfClient, err := clients.CloudFormation(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create CloudFormation client: %w", err)
	}

	// Try Metric Streams first, then API Polling
	for _, stackName := range []string{MetricStreamsStackName, APIPollingStackName} {
		// API Polling stack name has AWS Account ID suffix
		resolvedStackName := stackName
		if stackName == APIPollingStackName {
			// Get AWS account ID to build stack name
			stsClient, err := clients.STS(ctx)
			if err == nil {
				callerIdentity, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
				if err == nil {
					resolvedStackName = APIPollingStackName + "_" + aws.ToString(callerIdentity.Account)
				}
			}
		}

		resp, err := cfClient.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{
			StackName: aws.String(resolvedStackName),
		})
		if err != nil {
			// Stack not found — continue checking the other
			if isStackNotFound(err) {
				continue
			}
			return nil, fmt.Errorf("failed to describe stack %s: %w", resolvedStackName, err)
		}

		if len(resp.Stacks) > 0 {
			stack := resp.Stacks[0]
			status := string(stack.StackStatus)

			method := "api_polling"
			if stackName == MetricStreamsStackName {
				method = "metric_streams"
			}

			result := &types.IntegrationStatusResponse{
				StackName:   resolvedStackName,
				StackStatus: status,
				Method:      method,
				LastChecked: time.Now().UTC().Format(time.RFC3339),
			}

			switch {
			case isActiveStatus(stack.StackStatus):
				result.Status = "active"

				if method == "metric_streams" {
					// Metric Streams stack (newrelic-cloudformation-mstreams.yml) links the account automatically via custom resource
					// and has no outputs. So we can just set status to active.
					break
				}

				// Retrieve IAM Role ARN from stack outputs
				roleArn := getRoleArnFromStack(&stack)
				if roleArn != "" {
					if cfg.HasAPIKey() {
						err := checkAndLinkAccount(ctx, cfg, roleArn, method)
						if err != nil {
							result.Status = "error"
							result.Error = "AWS integration stack is active, but failed to link with New Relic account: " + err.Error()
						}
					} else {
						result.Status = "error"
						result.Error = "AWS integration stack is active, but New Relic User API Key is missing. Please configure it in settings."
					}
				} else {
					result.Status = "error"
					result.Error = "AWS integration stack is active, but IAM Role ARN could not be resolved from outputs."
				}
			case isInProgressStatus(stack.StackStatus):
				result.Status = "in_progress"
			default:
				result.Status = "error"
				if stack.StackStatusReason != nil {
					result.Error = *stack.StackStatusReason
				}
			}

			return result, nil
		}
	}

	// No stack found
	return &types.IntegrationStatusResponse{
		Status:      "not_setup",
		LastChecked: time.Now().UTC().Format(time.RFC3339),
	}, nil
}

// SetupIntegration deploys the NR integration CF stack.
//
// For Metric Streams: deploys the nested CF template with Firehose, S3, and metric stream.
// For API Polling: creates the integration IAM role.
func SetupIntegration(
	ctx context.Context,
	clients *awsclient.Factory,
	cfg *config.Config,
	req *types.IntegrationSetupRequest,
) (*types.IntegrationSetupResponse, error) {
	cfClient, err := clients.CloudFormation(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create CloudFormation client: %w", err)
	}

	// Get AWS account ID for stack parameters
	stsClient, err := clients.STS(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create STS client: %w", err)
	}

	callerIdentity, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to get AWS account ID: %w", err)
	}

	awsAccountID := aws.ToString(callerIdentity.Account)

	var stackName, templateURL string
	var params []cftypes.Parameter

	switch req.Method {
	case "metric_streams":
		if cfg.APIKey == "" {
			return nil, fmt.Errorf("New Relic User API Key is required for Metric Streams integration. Please configure it in settings.")
		}

		stackName = MetricStreamsStackName
		templateURL = MetricStreamsTemplateURL
		
		region := "US"
		if strings.ToLower(cfg.Region) == "eu" {
			region = "EU"
		} else if strings.ToLower(cfg.Region) == "jp" {
			region = "JP"
		}

		params = []cftypes.Parameter{
			{ParameterKey: aws.String("NewRelicAccountId"), ParameterValue: aws.String(cfg.AccountID)},
			{ParameterKey: aws.String("NewRelicAPIKey"), ParameterValue: aws.String(cfg.APIKey)},
			{ParameterKey: aws.String("NewRelicLicenseKey"), ParameterValue: aws.String(cfg.LicenseKey)},
			{ParameterKey: aws.String("NewRelicRegion"), ParameterValue: aws.String(region)},
			{ParameterKey: aws.String("IntegrationName"), ParameterValue: aws.String("serverless-orchestrator-metric-streams-" + awsAccountID)},
			{ParameterKey: aws.String("PollingIntegrationSlugs"), ParameterValue: aws.String("lambda")},
			{ParameterKey: aws.String("MetricCollectionMode"), ParameterValue: aws.String("PUSH")},
		}

	case "api_polling":
		if cfg.APIKey == "" {
			return nil, fmt.Errorf("New Relic User API Key is required for API Polling integration. Please configure it in settings.")
		}

		stackName = APIPollingStackName + "_" + awsAccountID
		templateURL = APIPollingTemplateURL
		
		region := "US"
		if strings.ToLower(cfg.Region) == "eu" {
			region = "EU"
		}

		params = []cftypes.Parameter{
			{ParameterKey: aws.String("NewRelicAccountId"), ParameterValue: aws.String(cfg.AccountID)},
			{ParameterKey: aws.String("NewRelicAPIKey"), ParameterValue: aws.String(cfg.APIKey)},
			{ParameterKey: aws.String("NewRelicRegion"), ParameterValue: aws.String(region)},
			{ParameterKey: aws.String("IntegrationName"), ParameterValue: aws.String("serverless-orchestrator-integration-" + awsAccountID)},
			{ParameterKey: aws.String("PollingIntegrationSlugs"), ParameterValue: aws.String("lambda")},
		}

	default:
		return nil, fmt.Errorf("unsupported integration method: %s", req.Method)
	}

	// Check if stack already exists
	describeResp, err := cfClient.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{
		StackName: aws.String(stackName),
	})
	if err == nil && len(describeResp.Stacks) > 0 {
		existingStack := describeResp.Stacks[0]
		status := existingStack.StackStatus

		// If stack is in a failed/rollbacked state, delete it first
		if isFailedStatus(status) {
			_, err = cfClient.DeleteStack(ctx, &cloudformation.DeleteStackInput{
				StackName: aws.String(stackName),
			})
			if err != nil {
				return nil, fmt.Errorf("failed to delete failed existing stack %s: %w. You can delete it manually via AWS CLI: aws cloudformation delete-stack --stack-name %s --retain-resources LogGroupManagerFunction", stackName, err, stackName)
			}

			// Wait for deletion to complete (up to 20 seconds)
			deleted := false
			for i := 0; i < 10; i++ {
				time.Sleep(2 * time.Second)
				_, err = cfClient.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{
					StackName: aws.String(stackName),
				})
				if err != nil && isStackNotFound(err) {
					deleted = true
					break
				}
			}
			if !deleted {
				return nil, fmt.Errorf("timeout waiting for deletion of failed existing stack %s. If it is stuck in ROLLBACK_FAILED, please delete it manually via AWS CLI: aws cloudformation delete-stack --stack-name %s --retain-resources LogGroupManagerFunction", stackName, stackName)
			}
		} else if status == cftypes.StackStatusDeleteInProgress {
			// Wait for deletion to complete
			deleted := false
			for i := 0; i < 10; i++ {
				time.Sleep(2 * time.Second)
				_, err = cfClient.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{
					StackName: aws.String(stackName),
				})
				if err != nil && isStackNotFound(err) {
					deleted = true
					break
				}
			}
			if !deleted {
				return nil, fmt.Errorf("another deletion of stack %s is in progress and timed out. You may force delete it manually via AWS CLI: aws cloudformation delete-stack --stack-name %s --retain-resources LogGroupManagerFunction", stackName, stackName)
			}
		} else if isInProgressStatus(status) {
			return &types.IntegrationSetupResponse{
				Status:  "in_progress",
				StackID: aws.ToString(existingStack.StackId),
			}, nil
		} else if !isFailedStatus(status) && status != cftypes.StackStatusDeleteComplete {
			// Stack is active/complete — no need to recreate it
			return &types.IntegrationSetupResponse{
				Status:  "active",
				StackID: aws.ToString(existingStack.StackId),
			}, nil
		}
	}

	// Create the CF stack
	createResp, err := cfClient.CreateStack(ctx, &cloudformation.CreateStackInput{
		StackName:   aws.String(stackName),
		TemplateURL: aws.String(templateURL),
		Parameters:  params,
		Capabilities: []cftypes.Capability{
			cftypes.CapabilityCapabilityIam,
			cftypes.CapabilityCapabilityNamedIam,
			cftypes.CapabilityCapabilityAutoExpand,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create CF stack %s: %w", stackName, err)
	}

	return &types.IntegrationSetupResponse{
		Status:  "in_progress",
		StackID: aws.ToString(createResp.StackId),
	}, nil
}

// RemoveIntegration deletes the NR integration CF stack and unlinks it from New Relic.
//
func RemoveIntegration(
	ctx context.Context,
	clients *awsclient.Factory,
	cfg *config.Config,
) error {
	cfClient, err := clients.CloudFormation(ctx)
	if err != nil {
		return fmt.Errorf("failed to create CloudFormation client: %w", err)
	}

	// Get AWS account ID
	stsClient, err := clients.STS(ctx)
	if err != nil {
		return fmt.Errorf("failed to create STS client: %w", err)
	}
	callerIdentity, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return fmt.Errorf("failed to get AWS account ID: %w", err)
	}
	awsAccountID := aws.ToString(callerIdentity.Account)

	// Unlink the AWS account from New Relic if User API Key is configured
	if cfg.HasAPIKey() {
		err := unlinkAccount(ctx, cfg, awsAccountID)
		if err != nil {
			fmt.Printf("warning: failed to unlink account from New Relic: %v\n", err)
		}
	}

	// Delete both stack types
	stackNames := []string{
		MetricStreamsStackName,
		APIPollingStackName + "_" + awsAccountID,
	}

	for _, name := range stackNames {
		_, err := cfClient.DeleteStack(ctx, &cloudformation.DeleteStackInput{
			StackName: aws.String(name),
		})
		if err != nil {
			if isStackNotFound(err) {
				continue
			}
			return fmt.Errorf("failed to delete stack %s: %w", name, err)
		}
	}

	return nil
}

// ─────────────────────────────────────────────────────────────
// NerdGraph Cloud Linking Helpers
// ─────────────────────────────────────────────────────────────

// getRoleArnFromStack scans CF outputs to locate any output formatted as an IAM Role ARN.
func getRoleArnFromStack(stack *cftypes.Stack) string {
	for _, output := range stack.Outputs {
		val := aws.ToString(output.OutputValue)
		if strings.HasPrefix(val, "arn:aws:iam::") && strings.Contains(val, ":role/") {
			return val
		}
	}
	return ""
}

// checkAndLinkAccount verifies if AWS account link exists in NR, and creates it if not.
func checkAndLinkAccount(ctx context.Context, cfg *config.Config, roleArn string, method string) error {
	client := NewNerdGraphClient(cfg)

	// Extract AWS Account ID from Role ARN
	parts := strings.Split(roleArn, ":")
	if len(parts) < 5 {
		return fmt.Errorf("invalid Role ARN: %s", roleArn)
	}
	awsAccountID := parts[4]

	acctIDInt := 0
	_, err := fmt.Sscan(cfg.AccountID, &acctIDInt)
	if err != nil {
		return fmt.Errorf("invalid New Relic Account ID: %s", cfg.AccountID)
	}

	// Query existing linked accounts
	query := fmt.Sprintf(`
{
  actor {
    account(id: %d) {
      cloud {
        linkedAccounts {
          id
          name
          externalId
        }
      }
    }
  }
}`, acctIDInt)

	var queryResp CloudQueryResponse
	err = client.Query(ctx, query, nil, &queryResp)
	if err != nil {
		return fmt.Errorf("failed to query linked accounts: %w", err)
	}

	for _, la := range queryResp.Actor.Account.Cloud.LinkedAccounts {
		if la.ExternalID == awsAccountID {
			// Link already exists
			return nil
		}
	}

	// Establish link
	linkName := "serverless-orchestrator-aws-" + awsAccountID
	mode := "PULL"
	if method == "metric_streams" {
		mode = "PUSH"
	}

	mutation := fmt.Sprintf(`
mutation {
  cloudLinkAccount(
    accountId: %d,
    accounts: {
      aws: [{
        name: "%s",
        arn: "%s",
        metricCollectionMode: %s
      }]
    }
  ) {
    linkedAccounts {
      id
      name
    }
    errors {
      message
    }
  }
}`, acctIDInt, linkName, roleArn, mode)

	var mutationResp struct {
		CloudLinkAccount struct {
			LinkedAccounts []struct {
				ID   interface{} `json:"id"`
				Name string      `json:"name"`
			} `json:"linkedAccounts"`
			Errors []struct {
				Message string `json:"message"`
			} `json:"errors"`
		} `json:"cloudLinkAccount"`
	}

	err = client.Query(ctx, mutation, nil, &mutationResp)
	if err != nil {
		return fmt.Errorf("failed to link account: %w", err)
	}

	if len(mutationResp.CloudLinkAccount.Errors) > 0 {
		return fmt.Errorf("failed to link account: %s", mutationResp.CloudLinkAccount.Errors[0].Message)
	}

	return nil
}

// unlinkAccount unlinks the AWS cloud link from the New Relic account.
func unlinkAccount(ctx context.Context, cfg *config.Config, awsAccountID string) error {
	client := NewNerdGraphClient(cfg)

	acctIDInt := 0
	_, err := fmt.Sscan(cfg.AccountID, &acctIDInt)
	if err != nil {
		return fmt.Errorf("invalid New Relic Account ID: %s", cfg.AccountID)
	}

	// Query existing linked accounts to find the ID
	query := fmt.Sprintf(`
{
  actor {
    account(id: %d) {
      cloud {
        linkedAccounts {
          id
          name
          externalId
        }
      }
    }
  }
}`, acctIDInt)

	var queryResp CloudQueryResponse
	err = client.Query(ctx, query, nil, &queryResp)
	if err != nil {
		return fmt.Errorf("failed to query linked accounts: %w", err)
	}

	var linkedID interface{}
	for _, la := range queryResp.Actor.Account.Cloud.LinkedAccounts {
		if la.ExternalID == awsAccountID {
			linkedID = la.ID
			break
		}
	}

	if linkedID == nil {
		// Already unlinked or never linked
		return nil
	}

	mutation := fmt.Sprintf(`
mutation {
  cloudUnlinkAccount(
    accountId: %d,
    accounts: {
      linkedAccountId: "%v"
    }
  ) {
    unlinkedAccounts {
      id
    }
    errors {
      message
    }
  }
}`, acctIDInt, linkedID)

	var mutationResp struct {
		CloudUnlinkAccount struct {
			UnlinkedAccounts []struct {
				ID interface{} `json:"id"`
			} `json:"unlinkedAccounts"`
			Errors []struct {
				Message string `json:"message"`
			} `json:"errors"`
		} `json:"cloudUnlinkAccount"`
	}

	err = client.Query(ctx, mutation, nil, &mutationResp)
	if err != nil {
		return fmt.Errorf("failed to unlink account: %w", err)
	}

	if len(mutationResp.CloudUnlinkAccount.Errors) > 0 {
		return fmt.Errorf("failed to unlink account: %s", mutationResp.CloudUnlinkAccount.Errors[0].Message)
	}

	return nil
}


// ─────────────────────────────────────────────────────────────
// Helper functions
// ─────────────────────────────────────────────────────────────

func isStackNotFound(err error) bool {
	return strings.Contains(err.Error(), "does not exist") ||
		strings.Contains(err.Error(), "ValidationError")
}

func isActiveStatus(status cftypes.StackStatus) bool {
	return status == cftypes.StackStatusCreateComplete ||
		status == cftypes.StackStatusUpdateComplete
}

func isInProgressStatus(status cftypes.StackStatus) bool {
	return status == cftypes.StackStatusCreateInProgress ||
		status == cftypes.StackStatusUpdateInProgress ||
		status == cftypes.StackStatusDeleteInProgress ||
		status == cftypes.StackStatusUpdateCompleteCleanupInProgress
}

func isFailedStatus(status cftypes.StackStatus) bool {
	return status == cftypes.StackStatusCreateFailed ||
		status == cftypes.StackStatusRollbackFailed ||
		status == cftypes.StackStatusRollbackComplete ||
		status == cftypes.StackStatusDeleteFailed ||
		status == cftypes.StackStatusUpdateRollbackFailed ||
		status == cftypes.StackStatusUpdateRollbackComplete
}
