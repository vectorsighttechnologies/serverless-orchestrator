<template>
  <div class="instrument-page">
    <!-- Shared App Header -->
    <header class="app-header">
      <div class="app-header-left">
        <div class="app-brand">
          <div class="app-logo-wrap">
            <img src="/logo.png" alt="Lambda Monitor" class="app-logo-img" />
          </div>
          <div>
            <span class="app-name">Lambda<span class="app-name-accent"> Monitor</span></span>
            <span class="app-connection">
              <span class="connection-dot connection-dot--live"></span>
              {{ connectionLabel }}
            </span>
          </div>
        </div>
      </div>
      <div class="app-header-right">
        <router-link to="/dashboard" class="btn btn-ghost btn-sm btn-icon" id="btn-back-dashboard">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="15 18 9 12 15 6"/></svg>
          Back to Dashboard
        </router-link>
      </div>
    </header>

    <!-- Breadcrumb -->
    <nav class="app-breadcrumb" aria-label="Breadcrumb">
      <router-link to="/provider" class="breadcrumb-link">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/></svg>
        Providers
      </router-link>
      <span class="breadcrumb-sep"><svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="9 18 15 12 9 6"/></svg></span>
      <router-link to="/dashboard" class="breadcrumb-link">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="18" height="18" rx="2"/><line x1="3" y1="9" x2="21" y2="9"/><line x1="9" y1="21" x2="9" y2="9"/></svg>
        Dashboard
      </router-link>
      <span class="breadcrumb-sep"><svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="9 18 15 12 9 6"/></svg></span>
      <span class="breadcrumb-current">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/></svg>
        Instrument Functions
      </span>
    </nav>

    <!-- Page Content -->
    <div class="instrument-content">
      <div class="animate-slide-up">

        <!-- Page Heading -->
        <div class="instrument-heading">
          <div>
            <h1 class="instrument-title">Instrument Functions</h1>
            <p class="instrument-subtitle">
              Instrumenting <strong class="text-accent">{{ selectedArns.length }} function{{ selectedArns.length !== 1 ? 's' : '' }}</strong> with New Relic observability
            </p>
          </div>
        </div>

        <div class="instrument-layout">
          <!-- Left: Config Columns -->
          <div class="config-column">

            <!-- Method -->
            <div class="section-block">
              <h2 class="section-title">
                <span class="section-title__num">1</span>
                Instrumentation Method
              </h2>
              <div class="option-list">
                <label class="option-card" :class="{ active: instrumentMethod === 'layer' }">
                  <input type="radio" class="radio" name="method" value="layer" v-model="instrumentMethod" />
                  <div class="option-card__icon">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="7" width="20" height="14" rx="2"/><path d="M16 21V5a2 2 0 00-2-2h-4a2 2 0 00-2 2v16"/></svg>
                  </div>
                  <div class="option-card__body">
                    <div class="option-card__title">Layer-Based <span class="badge badge-success" style="font-size:0.6rem;margin-left:6px">Recommended</span></div>
                    <p class="option-card__desc">Adds New Relic Lambda layer to your function with the extension sidecar. Zero code changes required.</p>
                  </div>
                </label>
                <label class="option-card" :class="{ active: instrumentMethod === 'log_ingestion' }">
                  <input type="radio" class="radio" name="method" value="log_ingestion" v-model="instrumentMethod" />
                  <div class="option-card__icon">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/></svg>
                  </div>
                  <div class="option-card__body">
                    <div class="option-card__title">Log Ingestion</div>
                    <p class="option-card__desc">Deploys a separate Lambda that forwards CloudWatch logs to New Relic. Good for unsupported runtimes.</p>
                  </div>
                </label>
              </div>
            </div>

            <!-- Mode (only for layer method) -->
            <div class="section-block" v-if="instrumentMethod === 'layer'">
              <h2 class="section-title">
                <span class="section-title__num">2</span>
                Telemetry Mode
              </h2>
              <div class="option-list">
                <label class="option-card" :class="{ active: instrumentMode === 'serverless' }">
                  <input type="radio" class="radio" name="mode" value="serverless" v-model="instrumentMode" />
                  <div class="option-card__icon">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/></svg>
                  </div>
                  <div class="option-card__body">
                    <div class="option-card__title">Serverless <span class="badge badge-success" style="font-size:0.6rem;margin-left:6px">Recommended</span></div>
                    <p class="option-card__desc">Extension-based, lightweight telemetry. Minimal cold-start overhead. Best for most workloads.</p>
                  </div>
                </label>
                <label class="option-card" :class="{ active: instrumentMode === 'apm' }">
                  <input type="radio" class="radio" name="mode" value="apm" v-model="instrumentMode" />
                  <div class="option-card__icon">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="20" x2="18" y2="10"/><line x1="12" y1="20" x2="12" y2="4"/><line x1="6" y1="20" x2="6" y2="14"/></svg>
                  </div>
                  <div class="option-card__body">
                    <div class="option-card__title">APM</div>
                    <p class="option-card__desc">Full distributed tracing with APM features. Higher fidelity data, slightly more overhead.</p>
                  </div>
                </label>
              </div>
            </div>
          </div>

          <!-- Right: Summary Panel -->
          <div class="summary-panel">
            <div class="summary-panel__inner">
              <h3 class="summary-panel__title">Summary</h3>

              <div class="summary-row-item">
                <span class="summary-row-item__label">Functions</span>
                <span class="summary-row-item__value text-accent">{{ selectedArns.length }}</span>
              </div>
              <div class="summary-row-item">
                <span class="summary-row-item__label">Method</span>
                <span class="mode-badge" :class="'mode-' + instrumentMethod">{{ methodLabel }}</span>
              </div>
              <div class="summary-row-item" v-if="instrumentMethod === 'layer'">
                <span class="summary-row-item__label">Mode</span>
                <span class="mode-badge" :class="'mode-' + instrumentMode">{{ modeLabel }}</span>
              </div>

              <div class="summary-divider"></div>

              <div class="fn-preview">
                <p class="fn-preview__label">Selected functions</p>
                <div class="fn-preview__list">
                  <div v-for="arn in selectedArns.slice(0, 5)" :key="arn" class="fn-preview__item font-mono">
                    {{ arnToName(arn) }}
                  </div>
                  <div v-if="selectedArns.length > 5" class="fn-preview__more">
                    +{{ selectedArns.length - 5 }} more
                  </div>
                </div>
              </div>

              <div class="summary-divider"></div>

              <button
                class="btn btn-primary instrument-submit-btn"
                @click="doInstrument"
                :disabled="isSubmitting"
                id="btn-confirm-instrument"
              >
                <span v-if="isSubmitting" class="spinner"></span>
                <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/></svg>
                {{ isSubmitting ? 'Instrumenting...' : `Instrument ${selectedArns.length} Function${selectedArns.length !== 1 ? 's' : ''}` }}
              </button>
              <button class="btn btn-ghost" style="width:100%;margin-top:var(--space-2)" @click="router.push('/dashboard')">
                Cancel
              </button>
            </div>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useConfig } from '@/composables/useConfig'
