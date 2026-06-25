export type Provider = 'newrelic' | 'datadog'

export type InstrumentationStatus = 'instrumented' | 'not_instrumented'
export type InstrumentationMode = 'serverless' | 'apm' | 'log_ingestion' | 'none'
export type InstrumentMethod = 'layer' | 'log_ingestion'
export type IntegrationMethod = 'metric_streams' | 'api_polling'
export type IntegrationStatus = 'active' | 'not_setup' | 'in_progress' | 'error'
export type ConfigSource = 'env_vars' | 'request' | 'none'

export interface RegionConfig {
    id?: string
    name?: string
    region: string
    apiGatewayUrl: string
    apiKey: string
}

export interface NRCredentials {
    licenseKey: string
    accountId: string
    apiKey: string
    region: 'us' | 'eu'
}

export interface HealthResponse {
    status: string
    config: {
        licenseKeyConfigured: boolean
        accountIdConfigured: boolean
        apiKeyConfigured: boolean
        region: string
        source: ConfigSource
    }
}

export interface LambdaFunction {
    name: string
    arn: string
    runtime: string
    region: string
    architecture: string
    status: InstrumentationStatus
    mode: InstrumentationMode
    layerVersion: string | null
    lastModified: string
    memorySize: number
    timeout: number
    codeSize: number
    invocations: number
    errorRate: number
    tags: Record<string, string>
}

export interface InstallRequest {
    functionArns: string[]
    method: InstrumentMethod
    mode: 'serverless' | 'apm'
}

export interface UninstallRequest {
    functionArns: string[]
}

export interface OperationResult {
    arn: string
    functionName: string
    success: boolean
    error?: string
}

export interface IntegrationInfo {
    status: IntegrationStatus
    method?: IntegrationMethod
    stackName?: string
    stackStatus?: string
    lastChecked: string
    error?: string
}

export interface IntegrationSetupRequest {
    method: IntegrationMethod
    includeLogs: boolean
}
