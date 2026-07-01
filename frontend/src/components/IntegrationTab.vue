<template>
  <div class="integration-tab animate-fade-in">
    <!-- Status Card -->
    <div class="integration-status-card card-flat">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-3">
          <h2 class="integration-title">AWS Integration</h2>
          <span :class="['badge', statusBadgeClass]">
            <span class="badge-dot"></span>
            {{ statusLabel }}
          </span>
        </div>
        <button
          v-if="integration.status === 'active'"
          class="btn btn-danger btn-sm"
          @click="handleRemove"
          :disabled="isRemoving"
          id="btn-remove-integration"
        >
          {{ isRemoving ? 'Removing...' : 'Remove Integration' }}
        </button>
      </div>

      <!-- Active Details -->
      <div v-if="integration.status === 'active'" class="status-details mt-6">
        <div class="detail-row">
          <span class="detail-label">Method</span>
          <span class="detail-value">{{ integration.method === 'metric_streams' ? 'Metric Streams (PUSH)' : 'API Polling (PULL)' }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">CF Stack</span>
          <span class="detail-value font-mono text-xs">{{ integration.stackName }} — {{ integration.stackStatus }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Last Checked</span>
          <span class="detail-value text-xs">{{ formatTime(integration.lastChecked) }}</span>
        </div>
      </div>
    </div>

    <!-- Setup Form (only if not setup) -->
    <div v-if="integration.status === 'not_setup' || integration.status === 'error'" class="setup-form mt-6">
      
      <!-- Datadog AWS Integration Form -->
      <div v-if="provider === 'datadog'" class="card-flat">
        <h3 class="setup-title">Setup Datadog AWS Integration</h3>
        <p class="setup-desc text-secondary text-sm">
          Connect your AWS account to Datadog for infrastructure-level metrics and system dashboard correlation.
          This will deploy a CloudFormation stack to provision the secure cross-account role.
        </p>

        <!-- Error Alert -->
        <div v-if="integration.status === 'error' && integration.error" class="alert alert-danger mt-4 mb-2">
          <span class="alert-icon">⚠️</span>
          <div class="flex-col">
            <strong style="font-weight: 600;">Integration Setup Failed</strong>
            <span class="text-xs mt-1" style="word-break: break-word; line-height: 1.4;">{{ integration.error }}</span>
          </div>
        </div>

        <div class="method-selector mt-6">
          <h4 class="modal-section-title">Integration Resources Created</h4>
          <div class="metric-streams-resources p-5" style="background: var(--bg-tertiary); border: 1px solid var(--border-primary); border-radius: var(--radius-md);">
            <ul class="resource-list">
              <li>Datadog cross-account IAM Role (default: <code>DatadogIntegrationRole</code>)</li>
              <li>Datadog security trust policy delegation</li>
            </ul>
          </div>
        </div>

        <div class="alert alert-warning mt-6">
          ⚠️ Deploys Datadog's official CloudFormation main template to authorize account integration.
        </div>

        <div class="mt-6 flex gap-3">
          <button
            class="btn btn-primary btn-lg"
            :disabled="isSetup"
            @click="handleSetup"
            id="btn-setup-integration"
          >
            <span v-if="isSetup" class="spinner"></span>
            {{ isSetup ? 'Setting up...' : 'Setup Integration' }}
          </button>
        </div>
      </div>

      <!-- New Relic Setup Form -->
      <div v-else class="card-flat">
        <h3 class="setup-title">Setup New Relic AWS Integration</h3>
        <p class="setup-desc text-secondary text-sm">
          Connect your AWS account to New Relic for infrastructure-level monitoring.
          Choose between real-time Metric Streams or legacy API Polling.
        </p>

        <!-- Error Alert -->
        <div v-if="integration.status === 'error' && integration.error" class="alert alert-danger mt-4 mb-2">
          <span class="alert-icon">⚠️</span>
          <div class="flex-col">
            <strong style="font-weight: 600;">Integration Setup Failed</strong>
            <span class="text-xs mt-1" style="word-break: break-word; line-height: 1.4;">{{ integration.error }}</span>
          </div>
        </div>

        <div class="method-selector mt-6">
          <h4 class="modal-section-title">Choose Integration Method</h4>

          <div class="radio-group">
            <label class="radio-card" :class="{ active: selectedMethod === 'metric_streams' }">
              <input type="radio" class="radio" v-model="selectedMethod" value="metric_streams" />
              <div class="radio-card__content">
                <div class="flex items-center gap-2">
                  <strong>Metric Streams</strong>
                  <span class="badge badge-success btn-sm">Recommended</span>
                </div>
                <p class="text-xs text-secondary mt-1">
                  Real-time metrics via CloudWatch Metric Streams → Kinesis Firehose → New Relic.
                  ~2 minute latency.
                </p>
                <div class="metric-streams-resources mt-4">
                  <p class="text-xs text-secondary" style="font-weight: 600;">Resources created:</p>
                  <ul class="resource-list">
                    <li>Kinesis Data Firehose delivery stream</li>
                    <li>CloudWatch Metric Stream</li>
                    <li>S3 backup bucket for failed deliveries</li>
                    <li>IAM role for NR access</li>
                    <li>NR linked account (via NerdGraph)</li>
                  </ul>
                </div>

                <label class="checkbox-option mt-4" v-if="selectedMethod === 'metric_streams'">
                  <input type="checkbox" class="checkbox" v-model="includeLogs" />
                  <span class="text-sm">Include Log Forwarding via Firehose</span>
                </label>
              </div>
            </label>

            <label class="radio-card" :class="{ active: selectedMethod === 'api_polling' }">
              <input type="radio" class="radio" v-model="selectedMethod" value="api_polling" />
              <div class="radio-card__content">
                <div class="flex items-center gap-2">
                  <strong>API Polling</strong>
                  <span class="badge badge-neutral btn-sm">Legacy</span>
                </div>
                <p class="text-xs text-secondary mt-1">
                  NR polls CloudWatch APIs every 5 minutes. Lower setup complexity, higher latency.
                </p>
                <div class="metric-streams-resources mt-4">
                  <p class="text-xs text-secondary" style="font-weight: 600;">Resources created:</p>
                  <ul class="resource-list">
                    <li>IAM role for NR read-only access</li>
                  </ul>
                </div>
              </div>
            </label>
          </div>
        </div>

        <div class="alert alert-warning mt-6">
          ⚠️ {{ selectedMethod === 'metric_streams'
            ? 'Metric Streams creates billable AWS resources (Kinesis Firehose, S3 bucket, CloudWatch Metric Stream).'
            : 'API Polling creates an IAM role allowing NR to read your CloudWatch metrics.'
          }}
        </div>

        <div class="mt-6 flex gap-3">
          <button
            class="btn btn-primary btn-lg"
            :disabled="isSetup"
            @click="handleSetup"
            id="btn-setup-integration"
          >
            <span v-if="isSetup" class="spinner"></span>
            {{ isSetup ? 'Setting up...' : 'Setup Integration' }}
          </button>
        </div>
      </div>
    </div>

    <!-- In Progress -->
    <div v-if="integration.status === 'in_progress'" class="in-progress mt-6">
      <div class="card-flat text-center" style="padding: var(--space-12);">
        <span class="spinner" style="width: 48px; height: 48px; border-width: 4px; display: inline-block;"></span>
        <p class="mt-4 text-primary font-medium" style="font-weight: 500; font-size: var(--font-size-md);">
          {{ isSetup ? 'Initializing deployment & cleaning up old resources...' : 'CloudFormation stack is being deployed...' }}
        </p>
        <p class="text-xs text-secondary mt-2">
          This process usually takes between 1 to 3 minutes.
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { api } from '@/services/api'
import type { IntegrationInfo, IntegrationMethod } from '@/types'
import { useConfig } from '@/composables/useConfig'

const { provider } = useConfig()

const integration = ref<IntegrationInfo>({
  status: 'not_setup',
  lastChecked: new Date().toISOString(),
})

const selectedMethod = ref<IntegrationMethod>('metric_streams')
const includeLogs = ref(false)
const isSetup = ref(false)
const isRemoving = ref(false)

let pollInterval: number | null = null

function startPolling() {
  if (pollInterval) return
  pollInterval = window.setInterval(async () => {
    await loadIntegrationStatus()
    if (integration.value.status !== 'in_progress') {
      stopPolling()
    }
  }, 8000)
}

function stopPolling() {
  if (pollInterval) {
    clearInterval(pollInterval)
    pollInterval = null
  }
}

const statusLabel = computed(() => {
  const labels: Record<string, string> = {
    active: 'Active',
    not_setup: 'Not Setup',
    in_progress: 'Deploying...',
    error: 'Error',
  }
  return labels[integration.value.status] || integration.value.status
})

const statusBadgeClass = computed(() => {
  const classes: Record<string, string> = {
    active: 'badge-success',
    not_setup: 'badge-neutral',
    in_progress: 'badge-warning',
    error: 'badge-danger',
  }
  return classes[integration.value.status] || 'badge-neutral'
})

function formatTime(iso: string): string {
  try {
    return new Date(iso).toLocaleString()
  } catch {
    return iso
  }
}

async function loadIntegrationStatus() {
  try {
    integration.value = await api.integrationStatus()
    if (integration.value.status === 'in_progress') {
      startPolling()
    } else {
      stopPolling()
    }
  } catch (err) {
    console.error('Failed to load integration status:', err)
  }
}

async function handleSetup() {
  isSetup.value = true
  // Instantly show the loading progress card
  integration.value.status = 'in_progress'
  try {
    await api.setupIntegration({
      method: selectedMethod.value,
      includeLogs: includeLogs.value,
    })
    startPolling()
  } catch (err) {
    console.error('Setup failed:', err)
    integration.value.status = 'error'
    if (err && typeof err === 'object' && 'message' in err) {
      integration.value.error = (err as { message: string }).message
    } else {
      integration.value.error = String(err)
    }
  } finally {
    isSetup.value = false
  }
}

async function handleRemove() {
  if (!confirm('Remove the AWS Integration? This will delete the CloudFormation stack and all related resources.')) return

  isRemoving.value = true
  try {
    await api.removeIntegration()
    await loadIntegrationStatus()
  } catch (err) {
    console.error('Remove failed:', err)
  } finally {
    isRemoving.value = false
  }
}

onMounted(() => {
  loadIntegrationStatus()
})

onUnmounted(() => {
  stopPolling()
})
</script>

<style scoped>
.integration-title {
  font-size: var(--font-size-xl);
  font-weight: 600;
}

.status-details {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  background: var(--bg-primary);
  border-radius: var(--radius-md);
  padding: var(--space-4);
}

.detail-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.detail-label {
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
}

.detail-value {
  font-size: var(--font-size-sm);
  color: var(--text-primary);
}

.setup-title {
  font-size: var(--font-size-lg);
  font-weight: 600;
  margin-bottom: var(--space-2);
}

.setup-desc {
  line-height: 1.6;
}

.modal-section-title {
  font-size: var(--font-size-sm);
  font-weight: 600;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: var(--space-3);
}

.radio-group {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.radio-card {
  display: flex;
  align-items: flex-start;
  gap: var(--space-3);
  padding: var(--space-4) var(--space-5);
  background: var(--bg-primary);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.radio-card:hover {
  border-color: var(--text-tertiary);
}

.radio-card.active {
  border-color: var(--accent);
  background: var(--accent-dim);
}

.radio-card__content {
  flex: 1;
}

.radio-card__content strong {
  font-size: var(--font-size-md);
}

.resource-list {
  list-style: none;
  padding: 0;
  margin-top: var(--space-2);
}

.resource-list li {
  font-size: var(--font-size-xs);
  color: var(--text-secondary);
  padding: 2px 0;
}

.resource-list li::before {
  content: '→ ';
  color: var(--text-tertiary);
}

.checkbox-option {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  cursor: pointer;
}

.mt-1 { margin-top: 4px; }
</style>
