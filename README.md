# Packulator Backend API

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

The pack calculation uses a greedy algorithm optimized for minimal pack count:

1. **Sort packs** by size (largest first)
2. **Greedy selection** - use largest packs that fit
3. **Optimization** - minimize total pack count
4. **Validation** - ensure complete coverage of order amount

Example for amount 1001 with packs [250, 500, 1000]:
- 1000 × 1 = 1000 (remainder: 1)  
- 250 × 1 = 250 (remainder: 0)
- **Result**: {1000: 1, 250: 1} = 1250 items in 2 packs

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## 📄 License

This project is licensed under the MIT License.
