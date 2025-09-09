<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useShotsStore, type Shot } from '@/stores/shots'

const props = defineProps<{ limit?: number }>()
const store = useShotsStore()

watch(
  () => props.limit,
  (n) => {
    if (n) store.fetchShots(n)
  },
)

type SortKey = keyof Shot
const sortKey = ref<SortKey>('brew_time')
const sortAsc = ref(false)
function setSort(key: SortKey) {
  if (sortKey.value === key) sortAsc.value = !sortAsc.value
  else {
    sortKey.value = key
    sortAsc.value = true
  }
}

function getComparableValue(shot: Shot, key: SortKey): number | string {
  const raw = shot[key] as unknown
  if (typeof raw === 'number') return raw
  if (typeof raw === 'string') {
    if (key === 'brew_time') {
      const ms = Date.parse(raw)
      return Number.isNaN(ms) ? raw : ms
    }
    return raw
  }
  return ''
}

const sortedRows = computed(() => {
  const arr = [...store.shots]
  arr.sort((a, b) => {
    const va = getComparableValue(a, sortKey.value)
    const vb = getComparableValue(b, sortKey.value)
    if (va < vb) return sortAsc.value ? -1 : 1
    if (va > vb) return sortAsc.value ? 1 : -1
    return 0
  })
  return arr
})

function formatDate(iso: string) {
  return new Date(iso).toLocaleString()
}
</script>

<template>
  <section>
    <header class="mb-2">
      <h2 class="m-0 text-lg font-semibold">Recent Shots</h2>
      <small v-if="store.loading">Loading…</small>
      <small v-else-if="store.error" class="text-red-700">{{ store.error }}</small>
    </header>

    <div class="overflow-auto border border-gray-200 rounded-lg">
      <table class="w-full border-collapse text-sm">
        <thead class="bg-gray-100">
          <tr>
            <th @click="setSort('brew_time')" class="p-2 cursor-pointer text-left select-none">
              Brew Time
            </th>
            <th
              @click="setSort('grind_size_actual')"
              class="p-2 cursor-pointer text-left select-none"
            >
              Grind Size
            </th>
            <th
              @click="setSort('brew_time_seconds')"
              class="p-2 cursor-pointer text-left select-none"
            >
              Brew Seconds
            </th>
            <th
              @click="setSort('peak_pressure_bar')"
              class="p-2 cursor-pointer text-left select-none"
            >
              Pressure (bar)
            </th>
            <th @click="setSort('dose_grams')" class="p-2 cursor-pointer text-left select-none">
              Dose (g)
            </th>
            <th @click="setSort('last_status')" class="p-2 cursor-pointer text-left select-none">
              Status
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in sortedRows" :key="row.shot_id" class="border-t border-gray-200">
            <td class="p-2 whitespace-nowrap">{{ formatDate(row.brew_time) }}</td>
            <td class="p-2">{{ row.grind_size_actual }}</td>
            <td class="p-2">{{ row.brew_time_seconds }}</td>
            <td class="p-2">{{ row.peak_pressure_bar }}</td>
            <td class="p-2">{{ row.dose_grams }}</td>
            <td class="p-2">
              <span
                :class="{
                  'text-green-600': row.last_status === 'ok',
                  'text-amber-600': row.last_status === 'warning',
                  'text-red-600': row.last_status !== 'ok' && row.last_status !== 'warning',
                }"
              >
                {{ row.last_status }}
              </span>
            </td>
          </tr>
          <tr v-if="!store.loading && !store.error && sortedRows.length === 0">
            <td colspan="6" class="p-3 text-center text-gray-500">No data</td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>
</template>
