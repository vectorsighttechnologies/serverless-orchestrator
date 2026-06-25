<template>
  <div class="inst-tab">
    <!-- Summary Cards -->
    <div class="summary-row">
      <div class="summary-card summary-card--total">
        <div class="summary-icon summary-icon--total">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/></svg>
        </div>
        <div class="summary-text">
          <span class="summary-num">{{ functions.length }}</span>
          <span class="summary-label">Total Functions</span>
        </div>
      </div>
      <div class="summary-card summary-card--instrumented">
        <div class="summary-icon summary-icon--instrumented">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 11.08V12a10 10 0 11-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
        </div>
        <div class="summary-text">
          <span class="summary-num">{{ instrumentedCount }}</span>
          <span class="summary-label">Instrumented</span>
          <div class="summary-progress">
            <div class="summary-progress-bar" :style="{ width: instrumentedPercent + '%' }"></div>
          </div>
          <span class="summary-pct">{{ instrumentedPercent }}%</span>
        </div>
      </div>
      <div class="summary-card summary-card--not">
        <div class="summary-icon summary-icon--not">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg>
        </div>
        <div class="summary-text">
          <span class="summary-num">{{ notInstrumentedCount }}</span>
          <span class="summary-label">Not Instrumented</span>
        </div>
      </div>
    </div>

    <!-- Search & Filters -->
    <div class="toolbar">
      <div class="search-box">
        <svg class="search-icon" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
        <input
          v-model="search"
          type="text"
          class="search-input"
          id="search-functions"
          placeholder="Search functions, ARNs..."
        />
        <button v-if="search" class="search-clear" @click="search = ''" title="Clear search">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>
      <div class="filters">
        <svg class="filter-icon" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3"/></svg>
        <span class="filter-label">Filters:</span>
        <select v-model="filterRegion" class="filter-select" id="filter-region">
          <option value="">All Regions</option>
          <option v-for="r in allRegions" :key="r" :value="r">{{ r }}</option>
        </select>
        <select v-model="filterRuntime" class="filter-select" id="filter-runtime">
          <option value="">All Runtimes</option>
          <option v-for="r in allRuntimes" :key="r" :value="r">{{ r }}</option>
        </select>
        <select v-model="filterStatus" class="filter-select" id="filter-status">
          <option value="">All Status</option>
          <option value="instrumented">Instrumented</option>
          <option value="not_instrumented">Not Instrumented</option>
        </select>
        <span class="result-count font-mono">{{ filteredFunctions.length }} / {{ functions.length }} functions</span>
      </div>
    </div>

    <!-- Actions Row -->
    <div class="actions-row">
      <span class="select-hint">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 11l3 3L22 4"/><path d="M21 12v7a2 2 0 01-2 2H5a2 2 0 01-2-2V5a2 2 0 012-2h11"/></svg>
        {{ selectedArns.length > 0 ? `${selectedArns.length} function${selectedArns.length !== 1 ? 's' : ''} selected` : 'Select functions to instrument' }}
      </span>
      <div class="action-buttons">
        <button
          class="action-btn action-btn--instrument"
          :disabled="selectedArns.length === 0"
          @click="showInstrumentModal = true"
          id="btn-instrument"
        >
          <span class="action-btn__icon">
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2"><polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/></svg>
          </span>
          <span class="action-btn__label">Instrument</span>
          <span v-if="selectedArns.length > 0" class="action-btn__count">{{ selectedArns.length }}</span>
        </button>
        <button
          class="action-btn action-btn--uninstrument"
          :disabled="selectedInstrumentedArns.length === 0"
          @click="doUninstrument"
          id="btn-uninstrument"
        >
          <span class="action-btn__icon">
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2"><polyline points="3 6 5 6 21 6"/><path d="M19 6l-1 14a2 2 0 01-2 2H8a2 2 0 01-2-2L5 6"/><path d="M10 11v6M14 11v6"/></svg>
          </span>
          <span class="action-btn__label">Uninstrument</span>
          <span v-if="selectedInstrumentedArns.length > 0" class="action-btn__count action-btn__count--danger">{{ selectedInstrumentedArns.length }}</span>
        </button>
      </div>
    </div>

    <!-- Function Table -->
    <div class="table-container">
      <table class="table">
        <thead>
          <tr>
            <th class="th-check">
              <input
                type="checkbox"
                class="checkbox"
                id="select-all"
                :checked="allSelected"
                @change="toggleAll"
              />
            </th>
            <th>FUNCTION NAME</th>
            <th>RUNTIME</th>
            <th>REGION</th>
            <th>STATUS</th>
            <th>MODE</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="fn in filteredFunctions"
            :key="fn.arn"
            :class="{ selected: selectedArns.includes(fn.arn) }"
          >
            <td class="td-check">
              <input
                type="checkbox"
                class="checkbox"
                :id="'checkbox-' + fn.name"
                :checked="selectedArns.includes(fn.arn)"
                :disabled="isOrchestrator(fn.name) || isBlocked"
                @change="toggleSelect(fn.arn)"
              />
            </td>
            <td class="td-name">
              <div class="fn-name font-mono">{{ fn.name }}</div>
              <div class="fn-arn text-xs text-secondary font-mono">{{ fn.arn }}</div>
            </td>
            <td>
              <span class="runtime-badge" :class="runtimeClass(fn.runtime)">{{ fn.runtime }}</span>
            </td>
            <td class="font-mono text-sm text-secondary">{{ fn.region }}</td>
            <td>
              <span class="status-indicator" :class="statusClass(fn)">
                <span class="status-dot"></span>
                {{ statusLabel(fn) }}
              </span>
            </td>
            <td>
              <span v-if="fn.mode !== 'none'" class="mode-badge" :class="'mode-' + fn.mode">
                {{ modeLabel(fn.mode) }}
              </span>
              <span v-else class="text-xs text-secondary">—</span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- ── Instrument Modal ── -->
    <ThemedModal
      v-model="showInstrumentModal"
      :title="'Instrument ' + selectedArns.length + ' function' + (selectedArns.length > 1 ? 's' : '')"
      subtitle="Choose how to instrument the selected Lambda functions."
      size="md"
    >
      <template #header-icon>
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/></svg>
      </template>

      <!-- Body content -->
      <div class="modal-section">
        <h3 class="modal-section-label">Method</h3>
        <div class="option-group">
          <label class="option-card" :class="{ active: instrumentMethod === 'layer' }">
            <input type="radio" name="method" value="layer" v-model="instrumentMethod" />
            <div class="option-icon option-icon--green">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="12 2 2 7 12 12 22 7 12 2"/><polyline points="2 17 12 22 22 17"/><polyline points="2 12 12 17 22 12"/></svg>
            </div>
            <div class="option-text">
              <strong>Layer-Based</strong>
              <p>Adds NR layer + extension sidecar. Zero code changes.</p>
            </div>
            <div class="option-check"><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="20 6 9 17 4 12"/></svg></div>
          </label>
          <label class="option-card" :class="{ active: instrumentMethod === 'log_ingestion' }">
            <input type="radio" name="method" value="log_ingestion" v-model="instrumentMethod" />
            <div class="option-icon option-icon--blue">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/></svg>
            </div>
            <div class="option-text">
              <strong>Log Ingestion</strong>
              <p>Deploys a Lambda to forward CloudWatch logs to NR.</p>
            </div>
            <div class="option-check"><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="20 6 9 17 4 12"/></svg></div>
          </label>
        </div>
      </div>

      <div class="modal-section" v-if="instrumentMethod === 'layer'">
        <h3 class="modal-section-label">Mode</h3>
        <div class="option-group">
          <label class="option-card" :class="{ active: instrumentMode === 'serverless' }">
            <input type="radio" name="mode" value="serverless" v-model="instrumentMode" />
            <div class="option-icon option-icon--green">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/></svg>
            </div>
            <div class="option-text">
              <strong>Serverless</strong>
              <span class="option-tag option-tag--green">Recommended</span>
              <p>Extension-based, lightweight telemetry with minimal overhead.</p>
            </div>
            <div class="option-check"><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="20 6 9 17 4 12"/></svg></div>
          </label>
          <label class="option-card" :class="{ active: instrumentMode === 'apm' }">
            <input type="radio" name="mode" value="apm" v-model="instrumentMode" />
            <div class="option-icon option-icon--blue">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
            </div>
            <div class="option-text">
              <strong>APM</strong>
              <p>Full distributed tracing with APM features and deeper insights.</p>
            </div>
            <div class="option-check"><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="20 6 9 17 4 12"/></svg></div>
          </label>
        </div>
      </div>

      <!-- Footer -->
      <template #footer>
        <button class="btn btn-secondary" @click="showInstrumentModal = false">Cancel</button>
        <button class="btn btn-primary btn-instrument-confirm" @click="doInstrument" id="btn-confirm-instrument">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/></svg>
          Instrument {{ selectedArns.length }} Function{{ selectedArns.length > 1 ? 's' : '' }}
        </button>
      </template>
    </ThemedModal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import ThemedModal from '@/components/ThemedModal.vue'
