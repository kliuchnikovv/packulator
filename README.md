# Packulator Backend API

[![CI Pipeline](https://github.com/kliuchnikovv/packulator/actions/workflows/ci.yml/badge.svg)](https://github.com/kliuchnikovv/packulator/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/kliuchnikovv/packulator/branch/main/graph/badge.svg)](https://codecov.io/gh/kliuchnikovv/packulator)
[![Go Report Card](https://goreportcard.com/badge/github.com/kliuchnikovv/packulator)](https://goreportcard.com/report/github.com/kliuchnikovv/packulator)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/kliuchnikovv/packulator)](go.mod)

Go-based HTTP API that calculates the number of shipping packs needed for customer orders.

## 🚀 Features

- **Pack calculation API** - Calculate optimal pack combinations for any order amount
- **Flexible pack configuration** - Pack sizes configurable via API without code changes  
- **PostgreSQL database** - Persistent storage for pack configurations
- **Health checks** - Built-in monitoring endpoints
- **Comprehensive tests** - Unit tests with high coverage
- **Railway deployment** - One-click deploy to production

## 📚 API Endpoints

### Pack Management
- `POST /packs/create` - Create new pack configuration
- `GET /packs/list` - List all available packs
- `GET /packs/id?id={id}` - Get specific pack by ID
- `GET /packs/hash?hash={hash}` - Get packs by version hash
- `DELETE /packs/delete?id={id}` - Delete pack configuration

### Pack Calculation  
- `GET /packaging/number_of_packages?amount={amount}&packs_hash={hash}` - Calculate pack combinations

### Health
- `GET /health/check` - Service health status

## 🛠️ Technology Stack

- **Backend**: Go 1.24+ with [Engi framework](https://github.com/kliuchnikovv/engi)
- **Database**: PostgreSQL with GORM
- **Deployment**: Railway with automatic PostgreSQL
- **Testing**: testify/mock for comprehensive unit tests

## 🚂 Quick Deploy

[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/template/go-api)

1. Fork this repository
2. Connect to [Railway](https://railway.app)
3. Add PostgreSQL database
4. Deploy automatically

See [DEPLOY.md](DEPLOY.md) for detailed instructions.

## 💻 Local Development

### Prerequisites
- Go 1.24+
- PostgreSQL database

### Setup
```bash
# Clone repository
git clone https://github.com/your-username/packulator.git
cd packulator

# Install dependencies
go mod download

# Set environment variables
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=your-password
export DB_NAME=packulator
export DB_SSL_MODE=disable

# Run application
go run cmd/main.go
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./internal/service -v
```

## 📖 API Usage Examples

### Create Pack Configuration
```bash
curl -X POST http://localhost:8080/packs/create \
  -H "Content-Type: application/json" \
  -d '{"packs": [250, 500, 1000, 2000, 5000]}'
```

### Calculate Pack Combinations
```bash
# First get the version hash from pack creation response
curl "http://localhost:8080/packaging/number_of_packages?amount=1001&packs_hash=abc123def456"
```

Response:
```json
{
  "250": 1,
  "500": 0, 
  "1000": 1
}
```

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP Client   │───▶│  Go Backend API │───▶│   PostgreSQL    │
│                 │    │                 │    │                 │
│ Frontend/Curl/  │    │ • Pack CRUD     │    │ • Pack configs  │
│ Postman/etc     │    │ • Calculations  │    │ • Version hashes│
│                 │    │ • Health checks │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 📁 Project Structure

```
├── cmd/main.go                 # Application entry point
├── internal/
│   ├── api/                    # HTTP handlers
│   │   ├── packs.go           # Pack CRUD endpoints  
│   │   ├── calculate.go       # Pack calculation endpoints
│   │   └── health.go          # Health check endpoints
│   ├── service/               # Business logic
│   │   ├── pack.go            # Pack management service
│   │   └── packaging.go       # Pack calculation algorithms
│   ├── store/                 # Data access layer
│   │   └── pack.go           # PostgreSQL operations
│   ├── model/                 # Data models
│   │   └── package.go        # Pack and request/response models
│   └── config/                # Configuration
│       ├── app.go            # Application config
│       └── database.go       # Database config
├── nixpacks.toml              # Railway build configuration  
├── railway.json               # Railway service configuration
├── Dockerfile                 # Container configuration
└── DEPLOY.md                  # Deployment guide
```

## 🔧 Configuration

Environment variables:
- `PORT` - Server port (default: 8080)  
- `HOST` - Server host (default: 0.0.0.0)
- `ENVIRONMENT` - App environment (development/production)
- `DB_HOST` - Database host
- `DB_PORT` - Database port (default: 5432)
- `DB_USER` - Database username
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name
- `DB_SSL_MODE` - SSL mode (disable/require)
- `LOG_LEVEL` - Logging level (debug/info/warn/error)
- `DEBUG` - Debug mode (true/false)

## 📊 Algorithm

The pack calculation uses a **dynamic programming algorithm** optimized for minimal overshoot and pack count:

### Algorithm Overview
1. **Dynamic Programming Table** - Build all possible pack combinations up to `amount + largest_pack`
2. **Optimal Selection** - For each sum, keep the variant with:
   - **Primary**: Minimal overshoot (excess over target amount)  
   - **Secondary**: Minimal pack count (when overshoot is equal)
3. **Result Selection** - Choose the best variant among all sums ≥ target amount

### Implementation Details
```go
// Each variant tracks:
type variant struct {
    numberOfPacks int64            // Total packs used
    overshoot     int64            // Amount exceeding target (sum - amount)
    combination   map[int64]int64  // Pack size → quantity mapping
}

// Priority function: less overshoot wins, then fewer packs
func isBetter(left, right *variant) bool {
    if left.overshoot < right.overshoot {
        return true
    } else if left.overshoot > right.overshoot {
        return false
    }
    return left.numberOfPacks < right.numberOfPacks
}
```

### Example Calculation
For **amount = 1001** with **packs = [250, 500, 1000]**:

**Step 1**: Build all combinations from 0 to 2001 (1001 + 1000)
**Step 2**: Find optimal variants for sums ≥ 1001:
- Sum 1250: {250: 1, 1000: 1} → overshoot = 249, packs = 2 ✅
- Sum 1500: {500: 3} → overshoot = 499, packs = 3  
- Sum 1750: {250: 3, 1000: 1} → overshoot = 749, packs = 4

**Result**: `{250: 1, 1000: 1}` = **2 packs with minimal overshoot of 249**

### Algorithm Benefits
- **Guaranteed Optimal**: Always finds the solution with minimal overshoot
- **Deterministic**: Same input always produces same output  
- **Efficient**: O(amount × pack_count) time complexity
- **Flexible**: Works with any pack size combination

## 🧪 Testing & CI/CD

### Test Coverage
- **Overall Coverage**: 81%+
- **API Layer**: 81.1% ✅
- **Service Layer**: 100% ✅  
- **Config Layer**: 100% ✅
- **Model Layer**: 100% ✅

### Running Tests Locally
```bash
# Run all tests
make test

# Run tests with coverage report
make test-coverage-html

# Run integration tests (requires PostgreSQL)
make test-integration

# Run all quality checks
make lint
make sec-scan
```

### CI Pipeline
The project uses GitHub Actions for continuous integration:
- ✅ **Linting** - golangci-lint with comprehensive rules
- ✅ **Testing** - Unit tests, integration tests, race detection
- ✅ **Security** - Gosec and Nancy vulnerability scanning  
- ✅ **Build** - Multi-architecture builds (Linux, macOS, Windows)
- ✅ **Coverage** - Automatic coverage reporting to Codecov

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## 📄 License

This project is licensed under the MIT License.
