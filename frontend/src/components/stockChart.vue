<template>
    <div class="bg-white rounded-lg shadow p-6">
      <div class="flex justify-between items-center mb-4">
        <h2 class="text-xl font-semibold">
          {{ symbol }} Stock Data - {{ timeRangeLabel }}
        </h2>
        <div class="flex space-x-2">
          <button 
            v-for="range in timeRanges" 
            :key="range.value"
            @click="updateRange(range.value)"
            class="px-3 py-1 text-sm rounded"
            :class="{
              'bg-blue-600 text-white': currentRange === range.value,
              'bg-gray-200 text-gray-700': currentRange !== range.value
            }"
          >
            {{ range.label }}
          </button>
          <button 
            @click="refreshData"
            class="px-3 py-1 bg-green-600 text-white text-sm rounded"
            :disabled="loading"
          >
            {{ loading ? 'Refreshing...' : 'Refresh Data' }}
          </button>
        </div>
      </div>
      <div v-if="latestEntry" class="mb-4 p-4 bg-gray-100 rounded text-sm">
  <h3 class="text-lg font-medium mb-2">Latest Stock Info ({{ latestEntry.date }})</h3>
  <div class="grid grid-cols-2 sm:grid-cols-3 gap-2">
    <div><strong>Open:</strong> {{ latestEntry.open.toFixed(2) }}</div>
    <div><strong>High:</strong> {{ latestEntry.high.toFixed(2) }}</div>
    <div><strong>Low:</strong> {{ latestEntry.low.toFixed(2) }}</div>
    <div><strong>Close:</strong> {{ latestEntry.close.toFixed(2) }}</div>
    <div><strong>Volume:</strong> {{ latestEntry.volume.toLocaleString() }}</div>
  </div>
</div>
  
      <div v-if="stockStore.loading" class="loading-indicator">
    <p>Loading data...</p>
    <p v-if="stockStore.error">(Note: {{ stockStore.error }})</p>
  </div>
  
      <div v-else-if="error" class="text-red-500 p-4">
        {{ error }}
      </div>
  
      <div v-else-if="stockStore.stockData.length === 0" class="text-gray-500 p-4">
  No data available for this Symbol, check that the symbol eists.
</div>
  
<div v-else class="h-96">
  <canvas ref="chartCanvas"></canvas>
</div>
    </div>
  </template>
  
  <script lang="ts">
  import { computed, defineComponent, onMounted, ref, watch, onBeforeUnmount, nextTick  } from 'vue'
  import { useStockStore } from '@/stores/stockStore'
  import {
  Chart,
  LineController,
  LineElement,
  PointElement,
  LinearScale,
  TimeScale,
  CategoryScale,
  Title,
  Tooltip,
  Legend,
  Filler
} from 'chart.js'

  Chart.register(
  LineController,
  LineElement,
  PointElement,
  LinearScale,
  TimeScale,
  CategoryScale,
  Title,
  Tooltip,
  Legend,
  Filler
)
  
  export default defineComponent({
    props: {
      symbol: {
        type: String,
        required: true
      }
    },
    setup(props) {
      const initialLoading = ref(true)
      const stockStore = useStockStore()
      const chartCanvas = ref<HTMLCanvasElement | null>(null)
      const chartInstance = ref<Chart | null>(null)
  
      const timeRanges = [
        { value: 'week', label: '1 Week' },
        { value: 'month', label: '1 Month' },
        { value: '6month', label: '6 Months' },
        { value: 'year', label: '1 Year' }
      ]
  
      const currentRange = ref('month')
      const timeRangeLabel = ref(timeRanges.find(r => r.value === currentRange.value)?.label || '')
  
      const destroyChart = () => {
        if (chartInstance.value) {
          chartInstance.value.destroy()
          chartInstance.value = null
        }
      }

      const latestEntry = computed(() => {
  return stockStore.stockData.length > 0
    ? stockStore.stockData[0]
    : null
})
  
      const renderChart = () => {
  if (!chartCanvas.value) {
    console.warn('Canvas not ready yet')
    return
  }

  const ctx = chartCanvas.value.getContext('2d')
  if (!ctx) {
    console.warn('Canvas context not available')
    return
  }

  // ✅ Destroy chart before creating a new one
  if (chartInstance.value) {
    chartInstance.value.destroy()
    chartInstance.value = null
  }

  if (stockStore.stockData.length === 0) return

  const labels = stockStore.stockData.map(d => d.date).reverse()
  const closePrices = stockStore.stockData.map(d => d.close).reverse()

  chartInstance.value = new Chart(ctx, {
    type: 'line',
    data: {
      labels,
      datasets: [{
        label: 'Closing Price',
        data: closePrices,
        borderColor: 'rgb(59, 130, 246)',
        backgroundColor: 'rgba(59, 130, 246, 0.1)',
        fill: true,
        tension: 0.1
      }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      scales: {
        y: {
          beginAtZero: false,
          title: {
            display: true,
            text: 'Price ($)'
          }
        },
        x: {
          title: {
            display: true,
            text: 'Date'
          }
        }
      }
    }
  })
}
  
      const updateRange = (range: string) => {
        currentRange.value = range
        timeRangeLabel.value = timeRanges.find(r => r.value === range)?.label || ''
        stockStore.fetchStockData(props.symbol, range)
      }
  
      const refreshData = async () => {
  await stockStore.refreshStockData()
}
  
      watch(() => stockStore.stockData, async (newData) => {
  if (Array.isArray(newData) && newData.length > 0) {
    await nextTick()  // ⏳ Wait for DOM (including canvas) to update
    if (chartCanvas.value) {
      renderChart()
    } else {
      console.warn('Canvas still not available after nextTick')
      setTimeout(() => {
        if (chartCanvas.value) {
          renderChart()
        }
      }, 100);
    }
  }
}, { deep: true })
  
onMounted(async () => {
      await stockStore.fetchStockData(props.symbol, currentRange.value)
      initialLoading.value = false
    })

    watch(() => props.symbol, async (newSymbol) => {
  initialLoading.value = true
  await stockStore.fetchStockData(newSymbol, currentRange.value)
  initialLoading.value = false
})
  
      onBeforeUnmount(() => {
        destroyChart()
      })
  
      return {
        chartCanvas,
        timeRanges,
        currentRange,
        timeRangeLabel,
        updateRange,
        refreshData,
        loading: stockStore.loading,
        error: stockStore.error,
        stockStore,
        initialLoading,
        latestEntry
      }
    }
  })
  </script>