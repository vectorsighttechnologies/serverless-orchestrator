# Security Policy

## Reporting a Vulnerability

The Serverless Orchestrator team takes security vulnerabilities seriously. We appreciate your efforts to responsibly disclose your findings.

**⚠️ Please do NOT report security vulnerabilities through public GitHub issues.**

Instead, please report them via email to:

📧 **[founders@vectorsight.tech](mailto:founders@vectorsight.tech)**

### What to Include

Please include the following information in your report:

- **Description** of the vulnerability and its potential impact
- **Steps to reproduce** or a proof-of-concept
- **Affected versions** or components (backend, frontend, Lambda)
- **Suggested fix** (if you have one)

### Response Timeline

- **Acknowledgment:** Within 48 hours of receiving your report
- **Initial Assessment:** Within 5 business days
- **Resolution Target:** Security patches within 30 days for critical issues

### Scope

The following areas are in scope for security reports:

- Authentication and authorization bypass
- Credential exposure or leakage (AWS keys, API keys, NR/DD secrets)
- SQL injection or command injection in the backend gateway
- Cross-site scripting (XSS) in the frontend dashboard
- Insecure data transmission or storage
- Lambda function privilege escalation
- CloudFormation template security misconfigurations

### Safe Harbor

We support safe harbor for security researchers who:

- Make a good faith effort to avoid privacy violations, data destruction, and service disruption
- Only interact with accounts you own or with explicit permission
- Do not exploit a vulnerability beyond what is necessary to confirm it
- Report vulnerabilities promptly and do not disclose them publicly until a fix is available

## Security Architecture

Serverless Orchestrator is designed with security in mind:

- **AES-256-GCM encryption** for credentials at rest in the SQLite database
- **JWT-based authentication** for all API endpoints
- **Zero-trust architecture** — sensitive credentials stay within your AWS account
- **API Gateway API Key protection** for Lambda orchestrator endpoints
- **No credential forwarding** — provider keys are configured as environment variables on the Lambda function, not transmitted through the UI

## Supported Versions

| Version | Supported |
|:--------|:---------:|
| Latest  | ✅        |

---

© 2026 [VectorSight Technologies](https://vectorsight.tech). All rights reserved.
