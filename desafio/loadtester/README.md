# LoadTester CLI - Teste de Carga em Go

Aplicação CLI desenvolvida em **Go** para realizar **testes de carga** em serviços web. Permite definir o número total de requisições e o nível de concorrência (requests simultâneos), gerando um relatório com o resultado da execução.

---

## 🛠 Pré-requisitos

- [Docker](https://www.docker.com/get-started) instalado
- Go 1.24+ (somente se quiser rodar local sem Docker)
- Git Bash ou terminal compatível (Windows/Linux/Mac)

---

## ⚙️ Funcionalidades

- Realiza requisições HTTP para uma URL fornecida.
- Distribui as requisições entre workers concorrentes.
- Garante que o número total de requests seja cumprido.
- Gera relatório ao final do teste contendo:
  - Tempo total de execução
  - Total de requests realizados
  - Requests com status HTTP 200
  - Distribuição de outros códigos de status HTTP
  - Contabiliza erros de requisição

---

## 📁 Estrutura do projeto

loadtester/
├── Dockerfile
├── go.mod
├── go.sum
└── main.go

yaml
Copiar código

---

## 🚀 Executando localmente (sem Docker)

1. Clone o projeto:

```bash
git clone https://github.com/kleber/loadtester.git
cd loadtester
Inicialize módulos Go:

bash
Copiar código
go mod tidy
Compile e rode:

bash
Copiar código
go build -o loadtester main.go
./loadtester --url=http://google.com --requests=1000 --concurrency=10
🐳 Executando via Docker
Build da imagem Docker:

bash
Copiar código
docker build -t loadtester .
Rodar o teste de carga:

bash
Copiar código
docker run --rm loadtester --url=http://google.com --requests=1000 --concurrency=10
✅ No Windows, Linux ou Mac, funciona sem precisar do -- adicional.

🔧 Parâmetros da CLI
Parâmetro	Tipo	Descrição
--url	string	URL do serviço a ser testado (obrigatório)
--requests	int	Número total de requests (padrão: 1)
--concurrency	int	Número de requisições simultâneas (padrão: 1)


docker run --rm loadtester --url=http://google.com --requests=1000 --concurrency=10
Iniciando teste de carga em http://google.com
Total de requests: 1000, Concurrency: 10

===== Relatório =====
Tempo total: 42.855623798s
Total de requests: 1000
Requests com status 200: 1000
Distribuição de outros códigos: