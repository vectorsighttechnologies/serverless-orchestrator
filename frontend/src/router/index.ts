import { createRouter, createWebHistory } from 'vue-router'
import { useConfig } from '@/composables/useConfig'

const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: '/',
            name: 'landing',
            component: () => import('@/pages/LandingPage.vue'),
        },
        {
            path: '/auth',
            name: 'auth',
            component: () => import('@/pages/AuthLogin.vue'),
        },
        {
            path: '/provider',
            name: 'provider',
            component: () => import('@/pages/ProviderSelection.vue'),
        },
        {
            path: '/login',
            name: 'login',
            component: () => import('@/pages/LoginConfig.vue'),
        },
        {
            path: '/dashboard',
            name: 'dashboard',
            component: () => import('@/pages/Dashboard.vue'),
        },
        {
            path: '/instrument',
            name: 'instrument',
            component: () => import('@/pages/InstrumentPage.vue'),
        },
    ],
    scrollBehavior(_to, _from, savedPosition) {
        if (savedPosition) return savedPosition
        return { top: 0 }
    },
})

router.beforeEach((to) => {
    const { isAuthenticated, provider, isConnected } = useConfig()

    // Public routes
    if (to.name === 'landing' || to.name === 'auth') {
        return
    }

    // Everything after auth requires authentication
    if (!isAuthenticated.value) {
        return { name: 'auth' }
    }

    if (to.name === 'login' && !provider.value) {
        return { name: 'provider' }
    }

    if (to.name === 'dashboard' && !isConnected.value) {
        return { name: 'login' }
    }

    if (to.name === 'instrument' && !isConnected.value) {
        return { name: 'login' }
    }
})

export default router
