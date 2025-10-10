Rate Limiter em Go

Aplicação em Go que implementa um rate limiter configurável, limitando requisições por IP ou por Token de acesso, com persistência no Redis. Funciona como middleware para servidores HTTP.

🛠 Pré-requisitos

Docker & Docker Compose

Go 1.24+ (para rodar localmente)

Rede local configurada (IP da máquina, ex: 192.168.0.16)

📂 Estrutura do projeto
rate-limiter/
├── Dockerfile
├── docker-compose.yml
├── .env
├── go.mod
├── go.sum
├── main.go
├── test-rate-limiter.sh
└── limiter/
    ├── limiter.go
    └── middleware.go

⚙️ Configuração (.env)
PORT=8080

# Limite por IP
IP_LIMIT=5
IP_BLOCK_SECONDS=10

# Limite por Token
TOKEN_LIMIT=10
TOKEN_BLOCK_SECONDS=10

# Redis
REDIS_ADDR=redis:6379
REDIS_PASSWORD=
REDIS_DB=0

📝 Docker
Dockerfile
# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o rate-limiter ./main.go

# Run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/rate-limiter .
EXPOSE 8080
CMD ["./rate-limiter"]

docker-compose.yml
services:
  redis:
    image: redis:7
    container_name: redis
    ports:
      - "6379:6379"

  app:
    build: .
    container_name: rate-limiter
    ports:
      - "8080:8080"
    env_file:
      - .env
    depends_on:
      - redis

🚀 Executando a aplicação

Build das imagens Docker:

docker-compose build


Subir containers:

docker-compose up -d


Verificar containers:

docker ps


Testar manualmente:

curl -i http://192.168.0.16:8080/test
curl -i -H "API_KEY: abc123" http://192.168.0.16:8080/test

🧪 Teste automatizado
test-rate-limiter.sh
#!/bin/bash
IP="192.168.0.16"
TOKEN="abc123"

echo "Testando limite por IP..."
for i in {1..6}; do
    curl -s -w "%{http_code}\n" http://$IP:8080/test &
done
wait
echo "-------------------------"

echo "Testando limite por Token..."
for i in {1..11}; do
    curl -s -w "%{http_code}\n" -H "API_KEY: $TOKEN" http://$IP:8080/test &
done
wait
echo "-------------------------"


Mostra 200 para requisições dentro do limite.

Mostra 429 quando o limite é ultrapassado.

🔹 Limites e bloqueio

IP_LIMIT: máximo de requisições por segundo por IP.

TOKEN_LIMIT: máximo de requisições por segundo por Token.

IP_BLOCK_SECONDS / TOKEN_BLOCK_SECONDS: tempo que IP/Token fica bloqueado após ultrapassar o limite.

Observação: limites por Token têm prioridade sobre limites por IP.

💡 Observações

Middleware modular: pode ser usado em qualquer servidor HTTP Go.

Persistência via Redis, com possibilidade de substituir por outro mecanismo (Strategy Pattern).

Mensagem de erro padrão quando o limite é atingido:

you have reached the maximum number of requests or actions allowed within a certain time frame


Código HTTP: 429 Too Many Requests

🛑 Parar containers
docker-compose down