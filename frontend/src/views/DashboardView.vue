<script setup lang="ts">
import GrindVsPressureScatter from '@/components/GrindVsPressureScatter.vue';
import Introduction from '../components/IntroductionSection.vue'
import ShotsList from '../components/ShotsList.vue'
</script>

<template>
  <main>
    <Introduction />
    <GrindVsPressureScatter :limit="50" />
    <ShotsList />
  </main>
</template>
