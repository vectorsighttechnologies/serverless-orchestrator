# Serverless Orchestrator

Serverless Orchestrator is a one-stop, self-hosted platform designed to simplify, automate, and manage AWS Lambda instrumentation and observability integrations (such as New Relic) across multiple AWS accounts and regions. 

It provides an intuitive dashboard UI to track function monitoring states, perform concurrent bulk instrumentation (attaching/detaching layers), and natively provision metric stream integrations via CloudFormation.

---

## Key Features

- **Multi-Account & Multi-Region Support:** Easily configure, persist, and switch between multiple AWS connections (Orchestrators) directly from the dashboard header.
- **Concurrent Bulk Instrumentation:** Attach or detach monitoring layers from dozens of Lambda functions simultaneously in seconds.
- **Automated Observability Setup:** Provision AWS-to-observability integrations (e.g. Kinesis Firehose + Cloudwatch Metric Streams) with single-click actions deploying nested CloudFormation templates.
- **NerdGraph Account Linking:** Automatically queries and registers the AWS account connection via NerdGraph GraphQL API, completely unlinking the accounts upon deletion.
- **Self-Cleaning Deployments:** Automatically detects and cleans up failed/stuck CloudFormation stacks before initiating a fresh integration.
- **Internal Exclusions:** Hides platform helper functions and the orchestrator itself to prevent accidental instrumentation.

---

## Repository Structure

- `frontend/` - Single Page Application built using **Vue 3**, **Vite**, **TypeScript**, and modern HSL styling.
- `backend/` - Secure HTTP API Gateway written in **Go**, utilizing **SQLite** (or PostgreSQL) to persist user connections and audit logs.
- `lambda/` - The serverless orchestrator engine in **Go**, packaged and deployed to AWS Lambda using the **AWS SAM CLI**.

---

## Requirements

- **Go 1.22+** (for compiling backend gateway & orchestrator engine)
- **Node.js 18+ & NPM** (for launching the Vue dashboard application)
- **Python 3.10+** (for packaging execution scripts)
- **AWS CLI & AWS SAM CLI** (configured with deployment permissions)

---

## Quick Start Guide

### 1. Deploy the Orchestrator Lambda

First, compile and deploy the orchestrator engine to the AWS region(s) you wish to manage:

```bash
cd lambda

# Package the Go Lambda bootstrap and SAM templates
python pack.py

# Deploy using SAM (interactive options)
sam deploy --guided
```

Upon successful deployment, copy the **OrchestratorApiUrl** and the **OrchestratorApiKey** from the outputs.

### 2. Launch the Backend Gateway

Initialize and run the backend server which handles routing, connection encryption, and audit logs:

```bash
cd ../backend

# Clean up and load dependencies
go mod tidy

# Start the Go server (runs on port 9000 by default)
go run ./cmd/server/main.go
```

The database file `serverless_orchestrator.db` will be initialized automatically in the root of the backend directory.

### 3. Launch the Frontend UI

Install node dependencies and run the dashboard UI locally:

```bash
cd ../frontend

# Install dependencies
npm install

# Start the development server
npm run dev
```

Open the printed URL (typically `http://localhost:5173`) in your browser. Register/login, go to **Settings** / **Connections**, add your newly deployed Orchestrator connection details, and begin managing your serverless functions!

---

## License

Licensed under the **Apache License, Version 2.0** (the "License"). You may obtain a copy of the License in the [LICENSE](./LICENSE) file or at:

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
