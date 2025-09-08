<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useShotsStore } from '@/stores/shots'
import { Scatter } from 'vue-chartjs'
import { Chart, LinearScale, PointElement, Tooltip, Legend, Title } from 'chart.js'

Chart.register(LinearScale, PointElement, Tooltip, Legend, Title)

const props = defineProps<{ limit?: number }>()

const store = useShotsStore()
onMounted(() => {
	if (store.shots.length === 0) store.fetchShots(props.limit ?? 200)
})

const data = computed(() => {
	const rows = props.limit ? store.shots.slice(0, props.limit) : store.shots
	return {
		datasets: [
			{
				label: 'Grind vs Peak Pressure',
				data: rows.map(s => ({ x: s.grind_size_actual, y: s.peak_pressure_bar })),
				pointRadius: 4,
				backgroundColor: rows.map(s =>
					s.last_status === 'ok' ? '#22c55e' :
					s.last_status === 'warning' ? '#f59e0b' : '#ef4444'
				),
			},
		],
	}
})

const options = {
	responsive: true,
	maintainAspectRatio: false,
	plugins: {
		legend: { display: false },
		tooltip: {
			callbacks: {
				label: (ctx: any) => `grind ${ctx.parsed.x}, pressure ${ctx.parsed.y} bar`,
			},
		},
	},
	scales: {
		x: { title: { display: true, text: 'Grind Size (actual)' } },
		y: { ticks: { stepSize: 0.5 },  title: { display: true, text: 'Peak Pressure (bar)' } },
	},
}
</script>

<template>
	<div>
		<Scatter class="h-96" :data="data" :options="options" />
	</div>
</template>
