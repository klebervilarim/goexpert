.PHONY: run test build docker-up docker-down

run:
	@echo "🚀 Iniciando servidor local na porta 8083..."
	WEATHERAPI_KEY=$${WEATHERAPI_KEY:?set WEATHERAPI_KEY} \
	go run ./cmd/api

test:
	@echo "🧪 Executando testes..."
	go test ./... -v

build:
	@echo "🏗️  Gerando binário local..."
	go build -o bin/api ./cmd/api

docker-up:
	@echo "🐳 Subindo container..."
	docker compose up --build

docker-down:
	@echo "🧹 Derrubando containers..."
	docker compose down
