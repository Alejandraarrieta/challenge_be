# Variables
APP_NAME := myapp
MAIN_FILE := cmd/main.go

.PHONY: start_db run start redis postgres stop_db stop_redis stop_postgres

## 🔌 Inicia Redis (solo si no está corriendo)
redis:
	@echo "🔌 Iniciando Redis..."
	@pgrep redis-server >/dev/null || (redis-server > /dev/null 2>&1 & echo "Redis iniciado ✅")

## 🗄️  Inicia PostgreSQL (solo si no está corriendo)
postgres:
	@echo "🗄️  Iniciando PostgreSQL..."
	@sudo service postgresql status | grep "active (running)" > /dev/null || sudo service postgresql start

## 🔃 Levanta Redis + PostgreSQL
start_db: redis postgres
	@echo "✅ Servicios de base de datos levantados"

## 🚀 Ejecuta la app
run:
	@echo "🚀 Ejecutando la aplicación..."
	go run $(MAIN_FILE)

## 🔃 Levanta DBs y corre la app
start: start_db run

## 🛑 Detiene Redis
stop_redis:
	@echo "🛑 Apagando Redis..."
	@pkill redis-server || echo "Redis ya estaba detenido"

## 🛑 Detiene PostgreSQL
stop_postgres:
	@echo "🛑 Apagando PostgreSQL..."
	@sudo service postgresql stop

## 🛑 Apaga Redis y PostgreSQL
stop_db: stop_redis stop_postgres
	@echo "⛔ Servicios de base de datos detenidos"

test:
	@echo "🧪 Ejecutando los tests..."
	go test ./... -v
