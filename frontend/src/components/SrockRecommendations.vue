<!-- StockRecommendations.vue -->
<template>
    <div class="mt-6 p-4 bg-white rounded-lg shadow">
      <h3 class="text-lg font-semibold mb-4">Stock Recommendations</h3>
      <div v-if="loading" class="text-gray-500">Loading recommendations...</div>
      <div v-else-if="recommendations.length === 0" class="text-gray-500">No data available</div>
      <table v-else class="w-full text-sm text-left">
        <thead>
          <tr>
            <th class="px-4 py-2">Symbol</th>
            <th class="px-4 py-2">Momentum</th>
            <th class="px-4 py-2">RSI</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="rec in recommendations" :key="rec.symbol" class="border-t">
            <td class="px-4 py-2 font-mono">{{ rec.symbol }}</td>
            <td class="px-4 py-2">{{ rec.momentum.toFixed(2) }}</td>
            <td class="px-4 py-2">{{ rec.rsi.toFixed(2) }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </template>
  
  <script lang="ts">
  import { defineComponent, onMounted } from 'vue'
  import { useStockStore } from '@/stores/stockStore'
  import { storeToRefs } from 'pinia'
  
  export default defineComponent({
    setup() {
      const store = useStockStore()
      const { stockRecommendations, loadingRecommendations } = storeToRefs(store)

  
      onMounted(async () => {
        await store.fetchStockRecommendations('short')
    console.log('Recommendations data:', store.stockRecommendations)
      })
  
      return {
        recommendations: stockRecommendations,
        loading: loadingRecommendations
      }
    }
  })
  </script>