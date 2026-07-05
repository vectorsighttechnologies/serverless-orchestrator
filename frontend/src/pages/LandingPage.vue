<template>
  <div class="landing">
    <!-- Nav -->
    <nav class="landing-nav">
      <div class="nav-inner container">
        <div class="nav-brand">
          <div class="nav-logo-wrap">
            <img src="/logo.png" alt="Serverless Orchestrator" class="nav-logo-img" />
          </div>
          <div>
            <span class="nav-name">Serverless<span class="nav-name-accent"> Orchestrator</span></span>
          </div>
        </div>
        <div class="nav-actions">
          <router-link to="/auth" class="btn btn-ghost" id="nav-login">Log In</router-link>
          <router-link to="/auth" class="btn btn-primary" id="nav-get-started">Get Started →</router-link>
        </div>
      </div>
    </nav>

    <!-- Hero -->
    <section class="hero">
      <div class="hero-bg">
        <div class="hero-orb hero-orb--1"></div>
        <div class="hero-orb hero-orb--2"></div>
        <div class="hero-orb hero-orb--3"></div>
        <div class="hero-grid-overlay"></div>
      </div>
      <div class="hero-content container animate-slide-up">
        <div class="hero-badge">
          <span class="badge badge-success">
            <span class="hero-badge-dot"></span>
            Now Available
          </span>
        </div>
        <h1 class="hero-title">
          One-Click Observability<br/>
          for <span class="hero-accent">AWS Lambda</span>
        </h1>
        <p class="hero-desc">
          Instrument your serverless fleet in minutes, not hours.
          Layer-based instrumentation, APM mode, Metric Streams,
          and log ingestion — all from a single pane of glass.
        </p>
        <div class="hero-cta">
          <router-link to="/auth" class="btn btn-primary btn-xl" id="hero-cta">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/></svg>
            Start Instrumenting Free
          </router-link>
          <a href="#how-it-works" class="btn btn-secondary btn-xl" id="hero-learn">
            See How It Works
          </a>
        </div>
        <div class="hero-stats">
          <div class="hero-stat">
            <span class="hero-stat__num">12+</span>
            <span class="hero-stat__label">Runtimes Supported</span>
          </div>
          <div class="hero-stat-divider"></div>
          <div class="hero-stat">
            <span class="hero-stat__num">< 2 min</span>
            <span class="hero-stat__label">Metric Stream Latency</span>
          </div>
          <div class="hero-stat-divider"></div>
          <div class="hero-stat">
            <span class="hero-stat__num">100%</span>
            <span class="hero-stat__label">Serverless Native</span>
          </div>
          <div class="hero-stat-divider"></div>
          <div class="hero-stat">
            <span class="hero-stat__num">0</span>
            <span class="hero-stat__label">Code Changes Required</span>
          </div>
        </div>
      </div>

      <!-- Floating terminal preview -->
      <div class="hero-terminal container animate-fade-in">
        <div class="terminal">
          <div class="terminal__header">
            <div class="terminal__dots">
              <span class="dot dot--red"></span>
              <span class="dot dot--yellow"></span>
              <span class="dot dot--green"></span>
            </div>
            <span class="terminal__title">Serverless Orchestrator — Instrumentation</span>
          </div>
          <div class="terminal__body">
            <div class="terminal__line"><span class="t-dim">$</span> <span class="t-cmd">select</span> provider <span class="t-accent">newrelic</span></div>
            <div class="terminal__line"><span class="t-dim">$</span> <span class="t-cmd">connect</span> --region us-east-1 --url https://api.example.execute-api...</div>
            <div class="terminal__line t-success">✓ Connected to orchestrator · 24 functions discovered</div>
            <div class="terminal__line"><span class="t-dim">$</span> <span class="t-cmd">instrument</span> --method layer --mode serverless --all</div>
            <div class="terminal__line t-success">✓ 24/24 functions instrumented with New Relic layer</div>
            <div class="terminal__line t-dim">→ View telemetry in New Relic One</div>
          </div>
        </div>
      </div>
    </section>

    <!-- Features -->
    <section class="features" id="features">
      <div class="container">
        <div class="section-header animate-fade-in">
          <div class="section-badge">Features</div>
          <h2 class="section-title">Everything You Need</h2>
          <p class="section-desc">
            From single-function instrumentation to full account-wide observability
          </p>
        </div>
        <div class="features-grid">
          <div class="feature-card animate-fade-in" v-for="(f, i) in features" :key="i" :class="'feature-card--' + f.color">
            <div class="feature-icon" :class="'feature-icon--' + f.color" v-html="f.svg"></div>
            <h3 class="feature-name">{{ f.title }}</h3>
            <p class="feature-desc">{{ f.desc }}</p>
            <div class="feature-tag" v-if="f.tag">{{ f.tag }}</div>
          </div>
        </div>
      </div>
    </section>

    <!-- How it works -->
    <section class="how-it-works" id="how-it-works">
      <div class="container">
        <div class="section-header animate-fade-in">
          <div class="section-badge">Process</div>
          <h2 class="section-title">How It Works</h2>
          <p class="section-desc">Three steps to full observability</p>
        </div>
        <div class="steps-row">
          <div class="step-card" v-for="(s, i) in steps" :key="i">
            <div class="step-connector" v-if="i < steps.length - 1"></div>
            <div class="step-number">{{ i + 1 }}</div>
            <div class="step-icon">{{ s.icon }}</div>
            <h3 class="step-title">{{ s.title }}</h3>
            <p class="step-desc">{{ s.desc }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Runtimes -->
    <section class="runtimes-section">
      <div class="container">
        <div class="section-header animate-fade-in">
          <div class="section-badge">Compatibility</div>
          <h2 class="section-title">All Major Runtimes</h2>
          <p class="section-desc">Broad runtime support — instrument any Lambda, zero code changes</p>
        </div>
        <div class="runtimes-grid">
          <div class="runtime-chip" v-for="rt in runtimes" :key="rt.name" :class="'rt-' + rt.lang">
            <span class="runtime-chip__icon">{{ rt.icon }}</span>
            <span>{{ rt.name }}</span>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA -->
    <section class="cta-section">
      <div class="container">
        <div class="cta-card animate-fade-in">
          <div class="cta-orb"></div>
          <div class="cta-body">
            <h2 class="cta-title">Ready to instrument your Lambdas?</h2>
            <p class="cta-desc">Start monitoring your serverless functions in under 5 minutes. No credit card required.</p>
            <router-link to="/auth" class="btn btn-primary btn-xl" id="cta-bottom">
              Get Started Now →
            </router-link>
          </div>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="landing-footer">
      <div class="container">
        <div class="footer-inner">
          <div class="footer-brand">
            <div class="footer-logo-wrap">
              <img src="/logo.png" alt="Serverless Orchestrator" class="footer-logo-img" />
            </div>
            <span>Serverless<span style="color: var(--accent)"> Orchestrator</span></span>
          </div>
          <div class="footer-links">
            <router-link to="/auth" class="footer-link">Log In</router-link>
            <router-link to="/auth" class="footer-link">Get Started</router-link>
          </div>
          <p class="footer-copy">
            © 2026 <a href="https://vectorsight.tech" class="footer-link" target="_blank" rel="noopener">VectorSight Technologies</a> · Apache 2.0 License · Open Source
          </p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
const features = [
  {
    svg: `<svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8"><polygon points="12 2 2 7 12 12 22 7 12 2"/><polyline points="2 17 12 22 22 17"/><polyline points="2 12 12 17 22 12"/></svg>`,
    color: 'green',
    title: 'Layer-Based Instrumentation',
    desc: 'Attach New Relic layers to any Lambda with zero code changes. Automatic instrumentation for Node.js, Python, Java, .NET, Ruby, Go and more.',
    tag: 'Zero Code Changes',
  },
  {
    svg: `<svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>`,
    color: 'blue',
    title: 'APM & Serverless Modes',
    desc: 'Choose lightweight serverless telemetry for minimal overhead, or full APM distributed tracing for deep application insights.',
    tag: null,
  },
  {
    svg: `<svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg>`,
    color: 'cyan',
    title: 'Metric Streams',
    desc: 'Real-time CloudWatch metrics via Kinesis Firehose. Under 2 minute latency, no polling. The modern recommended path for Lambda observability.',
    tag: '< 2 min latency',
  },
  {
    svg: `<svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/></svg>`,
    color: 'amber',
    title: 'Log Ingestion',
    desc: 'Forward CloudWatch logs directly to New Relic via a dedicated Lambda forwarder or a managed Firehose delivery stream.',
    tag: null,
  },
  {
    svg: `<svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8"><circle cx="12" cy="12" r="10"/><line x1="2" y1="12" x2="22" y2="12"/><path d="M12 2a15.3 15.3 0 010 20M12 2a15.3 15.3 0 000 20"/></svg>`,
    color: 'purple',
    title: 'Multi-Region Support',
    desc: 'Manage Lambda instrumentation across any number of AWS regions from a single UI. Configure per-region API endpoints with ease.',
    tag: 'All AWS Regions',
  },
  {
    svg: `<svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0110 0v4"/></svg>`,
    color: 'rose',
    title: 'Secure by Design',
    desc: 'API keys stay in your orchestrator. No sensitive credentials traverse the UI. Full zero-trust architecture — your data never leaves your AWS account.',
    tag: 'Zero-Trust',
  },
]

const steps = [
  {
    icon: '🚀',
    title: 'Deploy the Orchestrator',
    desc: 'Launch the Go-based orchestrator Lambda in your AWS account via a simple CloudFormation stack.',
  },
  {
    icon: '⚙️',
    title: 'Connect & Configure',
    desc: 'Point the UI to your API Gateway endpoint and configure your regions and New Relic credentials.',
  },
  {
    icon: '⚡',
    title: 'Instrument & Monitor',
    desc: 'Select your functions, choose a method, and start seeing data in New Relic within minutes.',
  },
]

const runtimes = [
  { name: 'Node.js 20.x', icon: '🟢', lang: 'node' },
  { name: 'Node.js 18.x', icon: '🟢', lang: 'node' },
  { name: 'Python 3.12', icon: '🐍', lang: 'python' },
  { name: 'Python 3.11', icon: '🐍', lang: 'python' },
  { name: 'Python 3.10', icon: '🐍', lang: 'python' },
  { name: 'Java 21', icon: '☕', lang: 'java' },
  { name: 'Java 17', icon: '☕', lang: 'java' },
  { name: '.NET 8', icon: '💜', lang: 'dotnet' },
  { name: 'Ruby 3.3', icon: '💎', lang: 'ruby' },
  { name: 'Go 1.x', icon: '🔵', lang: 'go' },
  { name: 'Custom Runtime', icon: '⚙️', lang: 'other' },
  { name: '+ More', icon: '✨', lang: 'other' },
]
</script>

<style scoped>
.landing {
  min-height: 100vh;
  background: var(--bg-primary);
}

/* ── Nav ── */
.landing-nav {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 200;
  background: rgba(10, 13, 20, 0.85);
  backdrop-filter: blur(20px);
  border-bottom: 1px solid var(--border-primary);
}

.nav-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-4) var(--space-6);
}

.nav-brand {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.nav-logo-wrap {
  width: 34px;
  height: 34px;
  border-radius: 9px;
  overflow: hidden;
  flex-shrink: 0;
  box-shadow: 0 0 0 1px rgba(28, 231, 131, 0.25), 0 0 12px rgba(28, 231, 131, 0.1);
  transition: box-shadow var(--transition-base);
}

.nav-logo-wrap:hover {
  box-shadow: 0 0 0 1px rgba(28, 231, 131, 0.5), 0 0 20px rgba(28, 231, 131, 0.15);
}

.nav-logo-img { width: 34px; height: 34px; display: block; object-fit: cover; }

.nav-name {
  font-weight: 700;
  font-size: var(--font-size-md);
  letter-spacing: -0.02em;
}

.nav-name-accent { color: var(--accent); }

.nav-actions {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

/* ── Hero ── */
.hero {
  position: relative;
  padding: 160px 0 80px;
  overflow: hidden;
}

.hero-bg {
  position: absolute;
  inset: 0;
  pointer-events: none;
}

.hero-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(100px);
}

.hero-orb--1 {
  width: 700px;
  height: 700px;
  background: var(--accent);
  opacity: 0.05;
  top: -250px;
  right: -150px;
}

.hero-orb--2 {
  width: 500px;
  height: 500px;
  background: #3b82f6;
  opacity: 0.04;
  bottom: -150px;
  left: -150px;
}

.hero-orb--3 {
  width: 300px;
  height: 300px;
  background: #a855f7;
  opacity: 0.03;
  top: 40%;
  left: 40%;
  transform: translate(-50%, -50%);
}

.hero-grid-overlay {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(255,255,255,0.015) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255,255,255,0.015) 1px, transparent 1px);
  background-size: 60px 60px;
}

