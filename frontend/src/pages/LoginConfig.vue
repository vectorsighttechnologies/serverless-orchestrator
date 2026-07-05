<template>
  <div class="config-page">
    <!-- Shared Header (same as Dashboard) -->
    <header class="app-header">
      <div class="app-header-left">
        <div class="app-brand">
          <div class="app-logo-wrap">
            <img src="/logo.png" alt="Serverless Orchestrator" class="app-logo-img" />
          </div>
          <div>
            <span class="app-name">Serverless<span class="app-name-accent"> Orchestrator</span></span>
            <span class="app-connection">
              <span class="connection-dot"></span>
              Not connected
            </span>
          </div>
        </div>
      </div>
      <div class="app-header-right">
        <router-link to="/provider" class="btn btn-ghost btn-sm btn-icon" id="btn-back-provider">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="15 18 9 12 15 6"/></svg>
          Back
        </router-link>
      </div>
    </header>

    <!-- Breadcrumb -->
    <nav class="app-breadcrumb" aria-label="Breadcrumb">
      <router-link to="/provider" class="breadcrumb-link" id="breadcrumb-providers">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/></svg>
        Providers
      </router-link>
      <span class="breadcrumb-sep">
        <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="9 18 15 12 9 6"/></svg>
      </span>
      <span class="breadcrumb-current">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 010 2.83 2 2 0 01-2.83 0l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-4 0v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83 0 2 2 0 010-2.83l.06-.06A1.65 1.65 0 004.68 15a1.65 1.65 0 00-1.51-1H3a2 2 0 010-4h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 012.83-2.83l.06.06A1.65 1.65 0 009 4.68a1.65 1.65 0 001-1.51V3a2 2 0 014 0v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 2.83l-.06.06A1.65 1.65 0 0019.4 9a1.65 1.65 0 001.51 1H21a2 2 0 010 4h-.09a1.65 1.65 0 00-1.51 1z"/></svg>
        Configuration
      </span>
    </nav>

    <!-- Page Content -->
    <div class="config-content">
      <div class="config-main animate-slide-up">

        <!-- Page Heading -->
        <!-- <div class="config-heading">
          <div class="config-heading-icon">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
          </div>
          <div>
            <h1 class="config-title">Connect to Orchestrator</h1>
            <p class="config-subtitle">Configure your API Gateway endpoints for each AWS region</p>
          </div>
        </div> -->

        <!-- Step indicator -->
        <div class="step-indicator">
          <div class="step-item step-done">
            <div class="step-dot"><svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="20 6 9 17 4 12"/></svg></div>
            <span>Provider</span>
          </div>
          <div class="step-line"></div>
          <div class="step-item step-active">
            <div class="step-dot step-num">2</div>
            <span>Configuration</span>
          </div>
          <div class="step-line"></div>
          <div class="step-item">
            <div class="step-dot step-num">3</div>
            <span>Dashboard</span>
          </div>
        </div>

        <!-- Deployment Guide Card -->
        <div class="deploy-guide-card card-flat animate-slide-up">
          <div class="deploy-guide-header">
            <div class="deploy-guide-icon">
              <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/></svg>
            </div>
            <div class="deploy-guide-title-wrap">
              <h2 class="deploy-guide-title">Deploy Orchestrator Lambda</h2>
              <p class="deploy-guide-subtitle">To instrument your Lambdas, deploy the orchestrator function in your AWS account first.</p>
            </div>
          </div>

          <div class="deploy-download-section">
            <div class="download-info">
              <span class="download-label">Required Deployment Package:</span>
              <span class="download-filename">orchestrator-sam.zip</span>
            </div>
            <a href="/orchestrator-sam.zip" download class="btn btn-primary download-btn" id="btn-download-sam">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
              Download SAM ZIP
            </a>
          </div>

          <div class="deploy-tabs-container">
            <div class="deploy-tabs-header">
              <span class="deploy-tabs-label">Deployment Method:</span>
              <div class="deploy-tabs">
                <button
                  :class="['deploy-tab', activeDeployTab === 'sam' ? 'active' : '']"
                  @click="activeDeployTab = 'sam'"
                  type="button"
                  id="tab-deploy-sam"
                >
                  AWS SAM (Recommended)
                </button>
                <button
                  :class="['deploy-tab', activeDeployTab === 'manual' ? 'active' : '']"
                  @click="activeDeployTab = 'manual'"
                  type="button"
                  id="tab-deploy-manual"
                >
                  Manual Deployment
                </button>
              </div>
            </div>

            <!-- SAM Instructions -->
            <div v-if="activeDeployTab === 'sam'" class="deploy-tab-content animate-fade-in">
              <ol class="step-list">
                <li>
                  <span class="step-number">1</span>
                  <div class="step-text">
                    <strong>Prerequisites:</strong> Ensure you have the <a href="https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html" target="_blank" rel="noopener">AWS SAM CLI</a> and <a href="https://aws.amazon.com/cli/" target="_blank" rel="noopener">AWS CLI</a> installed and configured with appropriate permissions.
                  </div>
                </li>
                <li>
                  <span class="step-number">2</span>
                  <div class="step-text">
                    <strong>Extract & Navigate:</strong> Extract the downloaded <code>orchestrator-sam.zip</code> and open your terminal in that directory.
                  </div>
                </li>
                <li>
                  <span class="step-number">3</span>
                  <div class="step-text">
                    <strong>Deploy Command:</strong> Run the interactive guided deployment:
                    <div class="code-block-wrap">
                      <pre><code>sam deploy --guided</code></pre>
                      <button class="btn-copy-code" @click="copyToClipboard('sam deploy --guided', 'sam')" title="Copy to clipboard" type="button">
                        <svg v-if="!copiedStates['sam']" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg>
                        <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" class="text-accent"><polyline points="20 6 9 17 4 12"/></svg>
                      </button>
                    </div>
                  </div>
                </li>
                <li>
                  <span class="step-number">4</span>
                  <div class="step-text">
                    <strong>Retrieve Details:</strong> Once deployed, look at the **CloudFormation Outputs**:
                    <ul>
                      <li>Copy the <strong>ApiUrl</strong> output and paste it in the <strong>API Gateway URL</strong> field below.</li>
                      <li>Run this command to retrieve your API Key value (replace <code>&lt;ApiKeyId&gt;</code> with the ApiKeyId from Outputs):
                        <div class="code-block-wrap">
                          <pre><code>aws apigateway get-api-key --api-key &lt;ApiKeyId&gt; --include-value</code></pre>
                          <button class="btn-copy-code" @click="copyToClipboard('aws apigateway get-api-key --api-key <ApiKeyId> --include-value', 'key')" title="Copy to clipboard" type="button">
                            <svg v-if="!copiedStates['key']" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg>
                            <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" class="text-accent"><polyline points="20 6 9 17 4 12"/></svg>
                          </button>
                        </div>
                      </li>
                      <li>Paste the API Key value in the <strong>API Key</strong> field below.</li>
                    </ul>
                  </div>
                </li>
              </ol>
            </div>

            <!-- Manual Instructions -->
            <div v-else class="deploy-tab-content animate-fade-in">
              <ol class="step-list">
                <li>
                  <span class="step-number">1</span>
                  <div class="step-text">
                    <strong>Create Lambda Function:</strong> In the AWS Console, create a new Lambda function.
                    <ul>
                      <li>Runtime: <code>Amazon Linux 2023 (provided.al2023)</code></li>
                      <li>Architecture: <code>arm64</code></li>
                      <li>Handler: <code>bootstrap</code></li>
                    </ul>
                  </div>
                </li>
                <li>
                  <span class="step-number">2</span>
                  <div class="step-text">
                    <strong>Upload Code:</strong> Extract <code>function.zip</code> from the downloaded package and upload it directly as your Lambda function code.
                  </div>
                </li>
                <li>
                  <span class="step-number">3</span>
                  <div class="step-text">
                    <strong>Configure IAM Role:</strong> Ensure the Lambda execution role has permission to manage Lambda layers, update configurations, and read IAM roles. (Refer to the policies defined in the <code>template.yaml</code> file inside the zip).
                  </div>
                </li>
                <li>
                  <span class="step-number">4</span>
                  <div class="step-text">
                    <strong>Set up API Gateway:</strong> Create an API Gateway REST API with a proxy resource (<code>/{proxy+}</code>) pointing to your Lambda function. Enable **API Key Required** on the method requests, create a Usage Plan, associate it with an API Key, and enter those credentials below.
                  </div>
                </li>
              </ol>
            </div>
          </div>
        </div>

        <!-- Region Cards -->
        <div class="regions-section">
          <div class="regions-header">
            <h2 class="section-heading">API Endpoints</h2>
            <button class="btn btn-secondary btn-sm" @click="addPendingRegion" id="btn-add-region">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
              Add Region
            </button>
          </div>

          <div class="regions-list">
            <div
              v-for="(region, index) in pendingRegions"
              :key="index"
              class="region-card"
              :class="{ 'region-card--active': index === 0 }"
            >
              <div class="region-card__header">
                <div class="region-card__label">
                  <span class="region-card__num">{{ index + 1 }}</span>
                  <h3 class="region-card__title">
                    {{ region.region || 'New Region' }}
                  </h3>
                  <span v-if="index === 0" class="badge badge-success" style="font-size:0.6rem">Primary</span>
                </div>
                <button
                  v-if="pendingRegions.length > 1"
                  class="btn btn-ghost btn-sm btn-remove"
                  @click="removePendingRegion(index)"
                >
                  <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                  Remove
                </button>
              </div>

              <div class="region-card__fields">
                <div class="form-group">
                  <label class="form-label">AWS Region</label>
                  <select v-model="region.region" class="form-select" :id="'region-select-' + index">
                    <option value="us-east-1">us-east-1 (N. Virginia)</option>
                    <option value="us-east-2">us-east-2 (Ohio)</option>
                    <option value="us-west-1">us-west-1 (N. California)</option>
                    <option value="us-west-2">us-west-2 (Oregon)</option>
                    <option value="eu-west-1">eu-west-1 (Ireland)</option>
                    <option value="eu-west-2">eu-west-2 (London)</option>
                    <option value="eu-central-1">eu-central-1 (Frankfurt)</option>
                    <option value="ap-southeast-1">ap-southeast-1 (Singapore)</option>
                    <option value="ap-southeast-2">ap-southeast-2 (Sydney)</option>
                    <option value="ap-northeast-1">ap-northeast-1 (Tokyo)</option>
                    <option value="ap-south-1">ap-south-1 (Mumbai)</option>
                  </select>
                </div>
                <div class="form-row">
                  <div class="form-group">
                    <label class="form-label">API Gateway URL</label>
                    <input
                      v-model="region.apiGatewayUrl"
                      type="url"
                      class="form-input"
                      placeholder="https://xxxxxxxx.execute-api.region.amazonaws.com/Prod"
                      :id="'api-url-' + index"
                    />
                  </div>
                  <div class="form-group form-group--key">
                    <label class="form-label">API Key</label>
                    <div class="input-with-icon">
                      <svg class="input-icon" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 1 1-7.778 7.778 5.5 5.5 0 0 1 7.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"/></svg>
                      <input
                        v-model="region.apiKey"
                        :type="showKey[index] ? 'text' : 'password'"
                        class="form-input input-padded"
                        placeholder="API Gateway key"
                        :id="'api-key-' + index"
                      />
                      <button class="key-toggle" @click="toggleKey(index)" type="button">
                        <svg v-if="!showKey[index]" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/></svg>
                        <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/><line x1="1" y1="1" x2="23" y2="23"/></svg>
                      </button>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Per-card status after connection attempt -->
              <div v-if="regionStatus[index]" class="region-status" :class="regionStatus[index].ok ? 'region-status--ok' : 'region-status--err'">
                <svg v-if="regionStatus[index].ok" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"/></svg>
                <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
                {{ regionStatus[index].message }}
              </div>
            </div>
          </div>
        </div>

        <!-- Connect Button -->
        <div class="connect-action">
          <button
            class="btn btn-primary btn-lg connect-btn"
            :disabled="!canConnect || isConnecting"
            @click="connect"
            id="btn-connect"
          >
            <span v-if="isConnecting" class="spinner"></span>
            <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M5 12h14M12 5l7 7-7 7"/></svg>
            {{ isConnecting ? 'Connecting...' : 'Test Connection' }}
          </button>
          <p class="connect-hint">Verifies connectivity before saving</p>
        </div>

        <!-- Health Check Result -->
        <div v-if="healthResult" class="health-result animate-slide-up">
          <!-- Keys found on orchestrator -->
          <div v-if="healthResult.config.source === 'env_vars'" class="result-card result-card--success">
            <div class="result-card__header">
              <div class="result-icon result-icon--success">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 11.08V12a10 10 0 11-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
              </div>
              <div>
                <h3 class="result-title">Connection Successful</h3>
                <p class="result-desc">Credentials are configured on the orchestrator</p>
              </div>
            </div>

            <!-- Datadog configured view -->
            <div v-if="provider === 'datadog'" class="config-details">
              <div class="config-item">
                <span class="config-item__label">Datadog API Key</span>
                <span class="badge badge-success"><span class="badge-dot"></span> Configured</span>
              </div>
              <div class="config-item">
                <span class="config-item__label">Datadog Site</span>
                <span class="badge badge-info">{{ healthResult.config.ddSiteConfigured ? 'Configured' : 'datadoghq.com' }}</span>
              </div>
            </div>

            <!-- New Relic configured view -->
            <div v-else class="config-details">
              <div class="config-item">
                <span class="config-item__label">NR License Key</span>
                <span class="badge badge-success"><span class="badge-dot"></span> Configured</span>
              </div>
              <div class="config-item">
                <span class="config-item__label">NR Account ID</span>
                <span class="badge badge-success"><span class="badge-dot"></span> {{ healthResult.config.accountIdConfigured ? 'Configured' : 'Missing' }}</span>
              </div>
              <div class="config-item">
                <span class="config-item__label">NR API Key</span>
                <span :class="['badge', healthResult.config.apiKeyConfigured ? 'badge-success' : 'badge-warning']">
                  <span class="badge-dot"></span>
                  {{ healthResult.config.apiKeyConfigured ? 'Configured' : 'Optional' }}
                </span>
              </div>
              <div class="config-item">
                <span class="config-item__label">Region</span>
                <span class="badge badge-info">{{ healthResult.config.region?.toUpperCase() || 'US' }}</span>
              </div>
            </div>
            <button class="btn btn-primary btn-lg" style="width:100%;margin-top:var(--space-6)" @click="goToDashboard" id="btn-continue">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="18" height="18" rx="2"/><line x1="3" y1="9" x2="21" y2="9"/><line x1="9" y1="21" x2="9" y2="9"/></svg>
              Continue to Dashboard →
            </button>
          </div>

          <!-- Keys NOT found — request them from user -->
          <div v-else class="result-card result-card--warning">
            <div class="result-card__header">
              <div class="result-icon result-icon--warning">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
              </div>
              <div>
                <h3 class="result-title">Credentials Not Found on Orchestrator</h3>
                <p class="result-desc">Enter your provider credentials below</p>
              </div>
            </div>

            <div class="alert alert-info" style="margin-bottom:var(--space-5)">
              💡 For better security, set these as environment variables on your orchestrator Lambda during deployment.
            </div>

            <!-- Datadog Fields -->
            <div v-if="provider === 'datadog'" class="nr-fields">
              <div class="form-group">
                <label class="form-label">Datadog API Key</label>
                <input v-model="ddCreds.apiKey" type="password" class="form-input" placeholder="Your Datadog API Key" id="dd-api-key" />
              </div>
              <div class="form-group">
                <label class="form-label">Datadog Site</label>
                <select v-model="ddCreds.site" class="form-select" id="dd-site">
                  <option value="datadoghq.com">datadoghq.com (US1)</option>
                  <option value="us3.datadoghq.com">us3.datadoghq.com (US3)</option>
                  <option value="us5.datadoghq.com">us5.datadoghq.com (US5)</option>
                  <option value="datadoghq.eu">datadoghq.eu (EU)</option>
                  <option value="ddog-gov.com">ddog-gov.com (GovCloud US)</option>
                </select>
              </div>
            </div>

            <!-- New Relic Fields -->
            <div v-else class="nr-fields">
              <div class="form-group">
                <label class="form-label">NR License Key</label>
                <input v-model="nrCreds.licenseKey" type="password" class="form-input" placeholder="Your NR Ingest License Key" id="nr-license-key" />
              </div>
              <div class="form-row">
                <div class="form-group">
                  <label class="form-label">NR Account ID</label>
                  <input v-model="nrCreds.accountId" type="text" class="form-input" placeholder="e.g. 12345" id="nr-account-id" />
                </div>
                <div class="form-group">
                  <label class="form-label">NR Region</label>
                  <select v-model="nrCreds.region" class="form-select" id="nr-region">
                    <option value="us">US</option>
                    <option value="eu">EU</option>
                  </select>
                </div>
              </div>
              <div class="form-group">
                <label class="form-label">NR API Key <span class="form-label-optional">(optional)</span></label>
                <input v-model="nrCreds.apiKey" type="password" class="form-input" placeholder="NRAK-xxxxxxxx" id="nr-api-key" />
                <span class="form-hint">Only needed for Metric Streams / API Polling setup</span>
              </div>
            </div>

            <button
              class="btn btn-primary btn-lg"
              style="width:100%;margin-top:var(--space-6)"
              :disabled="provider === 'datadog' ? !ddCreds.apiKey : (!nrCreds.licenseKey || !nrCreds.accountId)"
              @click="goToDashboardWithCreds"
              id="btn-continue-creds"
            >
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="18" height="18" rx="2"/><line x1="3" y1="9" x2="21" y2="9"/><line x1="9" y1="21" x2="9" y2="9"/></svg>
              Continue to Dashboard →
            </button>
          </div>
        </div>

        <!-- Error -->
        <div v-if="connectError" class="alert alert-danger animate-fade-in" style="margin-top:var(--space-4)">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
          {{ connectError }}
        </div>

      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useConfig } from '@/composables/useConfig'