import type { LambdaFunction, InstrumentMethod, InstrumentationMode } from '@/types'
import { api } from '@/services/api'

const functions = ref<LambdaFunction[]>([])
const search = ref('')
const filterStatus = ref('')
const filterRegion = ref('')
const filterRuntime = ref('')
const selectedArns = ref<string[]>([])
const showInstrumentModal = ref(false)
const instrumentMethod = ref<InstrumentMethod>('layer')
const instrumentMode = ref<'serverless' | 'apm'>('serverless')

const isLoading = ref(false)
const isOperating = ref(false)
const errorMessage = ref('')

async function fetchFunctions() {
  isLoading.value = true
  errorMessage.value = ''
  try {
    functions.value = await api.listFunctions()
  } catch (err) {
    errorMessage.value = err instanceof Error ? err.message : 'Failed to retrieve functions'
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchFunctions()
})

const allRegions = computed(() => {
  const regions = new Set(functions.value.map(f => f.region))
  return Array.from(regions).sort()
})

const allRuntimes = computed(() => {
  const runtimeSet = new Set(functions.value.map(f => f.runtime))
  return Array.from(runtimeSet).sort()
})
const filteredFunctions = computed(() => {
  return functions.value.filter(fn => {
    const q = search.value.toLowerCase()
    const matchesSearch = !q
      || fn.name.toLowerCase().includes(q)
      || fn.arn.toLowerCase().includes(q)

    const matchesStatus = !filterStatus.value || fn.status === filterStatus.value
    const matchesRegion = !filterRegion.value || fn.region === filterRegion.value
    const matchesRuntime = !filterRuntime.value || fn.runtime === filterRuntime.value

    return matchesSearch && matchesStatus && matchesRegion && matchesRuntime
  })
})

