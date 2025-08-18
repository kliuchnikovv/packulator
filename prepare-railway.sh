#!/bin/bash

# Prepare Packulator for Railway deployment
set -e

echo "🚂 Preparing Packulator for Railway deployment..."

# Check if we're in a git repository
if [ ! -d ".git" ]; then
    echo "❌ Not in a git repository. Initializing..."
    git init
    git add .
    git commit -m "Initial commit"
fi

# Check if frontend build works
echo "🔨 Testing frontend build..."
if [ -d "frontend" ]; then
    cd frontend
    if [ ! -d "node_modules" ]; then
        echo "📦 Installing frontend dependencies..."
        npm ci
    fi
    echo "🏗️  Building frontend..."
    npm run build
    cd ..
else
    echo "❌ Frontend directory not found!"
    exit 1
fi

# Check if Go build works
echo "🔨 Testing Go build..."
if [ ! -f "go.mod" ]; then
    echo "❌ go.mod not found!"
    exit 1
fi

echo "📦 Downloading Go dependencies..."
go mod download

echo "🏗️  Building Go application..."
CGO_ENABLED=0 go build -o main ./cmd/main.go

if [ ! -f "main" ]; then
    echo "❌ Go build failed!"
    exit 1
fi

# Clean up build artifact
rm -f main

# Check Railway config files
echo "⚙️  Checking Railway configuration..."

required_files=(
    "railway.json"
    "nixpacks.toml" 
    "Dockerfile.railway"
    ".env.railway"
)

for file in "${required_files[@]}"; do
    if [ ! -f "$file" ]; then
        echo "❌ Missing required file: $file"
        exit 1
    else
        echo "✅ Found: $file"
    fi
done

# Check for environment variables template
if [ -f ".env.railway" ]; then
    echo "✅ Environment variables template ready"
else
    echo "❌ .env.railway template missing"
    exit 1
fi

# Commit Railway configuration if changes exist
if ! git diff --cached --quiet || ! git diff --quiet; then
    echo "📝 Committing Railway configuration..."
    git add .
    git commit -m "Add Railway deployment configuration

- Added railway.json for Railway configuration
- Added nixpacks.toml for build process  
- Added Dockerfile.railway as fallback
- Added .env.railway environment template
- Ready for Railway deployment!"
    
    echo "✅ Changes committed!"
else
    echo "✅ No changes to commit"
fi

# Show next steps
echo ""
echo "🎉 Railway deployment preparation complete!"
echo ""
echo "📋 Next steps:"
echo "1. Push to GitHub: git push origin main"
echo "2. Go to railway.app and login with GitHub"
echo "3. Create new project from your GitHub repo"
echo "4. Add PostgreSQL database service"  
echo "5. Set environment variables from .env.railway"
echo ""
echo "🌍 Your app will be live at: https://your-app.up.railway.app"
echo ""
echo "📖 Full guide: cat RAILWAY_DEPLOY.md"