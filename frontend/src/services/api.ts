import type {
    HealthResponse,
    LambdaFunction,
    IntegrationInfo,
    InstallRequest,
    UninstallRequest,
    OperationResult,
    IntegrationSetupRequest,
    RegionConfig,
} from '@/types'

// Set USE_MOCK to false to connect to the Go Backend Gateway.
const USE_MOCK = false
const MOCK_DELAY = 600

import {
    mockHealthResponse,
    mockFunctions,
    mockIntegration,
} from '@/mock/data'

function delay(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms))
}

class ApiClient {
    private baseUrl = 'http://localhost:9000/api'
    private token = ''

    constructor() {
        // Read token from sessionStorage if present
        const savedToken = sessionStorage.getItem('lip_token')
        if (savedToken) {
            this.token = savedToken
        }
    }

    configure(baseUrl: string) {
        this.baseUrl = baseUrl.replace(/\/$/, '')
    }

    setToken(token: string) {
        this.token = token
        sessionStorage.setItem('lip_token', token)
    }

    clearToken() {
        this.token = ''
        sessionStorage.removeItem('lip_token')
    }

    private async request<T>(method: string, path: string, body?: unknown): Promise<T> {
        const headers: Record<string, string> = {
            'Content-Type': 'application/json',
        }

        if (this.token) {
            headers['Authorization'] = `Bearer ${this.token}`
        }

        const activeConnId = sessionStorage.getItem('lip_activeRegion_id')
        if (activeConnId) {
            headers['X-Connection-ID'] = activeConnId
        }

        const response = await fetch(`${this.baseUrl}${path}`, {
            method,
            headers,
            body: body ? JSON.stringify(body) : undefined,
        })

        if (!response.ok) {
            if (response.status === 401) {
                this.clearToken()
                sessionStorage.clear()
                // Redirect to login page cleanly
                window.location.href = '/auth'
                throw new Error("Session expired. Please log in again.")
            }
            const errorBody = await response.text()
            let errMsg = `API Error ${response.status}: ${errorBody}`
            try {
                const parsed = JSON.parse(errorBody)
                if (parsed.message) errMsg = parsed.message
            } catch {
                // Not JSON
            }
            throw new Error(errMsg)
        }

        return response.json()
    }

    // ─── Authentication Endpoints ───

    async register(email: string, password: string): Promise<{ message: string }> {
        return this.request('POST', '/auth/register', { email, password })
    }

    async login(email: string, password: string): Promise<{ accessToken: string; refreshToken: string; email: string }> {
        const res = await this.request<{ accessToken: string; refreshToken: string; email: string }>('POST', '/auth/login', { email, password })
        this.setToken(res.accessToken)
        return res
    }

    // ─── User Preferences Endpoints ───

    async getPreferences(): Promise<{
        selectedProvider: string
        nrAccountId: string
        nrRegion: string
        lambdaApiUrl: string
        hasNrApiKey: boolean
        hasNrLicenseKey: boolean
        hasLambdaApiKey: boolean
    }> {
        return this.request('GET', '/user/preferences')
    }

    async savePreferences(prefs: {
        selectedProvider: string
        nrAccountId: string
        nrApiKey?: string
        nrLicenseKey?: string
        nrRegion: string
        lambdaApiUrl: string
        lambdaApiKey?: string
    }): Promise<{ message: string }> {
        return this.request('PUT', '/user/preferences', prefs)
    }

    // ─── Lambda Proxy Endpoints ───

    async health(): Promise<HealthResponse> {
        if (USE_MOCK) {
            await delay(MOCK_DELAY)
            return { ...mockHealthResponse }
        }
        // In the backend gateway, /health matches /api/health
        return this.request('GET', '/health')
    }

    async listFunctions(): Promise<LambdaFunction[]> {
        if (USE_MOCK) {
            await delay(MOCK_DELAY)
            return [...mockFunctions]
        }
        // Proxies to Lambda Orchestrator /functions via backend
        const res = await this.request<{ functions: LambdaFunction[] }>('GET', '/functions')
        return res.functions || []
    }

    async installFunctions(req: InstallRequest): Promise<OperationResult[]> {
        if (USE_MOCK) {
            await delay(MOCK_DELAY * 2)
            return req.functionArns.map(arn => ({
                arn,
                functionName: arn.split(':').pop() || arn,
                success: true,
            }))
        }
        const res = await this.request<BatchOperationResponse>('POST', '/functions/install', req)
        return res.results
    }

    async uninstallFunctions(req: UninstallRequest): Promise<OperationResult[]> {
        if (USE_MOCK) {
            await delay(MOCK_DELAY * 2)
            return req.functionArns.map(arn => ({
                arn,
                functionName: arn.split(':').pop() || arn,
                success: true,
            }))
        }
        const res = await this.request<BatchOperationResponse>('POST', '/functions/uninstall', req)
        return res.results
    }

    // ─── Integration Endpoints (Stubbed/Proxy compatible) ───

    async integrationStatus(): Promise<IntegrationInfo> {
        if (USE_MOCK) {
            await delay(MOCK_DELAY)
            return { ...mockIntegration }
        }
        return this.request('GET', '/integration/status')
    }

    async setupIntegration(req: IntegrationSetupRequest): Promise<{ status: string; stackId?: string }> {
        if (USE_MOCK) {
            await delay(MOCK_DELAY * 3)
            return { status: 'in_progress', stackId: 'arn:aws:cloudformation:us-east-1:123456789:stack/NewRelicMetricStreams/xxx' }
        }
        return this.request('POST', '/integration/setup', req)
    }

    async removeIntegration(): Promise<{ status: string }> {
        if (USE_MOCK) {
            await delay(MOCK_DELAY * 2)
            return { status: 'deleting' }
        }
        return this.request('POST', '/integration/remove')
    }

    // ─── Connection Endpoints ───

    async getConnections(): Promise<RegionConfig[]> {
        return this.request('GET', '/user/connections')
    }

    async saveConnection(conn: RegionConfig): Promise<{ id: string; message: string }> {
        return this.request('POST', '/user/connections', conn)
    }

    async deleteConnection(id: string): Promise<{ message: string }> {
        return this.request('DELETE', `/user/connections/${id}`)
    }
}

// Temporary type to compile local proxy response helper
interface BatchOperationResponse {
    results: OperationResult[]
}

export const api = new ApiClient()