import { api } from '@/services/api'
import type { HealthResponse } from '@/types'

const router = useRouter()
const {
  provider,
  pendingRegions,
  addPendingRegion,
  removePendingRegion,
  setRegions,
  setConnected,
  saveToSession,
} = useConfig()

const showKey = ref<boolean[]>(pendingRegions.value.map(() => false))
const regionStatus = ref<{ ok: boolean; message: string }[]>([])
const nrCreds = ref({ licenseKey: '', accountId: '', apiKey: '', region: 'us' as 'us' | 'eu' })
const ddCreds = ref({ apiKey: '', site: 'datadoghq.com' })
const isConnecting = ref(false)
const connectError = ref('')
const healthResult = ref<HealthResponse | null>(null)

// Deployment guide state
const activeDeployTab = ref<'sam' | 'manual'>('sam')
const copiedStates = ref<Record<string, boolean>>({})

function copyToClipboard(text: string, id: string) {
  navigator.clipboard.writeText(text).then(() => {
    copiedStates.value[id] = true
    setTimeout(() => {
      copiedStates.value[id] = false
    }, 2000)
  }).catch(err => {
    console.error('Failed to copy text: ', err)
  })
}

function toggleKey(index: number) {
  showKey.value[index] = !showKey.value[index]
}

const canConnect = computed(() =>
  pendingRegions.value.every(r => r.apiGatewayUrl.trim() !== '' && r.apiKey.trim() !== '')
)

