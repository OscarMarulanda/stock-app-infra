# Stock App Infrastructure (Terraform)

This repo contains the infrastructure as code (IaC) for a stock recommendation app using Terraform and AWS.

## 🚀 Stack

- **Terraform** (>= 1.5)
- **AWS** (EC2, S3, SSM Parameter Store)
- **Amazon Linux 2**
- **Alpha Vantage API**

## 📦 Modules

- `frontend/`: Sets up an S3 bucket to host a Vue.js frontend as a static website.
- `backend/`: Provisions an EC2 instance for the Golang API server.
- `ssm/`: Stores the Alpha Vantage API key securely in SSM Parameter Store.

## 🔐 Sensitive Inputs

| Variable                  | Description                           |
|---------------------------|---------------------------------------|
| `alpha_vantage_api_key`   | Alpha Vantage API Key                 |
| `key_name`                | Existing EC2 Key Pair for SSH access  |
| `aws_region`              | AWS region (default: `us-east-2`)     |

## 📁 File Structure

. ├── backend/ │ └── main.tf ├── frontend/ │ └── main.tf ├── ssm/ │ └── main.tf ├── variables.tf ├── provider.tf ├── main.tf └── terraform.tfstate (ignored)

csharp
Copy
Edit

## 🧪 How to Deploy

```bash
terraform init
terraform plan
terraform apply



Stock Data Visualization Dashboard
A full-stack application that fetches stock market data from Alpha Vantage API, stores it in a database, and displays interactive charts with Vue.js.

Screenshot (Optional: Add screenshot later)

Features
View historical stock data (Open, High, Low, Close, Volume)
Select time ranges (1 week, 1 month, 6 months, 1 year)
Automatic data refresh from Alpha Vantage API
Local database caching to minimize API calls
Responsive design with Tailwind CSS
Tech Stack
Frontend:

Vue 3 (Composition API)
TypeScript
Pinia (State management)
Chart.js (Data visualization)
Tailwind CSS (Styling)
Axios (HTTP client)
Backend:

Go (Golang)
Gorilla Mux (Router)
CockroachDB/PostgreSQL (Database)
RESTy (HTTP client)
Prerequisites
Go 1.20+
Node.js 16+
CockroachDB/PostgreSQL
Alpha Vantage API key (free tier available)
Setup Instructions
1. Backend Setup
# Clone the repository
git clone https://github.com/yourusername/stock-visualizer.git
cd stock-visualizer/backend

# Install dependencies
go mod tidy

# Set up environment variables
cp .env.example .env
# Edit .env with your credentials

Configure .env:

env
ALPHA_VANTAGE_API_KEY=your_api_key_here
COCKROACHDB_DSN="postgresql://user:password@localhost:26257/defaultdb?sslmode=disable"
PORT=8080
2. Database Setup
Install CockroachDB (instructions)

Start the database:

bash
cockroach start-single-node --insecure --listen-addr=localhost
The Go backend will automatically create the required tables on first run.

3. Frontend Setup
bash
cd ../frontend

# Install dependencies
npm install

# Configure environment
cp .env.example .env
# Set your backend URL
VITE_API_BASE_URL=http://localhost:8080
Running the Application
Start Backend:

bash
cd backend
go run main.go
Start Frontend:

bash
cd ../frontend
npm run dev
Open your browser to:

http://localhost:5173
Project Structure
/backend
  ├── main.go            # Core server logic
  ├── go.mod             # Go dependencies
  └── .env               # Environment config

/frontend
  ├── src/
  │   ├── components/    # Vue components
  │   ├── stores/        # Pinia stores
  │   ├── App.vue        # Main app component
  │   └── main.ts        # App entry point
  ├── package.json
  └── vite.config.ts     # Build configuration
API Endpoints
Endpoint	Method	Description
/api/stocks/{symbol}	GET	Get stock data for time range
/api/stocks/{symbol}/refresh	POST	Force data refresh from Alpha Vantage
Troubleshooting
Port 8080 in use:

bash
# Windows
netstat -ano | findstr :8080
taskkill /PID [PID] /F

# Mac/Linux
lsof -i :8080
kill -9 [PID]
CORS Errors:

Verify VITE_API_BASE_URL matches your backend URL

Check CORS settings in main.go

License
MIT
