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
        const response = await api.get(`/stocks/${symbol}`, { // Use the custom axios instance
          params: { range }
        })
        this.stockData = response.data
        console.log(this.stockData)
      } catch (err) {
        this.error = err instanceof Error ? err.message : 'Failed to fetch stock data'
      } finally {
        this.loading = false
      }
    },
    async refreshStockData() {
      if (!this.currentSymbol) return
      
      this.loading = true
      this.error = null
      
      try {
        await api.post(`/stocks/${this.currentSymbol}/refresh`)  // Use the custom axios instance
        // After refresh, fetch the data again
        await this.fetchStockData(this.currentSymbol, this.currentRange)
      } catch (err) {
        this.error = err instanceof Error ? err.message : 'Failed to refresh stock data'
      } finally {
        this.loading = false
      }
    }
  }
})