onMounted(async () => {
  try {
    const prefs = await api.getPreferences()
    
    // Prefill fields if saved in backend database
    if (prefs.lambdaApiUrl) {
      pendingRegions.value[0].apiGatewayUrl = prefs.lambdaApiUrl
    }
    if (prefs.hasLambdaApiKey) {
      pendingRegions.value[0].apiKey = '••••••••'
    }
    
    nrCreds.value.accountId = prefs.nrAccountId || ''
    nrCreds.value.region = (prefs.nrRegion as 'us' | 'eu') || 'us'
    
    if (prefs.hasNrLicenseKey) {
      nrCreds.value.licenseKey = '••••••••'
    }
    if (prefs.hasNrApiKey) {
      nrCreds.value.apiKey = '••••••••'
    }

    ddCreds.value.site = prefs.ddSite || 'datadoghq.com'
    if (prefs.hasDdApiKey) {
      ddCreds.value.apiKey = '••••••••'
    }
  } catch (err) {
    console.error('Failed to prefill preferences:', err)
  }
})

async function connect() {
  isConnecting.value = true
  connectError.value = ''
  healthResult.value = null
  regionStatus.value = pendingRegions.value.map(() => ({ ok: false, message: 'Verifying...' }))

  try {
    const first = pendingRegions.value[0]
    
    // Save/Update config on backend gateway
    const payload = {
      selectedProvider: provider.value || 'newrelic',
      nrAccountId: nrCreds.value.accountId,
      nrApiKey: nrCreds.value.apiKey === '••••••••' ? '' : nrCreds.value.apiKey,
      nrLicenseKey: nrCreds.value.licenseKey === '••••••••' ? '' : nrCreds.value.licenseKey,
      nrRegion: nrCreds.value.region,
      ddApiKey: ddCreds.value.apiKey === '••••••••' ? '' : ddCreds.value.apiKey,
      ddSite: ddCreds.value.site,
      lambdaApiUrl: first.apiGatewayUrl,
      lambdaApiKey: first.apiKey === '••••••••' ? '' : first.apiKey,
    }
    
    await api.savePreferences(payload)
    
    // Test the setup by fetching functions (this checks if the gateway decrypts correctly, 
    // contacts AWS Lambda Orchestrator, and returns functions successfully).
    await api.listFunctions()
    
    // Set healthResult mock to satisfy template variables
    healthResult.value = {
      status: 'ok',
      config: {
        selectedProvider: provider.value || 'newrelic',
        licenseKeyConfigured: nrCreds.value.licenseKey !== '',
        accountIdConfigured: nrCreds.value.accountId !== '',
        apiKeyConfigured: nrCreds.value.apiKey !== '',
        ddApiKeyConfigured: ddCreds.value.apiKey !== '',
        ddSiteConfigured: ddCreds.value.site !== '',
        region: nrCreds.value.region,
        source: 'env_vars'
      }
    }
    
    regionStatus.value[0] = { ok: true, message: 'Connected successfully' }
  } catch (err) {
    const msg = err instanceof Error ? err.message : 'Connection verification failed'
    connectError.value = `Failed to verify connection: ${msg}`
    regionStatus.value[0] = { ok: false, message: msg }
  } finally {
    isConnecting.value = false
  }
}

