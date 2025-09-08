<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useShotsStore } from '@/stores/shots'
import { Line } from 'vue-chartjs'
import { Chart, LineElement, PointElement, LinearScale, CategoryScale, Tooltip, Legend, Title } from 'chart.js'

Chart.register(LineElement, PointElement, LinearScale, CategoryScale, Tooltip, Legend, Title)

const props = defineProps<{ limit?: number }>()

const store = useShotsStore()
onMounted(() => {
	if (store.shots.length === 0) store.fetchShots(props.limit ?? 100)
})

// Group by YYYY-MM-DD and compute average brew_time_seconds
const daily = computed(() => {
	const acc = new Map<string, { sum: number; count: number }>()
	const rows = props.limit ? store.shots.slice(0, props.limit) : store.shots
	for (const s of rows) {
		const t = new Date(s.brew_time)
		if (Number.isNaN(t.getTime())) continue
		const key = t.toISOString().slice(0, 10)
		const e = acc.get(key) ?? { sum: 0, count: 0 }
		e.sum += s.brew_time_seconds
		e.count++
		acc.set(key, e)
	}
	const dates = Array.from(acc.keys()).sort()
	const points = dates.map(d => +(acc.get(d)!.sum / acc.get(d)!.count).toFixed(2))
	return { dates, points }
})

const data = computed(() => ({
	labels: daily.value.dates,
	datasets: [
		{
			label: 'Avg Brew Time (s)',
			data: daily.value.points,
			borderColor: '#0ea5e9',
			backgroundColor: 'rgba(14,165,233,0.2)',
			tension: 0.25,
			fill: true,
		},
	],
}))

const options = {
	responsive: true,
	maintainAspectRatio: false,
	scales: {
		x: { title: { display: true, text: 'Date' } },
		y: { title: { display: true, text: 'Seconds' } },
	},
}
</script>

<template>
	<div class="h-96">
		<Line :data="data" :options="options" />
	</div>
</template>
