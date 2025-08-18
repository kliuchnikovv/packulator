# 🚂 Railway: Frontend + Backend отдельно

Деплой frontend (React) и backend (Go) как отдельные сервисы в Railway.

## 🎯 Архитектура

```
Frontend (React)     Backend (Go)        Database
Port 3000           Port 8080           PostgreSQL
     ↓                   ↓                   ↓
Railway Service 1   Railway Service 2   Railway Service 3
```

## 🚀 Деплой (2 сервиса)

### 1. Подготовьте репозиторий
```bash
git add .
git commit -m "Add separate Railway configs"
git push origin main
```

### 2. Создайте проект в Railway
- Перейдите на [railway.app](https://railway.app)
- **New Project** → **Deploy from GitHub repo**
- Выберите ваш `packulator` репозиторий

### 3. Деплой Backend (Go)

В созданном проекте:

1. **Переименуйте сервис** в `packulator-backend`
2. **Settings** → **Source** → **Root Directory**: `/` (корень)
3. **Variables** → добавьте:
   ```bash
   ENVIRONMENT=production
   LOG_LEVEL=info
   DEBUG=false
   PORT=8080
   ```
4. **Переместите файл**: `cp nixpacks-backend.toml nixpacks.toml`

### 4. Деплой Frontend (React)

В том же проекте:

1. **+ New** → **GitHub Repo** → выберите тот же репозиторий
2. **Переименуйте сервис** в `packulator-frontend`  
3. **Settings** → **Source** → **Root Directory**: `/frontend`
4. **Variables** → добавьте:
   ```bash
   PORT=3000
   REACT_APP_API_URL=https://packulator-backend.up.railway.app
   ```

### 5. Добавьте PostgreSQL
1. **+ New** → **Database** → **Add PostgreSQL** 
2. Railway автоматически свяжет с backend сервисом

## 📝 Конфигурации

### Backend nixpacks.toml (корень проекта):
```toml
[phases.setup]
nixPkgs = ["go_1_24"]

[phases.build] 
cmds = [
    "go mod download",
    "CGO_ENABLED=0 go build -o main ./cmd/main.go"
]

[start]
cmd = "./main"
```

### Frontend nixpacks.toml (в папке frontend/):
```toml
[phases.setup]
nixPkgs = ["nodejs_18"]

[phases.install]
cmds = ["npm ci"]

[phases.build]
cmds = ["npm run build"]

[start]
cmd = "npx serve -s build -l $PORT"
```

## 🌐 URL-ы после деплоя

- **Frontend**: `https://packulator-frontend.up.railway.app`
- **Backend API**: `https://packulator-backend.up.railway.app`
- **Health Check**: `https://packulator-backend.up.railway.app/health/check`

## ⚙️ Переменные окружения

### Backend сервис:
```bash
# Основные настройки
ENVIRONMENT=production
LOG_LEVEL=info  
DEBUG=false
PORT=8080

# База данных (автоматически от PostgreSQL сервиса)
DATABASE_URL=postgresql://...
DB_HOST=...
DB_PORT=5432
DB_USER=...
DB_PASSWORD=... 
DB_NAME=...
DB_SSL_MODE=require
```

### Frontend сервис:
```bash
# Сервер настройки
PORT=3000

# API подключение  
REACT_APP_API_URL=https://packulator-backend.up.railway.app

# CORS (если нужно)
REACT_APP_CORS_ENABLED=true
```

## 🔧 Настройка CORS в Backend

Добавьте в Go код поддержку CORS для frontend:

```go
// В файле internal/api/*.go или middleware
func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "https://packulator-frontend.up.railway.app")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}
```

## 🎯 Преимущества отдельных сервисов

### ✅ Плюсы:
- **Независимый деплой** - можно обновлять frontend и backend отдельно
- **Масштабирование** - можно масштабировать сервисы независимо
- **Разработка** - команды могут работать параллельно
- **Кэширование** - frontend лучше кэшируется CDN
- **Отладка** - проще изолировать проблемы

### ⚠️ Минусы:
- **CORS** - нужно настраивать Cross-Origin запросы
- **Стоимость** - 2 сервиса вместо 1 (~$10-15/месяц)
- **Сложность** - больше moving parts

## 🛠️ Локальная разработка

### Запуск backend:
```bash
# Настройте переменные окружения
export PORT=8080
export DB_HOST=localhost
# ... другие переменные

go run cmd/main.go
# Backend доступен на http://localhost:8080
```

### Запуск frontend:
```bash
cd frontend

# Настройте API URL для локального backend
export REACT_APP_API_URL=http://localhost:8080

npm start  
# Frontend доступен на http://localhost:3000
```

## 🔄 CI/CD

Railway автоматически деплоит оба сервиса при пуше:

```bash
git add .
git commit -m "Update frontend and backend" 
git push origin main

# 🔄 Автоматически деплоятся:
# 1. packulator-backend (Go) 
# 2. packulator-frontend (React)
```

## 💰 Стоимость

### Примерная стоимость на Railway:
- **Backend**: ~$5-8/месяц
- **Frontend**: ~$3-5/месяц  
- **PostgreSQL**: ~$2-3/месяц
- **Total**: ~$10-16/месяц

### Сравнение с монолитом:
- **1 сервис**: ~$5-8/месяц
- **2 сервиса**: ~$10-16/месяц

## ✅ Checklist деплоя

### Backend:
- [ ] `nixpacks.toml` в корне проекта
- [ ] Переменные окружения настроены
- [ ] PostgreSQL подключен  
- [ ] Health check работает: `/health/check`
- [ ] CORS настроен для frontend

### Frontend:  
- [ ] `nixpacks.toml` в папке `frontend/`
- [ ] `serve` добавлен в `package.json`
- [ ] `REACT_APP_API_URL` указывает на backend
- [ ] Build процесс работает: `npm run build`
- [ ] Static файлы отдаются корректно

### Проверка:
```bash
# Backend API
curl https://packulator-backend.up.railway.app/health/check

# Frontend  
curl https://packulator-frontend.up.railway.app

# CORS тест
curl -H "Origin: https://packulator-frontend.up.railway.app" \
     https://packulator-backend.up.railway.app/packs/list
```

## 🎉 Готово!

После деплоя у вас будет:
- 🌍 **Frontend**: React приложение на отдельном домене  
- ⚡ **Backend**: Go API с собственным доменом
- 🗄️ **Database**: PostgreSQL автоматически подключен
- 🔒 **HTTPS**: SSL сертификаты автоматически
- 📊 **Мониторинг**: Логи и метрики для каждого сервиса

**Ваши пользователи заходят на frontend URL и получают полное приложение!** ✨