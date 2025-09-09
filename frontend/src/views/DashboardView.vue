<script setup lang="ts">
import GrindVsPressureScatter from '@/components/GrindVsPressureScatter.vue';
import Introduction from '../components/IntroductionSection.vue'
import ShotsList from '../components/ShotsList.vue'
import BrewTimeTrend from '@/components/BrewTimeTrend.vue';
import { useShotsStore } from '@/stores/shots';
import { onBeforeMount } from 'vue';

const store = useShotsStore()

onBeforeMount(() => store.fetchShots())
</script>

<template>
  <main>
    <Introduction />
    <BrewTimeTrend :limit="200" />
    <GrindVsPressureScatter :limit="50" />
    <ShotsList />
  </main>
</template>
