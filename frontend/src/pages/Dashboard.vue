<template>
  <div class="dashboard">
    <!-- Top Bar -->
    <header class="dash-header">
      <div class="dash-header-left">
        <div class="dash-brand">
          <div class="dash-logo-wrap">
            <img src="/logo.png" alt="Lambda Monitor" class="dash-logo-img" />
          </div>
          <div>
            <span class="dash-name">Lambda<span class="dash-name-accent"> Monitor</span></span>
            <span class="dash-connection">
              <span class="connection-dot"></span>
              {{ connectionLabel }}
            </span>
          </div>
        </div>
        <!-- Active connection / Region select dropdown -->
        <div v-if="regions.length > 0" class="connection-selector">
          <select 
            :value="activeRegionIndex" 
            @change="handleConnectionChange" 
            class="connection-select"
            id="select-connection"
          >
            <option 
              v-for="(r, idx) in regions" 
              :key="idx" 
              :value="idx"
            >
              {{ r.name || `${r.region} (${getEndpointHost(r.apiGatewayUrl)})` }}
            </option>
          </select>
        </div>
      </div>
      <div class="dash-header-right">
        <span v-if="userEmail" class="user-chip">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
          {{ userEmail }}
        </span>
        <button class="btn btn-ghost btn-icon" id="btn-refresh" @click="refreshData" title="Refresh">
          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="23 4 23 10 17 10"/><polyline points="1 20 1 14 7 14"/><path d="M3.51 9a9 9 0 0114.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0020.49 15"/></svg>
          Refresh
        </button>
        <button class="btn btn-ghost btn-icon btn-disconnect" id="btn-disconnect" @click="handleLogout" title="Disconnect">
          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 21H5a2 2 0 01-2-2V5a2 2 0 012-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>
          Disconnect
        </button>
      </div>
    </header>

    <!-- Breadcrumb -->
    <nav class="breadcrumb" aria-label="Breadcrumb">
      <router-link to="/provider" class="breadcrumb-link" id="breadcrumb-providers">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/></svg>
        Providers
      </router-link>
      <span class="breadcrumb-sep">
        <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="9 18 15 12 9 6"/></svg>
      </span>
      <router-link to="/login" class="breadcrumb-link" id="breadcrumb-config">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 010 2.83 2 2 0 01-2.83 0l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-4 0v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83 0 2 2 0 010-2.83l.06-.06A1.65 1.65 0 004.68 15a1.65 1.65 0 00-1.51-1H3a2 2 0 010-4h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 012.83-2.83l.06.06A1.65 1.65 0 009 4.68a1.65 1.65 0 001-1.51V3a2 2 0 014 0v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 2.83l-.06.06A1.65 1.65 0 0019.4 9a1.65 1.65 0 001.51 1H21a2 2 0 010 4h-.09a1.65 1.65 0 00-1.51 1z"/></svg>
        Configuration
      </router-link>
      <span class="breadcrumb-sep">
        <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="9 18 15 12 9 6"/></svg>
      </span>
      <span class="breadcrumb-current">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"/><line x1="3" y1="9" x2="21" y2="9"/><line x1="9" y1="21" x2="9" y2="9"/></svg>
        Dashboard
      </span>
    </nav>

    <!-- Tabs -->
    <nav class="tabs" id="dashboard-tabs">
      <button
        class="tab"
        :class="{ active: activeTab === 'instrumentation' }"
        @click="activeTab = 'instrumentation'"
        id="tab-instrumentation"
      >
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/></svg>
        Lambda Instrumentation
      </button>
      <button
        class="tab"
        :class="{ active: activeTab === 'integration' }"
        @click="activeTab = 'integration'"
        id="tab-integration"
      >
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M10 13a5 5 0 007.54.54l3-3a5 5 0 00-7.07-7.07l-1.72 1.71"/><path d="M14 11a5 5 0 00-7.54-.54l-3 3a5 5 0 007.07 7.07l1.71-1.71"/></svg>
        AWS Integration
      </button>
    </nav>

    <!-- Tab Content -->
    <div class="tab-content" :key="activeRegionIndex">
      <InstrumentationTab v-if="activeTab === 'instrumentation'" />
      <IntegrationTab v-if="activeTab === 'integration'" />
    </div>
  </div>
</template>


<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useConfig } from '@/composables/useConfig'
import InstrumentationTab from '@/components/InstrumentationTab.vue'
import IntegrationTab from '@/components/IntegrationTab.vue'

const router = useRouter()
const { 
  activeRegionIndex, 
  regions, 
  activeRegion, 
  setActiveRegion, 
  saveToSession, 
  userEmail, 
  logout: doLogout 
} = useConfig()

const activeTab = ref<'instrumentation' | 'integration'>('instrumentation')

