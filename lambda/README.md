# Lambda Instrumentation Platform — Orchestrator Deployment

This package contains the AWS Serverless Application Model (SAM) deployment files for the Lambda Orchestrator function.

## Files Included
1. `template.yaml` - The SAM template containing all required AWS resources (Lambda Function, IAM Roles, API Gateway REST API).
2. `function.zip` - The precompiled Linux ARM64 Go binary for the orchestrator function.

## Prerequisites
- Installed and configured **AWS CLI** (with credentials for the target account/region).
- Installed **AWS SAM CLI**.
- AWS permissions to deploy CloudFormation stacks, Lambda functions, IAM roles, and API Gateways.

## Deployment Instructions

### Method 1: Using AWS SAM CLI (Recommended)

1. Extract the downloaded `orchestrator-sam.zip` into a folder.
2. Open a terminal in that folder.
3. Run the following command:
   ```bash
   sam deploy --guided
   ```
4. Follow the interactive prompts:
   - **Stack Name**: `lambda-instrumentation-orchestrator` (or any name you prefer)
   - **AWS Region**: Select the AWS region where your target Lambdas are (e.g., `us-east-1`).
   - **Parameter NewRelicLicenseKey**: (Optional) Enter your New Relic Ingest License Key.
   - **Parameter NewRelicAccountId**: (Optional) Enter your New Relic Account ID.
   - **Parameter NewRelicApiKey**: (Optional) Enter your New Relic User API Key.
   - **Parameter NewRelicRegion**: `us` (default) or `eu`.
   - **Confirm changes before deploy**: `y` or `n`
   - **Allow SAM CLI IAM role creation**: `y` (required to create the execution role)
   - **Disable rollback**: `n`
   - **OrchestratorFunction has no authentication. Is this okay?**: `y` (Note: The API Gateway endpoint *is* secured with an API Key, but SAM prompts this because the Lambda function itself doesn't have an auth provider directly configured).
   - **Save arguments to configuration file**: `y`
   - **SAM configuration file**: Press Enter (`samconfig.toml`)
   - **SAM configuration environment**: Press Enter (`default`)

5. Once the deployment finishes, the terminal will print **Outputs**:
   - **ApiUrl**: The endpoint URL (e.g. `https://xxxxxx.execute-api.us-east-1.amazonaws.com/Prod`). Copy this.
   - **ApiKeyId**: The ID of the API Key. Note that you need the actual key value.

6. Retrieve the API Key value:
   - Run the following AWS CLI command (replace `<ApiKeyId>` with the value from outputs):
     ```bash
     aws apigateway get-api-key --api-key <ApiKeyId> --include-value
     ```
     Look for the `value` field in the JSON response (e.g., `"value": "your-actual-api-key-here"`).
   - Alternatively, open the **AWS Console** -> **API Gateway** -> **API Keys**, click on the key named `<StackName>-api-key`, and click **Show** to copy the value.

7. Enter the **API Gateway URL** and the **API Key** into the configuration screen of the Lambda Monitor UI.
