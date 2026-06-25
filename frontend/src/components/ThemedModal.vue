<template>
  <Teleport to="body">
    <Transition name="themed-modal">
      <div
        v-if="modelValue"
        class="tm-overlay"
        @click.self="handleOverlayClick"
      >
        <div
          class="tm-dialog"
          :class="[sizeClass]"
          role="dialog"
          aria-modal="true"
        >
          <!-- Header -->
          <div class="tm-header">
            <div class="tm-header__left">
              <div v-if="$slots['header-icon']" class="tm-header__icon">
                <slot name="header-icon" />
              </div>
              <div class="tm-header__text">
                <h2 class="tm-header__title">
                  <slot name="title">{{ title }}</slot>
                </h2>
                <p v-if="subtitle || $slots.subtitle" class="tm-header__subtitle">
                  <slot name="subtitle">{{ subtitle }}</slot>
                </p>
              </div>
            </div>
            <button
              class="tm-header__close"
              @click="close"
              aria-label="Close modal"
            >
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                <line x1="18" y1="6" x2="6" y2="18" />
                <line x1="6" y1="6" x2="18" y2="18" />
              </svg>
            </button>
          </div>

          <!-- Body -->
          <div class="tm-body">
            <slot />
          </div>

          <!-- Footer -->
          <div v-if="$slots.footer" class="tm-footer">
            <slot name="footer" />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { watch, onUnmounted } from 'vue'

interface Props {
  modelValue: boolean
  title?: string
  subtitle?: string
  size?: 'sm' | 'md' | 'lg'
  closeOnOverlay?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  title: '',
  subtitle: '',
  size: 'md',
  closeOnOverlay: true,
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
}>()

const sizeClass = `tm-dialog--${props.size}`

function close() {
  emit('update:modelValue', false)
}

function handleOverlayClick() {
  if (props.closeOnOverlay) {
    close()
  }
}

// Lock body scroll when modal is open
function lockScroll() {
  document.body.style.overflow = 'hidden'
  document.body.style.paddingRight = `${window.innerWidth - document.documentElement.clientWidth}px`
}

function unlockScroll() {
  document.body.style.overflow = ''
  document.body.style.paddingRight = ''
}

watch(
  () => props.modelValue,
  (isOpen) => {
    if (isOpen) {
      lockScroll()
    } else {
      unlockScroll()
    }
  },
  { immediate: true }
)

onUnmounted(() => {
  unlockScroll()
})
</script>

<style>
/* ═══════════════════════════════════════════════════════════
   ThemedModal — Global styles (no scoped, since Teleported)
   ═══════════════════════════════════════════════════════════ */

/* Overlay */
.tm-overlay {
  position: fixed;
  inset: 0;
  z-index: 9000;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.70);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  padding: 24px;
}

/* Dialog container */
.tm-dialog {
  position: relative;
  display: flex;
  flex-direction: column;
  width: 100%;
  background: var(--bg-secondary, #161B22);
  border: 1px solid var(--border-primary, #30363D);
  border-radius: 16px;
  box-shadow:
    0 0 0 1px rgba(28, 231, 131, 0.06),
    0 24px 80px rgba(0, 0, 0, 0.65),
    0 8px 32px rgba(0, 0, 0, 0.4);
  overflow: hidden;
  max-height: calc(100vh - 48px);
}

/* Sizes */
.tm-dialog--sm { max-width: 400px; }
.tm-dialog--md { max-width: 520px; }
.tm-dialog--lg { max-width: 680px; }

/* ── Header ── */
.tm-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  padding: 20px 24px 16px;
  border-bottom: 1px solid var(--border-primary, #30363D);
  flex-shrink: 0;
}

.tm-header__left {
  display: flex;
  align-items: flex-start;
  gap: 14px;
  min-width: 0;
}

.tm-header__icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  background: var(--accent-dim, rgba(28, 231, 131, 0.13));
  border: 1px solid var(--accent-muted, rgba(28, 231, 131, 0.27));
  color: var(--accent, #1CE783);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  margin-top: 2px;
}

.tm-header__text {
  min-width: 0;
}

.tm-header__title {
  font-size: 1.1rem;
  font-weight: 700;
  color: var(--text-primary, #F0F6FC);
  margin: 0 0 2px;
  letter-spacing: -0.01em;
  line-height: 1.3;
}

.tm-header__subtitle {
  font-size: 0.75rem;
  color: var(--text-tertiary, #6E7681);
  margin: 0;
  line-height: 1.5;
}

.tm-header__close {
  background: none;
  border: none;
  color: var(--text-tertiary, #6E7681);
  cursor: pointer;
  padding: 6px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
  flex-shrink: 0;
  transition: color 0.15s ease, background 0.15s ease;
}

.tm-header__close:hover {
  color: var(--text-primary, #F0F6FC);
  background: var(--bg-hover, #292E3B);
}

/* ── Body ── */
.tm-body {
  padding: 20px 24px;
  overflow-y: auto;
  flex: 1 1 auto;
  min-height: 0;
  /* Custom scrollbar */
  scrollbar-width: thin;
  scrollbar-color: var(--border-hover, #484F58) transparent;
}

.tm-body::-webkit-scrollbar {
  width: 6px;
}

.tm-body::-webkit-scrollbar-track {
  background: transparent;
}

.tm-body::-webkit-scrollbar-thumb {
  background: var(--border-hover, #484F58);
  border-radius: 10px;
}

.tm-body::-webkit-scrollbar-thumb:hover {
  background: var(--text-tertiary, #6E7681);
}

/* ── Footer ── */
.tm-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 10px;
  padding: 14px 24px;
  border-top: 1px solid var(--border-primary, #30363D);
  background: var(--bg-tertiary, #1C2333);
  flex-shrink: 0;
}

/* ── Transitions ── */
.themed-modal-enter-active {
  transition: opacity 0.22s ease;
}

.themed-modal-leave-active {
  transition: opacity 0.18s ease;
}

.themed-modal-enter-active .tm-dialog {
  transition: transform 0.28s cubic-bezier(0.34, 1.4, 0.64, 1), opacity 0.22s ease;
}

.themed-modal-leave-active .tm-dialog {
  transition: transform 0.18s ease, opacity 0.15s ease;
}

.themed-modal-enter-from {
  opacity: 0;
}

.themed-modal-enter-from .tm-dialog {
  transform: scale(0.92) translateY(16px);
  opacity: 0;
}

.themed-modal-leave-to {
  opacity: 0;
}

.themed-modal-leave-to .tm-dialog {
  transform: scale(0.95) translateY(8px);
  opacity: 0;
}

/* ── Responsive ── */
@media (max-width: 600px) {
  .tm-overlay {
    padding: 12px;
    align-items: flex-end;
  }

  .tm-dialog {
    border-radius: 16px 16px 0 0;
    max-height: calc(100vh - 24px);
  }

  .tm-dialog--sm,
  .tm-dialog--md,
  .tm-dialog--lg {
    max-width: 100%;
  }

  .tm-header {
    padding: 16px 16px 12px;
  }

  .tm-body {
    padding: 16px;
  }

  .tm-footer {
    padding: 12px 16px;
  }
}
</style>