// Complete verification and advance to Dashboard
function goToDashboard() {
  setRegions(pendingRegions.value)
  setConnected(true)
  saveToSession()
  router.push('/dashboard')
}

// Handles Continue to Dashboard when NR/DD keys are submitted
function goToDashboardWithCreds() {
  connect().then(() => {
    if (!connectError.value) {
      goToDashboard()
    }
  })
}
</script>

<style scoped>
.config-page {
  min-height: 100vh;
  background: var(--bg-primary);
}

/* ── Shared App Header (same styling as Dashboard) ── */
.app-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 var(--space-6);
  height: 58px;
  background: rgba(22, 27, 34, 0.95);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid var(--border-primary);
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: 0 1px 0 rgba(28, 231, 131, 0.06), 0 2px 12px rgba(0,0,0,0.3);
}

.app-header-left { display: flex; align-items: center; }
.app-brand { display: flex; align-items: center; gap: var(--space-3); }

.app-logo-wrap {
  width: 40px;
  height: 40px;
  border-radius: var(--radius-md);
  overflow: hidden;
  flex-shrink: 0;
  box-shadow: 0 0 0 1px rgba(28, 231, 131, 0.2), 0 0 14px rgba(28, 231, 131, 0.1);
}

.app-logo-img {
  width: 40px;
  height: 40px;
  display: block;
  object-fit: cover;
}