const instrumentedCount = computed(() => functions.value.filter(f => f.status === 'instrumented').length)
const notInstrumentedCount = computed(() => functions.value.filter(f => f.status !== 'instrumented').length)
const instrumentedPercent = computed(() => {
  if (!functions.value.length) return 0
  return Math.round((instrumentedCount.value / functions.value.length) * 100)
})

const selectedInstrumentedArns = computed(() =>
  selectedArns.value.filter(arn => {
    const fn = functions.value.find(f => f.arn === arn)
    return fn && fn.status === 'instrumented'
  })
)

function isOrchestrator(name: string): boolean {
  return name.toLowerCase().includes('orchestrator')
}

const allSelected = computed(() => {
  const selectable = filteredFunctions.value.filter(f => !isOrchestrator(f.name))
  return selectable.length > 0 && selectable.every(f => selectedArns.value.includes(f.arn))
})

// Prevent selections or checks during loading
const isBlocked = computed(() => isLoading.value || isOperating.value)

function toggleAll() {
  if (isBlocked.value) return
  if (allSelected.value) {
    selectedArns.value = []
  } else {
    selectedArns.value = filteredFunctions.value
      .filter(f => !isOrchestrator(f.name))
      .map(f => f.arn)
  }
}

function toggleSelect(arn: string) {
  if (isBlocked.value) return
  const fn = functions.value.find(f => f.arn === arn)
  if (fn && isOrchestrator(fn.name)) return
  const idx = selectedArns.value.indexOf(arn)
  if (idx >= 0) selectedArns.value.splice(idx, 1)
  else selectedArns.value.push(arn)
}

function runtimeClass(runtime: string): string {
  if (runtime.startsWith('nodejs')) return 'rt-node'
  if (runtime.startsWith('python')) return 'rt-python'
  if (runtime.startsWith('java')) return 'rt-java'
  if (runtime.startsWith('dotnet')) return 'rt-dotnet'
  if (runtime.startsWith('ruby')) return 'rt-ruby'
  if (runtime.startsWith('go')) return 'rt-go'
  return 'rt-other'
}