.hero-content {
  position: relative;
  z-index: 1;
  text-align: center;
  max-width: 820px;
  margin: 0 auto;
}

.hero-badge {
  margin-bottom: var(--space-6);
}

.hero-badge-dot {
  display: inline-block;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--accent);
  margin-right: var(--space-1);
  animation: pulse-dot 2s ease infinite;
  vertical-align: middle;
}

@keyframes pulse-dot {
  0%, 100% { opacity: 1; box-shadow: 0 0 5px var(--accent); }
  50% { opacity: 0.7; box-shadow: 0 0 9px var(--accent); }
}

.hero-title {
  font-size: 3.8rem;
  font-weight: 800;
  line-height: 1.1;
  letter-spacing: -0.03em;
  margin-bottom: var(--space-6);
}

.hero-accent {
  background: linear-gradient(135deg, var(--accent) 0%, #3eda99 50%, #22d3ee 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.hero-desc {
  font-size: var(--font-size-lg);
  color: var(--text-secondary);
  line-height: 1.7;
  max-width: 600px;
  margin: 0 auto var(--space-8);
}

.hero-cta {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-4);
  margin-bottom: var(--space-12);
}

.btn-xl {
  padding: 14px 32px;
  font-size: var(--font-size-md);
  font-weight: 600;
  border-radius: var(--radius-lg);
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
}

