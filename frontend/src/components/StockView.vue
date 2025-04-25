<template>
    <div class="container mx-auto px-4 py-8">
      <div class="mb-6">
        <h1 class="text-2xl font-bold mb-2">Stock Market Data Visualization</h1>
        <p class="text-gray-600">View historical stock data for any symbol</p>
      </div>
  
      <div class="mb-6">
        <label for="symbol" class="block text-sm font-medium text-gray-700 mb-1">Stock Symbol</label>
        <div class="flex">
          <input
            v-model="symbol"
            type="text"
            id="symbol"
            placeholder="e.g. AAPL, MSFT, GOOGL"
            class="px-4 py-2 border rounded-l focus:outline-none focus:ring-2 focus:ring-blue-500"
            @keyup.enter="updateChart"
          />
          <button
            @click="updateChart"
            class="px-4 py-2 bg-blue-600 text-white rounded-r hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            Load
          </button>
        </div>
      </div>
  
      <StockChart v-if="symbol" :symbol="symbol" />
    </div>
  </template>
  
  <script lang="ts">
  import { defineComponent, ref } from 'vue'
  import StockChart from '@/components/stockChart.vue'
  
  export default defineComponent({
    name: 'StockView',
    components: {
      StockChart
    },
    setup() {
      const symbol = ref('AAPL') // Default symbol
  
      const updateChart = () => {
        // The StockChart component will react to prop changes automatically
      }
  
      return {
        symbol,
        updateChart
      }
    }
  })
  </script>