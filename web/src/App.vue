<script setup lang="ts">
import AuthFeedbackProvider from '@/components/auth/AuthFeedbackProvider.vue'
</script>

<template>
  <div id="app">
    <router-view />
    <AuthFeedbackProvider />
  </div>
</template>

<style scoped>
#app {
  width: 100%;
  min-height: 100vh;
}
</style>
