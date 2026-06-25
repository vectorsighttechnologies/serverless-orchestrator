<template>
  <router-view />
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useConfig } from '@/composables/useConfig'

const { loadFromSession, loadPreferencesFromDB, isAuthenticated } = useConfig()

onMounted(async () => {
  loadFromSession()
  if (isAuthenticated.value) {
    await loadPreferencesFromDB()
  }
})
</script>
