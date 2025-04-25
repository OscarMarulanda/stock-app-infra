<template>
  <div class="min-h-screen bg-gray-50">
    <StockView />
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import StockView from '@/components/StockView.vue'

export default defineComponent({
  name: 'App',
  components: {
    StockView
  }
})
</script>

<style>
/* Your global styles if needed */
</style>