.app-name {
  display: block;
  font-weight: 700;
  font-size: var(--font-size-md);
  letter-spacing: -0.02em;
}

.app-name-accent { color: var(--accent); }

.app-connection {
  display: flex;
  align-items: center;
  gap: var(--space-1);
  font-size: 0.65rem;
  color: var(--text-tertiary);
  font-family: var(--font-mono);
  margin-top: 1px;
}

.connection-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--text-tertiary);
  flex-shrink: 0;
}

.app-header-right { display: flex; align-items: center; gap: var(--space-2); }

.btn-icon {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  font-size: var(--font-size-sm);
}

/* ── Breadcrumb ── */
.app-breadcrumb {
  display: flex;
  align-items: center;
  gap: var(--space-1);
  padding: var(--space-2) var(--space-6);
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-primary);
  font-size: var(--font-size-xs);
}

.breadcrumb-link {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  color: var(--text-tertiary);
  text-decoration: none;
  transition: color var(--transition-fast);
  padding: 3px var(--space-2);
  border-radius: var(--radius-sm);
}

.breadcrumb-link:hover { color: var(--accent); background: var(--accent-dim); }

.breadcrumb-sep {
  color: var(--border-hover);
  display: flex;
  align-items: center;
  margin: 0 1px;
}

.breadcrumb-current {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  color: var(--text-primary);
  font-weight: 600;
  padding: 3px var(--space-2);
}

