<template>
  <div class="auth-page">
    <div class="auth-bg">
      <div class="auth-orb auth-orb--1"></div>
      <div class="auth-orb auth-orb--2"></div>
    </div>

    <div class="auth-card animate-slide-up">
      <router-link to="/" class="auth-back" id="auth-back-home">← Back to Home</router-link>

      <div class="auth-header">
        <div class="auth-logo">
          <img src="/logo.png" alt="Lambda Monitor" class="auth-logo-img" />
        </div>
        <h1 class="auth-title">{{ isRegisterMode ? 'Create Account' : 'Welcome Back' }}</h1>
        <p class="auth-subtitle">
          {{ isRegisterMode ? 'Sign up to configure the Instrumentation Platform' : 'Sign in to access the Instrumentation Platform' }}
        </p>
      </div>

      <!-- Registration / Login Form -->
      <form class="auth-form" @submit.prevent="handleSubmit">
        <div class="form-group">
          <label class="form-label" for="login-email">Email</label>
          <input
            v-model="email"
            type="email"
            class="form-input"
            id="login-email"
            placeholder="you@company.com"
            autocomplete="email"
            required
            :disabled="isLoading"
          />
        </div>

        <div class="form-group">
          <label class="form-label" for="login-password">Password</label>
          <input
            v-model="password"
            type="password"
            class="form-input"
            id="login-password"
            placeholder="••••••••"
            autocomplete="current-password"
            required
            :disabled="isLoading"
          />
        </div>

        <!-- Confirm Password for Registration -->
        <div v-if="isRegisterMode" class="form-group">
          <label class="form-label" for="register-confirm-password">Confirm Password</label>
          <input
            v-model="confirmPassword"
            type="password"
            class="form-input"
            id="register-confirm-password"
            placeholder="••••••••"
            autocomplete="new-password"
            required
            :disabled="isLoading"
          />
        </div>

        <div class="form-row" v-if="!isRegisterMode">
          <label class="checkbox-label">
            <input type="checkbox" class="checkbox" v-model="remember" :disabled="isLoading" />
            <span class="text-sm">Remember me</span>
          </label>
          <a href="#" class="forgot-link text-sm">Forgot password?</a>
        </div>

        <!-- Submit Button -->
        <button
          type="submit"
          class="btn btn-primary btn-lg auth-submit"
          :disabled="isLoading"
          id="btn-login"
        >
          <span v-if="isLoading" class="spinner"></span>
          {{ isLoading ? (isRegisterMode ? 'Creating Account...' : 'Signing in...') : (isRegisterMode ? 'Register' : 'Sign In') }}
        </button>

        <!-- Flash Status / Error Notifications -->
        <div v-if="error" class="alert alert-danger mt-4 animate-fade-in">
          {{ error }}
        </div>
        <div v-if="successMsg" class="alert alert-success mt-4 animate-fade-in">
          {{ successMsg }}
        </div>
      </form>

      <div class="auth-divider">
        <span>or</span>
      </div>

      <!-- Mode Toggle Options -->
      <div class="auth-alt">
        <button class="btn btn-secondary auth-sso-btn" id="btn-toggle-mode" @click="toggleMode" :disabled="isLoading">
          {{ isRegisterMode ? 'Already have an account? Sign In' : "Don't have an account? Sign Up" }}
        </button>
      </div>

      <p class="auth-footer" v-if="!isRegisterMode">
        Need assistance?
        <a href="#" class="auth-link" @click.prevent>Contact Support</a>
      </p>
    </div>

    <p class="auth-secure-note text-xs">
      🔒 Secured with AES-256-GCM encryption at rest
    </p>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useConfig } from '@/composables/useConfig'
import { api } from '@/services/api'

const router = useRouter()
const { setAuthenticated, saveToSession, loadPreferencesFromDB, isConnected } = useConfig()

const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const remember = ref(false)
const isLoading = ref(false)
const isRegisterMode = ref(false)
const error = ref('')
const successMsg = ref('')

function toggleMode() {
  isRegisterMode.value = !isRegisterMode.value
  error.value = ''
  successMsg.value = ''
  password.value = ''
  confirmPassword.value = ''
}