const connectionLabel = computed(() => {
  if (activeRegion.value) {
    try {
      const url = new URL(activeRegion.value.apiGatewayUrl)
      return url.hostname
    } catch { return activeRegion.value.apiGatewayUrl }
  }
  return 'Not connected'
})

function getEndpointHost(urlStr: string): string {
  try {
    const url = new URL(urlStr)
    return url.hostname.split('.')[0]
  } catch {
    return urlStr
  }
}

function handleConnectionChange(e: Event) {
  const select = e.target as HTMLSelectElement
  const index = Number(select.value)
  setActiveRegion(index)
  saveToSession()
}

function refreshData() {
  // In real app, re-fetch functions and integration data
}

function handleLogout() {
  doLogout()
  router.push('/')
}
</script>

<style scoped>
.dashboard {
  min-height: 100vh;
  background: var(--bg-primary);
}

/* ── Header ── */
.dash-header {
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

.dash-header-left {
  display: flex;
  align-items: center;
  gap: var(--space-4);
}

.dash-brand {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

/* ─ New Logo ─ */
.dash-logo-wrap {
  width: 40px;
  height: 40px;
  border-radius: var(--radius-md);
  overflow: hidden;
  flex-shrink: 0;
  box-shadow: 0 0 0 1px rgba(28, 231, 131, 0.2), 0 0 14px rgba(28, 231, 131, 0.1);
  transition: box-shadow var(--transition-base);
}

.dash-logo-wrap:hover {
  box-shadow: 0 0 0 1px rgba(28, 231, 131, 0.4), 0 0 20px rgba(28, 231, 131, 0.18);
}

.dash-logo-img {
  width: 40px;
  height: 40px;
  display: block;
  object-fit: cover;
}

.dash-name {
  display: block;
  font-weight: 700;
  font-size: var(--font-size-md);
  letter-spacing: -0.02em;
  color: var(--text-primary);
}

.dash-name-accent {
  color: var(--accent);
}

.dash-connection {
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
  background: var(--accent);
  flex-shrink: 0;
  box-shadow: 0 0 5px var(--accent);
  animation: pulse-dot 2s ease infinite;
}

@keyframes pulse-dot {
  0%, 100% { opacity: 1; box-shadow: 0 0 5px var(--accent); }
  50% { opacity: 0.7; box-shadow: 0 0 9px var(--accent); }
}

.dash-header-right {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.btn-icon {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  font-size: var(--font-size-sm);
}

.btn-disconnect:hover:not(:disabled) {
  color: var(--danger);
  background: var(--danger-dim);
}

/* Region selector & user chips */
.connection-selector {
  display: inline-flex;
  align-items: center;
}

.connection-select {
  background: rgba(22, 27, 34, 0.85);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-full);
  padding: 4px 14px;
  color: var(--accent);
  font-size: var(--font-size-xs);
  font-family: var(--font-mono);
  font-weight: 500;
  outline: none;
  cursor: pointer;
  transition: all var(--transition-fast);
  box-shadow: 0 0 0 1px rgba(28, 231, 131, 0.1);
}

.connection-select:hover {
  border-color: var(--accent);
  box-shadow: 0 0 8px rgba(28, 231, 131, 0.2);
}

.region-chip {
  display: inline-flex;
  align-items: center;
  gap: var(--space-1);
  padding: 4px 10px;
  background: var(--accent-dim);
  border: 1px solid var(--accent-muted);
  border-radius: var(--radius-full);
  font-size: var(--font-size-xs);
  font-family: var(--font-mono);
  color: var(--accent);
  font-weight: 500;
}

.user-chip {
  display: inline-flex;
  align-items: center;
  gap: var(--space-1);
  padding: 4px 10px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-full);
  font-size: var(--font-size-xs);
  color: var(--text-secondary);
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* ── Breadcrumb ── */
.breadcrumb {
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

.breadcrumb-link:hover {
  color: var(--accent);
  background: var(--accent-dim);
}

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

/* ── Tabs ── */
.tabs {
  display: flex;
  gap: 0;
  padding: 0 var(--space-6);
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-primary);
}

.tab {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-3) var(--space-5);
  font-family: var(--font-family);
  font-size: var(--font-size-sm);
  font-weight: 500;
  color: var(--text-tertiary);
  background: none;
  border: none;
  border-bottom: 2px solid transparent;
  cursor: pointer;
  transition: all var(--transition-fast);
  position: relative;
}

.tab:hover {
  color: var(--text-primary);
}

.tab.active {
  color: var(--accent);
  border-bottom-color: var(--accent);
}

/* ── Content ── */
.tab-content {
  padding: var(--space-6);
  max-width: 1400px;
  margin: 0 auto;
}
</style>