// Map styles for current execution status
function statusClass(fn: LambdaFunction): string {
  if (isOrchestrator(fn.name)) return 'status-protected'
  return fn.status === 'instrumented' ? 'status-active' : 'status-none'
}

function statusLabel(fn: LambdaFunction): string {
  if (isOrchestrator(fn.name)) return 'Protected'
  return fn.status === 'instrumented' ? 'Instrumented' : 'Not Instrumented'
}

function modeLabel(mode: InstrumentationMode): string {
  switch (mode) {
    case 'serverless': return 'Serverless'
    case 'apm': return 'APM'
    case 'log_ingestion': return 'Log Ingestion'
    default: return '—'
  }
}

async function doInstrument() {
  if (selectedArns.value.length === 0) return
  
  isOperating.value = true
  errorMessage.value = ''
  showInstrumentModal.value = false
  
  try {
    const payload = {
      functionArns: selectedArns.value,
      method: instrumentMethod.value,
      mode: instrumentMethod.value === 'layer' ? instrumentMode.value : 'serverless'
    }
    
    await api.installFunctions(payload)
    selectedArns.value = []
    
    // Invalidate state and re-sync
    await fetchFunctions()
  } catch (err) {
    errorMessage.value = err instanceof Error ? err.message : 'Bulk instrumentation failed'
  } finally {
    isOperating.value = false
  }
}

async function doUninstrument() {
  if (selectedInstrumentedArns.value.length === 0) return
  
  isOperating.value = true
  errorMessage.value = ''
  
  try {
    const payload = {
      functionArns: selectedInstrumentedArns.value
    }
    
    await api.uninstallFunctions(payload)
    selectedArns.value = []
    
    // Invalidate state and re-sync
    await fetchFunctions()
  } catch (err) {
    errorMessage.value = err instanceof Error ? err.message : 'Bulk uninstrumentation failed'
  } finally {
    isOperating.value = false
  }
}
</script>

<style scoped>
/* ── Summary Cards ── */
.summary-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--space-4);
  margin-bottom: var(--space-6);
}

.summary-card {
  display: flex;
  align-items: flex-start;
  gap: var(--space-4);
  padding: var(--space-5);
  border-radius: var(--radius-lg);
  border: 1px solid var(--border-primary);
  background: var(--bg-secondary);
  position: relative;
  overflow: hidden;
  transition: border-color var(--transition-base), transform var(--transition-fast);
}

.summary-card::after {
  content: '';
  position: absolute;
  top: 0; left: 0; right: 0;
  height: 2px;
  border-radius: var(--radius-lg) var(--radius-lg) 0 0;
}

.summary-card--total::after { background: var(--info); }
.summary-card--instrumented::after { background: var(--success); }
.summary-card--not::after { background: var(--text-tertiary); }

.summary-card:hover {
  border-color: var(--border-hover);
  transform: translateY(-1px);
}

.summary-text {
  flex: 1;
  min-width: 0;
}
.summary-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  border-radius: var(--radius-md);
  flex-shrink: 0;
}

.summary-icon--total { background: var(--info-dim); color: var(--info); }
.summary-icon--instrumented { background: var(--success-dim); color: var(--success); }
.summary-icon--not { background: rgba(139, 148, 158, 0.15); color: var(--text-secondary); }

.summary-num {
  display: block;
  font-size: var(--font-size-2xl);
  font-weight: 700;
  font-family: var(--font-mono);
  line-height: 1.1;
}

.summary-card--total .summary-num { color: var(--info); }
.summary-card--instrumented .summary-num { color: var(--success); }

.summary-label {
  font-size: var(--font-size-xs);
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

/* ── Toolbar ── */
.toolbar {
  display: flex;
  align-items: center;
  gap: var(--space-4);
  margin-bottom: var(--space-4);
  flex-wrap: wrap;
}

.search-box {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  background: var(--bg-primary);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-md);
  padding: var(--space-2) var(--space-3);
  min-width: 280px;
  transition: border-color var(--transition-fast);
}

.search-box:focus-within {
  border-color: var(--accent);
  box-shadow: 0 0 0 3px var(--accent-dim);
}

.search-icon {
  color: var(--text-tertiary);
  flex-shrink: 0;
}

