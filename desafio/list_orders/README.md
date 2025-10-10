# Orders App - Desafio List Orders

Aplicação em Go para **criação e listagem de orders** via **REST**, **gRPC** e **GraphQL**, usando **PostgreSQL via Docker**.

---

## 📂 Estrutura do Projeto

list_orders/
├── cmd/
│ └── main.go
├── internal/
│ └── order/
│ ├── grpc.go
│ ├── service.go
│ ├── repository.go
│ ├── handler.go
│ └── graphql.go
├── proto/
│ └── order.proto
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yaml
└── README.md


---

## 🛠 Pré-requisitos

- Docker & Docker Compose
- Go 1.24+ (para rodar localmente)
- Git Bash ou terminal compatível

> **Observação:** não é necessário instalar `psql` no Windows; podemos usar o container do Postgres.

---

## 🚀 Executando localmente (Go)

1. Baixe as dependências:

    ```bash
    go mod tidy
    ```

2. Execute a aplicação:

    ```bash
    go run ./cmd
    ```

3. Portas que a aplicação irá rodar:

    - REST: [http://localhost:8081/order](http://localhost:8081/order)
    - GraphQL: [http://localhost:8082/graphql](http://localhost:8082/graphql)
    - gRPC: `localhost:50051`

---

## 🐳 Executando com Docker

### Build da aplicação

```bash
docker build -t orders-app .
Rodar localmente



docker run -p 8081:8081 -p 8082:8082 -p 50051:50051 orders-app
Docker Compose
Exemplo de docker-compose.yaml:




version: "3.9"
services:
  db:
    image: postgres:15
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: orders
    ports:
      - "5433:5432"
    volumes:
      - db_data:/var/lib/postgresql/data

  app:
    build: .
    depends_on:
      - db
    ports:
      - "8081:8081" # REST
      - "50051:50051" # gRPC
      - "8082:8082" # GraphQL

volumes:
  db_data:
Rodar todos os serviços:




docker compose up --build
Parar:




docker compose down
⚙ Banco de Dados
Criar tabela e inserir dados de teste via Docker:




docker compose exec db psql -U postgres -d orders -c "
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    customer TEXT NOT NULL,
    amount NUMERIC NOT NULL
); 
INSERT INTO orders (customer, amount) VALUES 
('Alice', 99.9), 
('Kleber', 199.9) 
ON CONFLICT DO NOTHING;
"
🌐 Endpoints
REST
Criar order:




curl -X POST http://localhost:8081/order \
  -H "Content-Type: application/json" \
  -d '{"customer":"Alice","amount":99.9}'
Listar orders:




curl http://localhost:8081/order
Exemplo de retorno:




[
  { "id": 2, "customer": "Alice", "amount": 99.9 },
  { "id": 3, "customer": "Kleber", "amount": 199.9 }
]
GraphQL
Query para listar orders:




curl -X POST http://localhost:8082/graphql \
  -H "Content-Type: application/json" \
  -d '{"query":"{ listOrders { id customer amount } }"}'
Exemplo de retorno:




{
  "data": {
    "listOrders": [
      { "id": 2, "customer": "Alice", "amount": 99.9 },
      { "id": 3, "customer": "Kleber", "amount": 199.9 }
    ]
  }
}
gRPC
Listar orders via grpcurl:




grpcurl -plaintext localhost:50051 orderpb.OrderService/ListOrders
Exemplo de retorno:




{
  "orders": [
    { "id": 2, "customer": "Alice", "amount": 99.9 },
    { "id": 3, "customer": "Kleber", "amount": 199.9 }
  ]
}
A aplicação habilita Reflection, então grpcurl ou Evans podem listar os serviços sem passar .proto.

🔧 Configurações importantes
Conexão com Postgres (Go):

Dentro do container:




connStr := "host=db port=5432 user=postgres password=postgres dbname=orders sslmode=disable"
Local Windows:




connStr := "host=localhost port=5433 user=postgres password=postgres dbname=orders sslmode=disable"
📦 Go Modules (go.mod)