import type { InstrumentMethod } from '@/types'

const router = useRouter()
const route = useRoute()
const { activeRegion } = useConfig()

// Parse selected ARNs from query params
const selectedArns = computed<string[]>(() => {
  const q = route.query.arns
  if (!q) return []
  if (Array.isArray(q)) return q.filter(Boolean) as string[]
  return [q as string]
})

const instrumentMethod = ref<InstrumentMethod>('layer')
const instrumentMode = ref<'serverless' | 'apm'>('serverless')
const isSubmitting = ref(false)

const connectionLabel = computed(() => {
  if (activeRegion.value) {
    try {
      const url = new URL(activeRegion.value.apiGatewayUrl)
      return url.hostname
    } catch { return activeRegion.value.apiGatewayUrl }
  }
  return 'Not connected'
})

const methodLabel = computed(() => {
  return instrumentMethod.value === 'layer' ? 'Layer-Based' : 'Log Ingestion'
})

const modeLabel = computed(() => {
  return instrumentMode.value === 'serverless' ? 'Serverless' : 'APM'
})

function arnToName(arn: string): string {
  const parts = arn.split(':')
  return parts[parts.length - 1] || arn
}

async function doInstrument() {
  isSubmitting.value = true
  // Simulate the API call
  await new Promise(resolve => setTimeout(resolve, 1200))
  isSubmitting.value = false
  router.push('/dashboard')
}
</script>

<style scoped>
.instrument-page {
  min-height: 100vh;
  background: var(--bg-primary);
}

/* ── Shared Header (same as LoginConfig) ── */
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
  width: 40px; height: 40px;
  border-radius: var(--radius-md);
  overflow: hidden; flex-shrink: 0;
  box-shadow: 0 0 0 1px rgba(28,231,131,0.2), 0 0 14px rgba(28,231,131,0.1);
}
.app-logo-img { width: 40px; height: 40px; display: block; object-fit: cover; }
.app-name { display: block; font-weight: 700; font-size: var(--font-size-md); letter-spacing: -0.02em; }
.app-name-accent { color: var(--accent); }
.app-connection {
  display: flex; align-items: center; gap: var(--space-1);
  font-size: 0.65rem; color: var(--text-tertiary);
  font-family: var(--font-mono); margin-top: 1px;
}
.connection-dot {
  width: 6px; height: 6px;
  border-radius: 50%; background: var(--text-tertiary); flex-shrink: 0;
}
.connection-dot--live {
  background: var(--accent);
  box-shadow: 0 0 5px var(--accent);
  animation: pulse-dot 2s ease infinite;
}
@keyframes pulse-dot {
  0%, 100% { opacity: 1; box-shadow: 0 0 5px var(--accent); }
  50% { opacity: 0.7; box-shadow: 0 0 9px var(--accent); }
}
.app-header-right { display: flex; align-items: center; gap: var(--space-2); }
.btn-icon { display: inline-flex; align-items: center; gap: var(--space-2); font-size: var(--font-size-sm); }

