import { defineStore } from 'pinia'
import axios from 'axios'

const api = axios.create({
    baseURL: 'http://localhost:8080/api',  // Explicitly setting the backend URL for development
    timeout: 30000
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
    currentRange: 'month',
    isRefreshing: false
  }),
  actions: {
    async fetchStockData(symbol: string, range: string) {
      this.loading = true
      this.error = null
      this.currentSymbol = symbol
      this.currentRange = range
      
      try {
        // First try to get existing data with extended timeout
        const response = await api.get(`/stocks/${symbol}`, {
          params: { range },
          timeout: 30000
        })
        
        // If no data exists, trigger a refresh (with separate timeout)
        if (response.data.length === 0) {
          try {
            await api.post(`/stocks/${symbol}/refresh`, {}, { timeout: 30000 })
            // Fetch again after refresh
            const refreshedResponse = await api.get(`/stocks/${symbol}`, {
              params: { range },
              timeout: 30000
            })
            this.stockData = refreshedResponse.data
          } catch (refreshErr) {
            console.error('Refresh error:', refreshErr)
            // Even if refresh fails, try to get whatever data exists
            const fallbackResponse = await api.get(`/stocks/${symbol}`, {
              params: { range },
              timeout: 30000
            })
            this.stockData = fallbackResponse.data
          }
        } else {
          this.stockData = response.data
        }
      } catch (err: any) {
        console.error('API Error:', err)
        if (err.code === 'ECONNABORTED') {
          // On timeout, check if we have any data at all
          try {
            const fallback = await api.get(`/stocks/${symbol}`, {
              params: { range },
              timeout: 30000
            })
            if (fallback.data.length > 0) {
              this.stockData = fallback.data
              return // We got data despite the initial timeout
            }
          } catch (fallbackErr) {
            console.error('Fallback failed:', fallbackErr)
          }
          this.error = 'The request took too long. Data may still be loading...'
        } else {
          this.error = err.response?.data?.message || err.message || 'Failed to fetch stock data'
        }
      } finally {
        this.loading = false
      }
    },

    async refreshStockData() {
      this.isRefreshing = true
      this.loading = true
      this.error = null
      
      try {
        // Trigger refresh on backend
        await api.post(`/stocks/${this.currentSymbol}/refresh`, {}, { timeout: 30000 })
        
        // Fetch the updated data
        await this.fetchStockData(this.currentSymbol, this.currentRange)
        
      } catch (err: any) {
        console.error('Refresh error:', err)
        this.error = err.response?.data?.message || err.message || 'Failed to refresh stock data'
      } finally {
        this.loading = false
        this.isRefreshing = false
      }
    }
  }
})