async function handleSubmit() {
  error.value = ''
  successMsg.value = ''
  isLoading.value = true

  const emailVal = email.value.trim()
  const passVal = password.value

  if (isRegisterMode.value) {
    // ─── REGISTRATION FLOW ───
    if (passVal !== confirmPassword.value) {
      error.value = 'Passwords do not match.'
      isLoading.value = false
      return
    }

    try {
      await api.register(emailVal, passVal)
      successMsg.value = 'Account registered successfully! Logging you in...'
      
      // Auto login after successful signup
      const res = await api.login(emailVal, passVal)
      setAuthenticated(true, emailVal, res.accessToken)
      saveToSession()
      
      setTimeout(() => {
        router.push('/provider')
      }, 1000)
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Registration failed. Please try again.'
    }
  } else {
    // ─── LOGIN FLOW ───
    try {
      const res = await api.login(emailVal, passVal)
      setAuthenticated(true, emailVal, res.accessToken)
      await loadPreferencesFromDB()
      saveToSession()
      if (isConnected.value) {
        router.push('/dashboard')
      } else {
        router.push('/provider')
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Invalid username or password.'
    }
  }

  isLoading.value = false
}
</script>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--space-8);
  position: relative;
  overflow: hidden;
}

.auth-bg {
  position: absolute;
  inset: 0;
  pointer-events: none;
}

.auth-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(120px);
}

.auth-orb--1 {
  width: 400px;
  height: 400px;
  background: var(--accent);
  opacity: 0.05;
  top: 10%;
  right: 20%;
}

.auth-orb--2 {
  width: 300px;
  height: 300px;
  background: #6366f1;
  opacity: 0.04;
  bottom: 15%;
  left: 15%;
}

.auth-card {
  position: relative;
  z-index: 1;
  max-width: 420px;
  width: 100%;
  background: var(--bg-secondary);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-xl);
  padding: var(--space-8);
}

.auth-back {
  display: inline-block;
  font-size: var(--font-size-sm);
  color: var(--text-tertiary);
  text-decoration: none;
  margin-bottom: var(--space-6);
  transition: color var(--transition-fast);
}

.auth-back:hover {
  color: var(--accent);
}

.auth-header {
  text-align: center;
  margin-bottom: var(--space-8);
}

.auth-logo {
  width: 72px;
  height: 72px;
  margin: 0 auto var(--space-5);
  border-radius: 16px;
  overflow: hidden;
  background: var(--bg-primary);
  box-shadow:
    0 0 0 1px rgba(28, 231, 131, 0.2),
    0 0 24px rgba(28, 231, 131, 0.1),
    0 4px 16px rgba(0, 0, 0, 0.4);
}

.auth-logo-img {
  width: 100%;
  height: 100%;
  display: block;
  object-fit: cover;
}

.auth-title {
  font-size: var(--font-size-2xl);
  font-weight: 700;
  letter-spacing: -0.02em;
  margin-bottom: var(--space-2);
}

.auth-subtitle {
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.form-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  cursor: pointer;
}

.forgot-link {
  color: var(--accent);
  text-decoration: none;
  transition: opacity var(--transition-fast);
}

.forgot-link:hover {
  opacity: 0.8;
}

.auth-submit {
  width: 100%;
  margin-top: var(--space-2);
}

.auth-divider {
  display: flex;
  align-items: center;
  gap: var(--space-4);
  margin: var(--space-6) 0;
  color: var(--text-tertiary);
  font-size: var(--font-size-xs);
  text-transform: uppercase;
  letter-spacing: 0.08em;
}

.auth-divider::before,
.auth-divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: var(--border-primary);
}

.auth-sso-btn {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
}

.auth-footer {
  text-align: center;
  margin-top: var(--space-6);
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
}

.auth-link {
  color: var(--accent);
  text-decoration: none;
  font-weight: 500;
}

.auth-link:hover {
  text-decoration: underline;
}

.auth-secure-note {
  position: absolute;
  bottom: var(--space-6);
  color: var(--text-tertiary);
}
</style>
