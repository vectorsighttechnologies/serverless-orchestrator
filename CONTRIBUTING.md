# Contributing to Serverless Orchestrator

Thank you for your interest in contributing to **Serverless Orchestrator**! This guide will help you get started.

## 📋 Code of Conduct

By participating in this project, you agree to abide by our [Code of Conduct](CODE_OF_CONDUCT.md).

## ⚖️ Contributor License Agreement (CLA)

All contributors must sign our [Contributor License Agreement (CLA)](CLA.md) before their pull request can be merged. The CLA ensures that:

- Your contribution is **voluntary** with no expectation of compensation
- All **intellectual property rights** are irrevocably assigned to **VectorSight Technologies**
- You retain **no ownership claims** over the contribution once submitted
- Your contribution is your **original work** and doesn't violate third-party rights

### How signing works

1. When you open a PR, a bot will automatically comment asking you to sign
2. You reply with: **"I have read the CLA Document and I hereby sign the CLA"**
3. Your signature is recorded — you only need to sign **once** for all future PRs
4. The PR status check will pass and your contribution can be reviewed

> If you do not agree to the CLA terms, please do not submit contributions to this repository.

## 🚀 Getting Started

### Prerequisites

- **Go 1.22+** — Backend gateway & orchestrator Lambda compilation
- **Node.js 18+ & NPM** — Vue 3 frontend development
- **Python 3.10+** — Lambda packaging scripts
- **AWS CLI & AWS SAM CLI** — Cloud deployment and testing

### Local Development Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/vectorsighttechnologies/serverless-orchestrator.git
   cd serverless-orchestrator
   ```

2. **Start the backend gateway:**
   ```bash
   cd backend
   go mod tidy
   go run ./cmd/server/main.go
   ```

3. **Start the frontend dev server:**
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

4. **Deploy the orchestrator Lambda (optional, for end-to-end testing):**
   ```bash
   cd lambda
   python pack.py
   sam deploy --guided
   ```

## 🔀 How to Contribute

### Reporting Bugs

- Use the [Bug Report](https://github.com/vectorsighttechnologies/serverless-orchestrator/issues/new?template=bug_report.yml) template.
- Include steps to reproduce, expected behavior, and environment details.

### Suggesting Features

- Use the [Feature Request](https://github.com/vectorsighttechnologies/serverless-orchestrator/issues/new?template=feature_request.yml) template.
- Describe the problem, proposed solution, and alternatives considered.

### Submitting Pull Requests

1. **Fork** the repository and create a feature branch from `main`:
   ```bash
   git checkout -b feature/my-awesome-feature
   ```

2. **Make your changes** following the code style guidelines below.

3. **Test your changes:**
   - Backend: `cd backend && go build ./...`
   - Frontend: `cd frontend && npm run build`
   - Lambda: `cd lambda && go build ./cmd/lambda/`

4. **Commit with clear messages** using [Conventional Commits](https://www.conventionalcommits.org/):
   ```
   feat: add bulk export for instrumented functions
   fix: resolve race condition in concurrent Lambda updates
   docs: update deployment guide for new AWS regions
   ```

5. **Push and open a Pull Request** against `main`, filling out the PR template.

## 🏗️ Code Style Guidelines

### Go (Backend & Lambda)

- Follow standard Go formatting (`gofmt`).
- Use meaningful variable and function names.
- Add comments for exported functions and complex logic.
- Handle errors explicitly — no silent swallowing.

### TypeScript / Vue (Frontend)

- Use Vue 3 Composition API with `<script setup lang="ts">`.
- Follow existing component patterns and naming conventions.
- Use CSS custom properties (design tokens) defined in the global styles.
- Keep components focused and reusable.

### General

- Keep pull requests focused — one feature or fix per PR.
- Write descriptive commit messages.
- Update documentation when changing user-facing behavior.

## 📁 Project Structure

```
serverless-orchestrator/
├── backend/          # Go HTTP gateway (SQLite, auth, Lambda invoker)
│   ├── cmd/          # Entry points
│   └── internal/     # Business logic packages
├── frontend/         # Vue 3 + TypeScript SPA (Vite)
│   ├── src/          # Application source
│   └── public/       # Static assets
├── lambda/           # Go Lambda orchestrator (AWS engine)
│   ├── cmd/          # Lambda entry point
│   └── internal/     # AWS SDK integrations
└── .github/          # Issue/PR templates
```

## 📜 License

This project is licensed under the [Apache License 2.0](LICENSE). By contributing, you agree to the [Contributor License Agreement](#️-contributor-license-agreement-cla) above, which assigns all rights in your contribution to **VectorSight Technologies**.

© 2026 [VectorSight Technologies](https://vectorsight.tech). All rights reserved.

---

Thank you for helping make Serverless Orchestrator better! ⚡