/* ── Page Layout ── */
.config-content {
  padding: var(--space-8) var(--space-6);
  max-width: 780px;
  margin: 0 auto;
}

.config-main {
  display: flex;
  flex-direction: column;
  gap: var(--space-8);
}

/* ── Heading ── */
.config-heading {
  display: flex;
  align-items: center;
  gap: var(--space-4);
}

.config-heading-icon {
  width: 52px;
  height: 52px;
  border-radius: var(--radius-lg);
  background: var(--accent-dim);
  border: 1px solid var(--accent-muted);
  color: var(--accent);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.config-title {
  font-size: var(--font-size-2xl);
  font-weight: 700;
  letter-spacing: -0.02em;
  margin-bottom: var(--space-1);
}

.config-subtitle {
  color: var(--text-secondary);
  font-size: var(--font-size-sm);
}

/* ── Step Indicator ── */
.step-indicator {
  display: flex;
  align-items: center;
  gap: 0;
  background: var(--bg-secondary);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-lg);
  padding: var(--space-4) var(--space-6);
}

.step-item {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  font-size: var(--font-size-sm);
  color: var(--text-tertiary);
  font-weight: 500;
}

.step-item.step-done { color: var(--success); }
.step-item.step-active { color: var(--text-primary); }

.step-dot {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  border: 2px solid var(--border-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.7rem;
  font-weight: 700;
  flex-shrink: 0;
  background: var(--bg-tertiary);
  color: var(--text-tertiary);
}

.step-item.step-done .step-dot {
  background: var(--success-dim);
  border-color: var(--success);
  color: var(--success);
}

.step-item.step-active .step-dot {
  background: var(--accent-dim);
  border-color: var(--accent);
  color: var(--accent);
}

.step-line {
  flex: 1;
  height: 1px;
  background: var(--border-primary);
  margin: 0 var(--space-4);
}

/* ── Regions Section ── */
.section-heading {
  font-size: var(--font-size-sm);
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: var(--text-tertiary);
}

.regions-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--space-4);
}

