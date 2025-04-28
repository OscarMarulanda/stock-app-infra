import { defineStore } from 'pinia'
import axios from 'axios'

const api = axios.create({
    baseURL: 'http://localhost:8080/api',  // Explicitly setting the backend URL for development
    timeout: 10000
  })

interface StockData {
  date: string
  open: number
  high: number
  low: number
  close: number
  volume: number
}

export const useStockStore = defineStore('stock', {
  state: () => ({
    stockData: [] as StockData[],
    loading: false,
    error: null as string | null,
    currentSymbol: '',
    currentRange: 'month'
  }),
  actions: {
    async fetchStockData(symbol: string, range: string) {
      this.loading = true
      this.error = null
      this.currentSymbol = symbol
      this.currentRange = range
      
      try {
        // First try to get existing data
        const response = await api.get(`/stocks/${symbol}`, {
          params: { range }
        })
        
        // If no data exists, trigger a refresh
        if (response.data.length === 0) {
          await api.post(`/stocks/${symbol}/refresh`)
          // Fetch again after refresh
          const refreshedResponse = await api.get(`/stocks/${symbol}`, {
            params: { range }
          })
          this.stockData = refreshedResponse.data
        } else {
          this.stockData = response.data
        }
        
        console.log(this.stockData)
      } catch (err) {
        this.error = err instanceof Error ? err.message : 'Failed to fetch stock data'
      } finally {
        this.loading = false
      }
    }
  }
})