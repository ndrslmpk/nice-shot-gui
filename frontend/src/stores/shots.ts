import { ref } from 'vue'
import { defineStore } from 'pinia'
import type { Shot } from '@/models/shot'
import { fetchShots as fetchShotsApi } from '@/services/shotsService'

export type { Shot }

export const useShotsStore = defineStore('shots', () => {
  const shots = ref<Shot[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchShots(limit: number = 100): Promise<void> {
    loading.value = true
    error.value = null
    try {
      const data = await fetchShotsApi(limit)
      shots.value = data
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Unknown error'
      error.value = message
    } finally {
      loading.value = false
    }
  }

  return { shots, loading, error, fetchShots }
})