.search-input {
  background: none;
  border: none;
  outline: none;
  color: var(--text-primary);
  font-family: var(--font-mono);
  font-size: var(--font-size-sm);
  width: 100%;
}

.search-input::placeholder {
  color: var(--text-tertiary);
}

.search-clear {
  background: none;
  border: none;
  color: var(--text-tertiary);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2px;
  border-radius: 4px;
  flex-shrink: 0;
  transition: color var(--transition-fast), background var(--transition-fast);
}

.search-clear:hover {
  color: var(--text-primary);
  background: var(--bg-hover);
}

.filters {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  flex-wrap: wrap;
}

.filter-icon {
  color: var(--text-tertiary);
}

.filter-label {
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
  font-weight: 500;
}

.filter-select {
  background: var(--bg-tertiary);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-md);
  color: var(--text-primary);
  font-family: var(--font-family);
  font-size: var(--font-size-sm);
  padding: var(--space-1) var(--space-3);
  cursor: pointer;
  outline: none;
}

.filter-select:focus {
  border-color: var(--accent);
}

.result-count {
  font-size: var(--font-size-xs);
  color: var(--text-tertiary);
  margin-left: auto;
}

/* ── Actions Row ── */
.actions-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--space-4);
}

.select-hint {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: var(--text-tertiary);
  font-size: var(--font-size-sm);
}

.action-buttons {
  display: flex;
  gap: var(--space-2);
}

/* ── Action Buttons ── */
.action-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  border-radius: 10px;
  font-size: var(--font-size-sm);
  font-weight: 600;
  cursor: pointer;
  border: none;
  transition: all 0.18s ease;
  white-space: nowrap;
}

.action-btn:disabled {
  opacity: 0.38;
  cursor: not-allowed;
  pointer-events: none;
}

.action-btn__icon {
  display: flex;
  align-items: center;
  justify-content: center;
}

.action-btn__label {
  line-height: 1;
}

.action-btn__count {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 20px;
  height: 20px;
  padding: 0 6px;
  border-radius: 999px;
  font-size: 0.7rem;
  font-weight: 700;
  background: rgba(0, 0, 0, 0.25);
  color: inherit;
}

.action-btn__count--danger {
  background: rgba(255, 255, 255, 0.15);
}