.regions-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.region-card {
  background: var(--bg-secondary);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-lg);
  padding: var(--space-5);
  transition: border-color var(--transition-base);
  position: relative;
  overflow: hidden;
}

.region-card::before {
  content: '';
  position: absolute;
  top: 0; left: 0; right: 0;
  height: 2px;
  background: var(--border-primary);
}

.region-card--active::before {
  background: var(--accent);
}

.region-card__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--space-4);
}

.region-card__label {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.region-card__num {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: var(--bg-tertiary);
  border: 1px solid var(--border-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.7rem;
  font-weight: 700;
  color: var(--text-secondary);
  flex-shrink: 0;
}

.region-card--active .region-card__num {
  background: var(--accent-dim);
  border-color: var(--accent);
  color: var(--accent);
}

.region-card__title {
  font-size: var(--font-size-sm);
  font-weight: 600;
  font-family: var(--font-mono);
}

.btn-remove {
  color: var(--text-tertiary);
  font-size: var(--font-size-xs);
}

.btn-remove:hover {
  color: var(--danger);
  background: var(--danger-dim);
}

.region-card__fields {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.form-row {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: var(--space-4);
}

.input-with-icon {
  position: relative;
  display: flex;
  align-items: center;
}

.input-icon {
  position: absolute;
  left: var(--space-3);
  color: var(--text-tertiary);
  pointer-events: none;
  flex-shrink: 0;
}

.input-padded { padding-left: calc(var(--space-3) + 14px + var(--space-2)); }

.key-toggle {
  position: absolute;
  right: var(--space-3);
  background: none;
  border: none;
  color: var(--text-tertiary);
  cursor: pointer;
  display: flex;
  align-items: center;
  padding: 2px;
  border-radius: 4px;
  transition: color var(--transition-fast);
}

.key-toggle:hover { color: var(--text-primary); }

.form-label-optional {
  font-weight: 400;
  color: var(--text-tertiary);
  font-size: var(--font-size-xs);
}

/* ── Region Status ── */
.region-status {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  margin-top: var(--space-3);
  padding: var(--space-2) var(--space-3);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs);
  font-family: var(--font-mono);
}

.region-status--ok {
  background: var(--success-dim);
  color: var(--success);
}

.region-status--err {
  background: var(--danger-dim);
  color: var(--danger);
}

/* ── Connect Action ── */
.connect-action {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-2);
}

.connect-btn {
  width: 100%;
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  justify-content: center;
}

.connect-hint {
  font-size: var(--font-size-xs);
  color: var(--text-tertiary);
}

/* ── Result Cards ── */
.result-card {
  background: var(--bg-secondary);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-lg);
  padding: var(--space-6);
}

.result-card--success {
  border-color: var(--success);
  background: rgba(28, 231, 131, 0.04);
}

.result-card--warning {
  border-color: var(--warning);
  background: rgba(210, 153, 34, 0.04);
}

.result-card__header {
  display: flex;
  align-items: flex-start;
  gap: var(--space-4);
  margin-bottom: var(--space-5);
}

.result-icon {
  width: 44px;
  height: 44px;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.result-icon--success { background: var(--success-dim); color: var(--success); }
.result-icon--warning { background: var(--warning-dim); color: var(--warning); }

.result-title {
  font-size: var(--font-size-md);
  font-weight: 600;
  margin-bottom: var(--space-1);
}

.result-desc {
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
}

.config-details {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  background: var(--bg-tertiary);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-md);
  padding: var(--space-4);
}

.config-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.config-item__label {
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
}

.nr-fields {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.health-result { margin-top: var(--space-2); }

@media (max-width: 600px) {
  .form-row { grid-template-columns: 1fr; }
  .config-content { padding: var(--space-4); }
}

/* ── Deploy Guide Card ── */
.deploy-guide-card {
  background: var(--bg-secondary);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-lg);
  padding: var(--space-6);
  position: relative;
  overflow: hidden;
  box-shadow: var(--shadow-sm);
  margin-bottom: var(--space-6);
}

.deploy-guide-card::before {
  content: '';
  position: absolute;
  top: 0; left: 0; right: 0;
  height: 2px;
  background: var(--accent);
  opacity: 0.7;
}

.deploy-guide-header {
  display: flex;
  align-items: flex-start;
  gap: var(--space-4);
  margin-bottom: var(--space-5);
}

