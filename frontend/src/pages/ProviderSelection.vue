<template>
  <div class="provider-page">

    <!-- ── Toast ── -->
    <Transition name="toast">
      <div v-if="toast.show" class="toast" :class="'toast--' + toast.type">
        <div class="toast__icon">
          <svg v-if="toast.type === 'info'" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
          <svg v-else-if="toast.type === 'success'" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"/></svg>
        </div>
        <div class="toast__body">
          <p class="toast__title">{{ toast.title }}</p>
          <p class="toast__message">{{ toast.message }}</p>
        </div>
        <button class="toast__close" @click="dismissToast">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>
    </Transition>

    <!-- Shared App Header -->
    <header class="app-header">
      <div class="app-header-left">
        <div class="app-brand">
          <div class="app-logo-wrap">
            <img src="/logo.png" alt="Lambda Monitor" class="app-logo-img" />
          </div>
          <div>
            <span class="app-name">Lambda<span class="app-name-accent"> Monitor</span></span>
            <span class="app-tagline">Serverless Observability Platform</span>
          </div>
        </div>
      </div>
      <div class="app-header-right">
        <router-link to="/auth" class="btn btn-ghost btn-sm btn-icon" id="btn-back-auth">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="15 18 9 12 15 6"/></svg>
          Back
        </router-link>
      </div>
    </header>

    <!-- Breadcrumb -->
    <nav class="app-breadcrumb" aria-label="Breadcrumb">
      <span class="breadcrumb-current">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/></svg>
        Select Provider
      </span>
    </nav>

    <!-- Page Content -->
    <div class="provider-content">
      <div class="provider-main animate-slide-up">

        <!-- Heading -->
        <!-- <div class="provider-heading">
          <div class="provider-heading-icon">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/></svg>
          </div>
          <div>
            <h1 class="provider-title">Choose Your Provider</h1>
            <p class="provider-subtitle">Select the observability platform to instrument your AWS Lambda functions</p>
          </div>
        </div> -->

        <!-- Step Indicator -->
        <div class="step-indicator">
          <div class="step-item step-active">
            <div class="step-dot">1</div>
            <span>Provider</span>
          </div>
          <div class="step-line"></div>
          <div class="step-item">
            <div class="step-dot">2</div>
            <span>Configuration</span>
          </div>
          <div class="step-line"></div>
          <div class="step-item">
            <div class="step-dot">3</div>
            <span>Dashboard</span>
          </div>
        </div>

        <!-- Provider Cards -->
        <div class="providers-grid">

          <!-- ── New Relic ── -->
          <button
            class="provider-card provider-card--newrelic"
            @click="selectProvider('newrelic')"
            id="provider-newrelic"
          >
            <div class="card-glow card-glow--green"></div>
            <div class="card-top">
              <div class="card-logo">
                <svg viewBox="0 0 40 40" width="52" height="52" fill="none">
                  <path d="M34 13v14L20 35V21L34 13z" fill="#1CE783"/>
                  <path d="M20 21V7l14 7L20 21z" fill="#1CE783" opacity="0.7"/>
                  <path d="M20 7v14L6 13.5 20 7z" fill="#1CE783" opacity="0.4"/>
                </svg>
              </div>
              <span class="card-badge card-badge--available">
                <span class="card-badge-dot"></span>
                Available
              </span>
            </div>

            <div class="card-body">
              <h3 class="card-name">New Relic</h3>
              <p class="card-desc">
                Full-stack observability with Lambda layers, APM distributed
                tracing, Metric Streams, and log ingestion — all in one platform.
              </p>
            </div>

            <ul class="card-features">
              <li>
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="20 6 9 17 4 12"/></svg>
                Layer-based instrumentation
              </li>
              <li>
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="20 6 9 17 4 12"/></svg>
                APM &amp; Serverless telemetry modes
              </li>
              <li>
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="20 6 9 17 4 12"/></svg>
                Real-time Metric Streams
              </li>
              <li>
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="20 6 9 17 4 12"/></svg>
                Multi-region &amp; multi-account support
              </li>
            </ul>

            <div class="card-action">
              <span class="btn btn-primary card-cta">
                Select New Relic
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="5" y1="12" x2="19" y2="12"/><polyline points="12 5 19 12 12 19"/></svg>
              </span>
            </div>
          </button>

          <!-- ── Datadog (Coming Soon) ── -->
          <div class="provider-card provider-card--datadog" id="provider-datadog">
            <div class="card-glow card-glow--purple"></div>
            <div class="card-top">
              <div class="card-logo">
                <svg viewBox="0 0 40 40" width="52" height="52" fill="none">
                  <rect x="5" y="5" width="30" height="30" rx="6" fill="#632CA6" opacity="0.25"/>
                  <path d="M13 21l5 5 9-11" stroke="#9b72d4" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
              </div>
              <span class="card-badge card-badge--soon">Coming Soon</span>
            </div>

            <div class="card-body">
              <h3 class="card-name">Datadog</h3>
              <p class="card-desc">
                Lambda monitoring, distributed tracing, and log management
                powered by the Datadog agent and Lambda extension.
              </p>
            </div>

            <ul class="card-features card-features--muted">
              <li>
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="4"/></svg>
                Datadog Lambda extension
              </li>
              <li>
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="4"/></svg>
                Enhanced Lambda metrics
              </li>
              <li>
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="4"/></svg>
                Log management &amp; APM traces
              </li>
              <li>
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="4"/></svg>
                Serverless monitoring
              </li>
            </ul>

            <div class="card-action">
              <button class="notify-btn" @click.stop="handleNotifyMe" id="btn-notify-datadog">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M18 8A6 6 0 006 8c0 7-3 9-3 9h18s-3-2-3-9M13.73 21a2 2 0 01-3.46 0"/></svg>
                Notify Me
              </button>
            </div>
          </div>

        </div>

        <!-- Footer note -->
        <p class="provider-note">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
          Secure · Proprietary · Your keys never leave your AWS account
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useConfig } from '@/composables/useConfig'
import type { Provider } from '@/types'