/* Instrument — green gradient */
.action-btn--instrument {
  background: linear-gradient(135deg, #1ce783 0%, #0fd470 100%);
  color: #052b16;
  box-shadow: 0 2px 12px rgba(28, 231, 131, 0.3);
}

.action-btn--instrument:not(:disabled):hover {
  background: linear-gradient(135deg, #24f090 0%, #14e07a 100%);
  box-shadow: 0 4px 20px rgba(28, 231, 131, 0.45);
  transform: translateY(-1px);
}

.action-btn--instrument:not(:disabled):active {
  transform: translateY(0);
  box-shadow: 0 1px 6px rgba(28, 231, 131, 0.2);
}

/* Uninstrument — ghost red */
.action-btn--uninstrument {
  background: rgba(239, 68, 68, 0.08);
  color: #fc8181;
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.action-btn--uninstrument:not(:disabled):hover {
  background: rgba(239, 68, 68, 0.15);
  border-color: rgba(239, 68, 68, 0.6);
  color: #fca5a5;
  transform: translateY(-1px);
}

.action-btn--uninstrument:not(:disabled):active {
  transform: translateY(0);
}

/* ── Table ── */
.td-name {
  min-width: 260px;
}

.fn-name {
  font-weight: 600;
  font-size: var(--font-size-sm);
  margin-bottom: 2px;
}

.fn-arn {
  margin-bottom: 4px;
  opacity: 0.7;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 380px;
}

.summary-progress {
  width: 100%;
  height: 3px;
  background: var(--bg-tertiary);
  border-radius: var(--radius-full);
  margin-top: var(--space-2);
  overflow: hidden;
}

.summary-progress-bar {
  height: 100%;
  background: var(--success);
  border-radius: var(--radius-full);
  box-shadow: 0 0 6px rgba(28, 231, 131, 0.4);
  transition: width 0.6s ease;
}

.summary-pct {
  font-size: 0.6rem;
  font-family: var(--font-mono);
  color: var(--success);
  margin-top: 2px;
  display: block;
}


.runtime-badge {
  display: inline-block;
  padding: 3px 10px;
  border-radius: var(--radius-full);
  font-family: var(--font-mono);
  font-size: var(--font-size-xs);
  font-weight: 500;
  white-space: nowrap;
}

.rt-node { background: #1a472a; color: #4ade80; }
.rt-python { background: #2d1a50; color: #a78bfa; }
.rt-java { background: #3b1a1a; color: #fb923c; }
.rt-dotnet { background: #1a2a4b; color: #60a5fa; }
.rt-ruby { background: #3b1a2a; color: #fb7185; }
.rt-go { background: #1a3040; color: #22d3ee; }
.rt-other { background: var(--bg-tertiary); color: var(--text-secondary); }

/* ── Status Indicator ── */
.status-indicator {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  font-size: var(--font-size-sm);
  white-space: nowrap;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.status-active .status-dot { background: var(--accent); }
.status-none .status-dot { border: 1.5px solid var(--text-tertiary); background: transparent; }
.status-protected .status-dot { background: var(--warning); }

.status-active { color: var(--accent); }
.status-none { color: var(--text-tertiary); }
.status-protected { color: var(--warning); }

/* ── Mode Badge ── */
.mode-badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: var(--radius-full);
  font-size: var(--font-size-xs);
  font-weight: 500;
}

.mode-serverless { background: var(--accent-dim); color: var(--accent); }
.mode-apm { background: var(--info-dim); color: var(--info); }
.mode-log_ingestion { background: var(--warning-dim); color: var(--warning); }

.th-check, .td-check { width: 40px; text-align: center; }

/* ── Modal body sections (inside ThemedModal) ── */
.modal-section {
  margin-bottom: var(--space-5);
}

.modal-section:last-child {
  margin-bottom: 0;
}

.modal-section-label {
  font-size: 0.65rem;
  color: var(--text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.08em;
  font-weight: 700;
  margin-bottom: var(--space-3);
}

.option-group {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

/* Option card */
.option-card {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px 16px;
  border: 1px solid var(--border-primary);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.18s ease;
  background: var(--bg-tertiary);
  position: relative;
}

.option-card input[type="radio"] {
  display: none;
}

.option-card:hover {
  border-color: var(--border-hover);
  background: var(--bg-hover);
}

.option-card.active {
  border-color: var(--accent);
  background: rgba(28, 231, 131, 0.06);
  box-shadow: 0 0 0 1px rgba(28, 231, 131, 0.1);
}

.option-icon {
  width: 36px;
  height: 36px;
  border-radius: 9px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.option-icon--green {
  background: rgba(28, 231, 131, 0.12);
  color: var(--accent);
}

.option-icon--blue {
  background: rgba(59, 130, 246, 0.12);
  color: #60a5fa;
}

.option-text {
  flex: 1;
  min-width: 0;
}

.option-text strong {
  display: inline-block;
  font-size: var(--font-size-sm);
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 2px;
  margin-right: var(--space-2);
}

.option-text p {
  font-size: var(--font-size-xs);
  color: var(--text-secondary);
  margin: 0;
  line-height: 1.5;
}

.option-tag {
  display: inline-block;
  font-size: 0.6rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  padding: 2px 7px;
  border-radius: 999px;
  vertical-align: middle;
}

.option-tag--green {
  background: rgba(28, 231, 131, 0.15);
  color: var(--accent);
  border: 1px solid rgba(28, 231, 131, 0.3);
}

.option-check {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  border: 1.5px solid var(--border-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  color: transparent;
  flex-shrink: 0;
  transition: all 0.15s ease;
}

.option-card.active .option-check {
  background: var(--accent);
  border-color: var(--accent);
  color: #000;
}

/* ── Instrument confirm button ── */
.btn-instrument-confirm {
  background: linear-gradient(135deg, #1ce783 0%, #0fd470 100%);
  color: #052b16;
  border-color: transparent;
  font-weight: 600;
  box-shadow: 0 2px 12px rgba(28, 231, 131, 0.25);
}

.btn-instrument-confirm:hover:not(:disabled) {
  background: linear-gradient(135deg, #24f090 0%, #14e07a 100%);
  box-shadow: 0 4px 20px rgba(28, 231, 131, 0.4);
  transform: translateY(-1px);
}
</style>