.hero-stats {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-8);
  padding-top: var(--space-8);
  border-top: 1px solid var(--border-primary);
}

.hero-stat__num {
  display: block;
  font-size: var(--font-size-2xl);
  font-weight: 700;
  color: var(--accent);
}

.hero-stat__label {
  font-size: var(--font-size-xs);
  color: var(--text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

.hero-stat-divider {
  width: 1px;
  height: 40px;
  background: var(--border-primary);
}

/* ── Terminal ── */
.hero-terminal {
  position: relative;
  z-index: 1;
  max-width: 720px;
  margin: var(--space-16) auto 0;
}

.terminal {
  background: #0d1117;
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: var(--radius-xl);
  overflow: hidden;
  box-shadow: 0 20px 60px rgba(0,0,0,0.5), 0 0 0 1px rgba(28, 231, 131, 0.08);
}

.terminal__header {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-3) var(--space-4);
  background: rgba(255,255,255,0.04);
  border-bottom: 1px solid rgba(255,255,255,0.05);
}

.terminal__dots {
  display: flex;
  gap: 6px;
}

.dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.dot--red { background: #ff5f56; }
.dot--yellow { background: #ffbd2e; }
.dot--green { background: #27c93f; }

.terminal__title {
  font-size: var(--font-size-xs);
  color: var(--text-tertiary);
  font-family: var(--font-mono);
}

.terminal__body {
  padding: var(--space-5) var(--space-6);
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.terminal__line {
  font-family: var(--font-mono);
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
  line-height: 1.6;
}

.t-dim { color: var(--text-tertiary); }
.t-cmd { color: #79c0ff; }
.t-accent { color: var(--accent); }
.t-success { color: var(--success); }

/* ── Section Commons ── */
.section-header {
  text-align: center;
  margin-bottom: var(--space-12);
}

.section-badge {
  display: inline-block;
  background: var(--accent-dim);
  color: var(--accent);
  font-size: var(--font-size-xs);
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  padding: 4px 14px;
  border-radius: var(--radius-full);
  margin-bottom: var(--space-3);
  border: 1px solid var(--accent-muted);
}

.section-title {
  font-size: var(--font-size-2xl);
  font-weight: 700;
  margin-bottom: var(--space-3);
}

.section-desc {
  font-size: var(--font-size-md);
  color: var(--text-secondary);
}

/* ── Features ── */
.features {
  padding: 100px 0;
  border-top: 1px solid var(--border-primary);
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--space-5);
}

.feature-card {
  background: var(--bg-secondary);
  border: 1px solid var(--border-primary);
  border-radius: 18px;
  padding: 32px;
  transition: all var(--transition-base);
  position: relative;
  overflow: hidden;
}

.feature-card::before {
  content: '';
  position: absolute;
  top: 0; left: 0; right: 0;
  height: 2px;
  background: linear-gradient(90deg, var(--accent), transparent);
  opacity: 0;
  transition: opacity var(--transition-base);
}

.feature-card:hover {
  border-color: var(--border-hover);
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.feature-card:hover::before {
  opacity: 1;
}

.feature-icon {
  width: 52px;
  height: 52px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 20px;
  flex-shrink: 0;
}

.feature-icon--green  { background: rgba(28,231,131,0.12); color: #1ce783; }
.feature-icon--blue   { background: rgba(59,130,246,0.12);  color: #60a5fa; }
.feature-icon--cyan   { background: rgba(34,211,238,0.12);  color: #22d3ee; }
.feature-icon--amber  { background: rgba(251,191,36,0.12);   color: #fbbf24; }
.feature-icon--purple { background: rgba(168,85,247,0.12);  color: #c084fc; }
.feature-icon--rose   { background: rgba(251,113,133,0.12); color: #fb7185; }

/* top-border accent by color */
.feature-card--green::before  { background: linear-gradient(90deg, #1ce783, transparent); }
.feature-card--blue::before   { background: linear-gradient(90deg, #60a5fa, transparent); }
.feature-card--cyan::before   { background: linear-gradient(90deg, #22d3ee, transparent); }
.feature-card--amber::before  { background: linear-gradient(90deg, #fbbf24, transparent); }
.feature-card--purple::before { background: linear-gradient(90deg, #c084fc, transparent); }
.feature-card--rose::before   { background: linear-gradient(90deg, #fb7185, transparent); }

.feature-name {
  font-size: var(--font-size-md);
  font-weight: 600;
  margin-bottom: 10px;
}

.feature-desc {
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
  line-height: 1.75;
  margin-bottom: 20px;
}

.feature-tag {
  display: inline-block;
  background: color-mix(in srgb, currentColor 12%, transparent);
  color: currentColor;
  font-size: 0.65rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  padding: 3px 10px;
  border-radius: var(--radius-full);
  border: 1px solid color-mix(in srgb, currentColor 30%, transparent);
}

/* ── How it works ── */
.how-it-works {
  padding: 100px 0;
  background: var(--bg-secondary);
  border-top: 1px solid var(--border-primary);
  border-bottom: 1px solid var(--border-primary);
}

.steps-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--space-8);
  max-width: 900px;
  margin: 0 auto;
  position: relative;
}

.step-card {
  text-align: center;
  position: relative;
}

.step-connector {
  position: absolute;
  top: 24px;
  left: calc(50% + 36px);
  right: calc(-50% + 36px);
  height: 1px;
  background: linear-gradient(90deg, var(--accent), var(--border-primary));
  z-index: 0;
}

.step-number {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  background: var(--accent-dim);
  color: var(--accent);
  font-weight: 700;
  font-size: var(--font-size-lg);
  border-radius: 50%;
  margin-bottom: var(--space-4);
  border: 1px solid var(--accent);
  position: relative;
  z-index: 1;
}

.step-icon {
  font-size: 1.5rem;
  margin-bottom: var(--space-3);
}

.step-title {
  font-size: var(--font-size-md);
  font-weight: 600;
  margin-bottom: var(--space-2);
}

.step-desc {
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
  line-height: 1.6;
}

/* ── Runtimes ── */
.runtimes-section {
  padding: 80px 0;
  border-top: 1px solid var(--border-primary);
}

.runtimes-grid {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-3);
  justify-content: center;
}

.runtime-chip {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-4);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-full);
  font-size: var(--font-size-sm);
  font-family: var(--font-mono);
  background: var(--bg-secondary);
  transition: all var(--transition-fast);
  color: var(--text-secondary);
}

.runtime-chip:hover {
  border-color: var(--border-hover);
  color: var(--text-primary);
  transform: translateY(-1px);
}

.runtime-chip__icon { font-size: 1rem; }

/* Runtime colors on hover */
.rt-node:hover { border-color: #4ade80; color: #4ade80; }
.rt-python:hover { border-color: #a78bfa; color: #a78bfa; }
.rt-java:hover { border-color: #fb923c; color: #fb923c; }
.rt-dotnet:hover { border-color: #60a5fa; color: #60a5fa; }
.rt-ruby:hover { border-color: #fb7185; color: #fb7185; }
.rt-go:hover { border-color: #22d3ee; color: #22d3ee; }

/* ── Bottom CTA ── */
.cta-section {
  padding: 100px 0;
}

.cta-card {
  text-align: center;
  background: linear-gradient(135deg, rgba(28, 231, 131, 0.07) 0%, rgba(59, 130, 246, 0.04) 100%);
  border: 1px solid rgba(28, 231, 131, 0.2);
  border-radius: var(--radius-xl);
  padding: var(--space-16);
  position: relative;
  overflow: hidden;
}

.cta-orb {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 500px;
  height: 500px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(28, 231, 131, 0.05) 0%, transparent 70%);
  pointer-events: none;
}

.cta-body { position: relative; z-index: 1; }

.cta-title {
  font-size: var(--font-size-3xl);
  font-weight: 700;
  margin-bottom: var(--space-3);
}

.cta-desc {
  font-size: var(--font-size-md);
  color: var(--text-secondary);
  margin-bottom: var(--space-8);
}

/* ── Footer ── */
.landing-footer {
  border-top: 1px solid var(--border-primary);
  padding: var(--space-8) 0;
}

.footer-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-6);
}

.footer-logo-img { width: 28px; height: 28px; display: block; object-fit: cover; }

.footer-brand {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  font-weight: 600;
  font-size: var(--font-size-sm);
  color: var(--text-primary);
}

.footer-links {
  display: flex;
  gap: var(--space-6);
}

.footer-link {
  font-size: var(--font-size-sm);
  color: var(--text-tertiary);
  text-decoration: none;
  transition: color var(--transition-fast);
}

.footer-link:hover { color: var(--accent); }

.footer-copy {
  font-size: var(--font-size-xs);
  color: var(--text-tertiary);
}

/* ── Responsive ── */
@media (max-width: 900px) {
  .features-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .hero-title {
    font-size: 2.4rem;
  }
  .features-grid,
  .steps-row {
    grid-template-columns: 1fr;
  }
  .hero-stats {
    flex-direction: column;
    gap: var(--space-4);
  }
  .hero-stat-divider {
    width: 40px;
    height: 1px;
  }
  .hero-cta {
    flex-direction: column;
  }
  .footer-inner {
    flex-direction: column;
    text-align: center;
  }
  .step-connector {
    display: none;
  }
  .cta-title {
    font-size: var(--font-size-2xl);
  }
  .hero-terminal {
    display: none;
  }
}
</style>
