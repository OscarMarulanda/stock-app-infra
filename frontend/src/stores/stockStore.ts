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

interface StockRecommendation {
  symbol: string
  momentum: number
  rsi: number
}

export const useStockStore = defineStore('stock', {
  state: () => ({
    stockData: [] as StockData[],
    stockRecommendations: [] as StockRecommendation[],
    loadingRecommendations: false,
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
      this.stockData = [] // Reset data
      
      try {
        // First try to get existing data
        const response = await api.get(`/stocks/${symbol}`, {
          params: { range },
          timeout: 30000
        })
    
        // Handle case where data is null/undefined or empty array
        if (!response.data || response.data.length === 0) {
          try {
            // Try to refresh data
            await api.post(`/stocks/${symbol}/refresh`, {}, { timeout: 30000 })
            
            // Fetch again after refresh
            const refreshedResponse = await api.get(`/stocks/${symbol}`, {
              params: { range },
              timeout: 30000
            })
            
            if (!refreshedResponse.data || refreshedResponse.data.length === 0) {
              throw new Error(`No data available for symbol "${symbol}"`)
            }
            
            this.stockData = refreshedResponse.data
          } catch (refreshErr: any) {
            // Handle 404 specifically
            if (refreshErr.response?.status === 404) {
              throw new Error(`Stock symbol "${symbol}" not found`)
            }
            throw refreshErr // Re-throw other errors
          }
        } else {
          this.stockData = response.data
        }
      } catch (err: any) {
        console.error('API Error:', err)
        
        // Set user-friendly error message
        if (err.message.includes('not found') || err.response?.status === 404) {
          this.error = `Stock symbol "${symbol}" not found`
        } else if (err.message.includes('No data available')) {
          this.error = `No data available for symbol "${symbol}"`
        } else if (err.code === 'ECONNABORTED') {
          this.error = 'The request took too long. Please try again later.'
        } else {
          this.error = err.response?.data?.message || err.message || 'Failed to fetch stock data'
        }
        
        this.stockData = []
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
    },

    async fetchStockRecommendations(range: 'short' | 'medium' | 'long') {
      this.loadingRecommendations = true

      

      this.error = null
      this.stockRecommendations = []
    
      try {
        const response = await api.get('/stocksRecommend', {
          params: { range },
          timeout: 30000
        })
    
        if (!response.data || response.data.length === 0) {
          this.error = 'No stock recommendations found for the selected range.'
          return
        }
    
        this.stockRecommendations = response.data
      } catch (err: any) {
        console.error('Recommendation API Error:', err)
    
        if (err.code === 'ECONNABORTED') {
          this.error = 'The request timed out. Please try again.'
        } else if (err.response?.status === 400) {
          this.error = 'Invalid range selected.'
        } else {
          this.error = err.response?.data?.message || err.message || 'Failed to fetch recommendations.'
        }
      } finally {
        this.loadingRecommendations  = false
      }
    }
  }
})