/* ── Breadcrumb ── */
.app-breadcrumb {
  display: flex; align-items: center; gap: var(--space-1);
  padding: var(--space-2) var(--space-6);
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-primary);
  font-size: var(--font-size-xs);
}
.breadcrumb-link {
  display: inline-flex; align-items: center; gap: 5px;
  color: var(--text-tertiary); text-decoration: none;
  transition: color var(--transition-fast);
  padding: 3px var(--space-2); border-radius: var(--radius-sm);
}
.breadcrumb-link:hover { color: var(--accent); background: var(--accent-dim); }
.breadcrumb-sep { color: var(--border-hover); display: flex; align-items: center; margin: 0 1px; }
.breadcrumb-current {
  display: inline-flex; align-items: center; gap: 5px;
  color: var(--text-primary); font-weight: 600; padding: 3px var(--space-2);
}

/* ── Layout ── */
.instrument-content {
  padding: var(--space-8) var(--space-6);
  max-width: 1200px;
  margin: 0 auto;
}

.instrument-heading {
  margin-bottom: var(--space-8);
}

.instrument-title {
  font-size: var(--font-size-2xl);
  font-weight: 700;
  letter-spacing: -0.02em;
  margin-bottom: var(--space-2);
}

.instrument-subtitle {
  font-size: var(--font-size-md);
  color: var(--text-secondary);
}

.instrument-layout {
  display: grid;
  grid-template-columns: 1fr 340px;
  gap: var(--space-6);
  align-items: start;
}

/* ── Config Column ── */
.config-column {
  display: flex;
  flex-direction: column;
  gap: var(--space-6);
}

.section-block {
  background: var(--bg-secondary);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-lg);
  padding: var(--space-6);
}

.section-title {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  font-size: var(--font-size-sm);
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: var(--text-secondary);
  margin-bottom: var(--space-5);
}

.section-title__num {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: var(--accent-dim);
  border: 1px solid var(--accent);
  color: var(--accent);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.7rem;
  font-weight: 700;
  flex-shrink: 0;
}

.option-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.option-card {
  display: flex;
  align-items: flex-start;
  gap: var(--space-4);
  padding: var(--space-4) var(--space-5);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all var(--transition-fast);
  background: var(--bg-primary);
}

.option-card:hover { border-color: var(--border-hover); }

.option-card.active {
  border-color: var(--accent);
  background: var(--accent-dim);
}

.option-card__icon {
  width: 40px;
  height: 40px;
  border-radius: var(--radius-md);
  background: var(--bg-tertiary);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-secondary);
  flex-shrink: 0;
  transition: all var(--transition-fast);
}

.option-card.active .option-card__icon {
  background: rgba(28, 231, 131, 0.15);
  color: var(--accent);
}

.option-card__body { flex: 1; }

.option-card__title {
  font-weight: 600;
  font-size: var(--font-size-sm);
  margin-bottom: var(--space-1);
  display: flex;
  align-items: center;
}

.option-card__desc {
  font-size: var(--font-size-xs);
  color: var(--text-secondary);
  line-height: 1.6;
}

/* ── Summary Panel ── */
.summary-panel {
  position: sticky;
  top: calc(58px + 32px + var(--space-8));
}

.summary-panel__inner {
  background: var(--bg-secondary);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-lg);
  padding: var(--space-6);
}

.summary-panel__title {
  font-size: var(--font-size-sm);
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: var(--text-tertiary);
  margin-bottom: var(--space-5);
}

.summary-row-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-2) 0;
  font-size: var(--font-size-sm);
}

.summary-row-item__label { color: var(--text-secondary); }
.summary-row-item__value { font-weight: 600; }

.summary-divider {
  height: 1px;
  background: var(--border-primary);
  margin: var(--space-4) 0;
}

.fn-preview__label {
  font-size: var(--font-size-xs);
  color: var(--text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: var(--space-2);
}

.fn-preview__item {
  font-size: var(--font-size-xs);
  color: var(--text-secondary);
  padding: 3px 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.fn-preview__more {
  font-size: var(--font-size-xs);
  color: var(--accent);
  margin-top: var(--space-1);
}

.instrument-submit-btn {
  width: 100%;
  justify-content: center;
}

/* ── Mode Badges ── */
.mode-badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: var(--radius-full);
  font-size: var(--font-size-xs);
  font-weight: 500;
}

.mode-layer { background: var(--accent-dim); color: var(--accent); }
.mode-serverless { background: var(--accent-dim); color: var(--accent); }
.mode-apm { background: var(--info-dim); color: var(--info); }
.mode-log_ingestion { background: var(--warning-dim); color: var(--warning); }

@media (max-width: 768px) {
  .instrument-layout { grid-template-columns: 1fr; }
  .summary-panel { position: static; }
}
</style>
