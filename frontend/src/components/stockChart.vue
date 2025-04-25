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
  
      <div v-if="loading" class="flex justify-center items-center h-64">
        <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
      </div>
  
      <div v-else-if="error" class="text-red-500 p-4">
        {{ error }}
      </div>
  
      <div v-else-if="stockStore.stockData.length === 0" class="text-gray-500 p-4">
  No data available for this time range.
</div>
  
<div v-else class="h-96">
  <canvas ref="chartCanvas"></canvas>
</div>
    </div>
  </template>
  
  <script lang="ts">
  import { defineComponent, onMounted, ref, watch, onBeforeUnmount, nextTick  } from 'vue'
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
  
      const refreshData = () => {
        stockStore.refreshStockData()
      }
  
      watch(() => stockStore.stockData, async (newData) => {
  if (Array.isArray(newData) && newData.length > 0) {
    await nextTick()  // ⏳ Wait for DOM (including canvas) to update
    if (chartCanvas.value) {
      renderChart()
    } else {
      console.warn('Canvas still not available after nextTick')
    }
  }
}, { deep: true })
  
      onMounted(() => {
        stockStore.fetchStockData(props.symbol, currentRange.value)
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
      }
    }
  })
  </script>