const router = useRouter()
const { setProvider, saveToSession } = useConfig()

function selectProvider(p: Provider) {
  setProvider(p)
  saveToSession()
  router.push('/login')
}

// ── Toast ──
interface Toast {
  show: boolean
  type: 'info' | 'success' | 'error'
  title: string
  message: string
}

const toast = ref<Toast>({
  show: false,
  type: 'info',
  title: '',
  message: '',
})

let toastTimer: ReturnType<typeof setTimeout> | null = null

function showToast(type: Toast['type'], title: string, message: string) {
  if (toastTimer) clearTimeout(toastTimer)
  toast.value = { show: true, type, title, message }
  toastTimer = setTimeout(() => {
    toast.value.show = false
  }, 4000)
}

function dismissToast() {
  if (toastTimer) clearTimeout(toastTimer)
  toast.value.show = false
}

function handleNotifyMe() {
  showToast('info', "You're on the list!", "We'll notify you as soon as Datadog support launches.")
}
</script>

<style scoped>
.provider-page {
  min-height: 100vh;
  background: var(--bg-primary);
}

/* ━━ Toast ━━ */
.toast {
  position: fixed;
  bottom: 28px;
  right: 28px;
  z-index: 9999;
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 14px 16px;
  border-radius: 14px;
  border: 1px solid transparent;
  backdrop-filter: blur(20px);
  box-shadow: 0 8px 32px rgba(0,0,0,0.45), 0 2px 8px rgba(0,0,0,0.3);
  min-width: 300px;
  max-width: 380px;
}

.toast--info {
  background: rgba(22, 27, 40, 0.95);
  border-color: rgba(99, 44, 166, 0.5);
}

.toast--success {
  background: rgba(22, 40, 30, 0.95);
  border-color: rgba(28, 231, 131, 0.4);
}

