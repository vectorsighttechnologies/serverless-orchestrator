# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Datadog Lambda instrumentation and AWS integration support
- Datadog provider selection in the dashboard UI
- Environment variable configuration for Datadog (`DATADOG_API_KEY`, `DATADOG_SITE`)
- Environment variable configuration for New Relic on the orchestrator Lambda

### Changed
- Renamed project branding from "Lambda Monitor" to "Serverless Orchestrator"
- Updated all community files for open-source readiness (CONTRIBUTING, CODE_OF_CONDUCT, SECURITY)

## [0.1.0] - 2026-07-01

### Added
- Initial release of Serverless Orchestrator
- Vue 3 + TypeScript frontend dashboard with Vite
- Go backend gateway with encrypted SQLite credential store
- Go Lambda orchestrator for AWS Lambda instrumentation
- New Relic layer-based instrumentation (Node.js, Python, Java, .NET, Ruby)
- APM and Serverless telemetry modes
- CloudWatch Metric Streams provisioning via CloudFormation
- Log ingestion via Lambda forwarder
- Multi-region connection management
- JWT-based authentication with AES-256-GCM credential encryption
- Self-cleaning CloudFormation stack management
- SAM template for one-command Lambda deployment

---

© 2026 [VectorSight Technologies](https://vectorsight.tech). Licensed under [Apache 2.0](LICENSE).
