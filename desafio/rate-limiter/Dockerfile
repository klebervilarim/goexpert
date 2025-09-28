# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Copia arquivos de dependência
COPY go.mod go.sum ./
RUN go mod download

# Copia todo o projeto
COPY . .

# Compila
RUN go build -o rate-limiter ./main.go

# Run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/rate-limiter .
EXPOSE 8080
CMD ["./rate-limiter"]