.toast__icon {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.toast--info .toast__icon {
  background: rgba(99, 44, 166, 0.25);
  color: #a78bfa;
}

.toast--success .toast__icon {
  background: var(--success-dim);
  color: var(--accent);
}

.toast__body {
  flex: 1;
  min-width: 0;
}

.toast__title {
  font-weight: 600;
  font-size: var(--font-size-sm);
  color: var(--text-primary);
  margin-bottom: 2px;
}

.toast--info .toast__title { color: #c4b5fd; }

.toast__message {
  font-size: var(--font-size-xs);
  color: var(--text-secondary);
  line-height: 1.5;
}

.toast__close {
  background: none;
  border: none;
  color: var(--text-tertiary);
  cursor: pointer;
  padding: 2px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  line-height: 1;
  flex-shrink: 0;
  margin-top: 1px;
  transition: color var(--transition-fast);
}

.toast__close:hover { color: var(--text-primary); }

/* Toast transition */
.toast-enter-active {
  animation: toast-in 0.35s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
}
.toast-leave-active {
  animation: toast-out 0.25s ease forwards;
}

@keyframes toast-in {
  from { opacity: 0; transform: translateY(16px) scale(0.96); }
  to   { opacity: 1; transform: translateY(0) scale(1); }
}
@keyframes toast-out {
  from { opacity: 1; transform: translateY(0) scale(1); }
  to   { opacity: 0; transform: translateY(8px) scale(0.97); }
}

/* ━━ Shared App Header ━━ */
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

.app-logo-img { width: 40px; height: 40px; display: block; object-fit: cover; }

.app-name {
  display: block;
  font-weight: 700;
  font-size: var(--font-size-md);
  letter-spacing: -0.02em;
}
.app-name-accent { color: var(--accent); }

.app-tagline {
  display: block;
  font-size: 0.6rem;
  color: var(--text-tertiary);
  letter-spacing: 0.06em;
  text-transform: uppercase;
  margin-top: 1px;
}

.app-header-right { display: flex; align-items: center; gap: var(--space-2); }

.btn-icon {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  font-size: var(--font-size-sm);
}

/* ━━ Breadcrumb ━━ */
.app-breadcrumb {
  display: flex;
  align-items: center;
  gap: var(--space-1);
  padding: var(--space-2) var(--space-6);
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-primary);
  font-size: var(--font-size-xs);
}

.breadcrumb-current {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  color: var(--text-primary);
  font-weight: 600;
  padding: 3px var(--space-2);
}

/* ━━ Page Layout ━━ */
.provider-content {
  padding: var(--space-10) var(--space-6);
  max-width: 900px;
  margin: 0 auto;
}

.provider-main {
  display: flex;
  flex-direction: column;
  gap: var(--space-8);
}

/* ━━ Heading ━━ */
.provider-heading {
  display: flex;
  align-items: center;
  gap: var(--space-4);
}

.provider-heading-icon {
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

.provider-title {
  font-size: var(--font-size-2xl);
  font-weight: 700;
  letter-spacing: -0.02em;
  margin-bottom: var(--space-1);
}

.provider-subtitle {
  color: var(--text-secondary);
  font-size: var(--font-size-sm);
}

/* ━━ Step Indicator ━━ */
.step-indicator {
  display: flex;
  align-items: center;
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

/* ━━ Provider Grid ━━ */
.providers-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--space-6);
  align-items: stretch;
}

/* ━━ Provider Card ━━ */
.provider-card {
  position: relative;
  background: var(--bg-secondary);
  border: 1px solid var(--border-primary);
  border-radius: 18px;
  padding: 28px;
  text-align: left;
  transition: all 0.22s ease;
  font-family: var(--font-family);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* Glow blob */
.card-glow {
  position: absolute;
  top: -80px;
  right: -80px;
  width: 240px;
  height: 240px;
  border-radius: 50%;
  pointer-events: none;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.card-glow--green {
  background: radial-gradient(circle, rgba(28, 231, 131, 0.12) 0%, transparent 70%);
}

.card-glow--purple {
  background: radial-gradient(circle, rgba(99, 44, 166, 0.18) 0%, transparent 70%);
}

/* ─ New Relic card ─ */
.provider-card--newrelic {
  cursor: pointer;
}

.provider-card--newrelic:hover {
  border-color: rgba(28, 231, 131, 0.6);
  box-shadow: 0 0 0 1px rgba(28, 231, 131, 0.15), 0 12px 40px rgba(28, 231, 131, 0.1);
  transform: translateY(-3px);
}

.provider-card--newrelic:hover .card-glow--green {
  opacity: 1;
}

/* ─ Datadog card ─ */
.provider-card--datadog {
  cursor: pointer;
}

.provider-card--datadog:hover {
  border-color: rgba(99, 44, 166, 0.6);
  box-shadow: 0 0 0 1px rgba(99, 44, 166, 0.2), 0 12px 40px rgba(99, 44, 166, 0.1);
  transform: translateY(-3px);
}

.provider-card--datadog:hover .card-glow--purple {
  opacity: 1;
}

/* Card Top */
.card-top {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
}

.card-logo {
  display: flex;
  align-items: center;
}

/* Badges */
.card-badge {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 4px 10px;
  border-radius: var(--radius-full);
  font-size: 0.65rem;
  font-weight: 600;
  letter-spacing: 0.04em;
  text-transform: uppercase;
  flex-shrink: 0;
}

.card-badge--available {
  background: rgba(28, 231, 131, 0.15);
  border: 1px solid rgba(28, 231, 131, 0.35);
  color: var(--accent);
}

.card-badge-dot {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: var(--accent);
  animation: pulse-dot 2s ease infinite;
  flex-shrink: 0;
}

@keyframes pulse-dot {
  0%, 100% { opacity: 1; box-shadow: 0 0 4px var(--accent); }
  50% { opacity: 0.6; box-shadow: 0 0 8px var(--accent); }
}

.card-badge--soon {
  background: rgba(99, 44, 166, 0.15);
  border: 1px solid rgba(99, 44, 166, 0.35);
  color: #a78bfa;
}

/* Card Body */
.card-body {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.card-name {
  font-size: 1.3rem;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0;
  letter-spacing: -0.01em;
}

.card-desc {
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
  line-height: 1.7;
  margin: 0;
}

/* Features */
.card-features {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 10px;
  flex: 1;
}

.card-features li {
  display: flex;
  align-items: center;
  gap: 9px;
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
  line-height: 1.4;
}

.provider-card--newrelic .card-features svg {
  color: var(--accent);
  flex-shrink: 0;
}

.card-features--muted li {
  opacity: 0.65;
}

.card-features--muted svg {
  color: #9b72d4;
  flex-shrink: 0;
}

/* Card Action */
.card-action {
  padding-top: 6px;
  border-top: 1px solid var(--border-primary);
}

.card-cta {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  width: 100%;
  justify-content: center;
  padding: 11px 20px;
  font-size: var(--font-size-sm);
  font-weight: 600;
}

/* Notify Me button */
.notify-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 7px;
  width: 100%;
  padding: 11px 20px;
  border-radius: var(--radius-md);
  border: 1px solid rgba(99, 44, 166, 0.4);
  background: rgba(99, 44, 166, 0.1);
  color: #c4b5fd;
  font-size: var(--font-size-sm);
  font-weight: 600;
  cursor: pointer;
  transition: all 0.18s ease;
  font-family: var(--font-family);
}

.notify-btn:hover {
  background: rgba(99, 44, 166, 0.22);
  border-color: rgba(99, 44, 166, 0.65);
  color: #ddd6fe;
  box-shadow: 0 0 12px rgba(99, 44, 166, 0.2);
  transform: none;
}

/* ━━ Footer Note ━━ */
.provider-note {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  font-size: var(--font-size-xs);
  color: var(--text-tertiary);
}

/* ━━ Responsive ━━ */
@media (max-width: 680px) {
  .providers-grid {
    grid-template-columns: 1fr;
  }
  .provider-content {
    padding: var(--space-6) var(--space-4);
  }
}
</style>
