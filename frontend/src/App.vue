<script setup lang="ts">
import {  RouterView } from 'vue-router'
</script>

<template>
  <header>
    <link href="/src/style.css" rel="stylesheet">
  </header>

  <body>
    <div class="wrapper">
      <RouterView />
    </div>
  </body>

</template>

<style scoped>
</style>
