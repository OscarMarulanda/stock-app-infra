package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

type TimeSeriesResponse struct {
	MetaData struct {
		Information string `json:"1. Information"`
		Symbol      string `json:"2. Symbol"`
	} `json:"Meta Data"`
	TimeSeriesDaily map[string]struct {
		Open   string `json:"1. open"`
		High   string `json:"2. high"`
		Low    string `json:"3. low"`
		Close  string `json:"4. close"`
		Volume string `json:"5. volume"`
	} `json:"Time Series (Daily)"`

	ErrorMessage string `json:"Error Message"`
}

func main() {

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // For development, you can use "*" but restrict in production
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            true, // Enable this to see CORS logs
	})

	err := godotenv.Load("./.env")
	if err != nil {
		log.Println("Warning: No .env file found (using system env vars)")
	}

	dsn := os.Getenv("COCKROACHDB_DSN")
	if dsn == "" {
		log.Fatal("Missing required environment variable COCKROACHDB_DSN")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	initDB(db)

	r := mux.NewRouter()
	r.HandleFunc("/api/stocks/{symbol}", getStockDataHandler(db)).Methods("GET")
	r.HandleFunc("/api/stocksRecommend", getAllStockRecDataHandler(db)).Methods("GET")
	r.HandleFunc("/api/stocks/{symbol}/refresh", refreshStockDataHandler(db)).Methods("POST")

	handler := c.Handler(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func initDB(db *sql.DB) {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS stock_data (
		date DATE NOT NULL,
		symbol STRING NOT NULL,
		open FLOAT,
		high FLOAT,
		low FLOAT,
		close FLOAT,
		volume INT,
		last_updated TIMESTAMP DEFAULT now(),
		PRIMARY KEY (date, symbol)
	)
	`)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
}

func getStockDataHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		// Log incoming request
		log.Printf("🔍 Incoming GET request: %s %s", r.Method, r.URL.Path)
		log.Printf("Headers: %v", r.Header)

		vars := mux.Vars(r)
		symbol := vars["symbol"]
		timeRange := r.URL.Query().Get("range")

		// Log symbol and range
		log.Printf("📊 Symbol: %s, Time Range: %s", symbol, timeRange)

		if symbol == "" {
			http.Error(w, "Symbol is required", http.StatusBadRequest)
			log.Println("❌ Error: Symbol is required")
			return
		}

		if timeRange == "" {
			http.Error(w, "Range parameter is required", http.StatusBadRequest)
			log.Println("❌ Error: Range parameter is required")
			return
		}

		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM stock_data WHERE symbol = $1", symbol).Scan(&count)
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			log.Println("❌ Error: Database query failed", err)
			return
		}

		// If no data exists, fetch from API
		if count == 0 {
			_, err := fetchAndStoreData(db, symbol)
			if err != nil {
				if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "no data available") {
					http.Error(w, fmt.Sprintf("Stock symbol '%s' not found", symbol), http.StatusNotFound)
				} else {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				log.Println("❌ Error fetching initial data:", err)
				return
			}
		}

		var query string
		var args []interface{}

		switch timeRange {
		case "week":
			query = `SELECT date, open, high, low, close, volume 
					FROM stock_data 
					WHERE symbol = $1 
					ORDER BY date DESC 
					LIMIT 7`
			args = []interface{}{symbol}
		case "month":
			query = `SELECT date, open, high, low, close, volume 
					FROM stock_data 
					WHERE symbol = $1 
					ORDER BY date DESC 
					LIMIT 30`
			args = []interface{}{symbol}
		case "6month":
			query = `SELECT date, open, high, low, close, volume 
					FROM stock_data 
					WHERE symbol = $1 AND date >= CURRENT_DATE - INTERVAL '6 months' 
					ORDER BY date DESC`
			args = []interface{}{symbol}
		case "year":
			query = `SELECT date, open, high, low, close, volume 
					FROM stock_data 
					WHERE symbol = $1 AND date >= CURRENT_DATE - INTERVAL '1 year' 
					ORDER BY date DESC`
			args = []interface{}{symbol}
		default:
			http.Error(w, "Invalid time range specified", http.StatusBadRequest)
			log.Println("❌ Error: Invalid time range:", timeRange)
			return
		}

		rows, err := db.Query(query, args...)
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			log.Println("❌ Error: Database query failed", err)
			return
		}
		defer rows.Close()

		type StockData struct {
			Date   string  `json:"date"`
			Open   float64 `json:"open"`
			High   float64 `json:"high"`
			Low    float64 `json:"low"`
			Close  float64 `json:"close"`
			Volume int     `json:"volume"`
		}

		var data []StockData
		for rows.Next() {
			var sd StockData
			var date time.Time
			err := rows.Scan(&date, &sd.Open, &sd.High, &sd.Low, &sd.Close, &sd.Volume)
			if err != nil {
				http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
				log.Println("❌ Error: Scanning row", err)
				return
			}
			sd.Date = date.Format("2006-01-02")
			data = append(data, sd)
		}

		// Log response data
		log.Printf("✅ Successfully fetched %d stock data entries for symbol %s", len(data), symbol)

		// Send the response
		w.Header().Set("Content-Type", "application/json")

		if len(data) == 0 {
			http.Error(w, fmt.Sprintf("No data found for symbol '%s'", symbol), http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(data)
	}
}

func calculateMomentum(prices []float64) float64 {
	if len(prices) < 2 || prices[len(prices)-1] == 0 {
		return 0
	}
	return ((prices[0] / prices[len(prices)-1]) - 1) * 100
}

func calculateRSI(prices []float64) float64 {
	if len(prices) < 2 {
		return 0
	}

	var gains, losses float64
	for i := 1; i < len(prices); i++ {
		change := prices[i-1] - prices[i]
		if change > 0 {
			gains += change
		} else {
			losses -= change
		}
	}

	if gains+losses == 0 {
		return 50
	}

	avgGain := gains / float64(len(prices)-1)
	avgLoss := losses / float64(len(prices)-1)

	if avgLoss == 0 {
		return 100
	}

	rs := avgGain / avgLoss
	return 100 - (100 / (1 + rs))
}

func getAllSymbols(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT DISTINCT symbol FROM stock_data")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var symbols []string
	for rows.Next() {
		var symbol string
		if err := rows.Scan(&symbol); err != nil {
			return nil, err
		}
		symbols = append(symbols, symbol)
	}
	return symbols, nil
}

func getAllStockRecDataHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		// Log incoming request
		log.Printf("🔍 Incoming GET request: %s %s", r.Method, r.URL.Path)
		log.Printf("Headers: %v", r.Header)
		timeRange := r.URL.Query().Get("range")

		if timeRange == "" {
			http.Error(w, "Range parameter is required", http.StatusBadRequest)
			log.Println("❌ Error: Range parameter is required")
			return
		}

		symbols, err := getAllSymbols(db)
		if err != nil {
			http.Error(w, "Failed to retrieve symbols: "+err.Error(), http.StatusInternalServerError)
			return
		}

		type StockRecResponse struct {
			Symbol   string    `json:"symbol"`
			Data     []float64 `json:"data"`
			Momentum float64   `json:"momentum"`
			RSI      float64   `json:"rsi"`
		}

		var results []StockRecResponse

		// Log range
		log.Printf("📊 Time Range: %s", timeRange)

		for _, symbol := range symbols {
			var query string
			switch timeRange {
			case "short":
				query = `SELECT close FROM stock_data WHERE symbol = $1 ORDER BY date DESC LIMIT 7`
			case "medium":
				query = `SELECT close FROM stock_data WHERE symbol = $1 ORDER BY date DESC LIMIT 20`
			case "long":
				query = `SELECT close FROM stock_data WHERE symbol = $1 ORDER BY date DESC LIMIT 60`
			default:
				http.Error(w, "Invalid time range specified", http.StatusBadRequest)
				return
			}

			rows, err := db.Query(query, symbol)
			if err != nil {
				log.Printf("❌ Failed to query data for %s: %v", symbol, err)
				continue
			}

			var closes []float64
			for rows.Next() {
				var close float64
				if err := rows.Scan(&close); err != nil {
					log.Printf("❌ Failed to scan close for %s: %v", symbol, err)
					rows.Close()
					continue
				}
				closes = append(closes, close)
			}
			rows.Close()

			if len(closes) == 0 {
				continue
			}

			// Log response data

			// Send the response
			w.Header().Set("Content-Type", "application/json")

			result := StockRecResponse{
				Symbol:   symbol,
				Data:     closes,
				Momentum: calculateMomentum(closes),
				RSI:      calculateRSI(closes),
			}
			results = append(results, result)

		}
		sort.Slice(results, func(i, j int) bool {
			return results[i].Momentum > results[j].Momentum
		})
		json.NewEncoder(w).Encode(results)
		log.Printf("✅ Successfully fetched %d stock data recommendation entries", len(results))
	}
}

func refreshStockDataHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		symbol := vars["symbol"]

		newRecords, err := fetchAndStoreData(db, symbol)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":    "Stock data refreshed",
			"newRecords": newRecords,
		})
	}
}

func fetchAndStoreData(db *sql.DB, symbol string) (int, error) {
	apiKey := os.Getenv("ALPHA_VANTAGE_API_KEY")
	if apiKey == "" {
		return 0, fmt.Errorf("missing API key configuration")
	}

	baseURL := "https://www.alphavantage.co/query"

	// Fetch data from Alpha Vantage
	data, err := fetchStockData(resty.New(), baseURL, apiKey, symbol, "compact")
	if err != nil {
		return 0, fmt.Errorf("failed to fetch stock data: %w", err)
	}

	if data == nil || data.ErrorMessage != "" || len(data.TimeSeriesDaily) == 0 {
		return 0, fmt.Errorf("symbol '%s' not found or invalid", symbol)
	}

	// Get the latest date we already have in DB
	var latestDate sql.NullString
	err = db.QueryRow("SELECT MAX(date) FROM stock_data WHERE symbol = $1", symbol).Scan(&latestDate)
	if err != nil {
		return 0, fmt.Errorf("failed to query latest date: %w", err)
	}

	// Prepare the insert statement
	stmt, err := db.Prepare(`
        INSERT INTO stock_data (date, symbol, open, high, low, close, volume)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (date, symbol) DO UPDATE SET
            open = EXCLUDED.open,
            high = EXCLUDED.high,
            low = EXCLUDED.low,
            close = EXCLUDED.close,
            volume = EXCLUDED.volume,
            last_updated = now()
    `)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	newRecords := 0
	for dateStr, stockData := range data.TimeSeriesDaily {
		// Skip if we already have this data (unless we want to force update)
		if latestDate.Valid && dateStr <= latestDate.String {
			continue
		}

		// Convert string values to appropriate types
		open, err := strconv.ParseFloat(stockData.Open, 64)
		if err != nil {
			log.Printf("Failed to parse open price for %s: %v", dateStr, err)
			continue
		}

		high, err := strconv.ParseFloat(stockData.High, 64)
		if err != nil {
			log.Printf("Failed to parse high price for %s: %v", dateStr, err)
			continue
		}

		low, err := strconv.ParseFloat(stockData.Low, 64)
		if err != nil {
			log.Printf("Failed to parse low price for %s: %v", dateStr, err)
			continue
		}

		closeVal, err := strconv.ParseFloat(stockData.Close, 64)
		if err != nil {
			log.Printf("Failed to parse close price for %s: %v", dateStr, err)
			continue
		}

		volume, err := strconv.Atoi(stockData.Volume)
		if err != nil {
			log.Printf("Failed to parse volume for %s: %v", dateStr, err)
			continue
		}

		// Execute the insert
		_, err = stmt.Exec(
			dateStr,
			symbol,
			open,
			high,
			low,
			closeVal,
			volume,
		)
		if err != nil {
			log.Printf("Failed to insert data for %s: %v", dateStr, err)
			continue
		}

		newRecords++
	}

	return newRecords, nil
}

func fetchStockData(client *resty.Client, baseURL, apiKey, symbol, size string) (*TimeSeriesResponse, error) {
	// Implement retry logic for API calls
	var result *TimeSeriesResponse
	var err error

	for attempt := 1; attempt <= 3; attempt++ {
		resp, e := client.R().
			SetQueryParams(map[string]string{
				"function":   "TIME_SERIES_DAILY",
				"symbol":     symbol,
				"apikey":     apiKey,
				"outputsize": size,
			}).
			SetResult(&TimeSeriesResponse{}).
			Get(baseURL)

		if e != nil {
			err = fmt.Errorf("attempt %d: failed to get stock data: %v", attempt, e)
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}

		if resp.IsError() {
			err = fmt.Errorf("attempt %d: API error: %v", attempt, resp.Status())
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}

		result = resp.Result().(*TimeSeriesResponse)
		err = nil
		break
	}

	return result, err
}