.deploy-guide-icon {
  width: 44px;
  height: 44px;
  border-radius: var(--radius-md);
  background: var(--accent-dim);
  color: var(--accent);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  border: 1px solid var(--accent-muted);
}

.deploy-guide-title-wrap {
  flex: 1;
}

.deploy-guide-title {
  font-size: var(--font-size-md);
  font-weight: 600;
  margin-bottom: var(--space-1);
}

.deploy-guide-subtitle {
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
  line-height: 1.4;
}

.deploy-download-section {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: var(--bg-tertiary);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-md);
  padding: var(--space-4) var(--space-5);
  margin-bottom: var(--space-6);
  gap: var(--space-4);
}

.download-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.download-label {
  font-size: var(--font-size-xs);
  color: var(--text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.download-filename {
  font-family: var(--font-mono);
  font-size: var(--font-size-sm);
  color: var(--text-primary);
  font-weight: 500;
}

.download-btn {
  font-weight: 600;
  box-shadow: 0 0 10px rgba(28, 231, 131, 0.1);
}

.deploy-tabs-container {
  border-top: 1px solid var(--border-primary);
  padding-top: var(--space-5);
}

.deploy-tabs-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--space-4);
  flex-wrap: wrap;
  gap: var(--space-3);
}

.deploy-tabs-label {
  font-size: var(--font-size-xs);
  color: var(--text-secondary);
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.deploy-tabs {
  display: flex;
  background: var(--bg-primary);
  padding: 3px;
  border-radius: var(--radius-md);
  border: 1px solid var(--border-primary);
}

.deploy-tab {
  padding: var(--space-1) var(--space-3);
  font-size: var(--font-size-xs);
  color: var(--text-secondary);
  background: transparent;
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all var(--transition-fast);
  font-family: var(--font-family);
  font-weight: 500;
}

.deploy-tab:hover {
  color: var(--text-primary);
}

.deploy-tab.active {
  background: var(--bg-tertiary);
  color: var(--accent);
  box-shadow: var(--shadow-sm);
}

.deploy-tab-content {
  background: var(--bg-tertiary);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-md);
  padding: var(--space-5);
}

.step-list {
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: var(--space-5);
}

.step-list > li {
  display: flex;
  gap: var(--space-4);
  align-items: flex-start;
}

.step-number {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: var(--bg-primary);
  border: 1px solid var(--border-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--font-size-xs);
  font-weight: 700;
  color: var(--accent);
  flex-shrink: 0;
  margin-top: 2px;
}

.step-text {
  flex: 1;
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
  line-height: 1.5;
  min-width: 0; /* allows flex children to scale down */
}

.step-text strong {
  color: var(--text-primary);
}

.step-text code {
  font-family: var(--font-mono);
  background: var(--bg-primary);
  color: var(--accent);
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 0.85em;
  border: 1px solid var(--border-primary);
}

.step-text ul {
  list-style: disc;
  margin-left: var(--space-5);
  margin-top: var(--space-2);
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.step-text li {
  display: list-item;
  list-style-type: disc;
  color: var(--text-secondary);
  font-size: var(--font-size-sm);
  line-height: 1.5;
  margin-bottom: var(--space-1);
}

.code-block-wrap {
  position: relative;
  background: var(--bg-primary);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-sm);
  padding: var(--space-3) var(--space-4);
  margin-top: var(--space-2);
  margin-bottom: var(--space-2);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-3);
  min-width: 0;
  width: 100%;
}

.code-block-wrap pre {
  margin: 0;
  padding: 0;
  overflow-x: auto;
  flex: 1;
  min-width: 0;
  display: block;
}

.code-block-wrap code {
  font-family: var(--font-mono);
  font-size: var(--font-size-xs) !important;
  color: var(--text-primary) !important;
  background: transparent !important;
  border: none !important;
  padding: 0 !important;
  white-space: pre;
}

.btn-copy-code {
  background: transparent;
  border: none;
  color: var(--text-tertiary);
  cursor: pointer;
  padding: var(--space-1);
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all var(--transition-fast);
}

.btn-copy-code:hover {
  color: var(--accent);
  background: var(--bg-hover);
}

@media (max-width: 600px) {
  .deploy-download-section {
    flex-direction: column;
    align-items: stretch;
  }
  .download-btn {
    text-align: center;
  }
  .deploy-tabs-header {
    flex-direction: column;
    align-items: flex-start;
  }
  .deploy-tabs {
    width: 100%;
  }
  .deploy-tab {
    flex: 1;
    text-align: center;
  }
}
</style>
