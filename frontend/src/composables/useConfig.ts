import { ref, computed } from 'vue'
import type { Provider, RegionConfig, NRCredentials, ConfigSource } from '@/types'
import { api } from '@/services/api'

const provider = ref<Provider | null>(null)
const regions = ref<RegionConfig[]>([])
const pendingRegions = ref<RegionConfig[]>([
    { region: 'us-east-1', apiGatewayUrl: '', apiKey: '' }
])
const activeRegionIndex = ref(0)
const nrCredentials = ref<NRCredentials>({
    licenseKey: '',
    accountId: '',
    apiKey: '',
    region: 'us',
})
const configSource = ref<ConfigSource>('none')
const isConnected = ref(false)
const isAuthenticated = ref(false)
const userEmail = ref('')
const token = ref('')

const activeRegion = computed(() => regions.value[activeRegionIndex.value] || null)

// We set default API URL to our Go Backend Gateway
const backendUrl = ref('http://localhost:9000/api')

function setProvider(p: Provider) {
    provider.value = p
}

// Set up regions (in this unified backend model, we only have one main backend endpoint, 
// but we keep region lists for layout compatibility)
function addRegion(config: RegionConfig) {
    regions.value.push(config)
}

// Reset regions
function setRegions(configs: RegionConfig[]) {
    regions.value = configs
}

function removeRegion(index: number) {
    regions.value.splice(index, 1)
    if (activeRegionIndex.value >= regions.value.length) {
        activeRegionIndex.value = Math.max(0, regions.value.length - 1)
    }
}

function updateRegion(index: number, config: RegionConfig) {
    regions.value[index] = config
}

// Select active region dropdown
function setActiveRegion(index: number) {
    activeRegionIndex.value = index
}

function setPendingRegions(configs: RegionConfig[]) {
    pendingRegions.value = configs
}

function addPendingRegion() {
    pendingRegions.value.push({ region: 'us-east-1', apiGatewayUrl: '', apiKey: '' })
}

// Remove card
function removePendingRegion(index: number) {
    pendingRegions.value.splice(index, 1)
}

function setNRCredentials(creds: NRCredentials) {
    nrCredentials.value = creds
}

function setConfigSource(source: ConfigSource) {
    configSource.value = source
}

function setConnected(connected: boolean) {
    isConnected.value = connected
}

function setAuthenticated(auth: boolean, email = '', userToken = '') {
    isAuthenticated.value = auth
    userEmail.value = email
    if (userToken) {
        token.value = userToken
        api.setToken(userToken)
    }
}

function logout() {
    provider.value = null
    regions.value = []
    pendingRegions.value = [{ region: 'us-east-1', apiGatewayUrl: '', apiKey: '' }]
    activeRegionIndex.value = 0
    nrCredentials.value = { licenseKey: '', accountId: '', apiKey: '', region: 'us' }
    configSource.value = 'none'
    isConnected.value = false
    isAuthenticated.value = false
    userEmail.value = ''
    token.value = ''
    api.clearToken()
    sessionStorage.clear()
}

function saveToSession() {
    sessionStorage.setItem('lip_provider', provider.value || '')
    sessionStorage.setItem('lip_regions', JSON.stringify(regions.value))
    sessionStorage.setItem('lip_pendingRegions', JSON.stringify(pendingRegions.value))
    sessionStorage.setItem('lip_activeRegion', String(activeRegionIndex.value))
    sessionStorage.setItem('lip_configSource', configSource.value)
    sessionStorage.setItem('lip_connected', String(isConnected.value))
    sessionStorage.setItem('lip_authenticated', String(isAuthenticated.value))
    sessionStorage.setItem('lip_userEmail', userEmail.value)
    sessionStorage.setItem('lip_backendUrl', backendUrl.value)
    if (token.value) {
        sessionStorage.setItem('lip_token', token.value)
    }

    const active = activeRegion.value
    if (active && active.id) {
        sessionStorage.setItem('lip_activeRegion_id', active.id)
    } else {
        sessionStorage.removeItem('lip_activeRegion_id')
    }
}

function loadFromSession() {
    const savedProvider = sessionStorage.getItem('lip_provider')
    if (savedProvider) provider.value = savedProvider as Provider

    const savedRegions = sessionStorage.getItem('lip_regions')
    if (savedRegions) regions.value = JSON.parse(savedRegions)

    const savedPending = sessionStorage.getItem('lip_pendingRegions')
    if (savedPending) pendingRegions.value = JSON.parse(savedPending)

    const savedActiveRegion = sessionStorage.getItem('lip_activeRegion')
    if (savedActiveRegion) activeRegionIndex.value = Number(savedActiveRegion)

    const savedSource = sessionStorage.getItem('lip_configSource')
    if (savedSource) configSource.value = savedSource as ConfigSource

    const savedConnected = sessionStorage.getItem('lip_connected')
    if (savedConnected) isConnected.value = savedConnected === 'true'

    const savedAuth = sessionStorage.getItem('lip_authenticated')
    if (savedAuth) isAuthenticated.value = savedAuth === 'true'

    const savedEmail = sessionStorage.getItem('lip_userEmail')
    if (savedEmail) userEmail.value = savedEmail

    const savedBackendUrl = sessionStorage.getItem('lip_backendUrl')
    if (savedBackendUrl) {
        backendUrl.value = savedBackendUrl
        api.configure(savedBackendUrl)
    }

    const savedToken = sessionStorage.getItem('lip_token')
    if (savedToken) {
        token.value = savedToken
        api.setToken(savedToken)
    }
}

async function loadPreferencesFromDB() {
    try {
        const prefs = await api.getPreferences()
        
        provider.value = (prefs.selectedProvider as Provider) || 'newrelic'
        nrCredentials.value = {
            licenseKey: prefs.hasNrLicenseKey ? '••••••••' : '',
            accountId: prefs.nrAccountId || '',
            apiKey: prefs.hasNrApiKey ? '••••••••' : '',
            region: (prefs.nrRegion as any) || 'us',
        }

        // Load multiple AWS connections
        const conns = await api.getConnections()
        if (conns && conns.length > 0) {
            regions.value = conns
            isConnected.value = true
        } else if (prefs.lambdaApiUrl) {
            // Fallback for legacy DB configurations
            regions.value = [
                {
                    name: 'Default Connection',
                    region: prefs.nrRegion || 'us-east-1',
                    apiGatewayUrl: prefs.lambdaApiUrl,
                    apiKey: prefs.hasLambdaApiKey ? '••••••••' : ''
                }
            ]
            isConnected.value = true
        }

        saveToSession()
    } catch (err) {
        console.error('Failed to load preferences from DB:', err)
    }
}

export function useConfig() {
    return {
        provider,
        regions,
        pendingRegions,
        activeRegionIndex,
        activeRegion,
        nrCredentials,
        configSource,
        isConnected,
        isAuthenticated,
        userEmail,
        token,
        backendUrl,
        setProvider,
        addRegion,
        setRegions,
        removeRegion,
        updateRegion,
        setActiveRegion,
        setPendingRegions,
        addPendingRegion,
        removePendingRegion,
        setNRCredentials,
        setConfigSource,
        setConnected,
        setAuthenticated,
        logout,
        saveToSession,
        loadFromSession,
        loadPreferencesFromDB,
    }
}
